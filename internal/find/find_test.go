package find

import (
	"context"
	"math"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/resolve"
)

func testDoc(path string, fm map[string]any, body string) *corpus.Document {
	if fm == nil {
		fm = map[string]any{}
	}
	kind := corpus.KindSource
	if strings.HasPrefix(path, "wiki/") {
		kind = corpus.KindArticle
	}
	title, _ := fm["title"].(string)
	if title == "" {
		title = strings.TrimSuffix(path[strings.LastIndex(path, "/")+1:], ".md")
	}
	return &corpus.Document{Path: path, Title: title, Kind: kind, Frontmatter: fm, Body: body, BodyHash: corpus.BodyHash(body), Aliases: resolve.Aliases(fm)}
}

func TestScore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		p     float64
		spec  float64
		match bool
		depth float64
		want  float64
	}{
		{name: "bare probability", p: 0.8, want: 0.8},
		{name: "full specificity", p: 0.8, spec: 3, want: 0.8 * 1.25},
		{name: "kind match", p: 0.8, match: true, want: 0.8 * 1.2},
		{name: "full depth", p: 0.8, depth: 3, want: 0.8 * 1.1},
		{name: "everything", p: 1, spec: 1.5, match: true, depth: 1.5, want: 1 * 1.125 * 1.2 * 1.05},
		{name: "out-of-range levels are clamped", p: 1, spec: 7, depth: -2, want: 1.25},
	}
	for _, tt := range tests {
		if got := Score(tt.p, tt.spec, tt.match, tt.depth); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("%s: Score = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestKindMatches(t *testing.T) {
	t.Parallel()
	tests := []struct {
		kind  string
		genre string
		want  bool
	}{
		{KindDefinition, "reference_docs", true},
		{KindDefinition, "article_or_essay", true},
		{KindHowTo, "tutorial_or_guide", true},
		{KindHowTo, "paper", false},
		{KindComparison, "paper", true},
		{KindEvidence, "dataset_or_benchmark", true},
		{KindOpinion, "opinion_or_discussion", true},
		{KindOther, "paper", false},
		{KindDefinition, "", false},
	}
	for _, tt := range tests {
		if got := KindMatches(tt.kind, tt.genre); got != tt.want {
			t.Errorf("KindMatches(%q, %q) = %v, want %v", tt.kind, tt.genre, got, tt.want)
		}
	}
}

func TestRank(t *testing.T) {
	t.Parallel()
	guide := testDoc("raw/guide.md", map[string]any{"genre": "tutorial_or_guide", "depth": 2}, "")
	paper := testDoc("raw/paper.md", map[string]any{"genre": "paper"}, "")
	weak := testDoc("raw/weak.md", nil, "")
	low := testDoc("raw/low.md", nil, "")
	extra := testDoc("raw/extra.md", nil, "")
	lost := testDoc("raw/lost.md", nil, "")

	items := []Judged{
		{Doc: paper, PAnswers: 0.9, Specificity: 1, Decided: true},
		{Doc: guide, PAnswers: 0.85, Specificity: 3, Decided: true},
		{Doc: weak, PAnswers: 0.52, Decided: true},
		{Doc: low, PAnswers: 0.3, Decided: true},
		{Doc: extra, PAnswers: 0.7, Decided: true},
		{Doc: lost, Reason: "undecided:timeout"},
	}
	hits, drops := Rank(items, KindHowTo, 0.5, 0.35, 2)

	gotHits := make([]string, 0, len(hits))
	for _, hit := range hits {
		gotHits = append(gotHits, hit.Path)
	}
	// guide: 0.85 × 1.25 × 1.2 × (1 + 0.1·2/3) beats paper: 0.9 × (1 + 0.25/3).
	if want := []string{"raw/guide.md", "raw/paper.md"}; !reflect.DeepEqual(gotHits, want) {
		t.Fatalf("hits = %v, want %v", gotHits, want)
	}
	if hits[0].Depth == nil || *hits[0].Depth != 2 || hits[0].Kind != "tutorial_or_guide" {
		t.Fatalf("hit = %+v", hits[0])
	}

	reasons := map[string]string{}
	for _, drop := range drops {
		reasons[drop.Path] = drop.Reason
	}
	want := map[string]string{
		"raw/weak.md":  ReasonBelowWindow,
		"raw/low.md":   ReasonLowProbability,
		"raw/extra.md": ReasonLimit,
		"raw/lost.md":  "undecided:timeout",
	}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("drop reasons = %v, want %v", reasons, want)
	}
}

type stubVectors []string

func (s stubVectors) Search(context.Context, string, int) ([]string, error) { return s, nil }

func TestCandidatesExclusions(t *testing.T) {
	t.Parallel()
	docs := []*corpus.Document{
		testDoc("raw/kept.md", map[string]any{"title": "Vector search tuning", "genre": "tutorial_or_guide", "relevance": "core", "depth": 2, "concepts": []any{"[[Vector Database]]"}}, "How to tune vector search."),
		testDoc("raw/quarantined.md", map[string]any{"title": "Vector search spam", "triage": "quarantined", "genre": "tutorial_or_guide"}, "vector search vector search"),
		testDoc("raw/paper.md", map[string]any{"title": "Vector search paper", "genre": "paper", "relevance": "core", "depth": 3}, "vector search results"),
		testDoc("raw/shallow.md", map[string]any{"title": "Vector search intro", "genre": "tutorial_or_guide", "relevance": "core", "depth": 0.5}, "vector search basics"),
		testDoc("raw/adjacent.md", map[string]any{"title": "Vector search aside", "genre": "tutorial_or_guide", "relevance": "adjacent", "depth": 2}, "vector search"),
		testDoc("wiki/concepts/Vector Database.md", map[string]any{"title": "Vector Database", "genre": "tutorial_or_guide", "relevance": "core", "depth": 2}, "Stores embeddings."),
		testDoc("raw/unrelated.md", map[string]any{"title": "Cooking"}, "recipes"),
	}

	tests := []struct {
		name    string
		query   Query
		want    []string
		dropped map[string]string
	}{
		{
			name:  "no filters: quarantined dropped, concepts of top hits added",
			query: Query{Text: "vector search"},
			want:  []string{"raw/kept.md", "raw/paper.md", "raw/shallow.md", "raw/adjacent.md", "wiki/concepts/Vector Database.md"},
			dropped: map[string]string{
				"raw/quarantined.md": ReasonQuarantined,
			},
		},
		{
			name:  "kind, relevance and depth filters",
			query: Query{Text: "vector search", Kind: "tutorial_or_guide", Relevance: "core", MinDepth: 1},
			want:  []string{"raw/kept.md", "wiki/concepts/Vector Database.md"},
			dropped: map[string]string{
				"raw/quarantined.md": ReasonQuarantined,
				"raw/paper.md":       ReasonFilteredFacet,
				"raw/shallow.md":     ReasonFilteredFacet,
				"raw/adjacent.md":    ReasonFilteredFacet,
			},
		},
		{
			name:  "concept filter",
			query: Query{Text: "vector search", Concept: "[[Vector Database]]"},
			want:  []string{"raw/kept.md"},
		},
		{
			name:  "vector candidates are merged",
			query: Query{Text: "vector search", Kind: "tutorial_or_guide", Vectors: stubVectors{"raw/unrelated.md", "wiki/concepts/Vector Database"}},
			want:  []string{"raw/kept.md", "raw/shallow.md", "raw/adjacent.md", "wiki/concepts/Vector Database.md"},
			dropped: map[string]string{
				"raw/unrelated.md": ReasonFilteredFacet,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			candidates, drops := Candidates(context.Background(), docs, tt.query)
			got := make([]string, 0, len(candidates))
			for _, doc := range candidates {
				got = append(got, doc.Path)
			}
			slices.Sort(got)
			want := slices.Clone(tt.want)
			slices.Sort(want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("candidates = %v, want %v", got, want)
			}
			reasons := map[string]string{}
			for _, drop := range drops {
				reasons[drop.Path] = drop.Reason
			}
			for path, reason := range tt.dropped {
				if reasons[path] != reason {
					t.Fatalf("drop %s = %q, want %q (all: %v)", path, reasons[path], reason, reasons)
				}
			}
		})
	}
}

func TestFacets(t *testing.T) {
	t.Parallel()
	docs := []*corpus.Document{
		testDoc("raw/a.md", map[string]any{"genre": "paper", "relevance": "core", "depth": 2.4, "concepts": []any{"[[RAG]]", "[[Agents]]"}}, ""),
		testDoc("raw/b.md", map[string]any{"genre": "paper", "relevance": "adjacent", "depth": 3, "concepts": []any{"[[RAG]]"}}, ""),
		testDoc("raw/c.md", map[string]any{"genre": "tutorial_or_guide", "depth": 0.2}, ""),
		testDoc("raw/q.md", map[string]any{"genre": "paper", "triage": "quarantined"}, ""),
		testDoc("wiki/concepts/RAG.md", map[string]any{"genre": "reference_docs"}, ""),
	}
	report := Facets(docs)
	if report.Documents != 4 || report.Sources != 3 || report.Articles != 1 || report.Quarantined != 1 {
		t.Fatalf("totals = %+v", report)
	}
	wantGenre := []Count{{"paper", 2}, {"reference_docs", 1}, {"tutorial_or_guide", 1}}
	if !reflect.DeepEqual(report.Genre, wantGenre) {
		t.Fatalf("genre = %v", report.Genre)
	}
	wantRelevance := []Count{{FacetNone, 2}, {"adjacent", 1}, {"core", 1}}
	if !reflect.DeepEqual(report.Relevance, wantRelevance) {
		t.Fatalf("relevance = %v", report.Relevance)
	}
	wantDepth := []Count{{"0", 1}, {"1", 0}, {"2", 1}, {"3", 1}, {FacetNone, 1}}
	if !reflect.DeepEqual(report.Depth, wantDepth) {
		t.Fatalf("depth = %v", report.Depth)
	}
	wantConcepts := []Count{{"RAG", 2}, {"Agents", 1}}
	if !reflect.DeepEqual(report.Concepts, wantConcepts) {
		t.Fatalf("concepts = %v", report.Concepts)
	}
}
