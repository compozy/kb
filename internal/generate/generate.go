package generate

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/go-devstack/internal/adapter"
	"github.com/user/go-devstack/internal/graph"
	"github.com/user/go-devstack/internal/metrics"
	"github.com/user/go-devstack/internal/models"
	"github.com/user/go-devstack/internal/scanner"
	"github.com/user/go-devstack/internal/vault"
)

type scanWorkspaceFunc func(rootPath string, opts ...scanner.Option) (*models.ScannedWorkspace, error)
type normalizeGraphFunc func(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot
type computeMetricsFunc func(graph models.GraphSnapshot) models.MetricsResult
type renderDocumentsFunc func(
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
	topic models.TopicMetadata,
) []models.RenderedDocument
type renderBaseFilesFunc func(metrics models.MetricsResult) []models.BaseFile
type writeVaultFunc func(ctx context.Context, options vault.WriteVaultOptions) (vault.WriteVaultResult, error)

type runner struct {
	scanWorkspace   scanWorkspaceFunc
	adapters        []models.LanguageAdapter
	normalizeGraph  normalizeGraphFunc
	computeMetrics  computeMetricsFunc
	renderDocuments renderDocumentsFunc
	renderBaseFiles renderBaseFilesFunc
	writeVault      writeVaultFunc
	now             func() time.Time
	logger          *slog.Logger
}

// Generate runs the full repository-to-vault pipeline and returns a structured
// summary of the generated topic.
func Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
	return newRunner().Generate(ctx, opts)
}

func newRunner() runner {
	return runner{
		scanWorkspace: scanner.ScanWorkspace,
		adapters: []models.LanguageAdapter{
			adapter.TSAdapter{},
			adapter.GoAdapter{},
		},
		normalizeGraph:  graph.NormalizeGraph,
		computeMetrics:  metrics.ComputeMetrics,
		renderDocuments: vault.RenderDocuments,
		renderBaseFiles: vault.RenderBaseFiles,
		writeVault:      vault.WriteVault,
		now: func() time.Time {
			return time.Now().UTC()
		},
		logger: slog.Default(),
	}
}

// Generate runs the full repository-to-vault pipeline using the configured
// runner dependencies. Tests use this to substitute stage implementations.
func (r runner) Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
	if ctx == nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: nil context")
	}

	r = r.withDefaults()

	rootPath, vaultPath, err := resolvePaths(opts)
	if err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	totalStartedAt := r.now()
	topic := r.createTopicMetadata(rootPath, vaultPath, opts)
	timings := models.GenerationTimings{}

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("scan", "root_path", rootPath, "vault_path", vaultPath)
	stageStartedAt := r.now()
	scannedWorkspace, err := r.scanWorkspace(
		rootPath,
		scanner.WithOutputPath(vaultPath),
		scanner.WithIncludePatterns(opts.IncludePatterns...),
		scanner.WithExcludePatterns(opts.ExcludePatterns...),
	)
	timings.ScanMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	if err != nil {
		r.logStageFailed("scan", err, timings.ScanMillis)
		return models.GenerationSummary{}, fmt.Errorf("generate: scan workspace: %w", err)
	}

	languages := workspaceLanguages(scannedWorkspace)
	r.logStageCompleted(
		"scan",
		timings.ScanMillis,
		"files_scanned", len(scannedWorkspace.Files),
		"languages", languageNames(languages),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("select_adapters", "languages", languageNames(languages))
	stageStartedAt = r.now()
	selectedAdapters := selectAdapters(languages, r.adapters)
	timings.SelectAdaptersMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.logStageCompleted(
		"select_adapters",
		timings.SelectAdaptersMillis,
		"adapter_count", len(selectedAdapters),
		"adapters", adapterNames(selectedAdapters),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("parse", "adapter_count", len(selectedAdapters))
	stageStartedAt = r.now()
	parsedFiles := make([]models.ParsedFile, 0, len(scannedWorkspace.Files))
	for _, languageAdapter := range selectedAdapters {
		files := filesForAdapter(scannedWorkspace.Files, languageAdapter)
		if len(files) == 0 {
			continue
		}

		entries, err := languageAdapter.ParseFiles(files, rootPath)
		if err != nil {
			parseElapsed := elapsedMillis(r.now().Sub(stageStartedAt))
			r.logStageFailed("parse", err, parseElapsed)
			return models.GenerationSummary{}, fmt.Errorf("generate: parse files: %w", err)
		}

		parsedFiles = append(parsedFiles, entries...)
	}
	timings.ParseMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.logStageCompleted("parse", timings.ParseMillis, "parsed_files", len(parsedFiles))

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("normalize", "parsed_files", len(parsedFiles))
	stageStartedAt = r.now()
	graphSnapshot := r.normalizeGraph(rootPath, parsedFiles)
	timings.NormalizeMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.logStageCompleted(
		"normalize",
		timings.NormalizeMillis,
		"files", len(graphSnapshot.Files),
		"symbols", len(graphSnapshot.Symbols),
		"relations", len(graphSnapshot.Relations),
		"diagnostics", len(graphSnapshot.Diagnostics),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("metrics", "files", len(graphSnapshot.Files), "symbols", len(graphSnapshot.Symbols))
	stageStartedAt = r.now()
	metricResult := r.computeMetrics(graphSnapshot)
	timings.MetricsMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.logStageCompleted(
		"metrics",
		timings.MetricsMillis,
		"directories", len(metricResult.Directories),
		"file_metrics", len(metricResult.Files),
		"symbol_metrics", len(metricResult.Symbols),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("render", "topic_slug", topic.Slug)
	stageStartedAt = r.now()
	documents := r.renderDocuments(graphSnapshot, metricResult, topic)
	baseFiles := r.renderBaseFiles(metricResult)
	timings.RenderMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.logStageCompleted(
		"render",
		timings.RenderMillis,
		"documents", len(documents),
		"base_files", len(baseFiles),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.logStageStarted("write", "topic_path", topic.TopicPath)
	stageStartedAt = r.now()
	writeResult, err := r.writeVault(ctx, vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graphSnapshot,
		Documents: documents,
		BaseFiles: baseFiles,
	})
	timings.WriteMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	if err != nil {
		r.logStageFailed("write", err, timings.WriteMillis)
		return models.GenerationSummary{}, fmt.Errorf("generate: write vault: %w", err)
	}
	r.logStageCompleted(
		"write",
		timings.WriteMillis,
		"raw_documents", writeResult.RawDocumentsWritten,
		"wiki_documents", writeResult.WikiDocumentsWritten,
		"index_documents", writeResult.IndexDocumentsWritten,
	)

	timings.TotalMillis = elapsedMillis(r.now().Sub(totalStartedAt))

	return models.GenerationSummary{
		Command:               "generate",
		RootPath:              rootPath,
		VaultPath:             vaultPath,
		TopicPath:             topic.TopicPath,
		TopicSlug:             topic.Slug,
		FilesScanned:          len(scannedWorkspace.Files),
		FilesParsed:           len(graphSnapshot.Files),
		FilesSkipped:          len(scannedWorkspace.Files) - len(graphSnapshot.Files),
		SymbolsExtracted:      len(graphSnapshot.Symbols),
		RelationsEmitted:      len(graphSnapshot.Relations),
		RawDocumentsWritten:   writeResult.RawDocumentsWritten,
		WikiDocumentsWritten:  writeResult.WikiDocumentsWritten,
		IndexDocumentsWritten: writeResult.IndexDocumentsWritten,
		Timings:               timings,
		Diagnostics:           graphSnapshot.Diagnostics,
	}, nil
}

func (r runner) withDefaults() runner {
	if r.scanWorkspace == nil {
		r.scanWorkspace = scanner.ScanWorkspace
	}
	if len(r.adapters) == 0 {
		r.adapters = []models.LanguageAdapter{adapter.TSAdapter{}, adapter.GoAdapter{}}
	}
	if r.normalizeGraph == nil {
		r.normalizeGraph = graph.NormalizeGraph
	}
	if r.computeMetrics == nil {
		r.computeMetrics = metrics.ComputeMetrics
	}
	if r.renderDocuments == nil {
		r.renderDocuments = vault.RenderDocuments
	}
	if r.renderBaseFiles == nil {
		r.renderBaseFiles = vault.RenderBaseFiles
	}
	if r.writeVault == nil {
		r.writeVault = vault.WriteVault
	}
	if r.now == nil {
		r.now = func() time.Time {
			return time.Now().UTC()
		}
	}
	if r.logger == nil {
		r.logger = slog.Default()
	}

	return r
}

func resolvePaths(opts models.GenerateOptions) (string, string, error) {
	if strings.TrimSpace(opts.RootPath) == "" {
		return "", "", fmt.Errorf("root path is required")
	}

	rootPath, err := filepath.Abs(opts.RootPath)
	if err != nil {
		return "", "", fmt.Errorf("resolve root path: %w", err)
	}

	if strings.TrimSpace(opts.OutputPath) == "" {
		return rootPath, filepath.Join(rootPath, ".kodebase", "vault"), nil
	}

	outputPath, err := filepath.Abs(opts.OutputPath)
	if err != nil {
		return "", "", fmt.Errorf("resolve output path: %w", err)
	}

	return rootPath, outputPath, nil
}

func (r runner) createTopicMetadata(rootPath string, vaultPath string, opts models.GenerateOptions) models.TopicMetadata {
	slugSource := rootPath
	if strings.TrimSpace(opts.Topic) != "" {
		slugSource = opts.Topic
	}

	slug := vault.DeriveTopicSlug(slugSource)
	title := strings.TrimSpace(opts.Title)
	if title == "" {
		title = vault.DeriveTopicTitle(slug)
	}

	domainSource := slug
	if strings.TrimSpace(opts.Domain) != "" {
		domainSource = opts.Domain
	}

	return models.TopicMetadata{
		RootPath:  rootPath,
		Title:     title,
		Slug:      slug,
		Domain:    vault.DeriveTopicDomain(domainSource),
		Today:     r.now().Format("2006-01-02"),
		VaultPath: vaultPath,
		TopicPath: filepath.Join(vaultPath, slug),
	}
}

func workspaceLanguages(workspace *models.ScannedWorkspace) []models.SupportedLanguage {
	if workspace == nil {
		return nil
	}

	languages := make([]models.SupportedLanguage, 0, len(workspace.FilesByLanguage))
	for _, language := range models.SupportedLanguages() {
		if len(workspace.FilesByLanguage[language]) == 0 {
			continue
		}

		languages = append(languages, language)
	}

	return languages
}

func selectAdapters(languages []models.SupportedLanguage, adapters []models.LanguageAdapter) []models.LanguageAdapter {
	selected := make([]models.LanguageAdapter, 0, len(adapters))
	for _, languageAdapter := range adapters {
		for _, language := range languages {
			if !languageAdapter.Supports(language) {
				continue
			}

			selected = append(selected, languageAdapter)
			break
		}
	}

	return selected
}

func filesForAdapter(files []models.ScannedSourceFile, languageAdapter models.LanguageAdapter) []models.ScannedSourceFile {
	filtered := make([]models.ScannedSourceFile, 0, len(files))
	for _, file := range files {
		if !languageAdapter.Supports(file.Language) {
			continue
		}

		filtered = append(filtered, file)
	}

	return filtered
}

func languageNames(languages []models.SupportedLanguage) []string {
	names := make([]string, 0, len(languages))
	for _, language := range languages {
		names = append(names, string(language))
	}

	return names
}

func adapterNames(adapters []models.LanguageAdapter) []string {
	names := make([]string, 0, len(adapters))
	for _, languageAdapter := range adapters {
		names = append(names, fmt.Sprintf("%T", languageAdapter))
	}

	return names
}

func elapsedMillis(duration time.Duration) int64 {
	if duration < 0 {
		return 0
	}

	return duration.Milliseconds()
}

func (r runner) logStageStarted(stage string, attrs ...any) {
	r.logger.Info("generate stage started", append([]any{"stage", stage}, attrs...)...)
}

func (r runner) logStageCompleted(stage string, durationMillis int64, attrs ...any) {
	base := []any{"stage", stage, "duration_ms", durationMillis}
	r.logger.Info("generate stage completed", append(base, attrs...)...)
}

func (r runner) logStageFailed(stage string, err error, durationMillis int64) {
	r.logger.Error("generate stage failed", "stage", stage, "duration_ms", durationMillis, "error", err)
}
