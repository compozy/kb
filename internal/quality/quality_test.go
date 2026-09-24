package quality

import (
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func words(n int) string {
	return strings.TrimSpace(strings.Repeat("word ", n))
}

func codes(flags []Flag) []string {
	out := []string{}
	for _, flag := range flags {
		out = append(out, flag.Code)
	}
	return out
}

func TestCheck(t *testing.T) {
	t.Parallel()
	long := words(400)
	tests := []struct {
		name string
		in   Input
		want []string
	}{
		{name: "complete article", in: Input{Title: "How BM25 works", Body: long, SourceURL: "https://blog.example.com/posts/bm25"}, want: []string{}},
		{name: "site root", in: Input{Title: "Example blog", Body: long, SourceURL: "https://blog.example.com/"}, want: []string{NotAnArticle}},
		{name: "empty path", in: Input{Title: "Example", Body: long, SourceURL: "https://example.com"}, want: []string{NotAnArticle}},
		{name: "index.html", in: Input{Title: "Welcome", Body: long, SourceURL: "https://example.com/index.html"}, want: []string{NotAnArticle}},
		{name: "redirect dropped path", in: Input{Title: "Welcome", Body: long, RequestedURL: "https://example.com/post/1", FinalURL: "https://example.com/"}, want: []string{NotAnArticle}},
		{name: "final URL wins over source URL", in: Input{Title: "Post", Body: long, SourceURL: "https://example.com/", FinalURL: "https://example.com/post/1"}, want: []string{}},
		{name: "title equals reported site name", in: Input{Title: "Catapult", SiteName: "Catapult", Body: long, SourceURL: "https://catapult.com/products"}, want: []string{NotAnArticle}},
		{name: "title equals host label", in: Input{Title: "  catapult ", Body: long, SourceURL: "https://www.catapult.com/about"}, want: []string{NotAnArticle}},
		{name: "title longer than site name", in: Input{Title: "Catapult Vector review", SiteName: "Catapult", Body: long, SourceURL: "https://catapult.com/vector"}, want: []string{}},
		{name: "non-http source is ignored", in: Input{Title: "Notes", Body: long, SourceURL: "file:///tmp/notes.md"}, want: []string{}},
		{name: "http 404", in: Input{Title: "Post", Body: long, SourceURL: "https://example.com/p", StatusCode: 404}, want: []string{ErrorPage}},
		{name: "http 503", in: Input{Title: "Post", Body: long, StatusCode: 503}, want: []string{ErrorPage}},
		{name: "not found title", in: Input{Title: "Page Not Found | Example", Body: long}, want: []string{ErrorPage}},
		{name: "404 title", in: Input{Title: "404", Body: long}, want: []string{ErrorPage}},
		{name: "access denied title", in: Input{Title: "ACCESS DENIED", Body: long}, want: []string{ErrorPage}},
		{name: "404 inside a word does not match", in: Input{Title: "Model X4040 review", Body: long}, want: []string{}},
		{name: "thin body", in: Input{Title: "Short", Body: words(120)}, want: []string{Thin}},
		{name: "thin skipped", in: Input{Title: "Short", Body: words(20), SkipThin: true}, want: []string{}},
		{name: "every flag in fixed order", in: Input{Title: "Not Found", Body: "", SourceURL: "https://example.com/", StatusCode: 404}, want: []string{NotAnArticle, ErrorPage, Thin}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := codes(Check(tt.in))
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentWordsDropsNavigation(t *testing.T) {
	t.Parallel()
	chrome := "Home | Pricing | Blog | Careers"
	body := strings.Join([]string{
		"[Home](https://x.com) · [Blog](https://x.com/blog)",
		"- [[Wiki page]]",
		"https://example.com/a",
		"| --- | :---: |",
		chrome,
		"Five real words right here.",
		"A sentence with a [link](https://x.com) inside counts.",
	}, "\n")
	hostLines := map[string]int{NormalizeLine(chrome): 3}
	if got, want := ContentWords(body, hostLines), 5+7; got != want {
		t.Fatalf("ContentWords() = %d, want %d", got, want)
	}
	if got := ContentWords(body, map[string]int{NormalizeLine(chrome): 2}); got != 5+7+4 {
		t.Fatalf("a line on 2 sources is not chrome: ContentWords() = %d", got)
	}
}

func TestThinCountsAfterNavigationRemoval(t *testing.T) {
	t.Parallel()
	nav := strings.Repeat("[Link](https://x.com/a) [Other](https://x.com/b)\n", 400)
	flags := Check(Input{Title: "Listing", Body: nav + words(50)})
	if got := codes(flags); !reflect.DeepEqual(got, []string{Thin}) {
		t.Fatalf("Check() = %v, want thin", got)
	}
}

func TestHostLineCounts(t *testing.T) {
	t.Parallel()
	doc := func(path, url, body string) *corpus.Document {
		return &corpus.Document{Path: path, Body: body, Frontmatter: map[string]any{"source_url": url}}
	}
	docs := []*corpus.Document{
		doc("raw/a.md", "https://www.x.com/a", "Menu\nMenu\nAlpha text"),
		doc("raw/b.md", "https://x.com/b", "  menu  \nBeta text"),
		doc("raw/c.md", "https://y.com/c", "Menu\nGamma"),
		doc("raw/d.md", "", "Menu"),
	}
	counts := HostLineCounts(docs)
	if got := counts["x.com"]["menu"]; got != 2 {
		t.Fatalf("x.com menu = %d, want 2 (one per source, www stripped)", got)
	}
	if got := counts["y.com"]["menu"]; got != 1 {
		t.Fatalf("y.com menu = %d, want 1", got)
	}
	if len(counts) != 2 {
		t.Fatalf("hosts = %v, want only x.com and y.com", counts)
	}
}
