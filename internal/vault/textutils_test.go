package vault_test

import (
	"testing"

	"github.com/user/go-devstack/internal/vault"
)

func TestNormalizeComment(t *testing.T) {
	t.Parallel()

	raw := `/**
 * Builds the vault.
 *
 * Keeps structure stable.
 */`

	want := "Builds the vault.\n\nKeeps structure stable."
	if got := vault.NormalizeComment(raw); got != want {
		t.Fatalf("NormalizeComment() = %q, want %q", got, want)
	}
}

func TestExtractLeadingCommentFromGoSource(t *testing.T) {
	t.Parallel()

	source := `// Package main wires the CLI.
// It exposes the root command.
package main

func main() {}
`

	want := "Package main wires the CLI.\nIt exposes the root command."
	if got := vault.ExtractLeadingComment(source); got != want {
		t.Fatalf("ExtractLeadingComment() = %q, want %q", got, want)
	}
}

func TestExtractLeadingCommentFromTSSource(t *testing.T) {
	t.Parallel()

	source := `/**
 * Builds the topic graph.
 * Preserves the first block only.
 */
export function buildGraph() {}
`

	want := "Builds the topic graph.\nPreserves the first block only."
	if got := vault.ExtractLeadingComment(source); got != want {
		t.Fatalf("ExtractLeadingComment() = %q, want %q", got, want)
	}
}

func TestExtractLeadingCommentIgnoresNonLeadingComments(t *testing.T) {
	t.Parallel()

	source := `package main

// Not a leading comment.
func main() {}
`

	if got := vault.ExtractLeadingComment(source); got != "" {
		t.Fatalf("expected no leading comment, got %q", got)
	}
}

func TestStripQuotes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "double quotes", input: `"fmt"`, want: "fmt"},
		{name: "single quotes", input: "'./dep'", want: "./dep"},
		{name: "backticks", input: "`raw`", want: "raw"},
		{name: "mismatched quotes", input: `"dep'`, want: "dep"},
		{name: "empty", input: "", want: ""},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := vault.StripQuotes(testCase.input); got != testCase.want {
				t.Fatalf("StripQuotes(%q) = %q, want %q", testCase.input, got, testCase.want)
			}
		})
	}
}
