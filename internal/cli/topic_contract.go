package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/output"
	"github.com/compozy/kb/internal/scope"
	"github.com/compozy/kb/internal/session"
	ktopic "github.com/compozy/kb/internal/topic"
)

var (
	runDraftContract  = scope.Draft
	runImportClaude   = scope.ImportClaude
	runAcceptContract = scope.Accept
)

type topicContractOptions struct {
	Draft        bool
	ImportClaude bool
	Accept       bool
	Force        bool
	Yes          bool
	Flags        session.Flags
}

func newTopicContractCommand() *cobra.Command {
	options := &topicContractOptions{}
	command := &cobra.Command{
		Use:   "contract <topic> --draft | --import-claude | --accept [--force] [--yes]",
		Short: "Draft, import or accept the selection contract of a topic",
		Long: "--draft asks the generation model for a contract drafted from the topic's own collection (raw/ folders, " +
			"a stratified sample of source titles, article titles, CLAUDE.md scope and curation/screening files) and " +
			"writes it to topic.yaml under contract_draft.\n" +
			"--import-claude copies the `## Selection contract` section of the topic CLAUDE.md into contract_draft " +
			"(no model needed).\n" +
			"--accept runs the self-check (kept lines against out_of_scope lines; conflicts block unless --force) and the " +
			"impact preview (the relevance and quality gates of the draft over the topic's sources, a stratified sample " +
			"of 500 when larger), then activates the draft only after you type the topic slug (or pass --yes).",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTopicContractCommand(cmd, options, args[0])
		},
	}
	command.Flags().BoolVar(&options.Draft, "draft", false, "Generate a contract draft into topic.yaml contract_draft")
	command.Flags().BoolVar(&options.ImportClaude, "import-claude", false, "Import the CLAUDE.md ## Selection contract section into contract_draft")
	command.Flags().BoolVar(&options.Accept, "accept", false, "Self-check, preview and activate contract_draft")
	command.Flags().BoolVar(&options.Force, "force", false, "With --accept: continue despite self-check conflicts")
	command.Flags().BoolVar(&options.Yes, "yes", false, "With --accept: activate without typing the topic slug")
	command.Flags().Float64Var(&options.Flags.BudgetUSD, "budget", 0, "Spend ceiling in US$ for this run (default [decisions].budget_usd)")
	return command
}

func runTopicContractCommand(cmd *cobra.Command, options *topicContractOptions, topicSlug string) error {
	actions := 0
	for _, set := range []bool{options.Draft, options.ImportClaude, options.Accept} {
		if set {
			actions++
		}
	}
	if actions != 1 {
		return errors.New("topic contract: pass exactly one of --draft, --import-claude or --accept")
	}
	if (options.Force || options.Yes) && !options.Accept {
		return errors.New("topic contract: --force and --yes only apply to --accept")
	}
	switch {
	case options.ImportClaude:
		return runTopicContractImport(cmd, topicSlug)
	case options.Draft:
		return runTopicContractDraft(cmd, options, topicSlug)
	default:
		return runTopicContractAccept(cmd, options, topicSlug)
	}
}

// runTopicContractImport needs no decision model: it resolves the topic and
// parses its CLAUDE.md.
func runTopicContractImport(cmd *cobra.Command, topicSlug string) error {
	vaultPath, err := resolveCommandVaultPath(cmd, sessionGetwd, "kb topic contract")
	if err != nil {
		return err
	}
	info, err := ktopic.Resolve(vaultPath, topicSlug)
	if err != nil {
		return fmt.Errorf("topic contract: %w", err)
	}
	draft, err := runImportClaude(info.RootPath)
	if err != nil {
		return fmt.Errorf("topic contract: %w", err)
	}
	if err := writeContractDraft(cmd.OutOrStdout(), draft); err != nil {
		return err
	}
	stderr := cmd.ErrOrStderr()
	if err := draft.Validate(); err != nil {
		_, _ = fmt.Fprintf(stderr, "warning: the imported draft does not validate yet (--accept will refuse it): %v\n", err)
	}
	_, _ = fmt.Fprintf(stderr, "imported CLAUDE.md ## Selection contract into topic.yaml contract_draft (not active); review it, then run: kb topic contract %s --accept\n", info.Slug)
	return nil
}

func runTopicContractDraft(cmd *cobra.Command, options *topicContractOptions, topicSlug string) error {
	s, err := openSession(cmd, "kb topic contract", topicSlug, options.Flags)
	if err != nil {
		return err
	}
	stderr := cmd.ErrOrStderr()
	draft, inputs, err := runDraftContract(commandContext(cmd), s, scope.DraftOptions{})
	_, _ = fmt.Fprintf(stderr, "draft inputs: %d sources in %d raw/ folders, %d sampled titles, %d articles, %d curation/screening files\n",
		inputs.Sources, len(inputs.Folders), len(inputs.SourceSample), len(inputs.Articles), len(inputs.Screening))
	for _, dropped := range inputs.Dropped {
		_, _ = fmt.Fprintf(stderr, "dropped from the generated draft: %s\n", dropped)
	}
	s.WriteSummary(stderr)
	if err != nil {
		return fmt.Errorf("topic contract: %w", err)
	}
	if err := writeContractDraft(cmd.OutOrStdout(), draft); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(stderr, "wrote topic.yaml contract_draft (not active); edit it, then run: kb topic contract %s --accept\n", s.Topic.Slug)
	return nil
}

func runTopicContractAccept(cmd *cobra.Command, options *topicContractOptions, topicSlug string) error {
	s, err := openSession(cmd, "kb topic contract", topicSlug, options.Flags)
	if err != nil {
		return err
	}
	stdout, stderr := cmd.OutOrStdout(), cmd.ErrOrStderr()
	var writeErr error
	previewShown := false
	keep := func(err error) {
		if writeErr == nil {
			writeErr = err
		}
	}
	stdin := bufio.NewReader(cmd.InOrStdin())
	report, err := runAcceptContract(commandContext(cmd), s, scope.AcceptOptions{
		Force: options.Force,
		Yes:   options.Yes,
		Confirm: func(prompt string) (string, error) {
			_, _ = fmt.Fprint(stderr, prompt)
			line, err := stdin.ReadString('\n')
			if !strings.HasSuffix(line, "\n") {
				_, _ = fmt.Fprintln(stderr)
			}
			if errors.Is(err, io.EOF) {
				return line, nil
			}
			return line, err
		},
		Preview: scope.PreviewOptions{OnEstimate: func(documents int, estimate float64) {
			_, _ = fmt.Fprintf(stderr, "impact preview: judging %d sources, estimated US$ %.4f (US$ %.4f per document)\n", documents, estimate, scope.CostPerDocumentUSD)
		}},
		OnSelfCheck: func(conflicts scope.Conflicts) { keep(writeSelfCheck(stdout, conflicts)) },
		OnPreview: func(preview scope.PreviewReport) {
			keep(writePreview(stdout, preview))
			previewShown = true
			_, _ = fmt.Fprintf(stderr, "impact preview cost: US$ %.4f\n", preview.CostUSD)
		},
	})
	if report.Preview != nil && !previewShown {
		_, _ = fmt.Fprintf(stderr, "impact preview stopped after US$ %.4f\n", report.Preview.CostUSD)
	}
	s.WriteSummary(stderr)
	if writeErr != nil {
		return fmt.Errorf("topic contract: write output: %w", writeErr)
	}
	if err != nil {
		return fmt.Errorf("topic contract: %w", err)
	}
	_, _ = fmt.Fprintf(stderr, "contract activated for %s (topic.yaml contract, CLAUDE.md section re-rendered); next: kb classify %s\n", s.Topic.Slug, s.Topic.Slug)
	return nil
}

// writeContractDraft prints the draft as YAML followed by the CLAUDE.md
// section it renders to.
func writeContractDraft(w io.Writer, draft *contract.Contract) error {
	encoded, err := yaml.Marshal(map[string]any{"contract_draft": draft.Normalized()})
	if err != nil {
		return fmt.Errorf("topic contract: encode draft: %w", err)
	}
	if _, err := fmt.Fprintf(w, "%s\n%s", encoded, contract.RenderSection(draft)); err != nil {
		return fmt.Errorf("topic contract: write output: %w", err)
	}
	return nil
}

func writeSelfCheck(w io.Writer, conflicts scope.Conflicts) error {
	if _, err := fmt.Fprintf(w, "self-check: %d pairs, %d conflicts at P ≥ %.2f\n", conflicts.Pairs, len(conflicts.Conflicts), conflicts.Threshold); err != nil {
		return err
	}
	for _, undecided := range conflicts.Undecided {
		if _, err := fmt.Fprintf(w, "  undecided: %s\n", undecided); err != nil {
			return err
		}
	}
	if len(conflicts.Conflicts) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(conflicts.Conflicts))
	for _, conflict := range conflicts.Conflicts {
		rows = append(rows, map[string]any{
			"pair":     conflict.N,
			"p":        fmt.Sprintf("%.2f", conflict.P),
			"kept":     conflict.KeptField + ": " + conflict.Kept,
			"excluded": conflict.Excluded,
		})
	}
	_, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"pair", "p", "kept", "excluded"},
		Data:    rows,
	}))
	return err
}

func writePreview(w io.Writer, preview scope.PreviewReport) error {
	scopeLine := fmt.Sprintf("all %d sources", preview.Sources)
	if preview.Sampled {
		scopeLine = fmt.Sprintf("a stratified sample of %d of %d sources", preview.Judged, preview.Sources)
	}
	relevance := ""
	if preview.RelevanceOff {
		relevance = " (relevance off in topic.yaml: quality only)"
	}
	if _, err := fmt.Fprintf(w, "\nimpact preview over %s%s (quarantine: P(off_topic) ≥ %.2f or quality ≥ %.2f; review: P(off_topic) ≥ %.2f or quality ≥ %.2f)\n",
		scopeLine, relevance,
		preview.Thresholds["relevance_quarantine"], preview.Thresholds["quality_apply"],
		preview.Thresholds["relevance_review"], preview.Thresholds["quality_review"]); err != nil {
		return err
	}
	bands := []map[string]any{
		{"band": scope.BandKept, "documents": preview.Bands.Kept},
		{"band": scope.BandReview, "documents": preview.Bands.Review},
		{"band": "quarantine-would-be", "documents": preview.Bands.Quarantine},
		{"band": scope.BandUndecided, "documents": preview.Bands.Undecided},
	}
	if _, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{Format: output.OutputFormatTable, Columns: []string{"band", "documents"}, Data: bands})); err != nil {
		return err
	}

	folders := make([]map[string]any, 0, len(preview.Folders))
	for _, folder := range preview.Folders {
		folders = append(folders, map[string]any{
			"folder": folder.Folder, "judged": folder.Judged, "kept": folder.Kept,
			"review": folder.Review, "quarantine": folder.Quarantine, "undecided": folder.Undecided,
		})
	}
	if _, err := fmt.Fprintf(w, "\nper raw/ folder:\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"folder", "judged", "kept", "review", "quarantine", "undecided"},
		Data:    folders,
	})); err != nil {
		return err
	}

	top := make([]map[string]any, 0, len(preview.Top))
	for _, doc := range preview.Top {
		band := doc.Band
		if doc.Reason != "" {
			band += " (" + doc.Reason + ")"
		}
		top = append(top, map[string]any{"p_off_topic": fmt.Sprintf("%.2f", doc.POffTopic), "band": band, "title": doc.Title, "path": doc.Path})
	}
	if _, err := fmt.Fprintf(w, "\nhighest P(off_topic) (%d):\n", len(top)); err != nil {
		return err
	}
	if _, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"p_off_topic", "band", "title", "path"},
		Data:    top,
	})); err != nil {
		return err
	}

	if labels := preview.Labels; labels != nil {
		if _, err := fmt.Fprintf(w, "\nagainst %d relevance labels (%d negative): would-be quarantine precision %.2f (%d/%d), recall %.2f (%d/%d)\n",
			labels.Labels, labels.Negatives, labels.Precision, labels.TruePositive, labels.Predicted,
			labels.Recall, labels.TruePositive, labels.Negatives); err != nil {
			return err
		}
	}
	if len(preview.Undecided) > 0 {
		reasons := make([]string, 0, len(preview.Undecided))
		for reason, count := range preview.Undecided {
			reasons = append(reasons, fmt.Sprintf("%s ×%d", reason, count))
		}
		slices.Sort(reasons)
		if _, err := fmt.Fprintf(w, "\nundecided answers: %s\n", strings.Join(reasons, ", ")); err != nil {
			return err
		}
	}
	return nil
}
