package adapter

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestTSAdapterSupportsTSLikeLanguages(t *testing.T) {
	t.Parallel()

	adapter := TSAdapter{}

	for _, language := range []models.SupportedLanguage{models.LangTS, models.LangTSX, models.LangJS, models.LangJSX} {
		if !adapter.Supports(language) {
			t.Fatalf("expected TSAdapter to support %q", language)
		}
	}

	for _, language := range []models.SupportedLanguage{models.LangGo, models.LangRust, models.LangJava} {
		if adapter.Supports(language) {
			t.Fatalf("expected TSAdapter to reject %q", language)
		}
	}
}

func TestTSAdapterParsesSimpleTypeScriptFile(t *testing.T) {
	t.Parallel()

	parsed := parseSingleTSLikeFile(t, "sample.ts", `
// Module docs.
export interface Runner {
	run(): void
}

// Add sums two integers.
export function Add(left: number, right: number): number {
	return left + right
}
`)

	if parsed.File.ModuleDoc != "Module docs." {
		t.Fatalf("module doc = %q", parsed.File.ModuleDoc)
	}

	if len(parsed.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(parsed.Diagnostics))
	}

	runner := mustFindSymbol(t, parsed.Symbols, "Runner")
	if runner.SymbolKind != tsSymbolKindInterface {
		t.Fatalf("Runner kind = %q, want %q", runner.SymbolKind, tsSymbolKindInterface)
	}

	add := mustFindSymbol(t, parsed.Symbols, "Add")
	if add.SymbolKind != tsSymbolKindFunction {
		t.Fatalf("Add kind = %q, want %q", add.SymbolKind, tsSymbolKindFunction)
	}
	if add.Signature != "function Add(left: number, right: number): number" {
		t.Fatalf("Add signature = %q", add.Signature)
	}
	if add.DocComment != "Add sums two integers." {
		t.Fatalf("Add doc comment = %q", add.DocComment)
	}
	if !add.Exported {
		t.Fatal("expected Add to be marked exported")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, add.ID, models.RelContains) {
		t.Fatal("expected file contains relation for Add")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, add.ID, models.RelExports) {
		t.Fatal("expected file exports relation for Add")
	}
}

func TestTSAdapterParsesTSXComponent(t *testing.T) {
	t.Parallel()

	parsed := parseSingleTSLikeFile(t, "app.tsx", `
export default function App() {
	return <Button onClick={handleClick} />
}
`)

	app := mustFindSymbol(t, parsed.Symbols, "App")
	if app.SymbolKind != tsSymbolKindFunction {
		t.Fatalf("App kind = %q, want %q", app.SymbolKind, tsSymbolKindFunction)
	}
	if !app.Exported {
		t.Fatal("expected App to be exported")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, app.ID, models.RelExports) {
		t.Fatal("expected export relation for App")
	}
}

func TestTSAdapterParsesJavaScriptRequireImports(t *testing.T) {
	t.Parallel()

	parsedFiles := parseTSLikeSources(t, map[string]string{
		"dep.js": `
function work() {
	return 1
}

exports.work = work
`,
		"main.js": `
const dep = require("./dep")

function run() {
	return dep.work()
}

module.exports = run
`,
	})

	mainFile := mustFindParsedFile(t, parsedFiles, "main.js")
	depFile := mustFindParsedFile(t, parsedFiles, "dep.js")
	run := mustFindSymbol(t, mainFile.Symbols, "run")
	work := mustFindSymbol(t, depFile.Symbols, "work")

	if !hasRelation(mainFile.Relations, mainFile.File.ID, depFile.File.ID, models.RelImports) {
		t.Fatal("expected require() to create an imports relation")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, work.ID, models.RelReferences) {
		t.Fatal("expected namespace require binding to create a references relation")
	}
	if !hasRelation(mainFile.Relations, run.ID, work.ID, models.RelCalls) {
		t.Fatal("expected dep.work() to resolve to the exported work symbol")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, run.ID, models.RelExports) {
		t.Fatal("expected module.exports = run to create a default export relation")
	}
}

func TestTSAdapterExtractsClassAndMethodSymbols(t *testing.T) {
	t.Parallel()

	parsed := parseSingleTSLikeFile(t, "widget.ts", `
export class Widget {
	// Run the widget once.
	run(flag: boolean): number {
		if (flag) {
			return helper()
		}

		return 0
	}
}

function helper(): number {
	return 1
}
`)

	widget := mustFindSymbol(t, parsed.Symbols, "Widget")
	if widget.SymbolKind != tsSymbolKindClass {
		t.Fatalf("Widget kind = %q, want %q", widget.SymbolKind, tsSymbolKindClass)
	}

	run := mustFindSymbol(t, parsed.Symbols, "run")
	if run.SymbolKind != tsSymbolKindMethod {
		t.Fatalf("run kind = %q, want %q", run.SymbolKind, tsSymbolKindMethod)
	}
	if run.DocComment != "Run the widget once." {
		t.Fatalf("run doc comment = %q", run.DocComment)
	}
	if !hasRelation(parsed.Relations, run.ID, mustFindSymbol(t, parsed.Symbols, "helper").ID, models.RelCalls) {
		t.Fatal("expected method call relation to helper")
	}
}

func TestTSAdapterBuildsImportBindingsForDefaultNamedAndNamespaceImports(t *testing.T) {
	t.Parallel()

	parsedFiles := parseTSLikeSources(t, map[string]string{
		"dep.ts": `
export default function helper(): number {
	return 1
}

export function work(): number {
	return helper()
}

export const answer: number = 42
`,
		"main.ts": `
import helper, { work as runWork, answer } from "./dep"
import * as depNS from "./dep"

export default function Run(): number {
	helper()
	runWork()
	depNS.work()
	return answer
}
`,
	})

	depFile := mustFindParsedFile(t, parsedFiles, "dep.ts")
	mainFile := mustFindParsedFile(t, parsedFiles, "main.ts")

	helper := mustFindSymbol(t, depFile.Symbols, "helper")
	work := mustFindSymbol(t, depFile.Symbols, "work")
	answer := mustFindSymbol(t, depFile.Symbols, "answer")
	run := mustFindSymbol(t, mainFile.Symbols, "Run")

	if !hasRelation(mainFile.Relations, mainFile.File.ID, depFile.File.ID, models.RelImports) {
		t.Fatal("expected import relation to dep.ts")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, helper.ID, models.RelReferences) {
		t.Fatal("expected default import reference to helper")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, work.ID, models.RelReferences) {
		t.Fatal("expected named import reference to work")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, answer.ID, models.RelReferences) {
		t.Fatal("expected named import reference to answer")
	}
	if !hasRelation(mainFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected helper() to resolve to default import target")
	}
	if !hasRelation(mainFile.Relations, run.ID, work.ID, models.RelCalls) {
		t.Fatal("expected runWork()/depNS.work() to resolve to work")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, run.ID, models.RelExports) {
		t.Fatal("expected default export relation for Run")
	}
}

func TestTSAdapterHandlesNamedAndStarReExports(t *testing.T) {
	t.Parallel()

	parsedFiles := parseTSLikeSources(t, map[string]string{
		"dep.ts": `
export function helper(): number {
	return 1
}

export const extra = 2
`,
		"barrel.ts": `
export { helper as renamed } from "./dep"
export * from "./dep"
`,
		"main.ts": `
import { renamed, extra } from "./barrel"

export function run(): number {
	return renamed() + extra
}
`,
	})

	depFile := mustFindParsedFile(t, parsedFiles, "dep.ts")
	barrelFile := mustFindParsedFile(t, parsedFiles, "barrel.ts")
	mainFile := mustFindParsedFile(t, parsedFiles, "main.ts")

	helper := mustFindSymbol(t, depFile.Symbols, "helper")
	extra := mustFindSymbol(t, depFile.Symbols, "extra")
	run := mustFindSymbol(t, mainFile.Symbols, "run")

	if !hasRelation(barrelFile.Relations, barrelFile.File.ID, helper.ID, models.RelExports) {
		t.Fatal("expected named re-export relation for helper")
	}
	if !hasRelation(barrelFile.Relations, barrelFile.File.ID, extra.ID, models.RelExports) {
		t.Fatal("expected export-star relation for extra")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, helper.ID, models.RelReferences) {
		t.Fatal("expected imported renamed binding to reference helper")
	}
	if !hasRelation(mainFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected renamed() call to resolve through barrel re-export")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, extra.ID, models.RelReferences) {
		t.Fatal("expected extra import to resolve through barrel export star")
	}
}

func TestTSAdapterComputesCyclomaticComplexity(t *testing.T) {
	t.Parallel()

	parsed := parseSingleTSLikeFile(t, "complexity.ts", `
function Simple(): number {
	return 1
}

function Branchy(flag: boolean, other: boolean): number {
	return flag && other ? 1 : 0
}
`)

	if got := mustFindSymbol(t, parsed.Symbols, "Simple").CyclomaticComplexity; got != 1 {
		t.Fatalf("Simple complexity = %d, want 1", got)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "Branchy").CyclomaticComplexity; got != 3 {
		t.Fatalf("Branchy complexity = %d, want 3", got)
	}
}

func TestTSAdapterProducesDiagnosticsForParseErrors(t *testing.T) {
	t.Parallel()

	parsed := parseSingleTSLikeFile(t, "broken.ts", `
export function Broken( {
`)

	if len(parsed.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(parsed.Diagnostics))
	}

	diagnostic := parsed.Diagnostics[0]
	if diagnostic.Code != tsParseErrorCode {
		t.Fatalf("diagnostic code = %q, want %q", diagnostic.Code, tsParseErrorCode)
	}
	if diagnostic.Stage != models.StageParse {
		t.Fatalf("diagnostic stage = %q, want %q", diagnostic.Stage, models.StageParse)
	}
	if len(parsed.Symbols) != 0 {
		t.Fatalf("expected no symbols on parse error, got %d", len(parsed.Symbols))
	}
}

func parseSingleTSLikeFile(t *testing.T, relativePath string, source string) models.ParsedFile {
	t.Helper()

	parsedFiles := parseTSLikeSources(t, map[string]string{relativePath: source})
	if len(parsedFiles) != 1 {
		t.Fatalf("expected 1 parsed file, got %d", len(parsedFiles))
	}

	return parsedFiles[0]
}

func parseTSLikeSources(t *testing.T, sources map[string]string) []models.ParsedFile {
	t.Helper()

	dir := t.TempDir()
	paths := make([]string, 0, len(sources))
	for relativePath := range sources {
		paths = append(paths, relativePath)
	}
	sort.Strings(paths)

	files := make([]models.ScannedSourceFile, 0, len(paths))
	for _, relativePath := range paths {
		absolutePath := filepath.Join(dir, relativePath)
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", relativePath, err)
		}

		source := strings.TrimLeft(sources[relativePath], "\n")
		if err := os.WriteFile(absolutePath, []byte(source), 0o644); err != nil {
			t.Fatalf("write %s: %v", relativePath, err)
		}

		files = append(files, models.ScannedSourceFile{
			AbsolutePath: absolutePath,
			RelativePath: relativePath,
			Language:     languageForPath(t, relativePath),
		})
	}

	parsedFiles, err := (TSAdapter{}).ParseFiles(files, dir)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	return parsedFiles
}

func mustFindParsedFile(t *testing.T, parsedFiles []models.ParsedFile, relativePath string) models.ParsedFile {
	t.Helper()

	for _, parsedFile := range parsedFiles {
		if parsedFile.File.FilePath == relativePath {
			return parsedFile
		}
	}

	t.Fatalf("missing parsed file %q", relativePath)
	return models.ParsedFile{}
}

func languageForPath(t *testing.T, relativePath string) models.SupportedLanguage {
	t.Helper()

	switch filepath.Ext(relativePath) {
	case ".ts":
		return models.LangTS
	case ".tsx":
		return models.LangTSX
	case ".js":
		return models.LangJS
	case ".jsx":
		return models.LangJSX
	default:
		t.Fatalf("unsupported test file extension for %s", relativePath)
		return ""
	}
}
