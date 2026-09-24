package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	kfind "github.com/compozy/kb/internal/find"
	"github.com/compozy/kb/internal/output"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

type findCommandOptions struct {
	flags     session.Flags
	limit     int
	kind      string
	concept   string
	relevance string
	minDepth  float64
	json      bool
	explain   bool
	facets    bool
}

func newFindCommand() *cobra.Command {
	options := &findCommandOptions{}

	command := &cobra.Command{
		Use:   "find <topic> \"<question>\" | find --facets <topic>",
		Short: "Find the documents of a topic that answer a question, ranked, with named exclusions",
		Long: "Proposes candidates by code (facet filters, BM25, concepts), asks the decision model whether each\n" +
			"candidate answers the question and ranks the survivors. --explain lists every dropped candidate with\n" +
			"its reason. --facets prints the topic's facet counts without calling any model.",
		Args: func(cmd *cobra.Command, args []string) error {
			if options.facets {
				return cobra.ExactArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.facets {
				return runFindFacets(cmd, args[0], options.json)
			}
			return runFind(cmd, args[0], args[1], options)
		},
	}

	flags := command.Flags()
	flags.IntVar(&options.limit, "limit", kfind.DefaultLimit, "Maximum number of results")
	flags.StringVar(&options.kind, "kind", "", "Keep only documents with this genre")
	flags.StringVar(&options.concept, "concept", "", "Keep only documents listing this concept")
	flags.Float64Var(&options.minDepth, "min-depth", 0, "Keep only documents with depth >= N (0-3)")
	flags.StringVar(&options.relevance, "relevance", "", "Keep only documents with this relevance role")
	flags.BoolVar(&options.json, "json", false, "Print JSON instead of a table")
	flags.BoolVar(&options.explain, "explain", false, "Also list every dropped candidate with its reason")
	flags.BoolVar(&options.facets, "facets", false, "Print facet counts (genre, relevance, depth, concepts) without any model call")
	flags.Float64Var(&options.flags.BudgetUSD, "budget", 0, "Spend ceiling in US$ for this run (default [decisions].budget_usd)")
	return command
}

func runFind(cmd *cobra.Command, topicSlug, question string, options *findCommandOptions) error {
	if strings.TrimSpace(question) == "" {
		return fmt.Errorf("kb find: a question is required")
	}
	if options.limit < 1 {
		return fmt.Errorf("kb find: --limit must be >= 1. received %d", options.limit)
	}
	s, err := openSession(cmd, "kb find", topicSlug, options.flags)
	if err != nil {
		return err
	}
	result, err := kfind.Run(commandContext(cmd), s, kfind.Query{
		Text:      question,
		Limit:     options.limit,
		Kind:      options.kind,
		Concept:   options.concept,
		Relevance: options.relevance,
		MinDepth:  options.minDepth,
		Explain:   options.explain,
	})
	if err != nil {
		s.WriteSummary(cmd.ErrOrStderr())
		return fmt.Errorf("kb find: %w", err)
	}
	if err := writeFindResult(cmd.OutOrStdout(), result, options.json, options.explain); err != nil {
		return fmt.Errorf("kb find: write output: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "find: %d candidates judged, %d kept, question kind %s\n", result.Candidates, len(result.Hits), findValueOr(result.QuestionKind, "unknown"))
	s.WriteSummary(cmd.ErrOrStderr())
	return nil
}

func writeFindResult(w io.Writer, result kfind.Result, asJSON, explain bool) error {
	if asJSON {
		if result.Hits == nil {
			result.Hits = []kfind.Hit{}
		}
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	rows := make([]map[string]any, 0, len(result.Hits))
	for _, hit := range result.Hits {
		rows = append(rows, map[string]any{
			"path":      hit.Path,
			"title":     hit.Title,
			"rank":      fmt.Sprintf("%.3f", hit.Rank),
			"p_answers": fmt.Sprintf("%.2f", hit.PAnswers),
			"kind":      hit.Kind,
			"concepts":  strings.Join(hit.Concepts, ", "),
		})
	}
	if _, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"path", "title", "rank", "p_answers", "kind", "concepts"},
		Data:    rows,
	})); err != nil {
		return err
	}
	if !explain {
		return nil
	}

	dropped := make([]map[string]any, 0, len(result.Dropped))
	for _, drop := range result.Dropped {
		p := ""
		if drop.PAnswers != nil {
			p = fmt.Sprintf("%.2f", *drop.PAnswers)
		}
		dropped = append(dropped, map[string]any{"path": drop.Path, "title": drop.Title, "reason": drop.Reason, "p_answers": p})
	}
	if _, err := fmt.Fprintf(w, "\ndropped (%d):\n", len(dropped)); err != nil {
		return err
	}
	_, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"path", "title", "reason", "p_answers"},
		Data:    dropped,
	}))
	return err
}

// runFindFacets prints the facet counts of a topic. It reads frontmatter only
// and never opens a decision session, so it runs without the decision model.
func runFindFacets(cmd *cobra.Command, topicSlug string, asJSON bool) error {
	vaultPath, err := resolveCommandVaultPath(cmd, sessionGetwd, "kb find")
	if err != nil {
		return err
	}
	info, err := topic.Resolve(vaultPath, topicSlug)
	if err != nil {
		return fmt.Errorf("kb find: %w", err)
	}
	settings, err := contract.LoadSettings(info.RootPath)
	if err != nil {
		return fmt.Errorf("kb find: %w", err)
	}
	loaded, err := corpus.Load(info.RootPath, corpus.LoadOptions{Exclude: settings.Decisions.Exclude})
	if err != nil {
		return fmt.Errorf("kb find: %w", err)
	}
	report := kfind.Facets(loaded.Documents())

	w := cmd.OutOrStdout()
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			return fmt.Errorf("kb find: write output: %w", err)
		}
		return nil
	}
	if err := writeFacets(w, report); err != nil {
		return fmt.Errorf("kb find: write output: %w", err)
	}
	return nil
}

func writeFacets(w io.Writer, report kfind.FacetReport) error {
	if _, err := fmt.Fprintf(w, "documents: %d (%d sources, %d articles), %d quarantined\n", report.Documents, report.Sources, report.Articles, report.Quarantined); err != nil {
		return err
	}
	sections := []struct {
		name   string
		counts []kfind.Count
	}{
		{"genre", report.Genre},
		{"relevance", report.Relevance},
		{"depth", report.Depth},
		{"concepts", report.Concepts},
	}
	for _, section := range sections {
		rows := make([]map[string]any, 0, len(section.counts))
		for _, count := range section.counts {
			rows = append(rows, map[string]any{section.name: count.Value, "documents": count.Count})
		}
		if _, err := fmt.Fprintf(w, "\n%s:\n", section.name); err != nil {
			return err
		}
		if _, err := io.WriteString(w, output.FormatOutput(output.FormatOptions{
			Format:  output.OutputFormatTable,
			Columns: []string{section.name, "documents"},
			Data:    rows,
		})); err != nil {
			return err
		}
	}
	return nil
}

func findValueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
