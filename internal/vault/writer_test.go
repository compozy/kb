package vault_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/metrics"
	"github.com/compozy/kb/internal/models"
	topicpkg "github.com/compozy/kb/internal/topic"
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
		filepath.Join(topic.TopicPath, "wiki", "codebase"),
		filepath.Join(topic.TopicPath, "wiki", "codebase", "concepts"),
		filepath.Join(topic.TopicPath, "wiki", "codebase", "index"),
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
		filepath.Join(topic.TopicPath, "raw", "youtube", ".gitkeep"),
		filepath.Join(topic.TopicPath, "wiki", "concepts", ".gitkeep"),
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
	if !strings.Contains(claudeContent, "<!-- kb:codebase:start -->") || !strings.Contains(claudeContent, "<!-- kb:codebase:end -->") {
		t.Fatalf("expected managed codebase block in CLAUDE.md, got:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, vault.ToTopicWikiLink("demo-repo", vault.GetWikiIndexPath(vault.CodebaseDashboardTitle), vault.CodebaseDashboardTitle)) {
		t.Fatalf("expected codebase dashboard link in CLAUDE.md, got:\n%s", claudeContent)
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
	if got := strings.Count(secondLog, "## [2026-04-09] compile | refreshed codebase wiki"); got != 2 {
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

	expectedTotal := len(documents) + len(baseFiles) + 5
	if len(progress) != expectedTotal {
		t.Fatalf("progress events = %d, want %d", len(progress), expectedTotal)
	}

	progressPaths := make([]string, 0, len(progress))
	for _, update := range progress {
		progressPaths = append(progressPaths, filepath.ToSlash(update.Path))
	}
	for _, relativePath := range []string{
		filepath.ToSlash(filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetTopicIndexPath(vault.TopicDashboardTitle)))),
		filepath.ToSlash(filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetTopicIndexPath(vault.TopicConceptIndexTitle)))),
		filepath.ToSlash(filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetTopicIndexPath(vault.TopicSourceIndexTitle)))),
	} {
		if !contains(progressPaths, relativePath) {
			t.Fatalf("expected progress to report %q, got %#v", relativePath, progressPaths)
		}
	}

	last := progress[len(progress)-1]
	if last.Completed != expectedTotal || last.Total != expectedTotal {
		t.Fatalf("last progress event = %#v, want completed/total %d", last, expectedTotal)
	}
}

func TestWriteVaultResetsManagedCodebaseSubtreeOnly(t *testing.T) {
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
	manualIndexPath := filepath.Join(topic.TopicPath, filepath.FromSlash("wiki/index/Dashboard.md"))
	manualIndexBody := strings.Join([]string{
		"---",
		`domain: "demo-repo"`,
		`title: "Dashboard"`,
		`type: "index"`,
		`updated: "2026-04-09"`,
		"---",
		"",
		"# Dashboard",
		"",
		"Manual index that must survive regeneration.",
		"",
	}, "\n")
	if err := os.WriteFile(manualIndexPath, []byte(manualIndexBody), 0o644); err != nil {
		t.Fatalf("write manual index: %v", err)
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
	manualIndexContent := readFile(t, manualIndexPath)
	if !strings.Contains(manualIndexContent, "Manual index that must survive regeneration.") {
		t.Fatalf("manual index body was not preserved:\n%s", manualIndexContent)
	}
	if !strings.Contains(manualIndexContent, "<!-- kb:codebase-index:start -->") || !strings.Contains(manualIndexContent, "<!-- kb:codebase-index:end -->") {
		t.Fatalf("expected managed codebase index block in dashboard:\n%s", manualIndexContent)
	}
	if !strings.Contains(
		manualIndexContent,
		vault.ToTopicWikiLink("demo-repo", vault.GetWikiIndexPath(vault.CodebaseDashboardTitle), vault.CodebaseDashboardTitle),
	) {
		t.Fatalf("dashboard bridge missing codebase dashboard link:\n%s", manualIndexContent)
	}
}

func TestWriteVaultMigratesLegacyGeneratedIndexesAndRemovesLegacyManagedConcepts(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	if err := topicpkg.EnsureCurrentSkeleton(topic.TopicPath); err != nil {
		t.Fatalf("ensure topic skeleton: %v", err)
	}

	legacyConceptPath := filepath.Join(topic.TopicPath, "wiki", "concepts", "Kodebase - Directory Map.md")
	legacyConceptBody := strings.Join([]string{
		"---",
		`title: "Directory Map"`,
		`generator: "kodebase"`,
		"---",
		"",
		"# Directory Map",
		"",
		"Legacy generated concept that should be removed during upgrade.",
		"",
	}, "\n")
	if err := os.WriteFile(legacyConceptPath, []byte(legacyConceptBody), 0o644); err != nil {
		t.Fatalf("write legacy concept: %v", err)
	}

	manualLegacyConceptPath := filepath.Join(topic.TopicPath, "wiki", "concepts", "Kodebase - Manual Notes.md")
	manualLegacyConceptBody := strings.Join([]string{
		"---",
		`title: "Kodebase - Manual Notes"`,
		`type: "wiki"`,
		"---",
		"",
		"# Kodebase - Manual Notes",
		"",
		"Manually curated concept that must survive the upgrade cleanup.",
		"",
	}, "\n")
	if err := os.WriteFile(manualLegacyConceptPath, []byte(manualLegacyConceptBody), 0o644); err != nil {
		t.Fatalf("write manual legacy concept: %v", err)
	}

	for _, testCase := range []struct {
		title     string
		bodyLines []string
	}{
		{
			title: vault.TopicDashboardTitle,
			bodyLines: []string{
				"# Demo Repo - Dashboard",
				"",
				"Landing page for the generated Karpathy-compatible codebase topic.",
				"",
				"## Generated codebase articles",
			},
		},
		{
			title: vault.TopicConceptIndexTitle,
			bodyLines: []string{
				"# Demo Repo - Concept Index",
				"",
				"Alphabetical listing of every generated codebase wiki article in this topic.",
				"",
				"| Article | Summary |",
				"| ------- | ------- |",
			},
		},
		{
			title: vault.TopicSourceIndexTitle,
			bodyLines: []string{
				"# Demo Repo - Source Index",
				"",
				"Raw codebase snapshots currently cited by the generated codebase wiki.",
				"",
				"| Source | Cited by |",
				"| ------ | -------- |",
			},
		},
	} {
		targetPath := filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetTopicIndexPath(testCase.title)))
		legacyIndexBody := strings.Join(append([]string{
			"---",
			fmt.Sprintf(`title: %q`, testCase.title),
			`type: "index"`,
			`domain: "demo-repo"`,
			`updated: "2026-04-09"`,
			"---",
			"",
		}, testCase.bodyLines...), "\n") + "\n"
		if err := os.WriteFile(targetPath, []byte(legacyIndexBody), 0o644); err != nil {
			t.Fatalf("write legacy index %q: %v", testCase.title, err)
		}
	}

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	if _, err := os.Stat(legacyConceptPath); !os.IsNotExist(err) {
		t.Fatalf("expected legacy managed concept to be removed, stat err = %v", err)
	}
	assertFileExists(t, manualLegacyConceptPath)

	for _, testCase := range []struct {
		title            string
		legacyMarker     string
		expectedScaffold string
		expectedCodebase string
	}{
		{
			title:            vault.TopicDashboardTitle,
			legacyMarker:     "Landing page for the generated Karpathy-compatible codebase topic.",
			expectedScaffold: "Landing page for the Demo Repo knowledge base.",
			expectedCodebase: vault.CodebaseDashboardTitle,
		},
		{
			title:            vault.TopicConceptIndexTitle,
			legacyMarker:     "Alphabetical listing of every generated codebase wiki article in this topic.",
			expectedScaffold: "# Demo Repo — Concept Index",
			expectedCodebase: vault.CodebaseConceptIndexTitle,
		},
		{
			title:            vault.TopicSourceIndexTitle,
			legacyMarker:     "Raw codebase snapshots currently cited by the generated codebase wiki.",
			expectedScaffold: "# Demo Repo — Source Index",
			expectedCodebase: vault.CodebaseSourceIndexTitle,
		},
	} {
		content := readFile(t, filepath.Join(topic.TopicPath, filepath.FromSlash(vault.GetTopicIndexPath(testCase.title))))
		if strings.Contains(content, testCase.legacyMarker) {
			t.Fatalf("legacy marker still present in %q:\n%s", testCase.title, content)
		}
		if !strings.Contains(content, testCase.expectedScaffold) {
			t.Fatalf("expected scaffold content in %q:\n%s", testCase.title, content)
		}
		if !strings.Contains(
			content,
			vault.ToTopicWikiLink("demo-repo", vault.GetWikiIndexPath(testCase.expectedCodebase), testCase.expectedCodebase),
		) {
			t.Fatalf("expected codebase bridge in %q:\n%s", testCase.title, content)
		}
	}
}

func TestWriteVaultMigratesLegacyGeneratedClaudeBeforeManagedBlock(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	if err := os.MkdirAll(topic.TopicPath, 0o755); err != nil {
		t.Fatalf("create topic path: %v", err)
	}

	legacyClaude := strings.Join([]string{
		"# Demo Repo",
		"",
		"**Topic scope:** Generated codebase knowledge topic for `repo`. This topic stages raw code snapshots in `raw/codebase/` and compiles a starter wiki in `wiki/`.",
		"",
		"**Domain:** `demo-repo`",
		"",
		"This file is the schema document for the topic. The `kodebase` CLI manages `raw/codebase/`, `wiki/index/`, and wiki concept pages with `generator: kodebase` frontmatter. Everything else may be extended manually without being overwritten.",
		"",
		"## Audit log",
		"",
		"See [log.md](log.md) for the append-only record of ingest and compile operations.",
		"",
		"## Current wiki articles",
		"",
		"- " + vault.ToTopicWikiLink("demo-repo", "wiki/concepts/Kodebase - Directory Map.md", "Directory Map"),
		"- " + vault.ToTopicWikiLink("demo-repo", "wiki/concepts/Kodebase - Module Health.md", "Module Health"),
		"",
		"## Codebase corpus",
		"",
		"- Parsed files: 2",
		"- Parsed symbols: 3",
		"",
		"## Managed starter wiki",
		"",
		"- " + vault.ToTopicWikiLink("demo-repo", "wiki/index/Dashboard.md", "Dashboard"),
		"- " + vault.ToTopicWikiLink("demo-repo", "wiki/index/Concept Index.md", "Concept Index"),
		"- " + vault.ToTopicWikiLink("demo-repo", "wiki/index/Source Index.md", "Source Index"),
		"",
		"## Research gaps",
		"",
		"- Expand the starter wiki into architecture-level articles for the main subsystems.",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(topic.TopicPath, "CLAUDE.md"), []byte(legacyClaude), 0o644); err != nil {
		t.Fatalf("write legacy CLAUDE.md: %v", err)
	}

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	claudeContent := readFile(t, filepath.Join(topic.TopicPath, "CLAUDE.md"))
	if strings.Contains(claudeContent, "Generated codebase knowledge topic for `repo`") {
		t.Fatalf("legacy generated CLAUDE content was not replaced:\n%s", claudeContent)
	}
	if strings.Contains(claudeContent, "wiki/concepts/Kodebase - Directory Map.md") {
		t.Fatalf("legacy codebase links were not removed from CLAUDE.md:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "## Corpus inventory") {
		t.Fatalf("expected migrated CLAUDE.md to use the current scaffold template:\n%s", claudeContent)
	}
	if got := strings.Count(claudeContent, "<!-- kb:codebase:start -->"); got != 1 {
		t.Fatalf("expected exactly one managed block after migration, got %d\n%s", got, claudeContent)
	}
}

func TestWriteVaultPreservesManualClaudeContentAndUpdatesManagedBlock(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	if err := os.MkdirAll(topic.TopicPath, 0o755); err != nil {
		t.Fatalf("create topic path: %v", err)
	}

	manualClaude := strings.Join([]string{
		"# Demo Repo",
		"",
		"**Domain:** `demo-repo`",
		"",
		"Manual curated introduction.",
		"",
		"## Existing corpus",
		"- manual note",
		"",
		"<!-- kb:codebase:start -->",
		"old managed block",
		"<!-- kb:codebase:end -->",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(topic.TopicPath, "CLAUDE.md"), []byte(manualClaude), 0o644); err != nil {
		t.Fatalf("write manual CLAUDE.md: %v", err)
	}

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	claudeContent := readFile(t, filepath.Join(topic.TopicPath, "CLAUDE.md"))
	if !strings.Contains(claudeContent, "Manual curated introduction.") {
		t.Fatalf("manual CLAUDE.md content was lost:\n%s", claudeContent)
	}
	if strings.Contains(claudeContent, "old managed block") {
		t.Fatalf("stale managed block was not replaced:\n%s", claudeContent)
	}
	if got := strings.Count(claudeContent, "<!-- kb:codebase:start -->"); got != 1 {
		t.Fatalf("expected exactly one managed block start marker, got %d\n%s", got, claudeContent)
	}
	if !strings.Contains(claudeContent, vault.ToTopicWikiLink("demo-repo", vault.GetWikiIndexPath(vault.CodebaseDashboardTitle), vault.CodebaseDashboardTitle)) {
		t.Fatalf("managed block missing codebase dashboard link:\n%s", claudeContent)
	}
}

func TestWriteVaultAppendsManagedBlockWhenClaudeHasNoMarkers(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)
	if err := os.MkdirAll(topic.TopicPath, 0o755); err != nil {
		t.Fatalf("create topic path: %v", err)
	}

	manualClaude := strings.Join([]string{
		"# Demo Repo",
		"",
		"**Domain:** `demo-repo`",
		"",
		"Manual curated introduction without managed markers.",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(topic.TopicPath, "CLAUDE.md"), []byte(manualClaude), 0o644); err != nil {
		t.Fatalf("write manual CLAUDE.md: %v", err)
	}

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	claudeContent := readFile(t, filepath.Join(topic.TopicPath, "CLAUDE.md"))
	if !strings.Contains(claudeContent, "Manual curated introduction without managed markers.") {
		t.Fatalf("manual CLAUDE.md content was lost:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "<!-- kb:codebase:start -->") || !strings.Contains(claudeContent, "<!-- kb:codebase:end -->") {
		t.Fatalf("expected managed block markers to be appended:\n%s", claudeContent)
	}
}

func TestWriteVaultRejectsInvalidRenderedDocument(t *testing.T) {
	t.Parallel()

	topic, graph, _, baseFiles := testWriteVaultInputs(t)
	badDocument := models.RenderedDocument{
		Kind:         models.DocWiki,
		ManagedArea:  models.AreaWikiConcept,
		RelativePath: "wiki/codebase/concepts/Broken.md",
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

	if _, err := topicpkg.Info(writableTopic.VaultPath, writableTopic.Slug); err != nil {
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

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
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
