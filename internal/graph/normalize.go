// Package graph normalizes and deduplicates code relations from parsed files into a clean dependency graph.
package graph

import (
	"sort"

	"github.com/compozy/kb/internal/models"
)

type relationKey struct {
	fromID     string
	toID       string
	relation   models.RelationType
	confidence models.RelationConfidence
}

// NormalizeGraph merges parsed files into a single deterministically ordered graph snapshot.
func NormalizeGraph(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot {
	files := make([]models.GraphFile, 0, len(parsedFiles))
	symbols := make([]models.SymbolNode, 0)
	externalNodes := make([]models.ExternalNode, 0)
	relations := make([]models.RelationEdge, 0)
	diagnostics := make([]models.StructuredDiagnostic, 0)

	for _, parsedFile := range parsedFiles {
		if isRenderableFile(parsedFile) {
			files = append(files, parsedFile.File)
		}

		symbols = append(symbols, parsedFile.Symbols...)
		externalNodes = append(externalNodes, parsedFile.ExternalNodes...)
		relations = append(relations, parsedFile.Relations...)
		diagnostics = append(diagnostics, parsedFile.Diagnostics...)
	}

	uniqueFiles := uniqueByKey(files, func(file models.GraphFile) string {
		return file.ID
	})
	uniqueSymbols := uniqueByKey(symbols, func(symbol models.SymbolNode) string {
		return symbol.ID
	})
	uniqueExternalNodes := uniqueByKey(externalNodes, func(node models.ExternalNode) string {
		return node.ID
	})
	uniqueRelationEdges := uniqueByKey(relations, func(relation models.RelationEdge) relationKey {
		return relationKey{
			fromID:     relation.FromID,
			toID:       relation.ToID,
			relation:   relation.Type,
			confidence: relation.Confidence,
		}
	})

	sort.Slice(uniqueFiles, func(left int, right int) bool {
		return uniqueFiles[left].ID < uniqueFiles[right].ID
	})
	sort.Slice(uniqueSymbols, func(left int, right int) bool {
		return uniqueSymbols[left].ID < uniqueSymbols[right].ID
	})
	sort.Slice(uniqueExternalNodes, func(left int, right int) bool {
		return uniqueExternalNodes[left].ID < uniqueExternalNodes[right].ID
	})
	sort.Slice(uniqueRelationEdges, func(left int, right int) bool {
		return compareRelationEdges(uniqueRelationEdges[left], uniqueRelationEdges[right]) < 0
	})
	sort.Slice(diagnostics, func(left int, right int) bool {
		return compareDiagnostics(diagnostics[left], diagnostics[right]) < 0
	})

	return models.GraphSnapshot{
		RootPath:      rootPath,
		Files:         attachSymbolIDs(uniqueFiles, uniqueSymbols),
		Symbols:       uniqueSymbols,
		ExternalNodes: uniqueExternalNodes,
		Relations:     uniqueRelationEdges,
		Diagnostics:   diagnostics,
	}
}

func isRenderableFile(parsedFile models.ParsedFile) bool {
	if len(parsedFile.Diagnostics) == 0 {
		return true
	}

	return len(parsedFile.Symbols) > 0 || len(parsedFile.ExternalNodes) > 0 || len(parsedFile.Relations) > 0
}

func attachSymbolIDs(files []models.GraphFile, symbols []models.SymbolNode) []models.GraphFile {
	symbolIDsByFilePath := make(map[string][]string, len(files))
	for _, symbol := range symbols {
		symbolIDsByFilePath[symbol.FilePath] = append(symbolIDsByFilePath[symbol.FilePath], symbol.ID)
	}

	attached := make([]models.GraphFile, 0, len(files))
	for _, file := range files {
		updatedFile := file

		symbolIDs := append([]string{}, symbolIDsByFilePath[file.FilePath]...)
		if symbolIDs == nil {
			symbolIDs = []string{}
		}

		sort.Strings(symbolIDs)
		updatedFile.SymbolIDs = symbolIDs
		attached = append(attached, updatedFile)
	}

	return attached
}

func compareRelationEdges(left models.RelationEdge, right models.RelationEdge) int {
	switch {
	case left.FromID != right.FromID:
		return compareStrings(left.FromID, right.FromID)
	case left.ToID != right.ToID:
		return compareStrings(left.ToID, right.ToID)
	case left.Type != right.Type:
		return compareStrings(string(left.Type), string(right.Type))
	default:
		return compareStrings(string(left.Confidence), string(right.Confidence))
	}
}

func compareDiagnostics(left models.StructuredDiagnostic, right models.StructuredDiagnostic) int {
	switch {
	case left.Stage != right.Stage:
		return compareStrings(string(left.Stage), string(right.Stage))
	case left.FilePath != right.FilePath:
		return compareStrings(left.FilePath, right.FilePath)
	default:
		return compareStrings(left.Message, right.Message)
	}
}

func compareStrings(left string, right string) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func uniqueByKey[T any, K comparable](items []T, keyFn func(T) K) []T {
	seen := make(map[K]struct{}, len(items))
	unique := make([]T, 0, len(items))

	for _, item := range items {
		key := keyFn(item)
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, item)
	}

	return unique
}
