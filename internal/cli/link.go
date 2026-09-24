package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	klink "github.com/compozy/kb/internal/link"
	"github.com/compozy/kb/internal/session"
)

func newLinkCommand() *cobra.Command {
	var (
		flags  session.Flags
		all    bool
		dryRun bool
	)

	command := &cobra.Command{
		Use:   "link <topic>",
		Short: "Link a topic's documents to its wiki articles with the decision model",
		Long: "Proposes link candidates by code (mentions, BM25, shared structure), judges each document in one\n" +
			"decision request and writes frontmatter relations (both modes) and body links (decisions.mode: apply).\n" +
			"Incremental: documents whose body and candidate set are unchanged are skipped unless --all.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openSession(cmd, "kb link", args[0], flags)
			if err != nil {
				return err
			}
			report, err := klink.Run(commandContext(cmd), s, klink.Options{All: all, DryRun: dryRun})
			if err != nil {
				s.WriteSummary(cmd.ErrOrStderr())
				return fmt.Errorf("kb link: %w", err)
			}
			for _, line := range report.Lines() {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), line); err != nil {
					return fmt.Errorf("kb link: write output: %w", err)
				}
			}
			s.WriteSummary(cmd.ErrOrStderr())
			return nil
		},
	}

	command.Flags().BoolVar(&all, "all", false, "Relink every document, including unchanged ones")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Print what would change without writing files, state or review items")
	bindDecisionFlags(command, &flags)
	return command
}
