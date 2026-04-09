package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/go-devstack/internal/vault"
)

type blastRadiusRow struct {
	SymbolName             string
	SourcePath             string
	BlastRadius            int
	Centrality             float64
	ExternalReferenceCount int
	Smells                 []string
}

func newInspectBlastRadiusCommand(options *inspectSharedOptions) *cobra.Command {
	minimum := 0
	top := 0

	command := &cobra.Command{
		Use:   "blast-radius",
		Short: "Rank symbols by blast radius",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if minimum < 0 {
				return fmt.Errorf("--min must be >= 0. received %d", minimum)
			}
			if top < 0 {
				return fmt.Errorf("--top must be >= 0. received %d", top)
			}

			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toBlastRadiusOutput(context.Snapshot, minimum, top), nil
			})
		},
	}

	command.Flags().IntVar(&minimum, "min", 0, "Minimum blast radius threshold")
	command.Flags().IntVar(&top, "top", 0, "Maximum number of rows to return (0 for all)")
	return command
}

func toBlastRadiusOutput(snapshot vault.VaultSnapshot, minimum, top int) inspectOutput {
	rows := make([]blastRadiusRow, 0, len(snapshot.Symbols))
	for _, document := range snapshot.Symbols {
		blastRadius := inspectFrontmatterInt(document, "blast_radius")
		if blastRadius < minimum {
			continue
		}

		rows = append(rows, blastRadiusRow{
			SymbolName:             inspectFrontmatterString(document, "symbol_name"),
			SourcePath:             inspectFrontmatterString(document, "source_path"),
			BlastRadius:            blastRadius,
			Centrality:             inspectFrontmatterFloat(document, "centrality"),
			ExternalReferenceCount: inspectFrontmatterInt(document, "external_reference_count"),
			Smells:                 inspectFrontmatterStringArray(document, "smells"),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].BlastRadius != rows[j].BlastRadius {
			return rows[i].BlastRadius > rows[j].BlastRadius
		}
		if rows[i].SourcePath != rows[j].SourcePath {
			return rows[i].SourcePath < rows[j].SourcePath
		}
		return rows[i].SymbolName < rows[j].SymbolName
	})

	if top > 0 && top < len(rows) {
		rows = rows[:top]
	}

	return inspectOutput{
		Columns: []string{
			"symbol_name",
			"source_path",
			"blast_radius",
			"centrality",
			"external_reference_count",
			"smells",
		},
		Data: blastRadiusRowsToMaps(rows),
	}
}

func blastRadiusRowsToMaps(rows []blastRadiusRow) []map[string]any {
	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"symbol_name":              row.SymbolName,
			"source_path":              row.SourcePath,
			"blast_radius":             row.BlastRadius,
			"centrality":               row.Centrality,
			"external_reference_count": row.ExternalReferenceCount,
			"smells":                   row.Smells,
		})
	}
	return data
}
