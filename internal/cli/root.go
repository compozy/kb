package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// ExecuteContext runs the CLI with the provided context.
func ExecuteContext(ctx context.Context) error {
	return newRootCommand().ExecuteContext(ctx)
}

func newRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "kb",
		Short:         "Build and inspect topic-based knowledge bases",
		Long:          "kb scaffolds topic-based knowledge vaults, generates codebase snapshots, indexes collections, searches content, and inspects code intelligence.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	bindRootPersistentFlags(command)
	command.AddCommand(newTopicCommand())
	command.AddCommand(newIngestCommand())
	command.AddCommand(newGenerateCommand())
	command.AddCommand(newInspectCommand())
	command.AddCommand(newSearchCommand())
	command.AddCommand(newIndexCommand())
	command.AddCommand(newVersionCommand())

	return command
}
