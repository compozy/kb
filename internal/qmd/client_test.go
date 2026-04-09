package qmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSearchReturnsErrQMDUnavailableForMissingBinary(t *testing.T) {
	t.Parallel()

	client := NewClient(WithBinaryPath(filepath.Join(t.TempDir(), "missing-qmd")))

	_, err := client.Search(context.Background(), SearchOptions{
		Query: "authentication",
	})
	if !errors.Is(err, ErrQMDUnavailable) {
		t.Fatalf("Search error = %v, want ErrQMDUnavailable", err)
	}
}

func TestIndexAddConstructsExpectedArguments(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"collection add": addOutputFixture,
			"status":         statusOutputFixture,
		},
	})

	client := NewClient(
		WithBinaryPath(binaryPath),
		WithIndexName("task17-index"),
	)

	_, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		VaultPath:      "/tmp/vault",
		CollectionName: "docs",
	})
	if err != nil {
		t.Fatalf("Index returned error: %v", err)
	}

	invocations := readInvocationLog(t, logPath)
	if len(invocations) != 2 {
		t.Fatalf("invocations = %d, want 2", len(invocations))
	}

	expected := []string{"--index", "task17-index", "collection", "add", "/tmp/vault", "--name", "docs"}
	if !reflect.DeepEqual(invocations[0], expected) {
		t.Fatalf("add args = %#v, want %#v", invocations[0], expected)
	}
}

func TestIndexUpdateConstructsExpectedArguments(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"update": updateOutputFixture,
			"status": statusOutputFixture,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))

	_, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationUpdate,
		CollectionName: "docs",
	})
	if err != nil {
		t.Fatalf("Index returned error: %v", err)
	}

	invocations := readInvocationLog(t, logPath)
	if len(invocations) != 2 {
		t.Fatalf("invocations = %d, want 2", len(invocations))
	}

	expected := []string{"update", "docs"}
	if !reflect.DeepEqual(invocations[0], expected) {
		t.Fatalf("update args = %#v, want %#v", invocations[0], expected)
	}
}

func TestIndexWithContextAndEmbedRunsExpectedCommands(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"collection add": addOutputFixture,
			"context add":    "context added",
			"embed":          embedOutputFixture,
			"status": strings.ReplaceAll(statusOutputFixture,
				"Vectors:  0 embedded\n  Pending:  1 need embedding (run 'qmd embed')",
				"Vectors:  3 embedded\n  Pending:  0 need embedding (run 'qmd embed')",
			),
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))

	result, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		VaultPath:      "/tmp/vault",
		CollectionName: "docs",
		Context:        "Repository documentation",
		Embed:          true,
		ForceEmbed:     true,
	})
	if err != nil {
		t.Fatalf("Index returned error: %v", err)
	}
	if result.EmbedResult == nil || result.EmbedResult.ChunksEmbedded != 9 {
		t.Fatalf("embed result = %#v, want parsed embed summary", result.EmbedResult)
	}
	if !result.Status.HasVectorIndex {
		t.Fatalf("status = %#v, want vector index reported", result.Status)
	}

	invocations := readInvocationLog(t, logPath)
	if len(invocations) != 4 {
		t.Fatalf("invocations = %d, want 4", len(invocations))
	}

	expectedContext := []string{"context", "add", "qmd://docs/", "Repository documentation"}
	if !reflect.DeepEqual(invocations[1], expectedContext) {
		t.Fatalf("context args = %#v, want %#v", invocations[1], expectedContext)
	}

	expectedEmbed := []string{"embed", "-f"}
	if !reflect.DeepEqual(invocations[2], expectedEmbed) {
		t.Fatalf("embed args = %#v, want %#v", invocations[2], expectedEmbed)
	}
}

func TestIndexRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	client := NewClient(WithBinaryPath(writeFakeQMD(t, fakeQMDOptions{})))

	_, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		CollectionName: "docs",
	})
	if err == nil || !strings.Contains(err.Error(), "vault path is required") {
		t.Fatalf("add missing vault error = %v, want validation error", err)
	}

	_, err = client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		VaultPath:      "/tmp/vault",
		CollectionName: "docs",
		Embed:          false,
		ForceEmbed:     true,
	})
	if err == nil || !strings.Contains(err.Error(), "force embed requires embed=true") {
		t.Fatalf("force embed error = %v, want validation error", err)
	}

	_, err = client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperation("sync"),
		CollectionName: "docs",
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported operation") {
		t.Fatalf("invalid operation error = %v, want unsupported operation", err)
	}
}

func TestSearchHybridModeUsesQueryCommand(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"query": lexicalSearchJSONFixture,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))

	_, err := client.Search(context.Background(), SearchOptions{
		Query:      "authentication flow",
		Mode:       SearchModeHybrid,
		Collection: "docs",
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	invocations := readInvocationLog(t, logPath)
	if len(invocations) != 1 {
		t.Fatalf("invocations = %d, want 1", len(invocations))
	}

	expected := []string{"query", "--json", "-c", "docs", "authentication flow"}
	if !reflect.DeepEqual(invocations[0], expected) {
		t.Fatalf("hybrid args = %#v, want %#v", invocations[0], expected)
	}
}

func TestSearchAllOmitsLimitAndAcceptsModeAlias(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"search": lexicalSearchJSONFixture,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))

	_, err := client.Search(context.Background(), SearchOptions{
		Query: "all auth results",
		Mode:  SearchMode("lex"),
		All:   true,
		Limit: 25,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	invocations := readInvocationLog(t, logPath)
	expected := []string{"search", "--json", "--all", "all auth results"}
	if !reflect.DeepEqual(invocations[0], expected) {
		t.Fatalf("all-mode args = %#v, want %#v", invocations[0], expected)
	}
}

func TestSearchPassesLimitMinScoreAndFullFlags(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "args.log")
	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		LogPath: logPath,
		StdoutByCommand: map[string]string{
			"vsearch": lexicalSearchJSONFixture,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))
	minScore := 0.3

	_, err := client.Search(context.Background(), SearchOptions{
		Query:      "semantic auth",
		Mode:       SearchModeVector,
		Limit:      7,
		MinScore:   &minScore,
		Full:       true,
		Collection: "docs",
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	invocations := readInvocationLog(t, logPath)
	if len(invocations) != 1 {
		t.Fatalf("invocations = %d, want 1", len(invocations))
	}

	expected := []string{"vsearch", "--json", "-n", "7", "--min-score", "0.3", "--full", "-c", "docs", "semantic auth"}
	if !reflect.DeepEqual(invocations[0], expected) {
		t.Fatalf("vector args = %#v, want %#v", invocations[0], expected)
	}
}

func TestSearchParsesJSONAndNormalizesResults(t *testing.T) {
	t.Parallel()

	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		StdoutByCommand: map[string]string{
			"query": hybridSearchJSONFixture,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))
	minScore := 0.5

	results, err := client.Search(context.Background(), SearchOptions{
		Query:    "auth",
		Mode:     SearchModeHybrid,
		MinScore: &minScore,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	expected := []SearchResult{
		{
			DocID:   "#abc123",
			Path:    "docs/auth.md",
			Title:   "Auth",
			Snippet: "Best chunk preview",
			Score:   0.82,
		},
	}
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Search results = %#v, want %#v", results, expected)
	}
}

func TestSearchFullUsesBodyWhenPresent(t *testing.T) {
	t.Parallel()

	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		StdoutByCommand: map[string]string{
			"search": `[{"docid":"#body","score":0.4,"file":"qmd://docs/file.md","title":"Doc","body":"# Full body"}]`,
		},
	})

	client := NewClient(WithBinaryPath(binaryPath))

	results, err := client.Search(context.Background(), SearchOptions{
		Query: "Doc",
		Mode:  SearchModeLexical,
		Full:  true,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) != 1 || results[0].Snippet != "# Full body" {
		t.Fatalf("full results = %#v, want snippet from body", results)
	}
}

func TestSearchContextCancellationStopsRunningCommand(t *testing.T) {
	t.Parallel()

	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		ScriptBody: `
sleep 5
printf '[]\n'
`,
	})

	client := NewClient(WithBinaryPath(binaryPath))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Search(ctx, SearchOptions{
		Query: "slow",
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Search error = %v, want context deadline exceeded", err)
	}
}

func TestSearchFailureIncludesStderrDiagnostics(t *testing.T) {
	t.Parallel()

	binaryPath := writeFakeQMD(t, fakeQMDOptions{
		ScriptBody: `
printf 'backend exploded\n' >&2
exit 4
`,
	})

	client := NewClient(WithBinaryPath(binaryPath))

	_, err := client.Search(context.Background(), SearchOptions{
		Query: "broken",
	})
	if err == nil {
		t.Fatal("Search error = nil, want failure")
	}
	if !strings.Contains(err.Error(), "backend exploded") {
		t.Fatalf("Search error = %v, want stderr diagnostics", err)
	}
}

func TestParseUpdateResultParsesAddAndUpdateSummaries(t *testing.T) {
	t.Parallel()

	addResult, err := parseUpdateResult(addOutputFixture, IndexOperationAdd)
	if err != nil {
		t.Fatalf("parseUpdateResult(add) returned error: %v", err)
	}
	if addResult.Collections != 1 || addResult.Indexed != 1 || addResult.NeedsEmbedding != 1 {
		t.Fatalf("add result = %#v", addResult)
	}

	updateResult, err := parseUpdateResult(updateOutputFixture, IndexOperationUpdate)
	if err != nil {
		t.Fatalf("parseUpdateResult(update) returned error: %v", err)
	}
	if updateResult.Collections != 1 || updateResult.Unchanged != 1 {
		t.Fatalf("update result = %#v", updateResult)
	}
}

func TestParseEmbedResultParsesSuccessAndNoWork(t *testing.T) {
	t.Parallel()

	result, err := parseEmbedResult(embedOutputFixture)
	if err != nil {
		t.Fatalf("parseEmbedResult returned error: %v", err)
	}
	if result.DocsProcessed != 3 || result.ChunksEmbedded != 9 || result.Errors != 2 || result.DurationMs != 62500 {
		t.Fatalf("embed result = %#v", result)
	}

	emptyResult, err := parseEmbedResult("✓ All content hashes already have embeddings.")
	if err != nil {
		t.Fatalf("parseEmbedResult(no work) returned error: %v", err)
	}
	if emptyResult != (EmbedResult{}) {
		t.Fatalf("empty embed result = %#v, want zero value", emptyResult)
	}
}

func TestParseIndexStatusParsesCollectionsAndHealth(t *testing.T) {
	t.Parallel()

	status, err := parseIndexStatus(statusOutputFixture)
	if err != nil {
		t.Fatalf("parseIndexStatus returned error: %v", err)
	}

	if status.TotalDocuments != 3 || status.NeedsEmbedding != 1 || status.HasVectorIndex {
		t.Fatalf("status = %#v", status)
	}
	if len(status.Collections) != 2 {
		t.Fatalf("collections = %#v, want 2 entries", status.Collections)
	}
	if status.Collections[1].Name != "notes" || status.Collections[1].Documents != 1 {
		t.Fatalf("second collection = %#v", status.Collections[1])
	}
}

func TestParseIndexStatusAcceptsEmptyIndex(t *testing.T) {
	t.Parallel()

	status, err := parseIndexStatus(emptyStatusOutputFixture)
	if err != nil {
		t.Fatalf("parseIndexStatus(empty) returned error: %v", err)
	}

	if status.TotalDocuments != 0 || status.NeedsEmbedding != 0 || status.HasVectorIndex {
		t.Fatalf("status = %#v, want empty zero-value summary", status)
	}
	if len(status.Collections) != 0 {
		t.Fatalf("collections = %#v, want no collections", status.Collections)
	}
}

func TestParseHumanDurationMillisecondsParsesMultipleUnits(t *testing.T) {
	t.Parallel()

	milliseconds, err := parseHumanDurationMilliseconds("1h 2m 3.5s 10ms")
	if err != nil {
		t.Fatalf("parseHumanDurationMilliseconds returned error: %v", err)
	}

	if milliseconds != 3723510 {
		t.Fatalf("milliseconds = %d, want 3723510", milliseconds)
	}
}

func TestNormalizeSearchModeRejectsUnsupportedMode(t *testing.T) {
	t.Parallel()

	if _, _, err := normalizeSearchMode(SearchMode("bm25")); err == nil {
		t.Fatal("normalizeSearchMode error = nil, want unsupported mode")
	}
}

type fakeQMDOptions struct {
	LogPath         string
	ScriptBody      string
	StdoutByCommand map[string]string
}

func writeFakeQMD(t *testing.T, options fakeQMDOptions) string {
	t.Helper()

	scriptPath := filepath.Join(t.TempDir(), "qmd")
	var builder strings.Builder
	builder.WriteString("#!/bin/sh\n")
	builder.WriteString("set -eu\n")
	if options.LogPath != "" {
		builder.WriteString("if [ -n \"${QMD_LOG:-}\" ]; then :; fi\n")
		builder.WriteString("LOG_PATH=" + shellQuote(options.LogPath) + "\n")
		builder.WriteString("for arg in \"$@\"; do printf '%s\\n' \"$arg\" >> \"$LOG_PATH\"; done\n")
		builder.WriteString("printf '%s\\n' '---' >> \"$LOG_PATH\"\n")
	}

	if options.ScriptBody != "" {
		builder.WriteString(options.ScriptBody)
	} else {
		builder.WriteString("cmd=$1\n")
		builder.WriteString("sub=\"\"\n")
		builder.WriteString("if [ \"$cmd\" = \"--index\" ]; then shift 2; cmd=$1; fi\n")
		builder.WriteString("if [ \"$cmd\" = \"collection\" ] || [ \"$cmd\" = \"context\" ]; then sub=$2; fi\n")
		builder.WriteString("key=$cmd\n")
		builder.WriteString("if [ -n \"$sub\" ]; then key=\"$cmd $sub\"; fi\n")
		builder.WriteString("case \"$key\" in\n")
		for key, stdout := range options.StdoutByCommand {
			builder.WriteString("  " + shellQuote(key) + ")\n")
			builder.WriteString("    cat <<'EOF'\n")
			builder.WriteString(stdout)
			if !strings.HasSuffix(stdout, "\n") {
				builder.WriteString("\n")
			}
			builder.WriteString("EOF\n")
			builder.WriteString("    ;;\n")
		}
		builder.WriteString("  *)\n")
		builder.WriteString("    printf 'unexpected command: %s\\n' \"$key\" >&2\n")
		builder.WriteString("    exit 9\n")
		builder.WriteString("    ;;\n")
		builder.WriteString("esac\n")
	}

	if err := os.WriteFile(scriptPath, []byte(builder.String()), 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", scriptPath, err)
	}
	return scriptPath
}

func readInvocationLog(t *testing.T, path string) [][]string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	invocations := make([][]string, 0, 4)
	current := make([]string, 0, 8)

	flush := func() {
		if len(current) == 0 {
			return
		}
		invocations = append(invocations, append([]string(nil), current...))
		current = current[:0]
	}

	for _, line := range lines {
		if line == "---" {
			flush()
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		current = append(current, line)
	}
	flush()

	return invocations
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

const addOutputFixture = `Creating collection 'docs'...
Collection: /tmp/vault (**/*.md)

Indexed: 1 new, 0 updated, 0 unchanged, 0 removed

Run 'qmd embed' to update embeddings (1 unique hashes need vectors)
✓ Collection 'docs' created successfully`

const updateOutputFixture = `Updating 1 collection(s)...

[1/1] docs (**/*.md)
Collection: /tmp/vault (**/*.md)

Indexed: 0 new, 0 updated, 1 unchanged, 0 removed

✓ All collections updated.

Run 'qmd embed' to update embeddings (1 unique hashes need vectors)`

const embedOutputFixture = `✓ Done! Embedded 9 chunks from 3 documents in 1m 2.5s
⚠ 2 chunks failed`

const statusOutputFixture = `QMD Status

Index: /tmp/index.sqlite
Size:  92.0 KB

Documents
  Total:    3 files indexed
  Vectors:  0 embedded
  Pending:  1 need embedding (run 'qmd embed')
  Updated:  7m ago

Collections
  docs (qmd://docs/)
    Pattern:  **/*.md
    Files:    2 (updated 7m ago)
  notes (qmd://notes/)
    Pattern:  **/*.md
    Files:    1 (updated 3m ago)
`

const emptyStatusOutputFixture = `QMD Status

Index: /tmp/index.sqlite
Size:  4.0 KB

Documents
  Total:    0 files indexed
  Vectors:  0 embedded

No collections. Run 'qmd collection add .' to index markdown files.
`

const lexicalSearchJSONFixture = `[{"docid":"#body","score":0.4,"file":"qmd://docs/file.md","title":"Doc","snippet":"Preview"}]`

const hybridSearchJSONFixture = `[
  {
    "docid": "#abc123",
    "score": 0.82,
    "displayPath": "docs/auth.md",
    "file": "qmd://docs/auth.md",
    "title": "Auth",
    "bestChunk": "Best chunk preview",
    "body": "# Auth\n\nFull body"
  },
  {
    "docid": "#low",
    "score": 0.3,
    "file": "qmd://docs/low.md",
    "title": "Low",
    "snippet": "Low score snippet"
  }
]`
