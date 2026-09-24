package scope

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

var testNow = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

// testDraft is a valid contract draft used across tests.
var testDraft = &contract.Contract{
	Purpose:            "Collect material on typed decision models and on building software with them.",
	Core:               []string{"Typed decision models such as Jev"},
	Adjacent:           []string{"LLM-as-judge evaluation methods, kept as comparison"},
	CollectedOnPurpose: []string{"The brand audit pages collected for the audit"},
	OutOfScope:         []string{"Sports results with no decision-model content"},
}

// newTestTopic scaffolds topic "demo" in a temporary vault.
func newTestTopic(t *testing.T) (vault, root string) {
	t.Helper()
	vault = t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	return vault, info.RootPath
}

// openTestSession opens a session against the fake OpenRouter.
func openTestSession(t *testing.T, vault string, fake *fakes.OpenRouter) *session.Session {
	t.Helper()
	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = fake.URL
	cfg.Decisions.Concurrency = 4
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb topic contract", Now: func() time.Time { return testNow }})
	if err != nil {
		t.Fatalf("session.Open: %v", err)
	}
	return s
}

// fakeServer starts a fake OpenRouter closed at cleanup.
func fakeServer(t *testing.T, decide fakes.DecideFunc, generate fakes.GenerateFunc) *fakes.OpenRouter {
	t.Helper()
	fake := fakes.NewOpenRouter(decide, generate)
	t.Cleanup(fake.Close)
	return fake
}

// writeSource writes a raw source with a title and body.
func writeSource(t *testing.T, root, rel, title, body string) {
	t.Helper()
	content, err := frontmatter.Generate(map[string]any{"title": title, "type": "source", "stage": "raw", "source_kind": "article"}, body)
	if err != nil {
		t.Fatalf("frontmatter.Generate: %v", err)
	}
	writeFile(t, filepath.Join(root, filepath.FromSlash(rel)), content)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// docs builds in-memory documents at the given paths.
func docs(paths ...string) []*corpus.Document {
	out := make([]*corpus.Document, 0, len(paths))
	for _, p := range paths {
		out = append(out, &corpus.Document{Path: p, Title: strings.TrimSuffix(filepath.Base(p), ".md"), Kind: corpus.KindSource})
	}
	return out
}

// roleByTitle answers `role` off_topic (0.9) for documents whose state
// mentions "Offtopic", core (0.9) otherwise.
func roleByTitle(call fakes.Call, q fakes.Question) any {
	if q.ID != "role" {
		return nil
	}
	if strings.Contains(call.StateString(), "Offtopic") {
		return fakes.Pick("off_topic", 0.9)
	}
	return fakes.Pick("core", 0.9)
}
