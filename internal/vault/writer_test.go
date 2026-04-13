package vault_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/metrics"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/vault"
	"gopkg.in/yaml.v3"
)

func TestWriteVaultCreatesTopicSkeletonAndManagedFiles(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)

	result, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	})
	if err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	expectedCounts := countKinds(documents)
	if result != expectedCounts {
		t.Fatalf("write result = %#v, want %#v", result, expectedCounts)
	}

	for _, directoryPath := range []string{
		topic.TopicPath,
		filepath.Join(topic.TopicPath, "raw"),
		filepath.Join(topic.TopicPath, "raw", "articles"),
		filepath.Join(topic.TopicPath, "raw", "bookmarks"),
		filepath.Join(topic.TopicPath, "raw", "github"),
		filepath.Join(topic.TopicPath, "raw", "codebase"),
		filepath.Join(topic.TopicPath, "wiki"),
		filepath.Join(topic.TopicPath, "wiki", "concepts"),
		filepath.Join(topic.TopicPath, "wiki", "index"),
		filepath.Join(topic.TopicPath, "outputs"),
		filepath.Join(topic.TopicPath, "outputs", "briefings"),
		filepath.Join(topic.TopicPath, "outputs", "queries"),
		filepath.Join(topic.TopicPath, "outputs", "diagrams"),
		filepath.Join(topic.TopicPath, "outputs", "reports"),
		filepath.Join(topic.TopicPath, "bases"),
	} {
		assertDirExists(t, directoryPath)
	}

	rawPath := filepath.Join(topic.TopicPath, filepath.FromSlash("raw/codebase/files/src/alpha.ts.md"))
	rawContent := readFile(t, rawPath)
	frontmatter, body := parseFrontmatter(t, rawContent)

	if frontmatter["source_path"] != "src/alpha.ts" {
		t.Fatalf("source_path = %#v, want src/alpha.ts", frontmatter["source_path"])
	}
	if !strings.Contains(body, "[[demo-repo/raw/codebase/symbols/alpha--src-alpha-ts-l10|Alpha (function)]]") {
		t.Fatalf("expected raw file body to retain wiki links, got:\n%s", body)
	}

	wikiPath := filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetWikiConceptPath("Codebase Overview")))
	wikiContent := readFile(t, wikiPath)
	if !strings.Contains(wikiContent, "generator: \"kodebase\"") {
		t.Fatalf("expected managed wiki frontmatter in concept document, got:\n%s", wikiContent)
	}

	basePath := filepath.Join(topic.TopicPath, filepath.FromSlash("bases/symbol-explorer.base"))
	baseContent := readFile(t, basePath)
	var parsedBase map[string]interface{}
	if err := yaml.Unmarshal([]byte(baseContent), &parsedBase); err != nil {
		t.Fatalf("base definition did not parse as YAML: %v\n%s", err, baseContent)
	}
	if _, exists := parsedBase["views"]; !exists {
		t.Fatalf("expected base definition views, got %#v", parsedBase)
	}

	if target, err := os.Readlink(filepath.Join(topic.TopicPath, "AGENTS.md")); err != nil {
		t.Fatalf("expected AGENTS.md symlink: %v", err)
	} else if target != "CLAUDE.md" {
		t.Fatalf("AGENTS.md target = %q, want %q", target, "CLAUDE.md")
	}

	for _, gitkeepPath := range []string{
		filepath.Join(topic.TopicPath, "raw", "articles", ".gitkeep"),
		filepath.Join(topic.TopicPath, "raw", "bookmarks", ".gitkeep"),
		filepath.Join(topic.TopicPath, "raw", "github", ".gitkeep"),
		filepath.Join(topic.TopicPath, "outputs", "briefings", ".gitkeep"),
		filepath.Join(topic.TopicPath, "outputs", "queries", ".gitkeep"),
		filepath.Join(topic.TopicPath, "outputs", "diagrams", ".gitkeep"),
		filepath.Join(topic.TopicPath, "outputs", "reports", ".gitkeep"),
	} {
		assertFileExists(t, gitkeepPath)
	}
}

func TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	options := vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}

	if _, err := vault.WriteVault(context.Background(), options); err != nil {
		t.Fatalf("first WriteVault returned error: %v", err)
	}

	claudePath := filepath.Join(topic.TopicPath, "CLAUDE.md")
	claudeContent := readFile(t, claudePath)
	if !strings.Contains(claudeContent, "# Demo Repo") {
		t.Fatalf("expected topic title in CLAUDE.md, got:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "**Domain:** `demo-repo`") {
		t.Fatalf("expected domain metadata in CLAUDE.md, got:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "[[demo-repo/wiki/index/Dashboard|Dashboard]]") {
		t.Fatalf("expected dashboard link in CLAUDE.md, got:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, vault.ToTopicWikiLink("demo-repo", vault.GetWikiConceptPath("Codebase Overview"), "Codebase Overview")) {
		t.Fatalf("expected concept link in CLAUDE.md, got:\n%s", claudeContent)
	}

	logPath := filepath.Join(topic.TopicPath, "log.md")
	firstLog := readFile(t, logPath)
	if !strings.Contains(firstLog, "## [2026-04-09] bootstrap | topic scaffolded") {
		t.Fatalf("expected bootstrap entry in log, got:\n%s", firstLog)
	}
	if !strings.Contains(firstLog, "## [2026-04-09] ingest | codebase (4 files, 4 symbols)") {
		t.Fatalf("expected ingest entry in log, got:\n%s", firstLog)
	}

	if _, err := vault.WriteVault(context.Background(), options); err != nil {
		t.Fatalf("second WriteVault returned error: %v", err)
	}

	secondLog := readFile(t, logPath)
	if len(secondLog) <= len(firstLog) {
		t.Fatalf("expected second log write to append content")
	}
	if got := strings.Count(secondLog, "## [2026-04-09] ingest | codebase (4 files, 4 symbols)"); got != 2 {
		t.Fatalf("expected two ingest entries after second write, got %d\n%s", got, secondLog)
	}
	if got := strings.Count(secondLog, "## [2026-04-09] compile | refreshed starter wiki"); got != 2 {
		t.Fatalf("expected two compile entries after second write, got %d\n%s", got, secondLog)
	}
}

func TestWriteVaultReportsProgressForPersistedFiles(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	var progress []vault.WriteProgress

	_, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
		Progress: func(update vault.WriteProgress) {
			progress = append(progress, update)
		},
	})
	if err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	expectedTotal := len(documents) + len(baseFiles) + 2
	if len(progress) != expectedTotal {
		t.Fatalf("progress events = %d, want %d", len(progress), expectedTotal)
	}

	last := progress[len(progress)-1]
	if last.Completed != expectedTotal || last.Total != expectedTotal {
		t.Fatalf("last progress event = %#v, want completed/total %d", last, expectedTotal)
	}
}

func TestWriteVaultRemovesStaleManagedWikiConceptsOnly(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	options := vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}

	if _, err := vault.WriteVault(context.Background(), options); err != nil {
		t.Fatalf("initial WriteVault returned error: %v", err)
	}

	manualConceptPath := filepath.Join(topic.TopicPath, filepath.FromSlash("wiki/concepts/Manual Notes.md"))
	manualConceptBody := strings.Join([]string{
		"---",
		`title: "Manual Notes"`,
		"---",
		"",
		"# Manual Notes",
		"",
		"This page is unmanaged and should survive regeneration.",
		"",
	}, "\n")
	if err := os.WriteFile(manualConceptPath, []byte(manualConceptBody), 0o644); err != nil {
		t.Fatalf("write manual concept: %v", err)
	}

	filteredDocuments := filterOutDocument(documents, vault.GetWikiConceptPath("Dead Code Report"))
	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: filteredDocuments,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("second WriteVault returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetWikiConceptPath("Dead Code Report")))); !os.IsNotExist(err) {
		t.Fatalf("expected stale managed concept to be removed, stat err = %v", err)
	}

	assertFileExists(t, filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetWikiConceptPath("Codebase Overview"))))
	assertFileExists(t, manualConceptPath)
}

func TestWriteVaultRejectsInvalidRenderedDocument(t *testing.T) {
	t.Parallel()

	topic, graph, _, baseFiles := testWriteVaultInputs(t)
	badDocument := models.RenderedDocument{
		Kind:         models.DocWiki,
		ManagedArea:  models.AreaWikiConcept,
		RelativePath: "wiki/concepts/Broken.md",
		Body:         "# missing frontmatter",
	}

	_, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: []models.RenderedDocument{badDocument},
		BaseFiles: baseFiles,
	})
	if err == nil {
		t.Fatal("expected WriteVault to reject document without frontmatter")
	}
	if !strings.Contains(err.Error(), "missing YAML frontmatter") {
		t.Fatalf("expected frontmatter validation error, got %v", err)
	}
}

func TestWriteVaultPreservesCodebaseSkeletonWhenNoRawDocumentsAreRendered(t *testing.T) {
	t.Parallel()

	writableTopic, graph, documents, baseFiles := testWriteVaultInputs(t)
	nonRawDocuments := make([]models.RenderedDocument, 0, len(documents))
	for _, document := range documents {
		if document.Kind == models.DocRaw {
			continue
		}
		nonRawDocuments = append(nonRawDocuments, document)
	}

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     writableTopic,
		Graph:     graph,
		Documents: nonRawDocuments,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	for _, relativePath := range []string{
		"raw/codebase/files",
		"raw/codebase/symbols",
	} {
		assertDirExists(t, filepath.Join(writableTopic.TopicPath, filepath.FromSlash(relativePath)))
	}

	if _, err := topic.Info(writableTopic.VaultPath, writableTopic.Slug); err != nil {
		t.Fatalf("topic.Info returned error after empty raw write: %v", err)
	}
}

func testWriteVaultInputs(t *testing.T) (models.TopicMetadata, models.GraphSnapshot, []models.RenderedDocument, []models.BaseFile) {
	t.Helper()

	topic := testWritableTopicFixture(t)
	graph := testGraphFixture()
	metricResult := metrics.ComputeMetrics(graph)
	documents := vault.RenderDocuments(graph, metricResult, topic)
	baseFiles := vault.RenderBaseFiles(metricResult)

	return topic, graph, documents, baseFiles
}

func testWritableTopicFixture(t *testing.T) models.TopicMetadata {
	t.Helper()

	root := t.TempDir()
	rootPath := filepath.Join(root, "repo")
	vaultPath := filepath.Join(root, ".kb", "vault")
	if err := os.MkdirAll(rootPath, 0o755); err != nil {
		t.Fatalf("create root path: %v", err)
	}

	return models.TopicMetadata{
		RootPath:  rootPath,
		Slug:      "demo-repo",
		Title:     "Demo Repo",
		Domain:    "demo-repo",
		Today:     "2026-04-09",
		VaultPath: vaultPath,
		TopicPath: filepath.Join(vaultPath, "demo-repo"),
	}
}

func countKinds(documents []models.RenderedDocument) vault.WriteVaultResult {
	var counts vault.WriteVaultResult

	for _, document := range documents {
		switch document.Kind {
		case models.DocRaw:
			counts.RawDocumentsWritten++
		case models.DocWiki:
			counts.WikiDocumentsWritten++
		case models.DocIndex:
			counts.IndexDocumentsWritten++
		}
	}

	return counts
}

func filterOutDocument(documents []models.RenderedDocument, relativePath string) []models.RenderedDocument {
	filtered := make([]models.RenderedDocument, 0, len(documents))
	for _, document := range documents {
		if document.RelativePath == relativePath {
			continue
		}
		filtered = append(filtered, document)
	}
	return filtered
}

func readFile(t *testing.T, filePath string) string {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read %s: %v", filePath, err)
	}

	return string(content)
}

func assertDirExists(t *testing.T, directoryPath string) {
	t.Helper()

	info, err := os.Stat(directoryPath)
	if err != nil {
		t.Fatalf("stat %s: %v", directoryPath, err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", directoryPath)
	}
}

func assertFileExists(t *testing.T, filePath string) {
	t.Helper()

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("stat %s: %v", filePath, err)
	}
	if info.IsDir() {
		t.Fatalf("%s is a directory, want file", filePath)
	}
}
