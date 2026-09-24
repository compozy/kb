package actions

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/gate"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

var testNow = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

func TestLabelVerdict(t *testing.T) {
	t.Parallel()
	demote := review.Item{Queue: review.QueueLink, Action: map[string]any{"demote": true}}
	tests := []struct {
		item    review.Item
		verdict string
		want    string
	}{
		{review.Item{Queue: review.QueueGate}, Accept, review.VerdictPositive},
		{review.Item{Queue: review.QueueGate}, Reject, review.VerdictNegative},
		{review.Item{Queue: review.QueueSkip}, Accept, review.VerdictPositive},
		{review.Item{Queue: review.QueueRemove}, Accept, review.VerdictNegative},
		{review.Item{Queue: review.QueueRemove}, Reject, review.VerdictPositive},
		{review.Item{Queue: review.QueueRecapture}, Accept, review.VerdictNegative},
		{review.Item{Queue: review.QueueRecapture}, Reject, review.VerdictPositive},
		{review.Item{Queue: review.QueueLink}, Accept, review.VerdictPositive},
		{review.Item{Queue: review.QueueLink}, Reject, review.VerdictNegative},
		{demote, Accept, review.VerdictNegative},
		{demote, Reject, review.VerdictPositive},
		{review.Item{Queue: review.QueueContradiction}, Accept, review.VerdictPositive},
		{review.Item{Queue: review.QueueConceptProposal}, Reject, review.VerdictNegative},
		{review.Item{Queue: review.QueueOKFType}, Accept, review.VerdictPositive},
	}
	for _, tt := range tests {
		if got := LabelVerdict(tt.item, tt.verdict); got != tt.want {
			t.Errorf("%s %s (demote=%v) = %s, want %s", tt.verdict, tt.item.Queue, tt.item.Action["demote"], got, tt.want)
		}
	}
}

type env struct {
	root  string
	s     *session.Session
	store *review.Store
}

func newEnv(t *testing.T) *env {
	t.Helper()
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatal(err)
	}
	or := fakes.NewOpenRouter(nil, nil)
	t.Cleanup(or.Close)
	cfg := config.Default()
	cfg.OpenRouter.APIKey, cfg.OpenRouter.APIURL = "test-key", or.URL
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb review", Now: func() time.Time { return testNow }})
	if err != nil {
		t.Fatal(err)
	}
	return &env{root: info.RootPath, s: s, store: review.Open(info.RootPath, s.Now)}
}

func (e *env) write(t *testing.T, rel string, values map[string]any, body string) {
	t.Helper()
	content, err := frontmatter.Generate(values, body)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(e.root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (e *env) add(t *testing.T, item review.Item) review.Item {
	t.Helper()
	item.ID = review.ItemID(item.Queue, item.Subject, item.Target, item.Question)
	if _, err := e.store.Add(item); err != nil {
		t.Fatal(err)
	}
	stored, _, err := e.store.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	return stored
}

// fresh opens the review store again (a Store caches what it read).
func (e *env) fresh() *review.Store { return review.Open(e.root, nil) }

func (e *env) exists(rel string) bool {
	_, err := os.Stat(filepath.Join(e.root, filepath.FromSlash(rel)))
	return err == nil
}

func (e *env) frontmatter(t *testing.T, rel string) map[string]any {
	t.Helper()
	doc, err := corpus.ReadDocument(e.root, rel, corpus.KindSource)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Frontmatter
}

func TestApplyGateRejectThenAcceptRestores(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	e.write(t, "raw/articles/a.md", map[string]any{"title": "A", "source_kind": "article"}, "body\n")
	reject := e.add(t, review.Item{Queue: review.QueueGate, Purpose: "relevance", Subject: "raw/articles/a.md", Question: "role", Action: map[string]any{"reason": "off_topic"}})
	result, err := Apply(context.Background(), e.s, Deps{}, reject, Reject)
	if err != nil || result.Action != DidQuarantined || !e.exists("raw/_quarantine/articles/a.md") {
		t.Fatalf("reject = %+v, %v", result, err)
	}
	accept := e.add(t, review.Item{Queue: review.QueueGate, Purpose: "relevance", Subject: "raw/articles/a.md", Question: "role:again", Action: map[string]any{"reason": "off_topic", "quarantined": true}})
	result, err = Apply(context.Background(), e.s, Deps{}, accept, Accept)
	if err != nil || result.Action != DidRestored || !e.exists("raw/articles/a.md") {
		t.Fatalf("accept = %+v, %v", result, err)
	}
	values := e.frontmatter(t, "raw/articles/a.md")
	if values["triage"] != "kept" || values["triage_reason"] != nil {
		t.Fatalf("restored frontmatter = %v", values)
	}
	labels, _ := e.fresh().Labels()
	if len(labels) != 2 || labels[0].Verdict != review.VerdictNegative || labels[1].Verdict != review.VerdictPositive || labels[1].ItemID != accept.ID {
		t.Fatalf("labels = %+v", labels)
	}
	resolved, _, _ := e.fresh().Get(accept.ID)
	if _, err := Apply(context.Background(), e.s, Deps{}, resolved, Accept); !errors.Is(err, ErrNotPending) {
		t.Fatalf("a resolved item must not be applied again: %v", err)
	}
}

func TestApplyRemoveClosesSiblings(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	e.write(t, "raw/articles/r.md", map[string]any{"title": "Recipes", "source_kind": "article"}, "pasta\n")
	gated := e.add(t, review.Item{Queue: review.QueueGate, Purpose: "relevance", Subject: "raw/articles/r.md", Question: "role", Action: map[string]any{"reason": "off_topic"}})
	remove := e.add(t, review.Item{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/r.md", Question: "role", Action: map[string]any{"reason": "off_topic"}})

	result, err := Apply(context.Background(), e.s, Deps{}, remove, Accept)
	if err != nil || result.Action != DidQuarantined || !slices.Equal(result.Closed, []string{gated.ID}) {
		t.Fatalf("remove accept = %+v, %v", result, err)
	}
	sibling, _, _ := e.fresh().Get(gated.ID)
	if sibling.Status != review.StatusRejected {
		t.Fatalf("gate sibling = %+v", sibling)
	}
	labels, _ := e.fresh().Labels()
	if len(labels) != 1 || labels[0].Verdict != review.VerdictNegative || labels[0].Purpose != "relevance" {
		t.Fatalf("labels = %+v (siblings are not labelled)", labels)
	}
}

func TestApplyQueues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(t *testing.T, e *env) review.Item
		verdict string
		deps    Deps
		action  string
		check   func(t *testing.T, e *env, result Result)
	}{
		{
			name: "Should keep a remove item on reject without touching the file",
			setup: func(t *testing.T, e *env) review.Item {
				e.write(t, "raw/articles/k.md", map[string]any{"title": "K"}, "k\n")
				return e.add(t, review.Item{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/k.md", Question: "role"})
			},
			verdict: Reject, action: DidKept,
			check: func(t *testing.T, e *env, _ Result) {
				if !e.exists("raw/articles/k.md") || e.frontmatter(t, "raw/articles/k.md")["triage"] != nil {
					t.Fatal("reject must leave the source alone")
				}
			},
		},
		{
			name: "Should quarantine a recapture item without source_url",
			setup: func(t *testing.T, e *env) review.Item {
				e.write(t, "raw/articles/n.md", map[string]any{"title": "N"}, "n\n")
				return e.add(t, review.Item{Queue: review.QueueRecapture, Purpose: "quality", Subject: "raw/articles/n.md", Question: "code:thin", Action: map[string]any{"reason": "thin", "no_source": true}})
			},
			verdict: Accept, action: DidQuarantined,
			check: func(t *testing.T, e *env, result Result) {
				if result.Path != "raw/_quarantine/articles/n.md" {
					t.Fatalf("path = %s", result.Path)
				}
			},
		},
		{
			name: "Should dismiss a skipped item on reject",
			setup: func(t *testing.T, e *env) review.Item {
				return e.add(t, review.Item{Queue: review.QueueSkip, Purpose: "relevance", Subject: "sk-1", Question: "role_item", Action: map[string]any{"skipped_id": "sk-1"}})
			},
			verdict: Reject, action: DidDismissed,
		},
		{
			name: "Should rescue a skipped item on accept",
			setup: func(t *testing.T, e *env) review.Item {
				if err := gate.AppendSkipped(e.root, gate.SkippedRow{ID: "sk-2", URL: "https://example.com/x"}); err != nil {
					t.Fatal(err)
				}
				return e.add(t, review.Item{Queue: review.QueueSkip, Purpose: "relevance", Subject: "sk-2", Question: "role_item", Action: map[string]any{"skipped_id": "sk-2"}})
			},
			verdict: Accept, action: DidRescued,
			deps: Deps{Rescue: func(_ context.Context, row gate.SkippedRow) (string, error) {
				if row.URL != "https://example.com/x" {
					return "", errors.New("wrong row")
				}
				return "demo/raw/articles/x.md", nil
			}},
		},
		{
			name: "Should write contradicts on a contradiction accept",
			setup: func(t *testing.T, e *env) review.Item {
				e.write(t, "raw/articles/c.md", map[string]any{"title": "C"}, "c\n")
				return e.add(t, review.Item{Queue: review.QueueContradiction, Purpose: "link", Subject: "raw/articles/c.md", Target: "wiki/concepts/X.md", Question: "relation", Action: map[string]any{"relation": "contradicts", "target": "X"}})
			},
			verdict: Accept, action: DidWritten,
			check: func(t *testing.T, e *env, _ Result) {
				if got := frontmatter.GetStringSlice(e.frontmatter(t, "raw/articles/c.md"), "contradicts"); !slices.Equal(got, []string{"[[X]]"}) {
					t.Fatalf("contradicts = %v", got)
				}
			},
		},
		{
			name: "Should remove a kb-written relation on a demotion accept",
			setup: func(t *testing.T, e *env) review.Item {
				e.write(t, "raw/articles/d.md", map[string]any{"title": "D"}, "d\n")
				doc, err := corpus.ReadDocument(e.root, "raw/articles/d.md", corpus.KindSource)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := e.s.Writer.Apply(doc, map[string]any{"related": []string{"[[X]]", "[[Y]]"}}, e.s.StateMeta()); err != nil {
					t.Fatal(err)
				}
				return e.add(t, review.Item{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/d.md", Target: "wiki/concepts/X.md", Question: "demote:related", Action: map[string]any{"relation": "related", "target": "X", "demote": true}})
			},
			verdict: Accept, action: DidDemoted,
			check: func(t *testing.T, e *env, result Result) {
				if got := frontmatter.GetStringSlice(e.frontmatter(t, "raw/articles/d.md"), "related"); !slices.Equal(got, []string{"[[Y]]"}) {
					t.Fatalf("related = %v", got)
				}
				if result.Label.Verdict != review.VerdictNegative {
					t.Fatalf("label = %+v", result.Label)
				}
			},
		},
		{
			name: "Should leave a user-written relation on a demotion accept",
			setup: func(t *testing.T, e *env) review.Item {
				e.write(t, "raw/articles/u.md", map[string]any{"title": "U", "related": []string{"[[X]]"}}, "u\n")
				return e.add(t, review.Item{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/u.md", Target: "wiki/concepts/X.md", Question: "demote:related", Action: map[string]any{"relation": "related", "target": "X", "demote": true}})
			},
			verdict: Accept, action: DidNothing,
			check: func(t *testing.T, e *env, _ Result) {
				if got := frontmatter.GetStringSlice(e.frontmatter(t, "raw/articles/u.md"), "related"); !slices.Equal(got, []string{"[[X]]"}) {
					t.Fatalf("related = %v", got)
				}
			},
		},
		{
			name: "Should create a stub concept on a proposal accept",
			setup: func(t *testing.T, e *env) review.Item {
				return e.add(t, review.Item{Queue: review.QueueConceptProposal, Purpose: "concept", Subject: "wiki/concepts/Typed Judges.md", Question: "proposal", Action: map[string]any{"title": "Typed Judges", "criterion": "Documents about judges that answer typed questions."}})
			},
			verdict: Accept, action: DidCreated,
			check: func(t *testing.T, e *env, _ Result) {
				doc, err := corpus.ReadDocument(e.root, "wiki/concepts/Typed Judges.md", corpus.KindArticle)
				if err != nil || doc.Criterion() == "" || frontmatter.GetString(doc.Frontmatter, "stage") != "stub" {
					t.Fatalf("stub = %+v, %v", doc, err)
				}
			},
		},
		{
			name: "Should only record an OKF type verdict",
			setup: func(t *testing.T, e *env) review.Item {
				return e.add(t, review.Item{Queue: review.QueueOKFType, Purpose: "okf_type", Subject: "wiki/concepts/A.md", Target: "pattern", Question: "okf_type"})
			},
			verdict: Accept, action: DidRecorded,
			check: func(t *testing.T, _ *env, result Result) {
				if !strings.Contains(result.Message, "--type pattern") {
					t.Fatalf("message = %q", result.Message)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := newEnv(t)
			item := tt.setup(t, e)
			result, err := Apply(context.Background(), e.s, tt.deps, item, tt.verdict)
			if err != nil {
				t.Fatalf("Apply: %v", err)
			}
			if result.Action != tt.action {
				t.Fatalf("action = %s (%s), want %s", result.Action, result.Message, tt.action)
			}
			resolved, _, _ := e.fresh().Get(item.ID)
			labels, _ := e.fresh().Labels()
			if resolved.Status == review.StatusPending || len(labels) != 1 || labels[0].ItemID != item.ID || labels[0].Verdict != LabelVerdict(item, tt.verdict) {
				t.Fatalf("resolved = %+v labels = %+v", resolved, labels)
			}
			if tt.check != nil {
				tt.check(t, e, result)
			}
		})
	}
}

func TestApplyFailureLeavesItemPending(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	unknown := e.add(t, review.Item{Queue: "mystery", Subject: "x", Question: "q"})
	if _, err := Apply(context.Background(), e.s, Deps{}, unknown, Accept); err == nil {
		t.Fatal("unknown queues must fail")
	}
	skip := e.add(t, review.Item{Queue: review.QueueSkip, Subject: "sk-none", Question: "role_item", Action: map[string]any{"skipped_id": "sk-none"}})
	if _, err := Apply(context.Background(), e.s, Deps{}, skip, Accept); err == nil {
		t.Fatal("a rescue without a skipped row must fail")
	}
	if _, err := Apply(context.Background(), e.s, Deps{}, skip, "maybe"); err == nil {
		t.Fatal("invalid verdicts must fail")
	}
	pending, _ := e.fresh().Pending("")
	labels, _ := e.fresh().Labels()
	if len(pending) != 2 || len(labels) != 0 {
		t.Fatalf("pending = %d labels = %d", len(pending), len(labels))
	}
}

func TestInsertLink(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, body, text string
		start            int
		want, inserted   string
	}{
		{"Should insert at the recorded offset", "We use decision models.", "decision models", 7, "We use [[X|decision models]].", "decision models"},
		{"Should fall back to the first plain occurrence", "Intro. We use decision models.", "decision models", 3, "Intro. We use [[X|decision models]].", "decision models"},
		{"Should skip occurrences already inside a link", "[[Y|decision models]] and decision models", "decision models", 4, "[[Y|decision models]] and [[X|decision models]]", "decision models"},
		{"Should insert nothing when the text is gone", "Nothing here.", "decision models", 0, "Nothing here.", ""},
	}
	for _, tt := range tests {
		got, inserted := insertLink(tt.body, "X", tt.text, tt.start)
		if got != tt.want || inserted != tt.inserted {
			t.Errorf("%s: got %q (%q)", tt.name, got, inserted)
		}
	}
}
