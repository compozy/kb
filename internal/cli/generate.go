package cli

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/spf13/cobra"

	kgenerate "github.com/user/go-devstack/internal/generate"
	"github.com/user/go-devstack/internal/logger"
	"github.com/user/go-devstack/internal/models"
)

var runGenerate = kgenerate.Generate

func newGenerateCommand() *cobra.Command {
	options := models.GenerateOptions{}

	command := &cobra.Command{
		Use:   "generate <path>",
		Short: "Generate an Obsidian-compatible knowledge vault from a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.RootPath = args[0]

			log, err := logger.New("info", logger.WithWriter(cmd.ErrOrStderr()))
			if err != nil {
				return err
			}

			previous := slog.Default()
			slog.SetDefault(log)
			defer slog.SetDefault(previous)

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			summary, err := runGenerate(ctx, options)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)

			return encoder.Encode(summary)
		},
	}

	command.Flags().StringVar(&options.OutputPath, "output", "", "Vault root where the generated topic will be written")
	command.Flags().StringVar(&options.Topic, "topic", "", "Override the generated topic slug")
	command.Flags().StringVar(&options.Title, "title", "", "Override the generated topic title")
	command.Flags().StringVar(&options.Domain, "domain", "", "Override the generated topic domain")
	command.Flags().StringArrayVar(&options.IncludePatterns, "include", nil, "Re-include a path pattern that would otherwise be ignored; repeatable")
	command.Flags().StringArrayVar(&options.ExcludePatterns, "exclude", nil, "Exclude an additional path pattern from scanning; repeatable")
	command.Flags().BoolVar(&options.Semantic, "semantic", false, "Enable semantic analysis when the underlying adapters support it")

	return command
}
