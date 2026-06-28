package okf

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
)

func TestPromoteWritesConceptIndexAndLog(t *testing.T) {
	t.Parallel()

	t.Run("Should write promoted concept index and log entry", func(t *testing.T) {
		vaultPath := t.TempDir()
		sourceTopic := filepath.Join(vaultPath, "research")
		targetTopic := filepath.Join(vaultPath, "catalog")
		mkdirAll(t, filepath.Join(sourceTopic, "wiki", "concepts"))
		mkdirAll(t, targetTopic)
		writeFile(t, filepath.Join(sourceTopic, "CLAUDE.md"), "# Research\n")
		writeFile(t, filepath.Join(targetTopic, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(targetTopic, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Old Index\n")
		writeFile(t, filepath.Join(targetTopic, "log.md"), "# Directory Update Log\n")
		sourcePath := filepath.Join(sourceTopic, "wiki", "concepts", "Alpha Note.md")
		sourceContent := strings.Join([]string{
			"---",
			"title: Alpha Note",
			"type: wiki",
			"stage: compiled",
			"tags: [systems, alpha]",
			"---",
			"Alpha note explains the operational flow. See [[research/wiki/concepts/Beta Note|Beta]] and [[research/wiki/concepts/Alpha Note#details|details]].",
			"",
			"## Details",
		}, "\n")
		writeFile(t, sourcePath, sourceContent)

		result, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: sourcePath,
			VaultPath:     vaultPath,
			TargetTopic: models.TopicInfo{
				Slug:     "catalog",
				Mode:     models.TopicModeOKF,
				RootPath: targetTopic,
			},
			Type:        "Playbook",
			Description: "Operational alpha note.",
			Types:       []string{"Playbook"},
			Clock: func() time.Time {
				return time.Date(2026, 6, 27, 10, 11, 12, 0, time.UTC)
			},
		})
		if err != nil {
			t.Fatalf("Promote returned error: %v", err)
		}
		if result.WrittenPath != "alpha-note.md" {
			t.Fatalf("written path = %q, want alpha-note.md", result.WrittenPath)
		}
		if result.LinksRewritten != 2 {
			t.Fatalf("links rewritten = %d, want 2", result.LinksRewritten)
		}
		if len(result.UnresolvedLinks) != 1 || result.UnresolvedLinks[0] != "beta-note.md" {
			t.Fatalf("unresolved links = %#v, want beta-note.md", result.UnresolvedLinks)
		}
		if got := readFile(t, sourcePath); got != sourceContent {
			t.Fatalf("source document changed:\n%s", got)
		}

		values, body := parseMarkdown(t, filepath.Join(targetTopic, "alpha-note.md"))
		for key, want := range map[string]any{
			"type":        "Playbook",
			"title":       "Alpha Note",
			"description": "Operational alpha note.",
			"timestamp":   "2026-06-27T10:11:12Z",
		} {
			if got := values[key]; got != want {
				t.Fatalf("frontmatter[%s] = %#v, want %#v", key, got, want)
			}
		}
		if _, ok := values["stage"]; ok {
			t.Fatalf("wiki stage leaked into OKF frontmatter: %#v", values)
		}
		if !strings.Contains(body, "[Beta](beta-note.md)") || !strings.Contains(body, "[details](alpha-note.md#details)") {
			t.Fatalf("body missing transformed links:\n%s", body)
		}

		index := readFile(t, filepath.Join(targetTopic, "index.md"))
		for _, fragment := range []string{"okf_version: \"0.1\"", "## Playbook", "[Alpha Note](alpha-note.md)", "Operational alpha note."} {
			if !strings.Contains(index, fragment) {
				t.Fatalf("index.md missing %q:\n%s", fragment, index)
			}
		}
		log := readFile(t, filepath.Join(targetTopic, "log.md"))
		for _, fragment := range []string{"## 2026-06-27", "**Creation**", "[Alpha Note](alpha-note.md)", "`wiki/concepts/Alpha Note.md`"} {
			if !strings.Contains(log, fragment) {
				t.Fatalf("log.md missing %q:\n%s", fragment, log)
			}
		}
	})
}

func TestPromoteRejectsNonOKFTargetBeforeWriting(t *testing.T) {
	t.Parallel()

	t.Run("Should reject non-OKF target before writing", func(t *testing.T) {
		targetTopic := t.TempDir()
		_, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: "missing.md",
			TargetTopic: models.TopicInfo{
				Slug:     "wiki-topic",
				Mode:     models.TopicModeWiki,
				RootPath: targetTopic,
			},
			Type: "Reference",
		})
		if err == nil || !strings.Contains(err.Error(), "target topic must use mode okf") {
			t.Fatalf("error = %v, want non-OKF target rejection", err)
		}
		if entries, readErr := os.ReadDir(targetTopic); readErr != nil || len(entries) != 0 {
			t.Fatalf("target changed before rejection: entries=%d err=%v", len(entries), readErr)
		}
	})
}

func TestPromoteWarnsWhenSourceBodyHasNoSentenceFallback(t *testing.T) {
	t.Parallel()

	t.Run("Should warn when source body has no sentence fallback", func(t *testing.T) {
		vaultPath := t.TempDir()
		sourceTopic := filepath.Join(vaultPath, "research")
		targetTopic := filepath.Join(vaultPath, "catalog")
		mkdirAll(t, filepath.Join(sourceTopic, "wiki", "concepts"))
		mkdirAll(t, targetTopic)
		writeFile(t, filepath.Join(sourceTopic, "CLAUDE.md"), "# Research\n")
		writeFile(t, filepath.Join(targetTopic, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(targetTopic, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(targetTopic, "log.md"), "# Directory Update Log\n")
		sourcePath := filepath.Join(sourceTopic, "wiki", "concepts", "No Sentence.md")
		writeFile(t, sourcePath, "---\ntitle: No Sentence\n---\nbody text without punctuation")

		result, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: sourcePath,
			VaultPath:     vaultPath,
			TargetTopic: models.TopicInfo{
				Slug:     "catalog",
				Mode:     models.TopicModeOKF,
				RootPath: targetTopic,
			},
			Type: "Reference",
			Clock: func() time.Time {
				return time.Date(2026, 6, 27, 10, 11, 12, 0, time.UTC)
			},
		})
		if err != nil {
			t.Fatalf("Promote returned error: %v", err)
		}
		if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "source body has no sentence fallback") {
			t.Fatalf("warnings = %#v, want no sentence fallback warning", result.Warnings)
		}

		values, _ := parseMarkdown(t, filepath.Join(targetTopic, "no-sentence.md"))
		if got := values["description"]; got != "" {
			t.Fatalf("description = %#v, want empty string", got)
		}
	})
}

func TestCheckReportsConformanceAndStrictWarnings(t *testing.T) {
	t.Parallel()

	writeBundle := func(t *testing.T) string {
		t.Helper()
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(bundle, "log.md"), "# Directory Update Log\n\n## 2026-06-27\n* **Creation**: Created bundle.\n")
		writeFile(t, filepath.Join(bundle, "good.md"), "---\ntype: Playbook\ntitle: Good\ndescription: Good concept.\ntimestamp: 2026-06-27T10:00:00Z\n---\nBody with [broken](missing.md).\n")
		writeFile(t, filepath.Join(bundle, "missing-type.md"), "---\ntitle: Missing\n---\nBody.\n")
		writeFile(t, filepath.Join(bundle, "missing-fields.md"), "---\ntype: Unknown\n---\nBody.\n")
		return bundle
	}

	testCases := []struct {
		name       string
		strict     bool
		assertions []expectedIssue
	}{
		{
			name:   "Should report lenient conformance diagnostics",
			strict: false,
			assertions: []expectedIssue{
				{severity: models.SeverityError, filePath: "missing-type.md", target: "type", messageContains: "concept frontmatter must include non-empty type"},
				{severity: models.SeverityWarning, filePath: "missing-fields.md", target: "title", messageContains: `producer field "title" is missing`},
				{severity: models.SeverityWarning, filePath: "missing-fields.md", target: "type", messageContains: `type "Unknown" is outside the configured OKF vocabulary`},
			},
		},
		{
			name:   "Should promote warnings to errors in strict mode",
			strict: true,
			assertions: []expectedIssue{
				{severity: models.SeverityError, filePath: "missing-fields.md", target: "title", messageContains: `producer field "title" is missing`},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bundle := writeBundle(t)
			issues, err := Check(context.Background(), bundle, CheckOptions{
				Types:  []string{"Playbook"},
				Strict: testCase.strict,
			})
			if err != nil {
				t.Fatalf("Check returned error: %v", err)
			}
			for _, assertion := range testCase.assertions {
				assertIssue(t, issues, assertion)
			}
			if !HasErrors(issues) {
				t.Fatal("HasErrors = false, want true for missing type")
			}
		})
	}
}

func TestCheckAllowsLenientExternalBundleTraits(t *testing.T) {
	t.Parallel()

	t.Run("Should allow lenient external bundle traits", func(t *testing.T) {
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "index.md"), "# Index\n")
		writeFile(t, filepath.Join(bundle, "concept.md"), "---\ntype: External Type\n---\nSee [missing](missing.md).\n")

		issues, err := Check(context.Background(), bundle, CheckOptions{})
		if err != nil {
			t.Fatalf("Check returned error: %v", err)
		}
		if HasErrors(issues) {
			t.Fatalf("lenient external bundle should not have errors: %#v", issues)
		}
	})
}

func TestAllocateConceptPathReservesPath(t *testing.T) {
	t.Parallel()

	t.Run("Should reserve concept paths with exclusive creation", func(t *testing.T) {
		bundle := t.TempDir()
		firstRelativePath, firstAbsolutePath, firstFile, err := allocateConceptPath(bundle, "alpha")
		if err != nil {
			t.Fatalf("allocate first path: %v", err)
		}
		t.Cleanup(func() { _ = firstFile.Close() })
		if firstRelativePath != "alpha.md" {
			t.Fatalf("first relative path = %q, want alpha.md", firstRelativePath)
		}
		if _, err := os.Stat(firstAbsolutePath); err != nil {
			t.Fatalf("reserved first path does not exist: %v", err)
		}

		secondRelativePath, _, secondFile, err := allocateConceptPath(bundle, "alpha")
		if err != nil {
			t.Fatalf("allocate second path: %v", err)
		}
		t.Cleanup(func() { _ = secondFile.Close() })
		if secondRelativePath != "alpha-2.md" {
			t.Fatalf("second relative path = %q, want alpha-2.md", secondRelativePath)
		}
	})
}

func TestRegenerateIndexReturnsConceptParseErrors(t *testing.T) {
	t.Parallel()

	t.Run("Should return parse errors instead of dropping malformed concepts", func(t *testing.T) {
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(bundle, "bad.md"), "---\ntype: Playbook\n\nBody without closed frontmatter.\n")

		err := RegenerateIndex(bundle)
		if err == nil {
			t.Fatal("expected malformed concept frontmatter error")
		}
		if !strings.Contains(err.Error(), "parse concept bad.md") {
			t.Fatalf("error = %q, want parse concept bad.md", err.Error())
		}
	})
}

type expectedIssue struct {
	severity        models.DiagnosticSeverity
	filePath        string
	target          string
	messageContains string
}

func assertIssue(t *testing.T, issues []models.LintIssue, expected expectedIssue) {
	t.Helper()
	for _, issue := range issues {
		if issue.Severity == expected.severity &&
			issue.FilePath == expected.filePath &&
			issue.Target == expected.target &&
			strings.Contains(issue.Message, expected.messageContains) {
			return
		}
	}
	t.Fatalf(
		"missing issue severity=%s file=%s target=%s message containing %q in %#v",
		expected.severity,
		expected.filePath,
		expected.target,
		expected.messageContains,
		issues,
	)
}

func parseMarkdown(t *testing.T, filePath string) (map[string]any, string) {
	t.Helper()
	values, body, err := frontmatter.Parse(readFile(t, filePath))
	if err != nil {
		t.Fatalf("parse %s: %v", filePath, err)
	}
	return values, body
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
}

func writeFile(t *testing.T, filePath string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatalf("mkdir parent for %s: %v", filePath, err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", filePath, err)
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
