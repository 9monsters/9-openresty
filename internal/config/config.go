// Package config holds the runtime configuration shared by the CLI and the API.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Paths are the mutable OpenResty host locations. They default to the local
// brew layout but are always overridable, never hard-coded into logic.
type Paths struct {
	NginxBin      string `json:"nginxBin"`
	NginxConfDir  string `json:"nginxConfDir"`
	NginxConfPath string `json:"nginxConfPath"`
	NginxLogDir   string `json:"nginxLogDir"`
	NginxPIDPath  string `json:"nginxPIDPath"`
}

// DefaultPaths returns the local single-machine brew defaults.
func DefaultPaths() Paths {
	return Paths{
		NginxBin:      "/usr/local/opt/openresty/bin/openresty",
		NginxConfDir:  "/usr/local/etc/openresty",
		NginxConfPath: "/usr/local/etc/openresty/nginx.conf",
		NginxLogDir:   "/usr/local/var/log/nginx",
		NginxPIDPath:  "/usr/local/var/run/openresty.pid",
	}
}

// FileConfig is the on-disk persistent configuration format.
type FileConfig struct {
	Host  *string `json:"host,omitempty"`
	Port  *int    `json:"port,omitempty"`
	Paths *Paths  `json:"paths,omitempty"`
}

// Config is the resolved runtime configuration.
type Config struct {
	Host          string
	Port          int
	DatabaseDSN   string
	DatabaseDriver string
	Paths         Paths
}

// Addr returns the host:port the server binds to.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadFile reads a JSON config file and returns the parsed values.
func LoadFile(path string) (FileConfig, error) {
	var fc FileConfig
	if path == "" {
		return fc, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fc, nil
		}
		return fc, fmt.Errorf("read config file %s: %w", path, err)
	}
	if err := json.Unmarshal(raw, &fc); err != nil {
		return fc, fmt.Errorf("parse config file %s: %w", path, err)
	}
	return fc, nil
}

// Merge resolves the effective Config from a layered priority:
//  1. CLI flags (explicit, non-zero)
//  2. Config file values
//  3. Environment variables
//  4. Hard-coded defaults
type Merge struct {
	Host  string
	Port  int
	Paths Paths
	File  FileConfig
}

// Resolve produces the final Config.
func (m Merge) Resolve() Config {
	paths := DefaultPaths()
	if m.File.Paths != nil {
		applyPathOverride(&paths.NginxBin, m.File.Paths.NginxBin)
		applyPathOverride(&paths.NginxConfDir, m.File.Paths.NginxConfDir)
		applyPathOverride(&paths.NginxConfPath, m.File.Paths.NginxConfPath)
		applyPathOverride(&paths.NginxLogDir, m.File.Paths.NginxLogDir)
		applyPathOverride(&paths.NginxPIDPath, m.File.Paths.NginxPIDPath)
	}
	// CLI flags override file.
	applyPathOverride(&paths.NginxBin, m.Paths.NginxBin)
	applyPathOverride(&paths.NginxConfDir, m.Paths.NginxConfDir)
	applyPathOverride(&paths.NginxConfPath, m.Paths.NginxConfPath)
	applyPathOverride(&paths.NginxLogDir, m.Paths.NginxLogDir)
	applyPathOverride(&paths.NginxPIDPath, m.Paths.NginxPIDPath)

	cfg := Config{Paths: paths}

	// Host: CLI > file > env > default
	cfg.Host = firstNonEmpty(
		m.Host,
		ptrOrEmpty(m.File.Host),
		envOr("HOST", ""),
		"127.0.0.1",
	)

	// Port: CLI > file > env > default
	cfg.Port = firstNonZero(m.Port,
		ptrOrZero(m.File.Port),
		envInt("PORT", 3000),
	)

	// Database: driver defaults to postgres; DATABASE_URL env var wins for DSN.
	cfg.DatabaseDriver = envOr("DATABASE_DRIVER", "postgres")
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		cfg.DatabaseDSN = dsn
	} else if cfg.DatabaseDriver == "sqlite" {
		cfg.DatabaseDSN = envOr("SQLITE_PATH", "data/openresty-admin.db")
	} else {
		cfg.DatabaseDSN = "postgres://openresty:1aac%40Openresty@192.168.5.100:5432/openresty?sslmode=disable"
	}

	return cfg
}

func applyPathOverride(target *string, override string) {
	if override != "" {
		*target = override
	}
}

// --- helpers ---

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonZero(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

func ptrOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func ptrOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if p, err := strconv.Atoi(raw); err == nil {
			return p
		}
	}
	return fallback
}