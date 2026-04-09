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
		Use:           "kodebase",
		Short:         "Turn source code repositories into Obsidian knowledge vaults",
		Long:          "Kodebase analyzes source code repositories and generates Obsidian-compatible knowledge vaults with graph and metrics context.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	command.AddCommand(newGenerateCommand())
	command.AddCommand(newInspectCommand())
	command.AddCommand(newSearchCommand())
	command.AddCommand(newIndexCommand())
	command.AddCommand(newVersionCommand())

	return command
}
