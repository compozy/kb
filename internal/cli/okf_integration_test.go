//go:build integration

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	kokf "github.com/compozy/kb/internal/okf"
)

func TestCLIIntegrationOKFPromoteAndCheck(t *testing.T) {
	vaultRoot := t.TempDir()
	wikiTopic := scaffoldTopicForIntegration(t, vaultRoot, "research", "Research", "ops")
	okfTopic := runCLIJSON[models.TopicInfo](t,
		"topic", "new", "catalog", "Catalog", "ops",
		"--mode", "okf",
		"--vault", vaultRoot,
	)
	if okfTopic.Mode != models.TopicModeOKF {
		t.Fatalf("OKF topic mode = %q, want okf", okfTopic.Mode)
	}

	sourcePath := filepath.Join(wikiTopic.RootPath, "wiki", "concepts", "Alpha Note.md")
	writeFile(t, sourcePath, strings.Join([]string{
		"---",
		"title: Alpha Note",
		"type: wiki",
		"stage: compiled",
		"tags: [ops]",
		"---",
		"Alpha note explains the operations flow. See [[research/wiki/concepts/Beta Note|Beta]].",
	}, "\n"))

	promoteOutput := runCLI(t,
		"promote", sourcePath,
		"--to", okfTopic.Slug,
		"--type", "Playbook",
		"--description", "Operational flow.",
		"--vault", vaultRoot,
	)
	var result kokf.ConceptResult
	if err := json.Unmarshal([]byte(promoteOutput), &result); err != nil {
		t.Fatalf("promote output is not JSON: %v\n%s", err, promoteOutput)
	}
	if result.WrittenPath != "alpha-note.md" {
		t.Fatalf("written path = %q, want alpha-note.md", result.WrittenPath)
	}

	sourceAfter := readFile(t, sourcePath)
	if !strings.Contains(sourceAfter, "stage: compiled") {
		t.Fatalf("source was unexpectedly changed:\n%s", sourceAfter)
	}
	conceptPath := filepath.Join(okfTopic.RootPath, "alpha-note.md")
	values, body := readMarkdownDocument(t, conceptPath)
	if values["type"] != "Playbook" || values["description"] != "Operational flow." {
		t.Fatalf("unexpected concept frontmatter: %#v", values)
	}
	if !strings.Contains(body, "[Beta](beta-note.md)") {
		t.Fatalf("body missing relative OKF link:\n%s", body)
	}
	if _, ok := values["stage"]; ok {
		t.Fatalf("wiki stage leaked into OKF concept: %#v", values)
	}

	checkOutput := runCLI(t, "okf", "check", okfTopic.Slug, "--format", "json", "--vault", vaultRoot)
	var issues []models.LintIssue
	if err := json.Unmarshal([]byte(checkOutput), &issues); err != nil {
		t.Fatalf("check output is not JSON: %v\n%s", err, checkOutput)
	}
	if len(issues) != 0 {
		t.Fatalf("freshly promoted bundle has issues: %#v", issues)
	}

	index := readFile(t, filepath.Join(okfTopic.RootPath, "index.md"))
	if !strings.Contains(index, "## Playbook") || !strings.Contains(index, "[Alpha Note](alpha-note.md)") {
		t.Fatalf("index.md missing promoted concept:\n%s", index)
	}
	logContent := readFile(t, filepath.Join(okfTopic.RootPath, "log.md"))
	if !strings.Contains(logContent, "## "+frontmatter.DateLayout[:4]) && !strings.Contains(logContent, "**Creation**") {
		t.Fatalf("log.md missing promotion entry:\n%s", logContent)
	}
}

func TestCLIIntegrationOKFCheckFailsStrictWarnings(t *testing.T) {
	vaultRoot := t.TempDir()
	okfTopic := runCLIJSON[models.TopicInfo](t,
		"topic", "new", "catalog", "Catalog", "ops",
		"--mode", "okf",
		"--vault", vaultRoot,
	)
	writeFile(t, filepath.Join(okfTopic.RootPath, "concept.md"), "---\ntype: Unknown\n---\nBody.\n")

	errText := runCLIError(t, "okf", "check", okfTopic.Slug, "--strict", "--vault", vaultRoot)
	if !strings.Contains(errText, "found") {
		t.Fatalf("strict check error = %q", errText)
	}
}

func readFile(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read %s: %v", filePath, err)
	}
	return string(content)
}
