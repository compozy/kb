package adapter

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/user/kb/internal/models"
)

func TestGoAdapterSupportsOnlyGo(t *testing.T) {
	t.Parallel()

	adapter := GoAdapter{}

	if !adapter.Supports(models.LangGo) {
		t.Fatal("expected GoAdapter to support Go")
	}

	for _, language := range []models.SupportedLanguage{models.LangTS, models.LangTSX, models.LangJS, models.LangJSX} {
		if adapter.Supports(language) {
			t.Fatalf("expected GoAdapter to reject %q", language)
		}
	}
}

func TestGoAdapterParsesSimpleGoFile(t *testing.T) {
	t.Parallel()

	parsedFiles := parseGoSources(t, map[string]string{
		"sample.go": `
// Package sample demonstrates adapter output.
package sample

// Add sums two integers.
func Add(left int, right int) int {
	return left + right
}
`,
	})

	if len(parsedFiles) != 1 {
		t.Fatalf("expected 1 parsed file, got %d", len(parsedFiles))
	}

	parsed := parsedFiles[0]
	if parsed.File.ID != "file:sample.go" {
		t.Fatalf("file ID = %q, want %q", parsed.File.ID, "file:sample.go")
	}

	if parsed.File.ModuleDoc != "Package sample demonstrates adapter output." {
		t.Fatalf("module doc = %q", parsed.File.ModuleDoc)
	}

	if len(parsed.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(parsed.Diagnostics))
	}

	if len(parsed.File.SymbolIDs) != 2 {
		t.Fatalf("expected package and function symbols, got %d", len(parsed.File.SymbolIDs))
	}

	packageSymbol := mustFindSymbol(t, parsed.Symbols, "sample")
	if packageSymbol.SymbolKind != "package" {
		t.Fatalf("package symbol kind = %q, want %q", packageSymbol.SymbolKind, "package")
	}

	functionSymbol := mustFindSymbol(t, parsed.Symbols, "Add")
	if functionSymbol.SymbolKind != "function" {
		t.Fatalf("function symbol kind = %q, want %q", functionSymbol.SymbolKind, "function")
	}

	if functionSymbol.Signature != "func Add(left int, right int) int {" {
		t.Fatalf("signature = %q", functionSymbol.Signature)
	}

	if functionSymbol.DocComment != "Add sums two integers." {
		t.Fatalf("doc comment = %q", functionSymbol.DocComment)
	}

	if !hasRelation(parsed.Relations, parsed.File.ID, functionSymbol.ID, models.RelContains) {
		t.Fatal("expected file contains relation for Add")
	}
}

func TestGoAdapterSetsExportedFlags(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

func Exported() {}

func hidden() {}
`)

	if !mustFindSymbol(t, parsed.Symbols, "Exported").Exported {
		t.Fatal("expected Exported to be marked exported")
	}

	if mustFindSymbol(t, parsed.Symbols, "hidden").Exported {
		t.Fatal("expected hidden to be marked unexported")
	}
}

func TestGoAdapterExtractsTypeDeclarations(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

type Widget struct{}

type Runner interface {
	Run()
}

type Identifier = string
`)

	if got := mustFindSymbol(t, parsed.Symbols, "Widget").SymbolKind; got != "struct" {
		t.Fatalf("Widget kind = %q, want %q", got, "struct")
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Runner").SymbolKind; got != "interface" {
		t.Fatalf("Runner kind = %q, want %q", got, "interface")
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Identifier").SymbolKind; got != "type" {
		t.Fatalf("Identifier kind = %q, want %q", got, "type")
	}
}

func TestGoAdapterExtractsMethodSignatureWithReceiver(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

type Widget struct{}

func (w *Widget) Run() {}
`)

	method := mustFindSymbol(t, parsed.Symbols, "Run")
	if method.SymbolKind != "method" {
		t.Fatalf("method kind = %q, want %q", method.SymbolKind, "method")
	}

	if !strings.Contains(method.Signature, "(w *Widget) Run()") {
		t.Fatalf("method signature = %q, want receiver information", method.Signature)
	}
}

func TestGoAdapterExtractsImportRelations(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

import (
	"fmt"
	stringsAlias "strings"
)

func Use() {
	fmt.Println(stringsAlias.ToUpper("ok"))
}
`)

	if len(parsed.ExternalNodes) != 2 {
		t.Fatalf("expected 2 external nodes, got %d", len(parsed.ExternalNodes))
	}

	fmtNode := mustFindExternalNode(t, parsed.ExternalNodes, "fmt")
	if fmtNode.Label != "fmt" {
		t.Fatalf("fmt label = %q, want %q", fmtNode.Label, "fmt")
	}

	stringsNode := mustFindExternalNode(t, parsed.ExternalNodes, "strings")
	if stringsNode.Label != "stringsAlias (strings)" {
		t.Fatalf("strings label = %q", stringsNode.Label)
	}

	if !hasRelation(parsed.Relations, parsed.File.ID, fmtNode.ID, models.RelImports) {
		t.Fatal("expected import relation for fmt")
	}

	if !hasRelation(parsed.Relations, parsed.File.ID, stringsNode.ID, models.RelImports) {
		t.Fatal("expected import relation for strings")
	}
}

func TestGoAdapterExtractsCallRelationsForDirectIdentifiers(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

func helper() {}

func Run() {
	helper()
	println("ignored selector replacement")
}
`)

	helper := mustFindSymbol(t, parsed.Symbols, "helper")
	run := mustFindSymbol(t, parsed.Symbols, "Run")

	if !hasRelation(parsed.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected direct identifier call relation from Run to helper")
	}
}

func TestGoAdapterComputesCyclomaticComplexity(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

func Linear() {
	println("ok")
}

func Branching(flag bool, other bool, total int) {
	if flag {
		println("flag")
	} else if other {
		println("other")
	}

	for i := 0; i < total; i++ {
		println(i)
	}
}

func Logical(flag bool, other bool) {
	if flag && other {
		println("and")
	}
}
`)

	if got := mustFindSymbol(t, parsed.Symbols, "Linear").CyclomaticComplexity; got != 1 {
		t.Fatalf("Linear complexity = %d, want 1", got)
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Branching").CyclomaticComplexity; got != 4 {
		t.Fatalf("Branching complexity = %d, want 4", got)
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Logical").CyclomaticComplexity; got != 3 {
		t.Fatalf("Logical complexity = %d, want 3", got)
	}
}

func TestGoAdapterProducesDiagnosticsForParseErrors(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

func Broken( {
`)

	if len(parsed.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(parsed.Diagnostics))
	}

	diagnostic := parsed.Diagnostics[0]
	if diagnostic.Code != goParseErrorCode {
		t.Fatalf("diagnostic code = %q, want %q", diagnostic.Code, goParseErrorCode)
	}

	if diagnostic.Stage != models.StageParse {
		t.Fatalf("diagnostic stage = %q, want %q", diagnostic.Stage, models.StageParse)
	}

	if len(parsed.Symbols) != 0 {
		t.Fatalf("expected parse-error file to have no symbols, got %d", len(parsed.Symbols))
	}
}

func TestGoAdapterExtractsDocComments(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

// Widget stores parsed state.
type Widget struct{}

/*
Run executes the widget.
It supports multiple lines.
*/
func (w *Widget) Run() {}
`)

	if got := mustFindSymbol(t, parsed.Symbols, "Widget").DocComment; got != "Widget stores parsed state." {
		t.Fatalf("Widget doc comment = %q", got)
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Run").DocComment; got != "Run executes the widget.\nIt supports multiple lines." {
		t.Fatalf("Run doc comment = %q", got)
	}
}

func TestGoAdapterModuleDocUsesOnlyLeadingComment(t *testing.T) {
	t.Parallel()

	parsed := parseSingleGoFile(t, `
package sample

// Run executes the sample.
func Run() {}
`)

	if parsed.File.ModuleDoc != "" {
		t.Fatalf("module doc = %q, want empty string", parsed.File.ModuleDoc)
	}

	if got := mustFindSymbol(t, parsed.Symbols, "Run").DocComment; got != "Run executes the sample." {
		t.Fatalf("Run doc comment = %q", got)
	}
}

func parseSingleGoFile(t *testing.T, source string) models.ParsedFile {
	t.Helper()

	parsedFiles := parseGoSources(t, map[string]string{"sample.go": source})
	if len(parsedFiles) != 1 {
		t.Fatalf("expected 1 parsed file, got %d", len(parsedFiles))
	}

	return parsedFiles[0]
}

func parseGoSources(t *testing.T, sources map[string]string) []models.ParsedFile {
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
			Language:     models.LangGo,
		})
	}

	parsedFiles, err := (GoAdapter{}).ParseFiles(files, dir)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	return parsedFiles
}

func mustFindSymbol(t *testing.T, symbols []models.SymbolNode, name string) models.SymbolNode {
	t.Helper()

	for _, symbol := range symbols {
		if symbol.Name == name {
			return symbol
		}
	}

	t.Fatalf("missing symbol %q", name)
	return models.SymbolNode{}
}

func mustFindExternalNode(t *testing.T, nodes []models.ExternalNode, source string) models.ExternalNode {
	t.Helper()

	for _, node := range nodes {
		if node.Source == source {
			return node
		}
	}

	t.Fatalf("missing external node %q", source)
	return models.ExternalNode{}
}

func hasRelation(relations []models.RelationEdge, fromID string, toID string, relationType models.RelationType) bool {
	for _, relation := range relations {
		if relation.FromID == fromID && relation.ToID == toID && relation.Type == relationType {
			return true
		}
	}

	return false
}
