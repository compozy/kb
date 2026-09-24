package gate

import (
	"context"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

var testNow = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

func TestNormalizeURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, a, b string
		same       bool
	}{
		{"Should ignore scheme case, www and trailing slash", "HTTPS://WWW.Example.com/Posts/A/", "http://example.com/Posts/A", true},
		{"Should drop fragments and tracking parameters", "https://example.com/a?utm_source=x&id=2&fbclid=1#top", "https://example.com/a?id=2", true},
		{"Should sort the remaining query", "https://example.com/a?b=2&a=1", "https://example.com/a?a=1&b=2", true},
		{"Should keep path case", "https://example.com/Posts/A", "https://example.com/posts/a", false},
		{"Should keep meaningful parameters", "https://example.com/a?id=1", "https://example.com/a?id=2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NormalizeURL(tt.a) == NormalizeURL(tt.b); got != tt.same {
				t.Fatalf("NormalizeURL(%q)=%q, NormalizeURL(%q)=%q, same=%v", tt.a, NormalizeURL(tt.a), tt.b, NormalizeURL(tt.b), got)
			}
		})
	}
	if NormalizeURL("file:///tmp/a.md") != "" || NormalizeURL("not a url") != "" {
		t.Fatal("non-http URLs must normalize to empty")
	}
}

func TestPlatformID(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"https://www.youtube.com/watch?v=abcdefghijk&t=10": "youtube:abcdefghijk",
		"https://youtu.be/abcdefghijk?si=track":            "youtube:abcdefghijk",
		"https://m.youtube.com/shorts/abcdefghijk":         "youtube:abcdefghijk",
		"https://www.youtube.com/@channel/videos":          "",
		"https://www.instagram.com/reel/SHORTCODE123/":     "instagram:SHORTCODE123",
		"https://instagram.com/p/Cabc123/?igsh=x":          "instagram:Cabc123",
		"https://example.com/p/abc":                        "",
	}
	for raw, want := range tests {
		if got := PlatformID(raw); got != want {
			t.Errorf("PlatformID(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestIndexFindsDuplicatesIncludingQuarantined(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeDoc(t, root, "raw/articles/a.md", map[string]any{"title": "A", "source_url": "https://example.com/a"}, "alpha body\n")
	writeDoc(t, root, "raw/_quarantine/youtube/v.md", map[string]any{"title": "V", "video_id": "abcdefghijk"}, "transcript\n")
	writeDoc(t, root, "raw/_quarantine/articles/b.md", map[string]any{"title": "B"}, "bravo body\n")
	c, err := corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewIndex(root, c.Sources())
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		ref         Ref
		path, by    string
		quarantined bool
	}{
		{"Should match a normalized URL", Ref{URL: "https://www.example.com/a/?utm_medium=x"}, "raw/articles/a.md", "url", false},
		{"Should match a quarantined video by id", Ref{URL: "https://youtu.be/abcdefghijk"}, "raw/_quarantine/youtube/v.md", "platform_id", true},
		{"Should match a quarantined body hash", Ref{BodyHash: corpus.BodyHash("bravo body\n")}, "raw/_quarantine/articles/b.md", "body_hash", true},
	}
	for _, tt := range tests {
		match, ok := index.Find(tt.ref)
		if !ok || match.Path != tt.path || match.By != tt.by || match.Quarantined != tt.quarantined {
			t.Errorf("%s: Find = %+v, %v", tt.name, match, ok)
		}
	}
	if _, ok := index.Find(Ref{URL: "https://example.com/c", BodyHash: corpus.BodyHash("")}); ok {
		t.Error("an unknown URL with an empty body must not match")
	}
	if !strings.Contains((Match{Path: "raw/x.md", Quarantined: true, By: "url"}).Describe(), "duplicate of quarantined raw/x.md") {
		t.Error("Describe must name quarantined matches")
	}
}

func TestPrefetchAction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		off, kept   float64
		want        string
		quarantine  float64
		fetchThresh float64
	}{
		{"Should skip at P(off_topic) >= 0.8", 0.85, 0.1, ActionSkip, 0.8, 0.5},
		{"Should fetch when kept roles sum >= 0.5", 0.3, 0.6, ActionFetch, 0.8, 0.5},
		{"Should review in between", 0.5, 0.3, ActionReview, 0.8, 0.5},
		{"Should skip before fetch on ties", 0.8, 0.9, ActionSkip, 0.8, 0.5},
	}
	for _, tt := range tests {
		if got := PrefetchAction(tt.off, tt.kept, tt.quarantine, tt.fetchThresh); got != tt.want {
			t.Errorf("%s: got %s", tt.name, got)
		}
	}
}

func TestNeedsRefetch(t *testing.T) {
	t.Parallel()
	flag := []quality.Flag{{Code: quality.ErrorPage}}
	tests := []struct {
		name        string
		flags       []quality.Flag
		words       int
		refetchable bool
		want        bool
	}{
		{"Should refetch a flagged URL source", flag, 900, true, true},
		{"Should refetch a short URL source", nil, 120, true, true},
		{"Should not refetch a long clean source", nil, 400, true, false},
		{"Should never refetch a non-URL source", flag, 10, false, false},
	}
	for _, tt := range tests {
		if got := NeedsRefetch(tt.flags, tt.words, tt.refetchable); got != tt.want {
			t.Errorf("%s: got %v", tt.name, got)
		}
	}
}

func TestOutcomeBanding(t *testing.T) {
	t.Parallel()
	judgment := func(pOff float64, nouls map[string]float64) classify.GateJudgment {
		return classify.GateJudgment{RoleAsked: true, RoleDecided: true, POffTopic: pOff, Quality: nouls, QualityDecided: true, ReceiptKey: "k"}
	}
	tests := []struct {
		name          string
		j             classify.GateJudgment
		qualityMode   string
		relevanceMode string
		triage        string
		reason        string
		shadow        bool
	}{
		{"Should keep a clean document", judgment(0.1, map[string]float64{classify.NoulPaywall: 0.1}), "apply", "apply", TriageKept, "", false},
		{"Should quarantine a quality noul at apply", judgment(0.1, map[string]float64{classify.NoulErrorPage: 0.9}), "apply", "shadow", TriageQuarantined, "error_page", false},
		{"Should review a quality noul in the review band", judgment(0.1, map[string]float64{classify.NoulThin: 0.6}), "apply", "apply", TriageReview, "thin", false},
		{"Should quarantine off-topic in apply mode", judgment(0.9, nil), "apply", "apply", TriageQuarantined, "off_topic", false},
		{"Should only review off-topic in shadow mode", judgment(0.95, nil), "apply", "shadow", TriageReview, "off_topic", true},
		{"Should review the relevance review band", judgment(0.6, nil), "apply", "apply", TriageReview, "off_topic", false},
		{"Should prefer quarantine over review across purposes", judgment(0.9, map[string]float64{classify.NoulThin: 0.6}), "apply", "apply", TriageQuarantined, "off_topic", false},
		{"Should keep a quality quarantine in shadow as review", judgment(0.1, map[string]float64{classify.NoulPaywall: 0.9}), "shadow", "shadow", TriageReview, "paywall", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := Outcome{Triage: TriageKept}
			out.mergeQuality(tt.j, 0.8, 0.5, tt.qualityMode)
			out.mergeRelevance(tt.j, 0.8, 0.5, tt.relevanceMode)
			if out.Triage != tt.triage || out.Reason != tt.reason || (len(out.Shadow) > 0) != tt.shadow {
				t.Fatalf("outcome = %+v", out)
			}
		})
	}
	undecided := Outcome{Triage: TriageKept}
	undecided.mergeRelevance(classify.GateJudgment{RoleAsked: true, POffTopic: 0.99}, 0.8, 0.5, "apply")
	if undecided.Triage != TriageKept {
		t.Fatal("an undecided role must never move a source")
	}
}

func TestReviewItemAndTally(t *testing.T) {
	t.Parallel()
	if _, ok := ReviewItem(Outcome{Triage: TriageKept}, "raw/a.md", "A"); ok {
		t.Fatal("kept sources get no review item")
	}
	item, ok := ReviewItem(Outcome{Triage: TriageQuarantined, Reason: "thin", Stage: StageQuality, Purpose: "quality", Question: "code:thin", Probability: 1}, "raw/a.md", "A")
	if !ok || item.Queue != review.QueueGate || item.Action["quarantined"] != true || item.Action["reason"] != "thin" || item.Subject != "raw/a.md" {
		t.Fatalf("item = %+v", item)
	}
	var tally Tally
	for _, triage := range []string{TriageKept, TriageReview, TriageQuarantined, TriageSkipped, TriageDuplicateSkipped, TriageKept} {
		tally.Add(Outcome{Triage: triage})
	}
	if got := tally.Line(); got != "kept 2, review 1, quarantined 1, skipped 1, duplicates 1, undecided 0" {
		t.Fatalf("Line = %q", got)
	}
	tally.Add(Outcome{Triage: TriageReview, Reason: ReasonUndecided})
	tally.Add(Outcome{Triage: TriageKept, Excluded: true})
	if got := tally.Line(); got != "kept 3, review 2, quarantined 1, skipped 1, duplicates 1, undecided 1, excluded 1 (decisions.exclude, not judged)" {
		t.Fatalf("Line with undecided and excluded = %q", got)
	}

	undecided, ok := ReviewItem(Outcome{
		Triage: TriageReview, Reason: ReasonUndecided, Stage: StageUndecided, Purpose: "relevance", Question: "role",
		Undecided: []string{"role:timeout"},
	}, "raw/b.md", "B")
	if !ok || undecided.Queue != review.QueueGate || undecided.Action["reason"] != ReasonOffTopic || undecided.Action["undecided"] != "role:timeout" {
		t.Fatalf("undecided item = %+v", undecided)
	}
}

// TestMergeUndecided: a gate judgment that could not be made puts the
// source in review with triage_reason `undecided`, never `kept`; answers not
// checked because of decisions.exclude are no judgment, and a real verdict
// wins (spec §2.3).
func TestMergeUndecided(t *testing.T) {
	t.Parallel()
	decidedQuality := map[string]float64{classify.NoulThin: 0.1, classify.NoulPaywall: 0.1, classify.NoulErrorPage: 0.1}
	tests := []struct {
		name     string
		j        classify.GateJudgment
		flagged  bool
		start    Outcome
		triage   string
		reason   string
		question string
	}{
		{
			name:   "decided judgment stays kept",
			j:      classify.GateJudgment{RoleAsked: true, RoleDecided: true, Quality: decidedQuality, QualityDecided: true},
			start:  Outcome{Triage: TriageKept},
			triage: TriageKept,
		},
		{
			name:   "undecided role goes to review",
			j:      classify.GateJudgment{RoleAsked: true, Quality: decidedQuality, QualityDecided: true, Undecided: []string{"role:timeout"}},
			start:  Outcome{Triage: TriageKept},
			triage: TriageReview, reason: ReasonUndecided, question: "role",
		},
		{
			name:   "undecided quality nouls go to review",
			j:      classify.GateJudgment{Undecided: []string{"thin_or_boilerplate:budget", "paywall_or_login:budget"}},
			start:  Outcome{Triage: TriageKept},
			triage: TriageReview, reason: ReasonUndecided, question: "thin_or_boilerplate",
		},
		{
			name:   "state too large is undecided too",
			j:      classify.GateJudgment{RoleAsked: true, Undecided: []string{"role:not_checked:state_too_large", "thin_or_boilerplate:not_checked:state_too_large"}},
			start:  Outcome{Triage: TriageKept},
			triage: TriageReview, reason: ReasonUndecided, question: "role",
		},
		{
			name:   "excluded answers are no judgment",
			j:      classify.GateJudgment{RoleAsked: true, Undecided: []string{"role:not_checked:excluded", "thin_or_boilerplate:not_checked:excluded"}},
			start:  Outcome{Triage: TriageKept},
			triage: TriageKept,
		},
		{
			name:    "code flag already decided quality",
			j:       classify.GateJudgment{Undecided: []string{"thin_or_boilerplate:timeout"}},
			flagged: true,
			start:   Outcome{Triage: TriageKept},
			triage:  TriageKept,
		},
		{
			name:   "a real review verdict wins",
			j:      classify.GateJudgment{RoleAsked: true, Undecided: []string{"role:timeout"}},
			start:  Outcome{Triage: TriageReview, Reason: "thin"},
			triage: TriageReview, reason: "thin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := tt.start
			out.mergeUndecided(tt.j, tt.flagged)
			if out.Triage != tt.triage || out.Reason != tt.reason || out.Question != tt.question {
				t.Fatalf("outcome = %+v", out)
			}
		})
	}
}

func TestTitleSimilarityAndNearCandidates(t *testing.T) {
	t.Parallel()
	if got := TitleSimilarity("Typed judges handbook", "Typed judges handbook (2nd edition)"); got < 0.6 {
		t.Fatalf("similarity = %.2f", got)
	}
	if got := TitleSimilarity("Typed judges", "Pasta recipes"); got != 0 {
		t.Fatalf("similarity = %.2f", got)
	}
	doc := &corpus.Document{Path: "raw/articles/new.md", Title: "Typed judges handbook (2nd edition)", Body: "judges", Frontmatter: map[string]any{"source_url": "https://example.com/v2"}}
	sameHost := &corpus.Document{Path: "raw/articles/old.md", Title: "Typed judges handbook", Frontmatter: map[string]any{"source_url": "https://www.example.com/v1"}}
	otherHost := &corpus.Document{Path: "raw/articles/other.md", Title: "Typed judges handbook", Frontmatter: map[string]any{"source_url": "https://other.org/v1"}}
	candidates := NearCandidates(doc, []*corpus.Document{doc, sameHost, otherHost}, nil)
	if len(candidates) != 1 || candidates[0].Doc != sameHost || candidates[0].Origin != "title" {
		t.Fatalf("candidates = %+v", candidates)
	}
}

// evalEnv is a topic with a session on a fake OpenRouter.
type evalEnv struct {
	root string
	s    *session.Session
	or   *fakes.OpenRouter
}

func newEvalEnv(t *testing.T, decide fakes.DecideFunc, flags session.Flags, yaml string) *evalEnv {
	t.Helper()
	return newEvalEnvWith(t, decide, flags, yaml, nil)
}

// newEvalEnvWith is newEvalEnv with a setup run on the topic root before the
// session opens (extra banks, fixtures).
func newEvalEnvWith(t *testing.T, decide fakes.DecideFunc, flags session.Flags, yaml string, setup func(root string)) *evalEnv {
	t.Helper()
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatal(err)
	}
	if setup != nil {
		setup(info.RootPath)
	}
	if yaml != "" {
		c := &contract.Contract{Purpose: "Decision models research.", Core: []string{"decision models"}, OutOfScope: []string{"cooking recipes"}}
		if err := contract.SetContract(info.RootPath, c); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(info.RootPath, "topic.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, []byte(yaml)...), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	or := fakes.NewOpenRouter(decide, nil)
	t.Cleanup(or.Close)
	cfg := config.Default()
	cfg.OpenRouter.APIKey, cfg.OpenRouter.APIURL = "test-key", or.URL
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb test", Flags: flags, Now: func() time.Time { return testNow }})
	if err != nil {
		t.Fatal(err)
	}
	return &evalEnv{root: info.RootPath, s: s, or: or}
}

func (e *evalEnv) fetched(title, sourceURL, body string) Fetched {
	return Fetched{
		Title: title, Markdown: body,
		Build: func(title, markdown string) (*corpus.Document, error) {
			return &corpus.Document{
				Path: "raw/articles/new.md", Title: title, Kind: corpus.KindSource, Body: markdown, BodyHash: corpus.BodyHash(markdown),
				Frontmatter: map[string]any{"title": title, "source_url": sourceURL, "source_kind": "article"},
			}, nil
		},
	}
}

func longBody(word string, n int) string {
	return strings.Repeat(word+" words make a sentence here. ", n/6) + "\n"
}

func TestEvaluateStages(t *testing.T) {
	t.Parallel()
	offTopic := func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "role" {
			return fakes.Pick("off_topic", 0.95)
		}
		return nil
	}
	tests := []struct {
		name   string
		decide fakes.DecideFunc
		flags  session.Flags
		yaml   string
		url    string
		triage string
		reason string
		calls  int
	}{
		{"Should keep a clean source", nil, session.Flags{}, "", "https://example.com/a", TriageKept, "", 1},
		{"Should quarantine a root URL without any call", nil, session.Flags{}, "", "https://example.com/", TriageQuarantined, quality.NotAnArticle, 0},
		{"Should judge a flagged source in shadow and review it", nil, session.Flags{Decisions: session.ModeShadow}, "", "https://example.com/", TriageReview, quality.NotAnArticle, 1},
		{"Should quarantine off-topic when gates apply", offTopic, session.Flags{}, "decisions:\n  gates: apply\n", "https://example.com/a", TriageQuarantined, ReasonOffTopic, 1},
		{"Should only review off-topic by default", offTopic, session.Flags{}, "\n", "https://example.com/a", TriageReview, ReasonOffTopic, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := newEvalEnv(t, tt.decide, tt.flags, tt.yaml)
			checker, err := NewChecker(env.s)
			if err != nil {
				t.Fatal(err)
			}
			out, doc, err := checker.Evaluate(context.Background(), env.fetched("Decision engines", tt.url, longBody("decision", 400)))
			if err != nil {
				t.Fatal(err)
			}
			if out.Triage != tt.triage || out.Reason != tt.reason || doc == nil {
				t.Fatalf("outcome = %+v", out)
			}
			if got := len(env.or.Calls()); got != tt.calls {
				t.Fatalf("decision calls = %d, want %d", got, tt.calls)
			}
		})
	}
}

// TestEvaluateExclusionUndecidedAndSiteName: a source matching
// decisions.exclude is kept without any decision call; a failed gate
// judgment is written as review/undecided, never kept; a page whose title is
// the site name the fetcher reported is not an article (no call).
func TestEvaluateExclusionUndecidedAndSiteName(t *testing.T) {
	t.Parallel()
	brokenRole := func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "role" {
			return map[string]any{"type": "choice", "choice": "core"}
		}
		return nil
	}
	tests := []struct {
		name     string
		decide   fakes.DecideFunc
		yaml     string
		title    string
		site     string
		triage   string
		reason   string
		excluded bool
		calls    int
	}{
		{name: "excluded source is kept without a call", yaml: "decisions:\n  exclude:\n    - raw/articles/**\n", title: "Decision engines", triage: TriageKept, excluded: true},
		{name: "excluded glob that does not match still judges", yaml: "decisions:\n  exclude:\n    - raw/private/**\n", title: "Decision engines", triage: TriageKept, calls: 2},
		{name: "invalid role receipt is undecided review", decide: brokenRole, yaml: "\n", title: "Decision engines", triage: TriageReview, reason: ReasonUndecided, calls: 2},
		{name: "title equal to the reported site name", title: "Catapult", site: "Catapult", triage: TriageQuarantined, reason: quality.NotAnArticle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := newEvalEnv(t, tt.decide, session.Flags{}, tt.yaml)
			writeDoc(t, env.root, "raw/articles/other.md", map[string]any{
				"title": "Decision engines survey", "source_url": "https://example.com/survey", "source_kind": "article",
			}, longBody("survey", 400))
			checker, err := NewChecker(env.s)
			if err != nil {
				t.Fatal(err)
			}
			in := env.fetched(tt.title, "https://example.com/a", longBody("decision", 400))
			in.SiteName = tt.site
			out, _, err := checker.Evaluate(context.Background(), in)
			if err != nil {
				t.Fatal(err)
			}
			if out.Triage != tt.triage || out.Reason != tt.reason || out.Excluded != tt.excluded {
				t.Fatalf("outcome = %+v", out)
			}
			if tt.excluded && out.Evidence != ExcludedNote {
				t.Fatalf("excluded evidence = %q", out.Evidence)
			}
			if got := len(env.or.Calls()); got != tt.calls {
				t.Fatalf("decision calls = %d, want %d", got, tt.calls)
			}
		})
	}
}

func TestEvaluateNearDuplicatePolicy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		relation   string
		triage     string
		supersedes []string
	}{
		{"Should quarantine the same content", RelationSameContent, TriageQuarantined, nil},
		{"Should keep a new version and supersede the old one", RelationNewVersion, TriageKept, []string{"[[handbook]]"}},
		{"Should ignore related documents", "related", TriageKept, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := newEvalEnv(t, func(_ fakes.Call, q fakes.Question) any {
				if strings.HasPrefix(q.ID, "relation_") {
					return fakes.Pick(tt.relation, 0.9)
				}
				return nil
			}, session.Flags{}, "")
			writeDoc(t, env.root, "raw/articles/handbook.md", map[string]any{
				"title": "Typed judges handbook", "source_url": "https://example.com/v1", "source_kind": "article",
			}, longBody("handbook", 400))
			checker, err := NewChecker(env.s)
			if err != nil {
				t.Fatal(err)
			}
			out, _, err := checker.Evaluate(context.Background(), env.fetched("Typed judges handbook 2", "https://example.com/v2", longBody("judges", 400)))
			if err != nil {
				t.Fatal(err)
			}
			if out.Triage != tt.triage || !slices.Equal(out.Supersedes, tt.supersedes) {
				t.Fatalf("outcome = %+v", out)
			}
			if tt.triage == TriageQuarantined && (out.Reason != ReasonDuplicate || out.DuplicateOf != "raw/articles/handbook.md") {
				t.Fatalf("duplicate outcome = %+v", out)
			}
		})
	}
}

func TestPrefetchSkipsRecordsAndShadows(t *testing.T) {
	t.Parallel()
	decide := func(_ fakes.Call, q fakes.Question) any {
		switch q.ID {
		case "role_item_i1":
			return fakes.Pick("off_topic", 0.9)
		case "role_item_i2":
			return fakes.Pick("core", 0.9)
		case "role_item_i3":
			return fakes.Dist(map[string]float64{"general": 0.6, "off_topic": 0.3, "core": 0.1})
		}
		return nil
	}
	items := []Item{
		{URL: "https://food.example.com/pasta", Title: "Pasta"},
		{URL: "https://example.com/judges", Title: "Judges"},
		{URL: "https://example.com/misc", Title: "Misc"},
		{URL: "https://example.com/audit/x", Title: "Audit", Path: "raw/audit/x.md"},
	}
	for _, mode := range []string{session.ModeApply, session.ModeShadow} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			env := newEvalEnv(t, decide, session.Flags{}, "decisions:\n  gates: "+mode+"\n")
			env.s.Contract.CollectedOnPurposePaths = []string{"raw/audit/**"}
			decided, err := Prefetch(context.Background(), env.s, items, PrefetchOptions{Batch: "b1", Query: "list.txt"})
			if err != nil {
				t.Fatal(err)
			}
			actions := []string{decided[0].Action, decided[1].Action, decided[2].Action, decided[3].Action}
			wantFirst := ActionSkip
			if mode == session.ModeShadow {
				wantFirst = ActionReview
			}
			if !slices.Equal(actions, []string{wantFirst, ActionFetch, ActionReview, ActionFetch}) || decided[3].Reason != PrefetchCollectedPath {
				t.Fatalf("decisions = %+v", decided)
			}
			if len(env.or.Calls()) != 1 || strings.Contains(env.or.Calls()[0].StateString(), "Audit") {
				t.Fatalf("one request for the three judged items, got %d", len(env.or.Calls()))
			}
			rows, err := LoadSkipped(env.root)
			if err != nil {
				t.Fatal(err)
			}
			skips, _ := review.Open(env.root, nil).Pending(review.QueueSkip)
			if mode == session.ModeApply {
				if len(rows) != 1 || rows[0].ID != SkippedID(items[0].URL) || rows[0].Batch != "b1" || len(skips) != 1 {
					t.Fatalf("rows = %+v skips = %+v", rows, skips)
				}
				return
			}
			if len(rows) != 0 || len(skips) != 0 || decided[0].Shadow != ActionSkip {
				t.Fatalf("shadow must not skip: rows %+v skips %+v decision %+v", rows, skips, decided[0])
			}
		})
	}
}

// TestPrefetchExcludedItemsAndExtraQuestions: an item whose would-be path
// matches decisions.exclude is fetched without entering any request; the
// topic's extra relevance questions ride along once per item and never
// change the built-in verdict (spec §4.3, §14).
func TestPrefetchExcludedItemsAndExtraQuestions(t *testing.T) {
	t.Parallel()
	const extra = `{"id":"owner_relevance","version":"1","purpose":"relevance","guard":"G","questions":[` +
		`{"id":"vendor_item_{id}","type":"noul","instructions":"Evaluate only item {id}. Is it a vendor page?","source":"topic owner"}]}`
	env := newEvalEnvWith(t, func(_ fakes.Call, q fakes.Question) any {
		switch q.ID {
		case "role_item_i1":
			return fakes.Pick("core", 0.9)
		case "vendor_item_i1":
			return fakes.Noul(0.99)
		}
		return nil
	}, session.Flags{}, "decisions:\n  gates: apply\n  exclude:\n    - raw/private/**\n", func(root string) {
		dir := filepath.Join(root, ".decisions", "banks")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "owner.json"), []byte(extra), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	items := []Item{
		{URL: "https://example.com/secret", Title: "Secret memo", Path: "raw/private/secret-memo.md"},
		{URL: "https://example.com/judges", Title: "Judges", Path: "raw/articles/judges.md"},
	}
	decided, err := Prefetch(context.Background(), env.s, items, PrefetchOptions{Batch: "b1"})
	if err != nil {
		t.Fatal(err)
	}
	if decided[0].Action != ActionFetch || decided[0].Reason != PrefetchExcluded {
		t.Fatalf("excluded item = %+v", decided[0])
	}
	if decided[1].Action != ActionFetch || decided[1].Role != "core" {
		t.Fatalf("judged item = %+v", decided[1])
	}
	calls := env.or.Calls()
	if len(calls) != 1 || strings.Contains(calls[0].StateString(), "Secret memo") {
		t.Fatalf("the excluded item must stay out of the request: %d calls", len(calls))
	}
	if _, ok := calls[0].Questions["vendor_item_i1"]; !ok || len(calls[0].Questions) != 2 {
		t.Fatalf("questions = %v", slices.Sorted(maps.Keys(calls[0].Questions)))
	}
}

func TestSkippedRowsRoundTrip(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	row := SkippedRow{ID: SkippedID("https://example.com/a"), URL: "https://example.com/a", Title: "A", POffTopic: 0.9}
	if err := AppendSkipped(root, row); err != nil {
		t.Fatal(err)
	}
	if err := MarkRescued(root, row, "2026-09-24T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	got, err := FindSkipped(root, row.ID)
	if err != nil || !got.Rescued || got.Title != "A" {
		t.Fatalf("FindSkipped = %+v, %v", got, err)
	}
	if _, err := FindSkipped(root, "sk-missing"); err == nil {
		t.Fatal("unknown ids must fail")
	}
}

func writeDoc(t *testing.T, root, rel string, values map[string]any, body string) {
	t.Helper()
	content, err := frontmatter.Generate(values, body)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
