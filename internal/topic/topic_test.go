package topic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/kb/internal/frontmatter"
)

func TestNewCreatesTopicSkeletonAndTemplates(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	info, err := newWithDate(
		vaultPath,
		"rust-systems",
		"Rust Systems Programming",
		"rust",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("newWithDate returned error: %v", err)
	}

	if info.Slug != "rust-systems" {
		t.Fatalf("slug = %q, want rust-systems", info.Slug)
	}
	if info.Title != "Rust Systems Programming" {
		t.Fatalf("title = %q, want Rust Systems Programming", info.Title)
	}
	if info.Domain != "rust" {
		t.Fatalf("domain = %q, want rust", info.Domain)
	}
	if info.ArticleCount != 0 {
		t.Fatalf("article count = %d, want 0", info.ArticleCount)
	}
	if info.SourceCount != 0 {
		t.Fatalf("source count = %d, want 0", info.SourceCount)
	}
	if info.LastLogEntry != "## [2026-04-11] scaffold | rust-systems" {
		t.Fatalf("last log entry = %q", info.LastLogEntry)
	}

	topicPath := filepath.Join(vaultPath, "rust-systems")
	for _, relativePath := range []string{
		"raw/articles",
		"raw/bookmarks",
		"raw/codebase",
		"raw/codebase/files",
		"raw/codebase/symbols",
		"raw/github",
		"raw/youtube",
		"wiki/concepts",
		"wiki/index",
		"outputs/queries",
		"outputs/briefings",
		"outputs/diagrams",
		"outputs/reports",
		"bases",
	} {
		assertDirExists(t, filepath.Join(topicPath, filepath.FromSlash(relativePath)))
	}

	for _, relativePath := range []string{
		"raw/articles/.gitkeep",
		"raw/bookmarks/.gitkeep",
		"raw/codebase/files/.gitkeep",
		"raw/codebase/symbols/.gitkeep",
		"raw/github/.gitkeep",
		"raw/youtube/.gitkeep",
		"wiki/concepts/.gitkeep",
		"outputs/queries/.gitkeep",
		"outputs/briefings/.gitkeep",
		"outputs/diagrams/.gitkeep",
		"outputs/reports/.gitkeep",
		"bases/.gitkeep",
	} {
		assertFileExists(t, filepath.Join(topicPath, filepath.FromSlash(relativePath)))
	}

	dashboardValues, dashboardBody := parseFrontmatterFile(t, filepath.Join(topicPath, "wiki", "index", "Dashboard.md"))
	if got := dashboardValues["title"]; got != "Dashboard" {
		t.Fatalf("dashboard title = %#v, want Dashboard", got)
	}
	if got := dashboardValues["domain"]; got != "rust" {
		t.Fatalf("dashboard domain = %#v, want rust", got)
	}
	if got := dashboardValues["updated"]; got != "2026-04-11" {
		t.Fatalf("dashboard updated = %#v, want 2026-04-11", got)
	}
	if !strings.Contains(dashboardBody, "# Rust Systems Programming — Dashboard") {
		t.Fatalf("dashboard body missing substituted title:\n%s", dashboardBody)
	}
	if strings.Contains(dashboardBody, "TOPIC_") {
		t.Fatalf("dashboard body still contains placeholders:\n%s", dashboardBody)
	}

	conceptValues, conceptBody := parseFrontmatterFile(t, filepath.Join(topicPath, "wiki", "index", "Concept Index.md"))
	if got := conceptValues["title"]; got != "Concept Index" {
		t.Fatalf("concept index title = %#v, want Concept Index", got)
	}
	if got := conceptValues["domain"]; got != "rust" {
		t.Fatalf("concept index domain = %#v, want rust", got)
	}
	if !strings.Contains(conceptBody, "# Rust Systems Programming — Concept Index") {
		t.Fatalf("concept index body missing substituted title:\n%s", conceptBody)
	}

	sourceValues, sourceBody := parseFrontmatterFile(t, filepath.Join(topicPath, "wiki", "index", "Source Index.md"))
	if got := sourceValues["title"]; got != "Source Index" {
		t.Fatalf("source index title = %#v, want Source Index", got)
	}
	if got := sourceValues["domain"]; got != "rust" {
		t.Fatalf("source index domain = %#v, want rust", got)
	}
	if !strings.Contains(sourceBody, "# Rust Systems Programming — Source Index") {
		t.Fatalf("source index body missing substituted title:\n%s", sourceBody)
	}
}

func TestNewCreatesClaudeAndAgentsSymlink(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	_, err := newWithDate(
		vaultPath,
		"distributed-systems",
		"Distributed Systems",
		"distributed",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("newWithDate returned error: %v", err)
	}

	topicPath := filepath.Join(vaultPath, "distributed-systems")
	claudeContent := readFile(t, filepath.Join(topicPath, "CLAUDE.md"))
	if !strings.Contains(claudeContent, "# Distributed Systems") {
		t.Fatalf("CLAUDE.md missing topic title:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "**Domain:** `distributed`") {
		t.Fatalf("CLAUDE.md missing domain:\n%s", claudeContent)
	}
	if !strings.Contains(claudeContent, "collection `distributed-systems`") {
		t.Fatalf("CLAUDE.md missing slug substitution:\n%s", claudeContent)
	}

	target, err := os.Readlink(filepath.Join(topicPath, "AGENTS.md"))
	if err != nil {
		t.Fatalf("expected AGENTS.md symlink: %v", err)
	}
	if target != "CLAUDE.md" {
		t.Fatalf("AGENTS.md target = %q, want CLAUDE.md", target)
	}
}

func TestNewAppendsScaffoldEntryToLog(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	_, err := newWithDate(
		vaultPath,
		"go-runtime",
		"Go Runtime",
		"golang",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("newWithDate returned error: %v", err)
	}

	logContent := readFile(t, filepath.Join(vaultPath, "go-runtime", "log.md"))
	for _, fragment := range []string{
		"## [2026-04-11] bootstrap | topic scaffolded",
		"Topic `go-runtime` created via `new-topic.sh`. Domain: `golang`. Ready for ingest.",
		"## [2026-04-11] scaffold | go-runtime",
		"Topic `go-runtime` scaffolded via `kb topic new`. Domain: `golang`.",
	} {
		if !strings.Contains(logContent, fragment) {
			t.Fatalf("log.md missing %q:\n%s", fragment, logContent)
		}
	}
}

func TestNewReturnsErrorIfTopicExists(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	topicPath := filepath.Join(vaultPath, "existing-topic")
	if err := os.MkdirAll(topicPath, 0o755); err != nil {
		t.Fatalf("create existing topic: %v", err)
	}

	_, err := newWithDate(
		vaultPath,
		"existing-topic",
		"Existing Topic",
		"existing",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("newWithDate error = %v, want already exists", err)
	}
}

func TestNewCreatesTopicUsingExportedAPI(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	info, err := New(vaultPath, "kb-topic", "KB Topic", "kb")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if info.Slug != "kb-topic" {
		t.Fatalf("slug = %q, want kb-topic", info.Slug)
	}
	if info.Title != "KB Topic" {
		t.Fatalf("title = %q, want KB Topic", info.Title)
	}
	if info.Domain != "kb" {
		t.Fatalf("domain = %q, want kb", info.Domain)
	}
	if !strings.HasPrefix(info.LastLogEntry, "## [") || !strings.Contains(info.LastLogEntry, "scaffold | kb-topic") {
		t.Fatalf("last log entry = %q, want scaffold entry", info.LastLogEntry)
	}
}

func TestNewValidatesInputs(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "vault-file")
	if err := os.WriteFile(filePath, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write file-backed vault path: %v", err)
	}

	for _, tt := range []struct {
		name     string
		vault    string
		slug     string
		title    string
		domain   string
		contains string
	}{
		{
			name:     "empty vault path",
			vault:    "",
			slug:     "valid-topic",
			title:    "Valid Topic",
			domain:   "valid",
			contains: "vault path is required",
		},
		{
			name:     "invalid slug",
			vault:    t.TempDir(),
			slug:     "Invalid Topic",
			title:    "Valid Topic",
			domain:   "valid",
			contains: "topic slug must use lowercase alphanumerics",
		},
		{
			name:     "empty title",
			vault:    t.TempDir(),
			slug:     "valid-topic",
			title:    "",
			domain:   "valid",
			contains: "topic title is required",
		},
		{
			name:     "empty domain",
			vault:    t.TempDir(),
			slug:     "valid-topic",
			title:    "Valid Topic",
			domain:   "",
			contains: "topic domain is required",
		},
		{
			name:     "vault path is file",
			vault:    filePath,
			slug:     "valid-topic",
			title:    "Valid Topic",
			domain:   "valid",
			contains: "not a directory",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := newWithDate(
				tt.vault,
				tt.slug,
				tt.title,
				tt.domain,
				time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
			)
			if err == nil || !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("newWithDate error = %v, want substring %q", err, tt.contains)
			}
		})
	}
}

func TestListReturnsEmptySliceForVaultWithNoTopics(t *testing.T) {
	t.Parallel()

	topics, err := List(t.TempDir())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(topics) != 0 {
		t.Fatalf("topics length = %d, want 0", len(topics))
	}
}

func TestListReturnsTopicSlugsForMultipleTopics(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	for _, topicSpec := range []struct {
		slug   string
		title  string
		domain string
	}{
		{slug: "algorithms", title: "Algorithms", domain: "algorithms"},
		{slug: "operating-systems", title: "Operating Systems", domain: "systems"},
	} {
		if _, err := newWithDate(
			vaultPath,
			topicSpec.slug,
			topicSpec.title,
			topicSpec.domain,
			time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
		); err != nil {
			t.Fatalf("create topic %q: %v", topicSpec.slug, err)
		}
	}

	if err := os.MkdirAll(filepath.Join(vaultPath, "incomplete-topic", "wiki"), 0o755); err != nil {
		t.Fatalf("create incomplete topic: %v", err)
	}

	topics, err := List(vaultPath)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(topics) != 2 {
		t.Fatalf("topics length = %d, want 2", len(topics))
	}

	if topics[0].Slug != "algorithms" || topics[1].Slug != "operating-systems" {
		t.Fatalf("topic slugs = [%q %q], want [algorithms operating-systems]", topics[0].Slug, topics[1].Slug)
	}
	if topics[0].Title != "Algorithms" {
		t.Fatalf("first topic title = %q, want Algorithms", topics[0].Title)
	}
	if topics[1].Domain != "systems" {
		t.Fatalf("second topic domain = %q, want systems", topics[1].Domain)
	}
}

func TestListReturnsEmptySliceForMissingVaultPath(t *testing.T) {
	t.Parallel()

	vaultPath := filepath.Join(t.TempDir(), "missing")

	topics, err := List(vaultPath)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(topics) != 0 {
		t.Fatalf("topics length = %d, want 0", len(topics))
	}
}

func TestInfoReturnsTopicMetadataAndCounts(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if _, err := newWithDate(
		vaultPath,
		"systems-design",
		"Systems Design",
		"systems",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	topicPath := filepath.Join(vaultPath, "systems-design")
	writeFile(t, filepath.Join(topicPath, "wiki", "concepts", "Overview.md"), "# Overview\n")
	writeFile(t, filepath.Join(topicPath, "wiki", "concepts", "sub", "Patterns.md"), "# Patterns\n")
	writeFile(t, filepath.Join(topicPath, "wiki", "concepts", ".draft.md"), "# Hidden\n")

	writeFile(t, filepath.Join(topicPath, "raw", "articles", "cap.md"), "# Article\n")
	writeFile(t, filepath.Join(topicPath, "raw", "github", "repo.md"), "# Repo\n")
	writeFile(t, filepath.Join(topicPath, "raw", "bookmarks", "cluster.md"), "# Cluster\n")
	writeFile(t, filepath.Join(topicPath, "raw", "youtube", "talk.md"), "# Talk\n")
	writeFile(t, filepath.Join(topicPath, "raw", "codebase", "files", "main.md"), "# Main\n")
	writeFile(t, filepath.Join(topicPath, "raw", "github", ".cache.md"), "# Hidden\n")

	appendFile(t, filepath.Join(topicPath, "log.md"), "\n## [2026-04-12] ingest | sample\n\nImported a raw source.\n")

	info, err := Info(vaultPath, "systems-design")
	if err != nil {
		t.Fatalf("Info returned error: %v", err)
	}

	if info.Title != "Systems Design" {
		t.Fatalf("title = %q, want Systems Design", info.Title)
	}
	if info.Domain != "systems" {
		t.Fatalf("domain = %q, want systems", info.Domain)
	}
	if info.RootPath != topicPath {
		t.Fatalf("root path = %q, want %q", info.RootPath, topicPath)
	}
	if info.ArticleCount != 2 {
		t.Fatalf("article count = %d, want 2", info.ArticleCount)
	}
	if info.SourceCount != 5 {
		t.Fatalf("source count = %d, want 5", info.SourceCount)
	}
	if info.LastLogEntry != "## [2026-04-12] ingest | sample" {
		t.Fatalf("last log entry = %q, want ingest entry", info.LastLogEntry)
	}
}

func TestInfoRequiresSlug(t *testing.T) {
	t.Parallel()

	_, err := Info(t.TempDir(), "")
	if err == nil || !strings.Contains(err.Error(), "topic slug is required") {
		t.Fatalf("Info error = %v, want slug validation error", err)
	}
}

func TestInfoReturnsErrorForIncompleteTopicSkeleton(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(vaultPath, "broken-topic", "wiki", "index"), 0o755); err != nil {
		t.Fatalf("create broken topic: %v", err)
	}

	_, err := Info(vaultPath, "broken-topic")
	if err == nil || !strings.Contains(err.Error(), "missing the expected KB skeleton") {
		t.Fatalf("Info error = %v, want skeleton validation error", err)
	}
}

func TestInfoFallsBackWhenClaudeMetadataIsMissing(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if _, err := newWithDate(
		vaultPath,
		"plain-topic",
		"Plain Topic",
		"plain",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	topicPath := filepath.Join(vaultPath, "plain-topic")
	writeFile(t, filepath.Join(topicPath, "CLAUDE.md"), "schema document without explicit metadata\n")
	writeFile(t, filepath.Join(topicPath, "log.md"), "# Plain Topic - Log\n")

	info, err := Info(vaultPath, "plain-topic")
	if err != nil {
		t.Fatalf("Info returned error: %v", err)
	}

	if info.Title != "Plain Topic" {
		t.Fatalf("title = %q, want Plain Topic", info.Title)
	}
	if info.Domain != "plain-topic" {
		t.Fatalf("domain = %q, want plain-topic", info.Domain)
	}
	if info.LastLogEntry != "" {
		t.Fatalf("last log entry = %q, want empty string", info.LastLogEntry)
	}
}

func TestSubstituteValueReplacesNestedValues(t *testing.T) {
	t.Parallel()

	context := templateContext{
		Domain: "systems",
		Slug:   "systems-design",
		Title:  "Systems Design",
		Today:  "2026-04-11",
	}

	values := map[string]any{
		"title": []string{"TOPIC_TITLE", "TOPIC_DOMAIN"},
		"meta": map[string]any{
			"slug": "TOPIC_SLUG",
			"tags": []any{"TOPIC_DOMAIN", "YYYY-MM-DD"},
		},
	}

	got, ok := substituteValue(values, context).(map[string]any)
	if !ok {
		t.Fatalf("substituteValue type = %T, want map[string]any", got)
	}

	titleValues, ok := got["title"].([]string)
	if !ok {
		t.Fatalf("title type = %T, want []string", got["title"])
	}
	if titleValues[0] != "Systems Design" || titleValues[1] != "systems" {
		t.Fatalf("title values = %#v", titleValues)
	}

	metaValues, ok := got["meta"].(map[string]any)
	if !ok {
		t.Fatalf("meta type = %T, want map[string]any", got["meta"])
	}
	if metaValues["slug"] != "systems-design" {
		t.Fatalf("meta slug = %#v, want systems-design", metaValues["slug"])
	}

	tags, ok := metaValues["tags"].([]any)
	if !ok {
		t.Fatalf("tags type = %T, want []any", metaValues["tags"])
	}
	if tags[0] != "systems" || tags[1] != "2026-04-11" {
		t.Fatalf("tag values = %#v", tags)
	}
}

func TestHasTopicSkeletonRejectsPlainAgentsFile(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if _, err := newWithDate(
		vaultPath,
		"graphs",
		"Graphs",
		"graphs",
		time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
	); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	topicPath := filepath.Join(vaultPath, "graphs")
	agentsPath := filepath.Join(topicPath, "AGENTS.md")
	if err := os.Remove(agentsPath); err != nil {
		t.Fatalf("remove symlink: %v", err)
	}
	writeFile(t, agentsPath, "not a symlink")

	ok, err := hasTopicSkeleton(topicPath)
	if err != nil {
		t.Fatalf("hasTopicSkeleton returned error: %v", err)
	}
	if ok {
		t.Fatalf("hasTopicSkeleton = true, want false")
	}
}

func TestHasTopicSkeletonReturnsFalseForMissingTopic(t *testing.T) {
	t.Parallel()

	ok, err := hasTopicSkeleton(filepath.Join(t.TempDir(), "missing-topic"))
	if err != nil {
		t.Fatalf("hasTopicSkeleton returned error: %v", err)
	}
	if ok {
		t.Fatalf("hasTopicSkeleton = true, want false")
	}
}

func TestCountVisibleFilesSkipsHiddenDirectories(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "visible.md"), "# Visible\n")
	writeFile(t, filepath.Join(root, ".hidden", "hidden.md"), "# Hidden\n")

	count, err := countVisibleFiles(root)
	if err != nil {
		t.Fatalf("countVisibleFiles returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}

func TestHumanizeSlugHandlesEmptyAndHyphenatedValues(t *testing.T) {
	t.Parallel()

	if got := humanizeSlug(""); got != "Knowledge Base" {
		t.Fatalf("humanizeSlug(\"\") = %q, want Knowledge Base", got)
	}
	if got := humanizeSlug("distributed-systems"); got != "Distributed Systems" {
		t.Fatalf("humanizeSlug(distributed-systems) = %q, want Distributed Systems", got)
	}
}

func TestRenderTemplateReturnsErrorWhenAssetMissing(t *testing.T) {
	t.Parallel()

	_, err := renderTemplate("missing-template.md", templateContext{
		Domain: "systems",
		Slug:   "systems-design",
		Title:  "Systems Design",
		Today:  "2026-04-11",
	})
	if err == nil || !strings.Contains(err.Error(), "read template") {
		t.Fatalf("renderTemplate error = %v, want read template error", err)
	}
}

func TestInstallTemplatesReturnsErrorWhenTargetPathCannotBeWritten(t *testing.T) {
	t.Parallel()

	err := installTemplates(t.TempDir(), templateContext{
		Domain: "systems",
		Slug:   "systems-design",
		Title:  "Systems Design",
		Today:  "2026-04-11",
	})
	if err == nil || !strings.Contains(err.Error(), "write") {
		t.Fatalf("installTemplates error = %v, want write error", err)
	}
}

func TestAppendScaffoldEntryReturnsErrorWhenLogFileIsMissing(t *testing.T) {
	t.Parallel()

	err := appendScaffoldEntry(filepath.Join(t.TempDir(), "missing.log"), templateContext{
		Domain: "systems",
		Slug:   "systems-design",
		Title:  "Systems Design",
		Today:  "2026-04-11",
	})
	if err == nil || !strings.Contains(err.Error(), "open") {
		t.Fatalf("appendScaffoldEntry error = %v, want open error", err)
	}
}

func TestAppendScaffoldEntryReturnsErrorWhenWriteFails(t *testing.T) {
	t.Parallel()

	if _, err := os.Stat("/dev/full"); err != nil {
		t.Skip("/dev/full is not available")
	}

	err := appendScaffoldEntry("/dev/full", templateContext{
		Domain: "systems",
		Slug:   "systems-design",
		Title:  "Systems Design",
		Today:  "2026-04-11",
	})
	if err == nil || !strings.Contains(err.Error(), "write scaffold entry") {
		t.Fatalf("appendScaffoldEntry error = %v, want write error", err)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %q: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("%q is not a directory", path)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %q: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("%q is a directory, want file", path)
	}
}

func parseFrontmatterFile(t *testing.T, path string) (map[string]any, string) {
	t.Helper()

	values, body, err := frontmatter.Parse(readFile(t, path))
	if err != nil {
		t.Fatalf("parse frontmatter %q: %v", path, err)
	}

	return values, body
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}

	return string(content)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create dir for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}

func appendFile(t *testing.T, path, content string) {
	t.Helper()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open %q: %v", path, err)
	}

	if _, err := file.WriteString(content); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			t.Fatalf("append %q: %v (close error: %v)", path, err, closeErr)
		}
		t.Fatalf("append %q: %v", path, err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close %q: %v", path, err)
	}
}
