package cli

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/go-devstack/internal/vault"
)

func newInspectCircularDepsCommand(options *inspectSharedOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "circular-deps",
		Short: "List detected circular dependency cycles",
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
	cycles := detectInspectCircularDependencyCycles(snapshot)
	if len(cycles) == 0 {
		return inspectOutput{
			Columns: []string{"message"},
			Data: []map[string]any{
				{"message": "no circular dependencies found"},
			},
		}
	}

	data := make([]map[string]any, 0, len(cycles))
	for index, cycle := range cycles {
		data = append(data, map[string]any{
			"cycle": index + 1,
			"files": cycle,
		})
	}

	return inspectOutput{
		Columns: []string{"cycle", "files"},
		Data:    data,
	}
}

func detectInspectCircularDependencyCycles(snapshot vault.VaultSnapshot) [][]string {
	adjacency := make(map[string]map[string]struct{}, len(snapshot.Files))
	filePaths := make([]string, 0, len(snapshot.Files))
	linkLookup := make(map[string]string, len(snapshot.Files)*2)

	for _, document := range snapshot.Files {
		sourcePath := inspectFrontmatterString(document, "source_path")
		if sourcePath == "" {
			continue
		}

		adjacency[sourcePath] = map[string]struct{}{}
		filePaths = append(filePaths, sourcePath)

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

			addInspectStringSet(adjacency, sourcePath, targetPath)
		}
	}

	sort.Strings(filePaths)
	cyclesByKey := make(map[string][]string)

	for _, startPath := range filePaths {
		pathStack := []string{startPath}
		visited := map[string]struct{}{startPath: {}}

		var visit func(string)
		visit = func(currentPath string) {
			for _, neighborPath := range sortedInspectStringSet(adjacency[currentPath]) {
				if neighborPath == startPath {
					cycle := canonicalizeInspectCycle(append([]string(nil), pathStack...))
					cyclesByKey[strings.Join(cycle, " -> ")] = cycle
					continue
				}
				if neighborPath < startPath {
					continue
				}
				if _, seen := visited[neighborPath]; seen {
					continue
				}

				visited[neighborPath] = struct{}{}
				pathStack = append(pathStack, neighborPath)
				visit(neighborPath)
				pathStack = pathStack[:len(pathStack)-1]
				delete(visited, neighborPath)
			}
		}

		visit(startPath)
	}

	if len(cyclesByKey) == 0 {
		return [][]string{}
	}

	keys := make([]string, 0, len(cyclesByKey))
	for key := range cyclesByKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	cycles := make([][]string, 0, len(keys))
	for _, key := range keys {
		cycles = append(cycles, cyclesByKey[key])
	}

	return cycles
}

func addInspectStringSet(target map[string]map[string]struct{}, key, value string) {
	values, exists := target[key]
	if !exists {
		values = map[string]struct{}{}
		target[key] = values
	}

	values[value] = struct{}{}
}

func sortedInspectStringSet(values map[string]struct{}) []string {
	if len(values) == 0 {
		return []string{}
	}

	ordered := make([]string, 0, len(values))
	for value := range values {
		ordered = append(ordered, value)
	}
	sort.Strings(ordered)
	return ordered
}

func canonicalizeInspectCycle(cycle []string) []string {
	if len(cycle) <= 1 {
		return cycle
	}

	bestCycle := append([]string(nil), cycle...)
	bestKey := strings.Join(bestCycle, "\x00")

	for index := 1; index < len(cycle); index++ {
		rotatedCycle := append([]string(nil), cycle[index:]...)
		rotatedCycle = append(rotatedCycle, cycle[:index]...)
		rotatedKey := strings.Join(rotatedCycle, "\x00")
		if rotatedKey < bestKey {
			bestCycle = rotatedCycle
			bestKey = rotatedKey
		}
	}

	return bestCycle
}
