package resolve_test

import (
	"reflect"
	"testing"

	"github.com/compozy/kb/internal/resolve"
)

func TestBodyLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want []resolve.Link
	}{
		{
			name: "wikilink with display and heading",
			body: "See [[Note#Part|shown]] here.",
			want: []resolve.Link{{Target: "Note", Display: "shown", Heading: "Part", Start: 4, End: 23, Kind: resolve.KindWikilink}},
		},
		{
			name: "embed",
			body: "![[Diagram]]",
			want: []resolve.Link{{Target: "Diagram", Start: 0, End: 12, Kind: resolve.KindWikilink, Embed: true}},
		},
		{
			name: "markdown relative md link and unescape",
			body: "Read [the note](../wiki/My%20Note.md#intro).",
			want: []resolve.Link{{Target: "../wiki/My Note", Display: "the note", Heading: "intro", Start: 5, End: 43, Kind: resolve.KindMarkdown}},
		},
		{
			name: "markdown link without extension and angle brackets",
			body: "[a](folder/page) [b](<with space.md>)",
			want: []resolve.Link{
				{Target: "folder/page", Display: "a", Start: 0, End: 16, Kind: resolve.KindMarkdown},
				{Target: "with space", Display: "b", Start: 17, End: 37, Kind: resolve.KindMarkdown},
			},
		},
		{
			name: "external, anchor, mailto and image links are skipped",
			body: "[x](https://example.com/a.md) [y](#local) [z](mailto:a@b.c) ![i](img.png) [w](file.pdf)",
		},
		{
			name: "bracketed label followed by destination is a markdown link",
			body: "[[Opens in a new window]](https://example.com)",
		},
		{
			name: "fenced and inline code are skipped",
			body: "```\n[[InFence]]\n```\n`[[Inline]]` ``[[Double `tick`]]``\n~~~go\n[x](a.md)\n~~~\n[[After]]",
			want: []resolve.Link{{Target: "After", Start: 75, End: 84, Kind: resolve.KindWikilink}},
		},
		{
			name: "unterminated fence swallows the rest",
			body: "[[Before]]\n```\n[[Hidden]]\n",
			want: []resolve.Link{{Target: "Before", Start: 0, End: 10, Kind: resolve.KindWikilink}},
		},
		{
			name: "document order across kinds",
			body: "[m](m.md) then [[w]]",
			want: []resolve.Link{
				{Target: "m", Display: "m", Start: 0, End: 9, Kind: resolve.KindMarkdown},
				{Target: "w", Start: 15, End: 20, Kind: resolve.KindWikilink},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := resolve.BodyLinks(tc.body)
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("BodyLinks = %#v\nwant %#v", got, tc.want)
			}
			for _, link := range got {
				if tc.body[link.Start] != '[' && tc.body[link.Start] != '!' {
					t.Fatalf("link %#v does not start at markup", link)
				}
			}
		})
	}
}

func TestFrontmatterLinks(t *testing.T) {
	t.Parallel()

	values := map[string]any{
		"title":       "No [links] here",
		"up":          "[[Parent|P]]",
		"sources":     []string{"[[raw/articles/a]]", "https://example.com", "[[b#h]]"},
		"related":     []any{"[[X]]", 3, "text [[Y]] and [[Z]]"},
		"nested":      map[string]any{"k": "[[Ignored]]"},
		"informed_by": []string{},
	}
	got := resolve.FrontmatterLinks(values)
	want := map[string][]string{
		"up":      {"Parent"},
		"sources": {"raw/articles/a", "b"},
		"related": {"X", "Y", "Z"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FrontmatterLinks = %#v, want %#v", got, want)
	}
}

func TestCodeRanges(t *testing.T) {
	t.Parallel()

	text := "a `b` c\n```\ncode\n```\nd"
	got := resolve.CodeRanges(text)
	want := []resolve.Range{{Start: 2, End: 5}, {Start: 8, End: 21}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CodeRanges = %#v, want %#v", got, want)
	}
	if !resolve.InRanges(got, 3, 4) || resolve.InRanges(got, 0, 2) {
		t.Fatal("InRanges mismatch")
	}
}
