package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	// DefaultDotEnvPath is the local dotenv file loaded by the CLI entrypoint.
	DefaultDotEnvPath = ".env"

	// EnvConfigPath overrides the config file path.
	EnvConfigPath = "APP_CONFIG"

	// EnvFirecrawlAPIKey stores the Firecrawl API key override.
	EnvFirecrawlAPIKey = "FIRECRAWL_API_KEY"

	// EnvFirecrawlAPIURL stores the Firecrawl API URL override.
	EnvFirecrawlAPIURL = "FIRECRAWL_API_URL"

	// EnvOpenRouterAPIKey stores the OpenRouter API key override.
	EnvOpenRouterAPIKey = "OPENROUTER_API_KEY"

	// EnvOpenRouterAPIURL stores the OpenRouter API URL override.
	EnvOpenRouterAPIURL = "OPENROUTER_API_URL"
)

// ApplyEnvOverrides overlays config values that are sourced from environment
// variables at runtime.
func ApplyEnvOverrides(cfg *Config) {
	if cfg == nil {
		return
	}

	if value, ok := os.LookupEnv(EnvFirecrawlAPIKey); ok && value != "" {
		cfg.Firecrawl.APIKey = value
	}
	if value, ok := os.LookupEnv(EnvFirecrawlAPIURL); ok && value != "" {
		cfg.Firecrawl.APIURL = value
	}
	if value, ok := os.LookupEnv(EnvOpenRouterAPIKey); ok && value != "" {
		cfg.OpenRouter.APIKey = value
	}
	if value, ok := os.LookupEnv(EnvOpenRouterAPIURL); ok && value != "" {
		cfg.OpenRouter.APIURL = value
	}
}

// LoadDotEnvIfPresent loads a local dotenv file without overriding env vars
// already supplied by the shell or process manager.
func LoadDotEnvIfPresent(path string) error {
	if path == "" {
		path = DefaultDotEnvPath
	}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat dotenv %q: %w", path, err)
	}
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("load dotenv %q: %w", path, err)
	}
	return nil
}
