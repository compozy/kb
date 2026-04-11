package lint_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/user/go-devstack/internal/frontmatter"
	"github.com/user/go-devstack/internal/lint"
	"github.com/user/go-devstack/internal/models"
)

const (
	testTopicSlug = "systems-design"
	testDomain    = "systems"
)

func TestColumnsAndRowsAreFormatterCompatible(t *testing.T) {
	t.Parallel()

	columns := lint.Columns()
	wantColumns := []string{"severity", "kind", "filePath", "target", "message"}
	if !reflect.DeepEqual(columns, wantColumns) {
		t.Fatalf("Columns() = %#v, want %#v", columns, wantColumns)
	}

	rows := lint.Rows([]models.LintIssue{{
		Kind:     models.LintIssueKindDeadLink,
		Severity: models.SeverityError,
		FilePath: "wiki/concepts/Broken.md",
		Message:  "wikilink target does not exist",
		Target:   "Missing Article",
	}})
	if len(rows) != 1 {
		t.Fatalf("Rows() len = %d, want 1", len(rows))
	}

	row := rows[0]
	if row["severity"] != models.SeverityError {
		t.Fatalf("severity = %#v, want %q", row["severity"], models.SeverityError)
	}
	if row["kind"] != models.LintIssueKindDeadLink {
		t.Fatalf("kind = %#v, want %q", row["kind"], models.LintIssueKindDeadLink)
	}
	if row["filePath"] != "wiki/concepts/Broken.md" {
		t.Fatalf("filePath = %#v", row["filePath"])
	}
	if row["target"] != "Missing Article" {
		t.Fatalf("target = %#v", row["target"])
	}
}

func TestLintDetectsDeadLinksWithoutFlaggingValidDisplayLinks(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/source-note.md", sourceFrontmatter("Source Note", "article", "2026-04-10"), "# Source Note\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Target Concept.md", conceptFrontmatter("Target Concept", "2026-04-11", []string{"[[Source Note]]"}), "# Target Concept\n\nSee [[Source Note]].\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Target Concept|Display]]\n[[Missing Article]]\n")

	issues := mustLint(t, topicPath)
	if len(issues) != 1 {
		t.Fatalf("issues len = %d, want 1: %#v", len(issues), issues)
	}

	issue := issues[0]
	if issue.Kind != models.LintIssueKindDeadLink {
		t.Fatalf("kind = %q, want dead-link", issue.Kind)
	}
	if issue.FilePath != "wiki/index/Dashboard.md" {
		t.Fatalf("filePath = %q, want wiki/index/Dashboard.md", issue.FilePath)
	}
	if issue.Target != "Missing Article" {
		t.Fatalf("target = %q, want Missing Article", issue.Target)
	}
}

func TestLintDetectsOrphanArticles(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/source-note.md", sourceFrontmatter("Source Note", "article", "2026-04-10"), "# Source Note\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Orphan Concept.md", conceptFrontmatter("Orphan Concept", "2026-04-11", []string{"[[Source Note]]"}), "# Orphan Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "# Dashboard\n")

	issues := mustLint(t, topicPath)
	if len(issues) != 1 {
		t.Fatalf("issues len = %d, want 1: %#v", len(issues), issues)
	}

	issue := issues[0]
	if issue.Kind != models.LintIssueKindOrphan {
		t.Fatalf("kind = %q, want orphan", issue.Kind)
	}
	if issue.Severity != models.SeverityWarning {
		t.Fatalf("severity = %q, want warning", issue.Severity)
	}
	if issue.FilePath != "wiki/concepts/Orphan Concept.md" {
		t.Fatalf("filePath = %q", issue.FilePath)
	}
}

func TestLintDetectsMissingSourcesAndStaleContent(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/fresh-source.md", sourceFrontmatter("Fresh Source", "article", "2026-04-12"), "# Fresh Source\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Stale Concept.md", conceptFrontmatter("Stale Concept", "2026-04-11", []string{"[[Fresh Source]]", "[[Missing Source]]"}), "# Stale Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Stale Concept|Stale]]\n")

	issues := mustLint(t, topicPath)

	assertHasIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindMissingSource,
		Severity: models.SeverityError,
		FilePath: "wiki/concepts/Stale Concept.md",
		Target:   "Missing Source",
	})
	assertHasIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindStale,
		Severity: models.SeverityWarning,
		FilePath: "wiki/concepts/Stale Concept.md",
		Target:   "Fresh Source",
	})
}

func TestLintDetectsFormatViolations(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/source-note.md", sourceFrontmatter("Source Note", "article", "2026-04-10"), "# Source Note\n")
	writeMarkdownFile(t, topicPath, "raw/articles/missing-kind.md", sourceFrontmatter("Missing Kind", "", "2026-04-10"), "# Missing Kind\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Untitled Concept.md", conceptFrontmatter("", "2026-04-11", []string{"[[Source Note]]"}), "# Untitled Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Untitled Concept|Untitled]]\n")

	issues := mustLint(t, topicPath)

	assertHasIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindFormat,
		Severity: models.SeverityError,
		FilePath: "raw/articles/missing-kind.md",
		Target:   "source_kind",
	})
	assertHasIssue(t, issues, models.LintIssue{
		Kind:     models.LintIssueKindFormat,
		Severity: models.SeverityError,
		FilePath: "wiki/concepts/Untitled Concept.md",
		Target:   "title",
	})
}

func TestLintReturnsIssuesSortedBySeverityThenFilePath(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/current-source.md", sourceFrontmatter("Current Source", "article", "2026-04-12"), "# Current Source\n")
	writeMarkdownFile(t, topicPath, "raw/articles/a-source.md", sourceFrontmatter("A Source", "", "2026-04-10"), "# A Source\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Alpha Concept.md", conceptFrontmatter("Alpha Concept", "2026-04-11", []string{"[[Current Source]]"}), "# Alpha Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Bravo Concept.md", conceptFrontmatter("Bravo Concept", "2026-04-12", []string{"[[Current Source]]"}), "# Bravo Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Alpha Concept|Alpha]]\n[[Missing Link]]\n")

	issues := mustLint(t, topicPath)
	got := make([]string, 0, len(issues))
	for _, issue := range issues {
		got = append(got, issue.FilePath)
	}

	want := []string{
		"raw/articles/a-source.md",
		"wiki/index/Dashboard.md",
		"wiki/concepts/Alpha Concept.md",
		"wiki/concepts/Bravo Concept.md",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("issue order = %#v, want %#v", got, want)
	}
}

func TestLintReturnsEmptySliceForHealthyVault(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/source-note.md", sourceFrontmatter("Source Note", "article", "2026-04-10"), "# Source Note\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Healthy Concept.md", conceptFrontmatter("Healthy Concept", "2026-04-11", []string{"[[Source Note]]"}), "# Healthy Concept\n\nSee [[Source Note]].\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Healthy Concept|Healthy]]\n")

	issues := mustLint(t, topicPath)
	if len(issues) != 0 {
		t.Fatalf("issues = %#v, want empty slice", issues)
	}
}

func TestLintOnMixedVaultDetectsAllRequiredIssueKinds(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "raw/articles/current-source.md", sourceFrontmatter("Current Source", "article", "2026-04-12"), "# Current Source\n")
	writeMarkdownFile(t, topicPath, "raw/articles/broken-source.md", sourceFrontmatter("Broken Source", "", "2026-04-10"), "# Broken Source\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Outdated Concept.md", conceptFrontmatter("Outdated Concept", "2026-04-11", []string{"[[Current Source]]"}), "# Outdated Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Broken Concept.md", conceptFrontmatter("Broken Concept", "2026-04-11", []string{"[[Missing Source]]"}), "# Broken Concept\n\nSee [[Missing Article]].\n")
	writeMarkdownFile(t, topicPath, "wiki/concepts/Lonely Concept.md", conceptFrontmatter("Lonely Concept", "2026-04-12", []string{"[[Current Source]]"}), "# Lonely Concept\n")
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "[[systems-design/wiki/concepts/Outdated Concept|Outdated]]\n[[systems-design/wiki/concepts/Broken Concept|Broken]]\n")

	issues := mustLint(t, topicPath)

	assertKindPresent(t, issues, models.LintIssueKindDeadLink)
	assertKindPresent(t, issues, models.LintIssueKindOrphan)
	assertKindPresent(t, issues, models.LintIssueKindMissingSource)
	assertKindPresent(t, issues, models.LintIssueKindStale)
	assertKindPresent(t, issues, models.LintIssueKindFormat)
}

func TestSaveReportWritesLintReport(t *testing.T) {
	t.Parallel()

	topicPath := newTestTopic(t)
	writeMarkdownFile(t, topicPath, "wiki/index/Dashboard.md", indexFrontmatter("Dashboard"), "# Dashboard\n")

	reportPath, err := lint.SaveReport(topicPath, []models.LintIssue{{
		Kind:     models.LintIssueKindDeadLink,
		Severity: models.SeverityError,
		FilePath: "wiki/concepts/Broken.md",
		Message:  "wikilink target does not exist",
		Target:   "Missing Article",
	}}, time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("SaveReport returned error: %v", err)
	}

	wantPath := filepath.Join(topicPath, "outputs", "reports", "2026-04-11-lint.md")
	if reportPath != wantPath {
		t.Fatalf("reportPath = %q, want %q", reportPath, wantPath)
	}

	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}

	values, body, err := frontmatter.Parse(string(content))
	if err != nil {
		t.Fatalf("parse report frontmatter: %v", err)
	}

	if got := values["stage"]; got != "lint-report" {
		t.Fatalf("stage = %#v, want lint-report", got)
	}
	if got := values["domain"]; got != testDomain {
		t.Fatalf("domain = %#v, want %q", got, testDomain)
	}
	if got := values["issues_found"]; got != 1 {
		t.Fatalf("issues_found = %#v, want 1", got)
	}
	if got := values["issues_fixed"]; got != 0 {
		t.Fatalf("issues_fixed = %#v, want 0", got)
	}
	if body == "" {
		t.Fatal("expected report body content")
	}
}

func newTestTopic(t *testing.T) string {
	t.Helper()

	topicPath := filepath.Join(t.TempDir(), testTopicSlug)
	directories := []string{
		"raw/articles",
		"raw/bookmarks",
		"raw/github",
		"raw/youtube",
		"raw/codebase/files",
		"raw/codebase/symbols",
		"wiki/concepts",
		"wiki/index",
		"outputs/queries",
		"outputs/briefings",
		"outputs/diagrams",
		"outputs/reports",
	}

	for _, dir := range directories {
		if err := os.MkdirAll(filepath.Join(topicPath, filepath.FromSlash(dir)), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	return topicPath
}

func writeMarkdownFile(t *testing.T, topicPath, relativePath string, values map[string]any, body string) {
	t.Helper()

	document, err := frontmatter.Generate(values, body)
	if err != nil {
		t.Fatalf("Generate(%s): %v", relativePath, err)
	}

	absolutePath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("mkdir parent for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(absolutePath, []byte(document), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func sourceFrontmatter(title, sourceKind, scraped string) map[string]any {
	values := map[string]any{
		"type":   "source",
		"stage":  "raw",
		"domain": testDomain,
		"tags":   []string{testDomain, "raw", "article"},
	}
	if title != "" {
		values["title"] = title
	}
	if sourceKind != "" {
		values["source_kind"] = sourceKind
	}
	if scraped != "" {
		values["scraped"] = scraped
	}

	return values
}

func conceptFrontmatter(title, updated string, sources []string) map[string]any {
	values := map[string]any{
		"type":    "wiki",
		"stage":   "compiled",
		"domain":  testDomain,
		"tags":    []string{testDomain, "wiki", "concept"},
		"created": "2026-04-10",
		"updated": updated,
		"sources": sources,
	}
	if title != "" {
		values["title"] = title
	}

	return values
}

func indexFrontmatter(title string) map[string]any {
	return map[string]any{
		"title":   title,
		"type":    "index",
		"domain":  testDomain,
		"updated": "2026-04-11",
	}
}

func mustLint(t *testing.T, topicPath string) []models.LintIssue {
	t.Helper()

	issues, err := lint.Lint(topicPath)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}

	return issues
}

func assertHasIssue(t *testing.T, issues []models.LintIssue, want models.LintIssue) {
	t.Helper()

	for _, issue := range issues {
		if issue.Kind == want.Kind &&
			issue.Severity == want.Severity &&
			issue.FilePath == want.FilePath &&
			issue.Target == want.Target {
			return
		}
	}

	t.Fatalf("expected issue %+v in %#v", want, issues)
}

func assertKindPresent(t *testing.T, issues []models.LintIssue, want models.LintIssueKind) {
	t.Helper()

	for _, issue := range issues {
		if issue.Kind == want {
			return
		}
	}

	t.Fatalf("expected issue kind %q in %#v", want, issues)
}
