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

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"gopkg.in/yaml.v3"
)

const (
	claudeTemplatePath       = "assets/topic-claude-template.md"
	conceptIndexTemplatePath = "assets/concept-index-template.md"
	dashboardTemplatePath    = "assets/dashboard-template.md"
	logTemplatePath          = "assets/log-template.md"
	okfClaudeTemplatePath    = "assets/okf-claude-template.md"
	sourceIndexTemplatePath  = "assets/source-index-template.md"
	topicMarkerFile          = "CLAUDE.md"
	topicMetadataFileName    = "topic.yaml"
)

var (
	ErrTopicNotFound = errors.New("topic not found")

	domainPattern    = regexp.MustCompile("(?m)^\\*\\*Domain:\\*\\*\\s+`([^`]+)`")
	headingPattern   = regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`)
	topicSlugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	currentTopicDirectories = []string{
		"raw/articles",
		"raw/bookmarks",
		"raw/codebase",
		"raw/codebase/files",
		"raw/codebase/symbols",
		"raw/github",
		"raw/youtube",
		"wiki/codebase/concepts",
		"wiki/codebase/index",
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
		"wiki/codebase/concepts",
		"wiki/codebase/index",
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

	templateFS fs.FS = topicTemplateFS
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

// TopicRef is a validated vault-relative topic reference.
type TopicRef struct {
	Path     string
	Category string
	Leaf     string
}

type topicMetadataFile struct {
	Slug          string `yaml:"slug"`
	Title         string `yaml:"title"`
	Domain        string `yaml:"domain"`
	Mode          string `yaml:"mode,omitempty"`
	Category      string `yaml:"category,omitempty"`
	Path          string `yaml:"path,omitempty"`
	QMDCollection string `yaml:"qmd_collection,omitempty"`
}

// New scaffolds a new topic underneath the provided vault root.
func New(vaultPath, slug, title, domain string) (models.TopicInfo, error) {
	return NewWithMode(vaultPath, slug, title, domain, models.TopicModeWiki)
}

// NewWithMode scaffolds a new topic with an explicit lifecycle mode.
func NewWithMode(vaultPath, slug, title, domain string, mode models.TopicMode) (models.TopicInfo, error) {
	return newWithDateWithMode(vaultPath, slug, title, domain, mode, time.Now())
}

// List returns the topics discovered under the provided vault root.
func List(vaultPath string) ([]models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}

	topicSlugs, err := listTopicSlugs(cleanVaultPath)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}

	topics := make([]models.TopicInfo, 0, len(topicSlugs))
	for _, topicSlug := range topicSlugs {
		info, err := infoAtPath(filepath.Join(cleanVaultPath, filepath.FromSlash(topicSlug)), topicSlug)
		if err != nil {
			return nil, fmt.Errorf("list topics: read topic %q: %w", topicSlug, err)
		}
		topics = append(topics, info)
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Slug < topics[j].Slug
	})

	return topics, nil
}

// Resolve returns metadata for one marker-backed topic identified by path
// relative to the vault root.
func Resolve(vaultPath, topicRef string) (models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("resolve topic: %w", err)
	}

	cleanTopicRef, err := cleanTopicRef(topicRef)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("resolve topic: %w", err)
	}

	topicPath := filepath.Join(cleanVaultPath, filepath.FromSlash(cleanTopicRef))
	ok, err := hasTopicSkeleton(topicPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("resolve topic: inspect %q: %w", topicPath, err)
	}
	if !ok {
		return models.TopicInfo{}, fmt.Errorf("resolve topic: topic %q is missing %s: %w", cleanTopicRef, topicMarkerFile, ErrTopicNotFound)
	}

	info, err := infoAtPath(topicPath, cleanTopicRef)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("resolve topic: %w", err)
	}

	return info, nil
}

// Info returns topic metadata for one scaffolded topic.
func Info(vaultPath, slug string) (models.TopicInfo, error) {
	info, err := Resolve(vaultPath, slug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("topic info: %w", err)
	}

	return info, nil
}

func newWithDate(vaultPath, slug, title, domain string, now time.Time) (models.TopicInfo, error) {
	return newWithDateWithMode(vaultPath, slug, title, domain, models.TopicModeWiki, now)
}

func newWithDateWithMode(vaultPath, slug, title, domain string, mode models.TopicMode, now time.Time) (models.TopicInfo, error) {
	cleanVaultPath, err := normalizeVaultPath(vaultPath)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: %w", err)
	}

	topicRef, err := ParseTopicRef(slug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: %w", err)
	}
	cleanTitle := strings.TrimSpace(title)
	cleanDomain := strings.TrimSpace(domain)
	cleanMode, err := normalizeTopicMode(mode)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: %w", err)
	}

	switch {
	case cleanTitle == "":
		return models.TopicInfo{}, fmt.Errorf("new topic: topic title is required")
	case cleanDomain == "":
		return models.TopicInfo{}, fmt.Errorf("new topic: topic domain is required")
	}

	if err := ensureDirectory(cleanVaultPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure vault path: %w", err)
	}

	topicPath := filepath.Join(cleanVaultPath, filepath.FromSlash(topicRef.Path))
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
		Slug:   topicRef.Path,
		Title:  cleanTitle,
		Today:  now.Format(frontmatter.DateLayout),
	}

	if err := installTemplatesWithMode(topicPath, context, cleanMode); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: install templates: %w", err)
	}
	if err := writeMetadataFile(topicPath, topicMetadataForRef(topicRef, cleanTitle, cleanDomain, cleanMode)); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: write topic metadata: %w", err)
	}
	if err := ensureAgentsSymlink(topicPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure AGENTS.md symlink: %w", err)
	}
	if err := ensureGitkeeps(topicPath); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: ensure gitkeep files: %w", err)
	}
	logPath := filepath.Join(topicPath, "log.md")
	if cleanMode == models.TopicModeOKF {
		if err := appendOKFScaffoldEntry(logPath, context); err != nil {
			return models.TopicInfo{}, fmt.Errorf("new topic: append scaffold log entry: %w", err)
		}
	} else if err := appendScaffoldEntry(logPath, context); err != nil {
		return models.TopicInfo{}, fmt.Errorf("new topic: append scaffold log entry: %w", err)
	}

	info, err := infoAtPath(topicPath, topicRef.Path)
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

// ParseTopicRef validates and normalizes a topic reference relative to a vault.
func ParseTopicRef(topicRef string) (TopicRef, error) {
	trimmed := strings.TrimSpace(topicRef)
	if trimmed == "" {
		return TopicRef{}, fmt.Errorf("topic slug is required")
	}
	if filepath.IsAbs(trimmed) {
		return TopicRef{}, fmt.Errorf("topic slug must be relative: %q", topicRef)
	}

	input := filepath.ToSlash(trimmed)
	segments := strings.Split(input, "/")
	for _, segment := range segments {
		switch {
		case segment == "":
			return TopicRef{}, fmt.Errorf("topic slug cannot contain empty path segments: %q", topicRef)
		case segment == "..":
			return TopicRef{}, fmt.Errorf("topic slug cannot contain parent path segments: %q", topicRef)
		case segment == ".":
			return TopicRef{}, fmt.Errorf("topic slug cannot contain current path segments: %q", topicRef)
		case strings.HasPrefix(segment, "."):
			return TopicRef{}, fmt.Errorf("topic slug cannot reference hidden paths: %q", topicRef)
		case !topicSlugPattern.MatchString(segment):
			return TopicRef{}, fmt.Errorf(
				"topic slug must use lowercase alphanumerics separated by single hyphens: invalid segment %q in %q",
				segment,
				topicRef,
			)
		}
	}

	ref := strings.Join(segments, "/")
	leaf := segments[len(segments)-1]
	category := strings.Join(segments[:len(segments)-1], "/")
	return TopicRef{
		Path:     ref,
		Category: category,
		Leaf:     leaf,
	}, nil
}

func cleanTopicRef(topicRef string) (string, error) {
	ref, err := ParseTopicRef(topicRef)
	if err != nil {
		return "", err
	}
	return ref.Path, nil
}

func listTopicSlugs(vaultPath string) ([]string, error) {
	if _, err := os.Stat(vaultPath); errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("stat vault path %q: %w", vaultPath, err)
	}

	globs, err := topicGlobsForVault(vaultPath)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	for _, topicGlob := range globs {
		matches, err := filepath.Glob(filepath.Join(vaultPath, filepath.FromSlash(topicGlob)))
		if err != nil {
			return nil, fmt.Errorf("expand topic glob %q: %w", topicGlob, err)
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, fmt.Errorf("stat topic candidate %q: %w", match, err)
			}
			if !info.IsDir() {
				continue
			}

			relativePath, err := filepath.Rel(vaultPath, match)
			if err != nil {
				return nil, fmt.Errorf("rel topic candidate %q: %w", match, err)
			}
			slug := filepath.ToSlash(relativePath)
			if strings.HasPrefix(slug, ".") || strings.Contains(slug, "/.") {
				continue
			}

			ok, err := hasTopicSkeleton(match)
			if err != nil {
				return nil, fmt.Errorf("inspect %q: %w", match, err)
			}
			if ok {
				seen[slug] = struct{}{}
			}
		}
	}

	topics := make([]string, 0, len(seen))
	for slug := range seen {
		topics = append(topics, slug)
	}
	sort.Strings(topics)
	return topics, nil
}

func topicGlobsForVault(vaultPath string) ([]string, error) {
	configPath, found, err := config.DiscoverProjectConfigPath(vaultPath)
	if err != nil {
		return nil, err
	}
	if !found {
		return []string{"*"}, nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load project config %q: %w", configPath, err)
	}
	configVaultRoot, err := config.ResolveVaultRoot(configPath, cfg.Vault)
	if err != nil {
		return nil, err
	}
	cleanVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("resolve vault path %q: %w", vaultPath, err)
	}
	if filepath.Clean(configVaultRoot) != filepath.Clean(cleanVaultPath) {
		return []string{"*"}, nil
	}

	globs := make([]string, 0, len(cfg.Vault.TopicGlobs))
	for _, topicGlob := range cfg.Vault.TopicGlobs {
		globs = append(globs, filepath.ToSlash(filepath.Clean(filepath.FromSlash(strings.TrimSpace(topicGlob)))))
	}
	if len(globs) == 0 {
		return []string{"*"}, nil
	}
	return globs, nil
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
	for _, relativePath := range currentTopicDirectories {
		directoryPath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
		if err := os.MkdirAll(directoryPath, 0o755); err != nil {
			return fmt.Errorf("create %q: %w", directoryPath, err)
		}
	}

	return nil
}

// EnsureCurrentSkeleton creates the current topic directory structure without
// modifying topic marker files or templates.
func EnsureCurrentSkeleton(topicPath string) error {
	if strings.TrimSpace(topicPath) == "" {
		return fmt.Errorf("topic path is required")
	}

	if err := createTopicSkeleton(topicPath); err != nil {
		return err
	}
	if err := ensureTopicLog(topicPath); err != nil {
		return err
	}
	if err := ensureAgentsFile(topicPath); err != nil {
		return err
	}

	return nil
}

// RenderClaudeTemplate renders the default topic CLAUDE.md content for the
// provided topic metadata. Callers can append managed sections afterwards.
func RenderClaudeTemplate(slug, title, domain, today string) (string, error) {
	return renderTemplate(claudeTemplatePath, templateContext{
		Domain: strings.TrimSpace(domain),
		Slug:   strings.TrimSpace(slug),
		Title:  strings.TrimSpace(title),
		Today:  strings.TrimSpace(today),
	})
}

// RenderDashboardTemplate renders the default top-level Dashboard.md scaffold.
func RenderDashboardTemplate(slug, title, domain, today string) (string, error) {
	return renderTemplate(dashboardTemplatePath, templateContext{
		Domain: strings.TrimSpace(domain),
		Slug:   strings.TrimSpace(slug),
		Title:  strings.TrimSpace(title),
		Today:  strings.TrimSpace(today),
	})
}

// RenderConceptIndexTemplate renders the default top-level Concept Index scaffold.
func RenderConceptIndexTemplate(slug, title, domain, today string) (string, error) {
	return renderTemplate(conceptIndexTemplatePath, templateContext{
		Domain: strings.TrimSpace(domain),
		Slug:   strings.TrimSpace(slug),
		Title:  strings.TrimSpace(title),
		Today:  strings.TrimSpace(today),
	})
}

// RenderSourceIndexTemplate renders the default top-level Source Index scaffold.
func RenderSourceIndexTemplate(slug, title, domain, today string) (string, error) {
	return renderTemplate(sourceIndexTemplatePath, templateContext{
		Domain: strings.TrimSpace(domain),
		Slug:   strings.TrimSpace(slug),
		Title:  strings.TrimSpace(title),
		Today:  strings.TrimSpace(today),
	})
}

func installTemplates(topicPath string, context templateContext) error {
	return installTemplatesWithMode(topicPath, context, models.TopicModeWiki)
}

func installTemplatesWithMode(topicPath string, context templateContext, mode models.TopicMode) error {
	templates := topicTemplates
	if mode == models.TopicModeOKF {
		templates = []templateFile{
			{assetPath: okfClaudeTemplatePath, outputPath: "CLAUDE.md"},
		}
	}
	for _, file := range templates {
		rendered, err := renderTemplate(file.assetPath, context)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(topicPath, filepath.FromSlash(file.outputPath))
		if err := os.WriteFile(targetPath, []byte(rendered), 0o644); err != nil {
			return fmt.Errorf("write %q: %w", targetPath, err)
		}
	}

	if mode == models.TopicModeOKF {
		index, err := frontmatter.Generate(map[string]any{"okf_version": "0.1"}, "# OKF Bundle Index\n")
		if err != nil {
			return fmt.Errorf("generate OKF index: %w", err)
		}
		if err := os.WriteFile(filepath.Join(topicPath, "index.md"), []byte(index), 0o644); err != nil {
			return fmt.Errorf("write OKF index: %w", err)
		}
		if err := os.WriteFile(filepath.Join(topicPath, "log.md"), []byte("# Directory Update Log\n"), 0o644); err != nil {
			return fmt.Errorf("write OKF log: %w", err)
		}
	}

	return nil
}

// WriteMetadataFile writes the structured topic.yaml metadata file.
func WriteMetadataFile(topicPath, slug, title, domain string) error {
	topicRef, err := ParseTopicRef(slug)
	if err != nil {
		return fmt.Errorf("validate topic metadata slug: %w", err)
	}
	mode := models.TopicModeWiki
	existing, err := readTopicYAMLMetadata(filepath.Join(topicPath, topicMetadataFileName))
	if err != nil {
		return fmt.Errorf("read existing topic metadata: %w", err)
	}
	if existing.Mode != "" {
		mode, err = normalizeTopicMode(models.TopicMode(existing.Mode))
		if err != nil {
			return fmt.Errorf("validate existing topic mode: %w", err)
		}
	}
	return writeMetadataFile(topicPath, topicMetadataForRef(topicRef, title, domain, mode))
}

func topicMetadataForRef(topicRef TopicRef, title, domain string, mode models.TopicMode) topicMetadataFile {
	metadata := topicMetadataFile{
		Slug:   topicRef.Leaf,
		Title:  strings.TrimSpace(title),
		Domain: strings.TrimSpace(domain),
		Mode:   string(mode),
	}
	if topicRef.Category != "" {
		metadata.Category = topicRef.Category
		metadata.Path = topicRef.Path
		metadata.QMDCollection = topicRef.Leaf
	}
	return metadata
}

func writeMetadataFile(topicPath string, metadata topicMetadataFile) error {
	encoded, err := yaml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal topic metadata: %w", err)
	}

	metadataPath := filepath.Join(topicPath, topicMetadataFileName)
	if err := os.WriteFile(metadataPath, encoded, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", metadataPath, err)
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

func ensureAgentsFile(topicPath string) error {
	agentsPath := filepath.Join(topicPath, "AGENTS.md")
	if _, err := os.Lstat(agentsPath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("lstat %q: %w", agentsPath, err)
	}

	if err := os.Symlink("CLAUDE.md", agentsPath); err != nil {
		return fmt.Errorf("create AGENTS.md symlink: %w", err)
	}
	return nil
}

func ensureTopicLog(topicPath string) error {
	logPath := filepath.Join(topicPath, "log.md")
	if info, err := os.Stat(logPath); err == nil {
		if info.IsDir() {
			return fmt.Errorf("%q is a directory", logPath)
		}
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %q: %w", logPath, err)
	}

	if err := os.WriteFile(logPath, []byte("# Topic Log\n"), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", logPath, err)
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

func appendOKFScaffoldEntry(logPath string, context templateContext) error {
	content, err := os.ReadFile(logPath)
	if err != nil {
		return fmt.Errorf("read %q: %w", logPath, err)
	}

	entry := strings.Join([]string{
		fmt.Sprintf("## %s", context.Today),
		"",
		fmt.Sprintf("* **Initialization**: Created OKF bundle `%s` for `%s`.", context.Slug, context.Domain),
		"",
	}, "\n")

	trimmed := strings.TrimSpace(string(content))
	var next string
	if trimmed == "" {
		next = "# Directory Update Log\n\n" + entry
	} else {
		next = trimmed + "\n\n" + entry
	}
	if err := os.WriteFile(logPath, []byte(next), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", logPath, err)
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

	markerPath := filepath.Join(topicPath, topicMarkerFile)
	markerInfo, err := os.Stat(markerPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat %q: %w", markerPath, err)
	}

	return !markerInfo.IsDir(), nil
}

func infoAtPath(topicPath, slug string) (models.TopicInfo, error) {
	title, domain, mode, err := readTopicMetadata(topicPath, slug)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("read topic metadata: %w", err)
	}

	articleCount, sourceCount, err := topicCounts(topicPath, mode)
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("count topic files: %w", err)
	}

	lastLogEntry, err := readLastLogEntry(filepath.Join(topicPath, "log.md"))
	if err != nil {
		return models.TopicInfo{}, fmt.Errorf("read log entry: %w", err)
	}

	return models.TopicInfo{
		Slug:         slug,
		Title:        title,
		Domain:       domain,
		Mode:         mode,
		RootPath:     topicPath,
		ArticleCount: articleCount,
		SourceCount:  sourceCount,
		LastLogEntry: lastLogEntry,
	}, nil
}

func readTopicMetadata(claudePath, slug string) (string, string, models.TopicMode, error) {
	yamlMetadata, err := readTopicYAMLMetadata(filepath.Join(claudePath, topicMetadataFileName))
	if err != nil {
		return "", "", "", err
	}

	claudeTitle, claudeDomain, err := readClaudeMetadata(filepath.Join(claudePath, topicMarkerFile), slug)
	if err != nil {
		return "", "", "", err
	}

	title := firstNonEmpty(yamlMetadata.Title, claudeTitle, humanizeSlug(slug))
	domain := firstNonEmpty(yamlMetadata.Domain, claudeDomain, slug)
	mode, err := normalizeTopicMode(models.TopicMode(yamlMetadata.Mode))
	if err != nil {
		return "", "", "", err
	}

	return title, domain, mode, nil
}

func readTopicYAMLMetadata(metadataPath string) (topicMetadataFile, error) {
	content, err := os.ReadFile(metadataPath)
	if errors.Is(err, os.ErrNotExist) {
		return topicMetadataFile{}, nil
	}
	if err != nil {
		return topicMetadataFile{}, fmt.Errorf("read %q: %w", metadataPath, err)
	}

	var metadata topicMetadataFile
	if err := yaml.Unmarshal(content, &metadata); err != nil {
		return topicMetadataFile{}, fmt.Errorf("parse %q: %w", metadataPath, err)
	}

	metadata.Slug = strings.TrimSpace(metadata.Slug)
	metadata.Title = strings.TrimSpace(metadata.Title)
	metadata.Domain = strings.TrimSpace(metadata.Domain)
	metadata.Mode = strings.TrimSpace(metadata.Mode)
	return metadata, nil
}

func normalizeTopicMode(mode models.TopicMode) (models.TopicMode, error) {
	switch strings.TrimSpace(string(mode)) {
	case "", string(models.TopicModeWiki):
		return models.TopicModeWiki, nil
	case string(models.TopicModeOKF):
		return models.TopicModeOKF, nil
	default:
		return "", fmt.Errorf("topic mode must be wiki or okf: %q", mode)
	}
}

func topicCounts(topicPath string, mode models.TopicMode) (int, int, error) {
	if mode == models.TopicModeOKF {
		concepts, err := countOKFConceptFiles(topicPath)
		return concepts, 0, err
	}
	articleCount, err := countMarkdownFiles(filepath.Join(topicPath, "wiki", "concepts"))
	if err != nil {
		return 0, 0, fmt.Errorf("count wiki articles: %w", err)
	}
	sourceCount, err := countVisibleFiles(filepath.Join(topicPath, "raw"))
	if err != nil {
		return 0, 0, fmt.Errorf("count raw sources: %w", err)
	}
	return articleCount, sourceCount, nil
}

func countOKFConceptFiles(root string) (int, error) {
	return countFiles(root, func(entry fs.DirEntry) bool {
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") || strings.HasPrefix(name, ".") {
			return false
		}
		switch strings.ToLower(name) {
		case "index.md", "log.md", "claude.md", "agents.md", "readme.md", "license.md", "notice.md", "attribution.md":
			return false
		default:
			return true
		}
	})
}

func readClaudeMetadata(claudePath, slug string) (string, string, error) {
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func readLastLogEntry(logPath string) (string, error) {
	content, err := os.ReadFile(logPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read %q: %w", logPath, err)
	}

	last := ""
	for line := range strings.SplitSeq(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## [") || isOKFDateHeading(trimmed) {
			last = trimmed
		}
	}

	return last, nil
}

func isOKFDateHeading(line string) bool {
	if !strings.HasPrefix(line, "## ") {
		return false
	}
	value := strings.TrimSpace(strings.TrimPrefix(line, "## "))
	if len(value) != len(frontmatter.DateLayout) {
		return false
	}
	_, err := time.Parse(frontmatter.DateLayout, value)
	return err == nil
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
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
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
