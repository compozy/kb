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
	t.Setenv(EnvYouTubeYTDLPPath, "")
	t.Setenv(EnvYouTubeProxy, "")
	t.Setenv(EnvYouTubeCookiesFile, "")
	t.Setenv(EnvYouTubeUserAgent, "")
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
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log.level 'info', got %q", cfg.Log.Level)
	}
	if cfg.Vault.Root != "." {
		t.Errorf("expected default vault.root '.', got %q", cfg.Vault.Root)
	}
	if len(cfg.Vault.TopicGlobs) != 1 || cfg.Vault.TopicGlobs[0] != "*" {
		t.Errorf("expected default vault.topic_globs [*], got %#v", cfg.Vault.TopicGlobs)
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
	if cfg.YouTube.YTDLPPath != defaultYouTubeYTDLPPath {
		t.Errorf("expected default youtube.yt_dlp_path %q, got %q", defaultYouTubeYTDLPPath, cfg.YouTube.YTDLPPath)
	}
	if cfg.YouTube.RetryAttempts != defaultYouTubeRetryAttempts {
		t.Errorf("expected default youtube.retry_attempts %d, got %d", defaultYouTubeRetryAttempts, cfg.YouTube.RetryAttempts)
	}
	if cfg.YouTube.RetryBackoff != defaultYouTubeRetryBackoff {
		t.Errorf("expected default youtube.retry_backoff %q, got %q", defaultYouTubeRetryBackoff, cfg.YouTube.RetryBackoff)
	}
}

func TestLoadConfigRoundTrip(t *testing.T) {
	clearServiceEnv(t)

	content := `
[app]
name = "my-service"
env = "production"

[log]
level = "debug"

[vault]
root = "."
topic_globs = ["*", "harness/*"]

[firecrawl]
api_key = "firecrawl-key"
api_url = "https://firecrawl.internal"

[openrouter]
api_key = "openrouter-key"
api_url = "https://openrouter.internal/api"
stt_model = "acme/stt"

	[youtube]
	yt_dlp_path = "/opt/bin/yt-dlp"
	proxy = "http://proxy.internal:8080"
	cookies_file = "/tmp/youtube-cookies.txt"
	user_agent = "kb-test"
retry_attempts = 5
retry_backoff = "250ms"
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
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log.level 'debug', got %q", cfg.Log.Level)
	}
	if cfg.Vault.Root != "." {
		t.Errorf("expected vault.root '.', got %q", cfg.Vault.Root)
	}
	if len(cfg.Vault.TopicGlobs) != 2 || cfg.Vault.TopicGlobs[1] != "harness/*" {
		t.Errorf("expected vault.topic_globs to include harness/*, got %#v", cfg.Vault.TopicGlobs)
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
	if cfg.YouTube.YTDLPPath != "/opt/bin/yt-dlp" {
		t.Errorf("expected youtube.yt_dlp_path, got %q", cfg.YouTube.YTDLPPath)
	}
	if cfg.YouTube.Proxy != "http://proxy.internal:8080" {
		t.Errorf("expected youtube.proxy, got %q", cfg.YouTube.Proxy)
	}
	if cfg.YouTube.CookiesFile != "/tmp/youtube-cookies.txt" {
		t.Errorf("expected youtube.cookies_file, got %q", cfg.YouTube.CookiesFile)
	}
	if cfg.YouTube.UserAgent != "kb-test" {
		t.Errorf("expected youtube.user_agent, got %q", cfg.YouTube.UserAgent)
	}
	if cfg.YouTube.RetryAttempts != 5 {
		t.Errorf("expected youtube.retry_attempts 5, got %d", cfg.YouTube.RetryAttempts)
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
			name:   "invalid log level",
			mutate: func(c *Config) { c.Log.Level = "trace" },
		},
		{
			name:   "empty topic glob",
			mutate: func(c *Config) { c.Vault.TopicGlobs = []string{""} },
		},
		{
			name:   "invalid youtube retry attempts",
			mutate: func(c *Config) { c.YouTube.RetryAttempts = -1 },
		},
		{
			name:   "invalid youtube retry backoff",
			mutate: func(c *Config) { c.YouTube.RetryBackoff = "soon" },
		},
	}

	for _, tc := range testCases {
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
		{
			name:     "youtube yt-dlp path overrides toml",
			envKey:   EnvYouTubeYTDLPPath,
			envValue: "/env/bin/yt-dlp",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.YouTube.YTDLPPath != "/env/bin/yt-dlp" {
					t.Fatalf("expected youtube.yt_dlp_path to be overridden, got %q", cfg.YouTube.YTDLPPath)
				}
			},
		},
		{
			name:     "youtube proxy overrides toml",
			envKey:   EnvYouTubeProxy,
			envValue: "http://env.proxy:8080",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.YouTube.Proxy != "http://env.proxy:8080" {
					t.Fatalf("expected youtube.proxy to be overridden, got %q", cfg.YouTube.Proxy)
				}
			},
		},
		{
			name:     "youtube cookies file overrides toml",
			envKey:   EnvYouTubeCookiesFile,
			envValue: "/tmp/env-cookies.txt",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.YouTube.CookiesFile != "/tmp/env-cookies.txt" {
					t.Fatalf("expected youtube.cookies_file to be overridden, got %q", cfg.YouTube.CookiesFile)
				}
			},
		},
		{
			name:     "youtube user agent overrides toml",
			envKey:   EnvYouTubeUserAgent,
			envValue: "env-agent",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.YouTube.UserAgent != "env-agent" {
					t.Fatalf("expected youtube.user_agent to be overridden, got %q", cfg.YouTube.UserAgent)
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

	[youtube]
	yt_dlp_path = "/toml/bin/yt-dlp"
	proxy = "http://toml.proxy:8080"
	cookies_file = "/tmp/toml-cookies.txt"
	user_agent = "toml-agent"
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

func TestDiscoverProjectConfigPathWalksUp(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	configPath := filepath.Join(root, ProjectConfigFileName)
	if err := os.WriteFile(configPath, []byte("[vault]\nroot = \".\"\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	nested := filepath.Join(root, "harness", "goclaw")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	got, found, err := DiscoverProjectConfigPath(nested)
	if err != nil {
		t.Fatalf("DiscoverProjectConfigPath returned error: %v", err)
	}
	if !found {
		t.Fatal("expected config to be found")
	}
	if got != configPath {
		t.Fatalf("config path = %q, want %q", got, configPath)
	}
}
