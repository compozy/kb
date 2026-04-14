package cli

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	kgenerate "github.com/compozy/kb/internal/generate"
	"github.com/compozy/kb/internal/models"
)

var runGenerate = kgenerate.GenerateWithObserver

func supportedCodebaseLanguagesHelp() string {
	languages := append([]string(nil), models.SupportedLanguageNames()...)
	slices.Sort(languages)
	return fmt.Sprintf("Supported languages: %s.", strings.Join(languages, ", "))
}

func newGenerateCommand() *cobra.Command {
	options := models.GenerateOptions{}
	var legacyOutputPath string
	progressModeValue := string(generateProgressAuto)
	logFormatValue := string(generateLogFormatText)

	command := &cobra.Command{
		Use:    "generate <path>",
		Short:  "Generate an Obsidian-compatible knowledge vault from a repository",
		Long:   strings.Join([]string{"Generate an Obsidian-compatible knowledge vault from a repository.", "", supportedCodebaseLanguagesHelp()}, "\n"),
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			options.RootPath = args[0]

			options.VaultPath = commandVaultValue(cmd, options.VaultPath)
			if options.VaultPath == "" && legacyOutputPath != "" {
				options.VaultPath = legacyOutputPath
			}

			summary, err := runGeneratePipeline(cmd, options, progressModeValue, logFormatValue)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)

			return encoder.Encode(summary)
		},
	}

	command.Flags().StringVar(&legacyOutputPath, "output", "", "Deprecated alias for --vault")
	_ = command.Flags().MarkDeprecated("output", "use --vault instead")
	command.Flags().StringVar(&options.TopicSlug, "topic", "", "Override the generated topic slug")
	command.Flags().StringVar(&options.Title, "title", "", "Override the generated topic title")
	command.Flags().StringVar(&options.Domain, "domain", "", "Override the generated topic domain")
	command.Flags().StringArrayVar(&options.IncludePatterns, "include", nil, "Re-include a path pattern that would otherwise be ignored; repeatable")
	command.Flags().StringArrayVar(&options.ExcludePatterns, "exclude", nil, "Exclude an additional path pattern from scanning; repeatable")
	command.Flags().BoolVar(&options.DryRun, "dry-run", false, "Inspect scan and adapter selection without writing any files")
	command.Flags().BoolVar(&options.Semantic, "semantic", false, "Enable semantic analysis when the underlying adapters support it")
	command.Flags().StringVar(&progressModeValue, "progress", progressModeValue, "Progress rendering mode: auto, always, or never")
	command.Flags().StringVar(&logFormatValue, "log-format", logFormatValue, "Stderr event format: text or json")

	return command
}

func runGeneratePipeline(
	cmd *cobra.Command,
	options models.GenerateOptions,
	progressModeValue string,
	logFormatValue string,
) (models.GenerationSummary, error) {
	progressMode, err := parseGenerateProgressMode(progressModeValue)
	if err != nil {
		return models.GenerationSummary{}, err
	}

	logFormat, err := parseGenerateLogFormat(logFormatValue)
	if err != nil {
		return models.GenerationSummary{}, err
	}

	return runGenerate(
		commandContext(cmd),
		options,
		newGenerateObserver(cmd.ErrOrStderr(), progressMode, logFormat),
	)
}
