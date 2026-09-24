package qmd

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/corpus"
)

// candidatesStatus is a `qmd status` output listing collection demo with
// its index at indexPath.
func candidatesStatus(indexPath string) string {
	return "QMD Status\n\nIndex: " + indexPath + "\nSize:  92.0 KB\n\nDocuments\n  Total:    3 files indexed\n  Vectors:  3 embedded\n  Pending:  0 need embedding\n\nCollections\n  demo (qmd://demo/)\n    Pattern:  **/*.md\n    Files:    3 (updated 1m ago)\n"
}

// candidatesTopic writes a topic with one markdown file modified at mdTime
// and an index file modified at indexTime; it returns the topic root and the
// index path.
func candidatesTopic(t *testing.T, mdTime, indexTime time.Time) (string, string) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo")
	md := filepath.Join(root, "raw", "a.md")
	hidden := filepath.Join(root, ".decisions", "state.md")
	indexPath := filepath.Join(t.TempDir(), "index.sqlite")
	for _, file := range []string{md, hidden, indexPath} {
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for file, when := range map[string]time.Time{md: mdTime, indexPath: indexTime, hidden: mdTime.Add(time.Hour)} {
		if err := os.Chtimes(file, when, when); err != nil {
			t.Fatal(err)
		}
	}
	return root, indexPath
}

func TestOpenCandidatesFreshness(t *testing.T) {
	t.Parallel()
	base := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		missing   bool
		status    func(indexPath string) string
		mdTime    time.Time
		indexTime time.Time
		wantOn    bool
		wantNote  string
	}{
		{name: "fresh index newer than topic files (dot dirs ignored)", status: candidatesStatus, mdTime: base, indexTime: base.Add(time.Minute), wantOn: true},
		{name: "same mtime is fresh", status: candidatesStatus, mdTime: base, indexTime: base, wantOn: true},
		{name: "index older than a topic file", status: candidatesStatus, mdTime: base, indexTime: base.Add(-time.Minute), wantNote: "older than the topic files"},
		{name: "collection missing", status: func(indexPath string) string {
			return strings.ReplaceAll(candidatesStatus(indexPath), "demo", "other")
		}, mdTime: base, indexTime: base.Add(time.Minute), wantNote: `no qmd collection "demo"`},
		{name: "index path missing from status", status: func(indexPath string) string {
			return strings.ReplaceAll(candidatesStatus(indexPath), "Index: "+indexPath+"\n", "")
		}, mdTime: base, indexTime: base.Add(time.Minute), wantNote: "did not report the index path"},
		{name: "qmd not installed", missing: true, mdTime: base, indexTime: base, wantNote: "not installed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, indexPath := candidatesTopic(t, tt.mdTime, tt.indexTime)
			var client *QMDClient
			if tt.missing {
				client = NewClient(WithBinaryPath(filepath.Join(t.TempDir(), "missing-qmd")))
			} else {
				binary := writeFakeQMD(t, fakeQMDOptions{StdoutByCommand: map[string]string{"status": tt.status(indexPath)}})
				client = newFakeQMDClient(binary)
			}
			candidates, note := client.OpenCandidates(context.Background(), "demo", root)
			if (candidates != nil) != tt.wantOn {
				t.Fatalf("candidates on = %v, want %v (note %q)", candidates != nil, tt.wantOn, note)
			}
			if tt.wantOn && note != "" {
				t.Fatalf("note = %q, want none when on", note)
			}
			if !tt.wantOn && !strings.Contains(note, tt.wantNote) {
				t.Fatalf("note = %q, want it to contain %q", note, tt.wantNote)
			}
		})
	}
}

const candidateHits = `[
  {"docid":"#1","score":1,"file":"qmd://demo/raw/articles/self.md","title":"Self"},
  {"docid":"#2","score":0.5,"file":"qmd://demo/log.md","title":"log"},
  {"docid":"#3","score":0.33,"file":"qmd://demo/wiki/concepts/Decision%20Engines.md?index=st","title":"3"},
  {"docid":"#4","score":0.25,"file":"qmd://demo/wiki/index/Concept Index.md","title":"Index"},
  {"docid":"#5","score":0.2,"file":"qmd://demo/wiki/concepts/RAG.md","title":"RAG"},
  {"docid":"#6","score":0.16,"file":"qmd://demo/wiki/concepts/RAG.md","title":"RAG again"},
  {"docid":"#7","score":0.14,"file":"qmd://demo/wiki/concepts/Judges.md","title":"Judges"}
]`

// openFakeCandidates returns candidates over a fake qmd answering status and
// query, and the invocation log path.
func openFakeCandidates(t *testing.T, queryOutput string) (*Candidates, string) {
	t.Helper()
	now := time.Now()
	root, indexPath := candidatesTopic(t, now.Add(-time.Hour), now)
	logPath := filepath.Join(t.TempDir(), "args.log")
	binary := writeFakeQMD(t, fakeQMDOptions{LogPath: logPath, StdoutByCommand: map[string]string{
		"status": candidatesStatus(indexPath),
		"query":  queryOutput,
	}})
	candidates, note := newFakeQMDClient(binary).OpenCandidates(context.Background(), "demo", root)
	if candidates == nil {
		t.Fatalf("candidates off: %s", note)
	}
	return candidates, logPath
}

func TestCandidatesNeighboursDropsSelfAndHubs(t *testing.T) {
	t.Parallel()
	candidates, logPath := openFakeCandidates(t, candidateHits)
	doc := &corpus.Document{
		Path:        "raw/articles/self.md",
		Title:       "Self",
		Frontmatter: map[string]any{"summary": "A summary\nover lines."},
		Body:        "# Heading\n\nBody text.\n",
	}

	got, err := candidates.Neighbours(context.Background(), doc, 2)
	if err != nil {
		t.Fatalf("Neighbours: %v", err)
	}
	want := []string{"wiki/concepts/Decision Engines.md", "wiki/concepts/RAG.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Neighbours = %#v, want %#v", got, want)
	}

	invocations := readInvocationLog(t, logPath)
	query := invocations[len(invocations)-1]
	wantArgs := []string{"query", "--json", "-n", "12", "--no-rerank", "-c", "demo", "vec: Self A summary over lines. # Heading Body text."}
	if !reflect.DeepEqual(query, wantArgs) {
		t.Fatalf("query args = %#v, want %#v", query, wantArgs)
	}
}

func TestCandidatesSearchReturnsDocumentPaths(t *testing.T) {
	t.Parallel()
	candidates, logPath := openFakeCandidates(t, candidateHits)

	got, err := candidates.Search(context.Background(), "how do judges\ndecide?", 30)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	want := []string{"raw/articles/self.md", "log.md", "wiki/concepts/Decision Engines.md", "wiki/index/Concept Index.md", "wiki/concepts/RAG.md", "wiki/concepts/Judges.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Search = %#v, want %#v", got, want)
	}
	invocations := readInvocationLog(t, logPath)
	if last := invocations[len(invocations)-1]; last[len(last)-1] != "vec: how do judges decide?" {
		t.Fatalf("query text = %q", last[len(last)-1])
	}
}

func TestCandidatesQueryErrors(t *testing.T) {
	t.Parallel()
	candidates, _ := openFakeCandidates(t, "not json")

	if _, err := candidates.Search(context.Background(), "question", 5); err == nil || !strings.Contains(err.Error(), "parse JSON output") {
		t.Fatalf("Search error = %v, want a parse error", err)
	}
	if _, err := candidates.Search(context.Background(), "   ", 5); err == nil {
		t.Fatal("Search accepted an empty query")
	}
	if _, err := candidates.Search(context.Background(), "question", 0); err == nil {
		t.Fatal("Search accepted a zero limit")
	}
}

func TestSingleLineCapsLength(t *testing.T) {
	t.Parallel()
	long := strings.Repeat("word ", maxVectorQueryRunes)
	if got := []rune(singleLine(long)); len(got) != maxVectorQueryRunes {
		t.Fatalf("singleLine length = %d, want %d", len(got), maxVectorQueryRunes)
	}
	if got := singleLine(" a\n\tb  c "); got != "a b c" {
		t.Fatalf("singleLine = %q", got)
	}
}

func TestCandidatesDropQuarantinedHits(t *testing.T) {
	t.Parallel()
	candidates, _ := openFakeCandidates(t, quarantineHits)

	got, err := candidates.Search(context.Background(), "junk", 30)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if want := []string{"raw/articles/kept.md", "wiki/concepts/Quarantine.md"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Search = %#v, want %#v", got, want)
	}
}
