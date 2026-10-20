package nginx

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BackupEntry is one historical config snapshot.
type BackupEntry struct {
	ID        string `json:"id"`
	Filepath  string `json:"filepath"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Content   string `json:"content,omitempty"`
}

func (h *Host) backupDir() string {
	return filepath.Join(h.Paths().NginxConfDir, ".backups")
}

// BackupConfig stores a durable copy of a config file before it is modified.
func (h *Host) BackupConfig(rel, content, message string) (string, error) {
	dir := h.backupDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s%06d", strconvBase36(time.Now().UnixNano()), os.Getpid()%1000000)
	safeName := strings.NewReplacer("/", "_", "\\", "_").Replace(rel)
	backupPath := filepath.Join(dir, id+"_"+safeName)
	if err := os.WriteFile(backupPath, []byte(content), 0o600); err != nil {
		return "", err
	}
	meta, err := json.Marshal(BackupEntry{
		ID: id, Filepath: rel, Timestamp: time.Now().UTC().Format(time.RFC3339Nano), Message: message,
	})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(backupPath+".meta", meta, 0o600); err != nil {
		return "", err
	}
	return id, nil
}

// ListBackups returns the stored backups for a config path, newest first.
func (h *Host) ListBackups(rel string) []BackupEntry {
	dir := h.backupDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []BackupEntry{}
	}
	safeName := strings.NewReplacer("/", "_", "\\", "_").Replace(rel)
	backups := make([]BackupEntry, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), safeName) {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, entry.Name()+".meta"))
		if err != nil {
			continue
		}
		var backup BackupEntry
		if err := json.Unmarshal(raw, &backup); err != nil {
			continue
		}
		if content, err := os.ReadFile(filepath.Join(dir, entry.Name())); err == nil {
			backup.Content = string(content)
		}
		backups = append(backups, backup)
	}
	sort.Slice(backups, func(i, j int) bool { return backups[i].Timestamp > backups[j].Timestamp })
	return backups
}

func strconvBase36(value int64) string {
	const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
	if value == 0 {
		return "0"
	}
	out := make([]byte, 0, 16)
	for value > 0 {
		out = append([]byte{digits[value%36]}, out...)
		value /= 36
	}
	return string(out)
}
