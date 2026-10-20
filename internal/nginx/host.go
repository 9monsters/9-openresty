// Package nginx is the only Host Adapter allowed to touch the OpenResty
// installation: config files, the nginx binary and the log directory.
package nginx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"openresty-admin/internal/config"
)

// Change is one pending config file change. Delete removes the file.
type Change struct {
	Path    string
	Content string
	Delete  bool
}

// ConfigEntry is one entry in the config file listing.
type ConfigEntry struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Type  string `json:"type"`
	Size  int64  `json:"size,omitempty"`
	MTime string `json:"mtime,omitempty"`
}

// LogFile describes a log file on disk.
type LogFile struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	MTime string `json:"mtime"`
}

// Stats are the connection counters reported by the dashboard.
type Stats struct {
	ActiveConnections int64 `json:"activeConnections"`
	Accepts           int64 `json:"accepts"`
	Handled           int64 `json:"handled"`
	Requests          int64 `json:"requests"`
	Reading           int64 `json:"reading"`
	Writing           int64 `json:"writing"`
	Waiting           int64 `json:"waiting"`
}

// Host performs operations against the configured OpenResty installation.
type Host struct {
	mu    sync.RWMutex
	paths config.Paths
}

// NewHost builds a host adapter for the supplied paths.
func NewHost(paths config.Paths) *Host {
	return &Host{paths: paths}
}

// Paths returns the current host paths.
func (h *Host) Paths() config.Paths {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.paths
}

// UpdatePaths replaces the mutable host locations.
func (h *Host) UpdatePaths(paths config.Paths) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if paths.NginxBin != "" {
		h.paths.NginxBin = paths.NginxBin
	}
	if paths.NginxConfDir != "" {
		h.paths.NginxConfDir = paths.NginxConfDir
	}
	if paths.NginxConfPath != "" {
		h.paths.NginxConfPath = paths.NginxConfPath
	}
	if paths.NginxLogDir != "" {
		h.paths.NginxLogDir = paths.NginxLogDir
	}
	if paths.NginxPIDPath != "" {
		h.paths.NginxPIDPath = paths.NginxPIDPath
	}
}

// run executes nginx with argv (never a shell) and reports the exit status.
func (h *Host) run(args []string, sudo bool, timeout time.Duration) (bool, string) {
	paths := h.Paths()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := paths.NginxBin
	argv := args
	if sudo {
		name = "sudo"
		argv = append([]string{"-n", paths.NginxBin}, args...)
	}

	cmd := exec.CommandContext(ctx, name, argv...)
	output, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(output))
	if err != nil {
		if text == "" {
			text = err.Error()
		}
		return false, text
	}
	return true, text
}

// Version returns the nginx/openresty version banner.
func (h *Host) Version() (string, bool) {
	ok, output := h.run([]string{"-v"}, false, 5*time.Second)
	if !ok {
		return "Unknown", false
	}
	return output, strings.Contains(strings.ToLower(output), "openresty")
}

// Test validates the configured nginx configuration.
func (h *Host) Test() (bool, string) {
	paths := h.Paths()
	return h.run([]string{"-t", "-c", paths.NginxConfPath}, false, 30*time.Second)
}

// Reload signals nginx to reload, preferring passwordless sudo.
func (h *Host) Reload() (bool, string) {
	paths := h.Paths()
	if ok, output := h.run([]string{"-s", "reload", "-c", paths.NginxConfPath}, true, 30*time.Second); ok {
		return true, output
	} else if ok, direct := h.run([]string{"-s", "reload", "-c", paths.NginxConfPath}, false, 30*time.Second); ok {
		return true, direct
	} else {
		return false, strings.TrimSpace(output + "\n" + direct)
	}
}

// Start launches the nginx master process. Returns success and output.
func (h *Host) Start() (bool, string) {
	paths := h.Paths()
	return h.run([]string{"-c", paths.NginxConfPath}, false, 10*time.Second)
}

// Stop sends SIGQUIT to the nginx master process (graceful shutdown).
func (h *Host) Stop() (bool, string) {
	paths := h.Paths()
	// Try sudo first, then unprivileged.
	if ok, out := h.run([]string{"-s", "quit", "-c", paths.NginxConfPath}, true, 30*time.Second); ok {
		return true, out
	}
	return h.run([]string{"-s", "quit", "-c", paths.NginxConfPath}, false, 30*time.Second)
}

// Status reports whether the nginx master process is running.
func (h *Host) Status() string {
	paths := h.Paths()
	if raw, err := os.ReadFile(paths.NginxPIDPath); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(raw))); err == nil {
			// Signal 0 checks for the process without delivering a signal.
			if err := syscall.Kill(pid, 0); err == nil {
				return "running"
			}
		}
	}
	if out, err := exec.Command("pgrep", "-f", "nginx: master").Output(); err == nil && strings.TrimSpace(string(out)) != "" {
		return "running"
	}
	return "stopped"
}

// Stats reads stub_status, falling back to nil when it is not configured.
func (h *Host) Stats() *Stats {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1/nginx_status")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	text := string(body)
	if strings.Contains(text, "<html") {
		return nil
	}
	stats := &Stats{
		ActiveConnections: matchInt(text, `Active connections:\s*(\d+)`),
		Reading:           matchInt(text, `Reading:\s*(\d+)`),
		Writing:           matchInt(text, `Writing:\s*(\d+)`),
		Waiting:           matchInt(text, `Waiting:\s*(\d+)`),
	}
	if triple := regexpTriple.FindStringSubmatch(text); triple != nil {
		stats.Accepts = atoi(triple[1])
		stats.Handled = atoi(triple[2])
		stats.Requests = atoi(triple[3])
	}
	return stats
}

// ListConfigFiles walks the config root (or base) and lists config entries.
func (h *Host) ListConfigFiles(base string) []ConfigEntry {
	paths := h.Paths()
	root := paths.NginxConfDir
	if base != "" {
		if resolved, err := h.resolveConfigPath(base); err == nil {
			root = resolved
		} else {
			return []ConfigEntry{}
		}
	}
	if _, err := os.Stat(root); err != nil {
		return []ConfigEntry{}
	}

	entries := make([]ConfigEntry, 0)
	_ = filepath.WalkDir(root, func(path string, dir os.DirEntry, err error) error {
		if err != nil || path == root {
			return nil
		}
		name := dir.Name()
		if strings.HasPrefix(name, ".") || strings.HasSuffix(name, "~") || strings.HasSuffix(name, ".bak") {
			if dir.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if dir.IsDir() {
			entries = append(entries, ConfigEntry{Name: name, Path: rel, Type: "directory"})
			return nil
		}
		info, statErr := dir.Info()
		if statErr != nil {
			return nil
		}
		entries = append(entries, ConfigEntry{
			Name: name, Path: rel, Type: "file",
			Size: info.Size(), MTime: info.ModTime().UTC().Format(time.RFC3339),
		})
		return nil
	})
	return entries
}

// resolveConfigPath rejects escapes, symlinks and directories before any read
// or write. Only relative regular files inside the config root are allowed.
func (h *Host) resolveConfigPath(rel string) (string, error) {
	if rel == "" || filepath.IsAbs(rel) {
		return "", fmt.Errorf("configuration path must be relative")
	}
	for _, segment := range strings.FieldsFunc(rel, func(r rune) bool { return r == '/' || r == '\\' }) {
		if segment == ".." {
			return "", fmt.Errorf("configuration path must stay within the config root")
		}
	}

	root := filepath.Clean(h.Paths().NginxConfDir)
	full := filepath.Clean(filepath.Join(root, rel))
	if full != root && !strings.HasPrefix(full, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("configuration path escapes the config root")
	}

	current := root
	segments := strings.Split(strings.TrimPrefix(full, root+string(os.PathSeparator)), string(os.PathSeparator))
	for i, segment := range segments {
		current = filepath.Join(current, segment)
		info, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symbolic links are not allowed in configuration paths")
		}
		if i == len(segments)-1 {
			if !info.Mode().IsRegular() {
				return "", fmt.Errorf("configuration path must be a regular file")
			}
		} else if !info.IsDir() {
			return "", fmt.Errorf("configuration path is not a directory")
		}
	}
	return full, nil
}

// ReadConfig reads a config file, returning ok=false when it does not exist.
func (h *Host) ReadConfig(rel string) (string, bool) {
	content, ok, err := h.ReadForRelease(rel)
	if err != nil {
		return "", false
	}
	return content, ok
}

// ReadForRelease is ReadConfig with errors kept, for the release coordinator.
func (h *Host) ReadForRelease(rel string) (string, bool, error) {
	full, err := h.resolveConfigPath(rel)
	if err != nil {
		return "", false, err
	}
	raw, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return string(raw), true, nil
}

// WriteConfig writes a config file, snapshotting the previous revision to <file>.bak.
func (h *Host) WriteConfig(rel, content string) error {
	full, err := h.resolveConfigPath(rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(full); err == nil {
		if err := copyFile(full, full+".bak"); err != nil {
			return err
		}
	}
	return os.WriteFile(full, []byte(content), 0o644)
}

// DeleteConfig removes a single config file; directories are refused.
func (h *Host) DeleteConfig(rel string) error {
	full, err := h.resolveConfigPath(rel)
	if err != nil {
		return err
	}
	return os.Remove(full)
}

// ReadLog returns the last n non-empty lines of an access or error log.
func (h *Host) ReadLog(kind string, lines int) string {
	name := "access.log"
	if kind == "error" {
		name = "error.log"
	}
	content, err := h.ReadLogByPath(name, lines)
	if err != nil {
		return ""
	}
	return content
}

// LogFiles lists the *.log files in the log directory.
func (h *Host) LogFiles() []LogFile {
	dir := h.Paths().NginxLogDir
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []LogFile{}
	}
	files := make([]LogFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, LogFile{
			Name: entry.Name(), Path: filepath.Join(dir, entry.Name()),
			Size: info.Size(), MTime: info.ModTime().UTC().Format(time.RFC3339),
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	return files
}

// ReadLogByPath reads a log file relative to the log directory.
func (h *Host) ReadLogByPath(rel string, lines int) (string, error) {
	dir := filepath.Clean(h.Paths().NginxLogDir)
	full := filepath.Clean(filepath.Join(dir, rel))
	if full != dir && !strings.HasPrefix(full, dir+string(os.PathSeparator)) {
		return "", fmt.Errorf("log path escapes the log directory")
	}
	raw, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	return lastLines(string(raw), lines), nil
}

// StubStatusChanges builds the config changes that enable stub_status, so the
// caller can publish them through the release coordinator instead of writing
// files directly.
func (h *Host) StubStatusChanges() ([]Change, error) {
	paths := h.Paths()
	mainRel, err := filepath.Rel(filepath.Clean(paths.NginxConfDir), filepath.Clean(paths.NginxConfPath))
	if err != nil || strings.HasPrefix(mainRel, "..") || filepath.IsAbs(mainRel) {
		return nil, fmt.Errorf("main nginx config is outside the configured root")
	}
	mainRel = filepath.ToSlash(mainRel)

	mainConf, ok, err := h.ReadForRelease(mainRel)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("main nginx configuration does not exist")
	}

	const stubBlock = `# Auto-generated by OpenResty Admin
server {
    listen 127.0.0.1:80;
    server_name localhost;
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
}
`
	changes := make([]Change, 0, 2)
	if !strings.Contains(mainConf, "servers/*.conf") {
		includePattern := regexpHTTPBlock.FindStringIndex(mainConf)
		if includePattern == nil {
			return nil, fmt.Errorf("no http block exists in the main nginx config")
		}
		updated := mainConf[:includePattern[1]] + "\n    include servers/*.conf;" + mainConf[includePattern[1]:]
		changes = append(changes, Change{Path: mainRel, Content: updated})
	}
	changes = append(changes, Change{Path: "servers/status.conf", Content: stubBlock})
	return changes, nil
}

func copyFile(src, dst string) error {
	raw, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, raw, 0o644)
}

func lastLines(content string, lines int) string {
	if lines <= 0 {
		lines = 200
	}
	all := strings.Split(content, "\n")
	kept := make([]string, 0, lines)
	for _, line := range all {
		if strings.TrimSpace(line) != "" {
			kept = append(kept, line)
		}
	}
	if len(kept) > lines {
		kept = kept[len(kept)-lines:]
	}
	return strings.Join(kept, "\n")
}

func atoi(value string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}
