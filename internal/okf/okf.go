// Package okf implements Open Knowledge Format promotion and conformance checks.
package okf

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

const (
	OKFVersion = "0.1"
)

var (
	wikilinkPattern      = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)
	nonSentenceSpace     = regexp.MustCompile(`\s+`)
	markdownLinkPattern  = regexp.MustCompile(`\[[^\]]+\]\([^)]+\)`)
	markdownSyntaxRunes  = strings.NewReplacer("#", "", "*", "", "`", "", "_", "", ">", "", "|", " ")
	errTargetNotOKFTopic = errors.New("target topic must use mode okf")
)

// PromoteInput carries the data needed to promote one wiki document into an OKF bundle.
type PromoteInput struct {
	SourceDocPath string
	VaultPath     string
	TargetTopic   models.TopicInfo
	Type          string
	Description   string
	Types         []string
	Clock         func() time.Time
}

// ConceptResult reports the concept written by Promote.
type ConceptResult struct {
	WrittenPath     string   `json:"writtenPath"`
	Type            string   `json:"type"`
	LinksRewritten  int      `json:"linksRewritten"`
	UnresolvedLinks []string `json:"unresolvedLinks"`
	Warnings        []string `json:"warnings,omitempty"`
}

// CheckOptions configures OKF conformance checking.
type CheckOptions struct {
	Types  []string
	Strict bool
}

// Promote performs a mechanical, non-destructive wiki-to-OKF conversion.
func Promote(ctx context.Context, input PromoteInput) (ConceptResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return ConceptResult{}, fmt.Errorf("promote: %w", err)
	}
	if input.TargetTopic.Mode != models.TopicModeOKF {
		return ConceptResult{}, fmt.Errorf("promote: %w: %s", errTargetNotOKFTopic, input.TargetTopic.Slug)
	}
	conceptType := strings.TrimSpace(input.Type)
	if conceptType == "" {
		return ConceptResult{}, fmt.Errorf("promote: --type is required")
	}
	if strings.TrimSpace(input.TargetTopic.RootPath) == "" {
		return ConceptResult{}, fmt.Errorf("promote: target topic root path is required")
	}

	sourcePath, sourceTopicRoot, sourceTopicSlug, sourceRelativePath, err := resolveSourceDocument(input)
	if err != nil {
		return ConceptResult{}, fmt.Errorf("promote: %w", err)
	}
	_ = sourceTopicRoot

	sourceBytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return ConceptResult{}, fmt.Errorf("promote: read source %q: %w", sourcePath, err)
	}
	sourceValues, body, err := frontmatter.Parse(string(sourceBytes))
	if err != nil {
		return ConceptResult{}, fmt.Errorf("promote: parse source frontmatter: %w", err)
	}

	key := conceptKey(sourceRelativePath)
	writtenPath, absoluteTargetPath, err := allocateConceptPath(input.TargetTopic.RootPath, key)
	if err != nil {
		return ConceptResult{}, fmt.Errorf("promote: allocate concept path: %w", err)
	}

	warnings := typeWarnings(conceptType, input.Types)
	description, descriptionWarnings := resolveDescription(input.Description, body)
	warnings = append(warnings, descriptionWarnings...)
	title := strings.TrimSpace(frontmatter.GetString(sourceValues, "title"))
	if title == "" {
		title = vault.HumanizeSlug(key)
	}
	clock := input.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	now := clock().UTC()

	transformedBody, rewriteCount, unresolvedLinks := transformWikiLinks(body, sourceTopicSlug, input.TargetTopic.RootPath, path.Dir(writtenPath), writtenPath)
	values := map[string]any{
		"description": description,
		"timestamp":   now.Format(time.RFC3339),
		"title":       title,
		"type":        conceptType,
	}
	if tags := frontmatter.GetStringSlice(sourceValues, "tags"); len(tags) > 0 {
		values["tags"] = tags
	}

	document, err := frontmatter.Generate(values, transformedBody)
	if err != nil {
		return ConceptResult{}, fmt.Errorf("promote: generate concept frontmatter: %w", err)
	}
	if err := os.WriteFile(absoluteTargetPath, []byte(document), 0o644); err != nil {
		return ConceptResult{}, fmt.Errorf("promote: write concept %q: %w", absoluteTargetPath, err)
	}

	if err := RegenerateIndex(input.TargetTopic.RootPath); err != nil {
		return ConceptResult{}, fmt.Errorf("promote: regenerate index: %w", err)
	}
	if err := InsertLogEntry(input.TargetTopic.RootPath, now, title, writtenPath, sourceRelativePath); err != nil {
		return ConceptResult{}, fmt.Errorf("promote: update log: %w", err)
	}

	return ConceptResult{
		WrittenPath:     writtenPath,
		Type:            conceptType,
		LinksRewritten:  rewriteCount,
		UnresolvedLinks: unresolvedLinks,
		Warnings:        warnings,
	}, nil
}

// Check validates an OKF bundle and returns diagnostics.
func Check(ctx context.Context, bundlePath string, options CheckOptions) ([]models.LintIssue, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cleanBundlePath := strings.TrimSpace(bundlePath)
	if cleanBundlePath == "" {
		return nil, fmt.Errorf("okf check: bundle path is required")
	}
	info, err := os.Stat(cleanBundlePath)
	if err != nil {
		return nil, fmt.Errorf("okf check: stat bundle path %q: %w", cleanBundlePath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("okf check: bundle path must be a directory: %s", cleanBundlePath)
	}

	typeSet := normalizeTypeSet(options.Types)
	issues := make([]models.LintIssue, 0)
	err = filepath.WalkDir(cleanBundlePath, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if currentPath == cleanBundlePath {
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			return nil
		}

		relativePath, err := filepath.Rel(cleanBundlePath, currentPath)
		if err != nil {
			return fmt.Errorf("derive relative path for %q: %w", currentPath, err)
		}
		relativePath = vault.ToPosixPath(relativePath)
		switch {
		case isIndexFile(relativePath):
			issues = append(issues, checkIndexFile(cleanBundlePath, currentPath, relativePath, options.Strict)...)
		case isLogFile(relativePath):
			issues = append(issues, checkLogFile(currentPath, relativePath, options.Strict)...)
		case isExcludedMarkdown(entry.Name()):
			return nil
		default:
			issues = append(issues, checkConceptFile(currentPath, relativePath, typeSet, options.Strict)...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("okf check: %w", err)
	}

	sortIssues(issues)
	return issues, nil
}

// Columns returns the stable output column order for OKF diagnostics.
func Columns() []string {
	return []string{"severity", "kind", "filePath", "target", "message"}
}

// Rows converts issues into formatter rows.
func Rows(issues []models.LintIssue) []map[string]any {
	rows := make([]map[string]any, 0, len(issues))
	for _, issue := range issues {
		rows = append(rows, map[string]any{
			"severity": issue.Severity,
			"kind":     issue.Kind,
			"filePath": issue.FilePath,
			"target":   issue.Target,
			"message":  issue.Message,
		})
	}
	return rows
}

// HasErrors reports whether any diagnostic should fail the command.
func HasErrors(issues []models.LintIssue) bool {
	for _, issue := range issues {
		if issue.Severity == models.SeverityError {
			return true
		}
	}
	return false
}

// RegenerateIndex rewrites the root OKF index from concept frontmatter.
func RegenerateIndex(bundlePath string) error {
	concepts, okfVersion, err := loadConcepts(bundlePath)
	if err != nil {
		return err
	}
	sort.Slice(concepts, func(i, j int) bool {
		if concepts[i].Type != concepts[j].Type {
			return concepts[i].Type < concepts[j].Type
		}
		if concepts[i].Title != concepts[j].Title {
			return concepts[i].Title < concepts[j].Title
		}
		return concepts[i].RelativePath < concepts[j].RelativePath
	})

	var body strings.Builder
	body.WriteString("# OKF Bundle Index\n")
	currentType := ""
	for _, concept := range concepts {
		if concept.Type != currentType {
			currentType = concept.Type
			body.WriteString("\n## ")
			body.WriteString(currentType)
			body.WriteString("\n\n")
		}
		body.WriteString("* ")
		body.WriteString(vault.OKFLinkFormatter{}.Link("", concept.RelativePath, concept.Title))
		if concept.Description != "" {
			body.WriteString(" - ")
			body.WriteString(concept.Description)
		}
		body.WriteString("\n")
	}

	values := map[string]any{"okf_version": firstNonEmpty(okfVersion, OKFVersion)}
	content, err := frontmatter.Generate(values, body.String())
	if err != nil {
		return fmt.Errorf("generate index frontmatter: %w", err)
	}
	if err := os.WriteFile(filepath.Join(bundlePath, "index.md"), []byte(content), 0o644); err != nil {
		return fmt.Errorf("write index.md: %w", err)
	}
	return nil
}

// InsertLogEntry inserts a newest-first OKF log entry at the bundle root.
func InsertLogEntry(bundlePath string, when time.Time, title, conceptPath, sourcePath string) error {
	logPath := filepath.Join(bundlePath, "log.md")
	content, err := os.ReadFile(logPath)
	if errors.Is(err, os.ErrNotExist) {
		content = []byte("# Directory Update Log\n")
	} else if err != nil {
		return fmt.Errorf("read log.md: %w", err)
	}
	entry := fmt.Sprintf(
		"* **Creation**: Promoted %s from `%s`.",
		vault.OKFLinkFormatter{}.Link("", conceptPath, strings.TrimSpace(title)),
		sourcePath,
	)
	next := insertLogEntry(string(content), when.UTC().Format(frontmatter.DateLayout), entry)
	if err := os.WriteFile(logPath, []byte(next), 0o644); err != nil {
		return fmt.Errorf("write log.md: %w", err)
	}
	return nil
}

type conceptInfo struct {
	RelativePath string
	Type         string
	Title        string
	Description  string
}

func resolveSourceDocument(input PromoteInput) (absolutePath, topicRoot, topicSlug, topicRelativePath string, err error) {
	candidates := make([]string, 0, 3)
	source := strings.TrimSpace(input.SourceDocPath)
	if source == "" {
		return "", "", "", "", fmt.Errorf("source document path is required")
	}
	if filepath.IsAbs(source) {
		candidates = append(candidates, source)
	} else {
		candidates = append(candidates, source)
		if strings.TrimSpace(input.VaultPath) != "" {
			candidates = append(candidates, filepath.Join(input.VaultPath, filepath.FromSlash(source)))
		}
	}
	for _, candidate := range candidates {
		info, statErr := os.Stat(candidate)
		if statErr != nil {
			continue
		}
		if info.IsDir() {
			return "", "", "", "", fmt.Errorf("source document %q must be a file", candidate)
		}
		absolutePath, err := filepath.Abs(candidate)
		if err != nil {
			return "", "", "", "", fmt.Errorf("resolve source path %q: %w", candidate, err)
		}
		topicRoot, topicSlug = discoverTopicRoot(input.VaultPath, absolutePath)
		if topicRoot != "" {
			relative, relErr := filepath.Rel(topicRoot, absolutePath)
			if relErr != nil {
				return "", "", "", "", fmt.Errorf("derive source topic path: %w", relErr)
			}
			return absolutePath, topicRoot, topicSlug, vault.ToPosixPath(relative), nil
		}
		return absolutePath, "", "", vault.ToPosixPath(filepath.Base(absolutePath)), nil
	}
	return "", "", "", "", fmt.Errorf("source document not found: %s", source)
}

func discoverTopicRoot(vaultPath, sourcePath string) (string, string) {
	cleanVault := strings.TrimSpace(vaultPath)
	current := filepath.Dir(sourcePath)
	for {
		if cleanVault != "" && !vault.IsPathInside(cleanVault, current) {
			return "", ""
		}
		if isRegularFile(filepath.Join(current, "CLAUDE.md")) {
			if cleanVault == "" {
				return current, filepath.Base(current)
			}
			relative, err := filepath.Rel(cleanVault, current)
			if err != nil {
				return current, filepath.Base(current)
			}
			return current, vault.ToPosixPath(relative)
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", ""
}

func isRegularFile(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func allocateConceptPath(bundleRoot, key string) (string, string, error) {
	for suffix := 1; ; suffix++ {
		name := key
		if suffix > 1 {
			name = fmt.Sprintf("%s-%d", key, suffix)
		}
		relativePath := name + ".md"
		absolutePath := filepath.Join(bundleRoot, filepath.FromSlash(relativePath))
		if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
			return relativePath, absolutePath, nil
		} else if err != nil {
			return "", "", err
		}
	}
}

func conceptKey(documentPath string) string {
	base := path.Base(vault.StripMarkdownExtension(vault.ToPosixPath(documentPath)))
	return vault.SlugifySegment(base)
}

func typeWarnings(conceptType string, types []string) []string {
	typeSet := normalizeTypeSet(types)
	if len(typeSet) == 0 {
		return nil
	}
	if _, ok := typeSet[conceptType]; ok {
		return nil
	}
	return []string{fmt.Sprintf("type %q is outside the configured OKF vocabulary", conceptType)}
}

func normalizeTypeSet(types []string) map[string]struct{} {
	typeSet := make(map[string]struct{})
	for _, value := range types {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		typeSet[value] = struct{}{}
	}
	return typeSet
}

func resolveDescription(flagValue, body string) (string, []string) {
	if description := strings.TrimSpace(flagValue); description != "" {
		return description, nil
	}
	if sentence := firstBodySentence(body); sentence != "" {
		return sentence, nil
	}
	return "", []string{"description is empty because the source body has no sentence fallback"}
}

func firstBodySentence(body string) string {
	normalized := markdownSyntaxRunes.Replace(markdownLinkPattern.ReplaceAllString(body, " "))
	normalized = nonSentenceSpace.ReplaceAllString(strings.TrimSpace(normalized), " ")
	for index := 0; index < len(normalized); index++ {
		switch normalized[index] {
		case '.', '!', '?':
			if index == len(normalized)-1 || normalized[index+1] == ' ' {
				return strings.TrimSpace(normalized[:index+1])
			}
		}
	}
	return ""
}

func transformWikiLinks(body, sourceTopicSlug, bundleRoot, fromDir, currentPath string) (string, int, []string) {
	count := 0
	unresolved := make([]string, 0)
	formatter := vault.OKFLinkFormatter{}
	transformed := wikilinkPattern.ReplaceAllStringFunc(body, func(match string) string {
		inner := strings.TrimSuffix(strings.TrimPrefix(match, "[["), "]]")
		target, label := splitWikiLink(inner)
		targetPath := mappedTargetPath(target, sourceTopicSlug)
		if strings.TrimSpace(label) == "" {
			label = fallbackWikiLabel(target)
		}
		count++
		if stripFragment(targetPath) != stripFragment(currentPath) && !isRegularFile(filepath.Join(bundleRoot, filepath.FromSlash(stripFragment(targetPath)))) {
			unresolved = append(unresolved, targetPath)
		}
		return formatter.Link(fromDir, targetPath, label)
	})
	sort.Strings(unresolved)
	return transformed, count, dedupeStrings(unresolved)
}

func splitWikiLink(inner string) (string, string) {
	parts := strings.SplitN(inner, "|", 2)
	target := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return target, ""
	}
	return target, strings.TrimSpace(parts[1])
}

func mappedTargetPath(target, sourceTopicSlug string) string {
	target = vault.ToPosixPath(strings.TrimSpace(target))
	fragment := ""
	if before, after, found := strings.Cut(target, "#"); found {
		target = before
		fragment = "#" + after
	}
	target = strings.TrimPrefix(target, "/")
	if sourceTopicSlug != "" {
		prefix := strings.Trim(sourceTopicSlug, "/") + "/"
		target = strings.TrimPrefix(target, prefix)
	}
	return conceptKey(target) + ".md" + fragment
}

func stripFragment(target string) string {
	if before, _, found := strings.Cut(target, "#"); found {
		return before
	}
	return target
}

func fallbackWikiLabel(target string) string {
	target = stripFragment(target)
	return vault.HumanizeSlug(conceptKey(target))
}

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	deduped := values[:0]
	var previous string
	for index, value := range values {
		if index > 0 && value == previous {
			continue
		}
		deduped = append(deduped, value)
		previous = value
	}
	return append([]string(nil), deduped...)
}

func loadConcepts(bundlePath string) ([]conceptInfo, string, error) {
	okfVersion := OKFVersion
	indexPath := filepath.Join(bundlePath, "index.md")
	if content, err := os.ReadFile(indexPath); err == nil {
		values, _, parseErr := frontmatter.Parse(string(content))
		if parseErr != nil {
			return nil, "", fmt.Errorf("parse index.md: %w", parseErr)
		}
		if version := strings.TrimSpace(frontmatter.GetString(values, "okf_version")); version != "" {
			okfVersion = version
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, "", fmt.Errorf("read index.md: %w", err)
	}

	concepts := make([]conceptInfo, 0)
	err := filepath.WalkDir(bundlePath, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if currentPath == bundlePath {
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") || isExcludedMarkdown(entry.Name()) {
			return nil
		}
		relativePath, err := filepath.Rel(bundlePath, currentPath)
		if err != nil {
			return err
		}
		relativePath = vault.ToPosixPath(relativePath)
		if isIndexFile(relativePath) || isLogFile(relativePath) {
			return nil
		}
		content, err := os.ReadFile(currentPath)
		if err != nil {
			return err
		}
		values, _, err := frontmatter.Parse(string(content))
		if err != nil || len(values) == 0 {
			return nil
		}
		conceptType := strings.TrimSpace(frontmatter.GetString(values, "type"))
		title := strings.TrimSpace(frontmatter.GetString(values, "title"))
		if title == "" {
			title = vault.HumanizeSlug(conceptKey(relativePath))
		}
		concepts = append(concepts, conceptInfo{
			RelativePath: relativePath,
			Type:         firstNonEmpty(conceptType, "Concept"),
			Title:        title,
			Description:  strings.TrimSpace(frontmatter.GetString(values, "description")),
		})
		return nil
	})
	if err != nil {
		return nil, "", fmt.Errorf("scan concepts: %w", err)
	}
	return concepts, okfVersion, nil
}

func checkConceptFile(filePath, relativePath string, typeSet map[string]struct{}, strict bool) []models.LintIssue {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "", fmt.Sprintf("read concept: %v", err))}
	}
	values, _, err := frontmatter.Parse(string(content))
	if err != nil {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "", fmt.Sprintf("frontmatter is invalid: %v", err))}
	}
	if len(values) == 0 {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "type", "concept must start with YAML frontmatter")}
	}
	issues := make([]models.LintIssue, 0)
	conceptType := strings.TrimSpace(frontmatter.GetString(values, "type"))
	if conceptType == "" {
		issues = append(issues, newIssue(models.SeverityError, relativePath, "type", "concept frontmatter must include non-empty type"))
	}
	for _, field := range []string{"title", "description", "timestamp"} {
		if strings.TrimSpace(frontmatter.GetString(values, field)) == "" {
			issues = append(issues, newIssue(warningSeverity(strict), relativePath, field, fmt.Sprintf("producer field %q is missing", field)))
		}
	}
	if len(typeSet) > 0 && conceptType != "" {
		if _, ok := typeSet[conceptType]; !ok {
			issues = append(issues, newIssue(warningSeverity(strict), relativePath, "type", fmt.Sprintf("type %q is outside the configured OKF vocabulary", conceptType)))
		}
	}
	return issues
}

func checkIndexFile(bundlePath, filePath, relativePath string, strict bool) []models.LintIssue {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "", fmt.Sprintf("read index: %v", err))}
	}
	values, _, err := frontmatter.Parse(string(content))
	if err != nil {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "", fmt.Sprintf("index frontmatter is invalid: %v", err))}
	}
	if len(values) == 0 {
		return nil
	}
	if relativePath != "index.md" {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "frontmatter", "only the bundle-root index.md may declare okf_version frontmatter")}
	}
	for key := range values {
		if key != "okf_version" {
			return []models.LintIssue{newIssue(models.SeverityError, relativePath, key, "root index.md frontmatter may only contain okf_version")}
		}
	}
	_ = bundlePath
	_ = strict
	return nil
}

func checkLogFile(filePath, relativePath string, strict bool) []models.LintIssue {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []models.LintIssue{newIssue(models.SeverityError, relativePath, "", fmt.Sprintf("read log: %v", err))}
	}
	issues := make([]models.LintIssue, 0)
	for line := range strings.SplitSeq(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "## ") {
			continue
		}
		dateValue := strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
		if _, err := time.Parse(frontmatter.DateLayout, dateValue); err != nil {
			issues = append(issues, newIssue(models.SeverityError, relativePath, "date", fmt.Sprintf("log heading %q must use YYYY-MM-DD", trimmed)))
		}
	}
	_ = strict
	return issues
}

func newIssue(severity models.DiagnosticSeverity, filePath, target, message string) models.LintIssue {
	return models.LintIssue{
		Kind:     models.LintIssueKindFormat,
		Severity: severity,
		FilePath: filePath,
		Target:   target,
		Message:  message,
	}
}

func warningSeverity(strict bool) models.DiagnosticSeverity {
	if strict {
		return models.SeverityError
	}
	return models.SeverityWarning
}

func sortIssues(issues []models.LintIssue) {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Severity != issues[j].Severity {
			return issues[i].Severity < issues[j].Severity
		}
		if issues[i].FilePath != issues[j].FilePath {
			return issues[i].FilePath < issues[j].FilePath
		}
		if issues[i].Target != issues[j].Target {
			return issues[i].Target < issues[j].Target
		}
		return issues[i].Message < issues[j].Message
	})
}

func isIndexFile(relativePath string) bool {
	return strings.EqualFold(path.Base(vault.ToPosixPath(relativePath)), "index.md")
}

func isLogFile(relativePath string) bool {
	return strings.EqualFold(path.Base(vault.ToPosixPath(relativePath)), "log.md")
}

func isExcludedMarkdown(name string) bool {
	lower := strings.ToLower(strings.TrimSpace(name))
	switch lower {
	case "claude.md", "agents.md", "readme.md", "license.md", "notice.md", "attribution.md":
		return true
	default:
		return strings.HasPrefix(lower, "license.") || strings.HasPrefix(lower, "notice.") || strings.HasPrefix(lower, "attribution.")
	}
}

func insertLogEntry(content, dateHeading, entry string) string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		lines = []string{"# Directory Update Log"}
	}
	heading := "## " + dateHeading
	for index, line := range lines {
		if strings.TrimSpace(line) != heading {
			continue
		}
		next := append([]string{}, lines[:index+1]...)
		next = append(next, entry)
		next = append(next, lines[index+1:]...)
		return strings.TrimRight(strings.Join(next, "\n"), "\n") + "\n"
	}
	next := []string{strings.TrimRight(strings.Join(lines[:1], "\n"), "\n"), "", heading, "", entry}
	if len(lines) > 1 {
		rest := strings.TrimSpace(strings.Join(lines[1:], "\n"))
		if rest != "" {
			next = append(next, "", rest)
		}
	}
	return strings.TrimRight(strings.Join(next, "\n"), "\n") + "\n"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
