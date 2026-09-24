package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/fakes"
)

func TestLinkCommandDryRunAndApply(t *testing.T) {
	source := "---\ntitle: Source\n---\nWe run RAG in production.\n"
	vault := newLinkFindVault(t, map[string]string{
		"wiki/concepts/RAG.md": "---\ntitle: RAG\ncriterion: Retrieval before generation.\n---\nRAG.\n",
		"raw/articles/s.md":    source,
	})
	useLinkFindFakeSession(t, func(_ fakes.Call, q fakes.Question) any {
		switch {
		case strings.HasPrefix(q.ID, "should_link_"), strings.HasPrefix(q.ID, "mention_sense_"):
			return fakes.Noul(0.95)
		case strings.HasPrefix(q.ID, "relation_"):
			return fakes.Pick("related", 0.9)
		}
		return nil
	})
	sourcePath := filepath.Join(vault, "demo", "raw", "articles", "s.md")

	stdout, _, err := runLinkFindRoot(t, "link", "demo", "--vault", vault, "--dry-run")
	if err != nil {
		t.Fatalf("link --dry-run: %v", err)
	}
	if !strings.Contains(stdout, "link (dry run): 2 documents, 1 judged") || !strings.Contains(stdout, "would add related [[RAG]] to raw/articles/s.md") {
		t.Fatalf("dry-run output:\n%s", stdout)
	}
	if data, _ := os.ReadFile(sourcePath); string(data) != source {
		t.Fatalf("dry run wrote the source:\n%s", data)
	}

	stdout, stderr, err := runLinkFindRoot(t, "link", "demo", "--vault", vault, "--decisions", "apply")
	if err != nil {
		t.Fatalf("link: %v", err)
	}
	if !strings.Contains(stdout, "relations: related 1") || !strings.Contains(stdout, "body links: 1 inserted") {
		t.Fatalf("output:\n%s", stdout)
	}
	if !strings.Contains(stderr, "decisions:") {
		t.Fatalf("stderr lacks the run summary:\n%s", stderr)
	}
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "[[RAG|RAG]]") || !strings.Contains(string(data), "related:") {
		t.Fatalf("source after link:\n%s", data)
	}
}
