package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/topic"
)

const fileWriteBatchSize = 64

const (
	codebaseManagedBlockStart = "<!-- kb:codebase:start -->"
	codebaseManagedBlockEnd   = "<!-- kb:codebase:end -->"
	codebaseIndexBlockStart   = "<!-- kb:codebase-index:start -->"
	codebaseIndexBlockEnd     = "<!-- kb:codebase-index:end -->"
)

var (
	frontmatterBlockPattern = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---(?:\r?\n|$)`)
)

// WriteVaultOptions bundles the rendered vault payload written to disk.
type WriteVaultOptions struct {
	BaseFiles []models.BaseFile
	Documents []models.RenderedDocument
	Graph     models.GraphSnapshot
	Progress  func(WriteProgress)
	Topic     models.TopicMetadata
}

// WriteVaultResult reports how many managed markdown documents were written.
type WriteVaultResult struct {
	RawDocumentsWritten   int
	WikiDocumentsWritten  int
	IndexDocumentsWritten int
}

// WriteProgress reports one successful persisted file within the write stage.
type WriteProgress struct {
	Completed int
	Path      string
	Total     int
}

type fileWriteRequest struct {
	Body       string
	TargetPath string
}

type topicIndexBridgeSpec struct {
	LegacyMarkers  []string
	RenderTemplate func(slug, title, domain, today string) (string, error)
	Title          string
}

// WriteVault persists the rendered markdown and base files for a topic.
func WriteVault(ctx context.Context, options WriteVaultOptions) (WriteVaultResult, error) {
	if ctx == nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: nil context")
	}
	if err := validateTopic(options.Topic); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: %w", err)
	}

	if err := os.MkdirAll(options.Topic.VaultPath, 0o755); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: create vault path %q: %w", options.Topic.VaultPath, err)
	}
	if err := ensureTopicSkeleton(options.Topic.TopicPath); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: ensure topic skeleton: %w", err)
	}
	if err := resetManagedSubtrees(options.Topic.TopicPath); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: reset managed subtrees: %w", err)
	}

	renderedFiles, err := buildWriteRequests(options)
	if err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: %w", err)
	}
	if err := ensureDirectories(renderedFiles); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: ensure document directories: %w", err)
	}
	bridgeSpecs := topicIndexBridgeSpecs()
	progressReporter := newWriteProgressReporter(options.Progress, len(renderedFiles)+len(bridgeSpecs)+2)
	if err := writeFilesInBatches(ctx, renderedFiles, progressReporter.Report); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: persist rendered files: %w", err)
	}

	claudePath := filepath.Join(options.Topic.TopicPath, "CLAUDE.md")
	if err := writeTopicClaude(claudePath, options.Topic, options.Graph, options.Documents); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: write topic manifest: %w", err)
	}
	progressReporter.Report(claudePath)
	indexBridgePaths, err := writeTopicIndexBridges(options.Topic, bridgeSpecs)
	if err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: write topic index bridges: %w", err)
	}
	for _, bridgePath := range indexBridgePaths {
		progressReporter.Report(bridgePath)
	}
	if err := ensureAgentsSymlink(options.Topic.TopicPath); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: ensure topic agents symlink: %w", err)
	}
	if err := ensureTopicGitkeeps(options.Topic.TopicPath); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: ensure gitkeep files: %w", err)
	}

	counts := countWrittenDocuments(options.Documents)
	if err := appendLog(options.Topic, options.Graph, counts); err != nil {
		return WriteVaultResult{}, fmt.Errorf("write vault: append audit log: %w", err)
	}
	progressReporter.Report(filepath.Join(options.Topic.TopicPath, "log.md"))

	return counts, nil
}

func validateTopic(topic models.TopicMetadata) error {
	switch {
	case strings.TrimSpace(topic.VaultPath) == "":
		return fmt.Errorf("topic vault path is required")
	case strings.TrimSpace(topic.TopicPath) == "":
		return fmt.Errorf("topic path is required")
	case !IsPathInside(topic.VaultPath, topic.TopicPath):
		return fmt.Errorf("topic path %q must be inside vault path %q", topic.TopicPath, topic.VaultPath)
	default:
		return nil
	}
}

func ensureTopicSkeleton(topicPath string) error {
	return topic.EnsureCurrentSkeleton(topicPath)
}

func resetManagedSubtrees(topicPath string) error {
	for _, relativePath := range []string{
		filepath.Join("raw", "codebase"),
		filepath.Join("wiki", "codebase"),
	} {
		absolutePath := filepath.Join(topicPath, relativePath)
		if err := os.RemoveAll(absolutePath); err != nil {
			return fmt.Errorf("remove %q: %w", absolutePath, err)
		}
		if err := os.MkdirAll(absolutePath, 0o755); err != nil {
			return fmt.Errorf("recreate %q: %w", absolutePath, err)
		}
	}

	if err := removeLegacyManagedConcepts(topicPath); err != nil {
		return err
	}

	return topic.EnsureCurrentSkeleton(topicPath)
}

func removeLegacyManagedConcepts(topicPath string) error {
	conceptsPath := filepath.Join(topicPath, "wiki", "concepts")
	entries, err := os.ReadDir(conceptsPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return fmt.Errorf("read %q: %w", conceptsPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if !strings.HasPrefix(entry.Name(), legacyWikiConceptPrefix) {
			continue
		}

		targetPath := filepath.Join(conceptsPath, entry.Name())
		content, err := os.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("read legacy concept %q: %w", targetPath, err)
		}
		if !isKodebaseManagedMarkdown(string(content)) {
			continue
		}
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("remove legacy managed concept %q: %w", targetPath, err)
		}
	}

	return nil
}

func isKodebaseManagedMarkdown(markdown string) bool {
	values, _, err := frontmatter.Parse(markdown)
	if err != nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(frontmatter.GetString(values, "generator")), "kodebase")
}

func buildWriteRequests(options WriteVaultOptions) ([]fileWriteRequest, error) {
	requests := make([]fileWriteRequest, 0, len(options.Documents)+len(options.BaseFiles))

	for _, document := range options.Documents {
		cleanRelativePath, err := validateRenderedDocument(document)
		if err != nil {
			return nil, err
		}

		requests = append(requests, fileWriteRequest{
			Body:       document.Body,
			TargetPath: filepath.Join(options.Topic.TopicPath, filepath.FromSlash(cleanRelativePath)),
		})
	}

	for _, baseFile := range options.BaseFiles {
		cleanRelativePath, err := validateBaseFile(baseFile)
		if err != nil {
			return nil, err
		}

		requests = append(requests, fileWriteRequest{
			Body:       RenderBaseDefinition(baseFile.Definition),
			TargetPath: filepath.Join(options.Topic.TopicPath, filepath.FromSlash(cleanRelativePath)),
		})
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].TargetPath < requests[j].TargetPath
	})

	return requests, nil
}

func validateRenderedDocument(document models.RenderedDocument) (string, error) {
	cleanRelativePath, err := cleanTopicRelativePath(document.RelativePath)
	if err != nil {
		return "", fmt.Errorf("document %q has invalid relative path: %w", document.RelativePath, err)
	}

	expectedArea, expectedPrefix, err := expectedDocumentPlacement(document.Kind)
	if err != nil {
		return "", fmt.Errorf("document %q: %w", cleanRelativePath, err)
	}
	if document.ManagedArea != expectedArea {
		return "", fmt.Errorf(
			"document %q kind %q must use managed area %q, got %q",
			cleanRelativePath,
			document.Kind,
			expectedArea,
			document.ManagedArea,
		)
	}
	if !strings.HasPrefix(cleanRelativePath, expectedPrefix) {
		return "", fmt.Errorf(
			"document %q managed area %q must live under %q",
			cleanRelativePath,
			document.ManagedArea,
			expectedPrefix,
		)
	}
	if strings.TrimSpace(document.Body) == "" {
		return "", fmt.Errorf("document %q has empty body", cleanRelativePath)
	}

	match := frontmatterBlockPattern.FindStringSubmatch(document.Body)
	if match == nil {
		return "", fmt.Errorf("document %q is missing YAML frontmatter", cleanRelativePath)
	}
	if strings.TrimSpace(document.Body[len(match[0]):]) == "" {
		return "", fmt.Errorf("document %q is missing markdown body content", cleanRelativePath)
	}

	return cleanRelativePath, nil
}

func validateBaseFile(baseFile models.BaseFile) (string, error) {
	cleanRelativePath, err := cleanTopicRelativePath(baseFile.RelativePath)
	if err != nil {
		return "", fmt.Errorf("base file %q has invalid relative path: %w", baseFile.RelativePath, err)
	}
	if !strings.HasPrefix(cleanRelativePath, "bases/") {
		return "", fmt.Errorf("base file %q must live under \"bases/\"", cleanRelativePath)
	}
	if strings.TrimSpace(RenderBaseDefinition(baseFile.Definition)) == "" {
		return "", fmt.Errorf("base file %q rendered to an empty body", cleanRelativePath)
	}

	return cleanRelativePath, nil
}

func expectedDocumentPlacement(kind models.DocumentKind) (models.ManagedArea, string, error) {
	switch kind {
	case models.DocRaw:
		return models.AreaRawCodebase, "raw/codebase/", nil
	case models.DocWiki:
		return models.AreaWikiConcept, "wiki/codebase/concepts/", nil
	case models.DocIndex:
		return models.AreaWikiIndex, "wiki/codebase/index/", nil
	default:
		return "", "", fmt.Errorf("unsupported document kind %q", kind)
	}
}

func cleanTopicRelativePath(relativePath string) (string, error) {
	cleaned := path.Clean(ToPosixPath(relativePath))

	switch {
	case cleaned == "." || cleaned == "":
		return "", fmt.Errorf("path must not be empty")
	case strings.HasPrefix(cleaned, "/"):
		return "", fmt.Errorf("path must be relative")
	case cleaned == ".." || strings.HasPrefix(cleaned, "../"):
		return "", fmt.Errorf("path must not escape the topic root")
	default:
		return cleaned, nil
	}
}

func ensureDirectories(files []fileWriteRequest) error {
	directories := make(map[string]struct{}, len(files))
	for _, file := range files {
		directories[filepath.Dir(file.TargetPath)] = struct{}{}
	}

	sortedDirectories := make([]string, 0, len(directories))
	for directoryPath := range directories {
		sortedDirectories = append(sortedDirectories, directoryPath)
	}
	sort.Strings(sortedDirectories)

	for _, directoryPath := range sortedDirectories {
		if err := os.MkdirAll(directoryPath, 0o755); err != nil {
			return fmt.Errorf("create %q: %w", directoryPath, err)
		}
	}

	return nil
}

func writeFilesInBatches(ctx context.Context, files []fileWriteRequest, report func(string)) error {
	for index := 0; index < len(files); index += fileWriteBatchSize {
		if err := ctx.Err(); err != nil {
			return err
		}

		end := index + fileWriteBatchSize
		if end > len(files) {
			end = len(files)
		}

		for _, file := range files[index:end] {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := writeTextFile(file.TargetPath, file.Body); err != nil {
				return fmt.Errorf("write %q: %w", file.TargetPath, err)
			}
			if report != nil {
				report(file.TargetPath)
			}
		}
	}

	return nil
}

type writeProgressReporter struct {
	completed int
	report    func(WriteProgress)
	total     int
}

func newWriteProgressReporter(report func(WriteProgress), total int) *writeProgressReporter {
	return &writeProgressReporter{
		report: report,
		total:  total,
	}
}

func (r *writeProgressReporter) Report(path string) {
	if r == nil || r.report == nil {
		return
	}

	r.completed++
	r.report(WriteProgress{
		Completed: r.completed,
		Path:      path,
		Total:     r.total,
	})
}

func buildTopicClaudeManagedBlock(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	documents []models.RenderedDocument,
) string {
	conceptDocuments := make([]string, 0, len(documents))
	rawDocumentCount := 0

	for _, document := range documents {
		if document.ManagedArea == models.AreaRawCodebase {
			rawDocumentCount++
		}
		if document.ManagedArea != models.AreaWikiConcept {
			continue
		}

		conceptDocuments = append(conceptDocuments, renderedConceptTitle(document))
	}

	sort.Strings(conceptDocuments)

	rootLabel := topic.Slug
	if strings.TrimSpace(topic.RootPath) != "" {
		rootLabel = filepath.Base(topic.RootPath)
	}

	lines := []string{
		codebaseManagedBlockStart,
		"## Codebase Corpus (managed)",
		"",
		fmt.Sprintf(
			"Generated by `kb ingest codebase` from `%s`. Re-running refreshes `raw/codebase/`, `wiki/codebase/`, and the managed bridge sections in `wiki/index/`; manual content outside managed sections is preserved.",
			rootLabel,
		),
		"",
		fmt.Sprintf("- Parsed files: %d", len(graph.Files)),
		fmt.Sprintf("- Parsed symbols: %d", len(graph.Symbols)),
		fmt.Sprintf("- Raw codebase notes: %d", rawDocumentCount),
		"- `raw/codebase/files/` - file-level snapshots generated from source files",
		"- `raw/codebase/symbols/` - symbol-level snapshots generated from extracted declarations",
		"- `raw/codebase/indexes/` - generated directory and language inventories",
		"- `bases/` - generated Obsidian Bases views over the raw codebase notes",
		"",
		"### Generated codebase navigation",
		"",
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseDashboardTitle), CodebaseDashboardTitle),
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseConceptIndexTitle), CodebaseConceptIndexTitle),
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseSourceIndexTitle), CodebaseSourceIndexTitle),
		"",
		"### Generated codebase articles",
		"",
	}
	if len(conceptDocuments) == 0 {
		lines = append(lines, "_None generated._")
	} else {
		for _, articleTitle := range conceptDocuments {
			lines = append(lines, "- "+ToTopicWikiLink(topic.Slug, GetWikiConceptPath(articleTitle), articleTitle))
		}
	}
	lines = append(lines, "", codebaseManagedBlockEnd)

	return strings.Join(lines, "\n")
}

func buildTopicIndexBridgeManagedBlock(topic models.TopicMetadata) string {
	lines := []string{
		codebaseIndexBlockStart,
		"## Generated codebase analysis (managed)",
		"",
		"Latest `kb ingest codebase` output is published under `wiki/codebase/` so the top-level topic indexes can stay manually curated.",
		"",
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseDashboardTitle), CodebaseDashboardTitle),
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseConceptIndexTitle), CodebaseConceptIndexTitle),
		"- " + ToTopicWikiLink(topic.Slug, GetWikiIndexPath(CodebaseSourceIndexTitle), CodebaseSourceIndexTitle),
		"",
		codebaseIndexBlockEnd,
	}

	return strings.Join(lines, "\n")
}

func writeTopicClaude(
	claudePath string,
	topicInfo models.TopicMetadata,
	graph models.GraphSnapshot,
	documents []models.RenderedDocument,
) error {
	content, err := readOrCreateClaude(claudePath, topicInfo)
	if err != nil {
		return err
	}

	managedBlock := buildTopicClaudeManagedBlock(topicInfo, graph, documents)
	return writeTextFile(
		claudePath,
		upsertManagedSection(content, managedBlock, codebaseManagedBlockStart, codebaseManagedBlockEnd),
	)
}

func readOrCreateClaude(claudePath string, topicInfo models.TopicMetadata) (string, error) {
	content, err := os.ReadFile(claudePath)
	switch {
	case err == nil:
		markdown := string(content)
		if isLegacyGeneratedClaude(markdown) {
			return renderDefaultClaude(topicInfo)
		}
		return markdown, nil
	case !errors.Is(err, os.ErrNotExist):
		return "", fmt.Errorf("read %q: %w", claudePath, err)
	}

	return renderDefaultClaude(topicInfo)
}

func renderDefaultClaude(topicInfo models.TopicMetadata) (string, error) {
	rendered, err := topic.RenderClaudeTemplate(topicInfo.Slug, topicInfo.Title, topicInfo.Domain, topicInfo.Today)
	if err != nil {
		return "", fmt.Errorf("render default CLAUDE.md: %w", err)
	}

	return rendered, nil
}

func isLegacyGeneratedClaude(markdown string) bool {
	if strings.Contains(markdown, codebaseManagedBlockStart) {
		return false
	}

	for _, marker := range []string{
		"**Topic scope:** Generated codebase knowledge topic for `",
		"The `kodebase` CLI manages `raw/codebase/`, `wiki/index/`, and wiki concept pages with `generator: kodebase` frontmatter.",
		"## Managed starter wiki",
	} {
		if !strings.Contains(markdown, marker) {
			return false
		}
	}

	return true
}

func upsertManagedSection(content, managedBlock, startMarker, endMarker string) string {
	startIndex := strings.Index(content, startMarker)
	endIndex := strings.Index(content, endMarker)
	if startIndex >= 0 && endIndex > startIndex {
		endIndex += len(endMarker)
		before := strings.TrimRight(content[:startIndex], "\n")
		after := strings.TrimLeft(content[endIndex:], "\n")
		sections := make([]string, 0, 3)
		if before != "" {
			sections = append(sections, before)
		}
		sections = append(sections, managedBlock)
		if after != "" {
			sections = append(sections, after)
		}
		return strings.Join(sections, "\n\n") + "\n"
	}

	trimmed := strings.TrimRight(content, "\n")
	if trimmed == "" {
		return managedBlock + "\n"
	}

	return trimmed + "\n\n" + managedBlock + "\n"
}

func topicIndexBridgeSpecs() []topicIndexBridgeSpec {
	return []topicIndexBridgeSpec{
		{
			Title:          TopicDashboardTitle,
			RenderTemplate: topic.RenderDashboardTemplate,
			LegacyMarkers: []string{
				"Landing page for the generated Karpathy-compatible codebase topic.",
				"## Generated codebase articles",
			},
		},
		{
			Title:          TopicConceptIndexTitle,
			RenderTemplate: topic.RenderConceptIndexTemplate,
			LegacyMarkers: []string{
				"Alphabetical listing of every generated codebase wiki article in this topic.",
				"| Article | Summary |",
			},
		},
		{
			Title:          TopicSourceIndexTitle,
			RenderTemplate: topic.RenderSourceIndexTemplate,
			LegacyMarkers: []string{
				"Raw codebase snapshots currently cited by the generated codebase wiki.",
				"| Source | Cited by |",
			},
		},
	}
}

func writeTopicIndexBridges(topicInfo models.TopicMetadata, specs []topicIndexBridgeSpec) ([]string, error) {
	managedBlock := buildTopicIndexBridgeManagedBlock(topicInfo)
	writtenPaths := make([]string, 0, len(specs))

	for _, spec := range specs {
		targetPath := filepath.Join(topicInfo.TopicPath, filepath.FromSlash(GetTopicIndexPath(spec.Title)))
		content, err := readOrCreateTopicIndex(targetPath, topicInfo, spec)
		if err != nil {
			return nil, err
		}

		updated := upsertManagedSection(content, managedBlock, codebaseIndexBlockStart, codebaseIndexBlockEnd)
		if err := writeTextFile(targetPath, updated); err != nil {
			return nil, fmt.Errorf("write %q: %w", targetPath, err)
		}

		writtenPaths = append(writtenPaths, targetPath)
	}

	return writtenPaths, nil
}

func readOrCreateTopicIndex(targetPath string, topicInfo models.TopicMetadata, spec topicIndexBridgeSpec) (string, error) {
	content, err := os.ReadFile(targetPath)
	switch {
	case err == nil:
		markdown := string(content)
		if isLegacyGeneratedTopicIndex(markdown, spec.Title, spec.LegacyMarkers) {
			return spec.RenderTemplate(topicInfo.Slug, topicInfo.Title, topicInfo.Domain, topicInfo.Today)
		}
		return markdown, nil
	case !errors.Is(err, os.ErrNotExist):
		return "", fmt.Errorf("read %q: %w", targetPath, err)
	default:
		return spec.RenderTemplate(topicInfo.Slug, topicInfo.Title, topicInfo.Domain, topicInfo.Today)
	}
}

func isLegacyGeneratedTopicIndex(markdown, expectedTitle string, legacyMarkers []string) bool {
	values, body, err := frontmatter.Parse(markdown)
	if err != nil {
		return false
	}
	if frontmatter.GetString(values, "title") != expectedTitle {
		return false
	}
	for _, marker := range legacyMarkers {
		if !strings.Contains(body, marker) {
			return false
		}
	}

	return true
}

func renderedConceptTitle(document models.RenderedDocument) string {
	if title, ok := document.Frontmatter["title"].(string); ok && strings.TrimSpace(title) != "" {
		return strings.TrimSpace(title)
	}

	return stripWikiConceptFilePrefix(strings.TrimSuffix(path.Base(document.RelativePath), ".md"))
}

func ensureAgentsSymlink(topicPath string) error {
	agentsPath := filepath.Join(topicPath, "AGENTS.md")
	if err := os.RemoveAll(agentsPath); err != nil {
		return fmt.Errorf("remove existing agents file: %w", err)
	}
	if err := os.Symlink("CLAUDE.md", agentsPath); err != nil {
		return fmt.Errorf("create agents symlink: %w", err)
	}

	return nil
}

func ensureTopicGitkeeps(topicPath string) error {
	for _, directoryPath := range []string{
		filepath.Join(topicPath, "raw", "articles"),
		filepath.Join(topicPath, "raw", "bookmarks"),
		filepath.Join(topicPath, "raw", "github"),
		filepath.Join(topicPath, "raw", "youtube"),
		filepath.Join(topicPath, "wiki", "concepts"),
		filepath.Join(topicPath, "wiki", "index"),
		filepath.Join(topicPath, "wiki", "codebase", "concepts"),
		filepath.Join(topicPath, "wiki", "codebase", "index"),
		filepath.Join(topicPath, "outputs", "briefings"),
		filepath.Join(topicPath, "outputs", "queries"),
		filepath.Join(topicPath, "outputs", "diagrams"),
		filepath.Join(topicPath, "outputs", "reports"),
		filepath.Join(topicPath, "bases"),
	} {
		if err := ensureGitkeep(directoryPath); err != nil {
			return err
		}
	}

	return nil
}

func ensureGitkeep(directoryPath string) error {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("read %q: %w", directoryPath, err)
	}
	if len(entries) > 0 {
		return nil
	}

	if err := writeTextFile(filepath.Join(directoryPath, ".gitkeep"), ""); err != nil {
		return fmt.Errorf("write gitkeep for %q: %w", directoryPath, err)
	}

	return nil
}

func appendLog(topic models.TopicMetadata, graph models.GraphSnapshot, counts WriteVaultResult) error {
	logPath := filepath.Join(topic.TopicPath, "log.md")

	if _, err := os.Stat(logPath); errors.Is(err, os.ErrNotExist) {
		bootstrap := strings.Join([]string{
			fmt.Sprintf("# %s - Log", topic.Title),
			"",
			"Chronological, append-only record of every knowledge-base operation on this topic.",
			"",
			fmt.Sprintf("## [%s] bootstrap | topic scaffolded", topic.Today),
			"",
			fmt.Sprintf("Topic `%s` created by `kodebase generate`. Domain: `%s`.", topic.Slug, topic.Domain),
			"",
		}, "\n")

		if err := writeTextFile(logPath, bootstrap); err != nil {
			return fmt.Errorf("write bootstrap log: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("stat log file %q: %w", logPath, err)
	}

	entry := strings.Join([]string{
		fmt.Sprintf(
			"## [%s] ingest | codebase (%d files, %d symbols)",
			topic.Today,
			len(graph.Files),
			len(graph.Symbols),
		),
		"",
		fmt.Sprintf("Refreshed `%s/raw/codebase/` from `%s`.", topic.Slug, topic.RootPath),
		"",
		fmt.Sprintf(
			"## [%s] compile | refreshed codebase wiki (%d concept pages, %d index pages)",
			topic.Today,
			counts.WikiDocumentsWritten,
			counts.IndexDocumentsWritten,
		),
		"",
	}, "\n")

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file %q for append: %w", logPath, err)
	}

	if _, err := file.WriteString(entry); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("append log entry: %w (close error: %v)", err, closeErr)
		}
		return fmt.Errorf("append log entry: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close log file %q: %w", logPath, err)
	}

	return nil
}

func countWrittenDocuments(documents []models.RenderedDocument) WriteVaultResult {
	var result WriteVaultResult

	for _, document := range documents {
		switch document.Kind {
		case models.DocRaw:
			result.RawDocumentsWritten++
		case models.DocWiki:
			result.WikiDocumentsWritten++
		case models.DocIndex:
			result.IndexDocumentsWritten++
		}
	}

	return result
}

func writeTextFile(targetPath, body string) error {
	if err := os.WriteFile(targetPath, []byte(body), 0o644); err != nil {
		return err
	}

	return nil
}
