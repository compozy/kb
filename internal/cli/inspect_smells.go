package cli

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/go-devstack/internal/vault"
)

type smellRow struct {
	Kind       string
	Name       string
	SourcePath string
	SymbolKind string
	Smells     []string
}

func newInspectSmellsCommand(options *inspectSharedOptions) *cobra.Command {
	var smellType string

	command := &cobra.Command{
		Use:   "smells",
		Short: "List symbols and files with detected smells",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toSmellOutput(context.Snapshot, smellType), nil
			})
		},
	}

	command.Flags().StringVar(&smellType, "type", "", "Filter to a specific smell type")
	return command
}

func toSmellOutput(snapshot vault.VaultSnapshot, smellType string) inspectOutput {
	filter := strings.ToLower(strings.TrimSpace(smellType))
	rows := make([]smellRow, 0, len(snapshot.Symbols)+len(snapshot.Files))

	for _, document := range snapshot.Symbols {
		smells := inspectFrontmatterStringArray(document, "smells")
		if !includeSmellRow(smells, filter) {
			continue
		}

		rows = append(rows, smellRow{
			Kind:       "symbol",
			Name:       inspectFrontmatterString(document, "symbol_name"),
			SourcePath: inspectFrontmatterString(document, "source_path"),
			SymbolKind: inspectFrontmatterString(document, "symbol_kind"),
			Smells:     smells,
		})
	}

	for _, document := range snapshot.Files {
		smells := inspectFrontmatterStringArray(document, "smells")
		if !includeSmellRow(smells, filter) {
			continue
		}

		sourcePath := inspectFrontmatterString(document, "source_path")
		rows = append(rows, smellRow{
			Kind:       "file",
			Name:       sourcePath,
			SourcePath: sourcePath,
			Smells:     smells,
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
		Columns: []string{"kind", "name", "source_path", "symbol_kind", "smells"},
		Data:    smellRowsToMaps(rows),
	}
}

func includeSmellRow(smells []string, filter string) bool {
	if len(smells) == 0 {
		return false
	}
	if filter == "" {
		return true
	}

	for _, smell := range smells {
		if strings.EqualFold(smell, filter) {
			return true
		}
	}

	return false
}

func smellRowsToMaps(rows []smellRow) []map[string]any {
	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"kind":        row.Kind,
			"name":        row.Name,
			"source_path": row.SourcePath,
			"symbol_kind": row.SymbolKind,
			"smells":      row.Smells,
		})
	}
	return data
}
