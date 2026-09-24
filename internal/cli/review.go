package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/vault"
)

var reviewGetwd = os.Getwd

// reviewTopic is a resolved topic for the file-only review subcommands.
type reviewTopic struct {
	vaultPath string
	root      string
	slug      string
}

type reviewListOptions struct {
	Queue  string
	Format string
}

// newReviewCommand builds `kb review`. `kb review <topic>` is `kb review list
// <topic>`. list, import-links, import-labels and calibrate read files only
// and never need the decision model.
func newReviewCommand() *cobra.Command {
	listOptions := &reviewListOptions{}
	command := &cobra.Command{
		Use:   "review <topic>",
		Short: "List, import and calibrate a topic's review queue and labels",
		Long: "Review decisions that fell in the review band, import labels from human links and\n" +
			"screening files, and calibrate thresholds on a holdout. `kb review <topic>` lists pending items.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runReviewList(cmd, listOptions, args[0])
		},
	}
	bindReviewListFlags(command, listOptions)

	command.AddCommand(newReviewListCommand())
	command.AddCommand(newReviewImportLinksCommand())
	command.AddCommand(newReviewImportLabelsCommand())
	command.AddCommand(newReviewCalibrateCommand())
	command.AddCommand(newReviewVerdictCommand(actionsAccept), newReviewVerdictCommand(actionsReject))

	return command
}

func bindReviewListFlags(command *cobra.Command, options *reviewListOptions) {
	command.Flags().StringVar(&options.Queue, "queue", "", "Only list one queue: "+strings.Join(review.Queues, "|"))
	command.Flags().StringVar(&options.Format, "format", "table", "Output format (table|json)")
}

func newReviewListCommand() *cobra.Command {
	options := &reviewListOptions{}
	command := &cobra.Command{
		Use:   "list <topic>",
		Short: "List pending review items with counts per queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReviewList(cmd, options, args[0])
		},
	}
	bindReviewListFlags(command, options)
	return command
}

func runReviewList(cmd *cobra.Command, options *reviewListOptions, topicSlug string) error {
	format, err := reviewFormat(options.Format)
	if err != nil {
		return err
	}
	queue := strings.TrimSpace(options.Queue)
	if queue != "" && !slices.Contains(review.Queues, queue) {
		return fmt.Errorf("review list: unknown queue %q; expected one of %s", queue, strings.Join(review.Queues, ", "))
	}
	topic, err := resolveReviewTopic(cmd, "review list", topicSlug)
	if err != nil {
		return err
	}
	titles := map[string]string{}
	if c, err := corpus.Load(topic.root, corpus.LoadOptions{}); err == nil {
		for _, doc := range c.Documents() {
			titles[doc.Path] = doc.Title
		}
	}
	view, err := review.BuildListView(review.Open(topic.root, nil), queue, func(subject string) string { return titles[subject] })
	if err != nil {
		return fmt.Errorf("review list: %w", err)
	}
	if format == "json" {
		return writeReviewJSON(cmd.OutOrStdout(), view)
	}
	return review.RenderList(cmd.OutOrStdout(), view)
}

func newReviewImportLinksCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import-links <topic>",
		Short: "Import existing human-written links between topic documents as positive link labels",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			topic, err := resolveReviewTopic(cmd, "review import-links", args[0])
			if err != nil {
				return err
			}
			count, err := review.ImportLinks(topic.vaultPath, topic.root, nil)
			if err != nil {
				return fmt.Errorf("review import-links: %w", err)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "imported %d link labels into %s\n", count, topic.slug)
			return err
		},
	}
}

type reviewImportLabelsOptions struct {
	From           string
	IDField        string
	Match          string
	DecisionField  string
	KeepValues     []string
	DecidedByField string
}

func newReviewImportLabelsCommand() *cobra.Command {
	options := &reviewImportLabelsOptions{}
	command := &cobra.Command{
		Use:   "import-labels <topic>",
		Short: "Import relevance labels from a screening or curation file (JSONL, JSON or CSV)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReviewImportLabels(cmd, options, args[0])
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.From, "from", "", "Screening or curation file (.jsonl, .json or .csv)")
	flags.StringVar(&options.IDField, "id-field", "", "Field holding the source identifier")
	flags.StringVar(&options.Match, "match", review.MatchURL, "How the identifier matches a source: url|path|doi|pmcid")
	flags.StringVar(&options.DecisionField, "decision-field", "", "Field holding the screening decision")
	flags.StringSliceVar(&options.KeepValues, "keep-values", review.DefaultKeepValues, "Decision values meaning keep/relevant")
	flags.StringVar(&options.DecidedByField, "decided-by-field", "", "Field copied to the label's decided_by (default screening_stage or decided_by)")
	for _, name := range []string{"from", "id-field", "decision-field"} {
		_ = command.MarkFlagRequired(name)
	}
	return command
}

func runReviewImportLabels(cmd *cobra.Command, options *reviewImportLabelsOptions, topicSlug string) error {
	topic, err := resolveReviewTopic(cmd, "review import-labels", topicSlug)
	if err != nil {
		return err
	}
	report, err := review.ImportLabels(topic.root, review.ImportOptions{
		From:           options.From,
		IDField:        options.IDField,
		Match:          options.Match,
		DecisionField:  options.DecisionField,
		KeepValues:     options.KeepValues,
		DecidedByField: options.DecidedByField,
	})
	if err != nil {
		return fmt.Errorf("review import-labels: %w", err)
	}
	out := cmd.OutOrStdout()
	lines := []string{
		fmt.Sprintf("rows: %d, matched: %d, unmatched: %d", report.Rows, report.Matched, report.Unmatched),
		fmt.Sprintf("labels written: %d positive, %d negative (%d already imported)", report.Positive, report.Negative, report.Duplicates),
		"origin: " + report.Origin,
	}
	if len(report.UnmatchedIDs) > 0 {
		lines = append(lines, "unmatched ids:")
		for _, id := range report.UnmatchedIDs {
			if id == "" {
				id = "(empty)"
			}
			lines = append(lines, "  "+id)
		}
	}
	for _, line := range lines {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	return nil
}

type reviewCalibrateOptions struct {
	Write  bool
	Format string
}

func newReviewCalibrateCommand() *cobra.Command {
	options := &reviewCalibrateOptions{}
	command := &cobra.Command{
		Use:   "calibrate <topic>",
		Short: "Measure precision and recall per purpose from labels and choose thresholds on a dev/holdout split",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReviewCalibrate(cmd, options, args[0])
		},
	}
	command.Flags().BoolVar(&options.Write, "write", false, "Store the chosen thresholds in topic.yaml and the calibration in .decisions/calibration.json")
	command.Flags().StringVar(&options.Format, "format", "table", "Output format (table|json)")
	return command
}

func runReviewCalibrate(cmd *cobra.Command, options *reviewCalibrateOptions, topicSlug string) error {
	format, err := reviewFormat(options.Format)
	if err != nil {
		return err
	}
	topic, err := resolveReviewTopic(cmd, "review calibrate", topicSlug)
	if err != nil {
		return err
	}
	cfg, err := loadCLIConfig()
	if err != nil {
		return fmt.Errorf("review calibrate: %w", err)
	}
	settings, err := contract.LoadSettings(topic.root)
	if err != nil {
		return fmt.Errorf("review calibrate: %w", err)
	}
	var active *contract.Contract
	if settings.Accepted() {
		normalized := settings.Contract.Normalized()
		active = &normalized
	}
	report, err := review.Calibrate(topic.root, review.CalibrateOptions{
		ContractHash: active.Hash(),
		Thresholds:   decisions.Thresholds(nil).Merge(cfg.Decisions.Thresholds).Merge(settings.Decisions.Thresholds),
	})
	if err != nil {
		return fmt.Errorf("review calibrate: %w", err)
	}

	out := cmd.OutOrStdout()
	if format == "json" {
		err = writeReviewJSON(out, report)
	} else {
		err = review.RenderCalibration(out, report)
	}
	if err != nil || !options.Write {
		return err
	}
	switch err := review.WriteCalibration(topic.root, report); {
	case errors.Is(err, review.ErrNotEnoughLabels):
		_, err = fmt.Fprintln(cmd.ErrOrStderr(), "not enough labels: nothing written")
		return err
	case err != nil:
		return fmt.Errorf("review calibrate: %w", err)
	}
	_, err = fmt.Fprintln(cmd.ErrOrStderr(), "calibration written to topic.yaml decisions.thresholds and .decisions/calibration.json")
	return err
}

func resolveReviewTopic(cmd *cobra.Command, action, topicSlug string) (reviewTopic, error) {
	topicSlug = strings.TrimSpace(topicSlug)
	if topicSlug == "" {
		return reviewTopic{}, fmt.Errorf("%s: topic slug is required", action)
	}
	cwd, err := reviewGetwd()
	if err != nil {
		return reviewTopic{}, fmt.Errorf("%s: resolve cwd: %w", action, err)
	}
	resolved, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{CWD: cwd, Topic: topicSlug, Vault: commandVaultValue(cmd, "")})
	if err != nil {
		return reviewTopic{}, fmt.Errorf("%s: %w", action, err)
	}
	return reviewTopic{vaultPath: resolved.VaultPath, root: resolved.TopicPath, slug: resolved.TopicSlug}, nil
}

func reviewFormat(value string) (string, error) {
	switch format := strings.ToLower(strings.TrimSpace(value)); format {
	case "", "table":
		return "table", nil
	case "json":
		return "json", nil
	default:
		return "", fmt.Errorf(`invalid --format %q: expected "table" or "json"`, value)
	}
}

func writeReviewJSON(w io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("review: encode json: %w", err)
	}
	_, err = w.Write(append(data, '\n'))
	return err
}
