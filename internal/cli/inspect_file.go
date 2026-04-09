package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/go-devstack/internal/vault"
)

func newInspectFileCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "file <path>",
		Short: "Lookup a file by its exact source path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				sourcePath := strings.TrimSpace(args[0])
				if sourcePath == "" {
					return inspectOutput{}, fmt.Errorf("a file path is required")
				}

				return toFileLookupOutput(context.Snapshot, sourcePath)
			})
		},
	}

	return command
}

func toFileLookupOutput(snapshot vault.VaultSnapshot, sourcePath string) (inspectOutput, error) {
	document, err := findInspectFileBySourcePath(snapshot, sourcePath)
	if err != nil {
		return inspectOutput{}, err
	}

	return createInspectDetailOutput(
		inspectDetailEntry{Field: "relative_path", Value: document.RelativePath},
		inspectDetailEntry{Field: "source_path", Value: inspectFrontmatterString(document, "source_path")},
		inspectDetailEntry{Field: "language", Value: inspectFrontmatterString(document, "language")},
		inspectDetailEntry{Field: "symbol_count", Value: inspectFrontmatterInt(document, "symbol_count")},
		inspectDetailEntry{Field: "symbols", Value: inspectSymbolsForFile(snapshot, inspectFrontmatterString(document, "source_path"))},
		inspectDetailEntry{Field: "afferent_coupling", Value: inspectFrontmatterInt(document, "afferent_coupling")},
		inspectDetailEntry{Field: "efferent_coupling", Value: inspectFrontmatterInt(document, "efferent_coupling")},
		inspectDetailEntry{Field: "instability", Value: inspectFrontmatterFloat(document, "instability")},
		inspectDetailEntry{Field: "is_orphan_file", Value: inspectFrontmatterBool(document, "is_orphan_file")},
		inspectDetailEntry{Field: "is_god_file", Value: inspectFrontmatterBool(document, "is_god_file")},
		inspectDetailEntry{Field: "has_circular_dependency", Value: inspectFrontmatterBool(document, "has_circular_dependency")},
		inspectDetailEntry{Field: "smells", Value: inspectFrontmatterStringArray(document, "smells")},
		inspectDetailEntry{Field: "outgoing_relations", Value: createInspectRelationRows(document.OutgoingRelations)},
		inspectDetailEntry{Field: "backlinks", Value: createInspectRelationRows(document.Backlinks)},
	), nil
}
