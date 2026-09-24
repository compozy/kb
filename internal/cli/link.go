package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	klink "github.com/compozy/kb/internal/link"
	"github.com/compozy/kb/internal/qmd"
	"github.com/compozy/kb/internal/session"
)

// openQMDCandidates opens the optional qmd vector candidates of a topic
// (spec §9.1 step 4, §10 step 1): nil and a one-line note when qmd is
// missing or the topic collection is absent or stale. Tests replace it.
var openQMDCandidates = func(ctx context.Context, collection, topicRoot string) (*qmd.Candidates, string) {
	return qmd.NewClient().OpenCandidates(ctx, collection, topicRoot)
}

// qmdCandidates returns the topic's qmd vector candidates, or nil when off
// (--no-qmd, or not available or not fresh: a one-line note goes to
// stderr). onNote, when set, is printed when the candidates are on.
func qmdCandidates(cmd *cobra.Command, s *session.Session, disabled bool, onNote string) *qmd.Candidates {
	if disabled {
		return nil
	}
	candidates, note := openQMDCandidates(commandContext(cmd), s.Topic.Slug, s.Root())
	if candidates == nil {
		if note != "" {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), note)
		}
		return nil
	}
	if onNote != "" {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), onNote)
	}
	return candidates
}

func newLinkCommand() *cobra.Command {
	var (
		flags  session.Flags
		all    bool
		dryRun bool
		noQMD  bool
	)

	command := &cobra.Command{
		Use:   "link <topic>",
		Short: "Link a topic's documents to its wiki articles with the decision model",
		Long: "Proposes link candidates by code (mentions, BM25, shared structure), judges each document in one\n" +
			"decision request and writes frontmatter relations (both modes) and body links (decisions.mode: apply).\n" +
			"When qmd is installed and the topic collection is fresh (indexed after the newest topic file), its\n" +
			"vector neighbours are added as candidates; qmd scores never enter bands.\n" +
			"Incremental: documents whose body and candidate set are unchanged are skipped unless --all.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openSession(cmd, "kb link", args[0], flags)
			if err != nil {
				return err
			}
			options := klink.Options{All: all, DryRun: dryRun}
			onNote := "qmd candidates: on (collection " + s.Topic.Slug + ", about 1.2 s per document; --no-qmd skips them)"
			if candidates := qmdCandidates(cmd, s, noQMD, onNote); candidates != nil {
				options.Neighbours = candidates
			}
			report, err := klink.Run(commandContext(cmd), s, options)
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
	command.Flags().BoolVar(&noQMD, "no-qmd", false, "Do not add qmd vector neighbours as link candidates")
	bindDecisionFlags(command, &flags)
	return command
}
