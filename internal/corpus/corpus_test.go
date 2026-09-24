package corpus_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/corpus"
)

func writeTopicFile(t *testing.T, root, relative, content string) string {
	t.Helper()

	absolute := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", relative, err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}

	return absolute
}

func documentPaths(documents []*corpus.Document) []string {
	paths := make([]string, 0, len(documents))
	for _, document := range documents {
		paths = append(paths, document.Path)
	}

	return paths
}

func TestLoad(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTopicFile(t, root, "CLAUDE.md", "# Topic\n")
	writeTopicFile(t, root, "AGENTS.md", "# Agents\n")
	writeTopicFile(t, root, "log.md", "# Log\n")
	writeTopicFile(t, root, "raw/articles/a.md", "---\ntitle: Alpha\nsource_kind: article\n---\nAlpha body\n")
	writeTopicFile(t, root, "raw/articles/b.md", "---\nsource_kind: article\n---\n# Heading Title\n\nBeta body\n")
	writeTopicFile(t, root, "raw/articles/untitled.md", "plain body\n")
	writeTopicFile(t, root, "raw/articles/broken.md", "---\ntitle: [broken\n---\nbody\n")
	writeTopicFile(t, root, "raw/articles/unterminated.md", "---\ntitle: x\n")
	writeTopicFile(t, root, "raw/articles/notes.txt", "not markdown\n")
	writeTopicFile(t, root, "raw/_quarantine/articles/junk.md", "---\ntitle: Junk\n---\njunk\n")
	writeTopicFile(t, root, "raw/codebase/files/x.md", "---\ntitle: Code\n---\ncode\n")
	writeTopicFile(t, root, "raw/other/code-index.md", "---\ntitle: Index\nsource_kind: codebase-language-index\n---\n")
	writeTopicFile(t, root, "raw/private/secret.md", "---\ntitle: Secret\n---\n")
	writeTopicFile(t, root, "raw/private/deep/secret2.md", "---\ntitle: Secret2\n---\n")
	writeTopicFile(t, root, "wiki/concepts/Agents.md", "---\ntitle: Agents\naliases:\n    - AI agents\n---\nArticle body\n")
	writeTopicFile(t, root, "wiki/concepts/nested/Deep.md", "---\ntitle: Deep\n---\n")
	writeTopicFile(t, root, "wiki/codebase/concepts/Code.md", "---\ntitle: Code concept\n---\n")
	writeTopicFile(t, root, "wiki/index/Dashboard.md", "---\ntitle: Dashboard\n---\n")
	writeTopicFile(t, root, "outputs/queries/q.md", "---\ntitle: Q\n---\n")
	writeTopicFile(t, root, "bases/b.md", "---\ntitle: B\n---\n")
	writeTopicFile(t, root, ".decisions/x.md", "---\ntitle: X\n---\n")

	loaded, err := corpus.Load(root, corpus.LoadOptions{Exclude: []string{"raw/private/**"}})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	wantSources := []string{"raw/articles/a.md", "raw/articles/b.md", "raw/articles/untitled.md"}
	if got := documentPaths(loaded.Sources()); !reflect.DeepEqual(got, wantSources) {
		t.Fatalf("sources = %#v, want %#v", got, wantSources)
	}
	wantArticles := []string{"wiki/concepts/Agents.md", "wiki/concepts/nested/Deep.md"}
	if got := documentPaths(loaded.Articles()); !reflect.DeepEqual(got, wantArticles) {
		t.Fatalf("articles = %#v, want %#v", got, wantArticles)
	}
	if got := len(loaded.Documents()); got != 5 {
		t.Fatalf("documents = %d, want 5", got)
	}
	wantSkipped := []corpus.Skip{
		{Path: "raw/articles/broken.md", Reason: "frontmatter: invalid_yaml"},
		{Path: "raw/articles/unterminated.md", Reason: "frontmatter: missing_closing_delimiter"},
	}
	if !reflect.DeepEqual(loaded.Skipped, wantSkipped) {
		t.Fatalf("skipped = %#v, want %#v", loaded.Skipped, wantSkipped)
	}

	alpha := loaded.ByPath("raw/articles/a")
	if alpha == nil || alpha.Title != "Alpha" || alpha.Kind != corpus.KindSource || alpha.Body != "Alpha body\n" {
		t.Fatalf("alpha = %#v", alpha)
	}
	if alpha.BodyHash != corpus.BodyHash("Alpha body\n") || alpha.ModTime.IsZero() || !filepath.IsAbs(alpha.AbsPath) {
		t.Fatalf("alpha metadata = %#v", alpha)
	}
	if got := loaded.ByPath("raw/articles/b.md").Title; got != "Heading Title" {
		t.Fatalf("heading title = %q", got)
	}
	if got := loaded.ByPath("raw/articles/untitled.md").Title; got != "untitled" {
		t.Fatalf("stem title = %q", got)
	}
	agents := loaded.ByPath("wiki/concepts/Agents.md")
	if agents.Kind != corpus.KindArticle || !reflect.DeepEqual(agents.Aliases, []string{"AI agents"}) {
		t.Fatalf("agents = %#v", agents)
	}
	if got := loaded.ByBodyHash(corpus.BodyHash("Alpha body\n")); len(got) != 1 || got[0] != alpha {
		t.Fatalf("ByBodyHash = %#v", got)
	}
	if loaded.ByPath("raw/missing.md") != nil {
		t.Fatal("ByPath(missing) returned a document")
	}

	if _, err := corpus.Load(filepath.Join(root, "missing"), corpus.LoadOptions{}); err == nil {
		t.Fatal("Load(missing) returned nil error")
	}
	if _, err := corpus.Load(filepath.Join(root, "CLAUDE.md"), corpus.LoadOptions{}); err == nil {
		t.Fatal("Load(file) returned nil error")
	}
}

func TestDocumentHelpers(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTopicFile(t, root, "raw/articles/full.md", strings.Join([]string{
		"---",
		"title: Full",
		"summary: '  A summary.  '",
		"criterion: Must discuss things.",
		"concepts:",
		"    - '[[Agents]]'",
		"    - '[[wiki/concepts/Tools|Tools]]'",
		"entities: [OpenAI, MCP]",
		"questions:",
		"    - What is it?",
		"genre: paper",
		"depth: 2.5",
		"relevance: core",
		"triage: kept",
		"locked: true",
		"source_url: https://www.Example.com/post?id=1",
		"source_kind: article",
		"ingest_batch: url-2026-09-24-abc",
		"ingest_query: agents",
		"---",
		"body",
	}, "\n"))
	writeTopicFile(t, root, "raw/articles/empty.md", "---\nsummary:\ndepth: deep\n---\nbody\n")
	writeTopicFile(t, root, "raw/articles/intdepth.md", "---\ndepth: 3\nentities: Single\n---\nbody\n")

	loaded, err := corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	full := loaded.ByPath("raw/articles/full.md")
	checks := map[string][2]any{
		"Summary":    {full.Summary(), "A summary."},
		"Criterion":  {full.Criterion(), "Must discuss things."},
		"Concepts":   {full.Concepts(), []string{"Agents", "wiki/concepts/Tools"}},
		"Entities":   {full.Entities(), []string{"OpenAI", "MCP"}},
		"Questions":  {full.Questions(), []string{"What is it?"}},
		"Genre":      {full.Genre(), "paper"},
		"Relevance":  {full.Relevance(), "core"},
		"Triage":     {full.Triage(), "kept"},
		"Locked":     {full.Locked(), true},
		"SourceURL":  {full.SourceURL(), "https://www.Example.com/post?id=1"},
		"SourceKind": {full.SourceKind(), "article"},
	}
	for name, pair := range checks {
		if !reflect.DeepEqual(pair[0], pair[1]) {
			t.Errorf("%s = %#v, want %#v", name, pair[0], pair[1])
		}
	}
	if depth, ok := full.Depth(); !ok || depth != 2.5 {
		t.Errorf("Depth = %v, %v", depth, ok)
	}

	empty := loaded.ByPath("raw/articles/empty.md")
	if empty.Summary() != "" || empty.Locked() || len(empty.Concepts()) != 0 {
		t.Errorf("empty helpers = %q %v %#v", empty.Summary(), empty.Locked(), empty.Concepts())
	}
	if _, ok := empty.Depth(); ok {
		t.Error("Depth(deep) reported ok")
	}

	intDepth := loaded.ByPath("raw/articles/intdepth.md")
	if depth, ok := intDepth.Depth(); !ok || depth != 3 {
		t.Errorf("int Depth = %v, %v", depth, ok)
	}
	if got := intDepth.Entities(); !reflect.DeepEqual(got, []string{"Single"}) {
		t.Errorf("single-string entities = %#v", got)
	}

	wantProvenance := map[string]any{
		"path":         "raw/articles/full.md",
		"source_kind":  "article",
		"source_host":  "example.com",
		"ingest_batch": "url-2026-09-24-abc",
		"ingest_query": "agents",
	}
	if got := full.Provenance(); !reflect.DeepEqual(got, wantProvenance) {
		t.Fatalf("Provenance = %#v, want %#v", got, wantProvenance)
	}
	if got := empty.Provenance(); !reflect.DeepEqual(got, map[string]any{"path": "raw/articles/empty.md"}) {
		t.Fatalf("empty Provenance = %#v", got)
	}

	state := corpus.DocumentState(full, 100, nil)
	if state["title"] != "Full" || state["excerpt"] != "body" || !reflect.DeepEqual(state["provenance"], wantProvenance) {
		t.Fatalf("DocumentState = %#v", state)
	}
}

func TestReadDocument(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTopicFile(t, root, "raw/a.md", "---\ntitle: A\n---\nbody\n")
	document, err := corpus.ReadDocument(root, "raw/a.md", corpus.KindSource)
	if err != nil || document.Title != "A" || document.Path != "raw/a.md" {
		t.Fatalf("ReadDocument = %#v, %v", document, err)
	}
	if _, err := corpus.ReadDocument(root, "raw/missing.md", corpus.KindSource); err == nil {
		t.Fatal("ReadDocument(missing) returned nil error")
	}
}

func TestMatchGlob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"raw/brand-audit-*/**", "raw/brand-audit-2026/B016/firecrawl.md", true},
		{"raw/brand-audit-*/**", "raw/brand-audit-2026", true},
		{"raw/brand-audit-*/**", "raw/articles/x.md", false},
		{"raw/**", "raw/a/b/c.md", true},
		{"**/*.md", "a.md", true},
		{"**/*.md", "raw/a/b.md", true},
		{"**/*.md", "raw/a/b.txt", false},
		{"raw/*.md", "raw/a/b.md", false},
		{"raw/*.md", "raw/a.md", true},
		{"raw/?.md", "raw/a.md", true},
		{"raw/**/b.md", "raw/b.md", true},
		{"raw/**/b.md", "raw/x/y/b.md", true},
		{"raw/**/**/b.md", "raw/x/b.md", true},
		{"./raw/a.md", "/raw/a.md/", true},
		{"raw/[", "raw/[", false},
		{"", "", true},
		{"", "raw", false},
	}
	for _, tc := range tests {
		if got := corpus.MatchGlob(tc.pattern, tc.path); got != tc.want {
			t.Errorf("MatchGlob(%q, %q) = %v, want %v", tc.pattern, tc.path, got, tc.want)
		}
	}

	if !corpus.ValidGlob("raw/**/[ab]*.md") || corpus.ValidGlob("raw/[") {
		t.Fatal("ValidGlob mismatch")
	}
}

func TestHashes(t *testing.T) {
	t.Parallel()

	if corpus.BodyHash("a\r\nb\r\n") != corpus.BodyHash("a\nb\n") {
		t.Fatal("BodyHash must normalize CRLF")
	}
	if corpus.BodyHash("a") == corpus.BodyHash("b") {
		t.Fatal("BodyHash collision")
	}
	if len(corpus.BodyHash("")) != 64 {
		t.Fatal("BodyHash must be sha256 hex")
	}

	equal := [][2]any{
		{" summary ", "summary"},
		{[]any{"a", " b"}, []string{"a", "b"}},
		{2, 2.0},
		{int64(3), uint8(3)},
		{time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC), "2026-09-24"},
		{map[string]any{"b": 1, "a": "x "}, map[string]any{"a": "x", "b": 1.0}},
		{[]int{1, 2}, []any{1, 2}},
	}
	for _, pair := range equal {
		if corpus.ValueHash(pair[0]) != corpus.ValueHash(pair[1]) {
			t.Errorf("ValueHash(%#v) != ValueHash(%#v)", pair[0], pair[1])
		}
	}
	different := [][2]any{
		{[]string{"a", "b"}, []string{"b", "a"}},
		{"1", 1},
		{true, "true"},
		{nil, ""},
	}
	for _, pair := range different {
		if corpus.ValueHash(pair[0]) == corpus.ValueHash(pair[1]) {
			t.Errorf("ValueHash(%#v) == ValueHash(%#v)", pair[0], pair[1])
		}
	}
}

func TestExcerpt(t *testing.T) {
	t.Parallel()

	section := func(title string, words int, term string) string {
		return "## " + title + "\n" + strings.Repeat("filler ", words) + term + "\n"
	}
	body := "Intro paragraph.\n" +
		section("One", 30, "") +
		section("Two", 30, "") +
		section("Three", 30, "quantum") +
		section("Four", 30, "") +
		section("Five", 30, "")

	if got := corpus.Excerpt("short body", 100, nil); got != "short body" {
		t.Fatalf("short Excerpt = %q", got)
	}
	if got := corpus.Excerpt(body, 0, nil); got != body {
		t.Fatal("Excerpt with no budget must return the body")
	}

	got := corpus.Excerpt(body, 170, []string{"Quantum"})
	for _, want := range []string{"Intro paragraph.", "## One", "## Three", "[... 1 section omitted ...]"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Excerpt missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "## Two") {
		t.Fatalf("Excerpt kept a non-matching section before the matching one:\n%s", got)
	}
	if strings.Index(got, "## One") > strings.Index(got, "## Three") {
		t.Fatalf("Excerpt is not in document order:\n%s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "omitted ...]") {
		t.Fatalf("Excerpt must end with an omission marker when the tail is dropped:\n%s", got)
	}
	if corpus.EstimateTokens(got) > 170+20 {
		t.Fatalf("Excerpt exceeds budget: %d tokens", corpus.EstimateTokens(got))
	}

	noTerms := corpus.Excerpt(body, 170, nil)
	if !strings.Contains(noTerms, "## Two") || strings.Contains(noTerms, "## Five") {
		t.Fatalf("Excerpt without terms must fill in document order:\n%s", noTerms)
	}

	huge := "# Title\n" + strings.Repeat("word ", 2000) + "\n## Next\nmore\n"
	truncated := corpus.Excerpt(huge, 50, nil)
	if !strings.Contains(truncated, "[... truncated ...]") || !strings.Contains(truncated, "[... 1 section omitted ...]") {
		t.Fatalf("oversized head must be truncated with markers:\n%s", truncated)
	}

	fenced := "# T\n```\n# not a heading\n" + strings.Repeat("x ", 300) + "\n```\n## Real\n" + strings.Repeat("y ", 300) + "\n"
	if got := corpus.Excerpt(fenced, 200, nil); strings.Contains(got, "sections omitted") && !strings.Contains(got, "# not a heading") {
		t.Fatalf("heading inside a fence must not split sections:\n%s", got)
	}

	if corpus.EstimateTokens("abcd") != 1 || corpus.EstimateTokens("abcde") != 2 || corpus.EstimateTokens("ção") != 1 {
		t.Fatal("EstimateTokens must be runes/4 rounded up")
	}
}

func TestHead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		body string
		max  int
		want string
	}{
		{"See [[Target|shown]] and [[Other]] and ![[Embed#H]].", 100, "See shown and Other and Embed."},
		{"  ação longa  ", 3, "açã"},
		{"[[a]]bc", 2, "ab"},
	}
	for _, tc := range tests {
		if got := corpus.Head(tc.body, tc.max); got != tc.want {
			t.Errorf("Head(%q, %d) = %q, want %q", tc.body, tc.max, got, tc.want)
		}
	}
}
