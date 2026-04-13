package vault_test

import (
	"strings"
	"testing"

	"github.com/compozy/kb/internal/metrics"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
	"gopkg.in/yaml.v3"
)

func TestRenderDocumentsProducesRawWikiAndBaseSurfaces(t *testing.T) {
	t.Parallel()

	graph := testGraphFixture()
	metricResult := metrics.ComputeMetrics(graph)

	documents := vault.RenderDocuments(graph, metricResult, testTopicFixture())
	baseFiles := vault.RenderBaseFiles(metricResult)

	if len(documents) == 0 {
		t.Fatal("expected rendered documents")
	}

	if len(baseFiles) != 11 {
		t.Fatalf("expected 11 base files from the reference renderer, got %d", len(baseFiles))
	}

	if findDocument(t, documents, "raw/codebase/files/src/alpha.ts.md").Kind != models.DocRaw {
		t.Fatal("expected raw file document to have raw kind")
	}

	if findDocument(t, documents, vault.GetWikiConceptPath("Codebase Overview")).Kind != models.DocWiki {
		t.Fatal("expected concept article to have wiki kind")
	}

	if findDocument(t, documents, "wiki/index/Dashboard.md").Kind != models.DocIndex {
		t.Fatal("expected dashboard to have index kind")
	}
}

func TestRenderDocumentsRawFileFrontmatterAndBody(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, "raw/codebase/files/src/alpha.ts.md")

	if document.ManagedArea != models.AreaRawCodebase {
		t.Fatalf("managed area = %q, want %q", document.ManagedArea, models.AreaRawCodebase)
	}

	if got := document.Frontmatter["source_path"]; got != "src/alpha.ts" {
		t.Fatalf("source_path = %#v, want src/alpha.ts", got)
	}

	if got := document.Frontmatter["language"]; got != "ts" {
		t.Fatalf("language = %#v, want ts", got)
	}

	if !strings.Contains(document.Body, "title: \"Codebase File: src/alpha.ts\"") {
		t.Fatalf("expected file title frontmatter in body, got:\n%s", document.Body)
	}

	if !strings.Contains(document.Body, "[[demo-repo/raw/codebase/symbols/alpha--src-alpha-ts-l10|Alpha (function)]]") {
		t.Fatalf("expected symbol wiki-link in raw file document, got:\n%s", document.Body)
	}
}

func TestRenderDocumentsRawSymbolFrontmatterAndSignature(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, "raw/codebase/symbols/branchy--src-orphan-ts-l1.md")

	if got := document.Frontmatter["symbol_kind"]; got != "function" {
		t.Fatalf("symbol_kind = %#v, want function", got)
	}

	if got := document.Frontmatter["start_line"]; got != 1 {
		t.Fatalf("start_line = %#v, want 1", got)
	}

	if got := document.Frontmatter["end_line"]; got != 80 {
		t.Fatalf("end_line = %#v, want 80", got)
	}

	if !strings.Contains(document.Body, "cyclomatic_complexity: 12") {
		t.Fatalf("expected cyclomatic complexity in symbol frontmatter, got:\n%s", document.Body)
	}

	if !strings.Contains(document.Body, "function Branchy(input: string): string") {
		t.Fatalf("expected signature block in symbol document, got:\n%s", document.Body)
	}
}

func TestRenderDocumentsDirectoryIndexUsesWikiLinks(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, "raw/codebase/indexes/directories/src.md")

	expectedLinks := []string{
		"[[demo-repo/raw/codebase/files/src/alpha.ts|src/alpha.ts]]",
		"[[demo-repo/raw/codebase/files/src/beta.ts|src/beta.ts]]",
		"[[demo-repo/raw/codebase/files/src/orphan.ts|src/orphan.ts]]",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(document.Body, link) {
			t.Fatalf("expected directory index to contain %s, got:\n%s", link, document.Body)
		}
	}
}

func TestRenderDocumentsCodebaseOverviewContainsSummary(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, vault.GetWikiConceptPath("Codebase Overview"))

	if !strings.Contains(document.Body, "The corpus contains 4 parsed source files, 4 symbols, and 10 extracted relations.") {
		t.Fatalf("expected overview counts in body, got:\n%s", document.Body)
	}

	if !strings.Contains(document.Body, vault.ToTopicWikiLink("demo-repo", vault.GetWikiConceptPath("Module Health"), "Module Health")) {
		t.Fatalf("expected module health cross-reference in overview, got:\n%s", document.Body)
	}
}

func TestRenderDocumentsDependencyHotspotsListsTopFiles(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, vault.GetWikiConceptPath("Dependency Hotspots"))

	if !strings.Contains(document.Body, "[[demo-repo/raw/codebase/files/src/alpha.ts|src/alpha.ts]]") {
		t.Fatalf("expected alpha hotspot in dependency article, got:\n%s", document.Body)
	}
}

func TestRenderDocumentsCircularDependenciesListsGroups(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, vault.GetWikiConceptPath("Circular Dependencies"))

	if !strings.Contains(document.Body, "[[demo-repo/raw/codebase/files/src/alpha.ts|src/alpha.ts]] · [[demo-repo/raw/codebase/files/src/beta.ts|src/beta.ts]]") {
		t.Fatalf("expected cyclic group links in circular dependency article, got:\n%s", document.Body)
	}
}

func TestRenderDocumentsRawSourcesIncludeScrapedFrontmatter(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	wantScraped := "2026-04-09"

	for _, relativePath := range []string{
		"raw/codebase/files/src/alpha.ts.md",
		"raw/codebase/symbols/alpha--src-alpha-ts-l10.md",
		"raw/codebase/indexes/directories/src.md",
		"raw/codebase/indexes/languages/ts.md",
	} {
		document := findDocument(t, documents, relativePath)
		if got := document.Frontmatter["scraped"]; got != wantScraped {
			t.Fatalf("%s scraped = %#v, want %q", relativePath, got, wantScraped)
		}
	}
}

func TestRenderDocumentsCircularDependenciesKeepSourceEvidenceWithoutCycles(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, vault.GetWikiConceptPath("Circular Dependencies"))

	sources, ok := document.Frontmatter["sources"].([]string)
	if !ok {
		t.Fatalf("sources type = %T, want []string", document.Frontmatter["sources"])
	}
	if len(sources) == 0 {
		t.Fatal("expected circular dependency article to keep source evidence even when no cycles are detected")
	}
}

func TestRenderDocumentsDashboardLinksToAllConceptArticles(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, "wiki/index/Dashboard.md")

	for _, title := range []string{
		"Codebase Overview",
		"Directory Map",
		"Symbol Taxonomy",
		"Dependency Hotspots",
		"Complexity Hotspots",
		"Module Health",
		"Dead Code Report",
		"Code Smells",
		"Circular Dependencies",
		"High-Impact Symbols",
	} {
		link := vault.ToTopicWikiLink("demo-repo", vault.GetWikiConceptPath(title), title)
		if !strings.Contains(document.Body, link) {
			t.Fatalf("expected dashboard to link to %s, got:\n%s", title, document.Body)
		}
	}
}

func TestRenderBaseDefinitionsProduceValidYAML(t *testing.T) {
	t.Parallel()

	baseFiles := vault.RenderBaseFiles(metrics.ComputeMetrics(testGraphFixture()))
	for _, baseFile := range baseFiles {
		rendered := vault.RenderBaseDefinition(baseFile.Definition)

		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(rendered), &parsed); err != nil {
			t.Fatalf("base file %s did not parse as YAML: %v\n%s", baseFile.RelativePath, err, rendered)
		}

		if _, exists := parsed["views"]; !exists {
			t.Fatalf("base file %s missing views: %#v", baseFile.RelativePath, parsed)
		}
	}
}

func TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	for _, document := range documents {
		if strings.TrimSpace(document.Body) == "" {
			t.Fatalf("document %s has empty body", document.RelativePath)
		}

		switch document.Kind {
		case models.DocRaw, models.DocWiki, models.DocIndex:
		default:
			t.Fatalf("document %s has invalid kind %q", document.RelativePath, document.Kind)
		}

		frontmatter, markdownBody := parseFrontmatter(t, document.Body)
		if len(frontmatter) == 0 {
			t.Fatalf("document %s has empty frontmatter", document.RelativePath)
		}
		if strings.TrimSpace(markdownBody) == "" {
			t.Fatalf("document %s has empty markdown body", document.RelativePath)
		}
	}
}

func TestRenderDocumentsUseTopicWikiLinkSyntax(t *testing.T) {
	t.Parallel()

	documents := renderFixtureDocuments(t)
	document := findDocument(t, documents, "raw/codebase/files/commands/run.ts.md")

	if !strings.Contains(document.Body, "[[demo-repo/raw/codebase/symbols/main--commands-run-ts-l1|main (function)]]") {
		t.Fatalf("expected topic wiki-link syntax in raw document, got:\n%s", document.Body)
	}
}

func renderFixtureDocuments(t *testing.T) []models.RenderedDocument {
	t.Helper()

	graph := testGraphFixture()
	metricResult := metrics.ComputeMetrics(graph)
	return vault.RenderDocuments(graph, metricResult, testTopicFixture())
}

func findDocument(t *testing.T, documents []models.RenderedDocument, relativePath string) models.RenderedDocument {
	t.Helper()

	for _, document := range documents {
		if document.RelativePath == relativePath {
			return document
		}
	}

	t.Fatalf("document %s not found", relativePath)
	return models.RenderedDocument{}
}

func parseFrontmatter(t *testing.T, body string) (map[string]interface{}, string) {
	t.Helper()

	if !strings.HasPrefix(body, "---\n") {
		t.Fatalf("body missing leading frontmatter delimiter:\n%s", body)
	}

	remainder := strings.TrimPrefix(body, "---\n")
	index := strings.Index(remainder, "\n---\n")
	if index < 0 {
		t.Fatalf("body missing closing frontmatter delimiter:\n%s", body)
	}

	var parsed map[string]interface{}
	frontmatterBlock := remainder[:index]
	if err := yaml.Unmarshal([]byte(frontmatterBlock), &parsed); err != nil {
		t.Fatalf("frontmatter did not parse as YAML: %v\n%s", err, frontmatterBlock)
	}

	return parsed, remainder[index+5:]
}

func testTopicFixture() models.TopicMetadata {
	return models.TopicMetadata{
		RootPath:  "/repo",
		Slug:      "demo-repo",
		Title:     "Demo Repo",
		Domain:    "demo-repo",
		Today:     "2026-04-09",
		VaultPath: "/repo/.kb/vault",
		TopicPath: "/repo/.kb/vault/demo-repo",
	}
}

func testGraphFixture() models.GraphSnapshot {
	runFileID := "file:commands/run.ts"
	alphaFileID := "file:src/alpha.ts"
	betaFileID := "file:src/beta.ts"
	orphanFileID := "file:src/orphan.ts"

	mainID := "symbol:commands/run.ts:main:function:1:12"
	alphaID := "symbol:src/alpha.ts:Alpha:function:10:30"
	betaID := "symbol:src/beta.ts:Beta:function:5:18"
	branchyID := "symbol:src/orphan.ts:Branchy:function:1:80"

	return models.GraphSnapshot{
		RootPath: "/repo",
		Files: []models.GraphFile{
			{
				ID:        runFileID,
				NodeType:  "file",
				FilePath:  "commands/run.ts",
				Language:  models.LangTS,
				ModuleDoc: "Entrypoint for the topic.",
				SymbolIDs: []string{mainID},
			},
			{
				ID:        alphaFileID,
				NodeType:  "file",
				FilePath:  "src/alpha.ts",
				Language:  models.LangTS,
				ModuleDoc: "Coordinates alpha flows.",
				SymbolIDs: []string{alphaID},
			},
			{
				ID:        betaFileID,
				NodeType:  "file",
				FilePath:  "src/beta.ts",
				Language:  models.LangTS,
				ModuleDoc: "Supports alpha with beta helpers.",
				SymbolIDs: []string{betaID},
			},
			{
				ID:        orphanFileID,
				NodeType:  "file",
				FilePath:  "src/orphan.ts",
				Language:  models.LangTS,
				ModuleDoc: "Contains an unused branchy function.",
				SymbolIDs: []string{branchyID},
			},
		},
		Symbols: []models.SymbolNode{
			{
				ID:                   mainID,
				NodeType:             "symbol",
				Name:                 "main",
				SymbolKind:           "function",
				Language:             models.LangTS,
				FilePath:             "commands/run.ts",
				StartLine:            1,
				EndLine:              12,
				Signature:            "function main(): void",
				DocComment:           "Bootstraps the demo topic.",
				Exported:             true,
				CyclomaticComplexity: 3,
			},
			{
				ID:                   alphaID,
				NodeType:             "symbol",
				Name:                 "Alpha",
				SymbolKind:           "function",
				Language:             models.LangTS,
				FilePath:             "src/alpha.ts",
				StartLine:            10,
				EndLine:              30,
				Signature:            "function Alpha(input: string): string",
				DocComment:           "Coordinates the alpha path.",
				Exported:             true,
				CyclomaticComplexity: 6,
			},
			{
				ID:                   betaID,
				NodeType:             "symbol",
				Name:                 "Beta",
				SymbolKind:           "function",
				Language:             models.LangTS,
				FilePath:             "src/beta.ts",
				StartLine:            5,
				EndLine:              18,
				Signature:            "function Beta(input: string): string",
				DocComment:           "Provides beta helpers.",
				Exported:             true,
				CyclomaticComplexity: 4,
			},
			{
				ID:                   branchyID,
				NodeType:             "symbol",
				Name:                 "Branchy",
				SymbolKind:           "function",
				Language:             models.LangTS,
				FilePath:             "src/orphan.ts",
				StartLine:            1,
				EndLine:              80,
				Signature:            "function Branchy(input: string): string",
				DocComment:           "Contains many branches but no dependents.",
				Exported:             true,
				CyclomaticComplexity: 12,
			},
		},
		ExternalNodes: []models.ExternalNode{},
		Relations: []models.RelationEdge{
			{FromID: runFileID, ToID: mainID, Type: models.RelContains, Confidence: models.ConfidenceSemantic},
			{FromID: alphaFileID, ToID: alphaID, Type: models.RelContains, Confidence: models.ConfidenceSemantic},
			{FromID: betaFileID, ToID: betaID, Type: models.RelContains, Confidence: models.ConfidenceSemantic},
			{FromID: orphanFileID, ToID: branchyID, Type: models.RelContains, Confidence: models.ConfidenceSemantic},
			{FromID: runFileID, ToID: alphaFileID, Type: models.RelImports, Confidence: models.ConfidenceSemantic},
			{FromID: alphaFileID, ToID: betaFileID, Type: models.RelImports, Confidence: models.ConfidenceSemantic},
			{FromID: betaFileID, ToID: alphaFileID, Type: models.RelImports, Confidence: models.ConfidenceSemantic},
			{FromID: mainID, ToID: alphaID, Type: models.RelCalls, Confidence: models.ConfidenceSemantic},
			{FromID: alphaID, ToID: betaID, Type: models.RelCalls, Confidence: models.ConfidenceSemantic},
			{FromID: betaID, ToID: alphaID, Type: models.RelReferences, Confidence: models.ConfidenceSemantic},
		},
		Diagnostics: []models.StructuredDiagnostic{},
	}
}
