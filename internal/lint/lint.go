// Package lint checks knowledge base topics for structural issues such as missing frontmatter, broken wikilinks, and orphaned files.
package lint

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

var (
	fencedCodePattern      = regexp.MustCompile("(?s)```.*?```")
	inlineCodePattern      = regexp.MustCompile("`[^`\n]+`")
	leadingFrontmatterExpr = regexp.MustCompile(`(?s)^---\r?\n.*?\r?\n---\r?\n?`)
	wikilinkPattern        = regexp.MustCompile(`\[\[([^\[\]|#]+?)(?:\|[^\[\]]*?)?(?:#[^\[\]]*?)?\]\]`)
	linkTokenPattern       = regexp.MustCompile(`[\pL\pN]+`)
)

var linkStopwords = map[string]struct{}{
	"a":    {},
	"an":   {},
	"and":  {},
	"by":   {},
	"for":  {},
	"in":   {},
	"of":   {},
	"on":   {},
	"or":   {},
	"the":  {},
	"to":   {},
	"with": {},
}

var formatterColumns = []string{"severity", "kind", "filePath", "target", "message"}

const (
	javaParseErrorCode         = "JAVA_PARSE_ERROR"
	javaResolutionFallbackCode = "JAVA_RESOLUTION_FALLBACK"
)

// LintOptions configures optional lint checks.
type LintOptions struct {
	JavaGovernance JavaDiagnosticsGovernancePolicy
}

// JavaDiagnosticsGovernancePolicy configures threshold-based Java diagnostics checks.
type JavaDiagnosticsGovernancePolicy struct {
	MaxParseErrors      int
	MaxFallbackWarnings int
}

type schemaSpec struct {
	dateFields []string
	expected   map[string]string
	listFields []string
	required   []string
}

type vaultFile struct {
	absolutePath      string
	body              string
	frontmatter       map[string]any
	inTopic           bool
	links             []string
	parseErr          error
	relativePath      string
	vaultRelativePath string
}

type vaultState struct {
	domain     string
	files      []*vaultFile
	allFiles   []*vaultFile
	aliasIndex map[string][]*vaultFile
	pathIndex  map[string]*vaultFile
	stemIndex  map[string][]*vaultFile
	titleIndex map[string][]*vaultFile
	topicPath  string
	topicSlug  string
	vaultPath  string
}

// Lint walks one KB topic, validates structural issues, and returns sorted lint
// issues that can be formatted by the CLI layer.
func Lint(topicPath string) ([]models.LintIssue, error) {
	return LintWithOptions(topicPath, LintOptions{
		JavaGovernance: defaultJavaDiagnosticsGovernancePolicy(),
	})
}

// LintWithOptions walks one KB topic, validates structural issues, and applies
// optional governance policies.
func LintWithOptions(topicPath string, options LintOptions) ([]models.LintIssue, error) {
	state, err := loadVault(topicPath)
	if err != nil {
		return nil, err
	}

	issues := make([]models.LintIssue, 0)
	for _, file := range state.files {
		issues = append(issues, validateFile(file)...)
	}

	graphIssues, incoming := buildLinkGraph(state)
	issues = append(issues, graphIssues...)
	issues = append(issues, findOrphans(state, incoming)...)
	issues = append(issues, findSourceIssues(state)...)
	issues = append(issues, evaluateJavaDiagnosticsGovernance(state, options.JavaGovernance)...)

	sortIssues(issues)
	return issues, nil
}

func defaultJavaDiagnosticsGovernancePolicy() JavaDiagnosticsGovernancePolicy {
	return JavaDiagnosticsGovernancePolicy{
		MaxParseErrors:      0,
		MaxFallbackWarnings: -1,
	}
}

// Columns returns the stable column order for lint issue output.
func Columns() []string {
	return append([]string(nil), formatterColumns...)
}

// Rows converts lint issues into formatter-friendly row data.
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

// SaveReport writes a markdown lint report to outputs/reports/<date>-lint.md and
// returns the written absolute path.
func SaveReport(topicPath string, issues []models.LintIssue, now time.Time) (string, error) {
	state, err := loadVault(topicPath)
	if err != nil {
		return "", err
	}

	reportTime := now.UTC()
	if reportTime.IsZero() {
		reportTime = time.Now().UTC()
	}

	domain := state.domain
	if strings.TrimSpace(domain) == "" {
		domain = state.topicSlug
	}

	reportPath := filepath.Join(
		state.topicPath,
		filepath.FromSlash(path.Join("outputs", "reports", reportTime.Format(frontmatter.DateLayout)+"-lint.md")),
	)
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
		return "", fmt.Errorf("lint: create report directory: %w", err)
	}

	body := renderReport(state.topicSlug, issues, reportTime)
	document, err := frontmatter.Generate(map[string]any{
		"title":        fmt.Sprintf("Lint Report %s", reportTime.Format(frontmatter.DateLayout)),
		"type":         "output",
		"stage":        "lint-report",
		"domain":       domain,
		"tags":         []string{domain, "output", "lint-report"},
		"created":      reportTime.Format(frontmatter.DateLayout),
		"issues_found": len(issues),
		"issues_fixed": 0,
	}, body)
	if err != nil {
		return "", fmt.Errorf("lint: generate report frontmatter: %w", err)
	}

	if err := os.WriteFile(reportPath, []byte(document), 0o644); err != nil {
		return "", fmt.Errorf("lint: write report %q: %w", reportPath, err)
	}

	return reportPath, nil
}

func loadVault(topicPath string) (vaultState, error) {
	trimmedTopicPath := strings.TrimSpace(topicPath)
	if trimmedTopicPath == "" {
		return vaultState{}, fmt.Errorf("lint: topic path is required")
	}
	cleanTopicPath := filepath.Clean(trimmedTopicPath)

	info, err := os.Stat(cleanTopicPath)
	if err != nil {
		return vaultState{}, fmt.Errorf("lint: stat topic path %q: %w", cleanTopicPath, err)
	}
	if !info.IsDir() {
		return vaultState{}, fmt.Errorf("lint: topic path %q is not a directory", cleanTopicPath)
	}

	state := vaultState{
		files:      make([]*vaultFile, 0),
		allFiles:   make([]*vaultFile, 0),
		aliasIndex: make(map[string][]*vaultFile),
		pathIndex:  make(map[string]*vaultFile),
		stemIndex:  make(map[string][]*vaultFile),
		titleIndex: make(map[string][]*vaultFile),
		topicPath:  cleanTopicPath,
		topicSlug:  filepath.Base(cleanTopicPath),
		vaultPath:  filepath.Dir(cleanTopicPath),
	}

	topicRoots, err := discoverTopicRoots(state.vaultPath, cleanTopicPath)
	if err != nil {
		return vaultState{}, err
	}

	for _, topicRoot := range topicRoots {
		rootSlug := filepath.Base(topicRoot)
		if err := filepath.WalkDir(topicRoot, func(currentPath string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk %q: %w", currentPath, walkErr)
			}
			if entry.IsDir() {
				return nil
			}
			if strings.ToLower(filepath.Ext(entry.Name())) != ".md" {
				return nil
			}

			relativeToRoot, err := filepath.Rel(topicRoot, currentPath)
			if err != nil {
				return fmt.Errorf("resolve relative path for %q: %w", currentPath, err)
			}
			relativeToRoot = filepath.ToSlash(relativeToRoot)
			if relativeToRoot == "AGENTS.md" {
				return nil
			}

			content, err := os.ReadFile(currentPath)
			if err != nil {
				return fmt.Errorf("read %q: %w", currentPath, err)
			}

			values, _, parseErr := frontmatter.Parse(string(content))
			if parseErr != nil {
				values = map[string]any{}
			}

			body := markdownBody(string(content))
			vaultRelativePath := path.Join(rootSlug, relativeToRoot)
			file := &vaultFile{
				absolutePath:      currentPath,
				body:              body,
				frontmatter:       values,
				inTopic:           topicRoot == cleanTopicPath,
				links:             extractDocumentWikilinks(body, values),
				parseErr:          parseErr,
				relativePath:      vaultRelativePath,
				vaultRelativePath: vaultRelativePath,
			}
			state.allFiles = append(state.allFiles, file)

			if file.inTopic {
				file.relativePath = relativeToRoot
				state.files = append(state.files, file)
				if state.domain == "" {
					if domain := strings.TrimSpace(frontmatter.GetString(values, "domain")); domain != "" {
						state.domain = domain
					}
				}
			}

			return nil
		}); err != nil {
			return vaultState{}, err
		}
	}

	sort.Slice(state.files, func(i, j int) bool {
		return state.files[i].relativePath < state.files[j].relativePath
	})
	sort.Slice(state.allFiles, func(i, j int) bool {
		return state.allFiles[i].vaultRelativePath < state.allFiles[j].vaultRelativePath
	})

	for _, file := range state.allFiles {
		vaultRelativeNoExt := strings.TrimSuffix(file.vaultRelativePath, ".md")
		state.pathIndex[vaultRelativeNoExt] = file

		if file.inTopic {
			relativeNoExt := strings.TrimSuffix(file.relativePath, ".md")
			state.pathIndex[relativeNoExt] = file
			state.pathIndex[path.Join(state.topicSlug, relativeNoExt)] = file
		}

		stem := path.Base(vaultRelativeNoExt)
		if stem != "" {
			state.stemIndex[stem] = append(state.stemIndex[stem], file)
		}

		if title := strings.TrimSpace(frontmatter.GetString(file.frontmatter, "title")); title != "" {
			state.titleIndex[title] = append(state.titleIndex[title], file)
		}

		seenAliases := make(map[string]struct{})
		for _, alias := range linkAliasesForFile(file, state.topicSlug) {
			key := canonicalLinkKey(alias)
			if key == "" {
				continue
			}
			if _, exists := seenAliases[key]; exists {
				continue
			}
			seenAliases[key] = struct{}{}
			state.aliasIndex[key] = append(state.aliasIndex[key], file)
		}
	}

	return state, nil
}

func validateFile(file *vaultFile) []models.LintIssue {
	spec, ok := schemaForPath(file.relativePath)
	if !ok {
		return nil
	}

	if file.parseErr != nil {
		return []models.LintIssue{newIssue(
			models.LintIssueKindFormat,
			models.SeverityError,
			file.relativePath,
			fmt.Sprintf("frontmatter parse error: %v", file.parseErr),
			"",
		)}
	}

	issues := make([]models.LintIssue, 0)
	for _, field := range spec.required {
		value, exists := file.frontmatter[field]
		if !exists || isMissingValue(value) {
			issues = append(issues, newIssue(
				models.LintIssueKindFormat,
				models.SeverityError,
				file.relativePath,
				fmt.Sprintf("missing required frontmatter field %q", field),
				field,
			))
		}
	}

	for field, expected := range spec.expected {
		actual := strings.TrimSpace(frontmatter.GetString(file.frontmatter, field))
		if actual == "" {
			continue
		}
		if actual != expected {
			issues = append(issues, newIssue(
				models.LintIssueKindFormat,
				models.SeverityError,
				file.relativePath,
				fmt.Sprintf("frontmatter field %q must equal %q", field, expected),
				field,
			))
		}
	}

	for _, field := range spec.dateFields {
		value, exists := file.frontmatter[field]
		if !exists || isMissingValue(value) {
			continue
		}
		if frontmatter.GetTime(file.frontmatter, field).IsZero() {
			issues = append(issues, newIssue(
				models.LintIssueKindFormat,
				models.SeverityError,
				file.relativePath,
				fmt.Sprintf("frontmatter field %q must be a valid ISO date", field),
				field,
			))
		}
	}

	for _, field := range spec.listFields {
		value, exists := file.frontmatter[field]
		if !exists || isMissingValue(value) {
			continue
		}
		if frontmatter.GetStringSlice(file.frontmatter, field) == nil {
			issues = append(issues, newIssue(
				models.LintIssueKindFormat,
				models.SeverityError,
				file.relativePath,
				fmt.Sprintf("frontmatter field %q must be a list of strings", field),
				field,
			))
		}
	}

	return issues
}

func buildLinkGraph(state vaultState) ([]models.LintIssue, map[string]map[string]struct{}) {
	issues := make([]models.LintIssue, 0)
	incoming := make(map[string]map[string]struct{})
	seenDeadLinks := make(map[string]struct{})

	for _, file := range state.allFiles {
		for _, target := range file.links {
			resolved := state.resolveTarget(target, false)
			if resolved == nil {
				if !file.inTopic {
					continue
				}
				key := file.relativePath + "\x00" + target
				if _, exists := seenDeadLinks[key]; exists {
					continue
				}
				seenDeadLinks[key] = struct{}{}
				issues = append(issues, newIssue(
					models.LintIssueKindDeadLink,
					models.SeverityError,
					file.relativePath,
					"wikilink target does not exist",
					target,
				))
				continue
			}
			if !resolved.inTopic {
				continue
			}

			if incoming[resolved.relativePath] == nil {
				incoming[resolved.relativePath] = make(map[string]struct{})
			}
			incoming[resolved.relativePath][file.relativePath] = struct{}{}
		}
	}

	return issues, incoming
}

func discoverTopicRoots(vaultPath, topicPath string) ([]string, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("read vault path %q: %w", vaultPath, err)
	}

	roots := make([]string, 0, len(entries)+1)
	seen := map[string]struct{}{topicPath: {}}
	roots = append(roots, topicPath)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		candidate := filepath.Join(vaultPath, entry.Name())
		if _, exists := seen[candidate]; exists {
			continue
		}

		markerPath := filepath.Join(candidate, "CLAUDE.md")
		info, err := os.Stat(markerPath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("stat topic marker %q: %w", markerPath, err)
		}
		if info.IsDir() {
			continue
		}

		seen[candidate] = struct{}{}
		roots = append(roots, candidate)
	}

	sort.Strings(roots)
	return roots, nil
}

func findOrphans(state vaultState, incoming map[string]map[string]struct{}) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	for _, file := range state.files {
		if !isWikiConceptPath(file.relativePath) {
			continue
		}
		if len(incoming[file.relativePath]) > 0 {
			continue
		}
		issues = append(issues, newIssue(
			models.LintIssueKindOrphan,
			models.SeverityWarning,
			file.relativePath,
			"wiki article has no incoming wikilinks",
			"",
		))
	}

	return issues
}

func findSourceIssues(state vaultState) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	seenMissing := make(map[string]struct{})
	seenStale := make(map[string]struct{})

	for _, file := range state.files {
		if !isWikiConceptPath(file.relativePath) || file.parseErr != nil {
			continue
		}

		updated := frontmatter.GetTime(file.frontmatter, "updated")
		for _, sourceRef := range frontmatter.GetStringSlice(file.frontmatter, "sources") {
			normalizedRef := normalizeLinkTarget(sourceRef)
			if normalizedRef == "" || isHTTPURL(sourceRef) {
				continue
			}

			resolved := state.resolveTarget(normalizedRef, true)
			if resolved == nil {
				key := file.relativePath + "\x00" + normalizedRef
				if _, exists := seenMissing[key]; exists {
					continue
				}
				seenMissing[key] = struct{}{}
				issues = append(issues, newIssue(
					models.LintIssueKindMissingSource,
					models.SeverityError,
					file.relativePath,
					"source reference does not resolve to a raw file",
					normalizedRef,
				))
				continue
			}

			scraped := frontmatter.GetTime(resolved.frontmatter, "scraped")
			if updated.IsZero() || scraped.IsZero() || !updated.Before(scraped) {
				continue
			}

			key := file.relativePath + "\x00" + normalizedRef
			if _, exists := seenStale[key]; exists {
				continue
			}
			seenStale[key] = struct{}{}
			issues = append(issues, newIssue(
				models.LintIssueKindStale,
				models.SeverityWarning,
				file.relativePath,
				fmt.Sprintf("article updated %s before source scraped %s", updated.Format(frontmatter.DateLayout), scraped.Format(frontmatter.DateLayout)),
				normalizedRef,
			))
		}
	}

	return issues
}

type javaDiagnosticCount struct {
	Code       string
	Count      int
	MaxAllowed int
}

func evaluateJavaDiagnosticsGovernance(
	state vaultState,
	policy JavaDiagnosticsGovernancePolicy,
) []models.LintIssue {
	type aggregate struct {
		filePath      string
		parseErrors   int
		fallbackCount int
	}

	summary := aggregate{}
	for _, file := range state.files {
		if file.parseErr != nil {
			continue
		}
		if strings.TrimSpace(frontmatter.GetString(file.frontmatter, "source_kind")) != "codebase-language-index" {
			continue
		}
		if strings.TrimSpace(frontmatter.GetString(file.frontmatter, "language")) != string(models.LangJava) {
			continue
		}

		summary.filePath = file.relativePath
		summary.parseErrors += frontmatterInt(file.frontmatter, "java_parse_error_count")
		summary.fallbackCount += frontmatterInt(file.frontmatter, "java_resolution_fallback_count")
	}

	if summary.parseErrors == 0 && summary.fallbackCount == 0 {
		return nil
	}

	results := []javaDiagnosticCount{
		{
			Code:       javaParseErrorCode,
			Count:      summary.parseErrors,
			MaxAllowed: policy.MaxParseErrors,
		},
		{
			Code:       javaResolutionFallbackCode,
			Count:      summary.fallbackCount,
			MaxAllowed: policy.MaxFallbackWarnings,
		},
	}

	issues := make([]models.LintIssue, 0, len(results))
	for _, result := range results {
		if result.Count == 0 {
			continue
		}
		if result.MaxAllowed < 0 {
			continue
		}

		severity := models.SeverityWarning
		if result.MaxAllowed >= 0 && result.Count > result.MaxAllowed {
			severity = models.SeverityError
		}

		issues = append(issues, newIssue(
			models.LintIssueKindJavaDiagnosticGovernance,
			severity,
			summary.filePath,
			renderJavaGovernanceMessage(result, severity == models.SeverityError),
			result.Code,
		))
	}

	return issues
}

func renderJavaGovernanceMessage(result javaDiagnosticCount, blocking bool) string {
	status := "pass"
	if blocking {
		status = "fail"
	}

	payload, err := json.Marshal(map[string]any{
		"count":            result.Count,
		"diagnosticCode":   result.Code,
		"maxAllowed":       result.MaxAllowed,
		"status":           status,
		"thresholdEnabled": result.MaxAllowed >= 0,
	})
	if err != nil {
		return fmt.Sprintf(
			`{"diagnosticCode":"%s","count":%d,"maxAllowed":%d,"status":"%s","thresholdEnabled":%t}`,
			result.Code,
			result.Count,
			result.MaxAllowed,
			status,
			result.MaxAllowed >= 0,
		)
	}

	return string(payload)
}

func frontmatterInt(values map[string]any, key string) int {
	raw, exists := values[key]
	if !exists {
		return 0
	}

	switch typed := raw.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func schemaForPath(relativePath string) (schemaSpec, bool) {
	switch {
	case isWikiConceptPath(relativePath):
		return schemaSpec{
			dateFields: []string{"created", "updated"},
			expected: map[string]string{
				"type":  "wiki",
				"stage": "compiled",
			},
			listFields: []string{"tags", "sources"},
			required:   []string{"title", "type", "stage", "domain", "tags", "created", "updated", "sources"},
		}, true
	case strings.HasPrefix(relativePath, "raw/bookmarks/"):
		return schemaSpec{
			dateFields: []string{"created", "updated"},
			expected: map[string]string{
				"type":        "source",
				"stage":       "raw",
				"source_kind": string(models.SourceKindBookmarkCluster),
			},
			listFields: []string{"source_urls", "tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "status", "created", "updated", "source_urls", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "raw/github/"):
		return schemaSpec{
			dateFields: []string{"scraped"},
			expected: map[string]string{
				"type":        "source",
				"stage":       "raw",
				"source_kind": string(models.SourceKindGitHubREADME),
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "scraped", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "raw/youtube/"):
		return schemaSpec{
			dateFields: []string{"scraped"},
			expected: map[string]string{
				"type":        "source",
				"stage":       "raw",
				"source_kind": string(models.SourceKindYouTubeTranscript),
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "scraped", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "raw/codebase/files/"):
		return schemaSpec{
			dateFields: []string{"scraped"},
			expected: map[string]string{
				"type":        "source",
				"stage":       "raw",
				"source_kind": string(models.SourceKindCodebaseFile),
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "scraped", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "raw/codebase/symbols/"):
		return schemaSpec{
			dateFields: []string{"scraped"},
			expected: map[string]string{
				"type":        "source",
				"stage":       "raw",
				"source_kind": string(models.SourceKindCodebaseSymbol),
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "scraped", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "raw/"):
		return schemaSpec{
			dateFields: []string{"scraped"},
			expected: map[string]string{
				"type":  "source",
				"stage": "raw",
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "source_kind", "scraped", "tags"},
		}, true
	case strings.HasPrefix(relativePath, "outputs/queries/"):
		return schemaSpec{
			dateFields: []string{"created", "updated"},
			expected: map[string]string{
				"type":  "output",
				"stage": "query",
			},
			listFields: []string{"tags", "informed_by"},
			required:   []string{"title", "type", "stage", "domain", "tags", "created", "updated", "informed_by"},
		}, true
	case strings.HasPrefix(relativePath, "outputs/briefings/"):
		return schemaSpec{
			dateFields: []string{"created", "updated"},
			expected: map[string]string{
				"type":  "output",
				"stage": "briefing",
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "tags", "created", "updated"},
		}, true
	case strings.HasPrefix(relativePath, "outputs/diagrams/"):
		return schemaSpec{
			dateFields: []string{"created", "updated"},
			expected: map[string]string{
				"type":  "output",
				"stage": "diagram",
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "tags", "created", "updated"},
		}, true
	case strings.HasPrefix(relativePath, "outputs/reports/"):
		return schemaSpec{
			dateFields: []string{"created"},
			expected: map[string]string{
				"type":  "output",
				"stage": "lint-report",
			},
			listFields: []string{"tags"},
			required:   []string{"title", "type", "stage", "domain", "tags", "created", "issues_found", "issues_fixed"},
		}, true
	case strings.HasPrefix(relativePath, "wiki/index/"), strings.HasPrefix(relativePath, "wiki/codebase/index/"):
		return schemaSpec{
			dateFields: []string{"updated"},
			expected: map[string]string{
				"type": "index",
			},
			required: []string{"title", "type", "domain", "updated"},
		}, true
	default:
		return schemaSpec{}, false
	}
}

func renderReport(topicSlug string, issues []models.LintIssue, now time.Time) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "LINT REPORT - %s/ - %s\n\n", topicSlug, now.Format(frontmatter.DateLayout))
	for _, kind := range models.LintIssueKinds() {
		sectionIssues := issuesByKind(issues, kind)
		fmt.Fprintf(&builder, "%s (%d)\n", reportSectionTitle(kind), len(sectionIssues))
		if len(sectionIssues) == 0 {
			builder.WriteString("  (none)\n\n")
			continue
		}

		for _, issue := range sectionIssues {
			builder.WriteString("  - ")
			builder.WriteString(issue.FilePath)
			if issue.Target != "" {
				builder.WriteString(" -> ")
				builder.WriteString(issue.Target)
			}
			if issue.Message != "" {
				builder.WriteString(" : ")
				builder.WriteString(issue.Message)
			}
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	if len(issues) == 0 {
		builder.WriteString("OK — no issues found\n")
	}

	return builder.String()
}

func reportSectionTitle(kind models.LintIssueKind) string {
	switch kind {
	case models.LintIssueKindDeadLink:
		return "DEAD WIKILINKS"
	case models.LintIssueKindOrphan:
		return "ORPHAN ARTICLES"
	case models.LintIssueKindMissingSource:
		return "MISSING SOURCES"
	case models.LintIssueKindStale:
		return "STALE CONTENT"
	case models.LintIssueKindFormat:
		return "FORMAT VIOLATIONS"
	case models.LintIssueKindJavaDiagnosticGovernance:
		return "JAVA DIAGNOSTICS GOVERNANCE"
	default:
		return strings.ToUpper(string(kind))
	}
}

func issuesByKind(issues []models.LintIssue, kind models.LintIssueKind) []models.LintIssue {
	filtered := make([]models.LintIssue, 0)
	for _, issue := range issues {
		if issue.Kind == kind {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func markdownBody(markdown string) string {
	return leadingFrontmatterExpr.ReplaceAllString(markdown, "")
}

func extractWikilinks(text string) []string {
	clean := stripCode(text)
	matches := wikilinkPattern.FindAllStringSubmatch(clean, -1)
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		target := normalizeLinkTarget(match[1])
		if target == "" {
			continue
		}
		links = append(links, target)
	}

	return links
}

func extractDocumentWikilinks(body string, values map[string]any) []string {
	sourceKind := strings.TrimSpace(frontmatter.GetString(values, "source_kind"))
	switch sourceKind {
	case string(models.SourceKindCodebaseFile):
		return extractManagedCodebaseLinks(body, "Symbols", "Outgoing Relations", "Backlinks")
	case string(models.SourceKindCodebaseSymbol):
		parts := []string{
			extractBodyPreamble(body),
			vault.ExtractSection(body, "Outgoing Relations"),
			vault.ExtractSection(body, "Backlinks"),
		}
		return extractWikilinks(strings.Join(parts, "\n"))
	case "codebase-directory-index", "codebase-language-index":
		return extractManagedCodebaseLinks(body, "Files", "Symbols")
	default:
		return extractWikilinks(body)
	}
}

func extractManagedCodebaseLinks(body string, headings ...string) []string {
	parts := make([]string, 0, len(headings))
	for _, heading := range headings {
		section := vault.ExtractSection(body, heading)
		if strings.TrimSpace(section) == "" {
			continue
		}

		parts = append(parts, section)
	}

	return extractWikilinks(strings.Join(parts, "\n"))
}

func extractBodyPreamble(body string) string {
	normalizedBody := strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(normalizedBody, "\n")
	collected := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "## ") {
			break
		}

		collected = append(collected, line)
	}

	return strings.TrimSpace(strings.Join(collected, "\n"))
}

func stripCode(text string) string {
	text = fencedCodePattern.ReplaceAllString(text, "")
	return inlineCodePattern.ReplaceAllString(text, "")
}

func normalizeLinkTarget(target string) string {
	trimmed := strings.TrimSpace(target)
	trimmed = strings.TrimPrefix(trimmed, "[[")
	trimmed = strings.TrimSuffix(trimmed, "]]")
	if index := strings.IndexAny(trimmed, "|#"); index >= 0 {
		trimmed = trimmed[:index]
	}
	trimmed = strings.TrimSpace(trimmed)
	trimmed = strings.TrimPrefix(strings.ReplaceAll(trimmed, "\\", "/"), "./")
	trimmed = strings.Trim(trimmed, "/")
	return strings.TrimSuffix(trimmed, ".md")
}

func (state vaultState) resolveTarget(target string, rawOnly bool) *vaultFile {
	normalized := normalizeLinkTarget(target)
	if normalized == "" {
		return nil
	}

	if file := state.pathIndex[normalized]; isAcceptableTarget(file, rawOnly) {
		return file
	}
	if strings.HasPrefix(normalized, state.topicSlug+"/") {
		if file := state.pathIndex[strings.TrimPrefix(normalized, state.topicSlug+"/")]; isAcceptableTarget(file, rawOnly) {
			return file
		}
	}
	if file := pickCandidate(state.stemIndex[normalized], rawOnly); file != nil {
		return file
	}
	if file := pickCandidate(state.titleIndex[normalized], rawOnly); file != nil {
		return file
	}
	canonical := canonicalLinkKey(normalized)
	if canonical == "" {
		return nil
	}
	if file := pickCandidate(state.aliasIndex[canonical], rawOnly); file != nil {
		return file
	}
	if file := state.resolveCanonicalPrefixTarget(canonical, rawOnly); file != nil {
		return file
	}

	return nil
}

func canonicalLinkKey(target string) string {
	normalized := strings.ToLower(normalizeLinkTarget(target))
	tokens := linkTokenPattern.FindAllString(normalized, -1)
	if len(tokens) == 0 {
		return ""
	}

	filtered := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, isStopword := linkStopwords[token]; isStopword {
			continue
		}
		filtered = append(filtered, token)
	}
	if len(filtered) == 0 {
		filtered = tokens
	}

	return strings.Join(filtered, " ")
}

func linkAliasesForFile(file *vaultFile, topicSlug string) []string {
	aliases := make([]string, 0, 8)

	vaultRelativeNoExt := strings.TrimSuffix(file.vaultRelativePath, ".md")
	if vaultRelativeNoExt != "" {
		aliases = append(aliases, vaultRelativeNoExt, path.Base(vaultRelativeNoExt))
	}
	if file.inTopic {
		relativeNoExt := strings.TrimSuffix(file.relativePath, ".md")
		if relativeNoExt != "" {
			aliases = append(aliases, relativeNoExt, path.Join(topicSlug, relativeNoExt), path.Base(relativeNoExt))
		}
	}

	title := strings.TrimSpace(frontmatter.GetString(file.frontmatter, "title"))
	if title == "" {
		return aliases
	}

	aliases = append(aliases, title)
	if colon := strings.IndexRune(title, ':'); colon > 0 {
		aliases = append(aliases, strings.TrimSpace(title[:colon]))
	}

	return aliases
}

func (state vaultState) resolveCanonicalPrefixTarget(target string, rawOnly bool) *vaultFile {
	matches := make(map[*vaultFile]struct{})
	for key, candidates := range state.aliasIndex {
		if !hasCanonicalPrefixRelation(key, target) {
			continue
		}
		for _, candidate := range candidates {
			if !isAcceptableTarget(candidate, rawOnly) {
				continue
			}
			matches[candidate] = struct{}{}
		}
	}

	if len(matches) != 1 {
		return nil
	}

	for candidate := range matches {
		return candidate
	}

	return nil
}

func hasCanonicalPrefixRelation(candidate, target string) bool {
	if candidate == "" || target == "" {
		return false
	}
	if candidate == target {
		return true
	}
	return strings.HasPrefix(candidate, target+" ") || strings.HasPrefix(target, candidate+" ")
}

func pickCandidate(candidates []*vaultFile, rawOnly bool) *vaultFile {
	for _, candidate := range candidates {
		if isAcceptableTarget(candidate, rawOnly) {
			return candidate
		}
	}

	return nil
}

func isAcceptableTarget(file *vaultFile, rawOnly bool) bool {
	if file == nil {
		return false
	}
	if rawOnly && !strings.HasPrefix(file.relativePath, "raw/") {
		return false
	}

	return true
}

func isWikiConceptPath(relativePath string) bool {
	return strings.HasPrefix(relativePath, "wiki/concepts/") || strings.HasPrefix(relativePath, "wiki/codebase/concepts/")
}

func isHTTPURL(value string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
}

func isMissingValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	case []string:
		return len(typed) == 0
	case []any:
		return len(typed) == 0
	case map[string]any:
		return len(typed) == 0
	default:
		return false
	}
}

func newIssue(kind models.LintIssueKind, severity models.DiagnosticSeverity, filePath, message, target string) models.LintIssue {
	return models.LintIssue{
		Kind:     kind,
		Severity: severity,
		FilePath: filePath,
		Message:  message,
		Target:   target,
	}
}

func sortIssues(issues []models.LintIssue) {
	sort.Slice(issues, func(i, j int) bool {
		left := issues[i]
		right := issues[j]

		leftRank := severityRank(left.Severity)
		rightRank := severityRank(right.Severity)
		if leftRank != rightRank {
			return leftRank < rightRank
		}
		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		if left.Kind != right.Kind {
			return left.Kind < right.Kind
		}
		if left.Target != right.Target {
			return left.Target < right.Target
		}
		return left.Message < right.Message
	})
}

func severityRank(severity models.DiagnosticSeverity) int {
	switch severity {
	case models.SeverityError:
		return 0
	case models.SeverityWarning:
		return 1
	default:
		return 2
	}
}
