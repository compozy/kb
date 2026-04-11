package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return path
}

func clearServiceEnv(t *testing.T) {
	t.Helper()

	t.Setenv(EnvFirecrawlAPIKey, "")
	t.Setenv(EnvFirecrawlAPIURL, "")
	t.Setenv(EnvOpenRouterAPIKey, "")
	t.Setenv(EnvOpenRouterAPIURL, "")
}

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
	if cfg.Firecrawl.APIURL != defaultFirecrawlAPIURL {
		t.Errorf("expected default firecrawl.api_url %q, got %q", defaultFirecrawlAPIURL, cfg.Firecrawl.APIURL)
	}
	if cfg.OpenRouter.APIURL != defaultOpenRouterAPIURL {
		t.Errorf("expected default openrouter.api_url %q, got %q", defaultOpenRouterAPIURL, cfg.OpenRouter.APIURL)
	}
	if cfg.OpenRouter.STTModel != defaultOpenRouterSTTModel {
		t.Errorf("expected default openrouter.stt_model %q, got %q", defaultOpenRouterSTTModel, cfg.OpenRouter.STTModel)
	}
}

func TestLoadConfigRoundTrip(t *testing.T) {
	clearServiceEnv(t)

	content := `
[app]
name = "my-service"
env = "production"

[server]
host = "127.0.0.1"
port = 3000

[log]
level = "debug"

[firecrawl]
api_key = "firecrawl-key"
api_url = "https://firecrawl.internal"

[openrouter]
api_key = "openrouter-key"
api_url = "https://openrouter.internal/api"
stt_model = "acme/stt"
`
	path := writeConfigFile(t, content)

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
	if cfg.Firecrawl.APIKey != "firecrawl-key" {
		t.Errorf("expected firecrawl.api_key 'firecrawl-key', got %q", cfg.Firecrawl.APIKey)
	}
	if cfg.Firecrawl.APIURL != "https://firecrawl.internal" {
		t.Errorf("expected firecrawl.api_url 'https://firecrawl.internal', got %q", cfg.Firecrawl.APIURL)
	}
	if cfg.OpenRouter.APIKey != "openrouter-key" {
		t.Errorf("expected openrouter.api_key 'openrouter-key', got %q", cfg.OpenRouter.APIKey)
	}
	if cfg.OpenRouter.APIURL != "https://openrouter.internal/api" {
		t.Errorf("expected openrouter.api_url 'https://openrouter.internal/api', got %q", cfg.OpenRouter.APIURL)
	}
	if cfg.OpenRouter.STTModel != "acme/stt" {
		t.Errorf("expected openrouter.stt_model 'acme/stt', got %q", cfg.OpenRouter.STTModel)
	}
}

func TestLoadEmptyPathUsesDefaults(t *testing.T) {
	clearServiceEnv(t)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load with empty path: %v", err)
	}
	if cfg.App.Name != "app" {
		t.Errorf("expected default app.name 'app', got %q", cfg.App.Name)
	}
	if cfg.Firecrawl.APIURL != defaultFirecrawlAPIURL {
		t.Errorf("expected default firecrawl.api_url %q, got %q", defaultFirecrawlAPIURL, cfg.Firecrawl.APIURL)
	}
	if cfg.OpenRouter.APIURL != defaultOpenRouterAPIURL {
		t.Errorf("expected default openrouter.api_url %q, got %q", defaultOpenRouterAPIURL, cfg.OpenRouter.APIURL)
	}
	if cfg.OpenRouter.STTModel != defaultOpenRouterSTTModel {
		t.Errorf("expected default openrouter.stt_model %q, got %q", defaultOpenRouterSTTModel, cfg.OpenRouter.STTModel)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	clearServiceEnv(t)

	content := `
[app]
name = "test"
env = "development"
unknown_field = true
`
	path := writeConfigFile(t, content)

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

func TestLoadUsesFirecrawlDefaultsWhenSectionMissing(t *testing.T) {
	clearServiceEnv(t)

	path := writeConfigFile(t, `
[app]
name = "kb"
env = "development"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Firecrawl.APIURL != defaultFirecrawlAPIURL {
		t.Errorf("expected default firecrawl.api_url %q, got %q", defaultFirecrawlAPIURL, cfg.Firecrawl.APIURL)
	}
}

func TestLoadUsesOpenRouterDefaultsWhenSectionMissing(t *testing.T) {
	clearServiceEnv(t)

	path := writeConfigFile(t, `
[app]
name = "kb"
env = "development"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.OpenRouter.APIURL != defaultOpenRouterAPIURL {
		t.Errorf("expected default openrouter.api_url %q, got %q", defaultOpenRouterAPIURL, cfg.OpenRouter.APIURL)
	}
	if cfg.OpenRouter.STTModel != defaultOpenRouterSTTModel {
		t.Errorf("expected default openrouter.stt_model %q, got %q", defaultOpenRouterSTTModel, cfg.OpenRouter.STTModel)
	}
}

func TestLoadEnvOverridesServiceConfig(t *testing.T) {
	testCases := []struct {
		name     string
		envKey   string
		envValue string
		assert   func(*testing.T, Config)
	}{
		{
			name:     "firecrawl api key overrides toml",
			envKey:   EnvFirecrawlAPIKey,
			envValue: "env-firecrawl-key",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.Firecrawl.APIKey != "env-firecrawl-key" {
					t.Fatalf("expected firecrawl.api_key to be overridden, got %q", cfg.Firecrawl.APIKey)
				}
			},
		},
		{
			name:     "openrouter api key overrides toml",
			envKey:   EnvOpenRouterAPIKey,
			envValue: "env-openrouter-key",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.OpenRouter.APIKey != "env-openrouter-key" {
					t.Fatalf("expected openrouter.api_key to be overridden, got %q", cfg.OpenRouter.APIKey)
				}
			},
		},
		{
			name:     "firecrawl api url overrides toml",
			envKey:   EnvFirecrawlAPIURL,
			envValue: "https://env.firecrawl.dev",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.Firecrawl.APIURL != "https://env.firecrawl.dev" {
					t.Fatalf("expected firecrawl.api_url to be overridden, got %q", cfg.Firecrawl.APIURL)
				}
			},
		},
		{
			name:     "openrouter api url overrides toml",
			envKey:   EnvOpenRouterAPIURL,
			envValue: "https://env.openrouter.ai/api",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.OpenRouter.APIURL != "https://env.openrouter.ai/api" {
					t.Fatalf("expected openrouter.api_url to be overridden, got %q", cfg.OpenRouter.APIURL)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearServiceEnv(t)

			path := writeConfigFile(t, `
[app]
name = "kb"
env = "development"

[firecrawl]
api_key = "toml-firecrawl-key"
api_url = "https://toml.firecrawl.dev"

[openrouter]
api_key = "toml-openrouter-key"
api_url = "https://toml.openrouter.ai/api"
stt_model = "toml/stt"
`)

			t.Setenv(tc.envKey, tc.envValue)

			cfg, err := Load(path)
			if err != nil {
				t.Fatalf("load config: %v", err)
			}

			tc.assert(t, cfg)
		})
	}
}
