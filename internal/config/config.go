// Package config loads TOML configuration and environment-backed runtime overrides for the kb CLI.
package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	defaultFirecrawlAPIURL    = "https://api.firecrawl.dev"
	defaultOpenRouterAPIURL   = "https://openrouter.ai/api"
	defaultOpenRouterSTTModel = "google/gemini-2.5-flash"
)

// Config contains the complete TOML-backed runtime configuration.
type Config struct {
	App        AppConfig        `toml:"app"`
	Log        LogConfig        `toml:"log"`
	Firecrawl  FirecrawlConfig  `toml:"firecrawl"`
	OpenRouter OpenRouterConfig `toml:"openrouter"`
}

// AppConfig contains the application identity and environment.
type AppConfig struct {
	Name string `toml:"name"`
	Env  string `toml:"env"`
}

// LogConfig controls structured logging output.
type LogConfig struct {
	Level string `toml:"level"`
}

// FirecrawlConfig controls URL scraping API access.
type FirecrawlConfig struct {
	APIKey string `toml:"api_key"`
	APIURL string `toml:"api_url"`
}

// OpenRouterConfig controls the STT fallback provider.
type OpenRouterConfig struct {
	APIKey   string `toml:"api_key"`
	APIURL   string `toml:"api_url"`
	STTModel string `toml:"stt_model"`
}

// Default returns a sane starting configuration.
func Default() Config {
	return Config{
		App: AppConfig{
			Name: "app",
			Env:  "development",
		},
		Log: LogConfig{
			Level: "info",
		},
		Firecrawl: FirecrawlConfig{
			APIURL: defaultFirecrawlAPIURL,
		},
		OpenRouter: OpenRouterConfig{
			APIURL:   defaultOpenRouterAPIURL,
			STTModel: defaultOpenRouterSTTModel,
		},
	}
}

// Load reads and validates the TOML config file, then overlays runtime secrets
// from the environment.
func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		if err := decodeFile(path, &cfg); err != nil {
			return Config{}, err
		}
	}
	ApplyEnvOverrides(&cfg)
	return cfg, cfg.Validate()
}

func decodeFile(path string, cfg *Config) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat config %q: %w", path, err)
	}

	meta, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return fmt.Errorf("decode config %q: %w", path, err)
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		sort.Strings(keys)
		return fmt.Errorf("unknown config keys: %s", strings.Join(keys, ", "))
	}

	return nil
}

// Validate ensures the config is internally consistent before runtime startup.
func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	if err := c.Log.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate ensures application identity settings are usable.
func (c AppConfig) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("app.name is required")
	}
	switch strings.ToLower(strings.TrimSpace(c.Env)) {
	case "development", "staging", "production":
	default:
		return fmt.Errorf("app.env must be development, staging, or production: %q", c.Env)
	}
	return nil
}

// Validate ensures the log level is supported.
func (c LogConfig) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.Level)) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be debug, info, warn, or error: %q", c.Level)
	}
	return nil
}
