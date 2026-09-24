//go:build integration

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
)

// TestCLIIntegrationClassifyBootstrap runs the bootstrap flow of spec §5.2
// through the real commands and session opening: classify on a topic
// without articles, vocabulary draft and accept, then classify
// --only-missing fills concepts.
func TestCLIIntegrationClassifyBootstrap(t *testing.T) {
	fake := fakes.NewOpenRouter(func(_ fakes.Call, q fakes.Question) any {
		switch q.ID {
		case "role":
			return fakes.Pick("core", 0.9)
		case "kind":
			return fakes.Pick("article_or_essay", 0.9)
		case "enough_text_to_judge":
			return fakes.Noul(0.9)
		case "depth":
			return fakes.Level(1)
		case "primary_concept":
			criteria, _ := q.Criteria.(map[string]any)
			for option, description := range criteria {
				if text, _ := description.(string); strings.HasPrefix(text, "Decision engines") {
					return fakes.Pick(option, 0.9)
				}
			}
		}
		return nil
	}, func(schemaName, _, prompt string) any {
		switch schemaName {
		case "document_literals":
			return map[string]any{"summary": "Notes about decision engines and judges.", "entities": []string{}, "questions": []string{"What is a judge?", "How are decisions typed?", "Why calibrate?"}}
		case "concept_list":
			concepts := []any{map[string]any{"title": "Decision engines", "criterion": "Documents about engines that answer typed questions with probabilities."}}
			for i := 1; i < 12; i++ {
				concepts = append(concepts, map[string]any{"title": fmt.Sprintf("Subject %02d", i), "criterion": fmt.Sprintf("Documents about subject %d.", i)})
			}
			return map[string]any{"concepts": concepts}
		case "concept_aliases":
			return map[string]any{"aliases": []string{}}
		}
		return map[string]any{}
	})
	t.Cleanup(fake.Close)

	configPath := filepath.Join(t.TempDir(), "kb.toml")
	writeFile(t, configPath, "")
	t.Setenv("APP_CONFIG", configPath)
	t.Setenv("OPENROUTER_API_KEY", "test-key")
	t.Setenv("OPENROUTER_API_URL", fake.URL)

	vaultRoot := t.TempDir()
	info := scaffoldTopicForIntegration(t, vaultRoot, "demo", "Demo", "demo")
	sources := []string{"a", "b"}
	for _, name := range sources {
		writeMarkdownDocument(t, info.RootPath, "raw/articles/"+name+".md", map[string]any{
			"title": "Source " + name, "type": "source", "stage": "raw", "source_kind": "document",
		}, "# Source "+name+"\n\nDecision engines answer typed questions about documents.\n")
	}

	stdout, stderr := runCLIWithStreams(t, "classify", "demo", "--vault", vaultRoot, "--format", "json")
	var first classify.Report
	if err := json.Unmarshal([]byte(stdout), &first); err != nil {
		t.Fatalf("classify stdout is not a JSON report: %v\n%s", err, stdout)
	}
	if !first.NoVocabulary || first.Judged != 2 {
		t.Fatalf("first report = %+v", first)
	}
	if !strings.Contains(stderr, "no vocabulary: run `kb topic vocabulary demo --draft`") || !strings.Contains(stderr, "decisions:") {
		t.Fatalf("stderr misses the vocabulary hint or the run summary:\n%s", stderr)
	}

	draft := runCLI(t, "topic", "vocabulary", "demo", "--draft", "--vault", vaultRoot)
	if !strings.Contains(draft, "Decision engines") || !strings.Contains(draft, "mentions") {
		t.Fatalf("draft output:\n%s", draft)
	}
	accepted := runCLI(t, "topic", "vocabulary", "demo", "--accept", "--vault", vaultRoot)
	if lines := strings.Split(strings.TrimSpace(accepted), "\n"); len(lines) != 12 || lines[0] != "wiki/concepts/Decision engines.md" {
		t.Fatalf("accept output:\n%s", accepted)
	}
	if _, err := os.Stat(filepath.Join(info.RootPath, "wiki", "concepts", "Subject 01.md")); err != nil {
		t.Fatalf("stub article missing: %v", err)
	}

	table := runCLI(t, "classify", "demo", "--only-missing", "--vault", vaultRoot)
	if !strings.Contains(table, "written:concepts") {
		t.Fatalf("only-missing table:\n%s", table)
	}
	for _, name := range sources {
		values, _ := readMarkdownDocument(t, filepath.Join(info.RootPath, "raw", "articles", name+".md"))
		if got := frontmatter.GetStringSlice(values, "concepts"); len(got) != 1 || got[0] != "[[Decision engines]]" {
			t.Fatalf("%s concepts = %v", name, got)
		}
	}
}
