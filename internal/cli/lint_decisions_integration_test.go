//go:build integration

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/review"
)

// TestCLIIntegrationLintReadsDecisionRecords runs `kb lint` over a scaffolded
// topic whose `.decisions/` records, topic.yaml and kb-owned frontmatter were
// written the way decision-backed commands write them, and asserts every
// decision-aware lint kind surfaces without any model configured.
func TestCLIIntegrationLintReadsDecisionRecords(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "")
	vaultRoot := t.TempDir()
	topic := scaffoldTopicForIntegration(t, vaultRoot, "decisions-lint", "Decisions Lint", "knowledge")
	root := topic.RootPath

	source := func(title, scraped string, extra map[string]any) map[string]any {
		values := map[string]any{
			"title":       title,
			"type":        "source",
			"stage":       "raw",
			"domain":      "knowledge",
			"source_kind": "article",
			"scraped":     scraped,
			"tags":        []string{"knowledge", "raw", "article"},
		}
		for key, value := range extra {
			values[key] = value
		}
		return values
	}
	writeMarkdownDocument(t, root, "raw/articles/fresh-finding.md", source("Fresh Finding", "2026-05-02", map[string]any{
		"affects":   []string{"[[Retrieval]]"},
		"relevance": "core",
		"triage":    "kept",
	}), "A finding that changes the retrieval article. See [[junk-page]].\n")
	writeMarkdownDocument(t, root, "raw/articles/off-topic.md", source("Off Topic", "2026-04-01", map[string]any{
		"relevance": "off_topic",
		"triage":    "kept",
		"related":   []string{"[[Nowhere]]"},
	}), "Unrelated material.\n")
	writeMarkdownDocument(t, root, "raw/_quarantine/articles/junk-page.md", source("Junk Page", "2026-04-01", map[string]any{
		"triage":        "quarantined",
		"triage_reason": "error_page",
	}), "Not found.\n")
	writeMarkdownDocument(t, root, "wiki/concepts/Retrieval.md", map[string]any{
		"title":   "Retrieval",
		"type":    "wiki",
		"stage":   "compiled",
		"domain":  "knowledge",
		"tags":    []string{"knowledge", "wiki", "concept"},
		"created": "2026-04-01",
		"updated": "2026-04-15",
		"sources": []string{"[[fresh-finding]]"},
	}, "Retrieval overview.\n")
	writeMarkdownDocument(t, root, "wiki/concepts/Ranking.md", map[string]any{
		"title":     "Ranking",
		"type":      "wiki",
		"stage":     "stub",
		"domain":    "knowledge",
		"tags":      []string{"knowledge", "wiki", "stub"},
		"created":   "2026-04-01",
		"updated":   "2026-04-01",
		"sources":   []string{},
		"criterion": "Documents that explain how results are ordered.",
		"summary":   "A user note that kb never wrote.",
	}, "")
	writeMarkdownDocument(t, root, "wiki/index/Dashboard.md", map[string]any{
		"title":   "Dashboard",
		"type":    "index",
		"domain":  "knowledge",
		"updated": "2026-04-15",
	}, "[[Retrieval]]\n[[Ranking]]\n")

	if err := contract.SaveContractDraft(root, &contract.Contract{Purpose: "Retrieval research.", Core: []string{"retrieval methods"}}); err != nil {
		t.Fatalf("SaveContractDraft: %v", err)
	}

	stub, err := corpus.ReadDocument(root, "wiki/concepts/Ranking.md", corpus.KindArticle)
	if err != nil {
		t.Fatalf("ReadDocument: %v", err)
	}
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}
	if err := store.Put(corpus.StateRow{
		Path:     stub.Path,
		BodyHash: stub.BodyHash,
		Banks:    map[string]string{"classify": questions.MustLoad("classify").Version, "concept": questions.MustLoad("concept").Version},
		Written:  map[string]string{"criterion": corpus.ValueHash(stub.Criterion())},
	}); err != nil {
		t.Fatalf("Put: %v", err)
	}

	queue := review.Open(root, nil)
	for _, item := range []review.Item{
		{Queue: review.QueueContradiction, Purpose: "link", Subject: "raw/articles/fresh-finding.md", Target: "wiki/concepts/Retrieval.md", Question: "relation_1", Probability: 0.88},
		{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/fresh-finding.md", Target: "wiki/concepts/Ranking.md", Question: "should_link_2", Probability: 0.7},
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/off-topic.md", Question: "role", Probability: 0.9},
		{Queue: review.QueueRecapture, Purpose: "quality", Subject: "raw/articles/off-topic.md", Question: "thin_or_boilerplate", Probability: 0.8},
	} {
		if _, err := queue.Add(item); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}

	issues := runCLIJSON[[]models.LintIssue](t, "lint", topic.Slug, "--format", "json", "--vault", vaultRoot)

	for _, want := range []models.LintIssue{
		{Kind: models.LintIssueKindFrontmatterDeadLink, Severity: models.SeverityError, FilePath: "raw/articles/off-topic.md", Target: "Nowhere"},
		{Kind: models.LintIssueKindLinkToQuarantined, Severity: models.SeverityWarning, FilePath: "raw/articles/fresh-finding.md", Target: "junk-page"},
		{Kind: models.LintIssueKindNeedsCompile, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Retrieval.md", Target: "raw/articles/fresh-finding"},
		{Kind: models.LintIssueKindContradiction, Severity: models.SeverityWarning, FilePath: "raw/articles/fresh-finding.md", Target: "wiki/concepts/Retrieval.md"},
		{Kind: models.LintIssueKindPendingReview, Severity: models.SeverityInfo, FilePath: ".decisions/review.jsonl", Target: review.QueueLink},
		{Kind: models.LintIssueKindRecapturePending, Severity: models.SeverityInfo, FilePath: ".decisions/review.jsonl", Target: review.QueueRecapture},
		{Kind: models.LintIssueKindRemovePending, Severity: models.SeverityInfo, FilePath: ".decisions/review.jsonl", Target: review.QueueRemove},
		{Kind: models.LintIssueKindOffTopicKept, Severity: models.SeverityWarning, FilePath: "raw/articles/off-topic.md", Target: "relevance"},
		{Kind: models.LintIssueKindKeyConflict, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Ranking.md", Target: "summary"},
		{Kind: models.LintIssueKindCriterionMissing, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Retrieval.md", Target: "criterion"},
		{Kind: models.LintIssueKindContractMissing, Severity: models.SeverityWarning, FilePath: "topic.yaml", Target: "contract"},
		{Kind: models.LintIssueKindContractDraftPending, Severity: models.SeverityInfo, FilePath: "topic.yaml", Target: "contract_draft"},
		{Kind: models.LintIssueKindUnclassified, Severity: models.SeverityInfo, FilePath: ".decisions/state.jsonl", Target: "3"},
	} {
		assertHasLintIssue(t, issues, want)
	}
	for _, issue := range issues {
		switch {
		case issue.Kind == models.LintIssueKindDeadLink && issue.Target == "junk-page":
			t.Fatalf("link to a quarantined source reported as dead: %#v", issue)
		case issue.Kind == models.LintIssueKindStale && issue.FilePath == "wiki/concepts/Retrieval.md":
			t.Fatalf("affected pair reported as stale instead of needs-compile: %#v", issue)
		case issue.Kind == models.LintIssueKindFormat && issue.FilePath == "wiki/concepts/Ranking.md":
			t.Fatalf("stub article rejected: %#v", issue)
		case issue.Kind == models.LintIssueKindVocabularyMissing:
			t.Fatalf("topic with articles reported vocabulary-missing: %#v", issue)
		}
	}

	// lint is read-only: the records it read are unchanged and nothing new
	// appears under .decisions/.
	entries, err := os.ReadDir(filepath.Join(root, ".decisions"))
	if err != nil {
		t.Fatalf("read .decisions: %v", err)
	}
	for _, entry := range entries {
		if name := entry.Name(); name != "state.jsonl" && name != review.QueueFile {
			t.Fatalf("lint wrote %s under .decisions/", name)
		}
	}
}
