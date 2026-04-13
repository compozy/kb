package cli

import (
	"sort"

	"github.com/compozy/kb/internal/vault"
	"github.com/spf13/cobra"
)

type couplingRow struct {
	SourcePath            string
	AfferentCoupling      int
	EfferentCoupling      int
	Instability           float64
	HasCircularDependency bool
	Smells                []string
}

func newInspectCouplingCommand(options *inspectSharedOptions) *cobra.Command {
	var unstableOnly bool

	command := &cobra.Command{
		Use:   "coupling",
		Short: "Rank files by instability",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toCouplingOutput(context.Snapshot, unstableOnly), nil
			})
		},
	}

	command.Flags().BoolVar(&unstableOnly, "unstable", false, "Only show files with instability > 0.5")
	return command
}

func toCouplingOutput(snapshot vault.VaultSnapshot, unstableOnly bool) inspectOutput {
	rows := make([]couplingRow, 0, len(snapshot.Files))
	for _, document := range snapshot.Files {
		instability := inspectFrontmatterFloat(document, "instability")
		if unstableOnly && instability <= 0.5 {
			continue
		}

		rows = append(rows, couplingRow{
			SourcePath:            inspectFrontmatterString(document, "source_path"),
			AfferentCoupling:      inspectFrontmatterInt(document, "afferent_coupling"),
			EfferentCoupling:      inspectFrontmatterInt(document, "efferent_coupling"),
			Instability:           instability,
			HasCircularDependency: inspectFrontmatterBool(document, "has_circular_dependency"),
			Smells:                inspectFrontmatterStringArray(document, "smells"),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Instability != rows[j].Instability {
			return rows[i].Instability > rows[j].Instability
		}
		return rows[i].SourcePath < rows[j].SourcePath
	})

	return inspectOutput{
		Columns: []string{
			"source_path",
			"afferent_coupling",
			"efferent_coupling",
			"instability",
			"has_circular_dependency",
			"smells",
		},
		Data: couplingRowsToMaps(rows),
	}
}

func couplingRowsToMaps(rows []couplingRow) []map[string]any {
	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"source_path":             row.SourcePath,
			"afferent_coupling":       row.AfferentCoupling,
			"efferent_coupling":       row.EfferentCoupling,
			"instability":             row.Instability,
			"has_circular_dependency": row.HasCircularDependency,
			"smells":                  row.Smells,
		})
	}
	return data
}
