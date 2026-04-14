package adapter

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestRustAdapterSupportsOnlyRust(t *testing.T) {
	t.Parallel()

	adapter := RustAdapter{}
	if !adapter.Supports(models.LangRust) {
		t.Fatal("expected RustAdapter to support Rust")
	}

	for _, language := range []models.SupportedLanguage{models.LangTS, models.LangTSX, models.LangJS, models.LangJSX, models.LangGo} {
		if adapter.Supports(language) {
			t.Fatalf("expected RustAdapter to reject %q", language)
		}
	}
}

func TestRustAdapterParsesSimpleRustFile(t *testing.T) {
	t.Parallel()

	parsed := parseSingleRustFile(t, "src/lib.rs", `
//! crate docs

/// add docs
pub fn add(left: i32, right: i32) -> i32 {
    left + right
}

pub struct Widget;
`)

	if parsed.File.ID != "file:src/lib.rs" {
		t.Fatalf("file ID = %q, want file:src/lib.rs", parsed.File.ID)
	}
	if parsed.File.ModuleDoc != "crate docs" {
		t.Fatalf("module doc = %q, want crate docs", parsed.File.ModuleDoc)
	}

	add := mustFindSymbol(t, parsed.Symbols, "add")
	if add.SymbolKind != rustSymbolKindFunction {
		t.Fatalf("add kind = %q, want %q", add.SymbolKind, rustSymbolKindFunction)
	}
	if !add.Exported {
		t.Fatal("expected add to be exported")
	}
	if add.DocComment != "add docs" {
		t.Fatalf("add doc comment = %q, want add docs", add.DocComment)
	}

	widget := mustFindSymbol(t, parsed.Symbols, "Widget")
	if widget.SymbolKind != rustSymbolKindStruct {
		t.Fatalf("Widget kind = %q, want %q", widget.SymbolKind, rustSymbolKindStruct)
	}
	if !widget.Exported {
		t.Fatal("expected Widget to be exported")
	}

	if !hasRelation(parsed.Relations, parsed.File.ID, add.ID, models.RelContains) {
		t.Fatal("expected file contains relation for add")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, add.ID, models.RelExports) {
		t.Fatal("expected file exports relation for add")
	}
}

func TestRustAdapterParsesTypeDeclarations(t *testing.T) {
	t.Parallel()

	parsed := parseSingleRustFile(t, "src/lib.rs", `
pub enum State {
    Ready,
}

pub trait Runner {
    fn run(&self);
}

pub type Identifier = String;

pub union MaybeInt {
    value: i32,
}

pub const ANSWER: i32 = 42;
pub static NAME: &str = "kb";

macro_rules! helper_macro {
    () => {};
}
`)

	if got := mustFindSymbol(t, parsed.Symbols, "State").SymbolKind; got != rustSymbolKindEnum {
		t.Fatalf("State kind = %q, want %q", got, rustSymbolKindEnum)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "Runner").SymbolKind; got != rustSymbolKindTrait {
		t.Fatalf("Runner kind = %q, want %q", got, rustSymbolKindTrait)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "Identifier").SymbolKind; got != rustSymbolKindTypeAlias {
		t.Fatalf("Identifier kind = %q, want %q", got, rustSymbolKindTypeAlias)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "MaybeInt").SymbolKind; got != rustSymbolKindUnion {
		t.Fatalf("MaybeInt kind = %q, want %q", got, rustSymbolKindUnion)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "ANSWER").SymbolKind; got != rustSymbolKindConst {
		t.Fatalf("ANSWER kind = %q, want %q", got, rustSymbolKindConst)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "NAME").SymbolKind; got != rustSymbolKindStatic {
		t.Fatalf("NAME kind = %q, want %q", got, rustSymbolKindStatic)
	}
	if got := mustFindSymbol(t, parsed.Symbols, "helper_macro").SymbolKind; got != rustSymbolKindMacro {
		t.Fatalf("helper_macro kind = %q, want %q", got, rustSymbolKindMacro)
	}
}

func TestRustAdapterBuildsMethodAndCallRelations(t *testing.T) {
	t.Parallel()

	parsed := parseSingleRustFile(t, "src/lib.rs", `
struct Widget;

fn helper() {}

impl Widget {
    /// run docs
    pub fn run(&self) {
        helper();
    }
}
`)

	run := mustFindSymbol(t, parsed.Symbols, "run")
	if run.SymbolKind != rustSymbolKindMethod {
		t.Fatalf("run kind = %q, want %q", run.SymbolKind, rustSymbolKindMethod)
	}
	if !run.Exported {
		t.Fatal("expected run to be exported")
	}
	if run.DocComment != "run docs" {
		t.Fatalf("run doc comment = %q, want run docs", run.DocComment)
	}

	helper := mustFindSymbol(t, parsed.Symbols, "helper")
	if !hasRelation(parsed.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected call relation from run to helper")
	}
}

func TestRustAdapterResolvesUseAndPubUseAcrossFiles(t *testing.T) {
	t.Parallel()

	parsedFiles := parseRustSources(t, map[string]string{
		"Cargo.toml": `
[package]
name = "fixture"
version = "0.1.0"
edition = "2021"
`,
		"src/util.rs": `
pub struct Helper;

pub fn helper() {}

impl Helper {
    pub fn new() -> Self {
        Self
    }
}
`,
		"src/lib.rs": `
pub mod util;

pub use crate::util::{Helper as Alias, helper};

pub fn run() {
    helper();
    Alias::new();
}
`,
	})

	libFile := mustFindParsedFile(t, parsedFiles, "src/lib.rs")
	utilFile := mustFindParsedFile(t, parsedFiles, "src/util.rs")
	helperFn := mustFindSymbol(t, utilFile.Symbols, "helper")
	newMethod := mustFindSymbol(t, utilFile.Symbols, "new")
	run := mustFindSymbol(t, libFile.Symbols, "run")
	moduleSymbol := mustFindSymbol(t, libFile.Symbols, "util")

	if moduleSymbol.SymbolKind != rustSymbolKindModule {
		t.Fatalf("util symbol kind = %q, want %q", moduleSymbol.SymbolKind, rustSymbolKindModule)
	}
	if !hasRelation(libFile.Relations, libFile.File.ID, utilFile.File.ID, models.RelImports) {
		t.Fatal("expected lib.rs to import util.rs")
	}
	if !hasRelation(libFile.Relations, libFile.File.ID, helperFn.ID, models.RelReferences) {
		t.Fatal("expected lib.rs to reference helper() through use")
	}
	if !hasRelation(libFile.Relations, libFile.File.ID, helperFn.ID, models.RelExports) {
		t.Fatal("expected lib.rs to re-export helper() through pub use")
	}
	if !hasRelation(libFile.Relations, run.ID, helperFn.ID, models.RelCalls) {
		t.Fatal("expected run() to call helper() through use")
	}
	if !hasRelation(libFile.Relations, run.ID, newMethod.ID, models.RelCalls) {
		t.Fatal("expected run() to call Alias::new() through pub use alias")
	}
}

func TestRustAdapterResolvesPrivateChildModuleCallsWithinCrate(t *testing.T) {
	t.Parallel()

	parsedFiles := parseRustSources(t, map[string]string{
		"Cargo.toml": `
[package]
name = "fixture"
version = "0.1.0"
edition = "2021"
`,
		"src/util.rs": `
fn helper() {}
`,
		"src/lib.rs": `
mod util;

fn run() {
    util::helper();
}
`,
	})

	libFile := mustFindParsedFile(t, parsedFiles, "src/lib.rs")
	utilFile := mustFindParsedFile(t, parsedFiles, "src/util.rs")
	run := mustFindSymbol(t, libFile.Symbols, "run")
	helper := mustFindSymbol(t, utilFile.Symbols, "helper")

	if !hasRelation(libFile.Relations, libFile.File.ID, utilFile.File.ID, models.RelImports) {
		t.Fatal("expected lib.rs to import util.rs")
	}
	if !hasRelation(libFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected run() to call private util::helper() within the same crate")
	}
}

func TestRustAdapterResolvesPrivateUseBindingsWithinCrate(t *testing.T) {
	t.Parallel()

	parsedFiles := parseRustSources(t, map[string]string{
		"Cargo.toml": `
[package]
name = "fixture"
version = "0.1.0"
edition = "2021"
`,
		"src/util.rs": `
fn helper() {}
`,
		"src/lib.rs": `
mod util;

use crate::util::helper;

fn run() {
    helper();
}
`,
	})

	libFile := mustFindParsedFile(t, parsedFiles, "src/lib.rs")
	utilFile := mustFindParsedFile(t, parsedFiles, "src/util.rs")
	run := mustFindSymbol(t, libFile.Symbols, "run")
	helper := mustFindSymbol(t, utilFile.Symbols, "helper")

	if !hasRelation(libFile.Relations, libFile.File.ID, helper.ID, models.RelReferences) {
		t.Fatal("expected lib.rs to reference private helper() through use")
	}
	if !hasRelation(libFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected run() to call private helper() through use within the same crate")
	}
}

func TestRustAdapterProducesDiagnosticsForParseErrors(t *testing.T) {
	t.Parallel()

	parsed := parseSingleRustFile(t, "src/lib.rs", `
fn broken( {
`)

	if len(parsed.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(parsed.Diagnostics))
	}
	if parsed.Diagnostics[0].Code != rustParseErrorCode {
		t.Fatalf("diagnostic code = %q, want %q", parsed.Diagnostics[0].Code, rustParseErrorCode)
	}
	if len(parsed.Symbols) != 0 {
		t.Fatalf("expected no symbols for parse-error file, got %d", len(parsed.Symbols))
	}
}

func parseSingleRustFile(t *testing.T, relativePath string, source string) models.ParsedFile {
	t.Helper()

	parsedFiles := parseRustSources(t, map[string]string{relativePath: source})
	if len(parsedFiles) != 1 {
		t.Fatalf("expected 1 parsed file, got %d", len(parsedFiles))
	}

	return parsedFiles[0]
}

func parseRustSources(t *testing.T, sources map[string]string) []models.ParsedFile {
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

		if filepath.Ext(relativePath) != ".rs" {
			continue
		}

		files = append(files, models.ScannedSourceFile{
			AbsolutePath: absolutePath,
			RelativePath: relativePath,
			Language:     models.LangRust,
		})
	}

	parsedFiles, err := (RustAdapter{}).ParseFiles(files, dir)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	return parsedFiles
}
