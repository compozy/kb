package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
	ktopic "github.com/compozy/kb/internal/topic"
)

var runIngest = kingest.Ingest
var runIngestTopicInfo = ktopic.Info
var runIngestTopicNew = ktopic.New
var ingestGetwd = os.Getwd
var loadIngestConfig = loadCLIConfig
var existingYouTubeVideoIDs = kingest.ExistingYouTubeVideoIDs

type ingestTarget struct {
	TopicInfo models.TopicInfo
	VaultPath string
}

func newIngestCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest source material into an existing knowledge base topic",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.AddCommand(
		newIngestURLCommand(),
		newIngestFileCommand(),
		newIngestYouTubeCommand(),
		newIngestChannelCommand(),
		newIngestInstagramCommand(),
		newIngestCodebaseCommand(),
		newIngestBookmarksCommand(),
	)

	return command
}

func resolveIngestTarget(cmd *cobra.Command, action string, topicSlug string) (ingestTarget, error) {
	vaultPath, err := resolveCommandVaultPath(cmd, ingestGetwd, action)
	if err != nil {
		return ingestTarget{}, err
	}

	topicInfo, err := runIngestTopicInfo(vaultPath, strings.TrimSpace(topicSlug))
	if err != nil {
		return ingestTarget{}, fmt.Errorf("%s: %w", action, err)
	}

	return ingestTarget{
		TopicInfo: topicInfo,
		VaultPath: vaultPath,
	}, nil
}

func resolveIngestTargetWithMissingTopicHint(
	cmd *cobra.Command,
	action string,
	topicSlug string,
	title string,
	domain string,
) (ingestTarget, error) {
	target, err := resolveIngestTarget(cmd, action, topicSlug)
	if err == nil {
		return target, nil
	}
	if errors.Is(err, ktopic.ErrTopicNotFound) {
		return ingestTarget{}, fmt.Errorf("%w; create it with: %s", err, missingTopicCreateCommand(topicSlug, title, domain))
	}
	return ingestTarget{}, err
}

func missingTopicCreateCommand(topicSlug string, title string, domain string) string {
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		cleanTitle = "<title>"
	}
	cleanDomain := strings.TrimSpace(domain)
	if cleanDomain == "" {
		cleanDomain = "<domain>"
	}
	return fmt.Sprintf("kb topic new %s %q %s", strings.TrimSpace(topicSlug), cleanTitle, cleanDomain)
}

func requireTopicFlag(command *cobra.Command, topic *string) {
	command.Flags().StringVar(topic, "topic", "", "Target topic slug inside the vault")
	_ = command.MarkFlagRequired("topic")
}

func resolveInputFile(action string, sourcePath string) (string, error) {
	sourcePath = strings.TrimSpace(sourcePath)
	if sourcePath == "" {
		return "", fmt.Errorf("%s: source path is required", action)
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("%s: stat source path %q: %w", action, sourcePath, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s: source path %q must be a file", action, sourcePath)
	}

	return sourcePath, nil
}

func commandContext(cmd *cobra.Command) context.Context {
	ctx := cmd.Context()
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func loadCLIConfig() (kconfig.Config, error) {
	if err := kconfig.LoadDotEnvIfPresent(""); err != nil {
		return kconfig.Config{}, fmt.Errorf("load dotenv: %w", err)
	}

	cfgPath := strings.TrimSpace(os.Getenv(kconfig.EnvConfigPath))
	if cfgPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			if discoveredPath, found, err := kconfig.DiscoverProjectConfigPath(cwd); err != nil {
				return kconfig.Config{}, fmt.Errorf("discover project config: %w", err)
			} else if found {
				cfgPath = discoveredPath
			}
		} else {
			return kconfig.Config{}, fmt.Errorf("resolve cwd: %w", err)
		}
	}
	cfg, err := kconfig.Load(cfgPath)
	if err != nil {
		return kconfig.Config{}, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}

func writeJSON(cmd *cobra.Command, payload any) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(payload)
}
