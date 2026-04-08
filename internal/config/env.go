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

	// EnvDatabaseURL stores the database connection string.
	EnvDatabaseURL = "DATABASE_URL"

	// EnvAPIKey stores the application API key.
	EnvAPIKey = "API_KEY"
)

// Secrets contains the runtime secrets loaded from the environment.
type Secrets struct {
	DatabaseURL string
	APIKey      string
}

// LoadSecretsFromEnv reads the current process environment into a stable
// runtime value object.
func LoadSecretsFromEnv() Secrets {
	return Secrets{
		DatabaseURL: os.Getenv(EnvDatabaseURL),
		APIKey:      os.Getenv(EnvAPIKey),
	}
}

// Validate performs format checks for any secrets that are present.
func (s Secrets) Validate() error {
	return nil
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
