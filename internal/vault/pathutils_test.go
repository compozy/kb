package vault_test

import (
	"testing"

	"github.com/user/go-devstack/internal/models"
	"github.com/user/go-devstack/internal/vault"
)

func TestCreateFileIDDeterministic(t *testing.T) {
	t.Parallel()

	gotA := vault.CreateFileID("src/main.go")
	gotB := vault.CreateFileID("src/main.go")

	if gotA != gotB {
		t.Fatalf("expected deterministic file id, got %q and %q", gotA, gotB)
	}

	if gotA != "file:src/main.go" {
		t.Fatalf("unexpected file id %q", gotA)
	}
}

func TestCreateSymbolIDDeterministicAndUnique(t *testing.T) {
	t.Parallel()

	symbol := models.SymbolNode{
		FilePath:   "src/main.go",
		Name:       "RunServer",
		SymbolKind: "function",
		StartLine:  10,
		EndLine:    24,
	}

	other := symbol
	other.Name = "RunClient"

	gotA := vault.CreateSymbolID(symbol)
	gotB := vault.CreateSymbolID(symbol)
	gotOther := vault.CreateSymbolID(other)

	if gotA != gotB {
		t.Fatalf("expected deterministic symbol id, got %q and %q", gotA, gotB)
	}

	if gotA == gotOther {
		t.Fatalf("expected unique symbol ids, got %q", gotA)
	}
}

func TestToPosixPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "backslashes", input: `src\pkg\main.go`, want: "src/pkg/main.go"},
		{name: "trailing slash", input: "src/pkg/", want: "src/pkg"},
		{name: "root preserved", input: "/", want: "/"},
		{name: "windows drive root preserved", input: `C:\`, want: "C:/"},
		{name: "empty", input: "", want: ""},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := vault.ToPosixPath(testCase.input); got != testCase.want {
				t.Fatalf("ToPosixPath(%q) = %q, want %q", testCase.input, got, testCase.want)
			}
		})
	}
}

func TestIsPathInside(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		parent string
		target string
		want   bool
	}{
		{name: "child path", parent: "/repo", target: "/repo/src/main.go", want: true},
		{name: "same path", parent: "/repo/src", target: "/repo/src", want: true},
		{name: "outside path", parent: "/repo", target: "/other/main.go", want: false},
		{name: "relative child", parent: ".", target: "./src/main.go", want: true},
		{name: "empty paths", parent: "", target: "", want: true},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := vault.IsPathInside(testCase.parent, testCase.target); got != testCase.want {
				t.Fatalf("IsPathInside(%q, %q) = %t, want %t", testCase.parent, testCase.target, got, testCase.want)
			}
		})
	}
}

func TestDocumentPathDerivationHelpers(t *testing.T) {
	t.Parallel()

	symbol := models.SymbolNode{
		FilePath:  "src/services/api.go",
		Name:      "Run Server",
		StartLine: 42,
	}

	if got := vault.GetRawFileDocumentPath("src/main.go"); got != "raw/codebase/files/src/main.go.md" {
		t.Fatalf("unexpected raw file document path %q", got)
	}

	if got := vault.GetRawFileDocumentPath("/tmp/my file.go"); got != "raw/codebase/files/tmp/my file.go.md" {
		t.Fatalf("unexpected absolute raw file document path %q", got)
	}

	if got := vault.GetRawSymbolDocumentPath(symbol); got != "raw/codebase/symbols/run-server--src-services-api-go-l42.md" {
		t.Fatalf("unexpected raw symbol document path %q", got)
	}

	if got := vault.GetRawDirectoryIndexPath("."); got != "raw/codebase/indexes/directories/root.md" {
		t.Fatalf("unexpected root directory index path %q", got)
	}

	if got := vault.GetRawDirectoryIndexPath("src/services"); got != "raw/codebase/indexes/directories/src/services.md" {
		t.Fatalf("unexpected directory index path %q", got)
	}

	if got := vault.GetRawLanguageIndexPath("go"); got != "raw/codebase/indexes/languages/go.md" {
		t.Fatalf("unexpected language index path %q", got)
	}

	if got := vault.GetWikiConceptPath("Codebase Overview"); got != "wiki/concepts/Kodebase - Codebase Overview.md" {
		t.Fatalf("unexpected wiki concept path %q", got)
	}

	if got := vault.GetWikiIndexPath("Dashboard"); got != "wiki/index/Dashboard.md" {
		t.Fatalf("unexpected wiki index path %q", got)
	}

	if got := vault.GetBaseFilePath("symbol-explorer"); got != "bases/symbol-explorer.base" {
		t.Fatalf("unexpected base file path %q", got)
	}
}

func TestTopicHelpers(t *testing.T) {
	t.Parallel()

	if got := vault.SlugifySegment("  API Client v2  "); got != "api-client-v2" {
		t.Fatalf("unexpected slug %q", got)
	}

	if got := vault.HumanizeSlug("api-client-v2"); got != "Api Client V2" {
		t.Fatalf("unexpected humanized slug %q", got)
	}

	if got := vault.DeriveTopicSlug(`C:\Projects\Kodebase Go`); got != "kodebase-go" {
		t.Fatalf("unexpected topic slug %q", got)
	}

	if got := vault.DeriveTopicTitle("kodebase-go"); got != "Kodebase Go" {
		t.Fatalf("unexpected topic title %q", got)
	}

	if got := vault.DeriveTopicTitle(""); got != "Knowledge Base" {
		t.Fatalf("unexpected fallback topic title %q", got)
	}

	if got := vault.DeriveTopicDomain("kodebase-go"); got != "kodebase-go" {
		t.Fatalf("unexpected topic domain %q", got)
	}
}

func TestTopicWikiLinkHelpers(t *testing.T) {
	t.Parallel()

	if got := vault.StripMarkdownExtension("wiki/concepts/Kodebase - Codebase Overview.md"); got != "wiki/concepts/Kodebase - Codebase Overview" {
		t.Fatalf("unexpected stripped markdown path %q", got)
	}

	if got := vault.ToTopicWikiLink("topic-slug", "wiki/concepts/Kodebase - Codebase Overview.md", "Overview"); got != "[[topic-slug/wiki/concepts/Kodebase - Codebase Overview|Overview]]" {
		t.Fatalf("unexpected labeled wiki link %q", got)
	}

	if got := vault.ToTopicWikiLink("topic-slug", "wiki/index/Dashboard.md", ""); got != "[[topic-slug/wiki/index/Dashboard]]" {
		t.Fatalf("unexpected wiki link %q", got)
	}
}

func TestPathHelpersHandleEmptyInputs(t *testing.T) {
	t.Parallel()

	emptySymbol := models.SymbolNode{}

	if got := vault.CreateFileID(""); got != "file:" {
		t.Fatalf("unexpected empty file id %q", got)
	}

	if got := vault.CreateExternalID(""); got != "external:" {
		t.Fatalf("unexpected empty external id %q", got)
	}

	if got := vault.CreateSymbolID(emptySymbol); got != "symbol::item:item:0:0" {
		t.Fatalf("unexpected empty symbol id %q", got)
	}

	if got := vault.GetRawFileDocumentPath(""); got != "raw/codebase/files/.md" {
		t.Fatalf("unexpected empty raw file path %q", got)
	}

	if got := vault.GetRawDirectoryIndexPath(""); got != "raw/codebase/indexes/directories/root.md" {
		t.Fatalf("unexpected empty directory index path %q", got)
	}
}
