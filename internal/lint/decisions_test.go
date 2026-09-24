package lint_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/lint"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/review"
)

func TestLintFrontmatterLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, topicPath string)
		want  []models.LintIssue
		deny  []models.LintIssue
	}{
		{
			name: "dead relation link is a frontmatter-dead-link",
			setup: func(t *testing.T, topicPath string) {
				values := sourceFrontmatter("Notes", "article", "2026-04-10")
				values["related"] = []string{"[[Nowhere]]"}
				writeMarkdownFile(t, topicPath, "raw/articles/notes.md", values, "# Notes\n")
			},
			want: []models.LintIssue{{Kind: models.LintIssueKindFrontmatterDeadLink, Severity: models.SeverityError, FilePath: "raw/articles/notes.md", Target: "Nowhere"}},
		},
		{
			name: "frontmatter link counts as an incoming edge",
			setup: func(t *testing.T, topicPath string) {
				values := sourceFrontmatter("Notes", "article", "2026-04-10")
				values["concepts"] = []string{"[[Linked Concept]]"}
				writeMarkdownFile(t, topicPath, "raw/articles/notes.md", values, "# Notes\n")
				writeMarkdownFile(t, topicPath, "wiki/concepts/Linked Concept.md", conceptFrontmatter("Linked Concept", "2026-04-11", []string{"[[notes]]"}), "# Linked\n")
			},
			deny: []models.LintIssue{{Kind: models.LintIssueKindOrphan, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Linked Concept.md"}},
		},
		{
			name: "frontmatter link to a quarantined source",
			setup: func(t *testing.T, topicPath string) {
				writeMarkdownFile(t, topicPath, "raw/_quarantine/articles/junk.md", sourceFrontmatter("Junk", "article", "2026-04-10"), "# Junk\n")
				writeMarkdownFile(t, topicPath, "wiki/concepts/Topic.md", conceptFrontmatter("Topic", "2026-04-11", []string{"[[raw/articles/junk]]"}), "# Topic\n")
				writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[Topic]]\n")
			},
			want: []models.LintIssue{{Kind: models.LintIssueKindLinkToQuarantined, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Topic.md", Target: "raw/articles/junk"}},
			deny: []models.LintIssue{{Kind: models.LintIssueKindMissingSource, Severity: models.SeverityError, FilePath: "wiki/concepts/Topic.md", Target: "raw/articles/junk"}},
		},
		{
			name: "article sources stay with missing-source",
			setup: func(t *testing.T, topicPath string) {
				writeMarkdownFile(t, topicPath, "wiki/concepts/Topic.md", conceptFrontmatter("Topic", "2026-04-11", []string{"[[Gone]]"}), "# Topic\n")
				writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[Topic]]\n")
			},
			want: []models.LintIssue{{Kind: models.LintIssueKindMissingSource, Severity: models.SeverityError, FilePath: "wiki/concepts/Topic.md", Target: "Gone"}},
			deny: []models.LintIssue{{Kind: models.LintIssueKindFrontmatterDeadLink, Severity: models.SeverityError, FilePath: "wiki/concepts/Topic.md", Target: "Gone"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			tt.setup(t, topicPath)
			issues := mustLint(t, topicPath)
			for _, want := range tt.want {
				assertHasIssue(t, issues, want)
			}
			for _, deny := range tt.deny {
				assertNoIssue(t, issues, deny)
			}
		})
	}
}

func TestLintNeedsCompileReplacesStaleForAffectedPairs(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	affecting := sourceFrontmatter("Affecting", "article", "2026-04-12")
	affecting["affects"] = []string{"[[Target]]"}
	writeMarkdownFile(t, topicPath, "raw/articles/affecting.md", affecting, "# Affecting\n")
	writeMarkdownFile(t, topicPath, "raw/articles/plain.md", sourceFrontmatter("Plain", "article", "2026-04-12"), "# Plain\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Target.md", conceptFrontmatter("Target", "2026-04-11", []string{"[[affecting]]", "[[plain]]"}), "# Target\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[Target]]\n")

	issues := mustLint(t, topicPath)
	assertHasIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindNeedsCompile, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Target.md", Target: "raw/articles/affecting"})
	assertNoIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindStale, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Target.md", Target: "affecting"})
	assertHasIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindStale, Severity: models.SeverityWarning, FilePath: "wiki/concepts/Target.md", Target: "plain"})
}

func TestLintOwnedKeyShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		value    any
		severity models.DiagnosticSeverity // "" means valid
	}{
		{name: "triage enum", key: "triage", value: "kept"},
		{name: "triage bad", key: "triage", value: "maybe", severity: models.SeverityError},
		{name: "triage_reason bad", key: "triage_reason", value: "boring", severity: models.SeverityError},
		{name: "genre option", key: "genre", value: "paper"},
		{name: "genre other", key: "genre", value: "other"},
		{name: "genre bad", key: "genre", value: "novel", severity: models.SeverityError},
		{name: "depth number", key: "depth", value: 2.5},
		{name: "depth out of range", key: "depth", value: 4, severity: models.SeverityError},
		{name: "depth text", key: "depth", value: "deep", severity: models.SeverityError},
		{name: "relevance enum", key: "relevance", value: "collected_on_purpose"},
		{name: "relevance bad", key: "relevance", value: "meh", severity: models.SeverityError},
		{name: "quality bad", key: "quality", value: "great", severity: models.SeverityError},
		{name: "locked bool", key: "locked", value: true},
		{name: "locked text", key: "locked", value: "yes", severity: models.SeverityError},
		{name: "summary long", key: "summary", value: strings.Repeat("a", 401), severity: models.SeverityWarning},
		{name: "summary ok", key: "summary", value: "Short summary."},
		{name: "entities too many", key: "entities", value: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p"}, severity: models.SeverityWarning},
		{name: "entities not strings", key: "entities", value: []any{"a", 1}, severity: models.SeverityError},
		{name: "questions too few", key: "questions", value: []string{"Why?"}, severity: models.SeverityWarning},
		{name: "questions ok", key: "questions", value: []string{"Why?", "How?", "When?"}},
		{name: "related quoted wikilinks", key: "related", value: []string{"[[notes]]"}},
		{name: "related plain text", key: "related", value: []string{"notes"}, severity: models.SeverityError},
		{name: "concepts scalar", key: "concepts", value: "[[notes]]", severity: models.SeverityError},
		{name: "supersedes nested list", key: "supersedes", value: []any{[]any{"notes"}}, severity: models.SeverityError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			writeMarkdownFile(t, topicPath, "raw/articles/notes.md", sourceFrontmatter("Notes", "article", "2026-04-10"), "# Notes\n")
			values := sourceFrontmatter("Shaped", "article", "2026-04-10")
			values[tt.key] = tt.value
			writeMarkdownFile(t, topicPath, "raw/articles/shaped.md", values, "# Shaped\n")

			issues := mustLint(t, topicPath)
			want := models.LintIssue{Kind: models.LintIssueKindFormat, Severity: tt.severity, FilePath: "raw/articles/shaped.md", Target: tt.key}
			if tt.severity == "" {
				for _, issue := range issues {
					if issue.FilePath == "raw/articles/shaped.md" && issue.Target == tt.key {
						t.Fatalf("unexpected issue %#v", issue)
					}
				}
				return
			}
			assertHasIssue(t, issues, want)
		})
	}
}

func TestLintAcceptsStubArticles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		stage   string
		sources []string
		wantErr string // target of the expected format error, "" for none
	}{
		{name: "stub with empty sources", stage: "stub", sources: []string{}},
		{name: "compiled needs sources", stage: "compiled", sources: []string{}, wantErr: "sources"},
		{name: "unknown stage", stage: "draft", sources: []string{"[[notes]]"}, wantErr: "stage"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			writeMarkdownFile(t, topicPath, "raw/articles/notes.md", sourceFrontmatter("Notes", "article", "2026-04-10"), "# Notes\n")
			values := conceptFrontmatter("Stub", "2026-04-11", tt.sources)
			values["stage"] = tt.stage
			values["criterion"] = "Documents that discuss stubs."
			writeMarkdownFile(t, topicPath, "wiki/concepts/Stub.md", values, "")
			writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[Stub]]\n")

			issues := mustLint(t, topicPath)
			var formatIssues []models.LintIssue
			for _, issue := range issues {
				if issue.Kind == models.LintIssueKindFormat {
					formatIssues = append(formatIssues, issue)
				}
			}
			if tt.wantErr == "" {
				if len(formatIssues) != 0 {
					t.Fatalf("format issues = %#v, want none", formatIssues)
				}
				return
			}
			assertHasIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindFormat, Severity: models.SeverityError, FilePath: "wiki/concepts/Stub.md", Target: tt.wantErr})
		})
	}
}

func TestLintReviewQueueKinds(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/notes.md", sourceFrontmatter("Notes", "article", "2026-04-10"), "# Notes\n")
	store := review.Open(topicPath, nil)
	items := []review.Item{
		{Queue: review.QueueContradiction, Purpose: "link", Subject: "raw/articles/notes.md", Target: "wiki/concepts/A.md", Question: "relation_a", Probability: 0.9},
		{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/notes.md", Target: "wiki/concepts/B.md", Question: "should_link_b", Probability: 0.7},
		{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/notes.md", Target: "wiki/concepts/C.md", Question: "should_link_c", Probability: 0.65},
		{Queue: review.QueueRecapture, Purpose: "quality", Subject: "raw/articles/notes.md", Question: "thin_or_boilerplate", Probability: 0.8},
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/other.md", Question: "role", Probability: 0.85},
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/resolved.md", Question: "role", Probability: 0.9},
	}
	for _, item := range items {
		if _, err := store.Add(item); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}
	if _, err := store.Resolve(review.ItemID(review.QueueRemove, "raw/articles/resolved.md", "", "role"), review.StatusRejected); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	issues := mustLint(t, topicPath)
	assertHasIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindContradiction, Severity: models.SeverityWarning, FilePath: "raw/articles/notes.md", Target: "wiki/concepts/A.md"})
	assertIssueMessage(t, issues, models.LintIssueKindPendingReview, review.QueueLink, "2 pending item(s)")
	assertIssueMessage(t, issues, models.LintIssueKindRecapturePending, review.QueueRecapture, "1 pending item(s)")
	assertIssueMessage(t, issues, models.LintIssueKindRemovePending, review.QueueRemove, "1 pending item(s)")
	for _, issue := range issues {
		if issue.Kind == models.LintIssueKindPendingReview && issue.Target == review.QueueContradiction {
			t.Fatalf("contradictions are itemized, not counted: %#v", issue)
		}
	}
}

func TestLintOffTopicKept(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	kept := sourceFrontmatter("Kept", "article", "2026-04-10")
	kept["relevance"] = "off_topic"
	kept["triage"] = "kept"
	writeMarkdownFile(t, topicPath, "raw/articles/kept.md", kept, "# Kept\n")
	quarantined := sourceFrontmatter("Gone", "article", "2026-04-10")
	quarantined["relevance"] = "off_topic"
	quarantined["triage"] = "quarantined"
	writeMarkdownFile(t, topicPath, "raw/_quarantine/articles/gone.md", quarantined, "# Gone\n")

	issues := mustLint(t, topicPath)
	assertHasIssue(t, issues, models.LintIssue{Kind: models.LintIssueKindOffTopicKept, Severity: models.SeverityWarning, FilePath: "raw/articles/kept.md", Target: "relevance"})
	for _, issue := range issues {
		if issue.Kind == models.LintIssueKindOffTopicKept && issue.FilePath != "raw/articles/kept.md" {
			t.Fatalf("unexpected off-topic-kept %#v", issue)
		}
	}
}

func TestLintWorkflowKindsOnlyForAdoptedTopics(t *testing.T) {
	t.Parallel()

	accepted := &contract.Contract{Purpose: "Systems design notes.", Core: []string{"distributed systems"}}
	workflowKinds := []models.LintIssueKind{
		models.LintIssueKindContractMissing,
		models.LintIssueKindVocabularyMissing,
		models.LintIssueKindCriterionMissing,
		models.LintIssueKindUnclassified,
	}

	tests := []struct {
		name      string
		adopt     bool
		contract  *contract.Contract
		draft     *contract.Contract
		article   bool
		wantKinds []models.LintIssueKind
		denyKinds []models.LintIssueKind
	}{
		{name: "untouched topic keeps structural lint", denyKinds: workflowKinds},
		{
			name:      "adopted topic without contract or articles",
			adopt:     true,
			wantKinds: []models.LintIssueKind{models.LintIssueKindContractMissing, models.LintIssueKindVocabularyMissing, models.LintIssueKindUnclassified},
			denyKinds: []models.LintIssueKind{models.LintIssueKindCriterionMissing, models.LintIssueKindContractDraftPending},
		},
		{
			name:      "accepted contract and an article without criterion",
			contract:  accepted,
			article:   true,
			wantKinds: []models.LintIssueKind{models.LintIssueKindCriterionMissing, models.LintIssueKindUnclassified},
			denyKinds: []models.LintIssueKind{models.LintIssueKindContractMissing, models.LintIssueKindVocabularyMissing},
		},
		{
			name:      "draft pending adopts the topic",
			draft:     accepted,
			wantKinds: []models.LintIssueKind{models.LintIssueKindContractDraftPending, models.LintIssueKindContractMissing},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			writeMarkdownFile(t, topicPath, "raw/articles/notes.md", sourceFrontmatter("Notes", "article", "2026-04-10"), "# Notes\n")
			if tt.article {
				writeMarkdownFile(t, topicPath, "wiki/concepts/Topic.md", conceptFrontmatter("Topic", "2026-04-11", []string{"[[notes]]"}), "# Topic\n")
				writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[Topic]]\n")
			}
			if tt.adopt {
				writeRecord(t, topicPath, "receipts.jsonl", "")
			}
			if tt.contract != nil {
				if err := contract.SetContract(topicPath, tt.contract); err != nil {
					t.Fatalf("SetContract: %v", err)
				}
			}
			if tt.draft != nil {
				if err := contract.SaveContractDraft(topicPath, tt.draft); err != nil {
					t.Fatalf("SaveContractDraft: %v", err)
				}
			}

			issues := mustLint(t, topicPath)
			for _, kind := range tt.wantKinds {
				if !hasKind(issues, kind) {
					t.Fatalf("missing %s in %#v", kind, issues)
				}
			}
			for _, kind := range tt.denyKinds {
				if hasKind(issues, kind) {
					t.Fatalf("unexpected %s in %#v", kind, issues)
				}
			}
		})
	}
}

func TestLintUnclassifiedReasons(t *testing.T) {
	t.Parallel()

	accepted := &contract.Contract{Purpose: "Systems design notes.", Core: []string{"distributed systems"}}
	classify := questions.MustLoad("classify")
	relevance := questions.MustLoad("relevance")

	tests := []struct {
		name       string
		row        func(doc *corpus.Document) *corpus.StateRow
		wantReason string // "" means classified
	}{
		{name: "no row", row: func(*corpus.Document) *corpus.StateRow { return nil }, wantReason: "no state record"},
		{
			name: "current row",
			row: func(doc *corpus.Document) *corpus.StateRow {
				return &corpus.StateRow{BodyHash: doc.BodyHash, Contract: accepted.Hash(), Banks: map[string]string{"classify": classify.Version, "relevance": relevance.Version}}
			},
		},
		{
			name: "body changed",
			row: func(doc *corpus.Document) *corpus.StateRow {
				return &corpus.StateRow{BodyHash: "old", Contract: accepted.Hash(), Banks: map[string]string{"classify": classify.Version}}
			},
			wantReason: "body changed",
		},
		{
			name: "contract changed",
			row: func(doc *corpus.Document) *corpus.StateRow {
				return &corpus.StateRow{BodyHash: doc.BodyHash, Contract: "other", Banks: map[string]string{"classify": classify.Version}}
			},
			wantReason: "contract changed",
		},
		{
			name: "bank changed",
			row: func(doc *corpus.Document) *corpus.StateRow {
				return &corpus.StateRow{BodyHash: doc.BodyHash, Contract: accepted.Hash(), Banks: map[string]string{"classify": classify.Version, "relevance": "2000-01-01.1"}}
			},
			wantReason: "question bank changed",
		},
		{
			name: "row without classification",
			row: func(doc *corpus.Document) *corpus.StateRow {
				return &corpus.StateRow{BodyHash: doc.BodyHash, Contract: accepted.Hash(), Banks: map[string]string{"link": "x"}}
			},
			wantReason: "not classified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			writeMarkdownFile(t, topicPath, "raw/articles/notes.md", sourceFrontmatter("Notes", "article", "2026-04-10"), "# Notes\n")
			writeMarkdownFile(t, topicPath, "raw/codebase/files/main.go.md", sourceFrontmatter("main.go", string(models.SourceKindCodebaseFile), "2026-04-10"), "# main\n")
			if err := contract.SetContract(topicPath, accepted); err != nil {
				t.Fatalf("SetContract: %v", err)
			}
			doc, err := corpus.ReadDocument(topicPath, "raw/articles/notes.md", corpus.KindSource)
			if err != nil {
				t.Fatalf("ReadDocument: %v", err)
			}
			if row := tt.row(doc); row != nil {
				row.Path = doc.Path
				putStateRow(t, topicPath, *row)
			}

			issues := mustLint(t, topicPath)
			var found *models.LintIssue
			for index := range issues {
				if issues[index].Kind == models.LintIssueKindUnclassified {
					found = &issues[index]
				}
			}
			if tt.wantReason == "" {
				if found != nil {
					t.Fatalf("unexpected unclassified issue %#v", *found)
				}
				return
			}
			if found == nil {
				t.Fatalf("missing unclassified issue in %#v", issues)
			}
			if found.Severity != models.SeverityInfo || found.Target != "1" || !strings.Contains(found.Message, tt.wantReason+": 1") || !strings.Contains(found.Message, "raw/articles/notes.md") {
				t.Fatalf("unclassified issue = %#v, want reason %q for raw/articles/notes.md only", *found, tt.wantReason)
			}
		})
	}
}

func TestLintUnclassifiedAggregatesLargeTopics(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	for index := range 13 {
		name := "note-" + string(rune('a'+index))
		writeMarkdownFile(t, topicPath, "raw/articles/"+name+".md", sourceFrontmatter(name, "article", "2026-04-10"), "# "+name+"\n")
	}
	writeRecord(t, topicPath, "receipts.jsonl", "")

	issues := mustLint(t, topicPath)
	count := 0
	for _, issue := range issues {
		if issue.Kind != models.LintIssueKindUnclassified {
			continue
		}
		count++
		if issue.Target != "13" || !strings.Contains(issue.Message, "(+3 more)") {
			t.Fatalf("aggregated issue = %#v", issue)
		}
	}
	if count != 1 {
		t.Fatalf("unclassified issues = %d, want one aggregated issue", count)
	}
}

func TestLintKeyConflicts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		written  map[string]string
		conflict []string
	}{
		{name: "kb-written values", written: map[string]string{"summary": corpus.ValueHash("Short summary."), "genre": corpus.ValueHash("paper")}},
		{name: "edited value", written: map[string]string{"summary": corpus.ValueHash("Another summary."), "genre": corpus.ValueHash("paper")}, conflict: []string{"summary"}},
		{name: "never written", written: map[string]string{"genre": corpus.ValueHash("paper")}, conflict: []string{"summary"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicPath := newTestTopic(t)
			values := sourceFrontmatter("Notes", "article", "2026-04-10")
			values["summary"] = strings.Repeat("x", 450)
			values["genre"] = "paper"
			values["ingest_batch"] = "file-2026-04-10-abc"
			writeMarkdownFile(t, topicPath, "raw/articles/notes.md", values, "# Notes\n")
			if len(tt.conflict) == 0 {
				// A kb-written summary within limits.
				values["summary"] = "Short summary."
				writeMarkdownFile(t, topicPath, "raw/articles/notes.md", values, "# Notes\n")
			}
			doc, err := corpus.ReadDocument(topicPath, "raw/articles/notes.md", corpus.KindSource)
			if err != nil {
				t.Fatalf("ReadDocument: %v", err)
			}
			putStateRow(t, topicPath, corpus.StateRow{Path: doc.Path, BodyHash: doc.BodyHash, Banks: map[string]string{"classify": "x"}, Written: tt.written})

			issues := mustLint(t, topicPath)
			got := make([]string, 0)
			for _, issue := range issues {
				if issue.Kind == models.LintIssueKindKeyConflict {
					got = append(got, issue.Target)
				}
				// A user-owned key is not held to kb's schema.
				if issue.Kind == models.LintIssueKindFormat && issue.Target == "summary" {
					t.Fatalf("user-owned summary validated as kb-owned: %#v", issue)
				}
			}
			if strings.Join(got, ",") != strings.Join(tt.conflict, ",") {
				t.Fatalf("key conflicts = %v, want %v", got, tt.conflict)
			}
		})
	}
}

func TestSaveReportListsDecisionSections(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	reportPath, err := lint.SaveReport(topicPath, []models.LintIssue{{
		Kind:     models.LintIssueKindNeedsCompile,
		Severity: models.SeverityWarning,
		FilePath: "wiki/concepts/Target.md",
		Target:   "raw/articles/affecting",
		Message:  "recompile",
	}}, time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("SaveReport: %v", err)
	}
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	for _, section := range []string{"NEEDS COMPILE (1)", "DEAD FRONTMATTER LINKS (0)", "LINKS TO QUARANTINED SOURCES (0)", "UNCLASSIFIED DOCUMENTS (0)", "PENDING REMOVAL (0)"} {
		if !strings.Contains(string(content), section) {
			t.Fatalf("report missing %q:\n%s", section, content)
		}
	}
}

func assertNoIssue(t *testing.T, issues []models.LintIssue, deny models.LintIssue) {
	t.Helper()
	for _, issue := range issues {
		if issue.Kind == deny.Kind && issue.FilePath == deny.FilePath && (deny.Target == "" || issue.Target == deny.Target) {
			t.Fatalf("unexpected issue %#v in %#v", issue, issues)
		}
	}
}

func assertIssueMessage(t *testing.T, issues []models.LintIssue, kind models.LintIssueKind, target, fragment string) {
	t.Helper()
	for _, issue := range issues {
		if issue.Kind == kind && issue.Target == target {
			if issue.Severity != models.SeverityInfo || !strings.Contains(issue.Message, fragment) {
				t.Fatalf("issue %#v, want info with %q", issue, fragment)
			}
			return
		}
	}
	t.Fatalf("missing %s/%s in %#v", kind, target, issues)
}

func hasKind(issues []models.LintIssue, kind models.LintIssueKind) bool {
	for _, issue := range issues {
		if issue.Kind == kind {
			return true
		}
	}
	return false
}

func writeRecord(t *testing.T, topicPath, name, content string) {
	t.Helper()
	path := filepath.Join(topicPath, ".decisions", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func putStateRow(t *testing.T, topicPath string, row corpus.StateRow) {
	t.Helper()
	store, err := corpus.OpenState(topicPath)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}
	if err := store.Put(row); err != nil {
		t.Fatalf("Put: %v", err)
	}
}
