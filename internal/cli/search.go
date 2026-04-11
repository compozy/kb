package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/go-devstack/internal/output"
	"github.com/user/go-devstack/internal/qmd"
	"github.com/user/go-devstack/internal/vault"
)

type searchCommandClient interface {
	Search(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error)
}

type searchCommandOptions struct {
	Collection string
	Format     string
	Full       bool
	All        bool
	Lex        bool
	Vec        bool
	Limit      int
	MinScore   float64
	Topic      string
	Vault      string
}

var newSearchClient = func() searchCommandClient {
	return qmd.NewClient()
}

var resolveSearchVaultQuery = vault.ResolveVaultQuery
var searchGetwd = os.Getwd

func newSearchCommand() *cobra.Command {
	options := &searchCommandOptions{
		Format: string(output.OutputFormatTable),
		Limit:  10,
	}

	command := &cobra.Command{
		Use:   "search <query>",
		Short: "Search a generated vault with QMD hybrid, lexical, or vector queries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args[0], options)
		},
	}

	flags := command.Flags()
	flags.BoolVar(&options.Lex, "lex", false, "Use BM25 keyword search only")
	flags.BoolVar(&options.Vec, "vec", false, "Use vector similarity search only")
	flags.IntVar(&options.Limit, "limit", 10, "Maximum number of results to return")
	flags.Float64Var(&options.MinScore, "min-score", 0, "Minimum score threshold for returned matches")
	flags.BoolVar(&options.Full, "full", false, "Show the full matched document content instead of snippets")
	flags.BoolVar(&options.All, "all", false, "Return all matches above the minimum score threshold")
	flags.StringVar(&options.Collection, "collection", "", "Use an explicit QMD collection name instead of deriving from the topic")
	flags.StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	flags.StringVar(&options.Topic, "topic", "", "Topic slug used when deriving the collection name")

	return command
}

func runSearchCommand(cmd *cobra.Command, query string, options *searchCommandOptions) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return fmt.Errorf("a search query is required")
	}

	mode, err := resolveSearchMode(options.Lex, options.Vec)
	if err != nil {
		return err
	}

	format, err := parseSearchOutputFormat(options.Format)
	if err != nil {
		return err
	}

	if options.Limit < 1 {
		return fmt.Errorf("--limit must be >= 1. received %d", options.Limit)
	}

	var minScore *float64
	if cmd.Flags().Changed("min-score") {
		if options.MinScore < 0 {
			return fmt.Errorf("--min-score must be >= 0. received %v", options.MinScore)
		}
		minScore = &options.MinScore
	}

	collection, err := resolveSearchCollection(cmd, options)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	results, err := newSearchClient().Search(ctx, qmd.SearchOptions{
		Query:      query,
		Mode:       mode,
		Limit:      options.Limit,
		All:        options.All,
		MinScore:   minScore,
		Full:       options.Full,
		Collection: collection,
	})
	if err != nil {
		return wrapQMDCommandError("search", err)
	}

	_, err = cmd.OutOrStdout().Write([]byte(output.FormatOutput(output.FormatOptions{
		Format:  format,
		Columns: []string{"path", "score", "preview"},
		Data:    searchResultsToRows(results),
	})))
	if err != nil {
		return fmt.Errorf("search: write output: %w", err)
	}

	return nil
}

func resolveSearchMode(lexical, vector bool) (qmd.SearchMode, error) {
	if lexical && vector {
		return "", fmt.Errorf("choose at most one search mode flag: --lex or --vec")
	}

	if lexical {
		return qmd.SearchModeLexical, nil
	}
	if vector {
		return qmd.SearchModeVector, nil
	}

	return qmd.SearchModeHybrid, nil
}

func resolveSearchCollection(cmd *cobra.Command, options *searchCommandOptions) (string, error) {
	if collection := strings.TrimSpace(options.Collection); collection != "" {
		return collection, nil
	}

	cwd, err := searchGetwd()
	if err != nil {
		return "", fmt.Errorf("resolve cwd: %w", err)
	}

	resolvedVault, err := resolveSearchVaultQuery(vault.VaultQueryOptions{
		CWD:   cwd,
		Topic: strings.TrimSpace(options.Topic),
		Vault: commandVaultValue(cmd, options.Vault),
	})
	if err != nil {
		return "", err
	}

	return resolvedVault.TopicSlug, nil
}

func parseSearchOutputFormat(value string) (output.OutputFormat, error) {
	switch output.OutputFormat(strings.ToLower(strings.TrimSpace(value))) {
	case "", output.OutputFormatTable:
		return output.OutputFormatTable, nil
	case output.OutputFormatJSON:
		return output.OutputFormatJSON, nil
	case output.OutputFormatTSV:
		return output.OutputFormatTSV, nil
	default:
		return "", fmt.Errorf(`invalid --format %q: expected one of "table", "json", "tsv"`, value)
	}
}

func searchResultsToRows(results []qmd.SearchResult) []map[string]any {
	rows := make([]map[string]any, 0, len(results))
	for _, result := range results {
		rows = append(rows, map[string]any{
			"path":    result.Path,
			"score":   result.Score,
			"preview": result.Snippet,
		})
	}
	return rows
}

func wrapQMDCommandError(command string, err error) error {
	if errors.Is(err, qmd.ErrQMDUnavailable) {
		return fmt.Errorf(
			"%s: QMD is not available to kb. Install it with `%s` and ensure `qmd` is on PATH",
			command,
			qmd.InstallCommand,
		)
	}

	return fmt.Errorf("%s: %w", command, err)
}
