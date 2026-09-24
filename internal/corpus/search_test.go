package corpus_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func TestTokenize(t *testing.T) {
	t.Parallel()

	got := corpus.Tokenize("The Model-Context Protocol (MCP) é uma forma de integração, v2 a b")
	want := []string{"model", "context", "protocol", "mcp", "forma", "integração", "v2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize = %#v, want %#v", got, want)
	}
}

func TestBM25Search(t *testing.T) {
	t.Parallel()

	docs := []*corpus.Document{
		{Path: "a.md", Title: "Model Context Protocol", Body: "Protocol for tools."},
		{Path: "b.md", Title: "Cooking pasta", Body: "Boil water. Add pasta. Protocol mention once."},
		{Path: "c.md", Title: "Agents", Frontmatter: map[string]any{"summary": "Agents call tools over the model context protocol."}},
		{Path: "d.md", Title: "Unrelated", Body: "Nothing here."},
	}
	index := corpus.NewBM25(docs, corpus.DefaultFields, corpus.DefaultWeights())

	hits := index.Search("model context protocol", 0)
	gotPaths := make([]string, 0, len(hits))
	for _, hit := range hits {
		gotPaths = append(gotPaths, hit.Doc.Path)
		if hit.Score <= 0 {
			t.Fatalf("non-positive score in %#v", hit)
		}
	}
	if !reflect.DeepEqual(gotPaths, []string{"a.md", "c.md", "b.md"}) {
		t.Fatalf("Search paths = %#v", gotPaths)
	}
	if got := index.Search("model context protocol", 1); len(got) != 1 || got[0].Doc.Path != "a.md" {
		t.Fatalf("Search k=1 = %#v", got)
	}
	if got := index.Search("the of", 5); got != nil {
		t.Fatalf("stopword-only query = %#v", got)
	}
	if got := corpus.NewBM25(nil, corpus.DefaultFields, nil).Search("x", 1); got != nil {
		t.Fatalf("empty index Search = %#v", got)
	}

	tied := corpus.NewBM25([]*corpus.Document{{Path: "z.md", Title: "same"}, {Path: "y.md", Title: "same"}}, func(d *corpus.Document) map[string]string {
		return map[string]string{"title": d.Title}
	}, nil)
	tiedHits := tied.Search("same", 0)
	if len(tiedHits) != 2 || tiedHits[0].Doc.Path != "y.md" {
		t.Fatalf("ties must sort by path: %#v", tiedHits)
	}
}

func TestDefaultFields(t *testing.T) {
	t.Parallel()

	doc := &corpus.Document{
		Title:   "T",
		Aliases: []string{"A1", "A2"},
		Body:    "See [[x|y]]" + strings.Repeat("z", corpus.HeadChars),
		Frontmatter: map[string]any{
			"criterion": "C",
			"summary":   "S",
			"entities":  []string{"E1", "E2"},
			"questions": []string{"Q1"},
		},
	}
	fields := corpus.DefaultFields(doc)
	if fields["title"] != "T" || fields["aliases"] != "A1\nA2" || fields["criterion"] != "C" || fields["summary"] != "S" ||
		fields["entities"] != "E1\nE2" || fields["questions"] != "Q1" {
		t.Fatalf("DefaultFields = %#v", fields)
	}
	if !strings.HasPrefix(fields["head"], "See y") || len([]rune(fields["head"])) != corpus.HeadChars {
		t.Fatalf("head field = %q", fields["head"][:20])
	}
	if len(corpus.DefaultWeights()) != len(fields) {
		t.Fatal("DefaultWeights must cover every default field")
	}
}

func TestDictionaryMentions(t *testing.T) {
	t.Parallel()

	agents := &corpus.Document{Path: "wiki/concepts/Agents.md", Title: "AI Agents", Aliases: []string{"agentic systems", "AI"}}
	cot := &corpus.Document{Path: "wiki/concepts/CoT.md", Title: "Chain-of-Thought", Aliases: []string{"CoT", "ab"}}
	ml := &corpus.Document{Path: "wiki/concepts/ML.md", Title: "Machine Learning"}
	dictionary := corpus.NewDictionary([]*corpus.Document{agents, cot, ml})

	text := strings.Join([]string{
		"---",
		"title: Machine Learning notes",
		"---",
		"# Machine Learning heading",
		"We study ai agents and AGENTIC   systems. Then AI helps.",
		"Chain-of-thought prompting works; chain of thought does not match.",
		"Code `machine learning` and [[Machine Learning]] and [ML](machine learning.md) skip.",
		"URL https://example.com/machine-learning skip. Ai em português não conta.",
		"```",
		"machine learning in code",
		"```",
		"Final Machine",
		"Learning wraps. Machinelearning is one word. ab is short.",
	}, "\n")

	mentions := dictionary.Mentions(text)
	type got struct{ article, term, match string }
	gotMentions := make([]got, 0, len(mentions))
	for _, mention := range mentions {
		gotMentions = append(gotMentions, got{mention.Article.Path, mention.Term, text[mention.Start:mention.End]})
	}
	want := []got{
		{"wiki/concepts/Agents.md", "AI Agents", "ai agents"},
		{"wiki/concepts/Agents.md", "agentic systems", "AGENTIC   systems"},
		{"wiki/concepts/Agents.md", "AI", "AI"},
		{"wiki/concepts/CoT.md", "Chain-of-Thought", "Chain-of-thought"},
		{"wiki/concepts/ML.md", "Machine Learning", "Machine\nLearning"},
	}
	if !reflect.DeepEqual(gotMentions, want) {
		t.Fatalf("Mentions = %#v\nwant %#v", gotMentions, want)
	}

	if mentions[0].Sentence != "We study ai agents and AGENTIC   systems." {
		t.Fatalf("sentence = %q", mentions[0].Sentence)
	}
	if mentions[2].Sentence != "Then AI helps." {
		t.Fatalf("sentence = %q", mentions[2].Sentence)
	}

	long := "Start " + strings.Repeat("padding ", 100) + "Machine Learning " + strings.Repeat("tail ", 100)
	longMentions := dictionary.Mentions(long)
	if len(longMentions) != 1 || len([]rune(longMentions[0].Sentence)) > 300 || !strings.Contains(longMentions[0].Sentence, "Machine Learning") {
		t.Fatalf("long sentence = %#v", longMentions)
	}

	shared := corpus.NewDictionary([]*corpus.Document{
		{Path: "a.md", Title: "Shared Term"},
		{Path: "b.md", Title: "Other", Aliases: []string{"shared term"}},
	})
	if got := shared.Mentions("a shared term here"); len(got) != 2 {
		t.Fatalf("ambiguous term must report every article: %#v", got)
	}

	var empty *corpus.Dictionary
	if got := empty.Mentions("x"); got != nil {
		t.Fatalf("nil dictionary Mentions = %#v", got)
	}
}
