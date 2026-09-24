package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/actions"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/scope"
	"github.com/compozy/kb/internal/session"
)

// Verdict command names.
const (
	actionsAccept = actions.Accept
	actionsReject = actions.Reject
)

// reviewVerdictOptions are the flags of `kb review accept|reject`.
type reviewVerdictOptions struct {
	Topic      string
	Queue      string
	AllPurpose string
	Above      float64
	TopicWide  bool
	Format     string
	Flags      session.Flags
}

// reviewSelection is what selects the items of one verdict run.
type reviewSelection struct {
	IDs       []string
	Queue     string
	Purpose   string
	Above     float64
	AboveSet  bool
	TopicWide bool
}

// newReviewVerdictCommand builds `kb review accept` or `kb review reject`.
func newReviewVerdictCommand(verdict string) *cobra.Command {
	options := &reviewVerdictOptions{}
	short := "Accept review items: perform their action (restore, quarantine, rescue, recapture, write the link, create the concept)"
	if verdict == actionsReject {
		short = "Reject review items: dismiss them (quarantine a gated source; keep a remove/recapture source)"
	}
	command := &cobra.Command{
		Use:   verdict + " <id>...",
		Short: short,
		Long: short + ".\nSelect items by id, or in bulk with --queue and/or --all-purpose plus --above <p> or --topic-wide.\n" +
			"Every verdict is recorded as a label for calibration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReviewVerdict(cmd, verdict, options, reviewSelection{
				IDs: args, Queue: strings.TrimSpace(options.Queue), Purpose: strings.TrimSpace(options.AllPurpose),
				Above: options.Above, AboveSet: cmd.Flags().Changed("above"), TopicWide: options.TopicWide,
			})
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.Topic, "topic", "", "Topic slug")
	flags.StringVar(&options.Queue, "queue", "", "Only items of one queue: "+strings.Join(review.Queues, "|"))
	flags.StringVar(&options.AllPurpose, "all-purpose", "", "Only items of one decision purpose (relevance, quality, link, ...)")
	flags.Float64Var(&options.Above, "above", 0, "Bulk: pending items with probability ≥ p in the selected queue/purpose")
	flags.BoolVar(&options.TopicWide, "topic-wide", false, "Bulk: every pending item of the selected queue/purpose")
	flags.StringVar(&options.Format, "format", "table", "Output format (table|json)")
	bindDecisionFlags(command, &options.Flags)
	_ = command.MarkFlagRequired("topic")
	return command
}

// selectReviewItems returns the pending items a verdict applies to: the
// given ids (each must be pending and match the queue/purpose filters), or,
// in bulk, the pending items of the selected queue and/or purpose with
// probability ≥ --above (or all of them with --topic-wide). Bulk selection
// needs a queue or a purpose so a stray flag never touches every queue.
func selectReviewItems(store *review.Store, sel reviewSelection) ([]review.Item, error) {
	if sel.Queue != "" && !slices.Contains(review.Queues, sel.Queue) {
		return nil, fmt.Errorf("unknown queue %q; expected one of %s", sel.Queue, strings.Join(review.Queues, ", "))
	}
	matches := func(item review.Item) bool {
		return (sel.Queue == "" || item.Queue == sel.Queue) && (sel.Purpose == "" || item.Purpose == sel.Purpose)
	}
	if len(sel.IDs) > 0 {
		if sel.TopicWide || sel.AboveSet {
			return nil, errors.New("pass ids or a bulk selector (--above, --topic-wide), not both")
		}
		items := make([]review.Item, 0, len(sel.IDs))
		for _, id := range sel.IDs {
			item, ok, err := store.Get(strings.TrimSpace(id))
			if err != nil {
				return nil, err
			}
			switch {
			case !ok:
				return nil, fmt.Errorf("unknown review item %q", id)
			case item.Status != review.StatusPending:
				return nil, fmt.Errorf("review item %s is already %s", id, item.Status)
			case !matches(item):
				return nil, fmt.Errorf("review item %s is in queue %s (purpose %s), not the selected one", id, item.Queue, item.Purpose)
			}
			items = append(items, item)
		}
		return items, nil
	}
	if !sel.TopicWide && !sel.AboveSet {
		return nil, errors.New("select items by id, or in bulk with --above <p> or --topic-wide")
	}
	if sel.Queue == "" && sel.Purpose == "" {
		return nil, errors.New("bulk selection needs --queue or --all-purpose")
	}
	pending, err := store.Pending(sel.Queue)
	if err != nil {
		return nil, err
	}
	items := make([]review.Item, 0, len(pending))
	for _, item := range pending {
		if !matches(item) || (sel.AboveSet && !sel.TopicWide && item.Probability < sel.Above) {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func runReviewVerdict(cmd *cobra.Command, verdict string, options *reviewVerdictOptions, sel reviewSelection) error {
	action := "review " + verdict
	format, err := reviewFormat(options.Format)
	if err != nil {
		return err
	}
	s, err := openSession(cmd, "kb "+action, options.Topic, options.Flags)
	if err != nil {
		return err
	}
	store := review.Open(s.Root(), s.Now)
	items, err := selectReviewItems(store, sel)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	ctx := commandContext(cmd)
	scraper := newFirecrawlScraper(firecrawlConfig(s.Config.Firecrawl))
	var runner *kingest.Run
	target := ingestTarget{TopicInfo: s.Topic, VaultPath: s.VaultPath}
	deps := actions.Deps{
		Scrape: func(ctx context.Context, sourceURL string, opts firecrawl.ScrapeOptions) (*firecrawl.ScrapeResult, error) {
			return scraper.ScrapeWithOptions(ctx, sourceURL, opts)
		},
		Rescue: func(ctx context.Context, row gate.SkippedRow) (string, error) {
			if runner == nil {
				runner = kingest.NewRun(s, kingest.RunOptions{Command: "rescue", Batch: row.Batch, Query: row.Query})
			}
			result, err := rescueSkipped(ctx, runner, target, rescueFetchers{cfg: s.Config, scraper: scraper}, row)
			if err != nil {
				return "", err
			}
			if result.FilePath == "" {
				return "", fmt.Errorf("not written: %s %s", result.Triage, result.TriageReason)
			}
			return result.FilePath, nil
		},
	}

	results := make([]actions.Result, 0, len(items))
	var failures []error
	for _, item := range items {
		result, err := actions.Apply(ctx, s, deps, item, verdict)
		if err != nil {
			if ctx.Err() != nil {
				return err
			}
			failures = append(failures, err)
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error %s: %v\n", item.ID, err)
			continue
		}
		results = append(results, result)
	}

	out := cmd.OutOrStdout()
	if format == "json" {
		err = writeReviewJSON(out, results)
	} else {
		for _, result := range results {
			if _, err = fmt.Fprintf(out, "%sed %s (%s) %s: %s\n", verdict, result.ID, result.Queue, result.Subject, firstNonBlank(result.Message, result.Action)); err != nil {
				break
			}
			if err = writeManualRepairs(out, result.Manual); err != nil {
				break
			}
		}
		if err == nil && len(items) == 0 {
			_, err = fmt.Fprintln(out, "no pending items selected")
		}
	}
	if err != nil {
		return err
	}
	if format == "json" {
		// JSON output keeps stdout machine-readable; manual repairs are
		// in each result and repeated on stderr for the terminal.
		for _, result := range results {
			_ = writeManualRepairs(cmd.ErrOrStderr(), result.Manual)
		}
	}
	if movedSources(results) {
		folders, countErr := scope.SourceFolderCounts(s.Root())
		if countErr != nil {
			failures = append(failures, countErr)
		} else {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), scope.FolderCountsLine(folders))
		}
	}
	if runner != nil {
		summary, finishErr := runner.Finish(ctx)
		for _, line := range summary.Lines() {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), line)
		}
		failures = append(failures, finishErr)
	}
	s.WriteSummary(cmd.ErrOrStderr())
	if err := errors.Join(failures...); err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	return nil
}

// writeManualRepairs prints each restore entry that needs a manual repair
// (file, line or key, and the text that was removed), indented under its
// result line (spec §7.1).
func writeManualRepairs(w io.Writer, manual []actions.ManualRepair) error {
	for _, entry := range manual {
		for _, line := range entry.Lines() {
			if _, err := fmt.Fprintln(w, "  "+line); err != nil {
				return err
			}
		}
	}
	return nil
}

// movedSources reports whether a verdict quarantined or restored a source,
// after which the new source counts per raw/ folder are printed (spec §7.1).
func movedSources(results []actions.Result) bool {
	for _, result := range results {
		if result.Action == actions.DidQuarantined || result.Action == actions.DidRestored {
			return true
		}
	}
	return false
}
