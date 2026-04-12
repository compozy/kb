// Package metrics computes file-level, symbol-level, and directory-level code metrics including complexity, coupling, blast radius, and code smell detection.
package metrics

import (
	"math"
	"path"
	"sort"
	"strings"

	"github.com/user/kb/internal/models"
)

const (
	centralityDamping    = 0.85
	centralityIterations = 10
)

var dependencyRelationTypes = map[models.RelationType]struct{}{
	models.RelCalls:      {},
	models.RelReferences: {},
}

// ComputeMetrics derives symbol, file, directory, and circular dependency metrics for a graph snapshot.
func ComputeMetrics(graph models.GraphSnapshot) models.MetricsResult {
	result := models.MetricsResult{
		CircularDependencies: [][]string{},
		Directories:          map[string]models.DirectoryMetrics{},
		Files:                map[string]models.FileMetrics{},
		Symbols:              map[string]models.SymbolMetrics{},
	}

	symbolsByID := make(map[string]models.SymbolNode, len(graph.Symbols))
	filesByID := make(map[string]models.GraphFile, len(graph.Files))
	fileIDs := make(map[string]struct{}, len(graph.Files))
	symbolsByFilePath := make(map[string][]models.SymbolNode, len(graph.Files))
	incomingDependencyEdgesByTarget := make(map[string][]models.RelationEdge)
	outgoingDependencyEdgesBySource := make(map[string][]models.RelationEdge)
	reverseDependencySourcesByTarget := make(map[string]map[string]struct{})
	fileDependencies := make(map[string]map[string]struct{})
	reverseFileDependencies := make(map[string]map[string]struct{})

	for _, symbol := range graph.Symbols {
		symbolsByID[symbol.ID] = symbol
		symbolsByFilePath[symbol.FilePath] = append(symbolsByFilePath[symbol.FilePath], symbol)
	}

	for _, file := range graph.Files {
		filesByID[file.ID] = file
		fileIDs[file.ID] = struct{}{}
	}

	for _, relation := range graph.Relations {
		sourceFileID, hasSourceFile := getNodeFileID(relation.FromID, fileIDs, symbolsByID)
		targetFileID, hasTargetFile := getNodeFileID(relation.ToID, fileIDs, symbolsByID)

		if relation.Type == models.RelImports {
			if hasSourceFile && hasTargetFile && sourceFileID != targetFileID {
				addToSetMap(fileDependencies, sourceFileID, targetFileID)
				addToSetMap(reverseFileDependencies, targetFileID, sourceFileID)
			}
			continue
		}

		if !isDependencyRelationType(relation.Type) {
			continue
		}

		if _, exists := symbolsByID[relation.ToID]; exists {
			incomingDependencyEdgesByTarget[relation.ToID] = append(incomingDependencyEdgesByTarget[relation.ToID], relation)
			addToSetMap(reverseDependencySourcesByTarget, relation.ToID, relation.FromID)
		}

		outgoingDependencyEdgesBySource[relation.FromID] = append(outgoingDependencyEdgesBySource[relation.FromID], relation)

		if hasSourceFile && hasTargetFile && sourceFileID != targetFileID {
			addToSetMap(fileDependencies, sourceFileID, targetFileID)
			addToSetMap(reverseFileDependencies, targetFileID, sourceFileID)
		}
	}

	entryPointFileIDs := make(map[string]struct{})
	for _, file := range graph.Files {
		if isEntryPointFile(file, symbolsByFilePath[file.FilePath]) {
			entryPointFileIDs[file.ID] = struct{}{}
		}
	}

	circularDependencyGroups := detectCircularDependencyGroups(graph.Files, graph.Relations)
	result.CircularDependencies = circularDependencyGroups

	filesInCircularDependencies := make(map[string]struct{})
	for _, group := range circularDependencyGroups {
		for _, filePath := range group {
			filesInCircularDependencies[createFileID(filePath)] = struct{}{}
		}
	}

	dependencySourcesByTarget := make(map[string][]string, len(reverseDependencySourcesByTarget))
	for targetID, sourceIDs := range reverseDependencySourcesByTarget {
		dependencySourcesByTarget[targetID] = sortedSetValues(sourceIDs)
	}

	centralityBySymbolID := computeApproxCentrality(graph, symbolsByID)
	for _, symbol := range graph.Symbols {
		incomingDependencyEdges := incomingDependencyEdgesByTarget[symbol.ID]
		incomingDependencySources := dependencySourcesByTarget[symbol.ID]
		externalReferenceCount := 0

		for _, relation := range incomingDependencyEdges {
			sourceFileID, hasSourceFile := getNodeFileID(relation.FromID, fileIDs, symbolsByID)
			if hasSourceFile && sourceFileID != createFileID(symbol.FilePath) {
				externalReferenceCount++
			}
		}

		visitedDependents := make(map[string]struct{})
		queue := append([]string(nil), incomingDependencySources...)
		for index := 0; index < len(queue); index++ {
			currentID := queue[index]
			if currentID == "" {
				continue
			}
			if _, seen := visitedDependents[currentID]; seen {
				continue
			}

			visitedDependents[currentID] = struct{}{}

			if _, isSymbol := symbolsByID[currentID]; !isSymbol {
				continue
			}

			for _, parentID := range dependencySourcesByTarget[currentID] {
				if _, seen := visitedDependents[parentID]; !seen {
					queue = append(queue, parentID)
				}
			}
		}

		dependencyTargetsByFile := make(map[string]map[string]struct{})
		for _, relation := range outgoingDependencyEdgesBySource[symbol.ID] {
			targetSymbol, exists := symbolsByID[relation.ToID]
			if !exists || targetSymbol.ID == symbol.ID {
				continue
			}

			addToSetMap(dependencyTargetsByFile, targetSymbol.FilePath, targetSymbol.ID)
		}

		ownFileDependencyCount := len(dependencyTargetsByFile[symbol.FilePath])
		strongestForeignDependencyCount := 0
		for filePath, targetIDs := range dependencyTargetsByFile {
			if filePath == symbol.FilePath {
				continue
			}
			if len(targetIDs) > strongestForeignDependencyCount {
				strongestForeignDependencyCount = len(targetIDs)
			}
		}

		loc := symbol.EndLine - symbol.StartLine + 1
		if loc < 0 {
			loc = 0
		}

		isLongFunction := isFunctionLike(symbol.SymbolKind) && (loc > 50 || symbol.CyclomaticComplexity > 10)
		isDeadExport := symbol.Exported &&
			symbol.Name != "main" &&
			!containsKey(entryPointFileIDs, createFileID(symbol.FilePath)) &&
			externalReferenceCount == 0
		hasFeatureEnvy := isFunctionLike(symbol.SymbolKind) &&
			strongestForeignDependencyCount > ownFileDependencyCount &&
			strongestForeignDependencyCount > 0
		centrality := centralityBySymbolID[symbol.ID]
		smells := collectSortedSmells(
			conditionalSmell(isDeadExport, "dead-export"),
			conditionalSmell(isLongFunction, "long-function"),
			conditionalSmell(len(visitedDependents) > 20, "high-blast-radius"),
			conditionalSmell(centrality > 0.1, "bottleneck"),
			conditionalSmell(hasFeatureEnvy, "feature-envy"),
		)

		result.Symbols[symbol.ID] = models.SymbolMetrics{
			BlastRadius:            len(visitedDependents),
			Centrality:             centrality,
			DirectDependents:       len(incomingDependencySources),
			ExternalReferenceCount: externalReferenceCount,
			IsDeadExport:           isDeadExport,
			IsLongFunction:         isLongFunction,
			LOC:                    loc,
			Smells:                 smells,
		}
	}

	for _, file := range graph.Files {
		afferentCoupling := len(reverseFileDependencies[file.ID])
		efferentCoupling := len(fileDependencies[file.ID])
		_, isEntryPoint := entryPointFileIDs[file.ID]
		_, hasCircularDependency := filesInCircularDependencies[file.ID]
		isOrphanFile := !isEntryPoint && afferentCoupling == 0
		isGodFile := len(file.SymbolIDs) > 15 || efferentCoupling > 10

		result.Files[file.ID] = models.FileMetrics{
			AfferentCoupling:      afferentCoupling,
			EfferentCoupling:      efferentCoupling,
			HasCircularDependency: hasCircularDependency,
			Instability:           computeInstability(afferentCoupling, efferentCoupling),
			IsEntryPoint:          isEntryPoint,
			IsGodFile:             isGodFile,
			IsOrphanFile:          isOrphanFile,
			Smells: collectSortedSmells(
				conditionalSmell(isGodFile, "god-file"),
				conditionalSmell(isOrphanFile, "orphan-file"),
			),
		}
	}

	directoryOutgoingFiles := make(map[string]map[string]struct{})
	directoryIncomingFiles := make(map[string]map[string]struct{})
	directoryPaths := make(map[string]struct{})

	for _, file := range graph.Files {
		directoryPaths[getDirectoryPath(file.FilePath)] = struct{}{}
	}

	for sourceFileID, targetFileIDs := range fileDependencies {
		sourceFile, exists := filesByID[sourceFileID]
		if !exists {
			continue
		}

		sourceDirectoryPath := getDirectoryPath(sourceFile.FilePath)
		for targetFileID := range targetFileIDs {
			targetFile, exists := filesByID[targetFileID]
			if !exists {
				continue
			}

			targetDirectoryPath := getDirectoryPath(targetFile.FilePath)
			if sourceDirectoryPath == targetDirectoryPath {
				continue
			}

			addToSetMap(directoryOutgoingFiles, sourceDirectoryPath, targetFileID)
			addToSetMap(directoryIncomingFiles, targetDirectoryPath, sourceFileID)
		}
	}

	for _, directoryPath := range sortedSetValues(directoryPaths) {
		afferentCoupling := len(directoryIncomingFiles[directoryPath])
		efferentCoupling := len(directoryOutgoingFiles[directoryPath])
		result.Directories[directoryPath] = models.DirectoryMetrics{
			AfferentCoupling: afferentCoupling,
			EfferentCoupling: efferentCoupling,
			Instability:      computeInstability(afferentCoupling, efferentCoupling),
		}
	}

	return result
}

func computeApproxCentrality(
	graph models.GraphSnapshot,
	symbolsByID map[string]models.SymbolNode,
) map[string]float64 {
	symbolIDs := make([]string, 0, len(graph.Symbols))
	for _, symbol := range graph.Symbols {
		symbolIDs = append(symbolIDs, symbol.ID)
	}
	sort.Strings(symbolIDs)

	if len(symbolIDs) == 0 {
		return map[string]float64{}
	}

	symbolIndexByID := make(map[string]int, len(symbolIDs))
	incomingEdges := make([][]int, len(symbolIDs))
	outDegree := make([]int, len(symbolIDs))
	seenEdges := make(map[string]struct{})

	for index, symbolID := range symbolIDs {
		symbolIndexByID[symbolID] = index
	}

	for _, relation := range graph.Relations {
		if !isDependencyRelationType(relation.Type) {
			continue
		}
		if _, exists := symbolsByID[relation.FromID]; !exists {
			continue
		}
		if _, exists := symbolsByID[relation.ToID]; !exists {
			continue
		}

		fromIndex, hasFrom := symbolIndexByID[relation.FromID]
		toIndex, hasTo := symbolIndexByID[relation.ToID]
		if !hasFrom || !hasTo || fromIndex == toIndex {
			continue
		}

		edgeKey := relation.FromID + "\x00" + relation.ToID
		if _, seen := seenEdges[edgeKey]; seen {
			continue
		}

		seenEdges[edgeKey] = struct{}{}
		incomingEdges[toIndex] = append(incomingEdges[toIndex], fromIndex)
		outDegree[fromIndex]++
	}

	scores := make([]float64, len(symbolIDs))
	initialScore := 1 / float64(len(symbolIDs))
	for index := range scores {
		scores[index] = initialScore
	}

	for iteration := 0; iteration < centralityIterations; iteration++ {
		nextScores := make([]float64, len(symbolIDs))
		danglingMass := 0.0

		for index, degree := range outDegree {
			if degree == 0 {
				danglingMass += scores[index]
			}
		}

		danglingContribution := (centralityDamping * danglingMass) / float64(len(symbolIDs))
		baseScore := (1 - centralityDamping) / float64(len(symbolIDs))

		for targetIndex := range symbolIDs {
			incomingScore := 0.0
			for _, sourceIndex := range incomingEdges[targetIndex] {
				degree := outDegree[sourceIndex]
				if degree == 0 {
					degree = 1
				}
				incomingScore += scores[sourceIndex] / float64(degree)
			}

			nextScores[targetIndex] = baseScore + danglingContribution + centralityDamping*incomingScore
		}

		scores = nextScores
	}

	maxCentrality := 0.0
	for _, score := range scores {
		if score > maxCentrality {
			maxCentrality = score
		}
	}

	centralityBySymbolID := make(map[string]float64, len(symbolIDs))
	if maxCentrality == 0 {
		for _, symbolID := range symbolIDs {
			centralityBySymbolID[symbolID] = 0
		}
		return centralityBySymbolID
	}

	for index, symbolID := range symbolIDs {
		centralityBySymbolID[symbolID] = toMetricRatio(scores[index] / maxCentrality)
	}

	return centralityBySymbolID
}

func detectCircularDependencyGroups(files []models.GraphFile, relations []models.RelationEdge) [][]string {
	return FindCircularDependencyGroups(buildFileImportAdjacency(files, relations))
}

// FindCircularDependencyGroups returns deterministic strongly connected file
// groups for the provided adjacency list. Groups with a single node are
// ignored, so self-referential files are not treated as circular components.
func FindCircularDependencyGroups(adjacency map[string][]string) [][]string {
	nodeIDs := make([]string, 0, len(adjacency))
	for nodeID := range adjacency {
		nodeIDs = append(nodeIDs, nodeID)
	}
	sort.Strings(nodeIDs)

	if len(nodeIDs) == 0 {
		return [][]string{}
	}

	normalizedAdjacency := make(map[string][]string, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		neighbors := uniqueSortedStrings(adjacency[nodeID])
		filtered := make([]string, 0, len(neighbors))
		for _, neighborID := range neighbors {
			if neighborID == nodeID {
				continue
			}
			if _, exists := adjacency[neighborID]; !exists {
				continue
			}
			filtered = append(filtered, neighborID)
		}
		normalizedAdjacency[nodeID] = filtered
	}

	index := 0
	indexByNode := make(map[string]int, len(nodeIDs))
	lowLinkByNode := make(map[string]int, len(nodeIDs))
	onStack := make(map[string]bool, len(nodeIDs))
	stack := make([]string, 0, len(nodeIDs))
	components := make([][]string, 0)

	var strongConnect func(string)
	strongConnect = func(nodeID string) {
		indexByNode[nodeID] = index
		lowLinkByNode[nodeID] = index
		index++

		stack = append(stack, nodeID)
		onStack[nodeID] = true

		for _, neighborID := range normalizedAdjacency[nodeID] {
			if _, visited := indexByNode[neighborID]; !visited {
				strongConnect(neighborID)
				if lowLinkByNode[neighborID] < lowLinkByNode[nodeID] {
					lowLinkByNode[nodeID] = lowLinkByNode[neighborID]
				}
				continue
			}

			if onStack[neighborID] && indexByNode[neighborID] < lowLinkByNode[nodeID] {
				lowLinkByNode[nodeID] = indexByNode[neighborID]
			}
		}

		if lowLinkByNode[nodeID] != indexByNode[nodeID] {
			return
		}

		component := make([]string, 0)
		for {
			lastIndex := len(stack) - 1
			memberID := stack[lastIndex]
			stack = stack[:lastIndex]
			onStack[memberID] = false
			component = append(component, memberID)
			if memberID == nodeID {
				break
			}
		}

		if len(component) <= 1 {
			return
		}

		sort.Strings(component)
		components = append(components, component)
	}

	for _, nodeID := range nodeIDs {
		if _, visited := indexByNode[nodeID]; visited {
			continue
		}
		strongConnect(nodeID)
	}

	sort.Slice(components, func(i, j int) bool {
		return strings.Join(components[i], "\x00") < strings.Join(components[j], "\x00")
	})

	return components
}

func buildFileImportAdjacency(files []models.GraphFile, relations []models.RelationEdge) map[string][]string {
	adjacency := make(map[string][]string, len(files))

	for _, file := range files {
		adjacency[file.FilePath] = []string{}
	}

	for _, relation := range relations {
		if relation.Type != models.RelImports {
			continue
		}

		sourcePath := strings.TrimPrefix(relation.FromID, "file:")
		targetPath := strings.TrimPrefix(relation.ToID, "file:")
		if sourcePath == targetPath {
			continue
		}
		if _, exists := adjacency[sourcePath]; !exists {
			continue
		}
		if _, exists := adjacency[targetPath]; !exists {
			continue
		}

		adjacency[sourcePath] = append(adjacency[sourcePath], targetPath)
	}

	return adjacency
}

func uniqueSortedStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		unique[value] = struct{}{}
	}

	ordered := make([]string, 0, len(unique))
	for value := range unique {
		ordered = append(ordered, value)
	}
	sort.Strings(ordered)

	return ordered
}

func getNodeFileID(
	nodeID string,
	fileIDs map[string]struct{},
	symbolsByID map[string]models.SymbolNode,
) (string, bool) {
	if _, exists := fileIDs[nodeID]; exists {
		return nodeID, true
	}

	symbol, exists := symbolsByID[nodeID]
	if !exists {
		return "", false
	}

	return createFileID(symbol.FilePath), true
}

func getDirectoryPath(filePath string) string {
	return path.Dir(filePath)
}

func isEntryPointFile(file models.GraphFile, symbols []models.SymbolNode) bool {
	for _, segment := range strings.Split(file.FilePath, "/") {
		if segment == "commands" {
			return true
		}
	}

	for _, symbol := range symbols {
		if symbol.Name == "main" {
			return true
		}
	}

	return false
}

func computeInstability(afferentCoupling, efferentCoupling int) float64 {
	totalCoupling := afferentCoupling + efferentCoupling
	if totalCoupling == 0 {
		return 0
	}

	return toMetricRatio(float64(efferentCoupling) / float64(totalCoupling))
}

func toMetricRatio(value float64) float64 {
	return math.Round(value*10000) / 10000
}

func isDependencyRelationType(relationType models.RelationType) bool {
	_, exists := dependencyRelationTypes[relationType]
	return exists
}

func isFunctionLike(symbolKind string) bool {
	return symbolKind == "function" || symbolKind == "method"
}

func conditionalSmell(include bool, smell string) string {
	if !include {
		return ""
	}
	return smell
}

func collectSortedSmells(smells ...string) []string {
	filtered := make([]string, 0, len(smells))
	for _, smell := range smells {
		if smell != "" {
			filtered = append(filtered, smell)
		}
	}

	sort.Strings(filtered)
	return filtered
}

func createFileID(filePath string) string {
	return "file:" + filePath
}

func containsKey[K comparable](items map[K]struct{}, key K) bool {
	_, exists := items[key]
	return exists
}

func sortedSetValues[T ~string](items map[T]struct{}) []string {
	if len(items) == 0 {
		return []string{}
	}

	values := make([]string, 0, len(items))
	for value := range items {
		values = append(values, string(value))
	}
	sort.Strings(values)

	return values
}

func addToSetMap[K comparable, V comparable](items map[K]map[V]struct{}, key K, value V) {
	values := items[key]
	if values == nil {
		values = map[V]struct{}{}
		items[key] = values
	}

	values[value] = struct{}{}
}
