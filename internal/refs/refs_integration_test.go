//go:build integration

package refs

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	klint "github.com/compozy/kb/internal/lint"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/topic"
)

// quarantineTopic scaffolds a real topic whose source raw/articles/x.md is
// cited by a Source Index line (path-form link), an article's `sources:`
// list and an article body.
func quarantineTopic(t *testing.T) (vault, root string, touched []string) {
	t.Helper()
	vault = t.TempDir()
	info, err := topic.New(vault, "research", "Research", "research")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	root = info.RootPath
	writeFile(t, filepath.Join(root, "raw/articles/x.md"), "---\ntitle: X\nsource_url: https://example.com/x\n---\n# X\n\nBody.\n")
	writeFile(t, filepath.Join(root, "raw/articles/y.md"), "---\ntitle: Y\n---\nBody.\n")
	indexPath := filepath.Join(root, "wiki/index/Source Index.md")
	index := readString(t, indexPath)
	writeFile(t, indexPath, index+"\n- [[raw/articles/x]] — X capture\n- [[raw/articles/y]] — Y capture\n")
	articlePath := filepath.Join(root, "wiki/concepts/Topic.md")
	writeFile(t, articlePath, "---\ntitle: Topic\ntype: wiki\nstage: compiled\nsources:\n  - \"[[raw/articles/x]]\"\n  - \"[[raw/articles/y]]\"\n---\n# Topic\n\nAs [[x]] shows, and [[y]] agrees.\n")
	return vault, root, []string{indexPath, articlePath}
}

func TestQuarantineRoundTripIntegration(t *testing.T) {
	vault, root, touched := quarantineTopic(t)
	originals := map[string]string{}
	for _, path := range touched {
		originals[path] = readString(t, path)
	}

	result, err := Quarantine(context.Background(), Options{VaultPath: vault, TopicRoot: root, Path: "raw/articles/x.md", Reason: "off_topic"})
	if err != nil {
		t.Fatalf("Quarantine: %v", err)
	}
	index := readString(t, touched[0])
	article := readString(t, touched[1])
	if strings.Contains(index, "[[raw/articles/x]] — X capture") || !strings.Contains(index, "[[raw/articles/y]] — Y capture") {
		t.Fatalf("index after quarantine:\n%s", index)
	}
	if strings.Contains(article, `- "[[raw/articles/x]]"`) || !strings.Contains(article, "As [[x]] shows") {
		t.Fatalf("article after quarantine:\n%s", article)
	}
	recorded := map[string]bool{}
	ledger, err := LoadLedger(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range ledger {
		if entry.QuarantineID == result.ID && entry.Op == OpEdit {
			recorded[entry.File+"#"+entry.LineOrKey] = true
		}
	}
	if !recorded["wiki/index/Source Index.md#index"] || !recorded["wiki/concepts/Topic.md#sources"] {
		t.Fatalf("ledger = %+v", ledger)
	}
	if len(result.BodyRefs) != 1 || result.BodyRefs[0].File != "wiki/concepts/Topic.md" || result.BodyRefs[0].Text != "[[x]]" {
		t.Fatalf("body refs = %+v", result.BodyRefs)
	}

	issues, err := klint.Lint(root)
	if err != nil {
		t.Fatalf("lint: %v", err)
	}
	for _, issue := range issues {
		if issue.Kind == models.LintIssueKindDeadLink && issue.FilePath == "wiki/concepts/Topic.md" {
			t.Fatalf("body link to the quarantined source reported as dead link: %+v", issue)
		}
	}

	restored, err := Restore(context.Background(), Options{VaultPath: vault, TopicRoot: root, ID: result.ID})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(restored.Manual) != 0 {
		t.Fatalf("manual = %+v", restored.Manual)
	}
	for _, path := range touched {
		if got := readString(t, path); got != originals[path] {
			t.Errorf("%s not byte-identical after restore", path)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "raw/articles/x.md")); err != nil {
		t.Fatalf("source not moved back: %v", err)
	}
}

func TestRestoreAfterEditingTouchedFileIntegration(t *testing.T) {
	vault, root, touched := quarantineTopic(t)
	originalIndex := readString(t, touched[0])

	result, err := Quarantine(context.Background(), Options{VaultPath: vault, TopicRoot: root, Path: "raw/articles/x.md", Reason: "thin"})
	if err != nil {
		t.Fatalf("Quarantine: %v", err)
	}
	edited := strings.Replace(readString(t, touched[1]), "# Topic", "# Topic (edited)", 1)
	writeFile(t, touched[1], edited)

	restored, err := Restore(context.Background(), Options{VaultPath: vault, TopicRoot: root, ID: result.ID})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(restored.Manual) != 1 || restored.Manual[0].File != "wiki/concepts/Topic.md" || restored.Manual[0].LineOrKey != "sources" {
		t.Fatalf("manual = %+v", restored.Manual)
	}
	if got := readString(t, touched[1]); got != edited {
		t.Fatalf("edited article overwritten:\n%s", got)
	}
	if got := readString(t, touched[0]); got != originalIndex {
		t.Fatalf("index not restored:\n%s", got)
	}
}
