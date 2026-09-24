package link

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/resolve"
)

// testDoc builds an in-memory corpus document.
func testDoc(path string, kind corpus.Kind, fm map[string]any, body string) *corpus.Document {
	if fm == nil {
		fm = map[string]any{}
	}
	title, _ := fm["title"].(string)
	if title == "" {
		title = strings.TrimSuffix(path[strings.LastIndex(path, "/")+1:], ".md")
	}
	return &corpus.Document{
		Path: path, Title: title, Kind: kind, Frontmatter: fm, Body: body,
		BodyHash: corpus.BodyHash(body), Aliases: resolve.Aliases(fm),
	}
}

func article(stem string, fm map[string]any, body string) *corpus.Document {
	if fm == nil {
		fm = map[string]any{}
	}
	if _, ok := fm["title"]; !ok {
		fm["title"] = stem
	}
	return testDoc("wiki/concepts/"+stem+".md", corpus.KindArticle, fm, body)
}

type stubNeighbours []string

func (s stubNeighbours) Neighbours(context.Context, *corpus.Document, int) ([]string, error) {
	return s, nil
}

func candidatePaths(candidates []Candidate) []string {
	paths := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		paths = append(paths, candidate.Path)
	}
	return paths
}

func TestCandidates(t *testing.T) {
	t.Parallel()

	rag := article("Retrieval Augmented Generation", map[string]any{"aliases": []any{"RAG"}, "criterion": "Grounding generation in retrieved documents."}, "Body about retrieval.")
	vector := article("Vector Database", map[string]any{"criterion": "Databases indexing embeddings for similarity search."}, "Stores embeddings.")
	agents := article("Agents", map[string]any{"concepts": []any{"[[Planning]]", "[[Tool Use]]"}}, "Agents plan and call tools.")
	sourced := article("Evaluation", map[string]any{"sources": []any{"[[raw/articles/eval-paper]]"}}, "How to evaluate.")
	hub := testDoc("wiki/concepts/Concept Index.md", corpus.KindArticle, map[string]any{"title": "Concept Index"}, "RAG and Vector Database.")
	unrelated := article("Cooking", nil, "Recipes and kitchens.")

	articles := []*corpus.Document{rag, vector, agents, sourced, hub, unrelated}

	tests := []struct {
		name       string
		doc        *corpus.Document
		neighbours Neighbours
		want       []string
		wantAbsent []string
		origins    map[string]string
	}{
		{
			name:       "mention by alias comes first and hubs are dropped",
			doc:        testDoc("raw/articles/a.md", corpus.KindSource, nil, "We build RAG pipelines.\n\nThe Concept Index is not a target."),
			want:       []string{rag.Path},
			wantAbsent: []string{hub.Path, unrelated.Path},
			origins:    map[string]string{rag.Path: OriginMention},
		},
		{
			name:    "shared concepts (two or more)",
			doc:     testDoc("raw/articles/c.md", corpus.KindSource, map[string]any{"concepts": []any{"[[Planning]]", "[[Tool Use]]"}}, "Nothing lexical here."),
			want:    []string{agents.Path},
			origins: map[string]string{agents.Path: OriginStructure},
		},
		{
			name:    "article listing the source in sources",
			doc:     testDoc("raw/articles/eval-paper.md", corpus.KindSource, nil, "zzz"),
			want:    []string{sourced.Path},
			origins: map[string]string{sourced.Path: OriginStructure},
		},
		{
			name:    "existing body link is structure",
			doc:     testDoc("raw/articles/d.md", corpus.KindSource, nil, "See [[Cooking|recipes]]."),
			want:    []string{unrelated.Path},
			origins: map[string]string{unrelated.Path: OriginStructure},
		},
		{
			name:    "bm25 over criterion and summary",
			doc:     testDoc("raw/articles/e.md", corpus.KindSource, map[string]any{"summary": "similarity search over embeddings"}, "text"),
			want:    []string{vector.Path},
			origins: map[string]string{vector.Path: OriginLexical},
		},
		{
			name:       "article document never proposes itself",
			doc:        rag,
			wantAbsent: []string{rag.Path},
		},
		{
			name:       "semantic neighbours are optional extras",
			doc:        testDoc("raw/articles/f.md", corpus.KindSource, nil, "zzz"),
			neighbours: stubNeighbours{"wiki/concepts/Cooking", "wiki/concepts/missing"},
			want:       []string{unrelated.Path},
			origins:    map[string]string{unrelated.Path: OriginVector},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			generator := NewGenerator(articles, tt.neighbours)
			got := generator.Candidates(context.Background(), tt.doc)
			paths := candidatePaths(got)
			for index, want := range tt.want {
				if index >= len(paths) || paths[index] != want {
					t.Fatalf("candidates = %v, want prefix %v", paths, tt.want)
				}
			}
			for _, absent := range tt.wantAbsent {
				if slices.Contains(paths, absent) {
					t.Fatalf("candidates = %v, must not contain %s", paths, absent)
				}
			}
			for index, candidate := range got {
				if candidate.ID != fmt.Sprintf("c%d", index+1) {
					t.Fatalf("candidate %d id = %q", index, candidate.ID)
				}
				if want, ok := tt.origins[candidate.Path]; ok && candidate.Origins[0] != want {
					t.Fatalf("%s origins = %v, want first %s", candidate.Path, candidate.Origins, want)
				}
			}
		})
	}
}

func TestCandidatesCapAndHash(t *testing.T) {
	t.Parallel()
	articles := make([]*corpus.Document, 0, 30)
	var body strings.Builder
	for index := range 30 {
		name := fmt.Sprintf("Topic%02d", index)
		articles = append(articles, article(name, nil, "text"))
		body.WriteString(name + " ")
	}
	generator := NewGenerator(articles, nil)
	doc := testDoc("raw/articles/x.md", corpus.KindSource, nil, body.String())
	got := generator.Candidates(context.Background(), doc)
	if len(got) != MaxCandidates {
		t.Fatalf("len(candidates) = %d, want %d", len(got), MaxCandidates)
	}
	reversed := slices.Clone(got)
	slices.Reverse(reversed)
	if CandidateHash(got) != CandidateHash(reversed) {
		t.Fatal("CandidateHash depends on order")
	}
	if CandidateHash(got) == CandidateHash(got[1:]) {
		t.Fatal("CandidateHash ignores membership")
	}
}

func TestIsHub(t *testing.T) {
	t.Parallel()
	for path, want := range map[string]bool{
		"log.md":                         true,
		"CLAUDE.md":                      true,
		"wiki/concepts/Concept Index.md": true,
		"wiki/concepts/index.md":         true,
		"wiki/concepts/Indexing.md":      false,
		"wiki/concepts/Agents.md":        false,
		"raw/articles/a.md":              false,
	} {
		if got := IsHub(path); got != want {
			t.Errorf("IsHub(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestTargetsForm(t *testing.T) {
	t.Parallel()
	files := []resolve.File{
		{Path: "demo/wiki/concepts/RAG.md", TopicRel: "wiki/concepts/RAG.md", InTopic: true},
		{Path: "demo/wiki/concepts/Agents.md", TopicRel: "wiki/concepts/Agents.md", InTopic: true},
		{Path: "other/wiki/concepts/Agents.md"},
	}
	targets := NewTargets(resolve.NewIndex("demo", files), files, "demo")

	tests := []struct {
		path string
		want string
	}{
		{path: "wiki/concepts/RAG.md", want: "RAG"},
		{path: "wiki/concepts/Agents.md", want: "demo/wiki/concepts/Agents"},
	}
	for _, tt := range tests {
		if got := targets.Form(tt.path); got != tt.want {
			t.Errorf("Form(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
	if got := targets.Wikilink("wiki/concepts/RAG.md"); got != "[[RAG]]" {
		t.Errorf("Wikilink = %q", got)
	}
	if got := targets.Resolve("raw/a.md", "RAG"); got != "wiki/concepts/RAG.md" {
		t.Errorf("Resolve(RAG) = %q", got)
	}
	if got := targets.Resolve("raw/a.md", "other/wiki/concepts/Agents"); got != "" {
		t.Errorf("Resolve(other topic) = %q, want empty", got)
	}
}
