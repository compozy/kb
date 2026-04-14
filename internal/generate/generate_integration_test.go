//go:build integration

package generate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

func TestGenerateIntegrationBuildsVaultFromFixtureRepository(t *testing.T) {
	t.Parallel()

	fixtureRoot := filepath.Join("testdata", "fixture-go-repo")
	outputRoot := filepath.Join(t.TempDir(), "vault")
	generator := newRunner()

	summary, err := generator.Generate(context.Background(), models.GenerateOptions{
		RootPath:  fixtureRoot,
		VaultPath: outputRoot,
		TopicSlug: "fixture-go-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", summary.FilesScanned)
	}
	if summary.FilesParsed != 2 {
		t.Fatalf("FilesParsed = %d, want 2", summary.FilesParsed)
	}
	if summary.FilesSkipped != 0 {
		t.Fatalf("FilesSkipped = %d, want 0", summary.FilesSkipped)
	}
	if summary.SymbolsExtracted != 4 {
		t.Fatalf("SymbolsExtracted = %d, want 4", summary.SymbolsExtracted)
	}
	if summary.RawDocumentsWritten != 9 {
		t.Fatalf("RawDocumentsWritten = %d, want 9", summary.RawDocumentsWritten)
	}
	if summary.WikiDocumentsWritten != 10 {
		t.Fatalf("WikiDocumentsWritten = %d, want 10", summary.WikiDocumentsWritten)
	}
	if summary.IndexDocumentsWritten != 3 {
		t.Fatalf("IndexDocumentsWritten = %d, want 3", summary.IndexDocumentsWritten)
	}

	expectedPaths := []string{
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "main.go.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "internal", "greeter", "greeter.go.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "symbols", "hello--internal-greeter-greeter-go-l6.md"),
		filepath.Join(summary.TopicPath, filepath.FromSlash(vault.GetWikiConceptPath("Codebase Overview"))),
		filepath.Join(summary.TopicPath, filepath.FromSlash(vault.GetWikiIndexPath(vault.CodebaseDashboardTitle))),
		filepath.Join(summary.TopicPath, "bases", "module-health.base"),
		filepath.Join(summary.TopicPath, "CLAUDE.md"),
	}

	for _, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}

	logContent, err := os.ReadFile(filepath.Join(summary.TopicPath, "log.md"))
	if err != nil {
		t.Fatalf("read log.md: %v", err)
	}
	if got := string(logContent); !containsAll(got,
		"## [",
		"ingest | codebase (2 files, 4 symbols)",
	) {
		t.Fatalf("expected codebase ingest log entry, got:\n%s", got)
	}
}

func TestGenerateIntegrationScansRepositoryNestedInsideVaultButOutsideTopic(t *testing.T) {
	t.Parallel()

	vaultRoot := t.TempDir()
	repoRoot := filepath.Join(vaultRoot, ".resources", "nested-go-repo")
	if err := os.MkdirAll(filepath.Join(repoRoot, "internal", "greeter"), 0o755); err != nil {
		t.Fatalf("create nested repo: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module example.com/nested-go-repo\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "main.go"), []byte(strings.Join([]string{
		"package main",
		"",
		"import \"example.com/nested-go-repo/internal/greeter\"",
		"",
		"func main() {",
		"\tgreeter.Hello()",
		"}",
		"",
	}, "\n")), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "internal", "greeter", "greeter.go"), []byte(strings.Join([]string{
		"package greeter",
		"",
		"func Hello() string {",
		"\treturn \"hello\"",
		"}",
		"",
	}, "\n")), 0o644); err != nil {
		t.Fatalf("write greeter.go: %v", err)
	}

	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath:  repoRoot,
		VaultPath: vaultRoot,
		TopicSlug: "nested-go-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", summary.FilesScanned)
	}
	if summary.FilesParsed != 2 {
		t.Fatalf("FilesParsed = %d, want 2", summary.FilesParsed)
	}
	if _, err := os.Stat(filepath.Join(vaultRoot, "nested-go-repo", "raw", "codebase", "files", "main.go.md")); err != nil {
		t.Fatalf("expected nested-vault ingest output: %v", err)
	}
}

func TestGenerateIntegrationBuildsVaultFromRustWorkspace(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeRustWorkspaceFixture(t, repoRoot)

	outputRoot := filepath.Join(t.TempDir(), "vault")
	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath:  repoRoot,
		VaultPath: outputRoot,
		TopicSlug: "fixture-rust-workspace",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 3 {
		t.Fatalf("FilesScanned = %d, want 3", summary.FilesScanned)
	}
	if summary.FilesParsed != 3 {
		t.Fatalf("FilesParsed = %d, want 3", summary.FilesParsed)
	}
	if !containsAll(strings.Join(summary.DetectedLanguages, ","), "rust") {
		t.Fatalf("expected rust in detected languages, got %#v", summary.DetectedLanguages)
	}
	if summary.SymbolsExtracted == 0 {
		t.Fatalf("SymbolsExtracted = %d, want > 0", summary.SymbolsExtracted)
	}

	expectedPaths := []string{
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "core", "src", "lib.rs.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "core", "src", "util.rs.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "app", "src", "lib.rs.md"),
	}
	for _, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}
}

func writeRustWorkspaceFixture(t *testing.T, repoRoot string) {
	t.Helper()

	writeFixtureFile(t, repoRoot, "Cargo.toml", strings.Join([]string{
		"[workspace]",
		`members = ["crates/core", "crates/app"]`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/Cargo.toml", strings.Join([]string{
		"[package]",
		`name = "openfang-core"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/src/lib.rs", strings.Join([]string{
		"pub mod util;",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/src/util.rs", strings.Join([]string{
		"pub fn helper() {}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/app/Cargo.toml", strings.Join([]string{
		"[package]",
		`name = "openfang-app"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/app/src/lib.rs", strings.Join([]string{
		"use openfang_core::util::helper;",
		"",
		"pub fn run() {",
		"\thelper();",
		"}",
		"",
	}, "\n"))
}

func writeFixtureFile(t *testing.T, rootPath, relativePath, content string) {
	t.Helper()

	absolutePath := filepath.Join(rootPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}

	return true
}
