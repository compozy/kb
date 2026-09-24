package classify

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/frontmatter"
)

func TestValidateSummary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, summary, title string
		ok                   bool
		want                 string
	}{
		{name: "valid", summary: "Explains how BM25 ranks documents.", title: "BM25", ok: true},
		{name: "trimmed", summary: "  Explains BM25.  ", title: "BM25", ok: true},
		{name: "empty", summary: " ", title: "BM25"},
		{name: "too long without a sentence end", summary: strings.Repeat("a", MaxSummaryChars+1), title: "BM25"},
		{name: "over-long trimmed to last sentence", summary: "First sentence about BM25. " + strings.Repeat("b", MaxSummaryChars), title: "BM25", ok: true, want: "First sentence about BM25."},
		{name: "multi-line", summary: "One.\nTwo.", title: "BM25"},
		{name: "copy of the title", summary: "How BM25 works!", title: "How BM25 Works"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := validateSummary(tt.summary, tt.title)
			if (err == nil) != tt.ok {
				t.Fatalf("validateSummary(%q) err = %v, want ok %v", tt.summary, err, tt.ok)
			}
			if tt.want != "" && got != tt.want {
				t.Fatalf("validateSummary(%q) = %q, want %q", tt.summary, got, tt.want)
			}
		})
	}
}

func TestFilterEntities(t *testing.T) {
	t.Parallel()
	body := "We compare OpenAI, Anthropic and the Model Context Protocol."
	names := []string{"OpenAI", "openai", " Anthropic ", "Google", "model context protocol", ""}
	for range 20 {
		names = append(names, "OpenAI")
	}
	kept, dropped := filterEntities(names, body)
	if !slices.Equal(kept, []string{"OpenAI", "Anthropic", "model context protocol"}) || dropped != 1 {
		t.Fatalf("filterEntities = %v, dropped %d", kept, dropped)
	}

	many := make([]string, 0, 20)
	var text strings.Builder
	for i := range 20 {
		name := "Entity" + string(rune('A'+i))
		many = append(many, name)
		text.WriteString(name + " ")
	}
	kept, dropped = filterEntities(many, text.String())
	if len(kept) != MaxEntities || dropped != 5 {
		t.Fatalf("cap: kept %d, dropped %d", len(kept), dropped)
	}
}

func TestValidateQuestions(t *testing.T) {
	t.Parallel()
	if _, err := validateQuestions([]string{"A?", "B?"}); err == nil {
		t.Fatal("two questions must fail")
	}
	if _, err := validateQuestions([]string{"A?", "a?", "B?"}); err == nil {
		t.Fatal("duplicates do not count")
	}
	got, err := validateQuestions([]string{"A?", "B?", "C?", "D?", "E?", "F?", "G?", "multi\nline?"})
	if err != nil || len(got) != MaxQuestions {
		t.Fatalf("validateQuestions = %v, %v", got, err)
	}
}

func TestValidateCriterion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, criterion string
		ok              bool
	}{
		{name: "valid", criterion: "Documents that compare agent protocols such as MCP and A2A.", ok: true},
		{name: "title once", criterion: "Documents about Agent Protocols and their transports.", ok: true},
		{name: "title twice", criterion: "Agent Protocols documents that define agent protocols.", ok: false},
		{name: "empty", criterion: ""},
		{name: "too long", criterion: strings.Repeat("x", MaxCriterionChars+1)},
		{name: "multi-line", criterion: "One.\nTwo."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := validateCriterion(tt.criterion, "Agent Protocols"); (err == nil) != tt.ok {
				t.Fatalf("validateCriterion err = %v, want ok %v", err, tt.ok)
			}
		})
	}
}

func TestFilterAliases(t *testing.T) {
	t.Parallel()
	taken := map[string]bool{"llm evaluation": true, "evals": true}
	proposed := []string{
		"A2A", "a2a", "Agent-to-Agent", "Agent Protocols", "evals",
		"one two three four five six seven", "", "Protocolos de agentes",
		"P1", "P2", "P3", "P4", "P5", "P6",
	}
	got := filterAliases(proposed, "Agent Protocols", func(key string) bool { return taken[key] })
	want := []string{"A2A", "Agent-to-Agent", "Protocolos de agentes", "P1", "P2", "P3", "P4", "P5"}
	if !slices.Equal(got, want) {
		t.Fatalf("filterAliases = %v, want %v", got, want)
	}
}

func TestValidConceptsAndSanitize(t *testing.T) {
	t.Parallel()
	proposed := []proposedConcept{
		{Title: "Agent Protocols", Criterion: "Documents that compare protocols agents use."},
		{Title: "agent protocols", Criterion: "Duplicate title."},
		{Title: "RAG", Criterion: "Documents about retrieval."},
		{Title: "Existing", Criterion: "Collides with an article."},
		{Title: "Bad", Criterion: ""},
		{Title: "A/B: Testing?", Criterion: "Documents about experiments."},
		{Title: "A B Testing", Criterion: "Same file name as the previous one."},
		{Title: "one two three four five six seven eight nine", Criterion: "Too many words."},
	}
	got := validConcepts(proposed, map[string]bool{"existing": true}, 10)
	titles := make([]string, 0, len(got))
	for _, c := range got {
		titles = append(titles, c.Title)
	}
	if !slices.Equal(titles, []string{"Agent Protocols", "RAG", "A/B: Testing?"}) {
		t.Fatalf("validConcepts = %v", titles)
	}
	if got := validConcepts(proposed, nil, 1); len(got) != 1 {
		t.Fatalf("limit ignored: %d", len(got))
	}

	for title, want := range map[string]string{
		"A/B: Testing?":        "A B Testing",
		"  ..Hidden.. ":        "Hidden",
		"Multi   space\ttitle": "Multi space title",
		"[[Link]] #tag ^block": "Link tag block",
		"///":                  "",
	} {
		if got := sanitizeFileName(title); got != want {
			t.Errorf("sanitizeFileName(%q) = %q, want %q", title, got, want)
		}
	}
}

func TestCreateStubArticle(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	rel, err := CreateStubArticle(root, "demo", "A/B Testing", "Documents about controlled experiments.", testNow)
	if err != nil {
		t.Fatalf("CreateStubArticle: %v", err)
	}
	if rel != "wiki/concepts/A B Testing.md" {
		t.Fatalf("rel = %q", rel)
	}
	values, body, err := frontmatter.Parse(readFile(t, filepath.Join(root, filepath.FromSlash(rel))))
	if err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"title": "A/B Testing", "type": "wiki", "stage": "stub", "domain": "demo", "criterion": "Documents about controlled experiments."} {
		if got := frontmatter.GetString(values, key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
	if tags := frontmatter.GetStringSlice(values, "tags"); !slices.Equal(tags, []string{"demo", "wiki", "stub"}) {
		t.Errorf("tags = %v", tags)
	}
	if sources, ok := values["sources"]; !ok || len(frontmatter.GetStringSlice(values, "sources")) != 0 || sources == nil {
		t.Errorf("sources = %#v", values["sources"])
	}
	if got := frontmatter.GetTime(values, "created").Format("2006-01-02"); got != "2026-09-24" {
		t.Errorf("created = %q", got)
	}
	if body != "# A/B Testing\n" {
		t.Errorf("body = %q", body)
	}
	doc, err := corpus.ReadDocument(root, rel, corpus.KindArticle)
	if err != nil || !isStub(doc) {
		t.Fatalf("stub not recognized: %v", err)
	}

	before := readFile(t, filepath.Join(root, filepath.FromSlash(rel)))
	again, err := CreateStubArticle(root, "demo", "A/B Testing", "Other.", testNow)
	if !errors.Is(err, ErrArticleExists) || again != rel {
		t.Fatalf("second create = %q, %v; want ErrArticleExists", again, err)
	}
	if readFile(t, filepath.Join(root, filepath.FromSlash(rel))) != before {
		t.Fatal("an existing article must never be overwritten")
	}
	if _, err := CreateStubArticle(root, "demo", "///", "", testNow); err == nil {
		t.Fatal("unusable title must fail")
	}
	if _, err := os.Stat(filepath.Join(root, ConceptsDir)); err != nil {
		t.Fatal(err)
	}
}
