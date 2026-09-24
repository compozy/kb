package frontmatter_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/frontmatter"
)

func TestEditKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		updates []frontmatter.KeyUpdate
		want    string
	}{
		{
			name:    "replaces scalar and keeps comments and order",
			input:   "---\n# leading comment\ntitle: Example # inline\nsummary: old\n# between\ntags:\n    - a\n---\n# Body\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "new summary"}},
			want:    "---\n# leading comment\ntitle: Example # inline\nsummary: new summary\n# between\ntags:\n    - a\n---\n# Body\n",
		},
		{
			name:    "appends new keys in given order",
			input:   "---\ntitle: Example\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "triage", Value: "kept"}, {Key: "genre", Value: "paper"}},
			want:    "---\ntitle: Example\ntriage: kept\ngenre: paper\n---\nbody\n",
		},
		{
			name:    "quotes wikilink list values",
			input:   "---\ntitle: Example\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "related", Value: []string{"[[Alpha]]", "[[Beta|B]]"}}},
			want:    "---\ntitle: Example\nrelated:\n    - '[[Alpha]]'\n    - '[[Beta|B]]'\n---\nbody\n",
		},
		{
			name:    "replaces 4-space block list",
			input:   "---\ntags:\n    - a\n    - b\nconcepts:\n    - '[[Old]]'\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "concepts", Value: []string{"[[New]]"}}},
			want:    "---\ntags:\n    - a\n    - b\nconcepts:\n    - '[[New]]'\ntitle: T\n---\nbody\n",
		},
		{
			name:    "replaces column-0 block list",
			input:   "---\naliases:\n- A\n- B\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "aliases", Value: []string{"A", "B", "C"}}},
			want:    "---\naliases:\n    - A\n    - B\n    - C\ntitle: T\n---\nbody\n",
		},
		{
			name:    "replaces flow list",
			input:   "---\nconcepts: [\"[[A]]\", \"[[B]]\"]\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "concepts", Value: []string{"[[C]]"}}},
			want:    "---\nconcepts:\n    - '[[C]]'\ntitle: T\n---\nbody\n",
		},
		{
			name:    "replaces block scalar with blank lines inside",
			input:   "---\nsummary: |\n  line one\n\n  line two\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "short"}},
			want:    "---\nsummary: short\ntitle: T\n---\nbody\n",
		},
		{
			name:    "encodes multi-line string as literal block",
			input:   "---\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "first\nsecond"}},
			want:    "---\ntitle: T\nsummary: |-\n    first\n    second\n---\nbody\n",
		},
		{
			name:    "does not confuse key prefixes",
			input:   "---\ntags_extra: keep\ntags: [a]\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "tags", Value: []string{"b"}}},
			want:    "---\ntags_extra: keep\ntags:\n    - b\n---\nbody\n",
		},
		{
			name:    "deletes key block and keeps trailing blank line",
			input:   "---\ntitle: T\nsummary: >\n  folded text\n\n# comment\nextra: 1\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Delete: true}},
			want:    "---\ntitle: T\n\n# comment\nextra: 1\n---\nbody\n",
		},
		{
			name:    "delete of absent key is a no-op",
			input:   "---\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Delete: true}},
			want:    "---\ntitle: T\n---\nbody\n",
		},
		{
			name:    "replaces quoted keys",
			input:   "---\n\"summary\": old\n'genre': x\ntitle: T\n---\nbody\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "new"}, {Key: "genre", Value: "paper"}},
			want:    "---\nsummary: new\ngenre: paper\ntitle: T\n---\nbody\n",
		},
		{
			name:    "preserves CRLF line endings",
			input:   "---\r\ntitle: T\r\nsummary: old\r\n---\r\nbody\r\nline\r\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "new"}, {Key: "concepts", Value: []string{"[[A]]"}}},
			want:    "---\r\ntitle: T\r\nsummary: new\r\nconcepts:\r\n    - '[[A]]'\r\n---\r\nbody\r\nline\r\n",
		},
		{
			name:    "creates frontmatter when missing",
			input:   "# Body\n\ntext\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "s"}, {Key: "depth", Value: 2}},
			want:    "---\nsummary: s\ndepth: 2\n---\n# Body\n\ntext\n",
		},
		{
			name:    "no frontmatter and only deletes leaves document untouched",
			input:   "# Body\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Delete: true}},
			want:    "# Body\n",
		},
		{
			name:    "empty frontmatter gets appended keys",
			input:   "---\n---\nbody",
			updates: []frontmatter.KeyUpdate{{Key: "triage", Value: "kept"}},
			want:    "---\ntriage: kept\n---\nbody",
		},
		{
			name:    "quotes strings that look like other scalars",
			input:   "---\ntitle: T\n---\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "true"}, {Key: "genre", Value: "123"}},
			want:    "---\ntitle: T\nsummary: \"true\"\ngenre: \"123\"\n---\n",
		},
		{
			name:    "body bytes stay identical even when they look like yaml",
			input:   "---\ntitle: T\n---\nsummary: not frontmatter\n---\n",
			updates: []frontmatter.KeyUpdate{{Key: "summary", Value: "real"}},
			want:    "---\ntitle: T\nsummary: real\n---\nsummary: not frontmatter\n---\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := frontmatter.EditKeys(tc.input, tc.updates)
			if err != nil {
				t.Fatalf("EditKeys returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("EditKeys mismatch\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func TestEditKeysRoundTripPreservesUntouchedKeys(t *testing.T) {
	t.Parallel()

	input := strings.Join([]string{
		"---",
		"# kept comment",
		"title: \"Quoted: Title\"",
		"tags:",
		"    - kb",
		"    - raw",
		"summary: old",
		"notes: |",
		"  multi",
		"",
		"  line",
		"url: https://example.com/a:b",
		"concepts: ['[[A]]']",
		"---",
		"# Body [[Link]]",
		"",
	}, "\n")

	updates := []frontmatter.KeyUpdate{
		{Key: "summary", Value: "fresh"},
		{Key: "concepts", Value: []string{"[[A]]", "[[B]]"}},
		{Key: "aliases", Value: []string{"AI", "Inteligência: artificial"}},
	}
	got, err := frontmatter.EditKeys(input, updates)
	if err != nil {
		t.Fatalf("EditKeys returned error: %v", err)
	}

	values, body, err := frontmatter.Parse(got)
	if err != nil {
		t.Fatalf("Parse(edited) returned error: %v", err)
	}
	wantValues := map[string]any{
		"title":    "Quoted: Title",
		"tags":     []string{"kb", "raw"},
		"summary":  "fresh",
		"notes":    "multi\n\nline\n",
		"url":      "https://example.com/a:b",
		"concepts": []string{"[[A]]", "[[B]]"},
		"aliases":  []string{"AI", "Inteligência: artificial"},
	}
	if !reflect.DeepEqual(values, wantValues) {
		t.Fatalf("values = %#v, want %#v", values, wantValues)
	}
	_, originalBody, err := frontmatter.Parse(input)
	if err != nil {
		t.Fatalf("Parse(input) returned error: %v", err)
	}
	if body != originalBody {
		t.Fatalf("body changed: %q -> %q", originalBody, body)
	}

	for _, key := range []string{"title", "tags", "notes", "url"} {
		before, ok := frontmatter.TopLevelKeyText(input, key)
		if !ok {
			t.Fatalf("TopLevelKeyText(input, %q) missing", key)
		}
		after, ok := frontmatter.TopLevelKeyText(got, key)
		if !ok {
			t.Fatalf("TopLevelKeyText(edited, %q) missing", key)
		}
		if before != after {
			t.Fatalf("key %q raw text changed: %q -> %q", key, before, after)
		}
	}
	if !strings.HasPrefix(got, "---\n# kept comment\n") {
		t.Fatalf("leading comment lost: %q", got)
	}
}

func TestEditKeysErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		updates  []frontmatter.KeyUpdate
		wantKind frontmatter.ErrorKind
	}{
		{
			name:     "missing closing delimiter",
			input:    "---\ntitle: T\nbody\n",
			updates:  []frontmatter.KeyUpdate{{Key: "summary", Value: "x"}},
			wantKind: frontmatter.ErrorKindMissingClosingDelimiter,
		},
		{
			name:     "invalid existing yaml",
			input:    "---\ntitle: [unterminated\n---\n",
			updates:  []frontmatter.KeyUpdate{{Key: "summary", Value: "x"}},
			wantKind: frontmatter.ErrorKindInvalidYAML,
		},
		{
			name:     "unsupported value",
			input:    "---\ntitle: T\n---\n",
			updates:  []frontmatter.KeyUpdate{{Key: "summary", Value: func() {}}},
			wantKind: frontmatter.ErrorKindUnsupportedValue,
		},
		{
			name:     "empty key",
			input:    "---\ntitle: T\n---\n",
			updates:  []frontmatter.KeyUpdate{{Key: " ", Value: "x"}},
			wantKind: frontmatter.ErrorKindUnsupportedValue,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := frontmatter.EditKeys(tc.input, tc.updates)
			var fmErr *frontmatter.Error
			if !errors.As(err, &fmErr) {
				t.Fatalf("error = %v, want *frontmatter.Error", err)
			}
			if fmErr.Kind != tc.wantKind {
				t.Fatalf("kind = %q, want %q", fmErr.Kind, tc.wantKind)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantFront string
		wantBody  string
		wantHas   bool
	}{
		{name: "lf", input: "---\na: 1\n---\nbody", wantFront: "a: 1\n", wantBody: "body", wantHas: true},
		{name: "crlf", input: "---\r\na: 1\r\n---\r\nbody", wantFront: "a: 1\r\n", wantBody: "body", wantHas: true},
		{name: "no frontmatter", input: "# body", wantBody: "# body"},
		{name: "unterminated", input: "---\na: 1\n", wantBody: "---\na: 1\n"},
		{name: "empty", input: "---\n---\n", wantHas: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			front, body, has := frontmatter.Split(tc.input)
			if front != tc.wantFront || body != tc.wantBody || has != tc.wantHas {
				t.Fatalf("Split = (%q, %q, %v), want (%q, %q, %v)", front, body, has, tc.wantFront, tc.wantBody, tc.wantHas)
			}
		})
	}
}

func TestKeyBlocks(t *testing.T) {
	t.Parallel()

	input := "---\n# c\ntitle: T\ntags:\n- a\n- b\n\nsummary: >-\n  folded\n  text\n\"quoted key\": v\n---\nbody\n"
	blocks, err := frontmatter.KeyBlocks(input)
	if err != nil {
		t.Fatalf("KeyBlocks returned error: %v", err)
	}
	want := []frontmatter.KeyBlock{
		{Key: "title", Text: "title: T\n"},
		{Key: "tags", Text: "tags:\n- a\n- b\n"},
		{Key: "summary", Text: "summary: >-\n  folded\n  text\n"},
		{Key: "quoted key", Text: "\"quoted key\": v\n"},
	}
	if !reflect.DeepEqual(blocks, want) {
		t.Fatalf("KeyBlocks = %#v, want %#v", blocks, want)
	}

	if blocks, err := frontmatter.KeyBlocks("no frontmatter"); err != nil || blocks != nil {
		t.Fatalf("KeyBlocks(no frontmatter) = %#v, %v", blocks, err)
	}
	if _, err := frontmatter.KeyBlocks("---\ntitle: T\n"); err == nil {
		t.Fatal("KeyBlocks(unterminated) returned nil error")
	}
	if _, ok := frontmatter.TopLevelKeyText(input, "missing"); ok {
		t.Fatal("TopLevelKeyText(missing) reported ok")
	}
}
