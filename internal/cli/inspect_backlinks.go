package cli

import (
	"fmt"
	"strings"

	"github.com/compozy/kb/internal/vault"
	"github.com/spf13/cobra"
)

func newInspectBacklinksCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "backlinks <name-or-path>",
		Short: "Show backlinks for a matched symbol or file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				query := strings.TrimSpace(args[0])
				if query == "" {
					return inspectOutput{}, fmt.Errorf("a symbol name or file path is required")
				}

				return toBacklinksOutput(context.Snapshot, query)
			})
		},
	}

	return command
}

func toBacklinksOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
	document, _, err := resolveInspectEntity(snapshot, query)
	if err != nil {
		return inspectOutput{}, err
	}

	return inspectOutput{
		Columns: []string{"target_path", "type", "confidence"},
		Data:    createInspectRelationRows(document.Backlinks),
	}, nil
}
