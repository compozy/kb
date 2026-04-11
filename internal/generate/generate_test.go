package generate

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/user/go-devstack/internal/models"
	"github.com/user/go-devstack/internal/scanner"
	"github.com/user/go-devstack/internal/vault"
)

type fakeAdapter struct {
	name        string
	supported   map[models.SupportedLanguage]bool
	parseResult []models.ParsedFile
	parseErr    error
	calls       *[]string
}

func (a fakeAdapter) Supports(language models.SupportedLanguage) bool {
	return a.supported[language]
}

func (a fakeAdapter) ParseFiles(files []models.ScannedSourceFile, rootPath string) ([]models.ParsedFile, error) {
	if a.calls != nil {
		*a.calls = append(*a.calls, "parse:"+a.name)
	}

	if a.parseErr != nil {
		return nil, a.parseErr
	}

	return a.parseResult, nil
}

func (a fakeAdapter) ParseFilesWithProgress(
	files []models.ScannedSourceFile,
	rootPath string,
	report func(models.ScannedSourceFile),
) ([]models.ParsedFile, error) {
	parsedFiles, err := a.ParseFiles(files, rootPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if report != nil {
			report(file)
		}
	}

	return parsedFiles, nil
}

func TestRunnerGenerateCallsPipelineStagesInOrder(t *testing.T) {
	t.Parallel()

	var calls []string
	scannedWorkspace := &models.ScannedWorkspace{
		Files: []models.ScannedSourceFile{
			{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
		},
		FilesByLanguage: map[models.SupportedLanguage][]models.ScannedSourceFile{
			models.LangGo: {
				{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
			},
		},
	}

	graphSnapshot := models.GraphSnapshot{
		RootPath: "/repo",
		Files: []models.GraphFile{
			{ID: "file:main.go", FilePath: "main.go", Language: models.LangGo},
		},
		Symbols: []models.SymbolNode{
			{ID: "symbol:main.go:main:function:1:3", Name: "main", SymbolKind: "function", FilePath: "main.go"},
		},
		Relations: []models.RelationEdge{
			{FromID: "file:main.go", ToID: "symbol:main.go:main:function:1:3", Type: models.RelContains},
		},
	}

	generator := runner{
		scanWorkspace: func(rootPath string, opts ...scanner.Option) (*models.ScannedWorkspace, error) {
			calls = append(calls, "scan")
			return scannedWorkspace, nil
		},
		adapters: []models.LanguageAdapter{
			fakeAdapter{
				name:      "go",
				supported: map[models.SupportedLanguage]bool{models.LangGo: true},
				parseResult: []models.ParsedFile{
					{File: graphSnapshot.Files[0], Symbols: graphSnapshot.Symbols, Relations: graphSnapshot.Relations},
				},
				calls: &calls,
			},
		},
		normalizeGraph: func(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot {
			calls = append(calls, "normalize")
			return graphSnapshot
		},
		computeMetrics: func(graph models.GraphSnapshot) models.MetricsResult {
			calls = append(calls, "metrics")
			return models.MetricsResult{}
		},
		renderDocuments: func(
			graph models.GraphSnapshot,
			metrics models.MetricsResult,
			topic models.TopicMetadata,
		) []models.RenderedDocument {
			calls = append(calls, "render")
			return []models.RenderedDocument{
				{
					Kind:         models.DocWiki,
					ManagedArea:  models.AreaWikiConcept,
					RelativePath: vault.GetWikiConceptPath("Codebase Overview"),
					Frontmatter:  map[string]interface{}{"title": "Codebase Overview"},
					Body:         "---\ntitle: \"Codebase Overview\"\n---\n\n# Codebase Overview\n",
				},
			}
		},
		renderBaseFiles: func(metrics models.MetricsResult) []models.BaseFile {
			return []models.BaseFile{{RelativePath: "bases/module-health.base"}}
		},
		writeVault: func(ctx context.Context, options vault.WriteVaultOptions) (vault.WriteVaultResult, error) {
			calls = append(calls, "write")
			if options.Topic.VaultPath != "/vault" {
				t.Fatalf("topic vault path = %q, want /vault", options.Topic.VaultPath)
			}
			if options.Topic.TopicPath != "/vault/fixture" {
				t.Fatalf("topic path = %q, want /vault/fixture", options.Topic.TopicPath)
			}
			if options.Topic.Slug != "fixture" {
				t.Fatalf("topic slug = %q, want fixture", options.Topic.Slug)
			}
			return vault.WriteVaultResult{RawDocumentsWritten: 1, WikiDocumentsWritten: 1, IndexDocumentsWritten: 0}, nil
		},
		now: testClock(
			time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 1, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 1, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 1, 500000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 1, 500000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 0, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 100000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 100000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 200000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 200000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 300000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 300000000, time.UTC),
			time.Date(2026, 4, 9, 12, 0, 2, 400000000, time.UTC),
		),
	}

	summary, err := generator.Generate(context.Background(), models.GenerateOptions{
		RootPath:  "/repo",
		VaultPath: "/vault",
		TopicSlug: "fixture",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	expectedOrder := []string{"scan", "parse:go", "normalize", "metrics", "render", "write"}
	if !reflect.DeepEqual(calls, expectedOrder) {
		t.Fatalf("call order = %#v, want %#v", calls, expectedOrder)
	}

	if summary.FilesScanned != 1 || summary.FilesParsed != 1 || summary.SymbolsExtracted != 1 {
		t.Fatalf("unexpected summary counts: %#v", summary)
	}
	if summary.VaultPath != "/vault" || summary.TopicPath != "/vault/fixture" || summary.TopicSlug != "fixture" {
		t.Fatalf("unexpected summary paths: %#v", summary)
	}
	if summary.Timings.TotalMillis <= 0 {
		t.Fatalf("expected total timing to be recorded, got %#v", summary.Timings)
	}
}

func TestSelectAdaptersForGoOnlyWorkspace(t *testing.T) {
	t.Parallel()

	selected := selectAdapters(
		[]models.SupportedLanguage{models.LangGo},
		[]models.LanguageAdapter{
			fakeAdapter{name: "ts", supported: map[models.SupportedLanguage]bool{models.LangTS: true}},
			fakeAdapter{name: "go", supported: map[models.SupportedLanguage]bool{models.LangGo: true}},
		},
	)

	if len(selected) != 1 {
		t.Fatalf("selected %d adapters, want 1", len(selected))
	}
	if !selected[0].Supports(models.LangGo) {
		t.Fatalf("selected adapter does not support Go")
	}
}

func TestSelectAdaptersForMixedWorkspace(t *testing.T) {
	t.Parallel()

	selected := selectAdapters(
		[]models.SupportedLanguage{models.LangTS, models.LangGo},
		[]models.LanguageAdapter{
			fakeAdapter{name: "ts", supported: map[models.SupportedLanguage]bool{models.LangTS: true}},
			fakeAdapter{name: "go", supported: map[models.SupportedLanguage]bool{models.LangGo: true}},
		},
	)

	if len(selected) != 2 {
		t.Fatalf("selected %d adapters, want 2", len(selected))
	}
	if !selected[0].Supports(models.LangTS) || !selected[1].Supports(models.LangGo) {
		t.Fatalf("unexpected adapter selection order")
	}
}

func TestRunnerGenerateSummaryReportsCounts(t *testing.T) {
	t.Parallel()

	generator := runner{
		scanWorkspace: func(rootPath string, opts ...scanner.Option) (*models.ScannedWorkspace, error) {
			return &models.ScannedWorkspace{
				Files: []models.ScannedSourceFile{
					{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
					{AbsolutePath: "/repo/internal/helper.go", RelativePath: "internal/helper.go", Language: models.LangGo},
					{AbsolutePath: "/repo/web/index.ts", RelativePath: "web/index.ts", Language: models.LangTS},
				},
				FilesByLanguage: map[models.SupportedLanguage][]models.ScannedSourceFile{
					models.LangGo: {
						{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
						{AbsolutePath: "/repo/internal/helper.go", RelativePath: "internal/helper.go", Language: models.LangGo},
					},
					models.LangTS: {
						{AbsolutePath: "/repo/web/index.ts", RelativePath: "web/index.ts", Language: models.LangTS},
					},
				},
			}, nil
		},
		adapters: []models.LanguageAdapter{
			fakeAdapter{
				name:      "ts",
				supported: map[models.SupportedLanguage]bool{models.LangTS: true},
				parseResult: []models.ParsedFile{
					{File: models.GraphFile{ID: "file:web/index.ts", FilePath: "web/index.ts", Language: models.LangTS}},
				},
			},
			fakeAdapter{
				name:      "go",
				supported: map[models.SupportedLanguage]bool{models.LangGo: true},
				parseResult: []models.ParsedFile{
					{File: models.GraphFile{ID: "file:main.go", FilePath: "main.go", Language: models.LangGo}},
				},
			},
		},
		normalizeGraph: func(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot {
			return models.GraphSnapshot{
				RootPath: rootPath,
				Files: []models.GraphFile{
					{ID: "file:main.go", FilePath: "main.go", Language: models.LangGo},
					{ID: "file:web/index.ts", FilePath: "web/index.ts", Language: models.LangTS},
				},
				Symbols: []models.SymbolNode{
					{ID: "symbol:main", Name: "main", SymbolKind: "function", FilePath: "main.go"},
					{ID: "symbol:greet", Name: "greet", SymbolKind: "function", FilePath: "main.go"},
					{ID: "symbol:index", Name: "index", SymbolKind: "function", FilePath: "web/index.ts"},
					{ID: "symbol:helper", Name: "helper", SymbolKind: "function", FilePath: "web/index.ts"},
				},
				Relations: []models.RelationEdge{
					{FromID: "symbol:main", ToID: "symbol:greet", Type: models.RelCalls},
					{FromID: "symbol:index", ToID: "symbol:helper", Type: models.RelCalls},
					{FromID: "file:main.go", ToID: "file:web/index.ts", Type: models.RelImports},
					{FromID: "file:web/index.ts", ToID: "symbol:helper", Type: models.RelContains},
					{FromID: "file:main.go", ToID: "symbol:greet", Type: models.RelContains},
				},
			}
		},
		computeMetrics: func(graph models.GraphSnapshot) models.MetricsResult {
			return models.MetricsResult{}
		},
		renderDocuments: func(graph models.GraphSnapshot, metrics models.MetricsResult, topic models.TopicMetadata) []models.RenderedDocument {
			return []models.RenderedDocument{
				{
					Kind:         models.DocRaw,
					ManagedArea:  models.AreaRawCodebase,
					RelativePath: "raw/codebase/files/main.go.md",
					Frontmatter:  map[string]interface{}{"title": "main.go"},
					Body:         "---\ntitle: \"main.go\"\n---\n\n# main.go\n",
				},
			}
		},
		renderBaseFiles: func(metrics models.MetricsResult) []models.BaseFile {
			return nil
		},
		writeVault: func(ctx context.Context, options vault.WriteVaultOptions) (vault.WriteVaultResult, error) {
			return vault.WriteVaultResult{RawDocumentsWritten: 6, WikiDocumentsWritten: 10, IndexDocumentsWritten: 3}, nil
		},
		now: func() time.Time {
			return time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)
		},
	}

	summary, err := generator.Generate(context.Background(), models.GenerateOptions{
		RootPath: "/repo/demo-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 3 {
		t.Fatalf("FilesScanned = %d, want 3", summary.FilesScanned)
	}
	if summary.FilesParsed != 2 {
		t.Fatalf("FilesParsed = %d, want 2", summary.FilesParsed)
	}
	if summary.FilesSkipped != 1 {
		t.Fatalf("FilesSkipped = %d, want 1", summary.FilesSkipped)
	}
	if summary.SymbolsExtracted != 4 {
		t.Fatalf("SymbolsExtracted = %d, want 4", summary.SymbolsExtracted)
	}
	if summary.RelationsEmitted != 5 {
		t.Fatalf("RelationsEmitted = %d, want 5", summary.RelationsEmitted)
	}
	if summary.RawDocumentsWritten != 6 || summary.WikiDocumentsWritten != 10 || summary.IndexDocumentsWritten != 3 {
		t.Fatalf("unexpected document counts: %#v", summary)
	}
	if summary.TopicSlug != "demo-repo" {
		t.Fatalf("TopicSlug = %q, want demo-repo", summary.TopicSlug)
	}
}

func TestGenerateRequiresRootPath(t *testing.T) {
	t.Parallel()

	_, err := Generate(context.Background(), models.GenerateOptions{})
	if err == nil {
		t.Fatal("expected error for missing root path")
	}
	if !strings.Contains(err.Error(), "root path is required") {
		t.Fatalf("expected descriptive root path error, got %v", err)
	}
}

func TestResolveTargetUsesExplicitVaultPathAndTopicSlug(t *testing.T) {
	t.Parallel()

	target, err := resolveTarget(models.GenerateOptions{
		RootPath:  "/repo/source",
		VaultPath: "/vault/root",
		TopicSlug: "Custom Topic",
	})
	if err != nil {
		t.Fatalf("resolveTarget returned error: %v", err)
	}

	if target.RootPath != "/repo/source" {
		t.Fatalf("root path = %q, want /repo/source", target.RootPath)
	}
	if target.VaultPath != "/vault/root" {
		t.Fatalf("vault path = %q, want /vault/root", target.VaultPath)
	}
	if target.TopicSlug != "custom-topic" {
		t.Fatalf("topic slug = %q, want custom-topic", target.TopicSlug)
	}
}

func TestResolveTargetDefaultsVaultPathAndTopicSlugFromRootPath(t *testing.T) {
	t.Parallel()

	target, err := resolveTarget(models.GenerateOptions{RootPath: "/repo/source/demo-app"})
	if err != nil {
		t.Fatalf("resolveTarget returned error: %v", err)
	}

	if target.VaultPath != "/repo/source/demo-app/.kodebase/vault" {
		t.Fatalf("vault path = %q, want /repo/source/demo-app/.kodebase/vault", target.VaultPath)
	}
	if target.TopicSlug != "demo-app" {
		t.Fatalf("topic slug = %q, want demo-app", target.TopicSlug)
	}
}

func TestGenerateRespectsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Generate(ctx, models.GenerateOptions{RootPath: "."})
	if err == nil {
		t.Fatal("expected canceled context error")
	}
	if !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}

func TestRunnerGenerateEmitsParseAndWriteProgressEvents(t *testing.T) {
	t.Parallel()

	var events []Event
	observer := ObserverFunc(func(_ context.Context, event Event) {
		events = append(events, event)
	})

	generator := runner{
		scanWorkspace: func(rootPath string, opts ...scanner.Option) (*models.ScannedWorkspace, error) {
			return &models.ScannedWorkspace{
				Files: []models.ScannedSourceFile{
					{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
					{AbsolutePath: "/repo/helper.go", RelativePath: "helper.go", Language: models.LangGo},
				},
				FilesByLanguage: map[models.SupportedLanguage][]models.ScannedSourceFile{
					models.LangGo: {
						{AbsolutePath: "/repo/main.go", RelativePath: "main.go", Language: models.LangGo},
						{AbsolutePath: "/repo/helper.go", RelativePath: "helper.go", Language: models.LangGo},
					},
				},
			}, nil
		},
		adapters: []models.LanguageAdapter{
			fakeAdapter{
				name:      "go",
				supported: map[models.SupportedLanguage]bool{models.LangGo: true},
				parseResult: []models.ParsedFile{
					{File: models.GraphFile{ID: "file:main.go", FilePath: "main.go", Language: models.LangGo}},
					{File: models.GraphFile{ID: "file:helper.go", FilePath: "helper.go", Language: models.LangGo}},
				},
			},
		},
		normalizeGraph: func(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot {
			return models.GraphSnapshot{
				RootPath: rootPath,
				Files: []models.GraphFile{
					{ID: "file:main.go", FilePath: "main.go", Language: models.LangGo},
					{ID: "file:helper.go", FilePath: "helper.go", Language: models.LangGo},
				},
			}
		},
		computeMetrics: func(graph models.GraphSnapshot) models.MetricsResult {
			return models.MetricsResult{}
		},
		renderDocuments: func(graph models.GraphSnapshot, metrics models.MetricsResult, topic models.TopicMetadata) []models.RenderedDocument {
			return []models.RenderedDocument{
				{
					Kind:         models.DocRaw,
					ManagedArea:  models.AreaRawCodebase,
					RelativePath: "raw/codebase/files/main.go.md",
					Frontmatter:  map[string]interface{}{"title": "main.go"},
					Body:         "---\ntitle: \"main.go\"\n---\n\n# main.go\n",
				},
			}
		},
		renderBaseFiles: func(metrics models.MetricsResult) []models.BaseFile {
			return []models.BaseFile{{RelativePath: "bases/module-health.base"}}
		},
		writeVault: func(ctx context.Context, options vault.WriteVaultOptions) (vault.WriteVaultResult, error) {
			if options.Progress == nil {
				t.Fatal("expected write progress callback to be wired")
			}

			options.Progress(vault.WriteProgress{Completed: 1, Total: 4, Path: "raw/codebase/files/main.go.md"})
			options.Progress(vault.WriteProgress{Completed: 2, Total: 4, Path: "bases/module-health.base"})
			options.Progress(vault.WriteProgress{Completed: 3, Total: 4, Path: "CLAUDE.md"})
			options.Progress(vault.WriteProgress{Completed: 4, Total: 4, Path: "log.md"})

			return vault.WriteVaultResult{RawDocumentsWritten: 1}, nil
		},
		now: func() time.Time {
			return time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)
		},
	}

	if _, err := generator.GenerateWithObserver(context.Background(), models.GenerateOptions{RootPath: "/repo"}, observer); err != nil {
		t.Fatalf("GenerateWithObserver returned error: %v", err)
	}

	parseStarted := firstEvent(events, EventStageStarted, "parse")
	if parseStarted.Total != 2 {
		t.Fatalf("parse start total = %d, want 2", parseStarted.Total)
	}

	parseProgress := filterEvents(events, EventStageProgress, "parse")
	if len(parseProgress) != 2 {
		t.Fatalf("parse progress events = %d, want 2", len(parseProgress))
	}
	if parseProgress[0].Completed != 1 || parseProgress[1].Completed != 2 {
		t.Fatalf("unexpected parse progress events: %#v", parseProgress)
	}

	writeStarted := firstEvent(events, EventStageStarted, "write")
	if writeStarted.Total != 4 {
		t.Fatalf("write start total = %d, want 4", writeStarted.Total)
	}

	writeProgress := filterEvents(events, EventStageProgress, "write")
	if len(writeProgress) != 4 {
		t.Fatalf("write progress events = %d, want 4", len(writeProgress))
	}
	if writeProgress[3].Completed != 4 || writeProgress[3].Total != 4 {
		t.Fatalf("unexpected write progress events: %#v", writeProgress)
	}
}

func filterEvents(events []Event, kind EventKind, stage string) []Event {
	filtered := make([]Event, 0, len(events))
	for _, event := range events {
		if event.Kind == kind && event.Stage == stage {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func firstEvent(events []Event, kind EventKind, stage string) Event {
	for _, event := range events {
		if event.Kind == kind && event.Stage == stage {
			return event
		}
	}
	return Event{}
}

func testClock(instants ...time.Time) func() time.Time {
	index := 0

	return func() time.Time {
		if len(instants) == 0 {
			return time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)
		}
		if index >= len(instants) {
			return instants[len(instants)-1]
		}

		value := instants[index]
		index++
		return value
	}
}
