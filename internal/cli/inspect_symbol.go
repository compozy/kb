package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/vault"
	"github.com/spf13/cobra"
)

type inspectSymbolSummaryRow struct {
	Language   string
	SourcePath string
	StartLine  int
	Smells     []string
	SymbolKind string
	SymbolName string
}

func newInspectSymbolCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "symbol <name>",
		Short: "Lookup symbols by case-insensitive substring match",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				query := strings.TrimSpace(args[0])
				if query == "" {
					return inspectOutput{}, fmt.Errorf("a symbol name is required")
				}

				return toSymbolLookupOutput(context.Snapshot, query)
			})
		},
	}

	return command
}

func toSymbolLookupOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
	matches := vault.FindSymbolsByName(snapshot, query)
	if len(matches) == 0 {
		_, err := findSingleInspectSymbolMatch(snapshot, query)
		return inspectOutput{}, err
	}

	if len(matches) == 1 {
		return toSymbolDetailOutput(matches[0]), nil
	}

	return toSymbolSummaryOutput(matches), nil
}

func toSymbolSummaryOutput(documents []vault.VaultDocument) inspectOutput {
	rows := make([]inspectSymbolSummaryRow, 0, len(documents))
	for _, document := range documents {
		rows = append(rows, inspectSymbolSummaryRow{
			Language:   inspectFrontmatterString(document, "language"),
			SourcePath: inspectFrontmatterString(document, "source_path"),
			StartLine:  inspectFrontmatterInt(document, "start_line"),
			Smells:     inspectFrontmatterStringArray(document, "smells"),
			SymbolKind: inspectFrontmatterString(document, "symbol_kind"),
			SymbolName: inspectFrontmatterString(document, "symbol_name"),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].SymbolName != rows[j].SymbolName {
			return rows[i].SymbolName < rows[j].SymbolName
		}
		if rows[i].SourcePath != rows[j].SourcePath {
			return rows[i].SourcePath < rows[j].SourcePath
		}
		return rows[i].StartLine < rows[j].StartLine
	})

	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"symbol_name": row.SymbolName,
			"symbol_kind": row.SymbolKind,
			"source_path": row.SourcePath,
			"start_line":  row.StartLine,
			"language":    row.Language,
			"smells":      row.Smells,
		})
	}

	return inspectOutput{
		Columns: []string{"symbol_name", "symbol_kind", "source_path", "start_line", "language", "smells"},
		Data:    data,
	}
}

func toSymbolDetailOutput(document vault.VaultDocument) inspectOutput {
	signature := inspectSectionText(document, "Signature")
	if signature == "" {
		signature = "None"
	}

	return createInspectDetailOutput(
		inspectDetailEntry{Field: "relative_path", Value: document.RelativePath},
		inspectDetailEntry{Field: "symbol_name", Value: inspectFrontmatterString(document, "symbol_name")},
		inspectDetailEntry{Field: "symbol_kind", Value: inspectFrontmatterString(document, "symbol_kind")},
		inspectDetailEntry{Field: "source_path", Value: inspectFrontmatterString(document, "source_path")},
		inspectDetailEntry{Field: "language", Value: inspectFrontmatterString(document, "language")},
		inspectDetailEntry{Field: "exported", Value: inspectFrontmatterBool(document, "exported")},
		inspectDetailEntry{Field: "start_line", Value: inspectFrontmatterInt(document, "start_line")},
		inspectDetailEntry{Field: "end_line", Value: inspectFrontmatterInt(document, "end_line")},
		inspectDetailEntry{Field: "signature", Value: signature},
		inspectDetailEntry{Field: "loc", Value: inspectFrontmatterInt(document, "loc")},
		inspectDetailEntry{Field: "blast_radius", Value: inspectFrontmatterInt(document, "blast_radius")},
		inspectDetailEntry{Field: "centrality", Value: inspectFrontmatterFloat(document, "centrality")},
		inspectDetailEntry{Field: "cyclomatic_complexity", Value: inspectFrontmatterInt(document, "cyclomatic_complexity")},
		inspectDetailEntry{Field: "external_reference_count", Value: inspectFrontmatterInt(document, "external_reference_count")},
		inspectDetailEntry{Field: "is_dead_export", Value: inspectFrontmatterBool(document, "is_dead_export")},
		inspectDetailEntry{Field: "is_long_function", Value: inspectFrontmatterBool(document, "is_long_function")},
		inspectDetailEntry{Field: "smells", Value: inspectFrontmatterStringArray(document, "smells")},
		inspectDetailEntry{Field: "outgoing_relations", Value: createInspectRelationRows(document.OutgoingRelations)},
		inspectDetailEntry{Field: "backlinks", Value: createInspectRelationRows(document.Backlinks)},
	)
}
