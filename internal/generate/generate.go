package generate

import (
	"context"
	"fmt"
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

type progressAwareLanguageAdapter interface {
	ParseFilesWithProgress(
		files []models.ScannedSourceFile,
		rootPath string,
		report func(models.ScannedSourceFile),
	) ([]models.ParsedFile, error)
}

type runner struct {
	scanWorkspace   scanWorkspaceFunc
	adapters        []models.LanguageAdapter
	normalizeGraph  normalizeGraphFunc
	computeMetrics  computeMetricsFunc
	renderDocuments renderDocumentsFunc
	renderBaseFiles renderBaseFilesFunc
	writeVault      writeVaultFunc
	now             func() time.Time
	observer        Observer
}

// Generate runs the full repository-to-vault pipeline and returns a structured
// summary of the generated topic.
func Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
	return GenerateWithObserver(ctx, opts, nil)
}

// GenerateWithObserver runs the full pipeline and emits structured events to
// the provided observer while returning the final summary.
func GenerateWithObserver(ctx context.Context, opts models.GenerateOptions, observer Observer) (models.GenerationSummary, error) {
	return newRunner().GenerateWithObserver(ctx, opts, observer)
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
		observer: noopObserver{},
	}
}

// Generate runs the full repository-to-vault pipeline using the configured
// runner dependencies. Tests use this to substitute stage implementations.
func (r runner) Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
	return r.GenerateWithObserver(ctx, opts, nil)
}

// GenerateWithObserver runs the configured pipeline and reports structured
// events to the provided observer.
func (r runner) GenerateWithObserver(ctx context.Context, opts models.GenerateOptions, observer Observer) (models.GenerationSummary, error) {
	if ctx == nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: nil context")
	}

	r = r.withDefaults()
	if observer != nil {
		r.observer = observer
	}

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

	r.emitStageStarted(ctx, "scan", 0, "root_path", rootPath, "vault_path", vaultPath)
	stageStartedAt := r.now()
	scannedWorkspace, err := r.scanWorkspace(
		rootPath,
		scanner.WithOutputPath(vaultPath),
		scanner.WithIncludePatterns(opts.IncludePatterns...),
		scanner.WithExcludePatterns(opts.ExcludePatterns...),
	)
	timings.ScanMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	if err != nil {
		r.emitStageFailed(ctx, "scan", err, timings.ScanMillis, 0, 0)
		return models.GenerationSummary{}, fmt.Errorf("generate: scan workspace: %w", err)
	}

	languages := workspaceLanguages(scannedWorkspace)
	r.emitStageCompleted(
		ctx,
		"scan",
		timings.ScanMillis,
		0,
		0,
		"files_scanned", len(scannedWorkspace.Files),
		"languages", languageNames(languages),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.emitStageStarted(ctx, "select_adapters", 0, "languages", languageNames(languages))
	stageStartedAt = r.now()
	selectedAdapters := selectAdapters(languages, r.adapters)
	timings.SelectAdaptersMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.emitStageCompleted(
		ctx,
		"select_adapters",
		timings.SelectAdaptersMillis,
		0,
		0,
		"adapter_count", len(selectedAdapters),
		"adapters", adapterNames(selectedAdapters),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	parseTotal := len(scannedWorkspace.Files)
	parseCompleted := 0
	r.emitStageStarted(ctx, "parse", parseTotal, "adapter_count", len(selectedAdapters))
	stageStartedAt = r.now()
	parsedFiles := make([]models.ParsedFile, 0, len(scannedWorkspace.Files))
	reportParseProgress := func(file models.ScannedSourceFile) {
		parseCompleted++
		r.emitStageProgress(ctx, "parse", parseCompleted, parseTotal, "path", file.RelativePath)
	}
	for _, languageAdapter := range selectedAdapters {
		files := filesForAdapter(scannedWorkspace.Files, languageAdapter)
		if len(files) == 0 {
			continue
		}

		var entries []models.ParsedFile
		progressAdapter, supportsProgress := languageAdapter.(progressAwareLanguageAdapter)
		if supportsProgress {
			entries, err = progressAdapter.ParseFilesWithProgress(files, rootPath, reportParseProgress)
		} else {
			entries, err = languageAdapter.ParseFiles(files, rootPath)
			if err == nil && len(files) > 0 {
				parseCompleted += len(files)
				r.emitStageProgress(ctx, "parse", parseCompleted, parseTotal)
			}
		}
		if err != nil {
			parseElapsed := elapsedMillis(r.now().Sub(stageStartedAt))
			r.emitStageFailed(ctx, "parse", err, parseElapsed, parseCompleted, parseTotal)
			return models.GenerationSummary{}, fmt.Errorf("generate: parse files: %w", err)
		}

		parsedFiles = append(parsedFiles, entries...)
	}
	timings.ParseMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.emitStageCompleted(ctx, "parse", timings.ParseMillis, parseCompleted, parseTotal, "parsed_files", len(parsedFiles))

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.emitStageStarted(ctx, "normalize", 0, "parsed_files", len(parsedFiles))
	stageStartedAt = r.now()
	graphSnapshot := r.normalizeGraph(rootPath, parsedFiles)
	timings.NormalizeMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.emitStageCompleted(
		ctx,
		"normalize",
		timings.NormalizeMillis,
		0,
		0,
		"files", len(graphSnapshot.Files),
		"symbols", len(graphSnapshot.Symbols),
		"relations", len(graphSnapshot.Relations),
		"diagnostics", len(graphSnapshot.Diagnostics),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.emitStageStarted(ctx, "metrics", 0, "files", len(graphSnapshot.Files), "symbols", len(graphSnapshot.Symbols))
	stageStartedAt = r.now()
	metricResult := r.computeMetrics(graphSnapshot)
	timings.MetricsMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.emitStageCompleted(
		ctx,
		"metrics",
		timings.MetricsMillis,
		0,
		0,
		"directories", len(metricResult.Directories),
		"file_metrics", len(metricResult.Files),
		"symbol_metrics", len(metricResult.Symbols),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	r.emitStageStarted(ctx, "render", 0, "topic_slug", topic.Slug)
	stageStartedAt = r.now()
	documents := r.renderDocuments(graphSnapshot, metricResult, topic)
	baseFiles := r.renderBaseFiles(metricResult)
	timings.RenderMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	r.emitStageCompleted(
		ctx,
		"render",
		timings.RenderMillis,
		0,
		0,
		"documents", len(documents),
		"base_files", len(baseFiles),
	)

	if err := ctx.Err(); err != nil {
		return models.GenerationSummary{}, fmt.Errorf("generate: %w", err)
	}

	writeTotal := len(documents) + len(baseFiles) + 2
	r.emitStageStarted(ctx, "write", writeTotal, "topic_path", topic.TopicPath)
	stageStartedAt = r.now()
	writeResult, err := r.writeVault(ctx, vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graphSnapshot,
		Documents: documents,
		BaseFiles: baseFiles,
		Progress: func(progress vault.WriteProgress) {
			r.emitStageProgress(ctx, "write", progress.Completed, progress.Total, "path", progress.Path)
		},
	})
	timings.WriteMillis = elapsedMillis(r.now().Sub(stageStartedAt))
	if err != nil {
		r.emitStageFailed(ctx, "write", err, timings.WriteMillis, 0, writeTotal)
		return models.GenerationSummary{}, fmt.Errorf("generate: write vault: %w", err)
	}
	r.emitStageCompleted(
		ctx,
		"write",
		timings.WriteMillis,
		writeTotal,
		writeTotal,
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
	if r.observer == nil {
		r.observer = noopObserver{}
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

func (r runner) emitStageStarted(ctx context.Context, stage string, total int, attrs ...any) {
	r.observer.ObserveGenerateEvent(ctx, Event{
		Kind:   EventStageStarted,
		Stage:  stage,
		Fields: eventFields(attrs...),
		Total:  total,
	})
}

func (r runner) emitStageProgress(ctx context.Context, stage string, completed int, total int, attrs ...any) {
	r.observer.ObserveGenerateEvent(ctx, Event{
		Kind:      EventStageProgress,
		Stage:     stage,
		Completed: completed,
		Fields:    eventFields(attrs...),
		Total:     total,
	})
}

func (r runner) emitStageCompleted(ctx context.Context, stage string, durationMillis int64, completed int, total int, attrs ...any) {
	r.observer.ObserveGenerateEvent(ctx, Event{
		Kind:           EventStageCompleted,
		Stage:          stage,
		Completed:      completed,
		DurationMillis: durationMillis,
		Fields:         eventFields(attrs...),
		Total:          total,
	})
}

func (r runner) emitStageFailed(ctx context.Context, stage string, err error, durationMillis int64, completed int, total int) {
	r.observer.ObserveGenerateEvent(ctx, Event{
		Kind:           EventStageFailed,
		Stage:          stage,
		Completed:      completed,
		DurationMillis: durationMillis,
		Error:          err.Error(),
		Total:          total,
	})
}

func eventFields(attrs ...any) map[string]any {
	if len(attrs) == 0 {
		return nil
	}

	fields := make(map[string]any, len(attrs)/2)
	for index := 0; index+1 < len(attrs); index += 2 {
		key, ok := attrs[index].(string)
		if !ok || strings.TrimSpace(key) == "" {
			continue
		}
		fields[key] = attrs[index+1]
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}
