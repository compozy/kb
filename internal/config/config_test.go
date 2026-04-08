package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigHasValidDefaults(t *testing.T) {
	t.Parallel()

	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
	if cfg.App.Name != "app" {
		t.Errorf("expected default app.name 'app', got %q", cfg.App.Name)
	}
	if cfg.App.Env != "development" {
		t.Errorf("expected default app.env 'development', got %q", cfg.App.Env)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default server.host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default server.port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log.level 'info', got %q", cfg.Log.Level)
	}
}

func TestLoadConfigRoundTrip(t *testing.T) {
	t.Parallel()

	content := `
[app]
name = "my-service"
env = "production"

[server]
host = "127.0.0.1"
port = 3000

[log]
level = "debug"
`
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.App.Name != "my-service" {
		t.Errorf("expected app.name 'my-service', got %q", cfg.App.Name)
	}
	if cfg.App.Env != "production" {
		t.Errorf("expected app.env 'production', got %q", cfg.App.Env)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected server.host '127.0.0.1', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("expected server.port 3000, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log.level 'debug', got %q", cfg.Log.Level)
	}
}

func TestLoadEmptyPathUsesDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load with empty path: %v", err)
	}
	if cfg.App.Name != "app" {
		t.Errorf("expected default app.name 'app', got %q", cfg.App.Name)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	content := `
[app]
name = "test"
env = "development"
unknown_field = true
`
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unknown keys, got nil")
	}
}

func TestValidateRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		mutate func(*Config)
	}{
		{
			name:   "empty app name",
			mutate: func(c *Config) { c.App.Name = "" },
		},
		{
			name:   "whitespace app name",
			mutate: func(c *Config) { c.App.Name = "   " },
		},
		{
			name:   "invalid app env",
			mutate: func(c *Config) { c.App.Env = "local" },
		},
		{
			name:   "port zero",
			mutate: func(c *Config) { c.Server.Port = 0 },
		},
		{
			name:   "port negative",
			mutate: func(c *Config) { c.Server.Port = -1 },
		},
		{
			name:   "port too high",
			mutate: func(c *Config) { c.Server.Port = 70000 },
		},
		{
			name:   "invalid log level",
			mutate: func(c *Config) { c.Log.Level = "trace" },
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := Default()
			tc.mutate(&cfg)
			if err := cfg.Validate(); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestLoadDotEnvIfPresentLoadsValues(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("TEST_DOTENV_VAR=hello\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	if err := LoadDotEnvIfPresent(envPath); err != nil {
		t.Fatalf("load dotenv: %v", err)
	}
	if got := os.Getenv("TEST_DOTENV_VAR"); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestLoadDotEnvIfPresentMissingFileIsOK(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), ".env")
	if err := LoadDotEnvIfPresent(path); err != nil {
		t.Fatalf("missing .env should not error: %v", err)
	}
}
