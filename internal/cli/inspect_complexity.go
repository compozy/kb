package cli

import (
	"fmt"
	"sort"

	"github.com/compozy/kb/internal/vault"
	"github.com/spf13/cobra"
)

type complexityRow struct {
	SymbolName           string
	SymbolKind           string
	SourcePath           string
	CyclomaticComplexity int
	LOC                  int
	BlastRadius          int
	Smells               []string
}

func newInspectComplexityCommand(options *inspectSharedOptions) *cobra.Command {
	top := 20

	command := &cobra.Command{
		Use:   "complexity",
		Short: "Rank functions by cyclomatic complexity",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if top < 1 {
				return fmt.Errorf("--top must be >= 1. received %d", top)
			}

			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toComplexityOutput(context.Snapshot, top), nil
			})
		},
	}

	command.Flags().IntVar(&top, "top", 20, "Maximum number of rows to return")
	return command
}

func toComplexityOutput(snapshot vault.VaultSnapshot, top int) inspectOutput {
	rows := make([]complexityRow, 0, len(snapshot.Symbols))
	for _, document := range snapshot.Symbols {
		if !isFunctionLikeDocument(document) {
			continue
		}

		rows = append(rows, complexityRow{
			SymbolName:           inspectFrontmatterString(document, "symbol_name"),
			SymbolKind:           inspectFrontmatterString(document, "symbol_kind"),
			SourcePath:           inspectFrontmatterString(document, "source_path"),
			CyclomaticComplexity: inspectFrontmatterInt(document, "cyclomatic_complexity"),
			LOC:                  inspectFrontmatterInt(document, "loc"),
			BlastRadius:          inspectFrontmatterInt(document, "blast_radius"),
			Smells:               inspectFrontmatterStringArray(document, "smells"),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CyclomaticComplexity != rows[j].CyclomaticComplexity {
			return rows[i].CyclomaticComplexity > rows[j].CyclomaticComplexity
		}
		if rows[i].SourcePath != rows[j].SourcePath {
			return rows[i].SourcePath < rows[j].SourcePath
		}
		return rows[i].SymbolName < rows[j].SymbolName
	})

	if top < len(rows) {
		rows = rows[:top]
	}

	return inspectOutput{
		Columns: []string{
			"symbol_name",
			"symbol_kind",
			"source_path",
			"cyclomatic_complexity",
			"loc",
			"blast_radius",
			"smells",
		},
		Data: complexityRowsToMaps(rows),
	}
}

func complexityRowsToMaps(rows []complexityRow) []map[string]any {
	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"symbol_name":           row.SymbolName,
			"symbol_kind":           row.SymbolKind,
			"source_path":           row.SourcePath,
			"cyclomatic_complexity": row.CyclomaticComplexity,
			"loc":                   row.LOC,
			"blast_radius":          row.BlastRadius,
			"smells":                row.Smells,
		})
	}
	return data
}
