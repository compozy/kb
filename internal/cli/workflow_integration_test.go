//go:build integration

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

func TestCLIIntegrationScaffoldAndIngestFiles(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "systems-design", "Systems Design", "systems")

	textSourcePath := filepath.Join(t.TempDir(), "latency-budget.txt")
	writeFile(t, textSourcePath, strings.Join([]string{
		"Latency Budget",
		"",
		"Keep p95 under 200ms for user-facing requests.",
	}, "\n"))

	textResult := runCLIJSON[models.IngestResult](t,
		"ingest", "file", textSourcePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)

	if textResult.FilePath != "systems-design/raw/articles/latency-budget.md" {
		t.Fatalf("text ingest filePath = %q, want systems-design/raw/articles/latency-budget.md", textResult.FilePath)
	}
	if textResult.SourceType != models.SourceKindDocument {
		t.Fatalf("text ingest sourceType = %q, want %q", textResult.SourceType, models.SourceKindDocument)
	}

	textFrontmatter, textBody := readMarkdownDocument(t, filepath.Join(vaultRoot, filepath.FromSlash(textResult.FilePath)))
	assertRawDocumentFrontmatter(t, textFrontmatter, map[string]any{
		"title":       "Latency Budget",
		"type":        "source",
		"stage":       "raw",
		"domain":      "systems",
		"source_kind": string(models.SourceKindDocument),
		"source_path": filepath.ToSlash(filepath.Clean(textSourcePath)),
	})
	if !strings.Contains(textBody, "Keep p95 under 200ms") {
		t.Fatalf("expected text body to contain ingested content, got:\n%s", textBody)
	}

	csvSourcePath := filepath.Join(t.TempDir(), "service-levels.csv")
	writeFile(t, csvSourcePath, strings.Join([]string{
		"service,latency_ms",
		"api,120",
		"worker,250",
	}, "\n"))

	csvResult := runCLIJSON[models.IngestResult](t,
		"ingest", "file", csvSourcePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)

	if csvResult.FilePath != "systems-design/raw/articles/service-levels.md" {
		t.Fatalf("csv ingest filePath = %q, want systems-design/raw/articles/service-levels.md", csvResult.FilePath)
	}

	csvFrontmatter, csvBody := readMarkdownDocument(t, filepath.Join(vaultRoot, filepath.FromSlash(csvResult.FilePath)))
	assertRawDocumentFrontmatter(t, csvFrontmatter, map[string]any{
		"title":       "Service Levels",
		"type":        "source",
		"stage":       "raw",
		"domain":      "systems",
		"source_kind": string(models.SourceKindDocument),
		"source_path": filepath.ToSlash(filepath.Clean(csvSourcePath)),
	})
	for _, fragment := range []string{
		"| service | latency_ms |",
		"| --- | --- |",
		"| api | 120 |",
		"| worker | 250 |",
	} {
		if !strings.Contains(csvBody, fragment) {
			t.Fatalf("expected CSV markdown body to contain %q, got:\n%s", fragment, csvBody)
		}
	}
}

func TestCLIIntegrationScaffoldIngestCodebaseAndInspect(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "fixture-go-repo", "Fixture Go Repo", "golang")
	codebasePath := filepath.Join("..", "generate", "testdata", "fixture-go-repo")

	result := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", codebasePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
	)

	if result.Topic != topic.Slug {
		t.Fatalf("codebase ingest topic = %q, want %q", result.Topic, topic.Slug)
	}
	if result.SourceType != models.SourceKindCodebaseFile {
		t.Fatalf("codebase ingest sourceType = %q, want %q", result.SourceType, models.SourceKindCodebaseFile)
	}
	if result.FilePath != "fixture-go-repo/raw/codebase" {
		t.Fatalf("codebase ingest filePath = %q, want fixture-go-repo/raw/codebase", result.FilePath)
	}
	if result.Summary.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", result.Summary.FilesScanned)
	}
	if result.Summary.SymbolsExtracted != 4 {
		t.Fatalf("SymbolsExtracted = %d, want 4", result.Summary.SymbolsExtracted)
	}

	for _, relativePath := range []string{
		"raw/codebase/files/main.go.md",
		"raw/codebase/files/internal/greeter/greeter.go.md",
		"raw/codebase/symbols/hello--internal-greeter-greeter-go-l6.md",
	} {
		targetPath := filepath.Join(topic.RootPath, filepath.FromSlash(relativePath))
		if _, err := os.Stat(targetPath); err != nil {
			t.Fatalf("expected generated codebase document %q: %v", targetPath, err)
		}
	}

	for _, relativePath := range []string{
		vault.GetTopicIndexPath(vault.TopicDashboardTitle),
		vault.GetTopicIndexPath(vault.TopicConceptIndexTitle),
		vault.GetTopicIndexPath(vault.TopicSourceIndexTitle),
	} {
		content, err := os.ReadFile(filepath.Join(topic.RootPath, filepath.FromSlash(relativePath)))
		if err != nil {
			t.Fatalf("read top-level topic index %q: %v", relativePath, err)
		}

		for _, codebaseTitle := range []string{
			vault.CodebaseDashboardTitle,
			vault.CodebaseConceptIndexTitle,
			vault.CodebaseSourceIndexTitle,
		} {
			link := vault.ToTopicWikiLink(topic.Slug, vault.GetWikiIndexPath(codebaseTitle), codebaseTitle)
			if !strings.Contains(string(content), link) {
				t.Fatalf("expected %q to bridge to %q, got:\n%s", relativePath, codebaseTitle, string(content))
			}
		}
	}

	stdout := runCLI(t,
		"inspect", "smells",
		"--format", "json",
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)

	var rows []map[string]any
	if err := json.Unmarshal([]byte(stdout), &rows); err != nil {
		t.Fatalf("inspect smells did not return JSON: %v\n%s", err, stdout)
	}
	if len(rows) > 0 {
		for _, key := range []string{"kind", "name", "source_path", "symbol_kind", "smells"} {
			if _, ok := rows[0][key]; !ok {
				t.Fatalf("inspect smells row missing key %q: %#v", key, rows[0])
			}
		}
	}
}

func TestCLIIntegrationScaffoldIngestRustWorkspaceCodebase(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "fixture-rust-workspace", "Fixture Rust Workspace", "rust")
	repoRoot := t.TempDir()
	writeRustCodebaseFixture(t, repoRoot)

	result := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", repoRoot,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
	)

	if result.Topic != topic.Slug {
		t.Fatalf("codebase ingest topic = %q, want %q", result.Topic, topic.Slug)
	}
	if result.Summary.FilesScanned != 3 {
		t.Fatalf("FilesScanned = %d, want 3", result.Summary.FilesScanned)
	}
	if got := strings.Join(result.Summary.DetectedLanguages, ","); !strings.Contains(got, "rust") {
		t.Fatalf("expected rust in detected languages, got %#v", result.Summary.DetectedLanguages)
	}
	if result.Summary.SymbolsExtracted == 0 {
		t.Fatalf("SymbolsExtracted = %d, want > 0", result.Summary.SymbolsExtracted)
	}

	for _, relativePath := range []string{
		"raw/codebase/files/crates/core/src/lib.rs.md",
		"raw/codebase/files/crates/core/src/util.rs.md",
		"raw/codebase/files/crates/app/src/lib.rs.md",
	} {
		targetPath := filepath.Join(topic.RootPath, filepath.FromSlash(relativePath))
		if _, err := os.Stat(targetPath); err != nil {
			t.Fatalf("expected generated rust codebase document %q: %v", targetPath, err)
		}
	}

	issues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
	)
	if len(issues) != 0 {
		t.Fatalf("generated Rust content should pass lint, found %#v", issues)
	}
}

func TestCLIIntegrationCodebaseBootstrapCreatesDefaultVaultFromExternalRepo(t *testing.T) {
	repoRoot := t.TempDir()
	writeGoCodebaseFixture(t, repoRoot)

	withWorkingDirectory(t, t.TempDir(), func() {
		result := runCLIJSON[codebaseIngestResult](t,
			"ingest", "codebase", repoRoot,
			"--topic", "chat-sdk",
			"--title", "Chat SDK",
			"--domain", "messaging",
			"--progress", "never",
		)

		expectedVaultRoot := filepath.Join(repoRoot, ".kb", "vault")
		expectedTopicRoot := filepath.Join(expectedVaultRoot, "chat-sdk")

		if result.Topic != "chat-sdk" {
			t.Fatalf("codebase ingest topic = %q, want chat-sdk", result.Topic)
		}
		if result.Title != "Chat SDK" {
			t.Fatalf("codebase ingest title = %q, want Chat SDK", result.Title)
		}
		if result.Summary.VaultPath != expectedVaultRoot {
			t.Fatalf("vault path = %q, want %q", result.Summary.VaultPath, expectedVaultRoot)
		}
		if result.Summary.TopicPath != expectedTopicRoot {
			t.Fatalf("topic path = %q, want %q", result.Summary.TopicPath, expectedTopicRoot)
		}

		for _, relativePath := range []string{
			"raw/codebase/files/main.go.md",
			"raw/codebase/files/internal/greeter/greeter.go.md",
			"CLAUDE.md",
			"log.md",
		} {
			targetPath := filepath.Join(expectedTopicRoot, filepath.FromSlash(relativePath))
			if _, err := os.Stat(targetPath); err != nil {
				t.Fatalf("expected bootstrapped topic artifact %q: %v", targetPath, err)
			}
		}
	})

	withWorkingDirectory(t, repoRoot, func() {
		info := runCLIJSON[models.TopicInfo](t,
			"topic", "info", "chat-sdk",
		)
		if info.Title != "Chat SDK" || info.Domain != "messaging" {
			t.Fatalf("topic info = %#v", info)
		}

		rows := runCLIJSON[[]map[string]any](t,
			"inspect", "complexity",
			"--topic", "chat-sdk",
			"--format", "json",
		)
		if len(rows) == 0 {
			t.Fatal("expected inspect complexity rows after bootstrapped ingest")
		}
	})
}

func TestCLIIntegrationCodebaseBootstrapSupportsDeprecatedOutputAlias(t *testing.T) {
	repoRoot := t.TempDir()
	writeGoCodebaseFixture(t, repoRoot)
	vaultRoot := filepath.Join(t.TempDir(), "shared-vault")

	result := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", repoRoot,
		"--topic", "legacy-output",
		"--output", vaultRoot,
		"--progress", "never",
	)

	if result.Summary.VaultPath != vaultRoot {
		t.Fatalf("vault path = %q, want %q", result.Summary.VaultPath, vaultRoot)
	}
	if _, err := os.Stat(filepath.Join(vaultRoot, "legacy-output", "CLAUDE.md")); err != nil {
		t.Fatalf("expected bootstrapped topic under legacy output alias: %v", err)
	}
}

func TestCLIIntegrationEmptyCodebaseIngestFailsWithoutWritesAndKeepsTopicDiscoverable(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "empty-codebase", "Empty Codebase", "docs")
	codebasePath := filepath.Join("..", "..", "docs")

	errOutput := runCLIError(t,
		"ingest", "codebase", codebasePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
	)
	if !strings.Contains(errOutput, "no supported source files found") {
		t.Fatalf("unexpected empty ingest error:\n%s", errOutput)
	}

	for _, relativePath := range []string{
		"raw/codebase/files",
		"raw/codebase/symbols",
	} {
		info, err := os.Stat(filepath.Join(topic.RootPath, filepath.FromSlash(relativePath)))
		if err != nil {
			t.Fatalf("stat %q: %v", relativePath, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", relativePath)
		}
	}

	info := runCLIJSON[models.TopicInfo](t,
		"topic", "info", topic.Slug,
		"--vault", vaultRoot,
	)
	if info.Slug != topic.Slug {
		t.Fatalf("topic info slug = %q, want %q", info.Slug, topic.Slug)
	}

	listOutput := runCLI(t, "topic", "list", "--vault", vaultRoot)
	if !strings.Contains(listOutput, topic.Slug) {
		t.Fatalf("topic list output missing %q:\n%s", topic.Slug, listOutput)
	}
}

func TestCLIIntegrationCodebaseDryRunDoesNotWriteManagedArtifacts(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "fixture-go-repo", "Fixture Go Repo", "golang")
	codebasePath := filepath.Join("..", "generate", "testdata", "fixture-go-repo")

	result := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", codebasePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--dry-run",
	)

	if !result.Summary.DryRun {
		t.Fatalf("DryRun = %t, want true", result.Summary.DryRun)
	}
	if result.Summary.RawDocumentsWritten != 0 || result.Summary.WikiDocumentsWritten != 0 || result.Summary.IndexDocumentsWritten != 0 {
		t.Fatalf("expected dry-run write counts to stay zero, got %#v", result.Summary)
	}

	for _, relativePath := range []string{
		"raw/codebase/files/main.go.md",
		filepath.ToSlash(vault.GetWikiConceptPath("Codebase Overview")),
		filepath.ToSlash(vault.GetWikiIndexPath(vault.CodebaseDashboardTitle)),
	} {
		if _, err := os.Stat(filepath.Join(topic.RootPath, filepath.FromSlash(relativePath))); !os.IsNotExist(err) {
			t.Fatalf("expected dry-run to avoid writing %q, stat err = %v", relativePath, err)
		}
	}
}

func TestCLIIntegrationGeneratedContentPassesLint(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "rewrite-qa", "Rewrite QA", "engineering")
	codebasePath := filepath.Join("..", "generate", "testdata", "fixture-go-repo")

	bookmarkPath := filepath.Join(t.TempDir(), "bookmarks.md")
	writeFile(t, bookmarkPath, strings.Join([]string{
		"# QA Links",
		"",
		"- [Go Tour](https://go.dev/tour/)",
		"- [Tree-sitter](https://tree-sitter.github.io/tree-sitter/)",
	}, "\n"))

	runCLIJSON[models.IngestResult](t,
		"ingest", "bookmarks", bookmarkPath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)
	runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", codebasePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
	)

	issues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
	)
	if len(issues) != 0 {
		t.Fatalf("generated content should pass lint, found %#v", issues)
	}
}

func TestCLIIntegrationScaffoldIngestAndLint(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "knowledge-ops", "Knowledge Ops", "knowledge")

	sourceOnePath := filepath.Join(t.TempDir(), "source-note.txt")
	writeFile(t, sourceOnePath, strings.Join([]string{
		"Source Note",
		"",
		"Primary note for the topic.",
	}, "\n"))
	sourceTwoPath := filepath.Join(t.TempDir(), "architecture-map.txt")
	writeFile(t, sourceTwoPath, strings.Join([]string{
		"Architecture Map",
		"",
		"Secondary note for the topic.",
	}, "\n"))

	runCLIJSON[models.IngestResult](t,
		"ingest", "file", sourceOnePath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)
	runCLIJSON[models.IngestResult](t,
		"ingest", "file", sourceTwoPath,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
	)

	today := time.Now().UTC().Format(frontmatter.DateLayout)
	writeMarkdownDocument(t, topic.RootPath, "wiki/concepts/Broken Concept.md", map[string]any{
		"title":   "Broken Concept",
		"type":    "wiki",
		"stage":   "compiled",
		"domain":  "knowledge",
		"tags":    []string{"knowledge", "wiki", "concept"},
		"created": today,
		"updated": today,
		"sources": []string{"Source Note"},
	}, "Broken concept.\n\nSee [[Missing Page]].\n")
	writeMarkdownDocument(t, topic.RootPath, "wiki/concepts/Orphan Concept.md", map[string]any{
		"title":   "Orphan Concept",
		"type":    "wiki",
		"stage":   "compiled",
		"domain":  "knowledge",
		"tags":    []string{"knowledge", "wiki", "concept"},
		"created": today,
		"updated": today,
		"sources": []string{"Architecture Map"},
	}, "This concept is intentionally unlinked.\n")
	writeMarkdownDocument(t, topic.RootPath, "wiki/index/Dashboard.md", map[string]any{
		"title":   "Dashboard",
		"type":    "index",
		"domain":  "knowledge",
		"updated": today,
	}, "[[knowledge-ops/wiki/concepts/Broken Concept|Broken Concept]]\n")

	issues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
	)

	assertHasLintIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindDeadLink,
		Severity: models.SeverityError,
		FilePath: "wiki/concepts/Broken Concept.md",
		Target:   "Missing Page",
	})
	assertHasLintIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindOrphan,
		Severity: models.SeverityWarning,
		FilePath: "wiki/concepts/Orphan Concept.md",
	})
}

func writeRustCodebaseFixture(t *testing.T, repoRoot string) {
	t.Helper()

	writeFile(t, filepath.Join(repoRoot, "Cargo.toml"), strings.Join([]string{
		"[workspace]",
		`members = ["crates/core", "crates/app"]`,
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "crates", "core", "Cargo.toml"), strings.Join([]string{
		"[package]",
		`name = "openfang-core"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "crates", "core", "src", "lib.rs"), strings.Join([]string{
		"pub mod util;",
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "crates", "core", "src", "util.rs"), strings.Join([]string{
		"pub fn helper() {}",
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "crates", "app", "Cargo.toml"), strings.Join([]string{
		"[package]",
		`name = "openfang-app"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "crates", "app", "src", "lib.rs"), strings.Join([]string{
		"use openfang_core::util::helper;",
		"",
		"pub fn run() {",
		"\thelper();",
		"}",
		"",
	}, "\n"))
}

func writeGoCodebaseFixture(t *testing.T, repoRoot string) {
	t.Helper()

	writeFile(t, filepath.Join(repoRoot, "go.mod"), "module example.com/chat-sdk\n\ngo 1.22\n")
	writeFile(t, filepath.Join(repoRoot, "main.go"), strings.Join([]string{
		"package main",
		"",
		"import \"example.com/chat-sdk/internal/greeter\"",
		"",
		"func main() {",
		"\tgreeter.Hello()",
		"}",
		"",
	}, "\n"))
	writeFile(t, filepath.Join(repoRoot, "internal", "greeter", "greeter.go"), strings.Join([]string{
		"package greeter",
		"",
		"func Hello() string {",
		"\treturn \"hello\"",
		"}",
		"",
	}, "\n"))
}

func scaffoldTopicForIntegration(t *testing.T, vaultRoot, slug, title, domain string) models.TopicInfo {
	t.Helper()

	return runCLIJSON[models.TopicInfo](t,
		"topic", "new", slug, title, domain,
		"--vault", vaultRoot,
	)
}

func runCLI(t *testing.T, args ...string) string {
	t.Helper()

	command := newRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext(%q) returned error: %v\nstderr:\n%s", strings.Join(args, " "), err, stderr.String())
	}

	return stdout.String()
}

func runCLIError(t *testing.T, args ...string) string {
	t.Helper()

	command := newRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatalf("ExecuteContext(%q) unexpectedly succeeded", strings.Join(args, " "))
	}

	return err.Error() + "\n" + stderr.String()
}

func runCLIJSON[T any](t *testing.T, args ...string) T {
	t.Helper()

	stdout := runCLI(t, args...)
	var payload T
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout)
	}

	return payload
}

func withWorkingDirectory(t *testing.T, directory string, fn func()) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("chdir to %q: %v", directory, err)
	}
	defer func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working directory to %q: %v", previous, err)
		}
	}()

	fn()
}

func readMarkdownDocument(t *testing.T, documentPath string) (map[string]any, string) {
	t.Helper()

	content, err := os.ReadFile(documentPath)
	if err != nil {
		t.Fatalf("read document %q: %v", documentPath, err)
	}

	values, body, err := frontmatter.Parse(string(content))
	if err != nil {
		t.Fatalf("parse frontmatter for %q: %v", documentPath, err)
	}

	return values, body
}

func writeMarkdownDocument(t *testing.T, rootPath, relativePath string, values map[string]any, body string) {
	t.Helper()

	document, err := frontmatter.Generate(values, body)
	if err != nil {
		t.Fatalf("generate frontmatter for %q: %v", relativePath, err)
	}

	targetPath := filepath.Join(rootPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatalf("create parent directory for %q: %v", targetPath, err)
	}
	if err := os.WriteFile(targetPath, []byte(document), 0o644); err != nil {
		t.Fatalf("write markdown document %q: %v", targetPath, err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directory for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %q: %v", path, err)
	}
}

func assertRawDocumentFrontmatter(t *testing.T, values map[string]any, expected map[string]any) {
	t.Helper()

	for key, want := range expected {
		if got := values[key]; got != want {
			t.Fatalf("frontmatter[%q] = %#v, want %#v", key, got, want)
		}
	}

	if scraped := frontmatter.GetString(values, "scraped"); strings.TrimSpace(scraped) == "" {
		t.Fatal("expected scraped frontmatter field to be set")
	}
}

func assertHasLintIssue(t *testing.T, issues []models.LintIssue, want models.LintIssue) {
	t.Helper()

	for _, issue := range issues {
		if issue.Kind != want.Kind {
			continue
		}
		if issue.Severity != want.Severity {
			continue
		}
		if issue.FilePath != want.FilePath {
			continue
		}
		if want.Target != "" && issue.Target != want.Target {
			continue
		}

		return
	}

	t.Fatalf("missing lint issue %#v in %#v", want, issues)
}
