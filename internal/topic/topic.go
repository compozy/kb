// Package topic scaffolds, lists, and retrieves metadata for knowledge base topics within a vault.
package topic

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	repoassets "github.com/user/kb"
	"github.com/user/kb/internal/frontmatter"
	"github.com/user/kb/internal/models"
)

const (
	claudeTemplatePath       = ".agents/skills/karpathy-kb/assets/topic-claude-template.md"
	conceptIndexTemplatePath = ".agents/skills/karpathy-kb/assets/concept-index-template.md"
	dashboardTemplatePath    = ".agents/skills/karpathy-kb/assets/dashboard-template.md"
	logTemplatePath          = ".agents/skills/karpathy-kb/assets/log-template.md"
	sourceIndexTemplatePath  = ".agents/skills/karpathy-kb/assets/source-index-template.md"
)

var (
	domainPattern    = regexp.MustCompile("(?m)^\\*\\*Domain:\\*\\*\\s+`([^`]+)`")
	headingPattern   = regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`)
	topicSlugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	requiredTopicDirectories = []string{
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
	}

	gitkeepDirectories = []string{
		"raw/articles",
		"raw/bookmarks",
		"raw/codebase/files",
		"raw/codebase/symbols",
		"raw/github",
		"raw/youtube",
		"wiki/concepts",
		"outputs/queries",
		"outputs/briefings",
		"outputs/diagrams",
		"outputs/reports",
		"bases",
	}

	topicTemplates = []templateFile{
		{assetPath: dashboardTemplatePath, outputPath: "wiki/index/Dashboard.md"},
		{assetPath: conceptIndexTemplatePath, outputPath: "wiki/index/Concept Index.md"},
		{assetPath: sourceIndexTemplatePath, outputPath: "wiki/index/Source Index.md"},
		{assetPath: claudeTemplatePath, outputPath: "CLAUDE.md"},
		{assetPath: logTemplatePath, outputPath: "log.md"},
	}

	templateFS fs.FS = repoassets.KarpathyKBTemplates
)

type templateContext struct {
	Domain string
	Slug   string
	Title  string
	Today  string
}

type templateFile struct {
	assetPath  string
	outputPath string
}

// New scaffolds a new topic underneath the provided vault root.
func New(vaultPath, slug, title, domain string) (models.TopicInfo, error) {
	return newWithDate(vaultPath, slug, title, domain, time.Now())
}

// List returns the topics discovered under the provided vault root.
func List(vaultPath string) ([]models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}

	entries, err := os.ReadDir(cleanVaultPath)
	if errors.Is(err, os.ErrNotExist) {
		return []models.TopicInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list topics: read vault path %q: %w", cleanVaultPath, err)
	}

	topics := make([]models.TopicInfo, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		topicPath := filepath.Join(cleanVaultPath, entry.Name())
		ok, err := hasTopicSkeleton(topicPath)
		if err != nil {
			return nil, fmt.Errorf("list topics: inspect %q: %w", topicPath, err)
		}
		if !ok {
			continue
		}

		info, err := infoAtPath(topicPath, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("list topics: read topic %q: %w", entry.Name(), err)
		}

		topics = append(topics, info)
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Slug < topics[j].Slug
	})

	return topics, nil
}

// Info returns topic metadata for one scaffolded topic.
func Info(vaultPath, slug string) (models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("topic info: %w", err)
	}

	cleanSlug := strings.TrimSpace(slug)
	if cleanSlug == "" {
		return models.TopicInfo{}, fmt.Errorf("topic info: topic slug is required")
	}

	topicPath := filepath.Join(cleanVaultPath, cleanSlug)
	ok, err := hasTopicSkeleton(topicPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("topic info: inspect %q: %w", topicPath, err)
	}
	if !ok {
		return models.TopicInfo{}, fmt.Errorf("topic info: topic %q is missing the expected KB skeleton", cleanSlug)
	}

	info, err := infoAtPath(topicPath, cleanSlug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("topic info: %w", err)
	}

	return info, nil
}

func newWithDate(vaultPath, slug, title, domain string, now time.Time) (models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: %w", err)
	}

	cleanSlug := strings.TrimSpace(slug)
	cleanTitle := strings.TrimSpace(title)
	cleanDomain := strings.TrimSpace(domain)

	switch {
	case cleanSlug == "":
		return models.TopicInfo{}, fmt.Errorf("new topic: topic slug is required")
	case !topicSlugPattern.MatchString(cleanSlug):
		return models.TopicInfo{}, fmt.Errorf(
			"new topic: topic slug must use lowercase alphanumerics separated by single hyphens: %q",
			cleanSlug,
		)
	case cleanTitle == "":
		return models.TopicInfo{}, fmt.Errorf("new topic: topic title is required")
	case cleanDomain == "":
		return models.TopicInfo{}, fmt.Errorf("new topic: topic domain is required")
	}

	if err := ensureDirectory(cleanVaultPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure vault path: %w", err)
	}

	topicPath := filepath.Join(cleanVaultPath, cleanSlug)
	if _, err := os.Lstat(topicPath); err == nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: topic directory %q already exists", topicPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return models.TopicInfo{}, fmt.Errorf("new topic: inspect topic path %q: %w", topicPath, err)
	}

	if err := createTopicSkeleton(topicPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: create topic skeleton: %w", err)
	}

	context := templateContext{
		Domain: cleanDomain,
		Slug:   cleanSlug,
		Title:  cleanTitle,
		Today:  now.Format(frontmatter.DateLayout),
	}

	if err := installTemplates(topicPath, context); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: install templates: %w", err)
	}
	if err := ensureAgentsSymlink(topicPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure AGENTS.md symlink: %w", err)
	}
	if err := ensureGitkeeps(topicPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure gitkeep files: %w", err)
	}
	if err := appendScaffoldEntry(filepath.Join(topicPath, "log.md"), context); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: append scaffold log entry: %w", err)
	}

	info, err := infoAtPath(topicPath, cleanSlug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: %w", err)
	}

	return info, nil
}

func normalizeVaultPath(vaultPath string) (string, error) {
	trimmed := strings.TrimSpace(vaultPath)
	if trimmed == "" {
		return "", fmt.Errorf("vault path is required")
	}

	return filepath.Clean(trimmed), nil
}

func ensureDirectory(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %q: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", path)
	}

	return nil
}

func createTopicSkeleton(topicPath string) error {
	for _, relativePath := range requiredTopicDirectories {
		directoryPath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
		if err := os.MkdirAll(directoryPath, 0o755); err != nil {
			return fmt.Errorf("create %q: %w", directoryPath, err)
		}
	}

	return nil
}

func installTemplates(topicPath string, context templateContext) error {
	for _, file := range topicTemplates {
		rendered, err := renderTemplate(file.assetPath, context)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(topicPath, filepath.FromSlash(file.outputPath))
		if err := os.WriteFile(targetPath, []byte(rendered), 0o644); err != nil {
			return fmt.Errorf("write %q: %w", targetPath, err)
		}
	}

	return nil
}

func renderTemplate(assetPath string, context templateContext) (string, error) {
	source, err := fs.ReadFile(templateFS, assetPath)
	if err != nil {
		return "", fmt.Errorf("read template %q: %w", assetPath, err)
	}

	values, body, err := frontmatter.Parse(string(source))
	if err != nil {
		return "", fmt.Errorf("parse template %q frontmatter: %w", assetPath, err)
	}

	renderedBody := replacePlaceholders(body, context)
	if len(values) == 0 {
		return renderedBody, nil
	}

	renderedValues, ok := substituteValue(values, context).(map[string]any)
	if !ok {
		return "", fmt.Errorf("render template %q frontmatter: unexpected value type %T", assetPath, values)
	}

	rendered, err := frontmatter.Generate(renderedValues, renderedBody)
	if err != nil {
		return "", fmt.Errorf("generate template %q frontmatter: %w", assetPath, err)
	}

	return rendered, nil
}

func substituteValue(value any, context templateContext) any {
	switch typed := value.(type) {
	case string:
		return replacePlaceholders(typed, context)
	case []string:
		items := make([]string, len(typed))
		for index, item := range typed {
			items[index] = replacePlaceholders(item, context)
		}
		return items
	case []any:
		items := make([]any, len(typed))
		for index, item := range typed {
			items[index] = substituteValue(item, context)
		}
		return items
	case map[string]any:
		values := make(map[string]any, len(typed))
		for key, item := range typed {
			values[key] = substituteValue(item, context)
		}
		return values
	default:
		return value
	}
}

func replacePlaceholders(value string, context templateContext) string {
	replacer := strings.NewReplacer(
		"TOPIC_DOMAIN", context.Domain,
		"TOPIC_SLUG", context.Slug,
		"TOPIC_TITLE", context.Title,
		"YYYY-MM-DD", context.Today,
	)

	return replacer.Replace(value)
}

func ensureAgentsSymlink(topicPath string) error {
	agentsPath := filepath.Join(topicPath, "AGENTS.md")
	if err := os.RemoveAll(agentsPath); err != nil {
		return fmt.Errorf("remove existing AGENTS.md: %w", err)
	}
	if err := os.Symlink("CLAUDE.md", agentsPath); err != nil {
		return fmt.Errorf("create AGENTS.md symlink: %w", err)
	}

	return nil
}

func ensureGitkeeps(topicPath string) error {
	for _, relativePath := range gitkeepDirectories {
		directoryPath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
		entries, err := os.ReadDir(directoryPath)
		if err != nil {
			return fmt.Errorf("read %q: %w", directoryPath, err)
		}
		if len(entries) > 0 {
			continue
		}

		gitkeepPath := filepath.Join(directoryPath, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, nil, 0o644); err != nil {
			return fmt.Errorf("write %q: %w", gitkeepPath, err)
		}
	}

	return nil
}

func appendScaffoldEntry(logPath string, context templateContext) error {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %q: %w", logPath, err)
	}

	entry := strings.Join([]string{
		"",
		fmt.Sprintf("## [%s] scaffold | %s", context.Today, context.Slug),
		"",
		fmt.Sprintf("Topic `%s` scaffolded via `kb topic new`. Domain: `%s`.", context.Slug, context.Domain),
		"",
	}, "\n")

	if _, err := file.WriteString(entry); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("write scaffold entry: %w (close error: %v)", err, closeErr)
		}
		return fmt.Errorf("write scaffold entry: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %q: %w", logPath, err)
	}

	return nil
}

func hasTopicSkeleton(topicPath string) (bool, error) {
	info, err := os.Stat(topicPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat %q: %w", topicPath, err)
	}
	if !info.IsDir() {
		return false, nil
	}

	for _, relativePath := range requiredTopicDirectories {
		requiredPath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
		requiredInfo, err := os.Stat(requiredPath)
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("stat %q: %w", requiredPath, err)
		}
		if !requiredInfo.IsDir() {
			return false, nil
		}
	}

	for _, relativePath := range []string{"CLAUDE.md", "log.md"} {
		requiredPath := filepath.Join(topicPath, relativePath)
		requiredInfo, err := os.Stat(requiredPath)
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("stat %q: %w", requiredPath, err)
		}
		if requiredInfo.IsDir() {
			return false, nil
		}
	}

	agentsPath := filepath.Join(topicPath, "AGENTS.md")
	linkInfo, err := os.Lstat(agentsPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("lstat %q: %w", agentsPath, err)
	}
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}

	target, err := os.Readlink(agentsPath)
	if err != nil {
		return false, fmt.Errorf("readlink %q: %w", agentsPath, err)
	}
	if target != "CLAUDE.md" {
		return false, nil
	}

	return true, nil
}

func infoAtPath(topicPath, slug string) (models.TopicInfo, error) {
	title, domain, err := readTopicMetadata(filepath.Join(topicPath, "CLAUDE.md"), slug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("read topic metadata: %w", err)
	}

	articleCount, err := countMarkdownFiles(filepath.Join(topicPath, "wiki", "concepts"))
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("count wiki articles: %w", err)
	}

	sourceCount, err := countVisibleFiles(filepath.Join(topicPath, "raw"))
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("count raw sources: %w", err)
	}

	lastLogEntry, err := readLastLogEntry(filepath.Join(topicPath, "log.md"))
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("read log entry: %w", err)
	}

	return models.TopicInfo{
		Slug:         slug,
		Title:        title,
		Domain:       domain,
		RootPath:     topicPath,
		ArticleCount: articleCount,
		SourceCount:  sourceCount,
		LastLogEntry: lastLogEntry,
	}, nil
}

func readTopicMetadata(claudePath, slug string) (string, string, error) {
	content, err := os.ReadFile(claudePath)
	if err != nil {
		return "", "", fmt.Errorf("read %q: %w", claudePath, err)
	}

	title := humanizeSlug(slug)
	if match := headingPattern.FindStringSubmatch(string(content)); match != nil {
		title = strings.TrimSpace(match[1])
	}

	domain := slug
	if match := domainPattern.FindStringSubmatch(string(content)); match != nil {
		domain = strings.TrimSpace(match[1])
	}

	return title, domain, nil
}

func readLastLogEntry(logPath string) (string, error) {
	content, err := os.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", logPath, err)
	}

	last := ""
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## [") {
			last = trimmed
		}
	}

	return last, nil
}

func countMarkdownFiles(root string) (int, error) {
	return countFiles(root, func(entry fs.DirEntry) bool {
		name := entry.Name()
		return strings.HasSuffix(name, ".md") && !strings.HasPrefix(name, ".")
	})
}

func countVisibleFiles(root string) (int, error) {
	return countFiles(root, func(entry fs.DirEntry) bool {
		return !strings.HasPrefix(entry.Name(), ".")
	})
}

func countFiles(root string, include func(fs.DirEntry) bool) (int, error) {
	count := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if include(entry) {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk %q: %w", root, err)
	}

	return count, nil
}

func humanizeSlug(slug string) string {
	if slug == "" {
		return "Knowledge Base"
	}

	parts := strings.Split(slug, "-")
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}

		words = append(words, strings.ToUpper(part[:1])+part[1:])
	}
	if len(words) == 0 {
		return "Knowledge Base"
	}

	return strings.Join(words, " ")
}
