package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/kb/internal/vault"
)

func newInspectDepsCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "deps <name-or-path>",
		Short: "Show outgoing relations for a matched symbol or file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				query := strings.TrimSpace(args[0])
				if query == "" {
					return inspectOutput{}, fmt.Errorf("a symbol name or file path is required")
				}

				return toDependencyOutput(context.Snapshot, query)
			})
		},
	}

	return command
}

func toDependencyOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
	document, _, err := resolveInspectEntity(snapshot, query)
	if err != nil {
		return inspectOutput{}, err
	}

	return inspectOutput{
		Columns: []string{"target_path", "type", "confidence"},
		Data:    createInspectRelationRows(document.OutgoingRelations),
	}, nil
}
