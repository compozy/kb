package classify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

var testNow = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

// testContract is an accepted selection contract used across tests.
var testContract = &contract.Contract{
	Purpose:                 "Collect material on typed decision models and on building software with them.",
	Core:                    []string{"Typed decision models such as Jev"},
	Adjacent:                []string{"LLM-as-judge evaluation methods, kept as comparison"},
	CollectedOnPurpose:      []string{"The brand audit pages collected for the audit"},
	CollectedOnPurposePaths: []string{"raw/audit/**"},
	OutOfScope:              []string{"Sports results with no decision-model content"},
}

// newTestTopic scaffolds a wiki topic in a temporary vault.
func newTestTopic(t *testing.T) (vault, root string) {
	t.Helper()
	vault = t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	return vault, info.RootPath
}

// setTestContract activates testContract (or c) on the topic.
func setTestContract(t *testing.T, root string, c *contract.Contract) {
	t.Helper()
	if c == nil {
		c = testContract
	}
	if err := contract.SetContract(root, c); err != nil {
		t.Fatalf("SetContract: %v", err)
	}
}

// appendTopicYAML appends raw YAML to topic.yaml.
func appendTopicYAML(t *testing.T, root, text string) {
	t.Helper()
	path := filepath.Join(root, "topic.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte(text)...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// openTestSession opens a session against the fake OpenRouter at apiURL.
func openTestSession(t *testing.T, vault, apiURL string) *session.Session {
	t.Helper()
	return openTestSessionWith(t, vault, apiURL, session.Flags{})
}

// openTestSessionWith opens a session with per-run flags (budget, mode).
func openTestSessionWith(t *testing.T, vault, apiURL string, flags session.Flags) *session.Session {
	t.Helper()
	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = apiURL
	cfg.Decisions.Concurrency = 4
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb classify", Flags: flags, Now: func() time.Time { return testNow }})
	if err != nil {
		t.Fatalf("session.Open: %v", err)
	}
	return s
}

// writeDoc writes a markdown document with frontmatter under the topic.
func writeDoc(t *testing.T, root, rel string, values map[string]any, body string) string {
	t.Helper()
	content, err := frontmatter.Generate(values, body)
	if err != nil {
		t.Fatalf("frontmatter.Generate: %v", err)
	}
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// source writes a raw source.
func source(t *testing.T, root, rel, title, url, body string) string {
	t.Helper()
	values := map[string]any{"title": title, "type": "source", "stage": "raw", "source_kind": "article"}
	if url != "" {
		values["source_url"] = url
	}
	return writeDoc(t, root, rel, values, body)
}

// longBody returns a body of about n words mentioning subject.
func longBody(subject string, n int) string {
	var b strings.Builder
	b.WriteString("# " + subject + "\n\n")
	for i := range n / 10 {
		if i%5 == 0 {
			b.WriteString("\n")
		}
		b.WriteString("This paragraph discusses " + subject + " with concrete detail and examples here. ")
	}
	return b.String() + "\n"
}

// fakeServer starts a fake OpenRouter with the given rules.
func fakeServer(t *testing.T, decide fakes.DecideFunc, generate fakes.GenerateFunc) *fakes.OpenRouter {
	t.Helper()
	fake := fakes.NewOpenRouter(decide, generate)
	t.Cleanup(fake.Close)
	return fake
}

// readFile returns a file's content.
func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
