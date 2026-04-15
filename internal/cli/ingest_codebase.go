package cli

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/models"
	ktopic "github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/vault"
)

type codebaseIngestResult struct {
	Topic      string                   `json:"topic"`
	SourceType models.SourceKind        `json:"sourceType"`
	FilePath   string                   `json:"filePath"`
	Title      string                   `json:"title"`
	Summary    models.GenerationSummary `json:"summary"`
}

func newIngestCodebaseCommand() *cobra.Command {
	options := models.GenerateOptions{}
	var topic string
	var legacyOutputPath string
	progressModeValue := string(generateProgressAuto)
	logFormatValue := string(generateLogFormatText)

	command := &cobra.Command{
		Use:   "codebase <path>",
		Short: "Analyze a codebase and ingest the generated KB artifacts into a topic",
		Long: strings.Join([]string{
			"Analyze a codebase and ingest the generated KB artifacts into a topic.",
			"",
			supportedCodebaseLanguagesHelp(),
		}, "\n"),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveCodebaseIngestTarget(cmd, "ingest codebase", args[0], topic, options.Title, options.Domain, legacyOutputPath)
			if err != nil {
				return err
			}
			if strings.TrimSpace(legacyOutputPath) != "" && commandVaultValue(cmd, "") == "" {
				if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "Flag --output has been deprecated, use --vault instead"); err != nil {
					return fmt.Errorf("ingest codebase: write deprecation warning: %w", err)
				}
			}

			options.RootPath = args[0]
			options.VaultPath = target.VaultPath
			options.TopicSlug = target.TopicInfo.Slug
			options.Title = target.TopicInfo.Title
			options.Domain = target.TopicInfo.Domain

			summary, err := runGeneratePipeline(cmd, options, progressModeValue, logFormatValue)
			if err != nil {
				return fmt.Errorf("ingest codebase: %w", err)
			}

			return writeJSON(cmd, codebaseIngestResult{
				Topic:      target.TopicInfo.Slug,
				SourceType: models.SourceKindCodebaseFile,
				FilePath:   path.Join(target.TopicInfo.Slug, "raw", "codebase"),
				Title:      target.TopicInfo.Title,
				Summary:    summary,
			})
		},
	}

	requireTopicFlag(command, &topic)
	command.Flags().StringVar(&legacyOutputPath, "output", "", "Deprecated alias for --vault")
	command.Flags().StringVar(&options.Title, "title", "", "Override the generated topic title when bootstrapping a missing topic")
	command.Flags().StringVar(&options.Domain, "domain", "", "Override the generated topic domain when bootstrapping a missing topic")
	command.Flags().StringArrayVar(&options.IncludePatterns, "include", nil, "Re-include a path pattern that would otherwise be ignored; repeatable")
	command.Flags().StringArrayVar(&options.ExcludePatterns, "exclude", nil, "Exclude an additional path pattern from scanning; repeatable")
	command.Flags().BoolVar(&options.DryRun, "dry-run", false, "Inspect scan and adapter selection without writing any files")
	command.Flags().BoolVar(&options.Semantic, "semantic", false, "Enable semantic analysis when the underlying adapters support it")
	command.Flags().StringVar(&progressModeValue, "progress", progressModeValue, "Progress rendering mode: auto, always, or never")
	command.Flags().StringVar(&logFormatValue, "log-format", logFormatValue, "Stderr event format: text or json")

	return command
}

func resolveCodebaseIngestTarget(
	cmd *cobra.Command,
	action string,
	rootPath string,
	topicSlug string,
	title string,
	domain string,
	legacyOutputPath string,
) (ingestTarget, error) {
	vaultPath, err := resolveCodebaseVaultPath(cmd, action, rootPath, legacyOutputPath)
	if err != nil {
		return ingestTarget{}, err
	}

	cleanTopicSlug := strings.TrimSpace(topicSlug)
	topicInfo, err := runIngestTopicInfo(vaultPath, cleanTopicSlug)
	if err == nil {
		if strings.TrimSpace(title) != "" || strings.TrimSpace(domain) != "" {
			return ingestTarget{}, fmt.Errorf(
				"%s: --title and --domain are bootstrap-only and cannot be used when topic %q already exists",
				action,
				cleanTopicSlug,
			)
		}

		return ingestTarget{
			TopicInfo: topicInfo,
			VaultPath: vaultPath,
		}, nil
	}
	if !errors.Is(err, ktopic.ErrTopicNotFound) {
		return ingestTarget{}, fmt.Errorf("%s: %w", action, err)
	}

	return ingestTarget{
		TopicInfo: models.TopicInfo{
			Slug:     cleanTopicSlug,
			Title:    codebaseBootstrapTitle(cleanTopicSlug, title),
			Domain:   codebaseBootstrapDomain(cleanTopicSlug, domain),
			RootPath: filepath.Join(vaultPath, cleanTopicSlug),
		},
		VaultPath: vaultPath,
	}, nil
}

func resolveCodebaseVaultPath(cmd *cobra.Command, action string, rootPath string, legacyOutputPath string) (string, error) {
	explicitVaultPath := commandVaultValue(cmd, "")
	if explicitVaultPath == "" {
		explicitVaultPath = strings.TrimSpace(legacyOutputPath)
	}
	if explicitVaultPath != "" {
		absoluteVaultPath, err := filepath.Abs(explicitVaultPath)
		if err != nil {
			return "", fmt.Errorf("%s: resolve vault path %q: %w", action, explicitVaultPath, err)
		}
		return absoluteVaultPath, nil
	}

	absoluteRootPath, err := filepath.Abs(strings.TrimSpace(rootPath))
	if err != nil {
		return "", fmt.Errorf("%s: resolve root path %q: %w", action, rootPath, err)
	}

	return filepath.Join(absoluteRootPath, ".kb", "vault"), nil
}

func codebaseBootstrapTitle(topicSlug string, override string) string {
	if strings.TrimSpace(override) != "" {
		return strings.TrimSpace(override)
	}

	return vault.DeriveTopicTitle(topicSlug)
}

func codebaseBootstrapDomain(topicSlug string, override string) string {
	domainSource := topicSlug
	if strings.TrimSpace(override) != "" {
		domainSource = strings.TrimSpace(override)
	}

	return vault.DeriveTopicDomain(domainSource)
}
