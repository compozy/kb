//go:build integration

package vault_test

import (
	"strings"
	"testing"

	"github.com/compozy/kb/internal/metrics"
	"github.com/compozy/kb/internal/vault"
	"gopkg.in/yaml.v3"
)

func TestRenderDocumentsIntegrationProducesFullDocumentSet(t *testing.T) {
	graph := testGraphFixture()
	metricResult := metrics.ComputeMetrics(graph)

	documents := vault.RenderDocuments(graph, metricResult, testTopicFixture())
	baseFiles := vault.RenderBaseFiles(metricResult)

	if len(documents) != 24 {
		t.Fatalf("expected 24 rendered markdown documents, got %d", len(documents))
	}

	if len(baseFiles) != 11 {
		t.Fatalf("expected 11 rendered base files, got %d", len(baseFiles))
	}

	expectedPaths := []string{
		"raw/codebase/files/commands/run.ts.md",
		"raw/codebase/files/src/alpha.ts.md",
		"raw/codebase/symbols/alpha--src-alpha-ts-l10.md",
		"raw/codebase/indexes/directories/src.md",
		"raw/codebase/indexes/languages/ts.md",
		vault.GetWikiConceptPath("Codebase Overview"),
		vault.GetWikiConceptPath("Circular Dependencies"),
		"wiki/index/Dashboard.md",
	}

	for _, expectedPath := range expectedPaths {
		findDocument(t, documents, expectedPath)
	}

	for _, document := range documents {
		frontmatter, markdownBody := parseFrontmatter(t, document.Body)
		if len(frontmatter) == 0 {
			t.Fatalf("document %s has empty frontmatter", document.RelativePath)
		}
		if strings.TrimSpace(markdownBody) == "" {
			t.Fatalf("document %s has empty markdown body", document.RelativePath)
		}
	}

	for _, baseFile := range baseFiles {
		rendered := vault.RenderBaseDefinition(baseFile.Definition)
		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(rendered), &parsed); err != nil {
			t.Fatalf("base file %s did not parse as YAML: %v\n%s", baseFile.RelativePath, err, rendered)
		}
	}
}
