package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const (
	// DefaultDotEnvPath is the local dotenv file loaded by the CLI entrypoint.
	DefaultDotEnvPath = ".env"

	// ProjectConfigFileName is the default repository-local kb config file.
	ProjectConfigFileName = "kb.toml"

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

	// EnvYouTubeProxy stores the YouTube proxy URL override.
	EnvYouTubeProxy = "YOUTUBE_PROXY"

	// EnvYouTubeYTDLPPath stores the yt-dlp executable path override.
	EnvYouTubeYTDLPPath = "YOUTUBE_YT_DLP_PATH"

	// EnvYouTubeCookiesFile stores the YouTube cookies file override.
	EnvYouTubeCookiesFile = "YOUTUBE_COOKIES_FILE"

	// EnvYouTubeUserAgent stores the YouTube user agent override.
	EnvYouTubeUserAgent = "YOUTUBE_USER_AGENT"
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
	if value, ok := os.LookupEnv(EnvYouTubeYTDLPPath); ok && value != "" {
		cfg.YouTube.YTDLPPath = value
	}
	if value, ok := os.LookupEnv(EnvYouTubeProxy); ok && value != "" {
		cfg.YouTube.Proxy = value
	}
	if value, ok := os.LookupEnv(EnvYouTubeCookiesFile); ok && value != "" {
		cfg.YouTube.CookiesFile = value
	}
	if value, ok := os.LookupEnv(EnvYouTubeUserAgent); ok && value != "" {
		cfg.YouTube.UserAgent = value
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

// DiscoverProjectConfigPath walks up from cwd looking for kb.toml.
func DiscoverProjectConfigPath(cwd string) (string, bool, error) {
	resolvedCWD := strings.TrimSpace(cwd)
	if resolvedCWD == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			return "", false, fmt.Errorf("resolve cwd: %w", err)
		}
		resolvedCWD = currentPath
	}

	absolutePath, err := filepath.Abs(resolvedCWD)
	if err != nil {
		return "", false, fmt.Errorf("resolve config cwd %q: %w", cwd, err)
	}

	currentPath := absolutePath
	for {
		candidatePath := filepath.Join(currentPath, ProjectConfigFileName)
		info, err := os.Stat(candidatePath)
		switch {
		case err == nil && !info.IsDir():
			return candidatePath, true, nil
		case err == nil:
		case errors.Is(err, os.ErrNotExist):
		default:
			return "", false, fmt.Errorf("stat config %q: %w", candidatePath, err)
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", false, nil
		}
		currentPath = parentPath
	}
}

// ResolveVaultRoot resolves cfg.Root relative to the directory containing
// configPath when cfg.Root is relative.
func ResolveVaultRoot(configPath string, cfg VaultConfig) (string, error) {
	root := strings.TrimSpace(cfg.Root)
	if root == "" {
		root = defaultVaultRoot
	}
	if filepath.IsAbs(root) {
		return filepath.Clean(root), nil
	}

	configDir := filepath.Dir(configPath)
	absoluteRoot, err := filepath.Abs(filepath.Join(configDir, root))
	if err != nil {
		return "", fmt.Errorf("resolve vault root %q: %w", root, err)
	}
	return filepath.Clean(absoluteRoot), nil
}
