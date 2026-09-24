//go:build integration

package link

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

const (
	ragArticle = "wiki/concepts/Retrieval Augmented Generation.md"
	sourceA    = "raw/articles/source-a.md"
)

type linkVault struct {
	vault string
	root  string
}

func newLinkVault(t *testing.T, mode string) linkVault {
	t.Helper()
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	v := linkVault{vault: vault, root: info.RootPath}
	v.write(t, ragArticle, "---\ntitle: Retrieval Augmented Generation\naliases:\n  - RAG\ncriterion: Grounding a model's generation in retrieved documents.\n---\n# Retrieval Augmented Generation\n\nRetrieval before generation.\n")
	v.write(t, sourceA, "---\ntitle: Source A\nsource_kind: article\n---\n# Source A\n\n```\nRAG in code\n```\n\nWe run RAG in production. Later, RAG again.\n")
	if mode != "" {
		path := filepath.Join(v.root, "topic.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, []byte("decisions:\n  mode: "+mode+"\n")...), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return v
}

func (v linkVault) write(t *testing.T, rel, content string) {
	t.Helper()
	path := filepath.Join(v.root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (v linkVault) read(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(v.root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func (v linkVault) open(t *testing.T, fake *fakes.OpenRouter) *session.Session {
	t.Helper()
	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = fake.URL
	s, err := session.Open(session.Options{Config: cfg, VaultPath: v.vault, Topic: "demo", Command: "kb link"})
	if err != nil {
		t.Fatalf("session.Open: %v", err)
	}
	return s
}

// linkAll answers every link question positively.
func linkAll(_ fakes.Call, q fakes.Question) any {
	switch {
	case strings.HasPrefix(q.ID, "should_link_"), strings.HasPrefix(q.ID, "mention_sense_"):
		return fakes.Noul(0.95)
	case strings.HasPrefix(q.ID, "affects_"):
		return fakes.Noul(0.9)
	case strings.HasPrefix(q.ID, "relation_"):
		return fakes.Pick("extends", 0.9)
	}
	return nil
}

func newFake(t *testing.T, decide fakes.DecideFunc) *fakes.OpenRouter {
	t.Helper()
	fake := fakes.NewOpenRouter(decide, nil)
	t.Cleanup(fake.Close)
	return fake
}

func frontmatterList(t *testing.T, content, key string) []string {
	t.Helper()
	values, _, err := frontmatter.Parse(content)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return frontmatter.GetStringSlice(values, key)
}

func TestLinkApplyInsertsBodyLinkAndIsIncremental(t *testing.T) {
	v := newLinkVault(t, session.ModeApply)
	fake := newFake(t, linkAll)

	report, err := Run(context.Background(), v.open(t, fake), Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	content := v.read(t, sourceA)
	if got := frontmatterList(t, content, "extends"); len(got) != 1 || got[0] != "[[Retrieval Augmented Generation]]" {
		t.Fatalf("extends = %v\n%s", got, content)
	}
	if got := frontmatterList(t, content, "affects"); len(got) != 1 {
		t.Fatalf("affects = %v", got)
	}
	if !strings.Contains(content, "We run [[Retrieval Augmented Generation|RAG]] in production. Later, RAG again.") {
		t.Fatalf("body link not inserted at the first qualifying mention:\n%s", content)
	}
	if !strings.Contains(content, "```\nRAG in code\n```") {
		t.Fatalf("code block was touched:\n%s", content)
	}
	if report.Inserted != 1 || report.Applied["extends"] != 1 || report.Applied["affects"] != 1 {
		t.Fatalf("report = %+v", report)
	}
	rows := readJSONL(t, filepath.Join(v.root, ".decisions", InsertedLinksFile))
	if len(rows) != 1 || rows[0]["subject"] != sourceA || rows[0]["target"] != ragArticle || rows[0]["text"] != "RAG" || rows[0]["mode"] != "apply" {
		t.Fatalf("inserted-links rows = %v", rows)
	}
	for _, call := range fake.Calls() {
		if _, ok := call.Questions["affects_c1"]; ok && !strings.Contains(call.StateString(), `"criterion":"Grounding`) {
			t.Fatalf("candidate state lacks criterion: %s", call.StateString())
		}
	}

	calls := len(fake.Calls())
	again, err := Run(context.Background(), v.open(t, fake), Options{})
	if err != nil {
		t.Fatalf("second Run: %v", err)
	}
	if got := len(fake.Calls()); got != calls {
		t.Fatalf("incremental run made %d new calls", got-calls)
	}
	if again.Skipped[SkipUnchanged] == 0 || again.Judged != 0 {
		t.Fatalf("second report = %+v", again)
	}
	if v.read(t, sourceA) != content {
		t.Fatal("second run changed the document")
	}
}

func TestLinkShadowProposesBodyLinkButWritesRelation(t *testing.T) {
	v := newLinkVault(t, "")
	fake := newFake(t, linkAll)
	s := v.open(t, fake)

	report, err := Run(context.Background(), s, Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	content := v.read(t, sourceA)
	if strings.Contains(content, "[[Retrieval Augmented Generation|") {
		t.Fatalf("shadow mode inserted a body link:\n%s", content)
	}
	if got := frontmatterList(t, content, "extends"); len(got) != 1 {
		t.Fatalf("extends = %v (relation must be written in shadow too)", got)
	}
	if report.Proposals != 1 || report.Inserted != 0 {
		t.Fatalf("report = %+v", report)
	}
	items, err := review.Open(v.root, nil).Items(review.QueueLink, review.StatusPending)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Action["insert"] != true || items[0].Action["mention_text"] != "RAG" || items[0].Action["target"] != "Retrieval Augmented Generation" {
		t.Fatalf("queue = %+v", items)
	}
	if _, err := os.Stat(filepath.Join(v.root, ".decisions", InsertedLinksFile)); !os.IsNotExist(err) {
		t.Fatalf("shadow mode wrote inserted-links: %v", err)
	}
}

func TestLinkDryRunWritesNothing(t *testing.T) {
	v := newLinkVault(t, session.ModeApply)
	fake := newFake(t, linkAll)
	before := v.read(t, sourceA)

	report, err := Run(context.Background(), v.open(t, fake), Options{DryRun: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if v.read(t, sourceA) != before {
		t.Fatal("dry run changed the document")
	}
	for _, name := range []string{"state.jsonl", review.QueueFile, InsertedLinksFile} {
		if _, err := os.Stat(filepath.Join(v.root, ".decisions", name)); !os.IsNotExist(err) {
			t.Fatalf("dry run wrote %s", name)
		}
	}
	joined := strings.Join(report.Lines(), "\n")
	if !strings.Contains(joined, "would add extends [[Retrieval Augmented Generation]] to "+sourceA) || !strings.Contains(joined, "would insert") {
		t.Fatalf("dry-run lines:\n%s", joined)
	}
}

func TestLinkRespectsLocked(t *testing.T) {
	v := newLinkVault(t, session.ModeApply)
	locked := "---\ntitle: Source A\nlocked: true\n---\nWe run RAG in production.\n"
	v.write(t, sourceA, locked)
	fake := newFake(t, linkAll)

	report, err := Run(context.Background(), v.open(t, fake), Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if v.read(t, sourceA) != locked {
		t.Fatal("locked document was written")
	}
	if report.Skipped[SkipLocked] != 1 {
		t.Fatalf("report = %+v", report)
	}
	if fake.QuestionCount("affects_") != 0 {
		t.Fatal("a locked source was judged")
	}
}

func TestLinkRefusesChangedBody(t *testing.T) {
	v := newLinkVault(t, session.ModeApply)
	var once sync.Once
	edited := ""
	fake := newFake(t, func(call fakes.Call, q fakes.Question) any {
		if strings.Contains(call.StateString(), "Source A") {
			once.Do(func() {
				edited = v.read(t, sourceA) + "\nA human edit.\n"
				v.write(t, sourceA, edited)
			})
		}
		return linkAll(call, q)
	})

	report, err := Run(context.Background(), v.open(t, fake), Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got := v.read(t, sourceA); got != edited {
		t.Fatalf("changed document was overwritten:\n%s", got)
	}
	if report.Skipped[SkipChanged] != 1 || report.Inserted != 0 {
		t.Fatalf("report = %+v", report)
	}
}

func TestReversePassJudgesMentioningDocuments(t *testing.T) {
	v := newLinkVault(t, "")
	v.write(t, "wiki/concepts/Vector Database.md", "---\ntitle: Vector Database\ncriterion: Databases that index embeddings.\n---\nStores vectors.\n")
	v.write(t, "raw/articles/source-b.md", "---\ntitle: Source B\n---\nWe picked a vector database for search.\n")
	fake := newFake(t, linkAll)

	report, err := Reverse(context.Background(), v.open(t, fake), "wiki/concepts/Vector Database.md")
	if err != nil {
		t.Fatalf("Reverse: %v", err)
	}
	if report.Documents != 1 || report.Judged != 1 {
		t.Fatalf("report = %+v", report)
	}
	for _, call := range fake.Calls() {
		if _, ok := call.Questions["should_link_c2"]; ok {
			t.Fatal("reverse pass asked about more than one candidate")
		}
	}
	if got := frontmatterList(t, v.read(t, "raw/articles/source-b.md"), "extends"); len(got) != 1 || got[0] != "[[Vector Database]]" {
		t.Fatalf("extends = %v", got)
	}
	if strings.Contains(v.read(t, sourceA), "Vector Database") {
		t.Fatal("a document without a mention was linked")
	}
}

func readJSONL(t *testing.T, path string) []map[string]any {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()
	rows := make([]map[string]any, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := map[string]any{}
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		rows = append(rows, row)
	}
	return rows
}
