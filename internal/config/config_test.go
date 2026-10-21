package config

import (
	"os"
	"testing"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"HOST", "PORT", "DATABASE_URL", "DATABASE_DRIVER", "SQLITE_PATH"} {
		os.Unsetenv(k)
	}
}

func TestResolveDefaults(t *testing.T) {
	clearEnv(t)
	cfg := (Merge{}).Resolve()
	if cfg.Host != "127.0.0.1" {
		t.Errorf("default host = %q, want 127.0.0.1", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("default port = %d, want 3000", cfg.Port)
	}
	if cfg.DatabaseDriver != "postgres" {
		t.Errorf("default driver = %q, want postgres", cfg.DatabaseDriver)
	}
	if cfg.DatabaseDSN == "" {
		t.Error("expected a non-empty default DSN")
	}
}

func TestResolveEnvOverrides(t *testing.T) {
	clearEnv(t)
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("PORT", "8080")
	cfg := (Merge{}).Resolve()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("host = %q, want 0.0.0.0", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("port = %d, want 8080", cfg.Port)
	}
}

func TestResolveCLIOverridesEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("PORT", "8080")
	cfg := (Merge{Host: "10.0.0.9", Port: 9090}).Resolve()
	if cfg.Host != "10.0.0.9" {
		t.Errorf("CLI host should win: got %q", cfg.Host)
	}
	if cfg.Port != 9090 {
		t.Errorf("CLI port should win: got %d", cfg.Port)
	}
}

func TestResolveSQLiteDriver(t *testing.T) {
	clearEnv(t)
	t.Setenv("DATABASE_DRIVER", "sqlite")
	t.Setenv("SQLITE_PATH", "/tmp/openresty-test.db")
	cfg := (Merge{}).Resolve()
	if cfg.DatabaseDriver != "sqlite" {
		t.Errorf("driver = %q, want sqlite", cfg.DatabaseDriver)
	}
	if cfg.DatabaseDSN != "/tmp/openresty-test.db" {
		t.Errorf("dsn = %q, want sqlite path", cfg.DatabaseDSN)
	}
}

func TestResolveDatabaseURLWinsOverDriver(t *testing.T) {
	clearEnv(t)
	t.Setenv("DATABASE_DRIVER", "sqlite")
	t.Setenv("DATABASE_URL", "postgres://u:p@h/db")
	cfg := (Merge{}).Resolve()
	if cfg.DatabaseDSN != "postgres://u:p@h/db" {
		t.Errorf("DATABASE_URL should win: got %q", cfg.DatabaseDSN)
	}
}
