package cli

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"

	kgenerate "github.com/user/go-devstack/internal/generate"
	"github.com/user/go-devstack/internal/models"
)

var runGenerate = kgenerate.GenerateWithObserver

func newGenerateCommand() *cobra.Command {
	options := models.GenerateOptions{}
	var legacyOutputPath string
	progressModeValue := string(generateProgressAuto)
	logFormatValue := string(generateLogFormatText)

	command := &cobra.Command{
		Use:   "generate <path>",
		Short: "Generate an Obsidian-compatible knowledge vault from a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.RootPath = args[0]

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			progressMode, err := parseGenerateProgressMode(progressModeValue)
			if err != nil {
				return err
			}

			logFormat, err := parseGenerateLogFormat(logFormatValue)
			if err != nil {
				return err
			}
			if options.VaultPath == "" && legacyOutputPath != "" {
				options.VaultPath = legacyOutputPath
			}

			summary, err := runGenerate(ctx, options, newGenerateObserver(cmd.ErrOrStderr(), progressMode, logFormat))
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)

			return encoder.Encode(summary)
		},
	}

	command.Flags().StringVar(&options.VaultPath, "vault", "", "Vault root containing the generated topic")
	command.Flags().StringVar(&legacyOutputPath, "output", "", "Deprecated alias for --vault")
	_ = command.Flags().MarkDeprecated("output", "use --vault instead")
	command.Flags().StringVar(&options.TopicSlug, "topic", "", "Override the generated topic slug")
	command.Flags().StringVar(&options.Title, "title", "", "Override the generated topic title")
	command.Flags().StringVar(&options.Domain, "domain", "", "Override the generated topic domain")
	command.Flags().StringArrayVar(&options.IncludePatterns, "include", nil, "Re-include a path pattern that would otherwise be ignored; repeatable")
	command.Flags().StringArrayVar(&options.ExcludePatterns, "exclude", nil, "Exclude an additional path pattern from scanning; repeatable")
	command.Flags().BoolVar(&options.Semantic, "semantic", false, "Enable semantic analysis when the underlying adapters support it")
	command.Flags().StringVar(&progressModeValue, "progress", progressModeValue, "Progress rendering mode: auto, always, or never")
	command.Flags().StringVar(&logFormatValue, "log-format", logFormatValue, "Stderr event format: text or json")

	return command
}
