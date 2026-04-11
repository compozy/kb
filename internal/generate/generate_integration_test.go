//go:build integration

package generate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
	"github.com/user/go-devstack/internal/vault"
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
		filepath.Join(summary.TopicPath, "wiki", "index", "Dashboard.md"),
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

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}

	return true
}
