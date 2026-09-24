package ingest

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/gate"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

var runNow = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

func openRunSession(t *testing.T) (*session.Session, *fakes.OpenRouter) {
	t.Helper()
	vault := t.TempDir()
	if _, err := topic.New(vault, "demo", "Demo", "demo"); err != nil {
		t.Fatal(err)
	}
	fake := fakes.NewOpenRouter(nil, nil)
	t.Cleanup(fake.Close)
	cfg := config.Default()
	cfg.OpenRouter.APIKey, cfg.OpenRouter.APIURL = "test-key", fake.URL
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb ingest", Now: func() time.Time { return runNow }})
	if err != nil {
		t.Fatal(err)
	}
	return s, fake
}

func words(n int) string {
	return strings.Repeat("decision engines judge typed questions. ", n/5) + "\n"
}

func TestRunWritesOwnedKeysThroughTheWriter(t *testing.T) {
	t.Parallel()
	s, fake := openRunSession(t)
	run := NewRun(s, RunOptions{Command: "url", Batch: "url-2026-09-24-abc123", Query: "reading list"})
	result, err := run.Ingest(context.Background(), Options{
		SourceKind: models.SourceKindArticle, SourceURL: "https://example.com/posts/engines",
		Title: "Decision engines", Markdown: words(400),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Triage != gate.TriageKept || result.FilePath != "demo/raw/articles/decision-engines.md" {
		t.Fatalf("result = %+v", result)
	}
	rel := "raw/articles/decision-engines.md"
	doc, err := corpus.ReadDocument(s.Root(), rel, corpus.KindSource)
	if err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"ingest_batch": "url-2026-09-24-abc123", "ingest_query": "reading list", "triage": "kept"} {
		if got := frontmatter.GetString(doc.Frontmatter, key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
	row, ok := s.State.Get(rel)
	if !ok {
		t.Fatal("no state row for the ingested source")
	}
	for _, key := range []string{"ingest_batch", "ingest_query", "triage"} {
		if row.Written[key] != corpus.ValueHash(doc.Frontmatter[key]) {
			t.Errorf("state row must record %s as kb-written: %+v", key, row.Written)
		}
	}

	again, err := run.Ingest(context.Background(), Options{
		SourceKind: models.SourceKindArticle, SourceURL: "https://www.example.com/posts/engines/", Title: "Decision engines", Markdown: words(400),
	})
	if err != nil || again.Triage != gate.TriageDuplicateSkipped || again.DuplicateOf != rel {
		t.Fatalf("second ingest of the same URL in a run = %+v, %v", again, err)
	}

	summary, err := run.Finish(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if summary.Tally.Kept != 1 || summary.Tally.Duplicates != 1 || summary.Classify == nil || summary.Link == nil {
		t.Fatalf("summary = %+v", summary)
	}
	if asked := fake.QuestionCount("paywall_or_login"); asked != 1 {
		t.Fatalf("classify must reuse the gate judgment receipt; paywall_or_login asked %d times", asked)
	}
	logText, err := os.ReadFile(filepath.Join(s.Root(), "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"## [2026-09-24] ingest | url-2026-09-24-abc123 (url)",
		"kept 1, review 0, quarantined 0, skipped 0, duplicates 1, undecided 0 · decisions US$",
		"- kept: `Decision engines` → `demo/raw/articles/decision-engines.md`",
		"- duplicate-skipped (duplicate): `Decision engines` — duplicate of raw/articles/decision-engines.md (url)",
	} {
		if !strings.Contains(string(logText), want) {
			t.Errorf("log.md misses %q:\n%s", want, logText)
		}
	}
}

func TestRunForceSkipsGates(t *testing.T) {
	t.Parallel()
	s, _ := openRunSession(t)
	run := NewRun(s, RunOptions{Command: "url", Force: true})
	for range 2 {
		result, err := run.Ingest(context.Background(), Options{
			SourceKind: models.SourceKindArticle, SourceURL: "https://example.com/", Title: "Home", Markdown: "# Home\n",
		})
		if err != nil || result.Triage != gate.TriageKept {
			t.Fatalf("forced result = %+v, %v", result, err)
		}
	}
	if lines := ModeLines(s, true); len(lines) != 1 || !strings.Contains(lines[0], "--force") {
		t.Fatalf("mode lines = %v", lines)
	}
	if _, err := run.Finish(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestUniqueGatedSlugAvoidsQuarantine(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, rel := range []string{"raw/articles/a.md", "raw/_quarantine/articles/a-2.md"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	slug, err := uniqueGatedSlug(root, "articles", "a")
	if err != nil || slug != "a-3" {
		t.Fatalf("slug = %q, %v", slug, err)
	}
	if got := WouldBePath(models.SourceKindYouTubeTranscript, "My Talk!"); got != "raw/youtube/my-talk.md" {
		t.Fatalf("WouldBePath = %q", got)
	}
}
