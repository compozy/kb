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

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/firecrawl"
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

// TestApplyRecaptureKeepsLongerBody: recapture keeps the longer of the old
// and refetched bodies (spec §7 stage 4): a longer good refetch replaces the
// body and clears the classify banks of the state row so the next classify
// judges it; a refetch that is not longer leaves the old body, which is
// quarantined when it still fails and kept when it no longer does.
func TestApplyRecaptureKeepsLongerBody(t *testing.T) {
	t.Parallel()
	long := strings.Repeat("Typed judges answer questions with calibrated probabilities. ", 60)
	tests := []struct {
		name       string
		oldBody    string
		fresh      firecrawl.ScrapeResult
		action     string
		wantBody   string
		wantReason string
	}{
		{name: "longer good refetch replaces the body", oldBody: "Loading...\n", fresh: firecrawl.ScrapeResult{Markdown: long, StatusCode: 200}, action: DidRecaptured, wantBody: long},
		{name: "shorter refetch and a still thin capture quarantine", oldBody: "Loading the page, please wait.\n", fresh: firecrawl.ScrapeResult{Markdown: "x", StatusCode: 200}, action: DidQuarantined, wantReason: "thin"},
		{name: "longer refetch that still fails quarantines", oldBody: "Loading...\n", fresh: firecrawl.ScrapeResult{Markdown: "# Not Found\n\nThe page does not exist anymore.\n", Title: "Not Found", StatusCode: 404}, action: DidQuarantined, wantReason: "error_page"},
		{name: "shorter refetch keeps a good old body", oldBody: long, fresh: firecrawl.ScrapeResult{Markdown: "short", StatusCode: 200}, action: DidKept, wantBody: long},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := newEnv(t)
			rel := "raw/articles/p.md"
			e.write(t, rel, map[string]any{"title": "Typed judges", "source_kind": "article", "source_url": "https://example.com/posts/typed-judges"}, tt.oldBody)
			doc, err := corpus.ReadDocument(e.root, rel, corpus.KindSource)
			if err != nil {
				t.Fatal(err)
			}
			meta := e.s.StateMeta(classify.SourceBanks()...)
			meta.Banks["custom"] = "7"
			if _, err := e.s.Writer.Apply(doc, map[string]any{"quality": "thin"}, meta); err != nil {
				t.Fatal(err)
			}
			item := e.add(t, review.Item{Queue: review.QueueRecapture, Purpose: "quality", Subject: rel, Question: "code:thin", Action: map[string]any{"reason": "thin"}})
			fresh := tt.fresh
			deps := Deps{Scrape: func(context.Context, string, firecrawl.ScrapeOptions) (*firecrawl.ScrapeResult, error) {
				return &fresh, nil
			}}
			result, err := Apply(context.Background(), e.s, deps, item, Accept)
			if err != nil {
				t.Fatalf("Apply: %v", err)
			}
			if result.Action != tt.action {
				t.Fatalf("action = %s (%s), want %s", result.Action, result.Message, tt.action)
			}
			if tt.action == DidQuarantined {
				values := e.frontmatter(t, result.Path)
				if values["triage_reason"] != tt.wantReason {
					t.Fatalf("quarantined with %v, want %s", values["triage_reason"], tt.wantReason)
				}
				return
			}
			stored, err := corpus.ReadDocument(e.root, rel, corpus.KindSource)
			if err != nil {
				t.Fatal(err)
			}
			if stored.Body != tt.wantBody {
				t.Fatalf("body = %q", stored.Body)
			}
			row, ok := e.s.State.Get(rel)
			if !ok {
				t.Fatal("state row missing")
			}
			if tt.action == DidRecaptured {
				for _, bank := range classify.SourceBanks() {
					if _, set := row.Banks[bank.ID]; set {
						t.Fatalf("classify bank %s still recorded after recapture: %v", bank.ID, row.Banks)
					}
				}
				if row.Banks["custom"] != "7" || row.BodyHash != corpus.BodyHash(long) {
					t.Fatalf("row = %+v", row)
				}
				reopened, err := corpus.OpenState(e.root)
				if err != nil {
					t.Fatal(err)
				}
				persisted, _ := reopened.Get(rel)
				if persisted == nil || corpus.Unclassified(stored, persisted, "", map[string]string{"quality": classify.SourceBanks()[1].Version}) != true {
					t.Fatalf("a recaptured source must read as unclassified: %+v", persisted)
				}
			}
		})
	}
}

// TestApplyRestoreReportsManualRepairs: a restore whose touched file changed
// since the quarantine returns every entry left for manual repair with its
// file, line or key and the removed text (spec §7.1).
func TestApplyRestoreReportsManualRepairs(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	e.write(t, "raw/articles/a.md", map[string]any{"title": "A", "source_kind": "article"}, "body\n")
	index := filepath.Join(e.root, "wiki", "index", "Source Index.md")
	if err := os.MkdirAll(filepath.Dir(index), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(index, []byte("# Source Index\n\n- [[a]] — first source\n- [[b]] — second\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reject := e.add(t, review.Item{Queue: review.QueueGate, Purpose: "relevance", Subject: "raw/articles/a.md", Question: "role", Action: map[string]any{"reason": "off_topic"}})
	if result, err := Apply(context.Background(), e.s, Deps{}, reject, Reject); err != nil || result.Action != DidQuarantined {
		t.Fatalf("reject = %+v, %v", result, err)
	}
	data, err := os.ReadFile(index)
	if err != nil || strings.Contains(string(data), "[[a]]") {
		t.Fatalf("index after quarantine = %q, %v", data, err)
	}
	if err := os.WriteFile(index, append(data, []byte("- [[c]] — added later\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	accept := e.add(t, review.Item{Queue: review.QueueGate, Purpose: "relevance", Subject: "raw/articles/a.md", Question: "role:again", Action: map[string]any{"reason": "off_topic", "quarantined": true}})
	result, err := Apply(context.Background(), e.s, Deps{}, accept, Accept)
	if err != nil || result.Action != DidRestored {
		t.Fatalf("accept = %+v, %v", result, err)
	}
	if len(result.Manual) != 1 {
		t.Fatalf("manual = %+v", result.Manual)
	}
	manual := result.Manual[0]
	if manual.File != "wiki/index/Source Index.md" || !strings.Contains(manual.Before, "[[a]] — first source") || manual.LineOrKey == "" {
		t.Fatalf("manual entry = %+v", manual)
	}
	lines := strings.Join(manual.Lines(), "\n")
	if !strings.Contains(lines, "manual repair: wiki/index/Source Index.md") || !strings.Contains(lines, "- [[a]] — first source") {
		t.Fatalf("rendered = %q", lines)
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
