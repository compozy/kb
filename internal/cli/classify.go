package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/output"
	"github.com/compozy/kb/internal/session"
)

var runClassify = classify.Run

type classifyCommandOptions struct {
	OnlyMissing bool
	All         bool
	Format      string
	Flags       session.Flags
}

func newClassifyCommand() *cobra.Command {
	options := &classifyCommandOptions{Format: string(output.OutputFormatTable)}
	command := &cobra.Command{
		Use:   "classify <topic>",
		Short: "Judge every document of a topic on the classification facets",
		Long: "Classify writes genre, depth, relevance, quality, concepts, summary, entities and questions " +
			"(plus criterion and aliases on articles) as kb-owned frontmatter. It is incremental: documents " +
			"whose body, contract and question banks did not change are skipped without a call. Broken captures " +
			"and off-topic sources go to the recapture and remove review queues; nothing is quarantined.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClassifyCommand(cmd, options, args[0])
		},
	}
	flags := command.Flags()
	flags.BoolVar(&options.OnlyMissing, "only-missing", false, "Judge only documents without a state record or with a missing facet key")
	flags.BoolVar(&options.All, "all", false, "Re-judge every document, ignoring the state record")
	flags.StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	bindDecisionFlags(command, &options.Flags)
	return command
}

func runClassifyCommand(cmd *cobra.Command, options *classifyCommandOptions, topicSlug string) error {
	if options.OnlyMissing && options.All {
		return fmt.Errorf("classify: --only-missing and --all cannot be combined")
	}
	format, err := parseInspectOutputFormat(options.Format)
	if err != nil {
		return err
	}
	s, err := openSession(cmd, "kb classify", topicSlug, options.Flags)
	if err != nil {
		return err
	}

	report, runErr := runClassify(cmd.Context(), s, classify.Options{OnlyMissing: options.OnlyMissing, All: options.All})
	stderr := cmd.ErrOrStderr()
	for _, line := range report.Lines() {
		_, _ = fmt.Fprintln(stderr, line)
	}
	if report.NoVocabulary {
		_, _ = fmt.Fprintf(stderr, "no vocabulary: run `kb topic vocabulary %s --draft`\n", s.Topic.Slug)
	}
	s.WriteSummary(stderr)
	if runErr != nil {
		return fmt.Errorf("classify: %w", runErr)
	}
	return writeClassifyReport(cmd.OutOrStdout(), format, report)
}

func writeClassifyReport(w io.Writer, format output.OutputFormat, report classify.Report) error {
	if format == output.OutputFormatJSON {
		encoded, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("classify: encode report: %w", err)
		}
		if _, err := fmt.Fprintln(w, string(encoded)); err != nil {
			return fmt.Errorf("classify: write output: %w", err)
		}
		return nil
	}
	_, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  format,
		Columns: []string{"metric", "value"},
		Data:    report.Rows(),
	}))
	if err != nil {
		return fmt.Errorf("classify: write output: %w", err)
	}
	return nil
}
