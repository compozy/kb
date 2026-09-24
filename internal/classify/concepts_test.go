package classify

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func testArticle(name, title string, aliases []string, criterion, summary, body string) *corpus.Document {
	values := map[string]any{"title": title}
	if criterion != "" {
		values["criterion"] = criterion
	}
	if summary != "" {
		values["summary"] = summary
	}
	return &corpus.Document{Path: "wiki/concepts/" + name + ".md", Title: title, Kind: corpus.KindArticle, Aliases: aliases, Frontmatter: values, Body: body}
}

func stemLink(doc *corpus.Document) string { return linkStem(doc, nil, "") }

func candidatePaths(cs []*concept) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.article.Path)
	}
	return out
}

func TestVocabularyOptionText(t *testing.T) {
	t.Parallel()
	v := newVocabulary([]*corpus.Document{
		testArticle("rag", "Retrieval-Augmented Generation", []string{"RAG"}, "Documents about retrieving context to ground LLM answers.", "A summary.", "Body"),
		testArticle("evals", "LLM Evaluation", nil, "", "Measuring LLM output quality.", "Body"),
		testArticle("empty", "Empty Concept", nil, "", "", "Opening [[Link|text]] paragraph."),
	}, stemLink)

	want := map[string]string{
		"wiki/concepts/rag.md":   "Retrieval-Augmented Generation — RAG — Documents about retrieving context to ground LLM answers.",
		"wiki/concepts/evals.md": "LLM Evaluation — Measuring LLM output quality.",
		"wiki/concepts/empty.md": "Empty Concept — Opening text paragraph.",
	}
	for p, text := range want {
		if got := v.byPath[p].optionText(); got != text {
			t.Errorf("optionText(%s) = %q, want %q", p, got, text)
		}
	}
	if got := v.criterionMissing(); !slices.Equal(got, []string{"wiki/concepts/empty.md", "wiki/concepts/evals.md"}) {
		t.Fatalf("criterionMissing = %v", got)
	}
	if ids := []string{v.concepts[0].id, v.concepts[1].id, v.concepts[2].id}; !slices.Equal(ids, []string{"c01", "c02", "c03"}) {
		t.Fatalf("ids = %v (sorted by path)", ids)
	}
	if v.byPath["wiki/concepts/rag.md"].link != "[[rag]]" {
		t.Fatalf("link = %q", v.byPath["wiki/concepts/rag.md"].link)
	}
}

func TestVocabularyCandidates(t *testing.T) {
	t.Parallel()
	articles := []*corpus.Document{
		testArticle("rag", "Retrieval-Augmented Generation", []string{"RAG"}, "Retrieving context chunks for grounding.", "", ""),
		testArticle("evals", "LLM Evaluation", nil, "Benchmarks, judges and calibration of LLM outputs.", "", ""),
		testArticle("protocols", "Agent Protocols", []string{"MCP"}, "Protocols agents use to talk to tools.", "", ""),
		testArticle("cooking", "Cooking", nil, "Recipes and kitchens.", "", ""),
	}
	v := newVocabulary(articles, stemLink)
	doc := &corpus.Document{Path: "raw/a.md", Title: "Grounding answers", Kind: corpus.KindSource,
		Body: "We use RAG and MCP here. MCP again, and MCP once more. Calibration of judges matters for benchmarks.", Frontmatter: map[string]any{}}

	got := candidatePaths(v.candidates(doc))
	// Alias hits first (MCP ×3 before RAG ×1), then BM25 matches.
	if len(got) < 3 || got[0] != "wiki/concepts/protocols.md" || got[1] != "wiki/concepts/rag.md" {
		t.Fatalf("candidates = %v, want mentions first by count", got)
	}
	if !slices.Contains(got, "wiki/concepts/evals.md") {
		t.Fatalf("candidates = %v, want the BM25 match on calibration/judges/benchmarks", got)
	}
	if slices.Contains(got, "wiki/concepts/cooking.md") {
		t.Fatalf("candidates = %v, unrelated article must not be a candidate", got)
	}

	// An article never is its own candidate or option.
	self := articles[0]
	if slices.Contains(candidatePaths(v.candidates(self)), self.Path) || slices.Contains(candidatePaths(v.options(self)), self.Path) {
		t.Fatal("an article must not be offered as its own concept")
	}
	if len(v.options(doc)) != 4 || len(v.options(self)) != 3 {
		t.Fatalf("options = %d/%d", len(v.options(doc)), len(v.options(self)))
	}
}

func TestVocabularyCaps(t *testing.T) {
	t.Parallel()
	articles := make([]*corpus.Document, 0, 300)
	for i := range 300 {
		title := fmt.Sprintf("Topic%03d", i)
		criterion := "Generic filler criterion."
		if i == 299 {
			criterion = "Documents about zebrafish regeneration."
		}
		articles = append(articles, testArticle(fmt.Sprintf("t%03d", i), title, nil, criterion, "", ""))
	}
	v := newVocabulary(articles, stemLink)
	var body strings.Builder
	for i := range 30 {
		fmt.Fprintf(&body, "Topic%03d ", i)
	}
	body.WriteString("zebrafish regeneration")
	doc := &corpus.Document{Path: "raw/a.md", Title: "Zebrafish", Kind: corpus.KindSource, Body: body.String(), Frontmatter: map[string]any{}}

	options := v.options(doc)
	if len(options) != MaxConceptOptions {
		t.Fatalf("options = %d, want %d", len(options), MaxConceptOptions)
	}
	if !slices.Contains(candidatePaths(options), "wiki/concepts/t299.md") {
		t.Fatal("the best BM25 match must be among the capped options")
	}
	if got := v.candidates(doc); len(got) != MaxConceptCandidates {
		t.Fatalf("candidates = %d, want %d", len(got), MaxConceptCandidates)
	}
}

func TestLinkStem(t *testing.T) {
	t.Parallel()
	doc := &corpus.Document{Path: "wiki/concepts/Agent Protocols.md"}
	if got := linkStem(doc, func(string) bool { return false }, "demo"); got != "[[Agent Protocols]]" {
		t.Fatalf("linkStem = %q", got)
	}
	if got := linkStem(doc, func(string) bool { return true }, "demo"); got != "[[demo/wiki/concepts/Agent Protocols]]" {
		t.Fatalf("ambiguous linkStem = %q", got)
	}
}
