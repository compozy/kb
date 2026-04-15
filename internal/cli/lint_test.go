package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	klint "github.com/compozy/kb/internal/lint"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

func TestLintCommandAcceptsPositionalTopicSlug(t *testing.T) {
	originalRunLint := runLintEngine
	originalSaveReport := saveLintEngineReport
	originalResolve := resolveLintVaultQuery
	originalGetwd := lintGetwd
	t.Cleanup(func() {
		runLintEngine = originalRunLint
		saveLintEngineReport = originalSaveReport
		resolveLintVaultQuery = originalResolve
		lintGetwd = originalGetwd
	})

	var gotQuery vault.VaultQueryOptions
	lintGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveLintVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/demo-topic",
			TopicSlug: "demo-topic",
		}, nil
	}
	runLintEngine = func(topicPath string, options klint.LintOptions) ([]models.LintIssue, error) {
		return nil, nil
	}
	saveLintEngineReport = func(topicPath string, issues []models.LintIssue, now time.Time) (string, error) {
		return topicPath, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"lint", "demo-topic"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	wantQuery := vault.VaultQueryOptions{
		CWD:   "/workspace/repo",
		Topic: "demo-topic",
	}
	if gotQuery != wantQuery {
		t.Fatalf("vault query = %#v, want %#v", gotQuery, wantQuery)
	}
}

func TestLintCommandSupportsTableJSONAndTSVOutput(t *testing.T) {
	originalRunLint := runLintEngine
	originalSaveReport := saveLintEngineReport
	originalResolve := resolveLintVaultQuery
	originalGetwd := lintGetwd
	t.Cleanup(func() {
		runLintEngine = originalRunLint
		saveLintEngineReport = originalSaveReport
		resolveLintVaultQuery = originalResolve
		lintGetwd = originalGetwd
	})

	lintGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveLintVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/demo-topic",
			TopicSlug: "demo-topic",
		}, nil
	}
	runLintEngine = func(topicPath string, options klint.LintOptions) ([]models.LintIssue, error) {
		return []models.LintIssue{{
			Severity: models.SeverityError,
			Kind:     models.LintIssueKindDeadLink,
			FilePath: "wiki/concepts/Overview.md",
			Target:   "Missing Page",
			Message:  "dead wikilink",
		}}, nil
	}
	saveLintEngineReport = func(topicPath string, issues []models.LintIssue, now time.Time) (string, error) {
		return topicPath, nil
	}

	testCases := []struct {
		name       string
		format     string
		assertions func(t *testing.T, stdout string)
	}{
		{
			name:   "table",
			format: "table",
			assertions: func(t *testing.T, stdout string) {
				t.Helper()
				for _, fragment := range []string{"severity", "kind", "filePath", "dead-link", "dead wikilink"} {
					if !strings.Contains(stdout, fragment) {
						t.Fatalf("expected table output to contain %q, got:\n%s", fragment, stdout)
					}
				}
			},
		},
		{
			name:   "json",
			format: "json",
			assertions: func(t *testing.T, stdout string) {
				t.Helper()
				var decoded []map[string]any
				if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
					t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout)
				}
				if len(decoded) != 1 || decoded[0]["kind"] != "dead-link" {
					t.Fatalf("unexpected JSON payload %#v", decoded)
				}
			},
		},
		{
			name:   "tsv",
			format: "tsv",
			assertions: func(t *testing.T, stdout string) {
				t.Helper()
				if !strings.HasPrefix(stdout, "severity\tkind\tfilePath\ttarget\tmessage\n") {
					t.Fatalf("unexpected TSV header:\n%s", stdout)
				}
				if !strings.Contains(stdout, "error\tdead-link\twiki/concepts/Overview.md\tMissing Page\tdead wikilink") {
					t.Fatalf("unexpected TSV row:\n%s", stdout)
				}
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			command := newRootCommand()
			var stdout bytes.Buffer
			command.SetOut(&stdout)
			command.SetErr(new(bytes.Buffer))
			command.SetArgs([]string{"lint", "demo-topic", "--format", testCase.format})

			if err := command.ExecuteContext(context.Background()); err != nil {
				t.Fatalf("ExecuteContext returned error: %v", err)
			}

			testCase.assertions(t, stdout.String())
		})
	}
}

func TestLintCommandPassesJavaGovernanceThresholdFlags(t *testing.T) {
	originalRunLint := runLintEngine
	originalSaveReport := saveLintEngineReport
	originalResolve := resolveLintVaultQuery
	originalGetwd := lintGetwd
	t.Cleanup(func() {
		runLintEngine = originalRunLint
		saveLintEngineReport = originalSaveReport
		resolveLintVaultQuery = originalResolve
		lintGetwd = originalGetwd
	})

	lintGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveLintVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/demo-topic",
			TopicSlug: "demo-topic",
		}, nil
	}
	var gotOptions klint.LintOptions
	runLintEngine = func(topicPath string, options klint.LintOptions) ([]models.LintIssue, error) {
		gotOptions = options
		return nil, nil
	}
	saveLintEngineReport = func(topicPath string, issues []models.LintIssue, now time.Time) (string, error) {
		return topicPath, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"lint", "demo-topic",
		"--java-max-parse-errors", "2",
		"--java-max-fallback-warnings", "5",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.JavaGovernance.MaxParseErrors != 2 {
		t.Fatalf("MaxParseErrors = %d, want 2", gotOptions.JavaGovernance.MaxParseErrors)
	}
	if gotOptions.JavaGovernance.MaxFallbackWarnings != 5 {
		t.Fatalf("MaxFallbackWarnings = %d, want 5", gotOptions.JavaGovernance.MaxFallbackWarnings)
	}
}

func TestLintCommandSaveFlagWritesReport(t *testing.T) {
	originalRunLint := runLintEngine
	originalSaveReport := saveLintEngineReport
	originalResolve := resolveLintVaultQuery
	originalGetwd := lintGetwd
	originalNow := lintNow
	t.Cleanup(func() {
		runLintEngine = originalRunLint
		saveLintEngineReport = originalSaveReport
		resolveLintVaultQuery = originalResolve
		lintGetwd = originalGetwd
		lintNow = originalNow
	})

	var gotTopicPath string
	var gotIssues []models.LintIssue
	var gotNow time.Time

	lintGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveLintVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/demo-topic",
			TopicSlug: "demo-topic",
		}, nil
	}
	runLintEngine = func(topicPath string, options klint.LintOptions) ([]models.LintIssue, error) {
		return []models.LintIssue{{
			Severity: models.SeverityWarning,
			Kind:     models.LintIssueKindOrphan,
			FilePath: "wiki/concepts/Overview.md",
			Message:  "orphan article",
		}}, nil
	}
	lintNow = func() time.Time {
		return time.Date(2026, time.April, 11, 10, 30, 0, 0, time.UTC)
	}
	saveLintEngineReport = func(topicPath string, issues []models.LintIssue, now time.Time) (string, error) {
		gotTopicPath = topicPath
		gotIssues = append([]models.LintIssue(nil), issues...)
		gotNow = now
		return "/vault/demo-topic/outputs/reports/2026-04-11-lint.md", nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"lint", "demo-topic", "--save"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotTopicPath != "/vault/demo-topic" {
		t.Fatalf("report topic path = %q, want /vault/demo-topic", gotTopicPath)
	}
	if len(gotIssues) != 1 || gotIssues[0].Kind != models.LintIssueKindOrphan {
		t.Fatalf("report issues = %#v, want orphan issue", gotIssues)
	}
	if !gotNow.Equal(time.Date(2026, time.April, 11, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("report timestamp = %v, want fixed test time", gotNow)
	}
}

func TestLintCommandRejectsConflictingTopicSelectors(t *testing.T) {
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"lint", "demo-topic", "--topic", "other-topic"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected conflicting topic selector error")
	}
	if !strings.Contains(err.Error(), "conflicting topic selectors") {
		t.Fatalf("unexpected error %q", err)
	}
}

func TestLintCommandHelpShowsFlags(t *testing.T) {
	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"lint", "--help"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	for _, flag := range []string{
		"--format",
		"--save",
		"--topic",
		"--vault",
		"--java-max-parse-errors",
		"--java-max-fallback-warnings",
	} {
		if !strings.Contains(stdout.String(), flag) {
			t.Fatalf("expected help output to contain %q, got:\n%s", flag, stdout.String())
		}
	}
}
