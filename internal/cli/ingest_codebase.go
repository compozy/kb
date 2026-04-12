package cli

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/user/kb/internal/models"
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
	progressModeValue := string(generateProgressAuto)
	logFormatValue := string(generateLogFormatText)

	command := &cobra.Command{
		Use:   "codebase <path>",
		Short: "Analyze a codebase and ingest the generated KB artifacts into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest codebase", topic)
			if err != nil {
				return err
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
	command.Flags().StringArrayVar(&options.IncludePatterns, "include", nil, "Re-include a path pattern that would otherwise be ignored; repeatable")
	command.Flags().StringArrayVar(&options.ExcludePatterns, "exclude", nil, "Exclude an additional path pattern from scanning; repeatable")
	command.Flags().BoolVar(&options.Semantic, "semantic", false, "Enable semantic analysis when the underlying adapters support it")
	command.Flags().StringVar(&progressModeValue, "progress", progressModeValue, "Progress rendering mode: auto, always, or never")
	command.Flags().StringVar(&logFormatValue, "log-format", logFormatValue, "Stderr event format: text or json")

	return command
}
