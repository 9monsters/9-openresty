// Package release implements the atomic configuration publish coordinator.
// It is the ONLY path through which nginx config files are written, tested, and reloaded.
package release

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"openresty-admin/internal/model"
	"openresty-admin/internal/nginx"
)

// Host is the subset of the nginx Host Adapter the coordinator needs.
// *nginx.Host satisfies this interface; tests inject a fake.
type Host interface {
	ReadForRelease(rel string) (string, bool, error)
	WriteConfig(rel, content string) error
	DeleteConfig(rel string) error
	Test() (bool, string)
	Reload() (bool, string)
}

// Coordinator serialises all configuration changes through a single state machine.
// Only one release can be in flight at a time (process-level mutex).
type Coordinator struct {
	mu   sync.Mutex
	host Host
	db   *gorm.DB
	// probe runs the post-reload health check; overridable in tests.
	probe func([]ReleaseFile) bool
}

// New creates a Coordinator.
func New(host Host, db *gorm.DB) *Coordinator {
	c := &Coordinator{host: host, db: db}
	c.probe = c.healthCheck
	return c
}

// IsLocked reports whether a release is currently in progress.
func (c *Coordinator) IsLocked() bool {
	locked := !c.mu.TryLock()
	if !locked {
		c.mu.Unlock()
	}
	return locked
}

// ReleaseFile is one file change within a release, stored as JSON in ConfigRelease.Files.
type ReleaseFile struct {
	Path    string `json:"path"`
	Action  string `json:"action"` // create, update, delete
	Content string `json:"content,omitempty"`
	SHA256  string `json:"sha256,omitempty"`
}

// Publish executes the full release state machine for a set of changes.
// It snapshots, writes, validates, reloads, health-checks, and rolls back on failure.
func (c *Coordinator) Publish(changes []nginx.Change, source string) (*model.ConfigRelease, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	release := &model.ConfigRelease{
		Status:  "preparing",
		Source:  source,
		Summary: changeSummary(changes),
	}
	c.db.Create(release)

	files := make([]ReleaseFile, len(changes))
	for i, ch := range changes {
		action := "update"
		if ch.Delete {
			action = "delete"
		} else if _, ok, _ := c.host.ReadForRelease(ch.Path); !ok {
			action = "create"
		}
		files[i] = ReleaseFile{
			Path:    ch.Path,
			Action:  action,
			Content: ch.Content,
			SHA256:  sha256hex(ch.Content),
		}
	}
	filesJSON, _ := json.Marshal(files)
	release.Files = string(filesJSON)
	c.db.Save(release)

	slog.Info("release started", "release_id", release.ID, "files", len(files))

	// --- 1. Snapshot current state ---
	snapshots := make(map[string]string)
	for _, f := range files {
		if f.Action == "create" {
			continue // nothing to snapshot
		}
		content, ok, err := c.host.ReadForRelease(f.Path)
		if err != nil {
			return c.restoreAndFail(release, snapshots, "snapshot read failed: "+f.Path+": "+err.Error())
		}
		if ok {
			snapshots[f.Path] = content
		}
	}

	// --- 2. Write ---
	c.setStatus(release, "writing")
	for _, f := range files {
		if f.Action == "delete" {
			if err := c.host.DeleteConfig(f.Path); err != nil {
				return c.restoreAndFail(release, snapshots, "delete failed: "+f.Path+": "+err.Error())
			}
		} else {
			if err := c.host.WriteConfig(f.Path, f.Content); err != nil {
				return c.restoreAndFail(release, snapshots, "write failed: "+f.Path+": "+err.Error())
			}
		}
	}

	// --- 3. Validate (nginx -t) ---
	c.setStatus(release, "validating")
	ok, output := c.host.Test()
	outputStr := output
	release.Output = &outputStr
	c.db.Save(release)

	if !ok {
		restored := c.restoreSnapshots(snapshots)
		if !restored {
			return c.recoveryRequired(release, "nginx -t failed and snapshot restore failed")
		}
		// Re-test the restored config
		if retestOK, _ := c.host.Test(); !retestOK {
			return c.recoveryRequired(release, "nginx -t failed, restore succeeded but re-test still fails")
		}
		return c.finish(release, "validation_failed", "nginx -t: "+output)
	}

	// --- 4. Reload ---
	c.setStatus(release, "reloading")
	reloadOK, reloadOutput := c.host.Reload()
	combinedOutput := output + "\n--- reload ---\n" + reloadOutput
	release.Output = &combinedOutput
	c.db.Save(release)

	if !reloadOK {
		restored := c.restoreSnapshots(snapshots)
		if !restored {
			return c.recoveryRequired(release, "reload failed and snapshot restore failed")
		}
		c.host.Reload() // attempt to reload the restored config
		return c.finish(release, "reload_failed", "reload: "+reloadOutput)
	}

	// --- 5. Health check ---
	c.setStatus(release, "health_checking")
	if !c.probe(files) {
		restored := c.restoreSnapshots(snapshots)
		if !restored {
			return c.recoveryRequired(release, "health check failed and snapshot restore failed")
		}
		if retestOK, _ := c.host.Test(); !retestOK {
			return c.recoveryRequired(release, "health check failed, restore succeeded but re-test fails")
		}
		c.host.Reload() // reload the restored config
		return c.finish(release, "rolled_back", "health check failed, rolled back")
	}

	// --- Success ---
	return c.finish(release, "succeeded", "ok")
}

// restoreAndFail attempts to restore snapshots and marks the release as failed.
func (c *Coordinator) restoreAndFail(release *model.ConfigRelease, snapshots map[string]string, errMsg string) (*model.ConfigRelease, error) {
	c.restoreSnapshots(snapshots)
	if retestOK, _ := c.host.Test(); !retestOK {
		return c.recoveryRequired(release, errMsg+" (restore re-test failed)")
	}
	return c.finish(release, "rolled_back", errMsg)
}

// restoreSnapshots writes snapshot content back to disk. Returns false on failure.
func (c *Coordinator) restoreSnapshots(snapshots map[string]string) bool {
	for path, content := range snapshots {
		if err := c.host.WriteConfig(path, content); err != nil {
			slog.Error("restore snapshot failed", "path", path, "error", err)
			return false
		}
	}
	return true
}

// recoveryRequired marks the release as requiring manual intervention and logs loudly.
func (c *Coordinator) recoveryRequired(release *model.ConfigRelease, reason string) (*model.ConfigRelease, error) {
	slog.Error("RECOVERY_REQUIRED", "release_id", release.ID, "reason", reason)
	errStr := reason
	release.Error = &errStr
	return c.finish(release, "recovery_required", reason)
}

func (c *Coordinator) finish(release *model.ConfigRelease, status, msg string) (*model.ConfigRelease, error) {
	now := time.Now()
	release.Status = status
	release.CompletedAt = &now
	c.db.Save(release)
	slog.Info("release finished", "release_id", release.ID, "status", status)

	if status == "succeeded" {
		return release, nil
	}
	return release, fmt.Errorf("%s", msg)
}

func (c *Coordinator) setStatus(release *model.ConfigRelease, status string) {
	release.Status = status
	c.db.Save(release)
	slog.Info("release status", "release_id", release.ID, "status", status)
}

func changeSummary(changes []nginx.Change) string {
	paths := make([]string, len(changes))
	for i, ch := range changes {
		paths[i] = ch.Path
	}
	return strings.Join(paths, ", ")
}

func sha256hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
