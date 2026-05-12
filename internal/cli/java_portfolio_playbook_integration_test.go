//go:build integration

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestCLIIntegrationJavaPortfolioPlaybookCommandsAndSemantics(t *testing.T) {
	vaultRoot := t.TempDir()
	repoRoot := t.TempDir()
	writeJavaMultiModuleCodebaseFixture(t, repoRoot)

	const (
		topicSlug   = "java-portfolio-adoption"
		topicTitle  = "Java Portfolio Adoption"
		topicDomain = "java"
	)

	topic := runCLIJSON[models.TopicInfo](t,
		"topic", "new", topicSlug, topicTitle, topicDomain,
		"--vault", vaultRoot,
	)
	if topic.Slug != topicSlug {
		t.Fatalf("topic slug = %q, want %q", topic.Slug, topicSlug)
	}

	dryRunResult := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", repoRoot,
		"--topic", topicSlug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--log-format", "json",
		"--dry-run",
	)
	if !dryRunResult.Summary.DryRun {
		t.Fatalf("dry-run summary flag = %t, want true", dryRunResult.Summary.DryRun)
	}
	if got := strings.Join(dryRunResult.Summary.SelectedAdapters, ","); !strings.Contains(strings.ToLower(got), "javaadapter") {
		t.Fatalf("selected adapters = %#v, want java adapter", dryRunResult.Summary.SelectedAdapters)
	}

	fullRunResult := runCLIJSON[codebaseIngestResult](t,
		"ingest", "codebase", repoRoot,
		"--topic", topicSlug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--log-format", "json",
	)
	if fullRunResult.Summary.DryRun {
		t.Fatalf("full-run summary flag = %t, want false", fullRunResult.Summary.DryRun)
	}
	if got := fullRunResult.FilePath; got != filepath.ToSlash(filepath.Join(topicSlug, "raw", "codebase")) {
		t.Fatalf("full-run filePath = %q, want %q", got, filepath.ToSlash(filepath.Join(topicSlug, "raw", "codebase")))
	}
	if _, err := os.Stat(filepath.Join(vaultRoot, filepath.FromSlash(fullRunResult.FilePath))); err != nil {
		t.Fatalf("expected full-run to materialize codebase path %q: %v", fullRunResult.FilePath, err)
	}
}
