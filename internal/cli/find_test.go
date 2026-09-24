package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/fakes"
	kfind "github.com/compozy/kb/internal/find"
	"github.com/compozy/kb/internal/qmd"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

// newLinkFindVault scaffolds topic "demo" with the given documents.
func newLinkFindVault(t *testing.T, docs map[string]string) string {
	t.Helper()
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	for rel, content := range docs {
		path := filepath.Join(info.RootPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return vault
}

// useLinkFindFakeSession points openSession at a fake decisions server.
func useLinkFindFakeSession(t *testing.T, decide fakes.DecideFunc) *fakes.OpenRouter {
	t.Helper()
	fake := fakes.NewOpenRouter(decide, nil)
	t.Cleanup(fake.Close)
	useQMDCandidates(t, "")
	original := openSession
	t.Cleanup(func() { openSession = original })
	openSession = func(cmd *cobra.Command, command, topicSlug string, flags session.Flags) (*session.Session, error) {
		vaultPath, err := resolveCommandVaultPath(cmd, sessionGetwd, command)
		if err != nil {
			return nil, err
		}
		cfg := config.Default()
		cfg.OpenRouter.APIKey = "test-key"
		cfg.OpenRouter.APIURL = fake.URL
		return session.Open(session.Options{Config: cfg, VaultPath: vaultPath, Topic: topicSlug, Command: command, Flags: flags})
	}
	return fake
}

// useQMDCandidates points the optional qmd candidates at a fake qmd
// executable answering `status` (a fresh collection "demo") and `query`
// with queryJSON; an empty queryJSON turns qmd candidates off. It returns
// the fake's argument log path.
func useQMDCandidates(t *testing.T, queryJSON string) string {
	t.Helper()
	original := openQMDCandidates
	t.Cleanup(func() { openQMDCandidates = original })
	if queryJSON == "" {
		openQMDCandidates = func(context.Context, string, string) (*qmd.Candidates, string) { return nil, "" }
		return ""
	}
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.sqlite")
	if err := os.WriteFile(indexPath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	future := time.Now().Add(time.Hour)
	if err := os.Chtimes(indexPath, future, future); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(dir, "args.log")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + logPath + "'\ncase \"$1\" in\n" +
		"status) printf 'QMD Status\\n\\nIndex: " + indexPath + "\\n\\nDocuments\\n  Total:    1 files indexed\\n\\nCollections\\n  demo (qmd://demo/)\\n    Pattern:  **/*.md\\n    Files:    1 (updated 1m ago)\\n' ;;\n" +
		"query) cat <<'EOF'\n" + queryJSON + "\nEOF\n;;\n*) exit 9 ;;\nesac\n"
	binary := filepath.Join(dir, "qmd")
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	openQMDCandidates = func(ctx context.Context, collection, topicRoot string) (*qmd.Candidates, string) {
		return qmd.NewClient(qmd.WithBinaryPath(binary)).OpenCandidates(ctx, collection, topicRoot)
	}
	return logPath
}

func TestFindCommandAddsQMDVectorCandidates(t *testing.T) {
	vault := newLinkFindVault(t, map[string]string{
		"raw/articles/guide.md":  "---\ntitle: Tuning HNSW\n---\nTune HNSW.\n",
		"raw/articles/vector.md": "---\ntitle: Graph index parameters\n---\nEfConstruction and M trade recall for memory.\n",
	})
	useLinkFindFakeSession(t, func(_ fakes.Call, q fakes.Question) any {
		if strings.HasPrefix(q.ID, "answers_") {
			return fakes.Noul(0.9)
		}
		return nil
	})
	logPath := useQMDCandidates(t, `[{"docid":"#1","score":1,"file":"qmd://demo/raw/articles/vector.md","title":"x"}]`)

	stdout, _, err := runLinkFindRoot(t, "find", "demo", "tune HNSW", "--vault", vault, "--json")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	var result kfind.Result
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode: %v\n%s", err, stdout)
	}
	if result.Candidates != 2 {
		t.Fatalf("candidates = %d, want 2 (BM25 hit + qmd vector hit): %+v", result.Candidates, result)
	}
	log, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(log), "query --json -n 90 --no-rerank -c demo vec: tune HNSW") {
		t.Fatalf("qmd invocations:\n%s", log)
	}

	if err := os.Remove(logPath); err != nil {
		t.Fatal(err)
	}
	stdout, _, err = runLinkFindRoot(t, "find", "demo", "tune HNSW", "--vault", vault, "--json", "--no-qmd")
	if err != nil {
		t.Fatalf("find --no-qmd: %v", err)
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode: %v\n%s", err, stdout)
	}
	if result.Candidates != 1 {
		t.Fatalf("--no-qmd candidates = %d, want 1", result.Candidates)
	}
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("--no-qmd still ran qmd (log stat err %v)", err)
	}
}

func runLinkFindRoot(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	command := newRootCommand()
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)
	err := command.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

func TestFindFacetsNeedsNoDecisionModel(t *testing.T) {
	vault := newLinkFindVault(t, map[string]string{
		"raw/articles/a.md": "---\ntitle: A\ngenre: paper\nrelevance: core\ndepth: 2\nconcepts:\n  - \"[[RAG]]\"\n---\nBody.\n",
		"raw/articles/b.md": "---\ntitle: B\ngenre: paper\n---\nBody.\n",
	})
	original := openSession
	t.Cleanup(func() { openSession = original })
	openSession = func(*cobra.Command, string, string, session.Flags) (*session.Session, error) {
		t.Fatal("--facets opened a decision session")
		return nil, nil
	}

	stdout, _, err := runLinkFindRoot(t, "find", "--facets", "demo", "--vault", vault, "--json")
	if err != nil {
		t.Fatalf("find --facets: %v", err)
	}
	var report kfind.FacetReport
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("decode %q: %v", stdout, err)
	}
	if report.Documents != 2 || report.Genre[0] != (kfind.Count{Value: "paper", Count: 2}) || report.Concepts[0] != (kfind.Count{Value: "RAG", Count: 1}) {
		t.Fatalf("report = %+v", report)
	}

	table, _, err := runLinkFindRoot(t, "find", "--facets", "demo", "--vault", vault)
	if err != nil {
		t.Fatalf("find --facets table: %v", err)
	}
	for _, want := range []string{"documents: 2 (2 sources, 0 articles)", "genre:", "relevance:", "depth:", "concepts:"} {
		if !strings.Contains(table, want) {
			t.Fatalf("table output lacks %q:\n%s", want, table)
		}
	}
}

func TestFindCommandPrintsRankedResultsAndExplain(t *testing.T) {
	vault := newLinkFindVault(t, map[string]string{
		"raw/articles/guide.md": "---\ntitle: Tuning HNSW\ngenre: tutorial_or_guide\n---\nTune HNSW.\n",
		"raw/articles/spam.md":  "---\ntitle: HNSW spam\ntriage: quarantined\n---\nHNSW.\n",
	})
	useLinkFindFakeSession(t, func(_ fakes.Call, q fakes.Question) any {
		if strings.HasPrefix(q.ID, "answers_") {
			return fakes.Noul(0.9)
		}
		return nil
	})

	stdout, stderr, err := runLinkFindRoot(t, "find", "demo", "tune HNSW", "--vault", vault, "--explain")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	for _, want := range []string{"raw/articles/guide.md", "Tuning HNSW", "0.90", "dropped (1):", "raw/articles/spam.md", "quarantined"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("output lacks %q:\n%s", want, stdout)
		}
	}
	if !strings.Contains(stderr, "find: 1 candidates judged, 1 kept") || !strings.Contains(stderr, "decisions:") {
		t.Fatalf("stderr lacks the run summary:\n%s", stderr)
	}

	jsonOut, _, err := runLinkFindRoot(t, "find", "demo", "tune HNSW", "--vault", vault, "--json")
	if err != nil {
		t.Fatalf("find --json: %v", err)
	}
	var result kfind.Result
	if err := json.Unmarshal([]byte(jsonOut), &result); err != nil {
		t.Fatalf("decode: %v\n%s", err, jsonOut)
	}
	if len(result.Hits) != 1 || result.Hits[0].Path != "raw/articles/guide.md" || result.Dropped != nil {
		t.Fatalf("result = %+v", result)
	}
}

func TestFindCommandArgs(t *testing.T) {
	if _, _, err := runLinkFindRoot(t, "find", "demo"); err == nil {
		t.Fatal("find without a question must fail")
	}
	if _, _, err := runLinkFindRoot(t, "find", "--facets", "demo", "extra"); err == nil {
		t.Fatal("find --facets with two args must fail")
	}
}
