package cli

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/go-devstack/internal/vault"
)

type deadCodeRow struct {
	Kind       string
	Name       string
	SourcePath string
	SymbolKind string
	Reason     string
	Smells     []string
}

func newInspectDeadCodeCommand(options *inspectSharedOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "dead-code",
		Short: "List dead exports and orphan files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toDeadCodeOutput(context.Snapshot), nil
			})
		},
	}
}

func toDeadCodeOutput(snapshot vault.VaultSnapshot) inspectOutput {
	rows := make([]deadCodeRow, 0, len(snapshot.Symbols)+len(snapshot.Files))

	for _, document := range snapshot.Symbols {
		if !inspectFrontmatterBool(document, "is_dead_export") {
			continue
		}

		rows = append(rows, deadCodeRow{
			Kind:       "symbol",
			Name:       inspectFrontmatterString(document, "symbol_name"),
			SourcePath: inspectFrontmatterString(document, "source_path"),
			SymbolKind: inspectFrontmatterString(document, "symbol_kind"),
			Reason:     "dead-export",
			Smells:     inspectFrontmatterStringArray(document, "smells"),
		})
	}

	for _, document := range snapshot.Files {
		if !inspectFrontmatterBool(document, "is_orphan_file") {
			continue
		}

		sourcePath := inspectFrontmatterString(document, "source_path")
		rows = append(rows, deadCodeRow{
			Kind:       "file",
			Name:       sourcePath,
			SourcePath: sourcePath,
			Reason:     "orphan-file",
			Smells:     inspectFrontmatterStringArray(document, "smells"),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Kind != rows[j].Kind {
			return rows[i].Kind < rows[j].Kind
		}
		if rows[i].SourcePath != rows[j].SourcePath {
			return rows[i].SourcePath < rows[j].SourcePath
		}
		return rows[i].Name < rows[j].Name
	})

	return inspectOutput{
		Columns: []string{"kind", "name", "source_path", "symbol_kind", "reason", "smells"},
		Data:    deadCodeRowsToMaps(rows),
	}
}

func deadCodeRowsToMaps(rows []deadCodeRow) []map[string]any {
	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"kind":        row.Kind,
			"name":        row.Name,
			"source_path": row.SourcePath,
			"symbol_kind": row.SymbolKind,
			"reason":      row.Reason,
			"smells":      row.Smells,
		})
	}
	return data
}
