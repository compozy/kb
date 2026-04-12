package graph

import (
	"reflect"
	"testing"

	"github.com/user/kb/internal/models"
)

func TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles(t *testing.T) {
	t.Parallel()

	snapshot := NormalizeGraph("/workspace", nil)

	if snapshot.RootPath != "/workspace" {
		t.Fatalf("root path = %q, want %q", snapshot.RootPath, "/workspace")
	}

	assertEmptyGraphSnapshot(t, snapshot)
}

func TestNormalizeGraphPassesThroughSingleParsedFile(t *testing.T) {
	t.Parallel()

	parsedFile := parsedFileFixture("alpha.ts")
	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{parsedFile})

	if len(snapshot.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(snapshot.Files))
	}

	if got := snapshot.Files[0].ID; got != parsedFile.File.ID {
		t.Fatalf("file ID = %q, want %q", got, parsedFile.File.ID)
	}

	if got := snapshot.Symbols[0].ID; got != parsedFile.Symbols[0].ID {
		t.Fatalf("first symbol ID = %q, want %q", got, parsedFile.Symbols[0].ID)
	}

	if !reflect.DeepEqual(snapshot.ExternalNodes, parsedFile.ExternalNodes) {
		t.Fatalf("external nodes = %#v, want %#v", snapshot.ExternalNodes, parsedFile.ExternalNodes)
	}

	if !reflect.DeepEqual(snapshot.Relations, parsedFile.Relations) {
		t.Fatalf("relations = %#v, want %#v", snapshot.Relations, parsedFile.Relations)
	}

	if !reflect.DeepEqual(snapshot.Diagnostics, parsedFile.Diagnostics) {
		t.Fatalf("diagnostics = %#v, want %#v", snapshot.Diagnostics, parsedFile.Diagnostics)
	}

	expectedSymbolIDs := []string{parsedFile.Symbols[0].ID, parsedFile.Symbols[1].ID}
	if !reflect.DeepEqual(snapshot.Files[0].SymbolIDs, expectedSymbolIDs) {
		t.Fatalf("symbol IDs = %#v, want %#v", snapshot.Files[0].SymbolIDs, expectedSymbolIDs)
	}
}

func TestNormalizeGraphDeduplicatesFilesSymbolsExternalNodesAndRelations(t *testing.T) {
	t.Parallel()

	first := parsedFileFixture("alpha.ts")
	second := parsedFileFixture("beta.ts")

	duplicateFile := second.File
	duplicateFile.ID = first.File.ID
	duplicateFile.FilePath = first.File.FilePath

	duplicateSymbol := second.Symbols[0]
	duplicateSymbol.ID = first.Symbols[0].ID
	duplicateSymbol.FilePath = first.Symbols[0].FilePath

	duplicateExternal := second.ExternalNodes[0]
	duplicateExternal.ID = first.ExternalNodes[0].ID

	duplicateRelation := second.Relations[0]
	duplicateRelation.FromID = first.Relations[0].FromID
	duplicateRelation.ToID = first.Relations[0].ToID
	duplicateRelation.Type = first.Relations[0].Type
	duplicateRelation.Confidence = first.Relations[0].Confidence

	second.File = duplicateFile
	second.Symbols[0] = duplicateSymbol
	second.ExternalNodes[0] = duplicateExternal
	second.Relations[0] = duplicateRelation

	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{first, second})

	if len(snapshot.Files) != 1 {
		t.Fatalf("expected 1 deduplicated file, got %d", len(snapshot.Files))
	}

	if len(snapshot.Symbols) != 3 {
		t.Fatalf("expected 3 deduplicated symbols, got %d", len(snapshot.Symbols))
	}

	if len(snapshot.ExternalNodes) != 1 {
		t.Fatalf("expected 1 deduplicated external node, got %d", len(snapshot.ExternalNodes))
	}

	if len(snapshot.Relations) != 3 {
		t.Fatalf("expected 3 deduplicated relations, got %d", len(snapshot.Relations))
	}
}

func TestNormalizeGraphAttachesSymbolIDsToParentFiles(t *testing.T) {
	t.Parallel()

	first := parsedFileFixture("zeta.ts")
	second := parsedFileFixture("alpha.ts")

	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{first, second})

	if len(snapshot.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(snapshot.Files))
	}

	alphaFile := snapshot.Files[0]
	zetaFile := snapshot.Files[1]

	expectedAlphaSymbols := []string{"symbol:alpha.ts:AlphaService", "symbol:alpha.ts:alphaHelper"}
	expectedZetaSymbols := []string{"symbol:zeta.ts:ZetaService", "symbol:zeta.ts:zetaHelper"}

	if !reflect.DeepEqual(alphaFile.SymbolIDs, expectedAlphaSymbols) {
		t.Fatalf("alpha symbol IDs = %#v, want %#v", alphaFile.SymbolIDs, expectedAlphaSymbols)
	}

	if !reflect.DeepEqual(zetaFile.SymbolIDs, expectedZetaSymbols) {
		t.Fatalf("zeta symbol IDs = %#v, want %#v", zetaFile.SymbolIDs, expectedZetaSymbols)
	}
}

func TestNormalizeGraphSortsCollectionsDeterministically(t *testing.T) {
	t.Parallel()

	first := parsedFileFixture("zeta.ts")
	second := parsedFileFixture("alpha.ts")

	first.Symbols[0].ID = "symbol:zeta.ts:002"
	first.Symbols[1].ID = "symbol:zeta.ts:001"
	first.ExternalNodes[0].ID = "external:pkg-zeta"
	second.ExternalNodes[0].ID = "external:pkg-alpha"

	first.Relations[0] = models.RelationEdge{
		FromID:     "symbol:zeta.ts:002",
		ToID:       "external:pkg-zeta",
		Type:       models.RelImports,
		Confidence: models.ConfidenceSyntactic,
	}
	second.Relations[0] = models.RelationEdge{
		FromID:     "symbol:alpha.ts:001",
		ToID:       "external:pkg-alpha",
		Type:       models.RelCalls,
		Confidence: models.ConfidenceSemantic,
	}

	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{first, second})

	assertOrderedIDs(t, []models.GraphFile{snapshot.Files[0], snapshot.Files[1]}, []string{"file:alpha.ts", "file:zeta.ts"})
	assertOrderedSymbolIDs(t, snapshot.Symbols, []string{"symbol:alpha.ts:AlphaService", "symbol:alpha.ts:alphaHelper", "symbol:zeta.ts:001", "symbol:zeta.ts:002"})
	assertOrderedExternalIDs(t, snapshot.ExternalNodes, []string{"external:pkg-alpha", "external:pkg-zeta"})

	expectedRelationOrder := []models.RelationEdge{
		second.Relations[1],
		first.Relations[1],
		second.Relations[0],
		first.Relations[0],
	}
	if !reflect.DeepEqual(snapshot.Relations, expectedRelationOrder) {
		t.Fatalf("relations = %#v, want %#v", snapshot.Relations, expectedRelationOrder)
	}
}

func TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage(t *testing.T) {
	t.Parallel()

	first := parsedFileFixture("alpha.ts")
	second := parsedFileFixture("beta.ts")

	first.Diagnostics = []models.StructuredDiagnostic{
		{Stage: models.StageValidate, FilePath: "zeta.ts", Message: "validate later"},
		{Stage: models.StageParse, FilePath: "beta.ts", Message: "parse beta"},
	}
	second.Diagnostics = []models.StructuredDiagnostic{
		{Stage: models.StageParse, FilePath: "", Message: "parse no file"},
		{Stage: models.StageParse, FilePath: "alpha.ts", Message: "parse alpha"},
		{Stage: models.StageRender, FilePath: "alpha.ts", Message: "render alpha"},
	}

	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{first, second})

	expected := []models.StructuredDiagnostic{
		{Stage: models.StageParse, FilePath: "", Message: "parse no file"},
		{Stage: models.StageParse, FilePath: "alpha.ts", Message: "parse alpha"},
		{Stage: models.StageParse, FilePath: "beta.ts", Message: "parse beta"},
		{Stage: models.StageRender, FilePath: "alpha.ts", Message: "render alpha"},
		{Stage: models.StageValidate, FilePath: "zeta.ts", Message: "validate later"},
	}
	if !reflect.DeepEqual(snapshot.Diagnostics, expected) {
		t.Fatalf("diagnostics = %#v, want %#v", snapshot.Diagnostics, expected)
	}
}

func TestNormalizeGraphOmitsDiagnosticOnlyFiles(t *testing.T) {
	t.Parallel()

	parsedFile := models.ParsedFile{
		File: models.GraphFile{
			ID:       "file:broken.ts",
			NodeType: "file",
			FilePath: "broken.ts",
			Language: models.LangTS,
		},
		Symbols:       []models.SymbolNode{},
		ExternalNodes: []models.ExternalNode{},
		Relations:     []models.RelationEdge{},
		Diagnostics: []models.StructuredDiagnostic{
			{Stage: models.StageParse, FilePath: "broken.ts", Message: "syntax error"},
		},
	}

	snapshot := NormalizeGraph("/workspace", []models.ParsedFile{parsedFile})

	if len(snapshot.Files) != 0 {
		t.Fatalf("expected diagnostic-only file to be omitted, got %d files", len(snapshot.Files))
	}

	if len(snapshot.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(snapshot.Diagnostics))
	}
}

func assertEmptyGraphSnapshot(t *testing.T, snapshot models.GraphSnapshot) {
	t.Helper()

	if len(snapshot.Files) != 0 {
		t.Fatalf("expected no files, got %d", len(snapshot.Files))
	}
	if len(snapshot.Symbols) != 0 {
		t.Fatalf("expected no symbols, got %d", len(snapshot.Symbols))
	}
	if len(snapshot.ExternalNodes) != 0 {
		t.Fatalf("expected no external nodes, got %d", len(snapshot.ExternalNodes))
	}
	if len(snapshot.Relations) != 0 {
		t.Fatalf("expected no relations, got %d", len(snapshot.Relations))
	}
	if len(snapshot.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(snapshot.Diagnostics))
	}
}

func assertOrderedIDs(t *testing.T, files []models.GraphFile, expected []string) {
	t.Helper()

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for index, file := range files {
		if file.ID != expected[index] {
			t.Fatalf("file %d ID = %q, want %q", index, file.ID, expected[index])
		}
	}
}

func assertOrderedSymbolIDs(t *testing.T, symbols []models.SymbolNode, expected []string) {
	t.Helper()

	if len(symbols) != len(expected) {
		t.Fatalf("expected %d symbols, got %d", len(expected), len(symbols))
	}

	for index, symbol := range symbols {
		if symbol.ID != expected[index] {
			t.Fatalf("symbol %d ID = %q, want %q", index, symbol.ID, expected[index])
		}
	}
}

func assertOrderedExternalIDs(t *testing.T, nodes []models.ExternalNode, expected []string) {
	t.Helper()

	if len(nodes) != len(expected) {
		t.Fatalf("expected %d external nodes, got %d", len(expected), len(nodes))
	}

	for index, node := range nodes {
		if node.ID != expected[index] {
			t.Fatalf("external node %d ID = %q, want %q", index, node.ID, expected[index])
		}
	}
}

func parsedFileFixture(relativePath string) models.ParsedFile {
	baseName := relativePath[:len(relativePath)-3]

	return models.ParsedFile{
		File: models.GraphFile{
			ID:        "file:" + relativePath,
			NodeType:  "file",
			FilePath:  relativePath,
			Language:  models.LangTS,
			ModuleDoc: "module doc for " + relativePath,
			SymbolIDs: []string{"stale"},
		},
		Symbols: []models.SymbolNode{
			{
				ID:         "symbol:" + relativePath + ":" + capitalize(baseName) + "Service",
				NodeType:   "symbol",
				Name:       capitalize(baseName) + "Service",
				SymbolKind: "function",
				Language:   models.LangTS,
				FilePath:   relativePath,
				StartLine:  1,
				EndLine:    3,
				Exported:   true,
			},
			{
				ID:         "symbol:" + relativePath + ":" + baseName + "Helper",
				NodeType:   "symbol",
				Name:       baseName + "Helper",
				SymbolKind: "function",
				Language:   models.LangTS,
				FilePath:   relativePath,
				StartLine:  4,
				EndLine:    5,
				Exported:   false,
			},
		},
		ExternalNodes: []models.ExternalNode{
			{
				ID:       "external:pkg-" + baseName,
				NodeType: "external",
				Source:   "pkg/" + baseName,
				Label:    "pkg/" + baseName,
			},
		},
		Relations: []models.RelationEdge{
			{
				FromID:     "file:" + relativePath,
				ToID:       "external:pkg-" + baseName,
				Type:       models.RelImports,
				Confidence: models.ConfidenceSyntactic,
			},
			{
				FromID:     "file:" + relativePath,
				ToID:       "symbol:" + relativePath + ":" + capitalize(baseName) + "Service",
				Type:       models.RelContains,
				Confidence: models.ConfidenceSemantic,
			},
		},
		Diagnostics: []models.StructuredDiagnostic{
			{
				Code:     "fixture",
				Severity: models.SeverityWarning,
				Stage:    models.StageParse,
				Message:  "fixture diagnostic for " + relativePath,
				FilePath: relativePath,
				Language: models.LangTS,
			},
		},
	}
}

func capitalize(value string) string {
	if value == "" {
		return ""
	}

	return string(value[0]-32) + value[1:]
}
