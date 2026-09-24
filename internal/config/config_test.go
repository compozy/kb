package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
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
	t.Setenv(EnvOpenAIAPIKey, "")
	t.Setenv(EnvOpenAIAPIURL, "")
	t.Setenv(EnvSTTProvider, "")
	t.Setenv(EnvSTTModel, "")
	t.Setenv(EnvYouTubeYTDLPPath, "")
	t.Setenv(EnvYouTubeProxy, "")
	t.Setenv(EnvYouTubeCookiesFile, "")
	t.Setenv(EnvYouTubeUserAgent, "")
	t.Setenv(EnvYouTubeCaptionLanguages, "")
	t.Setenv(EnvDecisionsModel, "")
	t.Setenv(EnvGenerationModel, "")
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
	t.Run("Should default OKF types to empty", func(t *testing.T) {
		if len(cfg.OKF.Types) != 0 {
			t.Errorf("expected default okf.types to be empty, got %#v", cfg.OKF.Types)
		}
	})
	if cfg.Firecrawl.APIURL != defaultFirecrawlAPIURL {
		t.Errorf("expected default firecrawl.api_url %q, got %q", defaultFirecrawlAPIURL, cfg.Firecrawl.APIURL)
	}
	if cfg.OpenRouter.APIURL != defaultOpenRouterAPIURL {
		t.Errorf("expected default openrouter.api_url %q, got %q", defaultOpenRouterAPIURL, cfg.OpenRouter.APIURL)
	}
	if cfg.OpenRouter.STTModel != defaultOpenRouterSTTModel {
		t.Errorf("expected default openrouter.stt_model %q, got %q", defaultOpenRouterSTTModel, cfg.OpenRouter.STTModel)
	}
	if cfg.STT.Provider != defaultSTTProvider {
		t.Errorf("expected default stt.provider %q, got %q", defaultSTTProvider, cfg.STT.Provider)
	}
	if cfg.STT.APIURL != defaultOpenAIAPIURL {
		t.Errorf("expected default stt.api_url %q, got %q", defaultOpenAIAPIURL, cfg.STT.APIURL)
	}
	if cfg.STT.Model != defaultSTTModel {
		t.Errorf("expected default stt.model %q, got %q", defaultSTTModel, cfg.STT.Model)
	}
	if cfg.STT.AudioFormat != defaultSTTAudioFormat {
		t.Errorf("expected default stt.audio_format %q, got %q", defaultSTTAudioFormat, cfg.STT.AudioFormat)
	}
	if cfg.STT.Concurrency != defaultSTTConcurrency {
		t.Errorf("expected default stt.concurrency %d, got %d", defaultSTTConcurrency, cfg.STT.Concurrency)
	}
	if cfg.YouTube.YTDLPPath != defaultYouTubeYTDLPPath {
		t.Errorf("expected default youtube.yt_dlp_path %q, got %q", defaultYouTubeYTDLPPath, cfg.YouTube.YTDLPPath)
	}
	if cfg.YouTube.Transcription != defaultYouTubeTranscription {
		t.Errorf("expected default youtube.transcription %q, got %q", defaultYouTubeTranscription, cfg.YouTube.Transcription)
	}
	if cfg.YouTube.RetryAttempts != defaultYouTubeRetryAttempts {
		t.Errorf("expected default youtube.retry_attempts %d, got %d", defaultYouTubeRetryAttempts, cfg.YouTube.RetryAttempts)
	}
	if cfg.YouTube.RetryBackoff != defaultYouTubeRetryBackoff {
		t.Errorf("expected default youtube.retry_backoff %q, got %q", defaultYouTubeRetryBackoff, cfg.YouTube.RetryBackoff)
	}
	if !reflect.DeepEqual(cfg.YouTube.CaptionLanguages, []string{"orig"}) {
		t.Errorf("expected default youtube.caption_languages [orig], got %#v", cfg.YouTube.CaptionLanguages)
	}
	if cfg.YouTube.AllowTranslatedCaptions {
		t.Error("expected youtube.allow_translated_captions to default false")
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

[okf]
types = ["Playbook", " Reference ", ""]

[firecrawl]
api_key = "firecrawl-key"
api_url = "https://firecrawl.internal"

[openrouter]
api_key = "openrouter-key"
api_url = "https://openrouter.internal/api"
stt_model = "acme/stt"

[stt]
provider = "openai"
api_key = "openai-key"
api_url = "https://openai.internal"
model = "gpt-4o-mini-transcribe"
language = "pt"
prompt = "Technical conference talk."
audio_format = "mp3"
chunk_duration = "5m"
max_chunk_bytes = 123456
concurrency = 3
ffmpeg_path = "/opt/bin/ffmpeg"

	[youtube]
	yt_dlp_path = "/opt/bin/yt-dlp"
	proxy = "http://proxy.internal:8080"
	cookies_file = "/tmp/youtube-cookies.txt"
	user_agent = "kb-test"
transcription = "auto"
retry_attempts = 5
retry_backoff = "250ms"
caption_languages = ["orig", " pt "]
allow_translated_captions = true
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
	t.Run("Should normalize OKF type vocabulary", func(t *testing.T) {
		if !reflect.DeepEqual(cfg.OKF.Types, []string{"Playbook", "Reference"}) {
			t.Errorf("expected normalized okf.types, got %#v", cfg.OKF.Types)
		}
	})
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
	if cfg.STT.APIKey != "openai-key" {
		t.Errorf("expected stt.api_key 'openai-key', got %q", cfg.STT.APIKey)
	}
	if cfg.STT.APIURL != "https://openai.internal" {
		t.Errorf("expected stt.api_url 'https://openai.internal', got %q", cfg.STT.APIURL)
	}
	if cfg.STT.Model != "gpt-4o-mini-transcribe" {
		t.Errorf("expected stt.model 'gpt-4o-mini-transcribe', got %q", cfg.STT.Model)
	}
	if cfg.STT.Language != "pt" {
		t.Errorf("expected stt.language 'pt', got %q", cfg.STT.Language)
	}
	if cfg.STT.Prompt != "Technical conference talk." {
		t.Errorf("expected stt.prompt, got %q", cfg.STT.Prompt)
	}
	if cfg.STT.ChunkDuration != "5m" || cfg.STT.MaxChunkBytes != 123456 || cfg.STT.Concurrency != 3 {
		t.Errorf("unexpected stt chunk settings: %#v", cfg.STT)
	}
	if cfg.YouTube.YTDLPPath != "/opt/bin/yt-dlp" {
		t.Errorf("expected youtube.yt_dlp_path, got %q", cfg.YouTube.YTDLPPath)
	}
	if cfg.YouTube.Transcription != "auto" {
		t.Errorf("expected youtube.transcription auto, got %q", cfg.YouTube.Transcription)
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
	if !reflect.DeepEqual(cfg.YouTube.CaptionLanguages, []string{"orig", "pt"}) {
		t.Errorf("expected youtube.caption_languages, got %#v", cfg.YouTube.CaptionLanguages)
	}
	if !cfg.YouTube.AllowTranslatedCaptions {
		t.Error("expected youtube.allow_translated_captions true")
	}
}

func TestLoadDecisionsAndGenerationSections(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		assert  func(*testing.T, Config)
	}{
		{
			name:    "Should default decisions, generation and refetch settings when sections are missing",
			content: "[app]\nname = \"kb\"\nenv = \"development\"\n",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				want := DecisionsConfig{
					Model:         "typesafe/jev-1.13",
					Deadline:      "20s",
					Retries:       2,
					Concurrency:   8,
					MaxStateBytes: 98304,
					BudgetUSD:     1.0,
					Mode:          DecisionModeShadow,
				}
				if !reflect.DeepEqual(cfg.Decisions, want) {
					t.Fatalf("decisions = %#v, want %#v", cfg.Decisions, want)
				}
				wantGen := GenerationConfig{
					Model:           "xiaomi/mimo-v2.6-flash",
					FallbackModel:   "deepseek/deepseek-v4-flash",
					Deadline:        "90s",
					Retries:         1,
					SummaryLanguage: "en",
				}
				if cfg.Generation != wantGen {
					t.Fatalf("generation = %#v, want %#v", cfg.Generation, wantGen)
				}
				if cfg.Firecrawl.RefetchWaitMS != 3000 || cfg.Firecrawl.RefetchOnlyMainContent {
					t.Fatalf("unexpected firecrawl refetch defaults: %#v", cfg.Firecrawl)
				}
				deadline, err := cfg.Decisions.DeadlineDuration()
				if err != nil || deadline != 20*time.Second {
					t.Fatalf("decisions deadline = %v, %v", deadline, err)
				}
				genDeadline, err := cfg.Generation.DeadlineDuration()
				if err != nil || genDeadline != 90*time.Second {
					t.Fatalf("generation deadline = %v, %v", genDeadline, err)
				}
			},
		},
		{
			name: "Should load every decisions, generation and refetch key from TOML",
			content: `
[firecrawl]
refetch_wait_ms = 5000
refetch_only_main_content = true

[decisions]
model = "typesafe/jev-1.14"
deadline = "15s"
retries = 0
concurrency = 4
max_state_bytes = 65536
budget_usd = 2.5
mode = "APPLY"

[decisions.thresholds]
link_apply = 0.9
relevance_quarantine = 0.75

[generation]
model = "acme/gen"
fallback_model = ""
deadline = "30s"
retries = 2
summary_language = "pt"
`,
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				want := DecisionsConfig{
					Model:         "typesafe/jev-1.14",
					Deadline:      "15s",
					Retries:       0,
					Concurrency:   4,
					MaxStateBytes: 65536,
					BudgetUSD:     2.5,
					Mode:          DecisionModeApply,
					Thresholds:    map[string]float64{"link_apply": 0.9, "relevance_quarantine": 0.75},
				}
				if !reflect.DeepEqual(cfg.Decisions, want) {
					t.Fatalf("decisions = %#v, want %#v", cfg.Decisions, want)
				}
				wantGen := GenerationConfig{Model: "acme/gen", Deadline: "30s", Retries: 2, SummaryLanguage: "pt"}
				if cfg.Generation != wantGen {
					t.Fatalf("generation = %#v, want %#v", cfg.Generation, wantGen)
				}
				if cfg.Firecrawl.RefetchWaitMS != 5000 || !cfg.Firecrawl.RefetchOnlyMainContent {
					t.Fatalf("unexpected firecrawl refetch settings: %#v", cfg.Firecrawl)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearServiceEnv(t)
			cfg, err := Load(writeConfigFile(t, tc.content))
			if err != nil {
				t.Fatalf("load config: %v", err)
			}
			tc.assert(t, cfg)
		})
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
	if cfg.STT.Model != defaultSTTModel {
		t.Errorf("expected default stt.model %q, got %q", defaultSTTModel, cfg.STT.Model)
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
		name            string
		mutate          func(*Config)
		wantErrContains string
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
		{
			name:   "invalid stt provider",
			mutate: func(c *Config) { c.STT.Provider = "local" },
		},
		{
			name:   "invalid stt audio format",
			mutate: func(c *Config) { c.STT.AudioFormat = "flac" },
		},
		{
			name:   "invalid stt concurrency",
			mutate: func(c *Config) { c.STT.Concurrency = -1 },
		},
		{
			name:   "invalid youtube transcription",
			mutate: func(c *Config) { c.YouTube.Transcription = "maybe" },
		},
		{
			name:            "empty youtube caption languages",
			mutate:          func(c *Config) { c.YouTube.CaptionLanguages = []string{" "} },
			wantErrContains: "youtube.caption_languages",
		},
		{
			name:            "negative firecrawl refetch wait",
			mutate:          func(c *Config) { c.Firecrawl.RefetchWaitMS = -1 },
			wantErrContains: "firecrawl.refetch_wait_ms",
		},
		{
			name:            "empty decisions model",
			mutate:          func(c *Config) { c.Decisions.Model = " " },
			wantErrContains: "decisions.model",
		},
		{
			name:            "invalid decisions deadline",
			mutate:          func(c *Config) { c.Decisions.Deadline = "soon" },
			wantErrContains: "decisions.deadline",
		},
		{
			name:            "non-positive decisions deadline",
			mutate:          func(c *Config) { c.Decisions.Deadline = "0s" },
			wantErrContains: "decisions.deadline",
		},
		{
			name:            "negative decisions retries",
			mutate:          func(c *Config) { c.Decisions.Retries = -1 },
			wantErrContains: "decisions.retries",
		},
		{
			name:            "zero decisions concurrency",
			mutate:          func(c *Config) { c.Decisions.Concurrency = 0 },
			wantErrContains: "decisions.concurrency",
		},
		{
			name:            "tiny decisions max state bytes",
			mutate:          func(c *Config) { c.Decisions.MaxStateBytes = 10 },
			wantErrContains: "decisions.max_state_bytes",
		},
		{
			name:            "zero decisions budget",
			mutate:          func(c *Config) { c.Decisions.BudgetUSD = 0 },
			wantErrContains: "decisions.budget_usd",
		},
		{
			name:            "invalid decisions mode",
			mutate:          func(c *Config) { c.Decisions.Mode = "live" },
			wantErrContains: "decisions.mode",
		},
		{
			name:            "unknown decisions threshold",
			mutate:          func(c *Config) { c.Decisions.Thresholds = map[string]float64{"link_aply": 0.9} },
			wantErrContains: "decisions.thresholds.link_aply is not a known threshold",
		},
		{
			name:            "out of range decisions threshold",
			mutate:          func(c *Config) { c.Decisions.Thresholds = map[string]float64{"link_apply": 1.5} },
			wantErrContains: "decisions.thresholds.link_apply must be between 0 and 1",
		},
		{
			name:            "empty generation model",
			mutate:          func(c *Config) { c.Generation.Model = "" },
			wantErrContains: "generation.model",
		},
		{
			name:            "invalid generation deadline",
			mutate:          func(c *Config) { c.Generation.Deadline = "-5s" },
			wantErrContains: "generation.deadline",
		},
		{
			name:            "negative generation retries",
			mutate:          func(c *Config) { c.Generation.Retries = -2 },
			wantErrContains: "generation.retries",
		},
		{
			name:            "empty generation summary language",
			mutate:          func(c *Config) { c.Generation.SummaryLanguage = "" },
			wantErrContains: "generation.summary_language",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := Default()
			tc.mutate(&cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if tc.wantErrContains != "" && !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Fatalf("error = %q, want %q", err.Error(), tc.wantErrContains)
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
			name:     "openai api key overrides stt config",
			envKey:   EnvOpenAIAPIKey,
			envValue: "env-openai-key",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.STT.APIKey != "env-openai-key" {
					t.Fatalf("expected stt.api_key to be overridden, got %q", cfg.STT.APIKey)
				}
			},
		},
		{
			name:     "stt provider overrides toml",
			envKey:   EnvSTTProvider,
			envValue: "openrouter",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.STT.Provider != "openrouter" {
					t.Fatalf("expected stt.provider to be overridden, got %q", cfg.STT.Provider)
				}
			},
		},
		{
			name:     "stt model overrides toml",
			envKey:   EnvSTTModel,
			envValue: "gpt-4o-mini-transcribe",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.STT.Model != "gpt-4o-mini-transcribe" {
					t.Fatalf("expected stt.model to be overridden, got %q", cfg.STT.Model)
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
		{
			name:     "youtube caption languages override toml",
			envKey:   EnvYouTubeCaptionLanguages,
			envValue: "orig, pt , es",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if !reflect.DeepEqual(cfg.YouTube.CaptionLanguages, []string{"orig", "pt", "es"}) {
					t.Fatalf("expected youtube.caption_languages to be overridden, got %#v", cfg.YouTube.CaptionLanguages)
				}
			},
		},
		{
			name:     "decisions model env overrides default",
			envKey:   EnvDecisionsModel,
			envValue: " typesafe/jev-1.14 ",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.Decisions.Model != "typesafe/jev-1.14" {
					t.Fatalf("expected decisions.model to be overridden, got %q", cfg.Decisions.Model)
				}
			},
		},
		{
			name:     "generation model env overrides default",
			envKey:   EnvGenerationModel,
			envValue: "acme/gen-1",
			assert: func(t *testing.T, cfg Config) {
				t.Helper()
				if cfg.Generation.Model != "acme/gen-1" {
					t.Fatalf("expected generation.model to be overridden, got %q", cfg.Generation.Model)
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

[stt]
provider = "openai"
api_key = "toml-openai-key"
api_url = "https://toml.openai.internal"
model = "gpt-4o-transcribe"

	[youtube]
	yt_dlp_path = "/toml/bin/yt-dlp"
	proxy = "http://toml.proxy:8080"
		cookies_file = "/tmp/toml-cookies.txt"
		user_agent = "toml-agent"
		caption_languages = ["orig", "en"]
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
