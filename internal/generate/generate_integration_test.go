//go:build integration

package generate

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestGenerateIntegrationBuildsVaultFromFixtureRepository(t *testing.T) {
	t.Parallel()

	fixtureRoot := filepath.Join("testdata", "fixture-go-repo")
	outputRoot := filepath.Join(t.TempDir(), "vault")
	generator := newRunner()
	generator.logger = slog.New(slog.NewTextHandler(io.Discard, nil))

	summary, err := generator.Generate(context.Background(), models.GenerateOptions{
		RootPath:   fixtureRoot,
		OutputPath: outputRoot,
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
		filepath.Join(summary.TopicPath, "wiki", "concepts", "Codebase Overview.md"),
		filepath.Join(summary.TopicPath, "wiki", "index", "Dashboard.md"),
		filepath.Join(summary.TopicPath, "bases", "module-health.base"),
		filepath.Join(summary.TopicPath, "CLAUDE.md"),
	}

	for _, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}
}
