package cli

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/kb/internal/metrics"
	"github.com/user/kb/internal/vault"
)

func newInspectCircularDepsCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "circular-deps",
		Short: "List files that participate in circular dependencies",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspectCommand(cmd, options, func(context inspectContext) (inspectOutput, error) {
				return toCircularDepsOutput(context.Snapshot), nil
			})
		},
	}

	return command
}

func toCircularDepsOutput(snapshot vault.VaultSnapshot) inspectOutput {
	rows := circularDependencyRows(snapshot)
	if len(rows) == 0 {
		return inspectOutput{
			Columns: []string{"message"},
			Data: []map[string]any{
				{"message": "no circular dependencies found"},
			},
		}
	}

	return inspectOutput{
		Columns: []string{"source_path", "afferent_coupling", "efferent_coupling", "instability", "smells"},
		Data:    rows,
	}
}

func circularDependencyRows(snapshot vault.VaultSnapshot) []map[string]any {
	rows, hasPersistedFlags := circularDependencyRowsFromFlags(snapshot)
	if hasPersistedFlags {
		sortCircularDependencyRows(rows)
		return rows
	}

	rows = circularDependencyRowsFromFallback(snapshot)
	sortCircularDependencyRows(rows)
	return rows
}

func circularDependencyRowsFromFlags(snapshot vault.VaultSnapshot) ([]map[string]any, bool) {
	rows := make([]map[string]any, 0)
	hasPersistedFlags := false

	for _, document := range snapshot.Files {
		if _, exists := document.Frontmatter["has_circular_dependency"]; exists {
			hasPersistedFlags = true
		}
		if !inspectFrontmatterBool(document, "has_circular_dependency") {
			continue
		}

		rows = append(rows, toCircularDependencyRow(document))
	}

	return rows, hasPersistedFlags
}

func circularDependencyRowsFromFallback(snapshot vault.VaultSnapshot) []map[string]any {
	documentsBySourcePath := make(map[string]vault.VaultDocument, len(snapshot.Files))
	for _, document := range snapshot.Files {
		sourcePath := inspectFrontmatterString(document, "source_path")
		if sourcePath == "" {
			continue
		}
		documentsBySourcePath[sourcePath] = document
	}

	circularFiles := make(map[string]struct{})
	for _, group := range metrics.FindCircularDependencyGroups(buildInspectImportAdjacency(snapshot)) {
		for _, sourcePath := range group {
			circularFiles[sourcePath] = struct{}{}
		}
	}

	rows := make([]map[string]any, 0, len(circularFiles))
	for sourcePath := range circularFiles {
		document, exists := documentsBySourcePath[sourcePath]
		if !exists {
			continue
		}
		rows = append(rows, toCircularDependencyRow(document))
	}

	return rows
}

func buildInspectImportAdjacency(snapshot vault.VaultSnapshot) map[string][]string {
	adjacency := make(map[string][]string, len(snapshot.Files))
	linkLookup := make(map[string]string, len(snapshot.Files)*2)

	for _, document := range snapshot.Files {
		sourcePath := inspectFrontmatterString(document, "source_path")
		if sourcePath == "" {
			continue
		}

		adjacency[sourcePath] = []string{}

		relativeLinkPath := strings.TrimSuffix(document.RelativePath, ".md")
		linkLookup[relativeLinkPath] = sourcePath
		if snapshot.TopicSlug != "" {
			linkLookup[snapshot.TopicSlug+"/"+relativeLinkPath] = sourcePath
		}
	}

	for _, document := range snapshot.Files {
		sourcePath := inspectFrontmatterString(document, "source_path")
		if sourcePath == "" {
			continue
		}

		for _, relation := range document.OutgoingRelations {
			if relation.Type != "imports" {
				continue
			}

			targetPath, exists := linkLookup[relation.TargetPath]
			if !exists || targetPath == sourcePath {
				continue
			}

			adjacency[sourcePath] = append(adjacency[sourcePath], targetPath)
		}
	}

	return adjacency
}

func toCircularDependencyRow(document vault.VaultDocument) map[string]any {
	return map[string]any{
		"afferent_coupling": inspectFrontmatterInt(document, "afferent_coupling"),
		"efferent_coupling": inspectFrontmatterInt(document, "efferent_coupling"),
		"instability":       inspectFrontmatterFloat(document, "instability"),
		"smells":            inspectFrontmatterStringArray(document, "smells"),
		"source_path":       inspectFrontmatterString(document, "source_path"),
	}
}

func sortCircularDependencyRows(rows []map[string]any) {
	sort.Slice(rows, func(i, j int) bool {
		left, _ := rows[i]["source_path"].(string)
		right, _ := rows[j]["source_path"].(string)
		return left < right
	})
}
