package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/fakes"
	kfind "github.com/compozy/kb/internal/find"
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
