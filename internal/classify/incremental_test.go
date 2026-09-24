package classify

import (
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func TestSkipReason(t *testing.T) {
	t.Parallel()
	vault, root := newTestTopic(t)
	setTestContract(t, root, nil)
	s := openTestSession(t, vault, "http://127.0.0.1:1")
	complete := map[string]any{
		"genre": "paper", "summary": "S.", "entities": []string{"E"}, "questions": []string{"Q?"},
		"relevance": "core", "concepts": []string{},
	}
	doc := &corpus.Document{Path: "raw/a.md", Kind: corpus.KindSource, BodyHash: "h1", Frontmatter: complete}
	current := &corpus.StateRow{Path: "raw/a.md", BodyHash: "h1", Contract: s.ContractHash(), Banks: bankVersions(SourceBanks())}
	staleBank := &corpus.StateRow{Path: "raw/a.md", BodyHash: "h1", Contract: s.ContractHash(), Banks: map[string]string{"relevance": "old"}}
	changedBody := &corpus.StateRow{Path: "raw/a.md", BodyHash: "h0", Contract: s.ContractHash(), Banks: bankVersions(SourceBanks())}
	otherContract := &corpus.StateRow{Path: "raw/a.md", BodyHash: "h1", Contract: "other", Banks: bankVersions(SourceBanks())}
	missing := &corpus.Document{Path: "raw/a.md", Kind: corpus.KindSource, BodyHash: "h1", Frontmatter: map[string]any{"genre": "paper"}}
	locked := &corpus.Document{Path: "raw/a.md", Kind: corpus.KindSource, BodyHash: "h1", Frontmatter: map[string]any{"locked": true}}
	stub := &corpus.Document{Path: "wiki/concepts/x.md", Kind: corpus.KindArticle, BodyHash: "h1", Frontmatter: map[string]any{"stage": "stub"}}
	article := &corpus.Document{Path: "wiki/concepts/y.md", Kind: corpus.KindArticle, BodyHash: "h1", Frontmatter: map[string]any{}}
	articleRow := &corpus.StateRow{Path: "wiki/concepts/y.md", BodyHash: "h1", Contract: "old-contract", Banks: bankVersions(ArticleBanks())}

	tests := []struct {
		name string
		doc  *corpus.Document
		row  *corpus.StateRow
		opts Options
		want string
	}{
		{name: "no state row", doc: doc, want: ""},
		{name: "current row", doc: doc, row: current, want: "unchanged"},
		{name: "bank version changed", doc: doc, row: staleBank, want: ""},
		{name: "body changed", doc: doc, row: changedBody, want: ""},
		{name: "contract changed", doc: doc, row: otherContract, want: ""},
		{name: "all ignores the state row", doc: doc, row: current, opts: Options{All: true}, want: ""},
		{name: "only-missing with every key", doc: doc, row: current, opts: Options{OnlyMissing: true}, want: "complete"},
		{name: "only-missing with a stale bank but every key", doc: doc, row: staleBank, opts: Options{OnlyMissing: true}, want: "complete"},
		{name: "only-missing with a missing key", doc: missing, row: current, opts: Options{OnlyMissing: true}, want: ""},
		{name: "only-missing without row", doc: doc, opts: Options{OnlyMissing: true}, want: ""},
		{name: "locked", doc: locked, opts: Options{All: true}, want: "locked"},
		{name: "stub article", doc: stub, want: "stub"},
		{name: "article ignores contract changes", doc: article, row: articleRow, want: "unchanged"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, err := corpus.Load(root, corpus.LoadOptions{})
			if err != nil {
				t.Fatal(err)
			}
			r := newRun(s, c, tt.opts)
			r.vocab = newVocabulary([]*corpus.Document{testArticle("z", "Zed", nil, "Crit.", "", "")}, stemLink)
			if tt.row != nil {
				r.rows[tt.doc.Path] = tt.row
			}
			if got := r.skipReason(tt.doc); got != tt.want {
				t.Fatalf("skipReason = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNeedsLiteral(t *testing.T) {
	t.Parallel()
	kbValue := "Written by kb."
	doc := func(value any, bodyHash string) *corpus.Document {
		values := map[string]any{}
		if value != nil {
			values["summary"] = value
		}
		return &corpus.Document{Path: "raw/a.md", BodyHash: bodyHash, Frontmatter: values}
	}
	row := &corpus.StateRow{Path: "raw/a.md", BodyHash: "old", Written: map[string]string{"summary": corpus.ValueHash(kbValue)}}
	tests := []struct {
		name string
		doc  *corpus.Document
		row  *corpus.StateRow
		want bool
	}{
		{name: "missing", doc: doc(nil, "new"), want: true},
		{name: "empty string", doc: doc("  ", "new"), want: true},
		{name: "present without row (user-owned)", doc: doc("User summary.", "new"), want: false},
		{name: "kb-written, body unchanged", doc: doc(kbValue, "old"), row: row, want: false},
		{name: "kb-written, body changed", doc: doc(kbValue, "new"), row: row, want: true},
		{name: "user-edited, body changed", doc: doc("Edited.", "new"), row: row, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := needsLiteral(tt.doc, tt.row, "summary"); got != tt.want {
				t.Fatalf("needsLiteral = %v, want %v", got, tt.want)
			}
		})
	}
}
