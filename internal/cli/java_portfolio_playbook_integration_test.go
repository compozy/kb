//go:build integration

package cli

import (
	"encoding/json"
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

	dryRunStdout, dryRunStderr := runCLIWithStreams(t,
		"ingest", "codebase", repoRoot,
		"--topic", topicSlug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--log-format", "json",
		"--dry-run",
	)

	dryRunPayload := decodeJSONMap(t, []byte(dryRunStdout))
	assertCodebaseIngestContractShape(t, dryRunPayload)
	assertCodebaseIngestContractSemantics(t, dryRunPayload, true)

	var dryRunResult codebaseIngestResult
	if err := json.Unmarshal([]byte(dryRunStdout), &dryRunResult); err != nil {
		t.Fatalf("unmarshal dry-run payload: %v\n%s", err, dryRunStdout)
	}
	assertJavaCodebaseSummary(t, dryRunResult.Summary, 6, 10)
	if got := strings.Join(dryRunResult.Summary.SelectedAdapters, ","); !strings.Contains(strings.ToLower(got), "javaadapter") {
		t.Fatalf("selected adapters = %#v, want java adapter", dryRunResult.Summary.SelectedAdapters)
	}

	parseCompletedDryRun := findJSONStageCompletedEvent(t, dryRunStderr, "parse")
	if got := eventFieldInt(t, parseCompletedDryRun, "java_files_processed"); got < 1 {
		t.Fatalf("java_files_processed = %d, want >= 1", got)
	}
	if got := eventFieldInt(t, parseCompletedDryRun, "java_fallback_count"); got < 0 {
		t.Fatalf("java_fallback_count = %d, want >= 0", got)
	}
	if got := eventFieldInt(t, parseCompletedDryRun, "java_unresolved_count"); got < 0 {
		t.Fatalf("java_unresolved_count = %d, want >= 0", got)
	}

	fullRunStdout, fullRunStderr := runCLIWithStreams(t,
		"ingest", "codebase", repoRoot,
		"--topic", topicSlug,
		"--vault", vaultRoot,
		"--progress", "never",
		"--log-format", "json",
	)

	fullRunPayload := decodeJSONMap(t, []byte(fullRunStdout))
	assertCodebaseIngestContractShape(t, fullRunPayload)
	assertCodebaseIngestContractSemantics(t, fullRunPayload, false)

	var fullRunResult codebaseIngestResult
	if err := json.Unmarshal([]byte(fullRunStdout), &fullRunResult); err != nil {
		t.Fatalf("unmarshal full-run payload: %v\n%s", err, fullRunStdout)
	}
	assertJavaCodebaseSummary(t, fullRunResult.Summary, 6, 10)

	parseCompletedFullRun := findJSONStageCompletedEvent(t, fullRunStderr, "parse")
	if got := eventFieldInt(t, parseCompletedFullRun, "java_files_processed"); got < 1 {
		t.Fatalf("java_files_processed = %d, want >= 1", got)
	}

	issues := runCLIJSON[[]models.LintIssue](t,
		"lint", topicSlug,
		"--vault", vaultRoot,
		"--format", "json",
	)
	if len(issues) != 0 {
		t.Fatalf("lint issues = %#v, want none", issues)
	}
}
