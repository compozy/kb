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
	kgenerate "github.com/compozy/kb/internal/generate"
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

func TestCLIIntegrationScaffoldIngestJavaWorkspaceCodebase(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "fixture-java-workspace", "Fixture Java Workspace", "java")
	repoRoot := t.TempDir()
	writeJavaMultiModuleCodebaseFixture(t, repoRoot)

	result := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", repoRoot,
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
	assertJavaCodebaseSummary(t, result.Summary, 6, 10)
	if got, want := result.Summary.FilesParsed, 6; got != want {
		t.Fatalf("FilesParsed = %d, want %d", got, want)
	}
	if got, want := strings.Join(result.Summary.SelectedAdapters, ","), "adapter.JavaAdapter"; !strings.EqualFold(got, want) {
		t.Fatalf("SelectedAdapters = %#v, want [%s]", result.Summary.SelectedAdapters, want)
	}

	for _, relativePath := range []string{
		"raw/codebase/files/shared-a/src/main/java/com/acme/shareda/Helper.java.md",
		"raw/codebase/files/shared-b/src/main/java/com/acme/sharedb/Helper.java.md",
		"raw/codebase/files/shared-b/src/main/java/com/acme/sharedb/Outer.java.md",
		"raw/codebase/files/shared-b/src/main/java/com/acme/sharedb/Tooling.java.md",
		"raw/codebase/files/app/src/main/java/com/acme/app/Runner.java.md",
		"raw/codebase/files/app/src/main/java/com/acme/app/AppMain.java.md",
	} {
		targetPath := filepath.Join(topic.RootPath, filepath.FromSlash(relativePath))
		if _, err := os.Stat(targetPath); err != nil {
			t.Fatalf("expected generated java codebase document %q: %v", targetPath, err)
		}
	}

	assertGeneratedSymbolContains(t, filepath.Join(topic.RootPath, "raw", "codebase", "symbols"), "helper")
	assertGeneratedSymbolContains(t, filepath.Join(topic.RootPath, "raw", "codebase", "symbols"), "outer")
	assertGeneratedSymbolContains(t, filepath.Join(topic.RootPath, "raw", "codebase", "symbols"), "assistnested")
	assertGeneratedSymbolContains(t, filepath.Join(topic.RootPath, "raw", "codebase", "symbols"), "appmain")

	issues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
	)
	if len(issues) != 0 {
		t.Fatalf("generated Java content should pass lint, found %#v", issues)
	}
}

func TestCLIIntegrationLintJavaDiagnosticsGovernanceWithControlledCounts(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "java-governance", "Java Governance", "java")

	writeMarkdownDocument(t, topic.RootPath, "raw/codebase/index/java.md", map[string]any{
		"title":                          "Language Snapshot: java",
		"type":                           "source",
		"stage":                          "raw",
		"domain":                         "java",
		"source_kind":                    "codebase-language-index",
		"scraped":                        "2026-04-12",
		"tags":                           []string{"java", "raw", "codebase", "language-index", "java"},
		"language":                       "java",
		"java_diagnostic_total_count":    4,
		"java_parse_error_count":         1,
		"java_resolution_fallback_count": 3,
	}, strings.Join([]string{
		"# Language Snapshot: java",
		"",
		"## Java Diagnostics",
		"- Total diagnostics: 4",
		"- JAVA_PARSE_ERROR: 1",
		"- JAVA_RESOLUTION_FALLBACK: 3",
	}, "\n"))

	defaultIssues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
	)
	assertHasLintIssue(t, defaultIssues, models.LintIssue{
		Kind:     models.LintIssueKindJavaDiagnosticGovernance,
		Severity: models.SeverityError,
		FilePath: "raw/codebase/index/java.md",
		Target:   "JAVA_PARSE_ERROR",
	})
	for _, issue := range defaultIssues {
		if issue.Kind == models.LintIssueKindJavaDiagnosticGovernance && issue.Target == "JAVA_RESOLUTION_FALLBACK" {
			t.Fatalf("fallback governance should remain disabled by default, got %#v", defaultIssues)
		}
	}

	thresholdIssues := runCLIJSON[[]models.LintIssue](t,
		"lint", topic.Slug,
		"--format", "json",
		"--vault", vaultRoot,
		"--java-max-fallback-warnings", "2",
	)
	assertHasLintIssue(t, thresholdIssues, models.LintIssue{
		Kind:     models.LintIssueKindJavaDiagnosticGovernance,
		Severity: models.SeverityError,
		FilePath: "raw/codebase/index/java.md",
		Target:   "JAVA_RESOLUTION_FALLBACK",
	})
}

func TestCLIIntegrationJavaIngestJSONContractStableAcrossModes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		dryRun       bool
		expectWrites bool
	}{
		{
			name:         "full-run",
			dryRun:       false,
			expectWrites: true,
		},
		{
			name:         "dry-run",
			dryRun:       true,
			expectWrites: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			vaultRoot := t.TempDir()
			repoRoot := t.TempDir()
			writeJavaMultiModuleCodebaseFixture(t, repoRoot)

			topicName := "fixture-java-contract"
			if testCase.dryRun {
				topicName += "-dry"
			}
			topic := scaffoldTopicForIntegration(t, vaultRoot, topicName, "Fixture Java Contract", "java")

			args := []string{
				"ingest", "codebase", repoRoot,
				"--topic", topic.Slug,
				"--vault", vaultRoot,
				"--progress", "never",
			}
			if testCase.dryRun {
				args = append(args, "--dry-run")
			}

			stdout, _ := runCLIWithStreams(t, args...)
			payload := decodeJSONMap(t, []byte(stdout))
			assertCodebaseIngestContractShape(t, payload)
			assertCodebaseIngestContractSemantics(t, payload, testCase.dryRun)

			var typedResult codebaseIngestResult
			if err := json.Unmarshal([]byte(stdout), &typedResult); err != nil {
				t.Fatalf("stdout did not contain typed JSON payload: %v\n%s", err, stdout)
			}
			if typedResult.Topic != topic.Slug {
				t.Fatalf("topic = %q, want %q", typedResult.Topic, topic.Slug)
			}
			if typedResult.SourceType != models.SourceKindCodebaseFile {
				t.Fatalf("sourceType = %q, want %q", typedResult.SourceType, models.SourceKindCodebaseFile)
			}
			if typedResult.Summary.DryRun != testCase.dryRun {
				t.Fatalf("summary.dryRun = %t, want %t", typedResult.Summary.DryRun, testCase.dryRun)
			}
			if testCase.expectWrites {
				assertJavaCodebaseSummary(t, typedResult.Summary, 6, 10)
			}
		})
	}
}

func TestCLIIntegrationJavaIngestJSONLogsIncludeTelemetry(t *testing.T) {
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "fixture-java-telemetry", "Fixture Java Telemetry", "java")
	repoRoot := t.TempDir()
	writeJavaMultiModuleCodebaseFixture(t, repoRoot)

	stdout, stderr := runCLIWithStreams(t,
		"ingest", "codebase", repoRoot,
		"--topic", topic.Slug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--log-format", "json",
	)

	var result codebaseIngestResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout did not contain JSON summary: %v\n%s", err, stdout)
	}
	if err := validateJavaCodebaseSummary(result.Summary, 6, 10); err != nil {
		t.Fatal(err)
	}

	parseCompleted := findJSONStageCompletedEvent(t, stderr, "parse")
	if got := eventFieldInt(t, parseCompleted, "java_files_processed"); got < 1 {
		t.Fatalf("java_files_processed = %d, want >= 1", got)
	}
	if got := eventFieldInt(t, parseCompleted, "java_parse_duration_millis"); got < 0 {
		t.Fatalf("java_parse_duration_millis = %d, want >= 0", got)
	}
	if got := eventFieldString(t, parseCompleted, "java_resolver_mode"); got != "deep" && got != "fallback" {
		t.Fatalf("java_resolver_mode = %q, want deep or fallback", got)
	}
	if got := eventFieldInt(t, parseCompleted, "java_fallback_count"); got < 0 {
		t.Fatalf("java_fallback_count = %d, want >= 0", got)
	}
	if got := eventFieldInt(t, parseCompleted, "java_unresolved_count"); got < 0 {
		t.Fatalf("java_unresolved_count = %d, want >= 0", got)
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

func assertGeneratedSymbolContains(t *testing.T, symbolsDirectory string, expectedFragment string) {
	t.Helper()

	entries, err := os.ReadDir(symbolsDirectory)
	if err != nil {
		t.Fatalf("read generated symbols in %q: %v", symbolsDirectory, err)
	}

	expected := strings.ToLower(expectedFragment)
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Name()), expected) {
			return
		}
	}

	t.Fatalf("expected generated symbol artifact containing %q in %q", expectedFragment, symbolsDirectory)
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

	stdout, _ := runCLIWithStreams(t, args...)
	return stdout
}

func runCLIWithStreams(t *testing.T, args ...string) (string, string) {
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

	return stdout.String(), stderr.String()
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

func findJSONStageCompletedEvent(t *testing.T, stderrOutput string, stage string) kgenerate.Event {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(stderrOutput), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event kgenerate.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("stderr line was not valid JSON event: %v\nline=%s", err, line)
		}
		if event.Kind == kgenerate.EventStageCompleted && event.Stage == stage {
			return event
		}
	}

	t.Fatalf("missing stage_completed event for stage %q in stderr:\n%s", stage, stderrOutput)
	return kgenerate.Event{}
}

func eventFieldInt(t *testing.T, event kgenerate.Event, key string) int {
	t.Helper()

	value, ok := event.Fields[key]
	if !ok {
		t.Fatalf("event field %q missing from %#v", key, event.Fields)
	}

	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case int64:
		return int(typed)
	default:
		t.Fatalf("event field %q has unsupported numeric type %T (%#v)", key, value, value)
	}

	return 0
}

func eventFieldString(t *testing.T, event kgenerate.Event, key string) string {
	t.Helper()

	value, ok := event.Fields[key]
	if !ok {
		t.Fatalf("event field %q missing from %#v", key, event.Fields)
	}

	typed, ok := value.(string)
	if !ok {
		t.Fatalf("event field %q has unsupported string type %T (%#v)", key, value, value)
	}

	return typed
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
