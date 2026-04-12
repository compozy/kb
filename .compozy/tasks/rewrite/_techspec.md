# Kodebase Go Port — Implementation Plan

> **For Hermes:** Use subagent-driven-development skill to implement this plan task-by-task.

**Goal:** Port the kodebase TypeScript CLI to Go — a tool that turns source code repositories into Karpathy-style Obsidian knowledge vaults with rich code metrics.

**Architecture:** Single-binary Go CLI with `cobra` for command routing. Internal packages mirror the TypeScript layers: models, scanner, adapters (tree-sitter), metrics engine, vault renderer/writer, and QMD shell integration. Follows kb conventions (config, logger, Makefile verify gate).

**Tech Stack:** Go 1.24, cobra (CLI), tree-sitter/go-tree-sitter + tree-sitter-typescript + tree-sitter-javascript + tree-sitter-go (parsing), slog (logging), BurntSushi/toml (config), mage (build).

**Reference:** `~/dev/projects/kodebase` (TypeScript source — ~6,700 LOC, 24 files)

---

## Phase 0: Project Bootstrap & CLI Skeleton

### Task 0.1: Rename module and update cmd entry point

**Objective:** Rename the Go module from `github.com/user/kb` to `github.com/pedronauck/kodebase` and update cmd to use cobra for subcommand routing.

**Files:**

- Modify: `go.mod`
- Modify: `cmd/app/main.go`
- Create: `internal/cli/root.go`
- Create: `internal/cli/generate.go`
- Create: `internal/cli/inspect.go`
- Create: `internal/cli/search.go`
- Create: `internal/cli/index.go`
- Create: `internal/cli/version.go`

**Step 1: Install cobra**

```bash
cd ~/dev/projects/kodebase-go
go get github.com/spf13/cobra@latest
```

**Step 2: Update go.mod**

```go
module github.com/pedronauck/kodebase

go 1.24
```

**Step 3: Create root command**

Create `internal/cli/root.go`:

```go
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kodebase",
	Short: "Turn source code repositories into Obsidian knowledge vaults",
	Long:  "Kodebase analyzes source code and generates Karpathy-style Obsidian knowledge vaults with rich code metrics.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Logger/config setup will go here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(indexCmd)
}
```

**Step 4: Create stub subcommands**

Each file in `internal/cli/` — generate.go, inspect.go, search.go, index.go, version.go — defines a `*cobra.Command` with `RunE` that prints "not implemented" and returns nil. Example for generate.go:

```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate <path>",
	Short: "Generate an Obsidian knowledge vault from source code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("generate: not implemented")
		return nil
	},
}
```

**Step 5: Wire main.go**

```go
package main

import "github.com/pedronauck/kodebase/internal/cli"

func main() {
	cli.Execute()
}
```

**Step 6: Verify**

```bash
cd ~/dev/projects/kodebase-go && go build ./... && ./app version
# Expected: version output
./app generate .
# Expected: "generate: not implemented"
./app inspect
# Expected: inspect subcommand help
```

**Step 7: Commit**

```bash
git add -A && git commit -m "feat: bootstrap project with cobra CLI skeleton"
```

---

## Phase 1: Core Models

### Task 1.1: Define domain models

**Objective:** Port all TypeScript interfaces/types to Go structs in a single models package.

**Files:**

- Create: `internal/kodebase/models/models.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/models.ts` (220 lines)

**Step 1: Write models.go**

Port all types:

- `SupportedLanguage` — string type with constants: LangTS, LangTSX, LangJS, LangJSX, LangGo
- `RelationType` — string type with constants: RelImports, RelExports, RelCalls, RelReferences, RelDeclares, RelContains
- `RelationConfidence` — string type: ConfidenceSemantic, ConfidenceSyntactic
- `DiagnosticSeverity` — string type: SeverityWarning, SeverityError
- `DiagnosticStage` — string type: StageScan, StageParse, StageRender, StageWrite, StageValidate
- `DocumentKind` — string type: DocRaw, DocWiki, DocIndex
- `ManagedArea` — string type: AreaRawCodebase, AreaWikiConcept, AreaWikiIndex
- `BaseViewType` — string type: ViewTable, ViewCards, ViewList
- `StructuredDiagnostic` struct
- `GraphFile` struct (id, nodeType, filePath, language, moduleDoc, symbolIds)
- `SymbolNode` struct (id, nodeType, name, symbolKind, language, filePath, startLine, endLine, signature, docComment, exported, cyclomaticComplexity)
- `ExternalNode` struct (id, nodeType, source, label)
- `RelationEdge` struct (fromId, toId, type, confidence)
- `ParsedFile` struct (file, symbols, externalNodes, relations, diagnostics)
- `GraphSnapshot` struct (rootPath, files, symbols, externalNodes, relations, diagnostics)
- `ScannedSourceFile` struct (absolutePath, relativePath, language)
- `ScannedWorkspace` struct (files, filesByLanguage)
- `RenderedDocument` struct (kind, managedArea, relativePath, frontmatter, body)
- `TopicMetadata` struct
- `GenerateOptions` struct
- `GenerationSummary` struct
- `LanguageAdapter` interface
- `SymbolMetrics` struct
- `FileMetrics` struct
- `DirectoryMetrics` struct
- `MetricsResult` struct
- `BaseFilter`, `BaseProperty`, `BaseGroupBy`, `BaseView`, `BaseDefinition`, `BaseFile` structs

**Step 2: Write test**

Create `internal/kodebase/models/models_test.go`:

```go
package models_test

import (
	"testing"

	"github.com/pedronauck/kodebase/internal/kodebase/models"
)

func TestSupportedLanguages(t *testing.T) {
	langs := models.SupportedLanguages()
	if len(langs) == 0 {
		t.Fatal("expected non-empty supported languages")
	}
}
```

**Step 3: Verify**

```bash
go test ./internal/kodebase/models/...
```

**Step 4: Commit**

```bash
git add -A && git commit -m "feat: add domain models package"
```

---

## Phase 2: Workspace Scanner

### Task 2.1: Implement workspace scanner

**Objective:** Port `scan-workspace.ts` — discover and filter source files using gitignore rules.

**Files:**

- Create: `internal/kodebase/scanner/scanner.go`
- Create: `internal/kodebase/scanner/scanner_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/scan-workspace.ts` (313 lines)

**Dependencies:** `github.com/sabhiram/go-gitignore` (gitignore parsing)

**Step 1: Install dependency**

```bash
go get github.com/sabhiram/go-gitignore@latest
```

**Step 2: Implement scanner**

Port logic:

- `ScanWorkspace(rootPath string, opts ScanOptions) (*ScannedWorkspace, error)`
- Load `.gitignore` from rootPath
- Apply default exclusions (node_modules, .git, dist, vendor, .kodebase, etc.)
- Apply user include/exclude patterns
- Walk directory tree, classify files by extension to `SupportedLanguage`
- Return `ScannedWorkspace` with files grouped by language

**Step 3: Write tests**

- Test with temp dir containing mixed files
- Test gitignore filtering
- Test include/exclude patterns

**Step 4: Verify**

```bash
go test ./internal/kodebase/scanner/... -v
```

**Step 5: Commit**

```bash
git add -A && git commit -m "feat: implement workspace scanner with gitignore support"
```

---

## Phase 3: Tree-Sitter Adapters

### Task 3.1: Set up tree-sitter infrastructure

**Objective:** Install tree-sitter Go bindings and grammar packages.

**Files:**

- Modify: `go.mod`

**Step 1: Install tree-sitter packages**

```bash
go get github.com/tree-sitter/go-tree-sitter@latest
go get github.com/tree-sitter/tree-sitter-go/bindings/go@latest
go get github.com/tree-sitter/tree-sitter-typescript/bindings/go@latest
go get github.com/tree-sitter/tree-sitter-javascript/bindings/go@latest
```

**Step 2: Verify compilation**

```bash
go build ./...
```

**Step 3: Commit**

```bash
git add -A && git commit -m "feat: add tree-sitter dependencies for Go, TS, JS parsing"
```

---

### Task 3.2: Implement Go language adapter

**Objective:** Port `go-adapter.ts` — parse Go source files using tree-sitter-go.

**Files:**

- Create: `internal/kodebase/adapter/go_adapter.go`
- Create: `internal/kodebase/adapter/go_adapter_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/adapters/go-adapter.ts` (317 lines)

**Step 1: Implement Go adapter**

Port logic:

- `type GoAdapter struct{}` implementing `LanguageAdapter` interface
- `Supports(lang)` returns true for LangGo
- `ParseFiles(files, rootPath)` parses each .go file:
  - Extract package declarations
  - Extract function/method declarations (name, signature, doc comment, exported, lines)
  - Extract type declarations (struct, interface, type alias)
  - Extract imports and build import relations
  - Extract function calls and build call relations
  - Compute cyclomatic complexity per function
  - Return `[]ParsedFile`

**Step 2: Write tests**

- Test parsing a simple Go file with functions and types
- Test import extraction
- Test call relation extraction
- Test complexity computation

**Step 3: Verify**

```bash
go test ./internal/kodebase/adapter/... -v -run TestGoAdapter
```

**Step 4: Commit**

```bash
git add -A && git commit -m "feat: implement Go language adapter with tree-sitter"
```

---

### Task 3.3: Implement TypeScript/JavaScript adapter

**Objective:** Port `oxc-typescript-adapter.ts` — parse TS/JS files using tree-sitter-typescript and tree-sitter-javascript.

**Files:**

- Create: `internal/kodebase/adapter/ts_adapter.go`
- Create: `internal/kodebase/adapter/ts_adapter_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/adapters/oxc-typescript-adapter.ts` (936 lines)

**Step 1: Implement TS adapter**

Port logic:

- `type TSAdapter struct{}` implementing `LanguageAdapter` interface
- `Supports(lang)` returns true for LangTS, LangTSX, LangJS, LangJSX
- `ParseFiles(files, rootPath)` parses each TS/JS file:
  - Select correct grammar (typescript vs tsx vs javascript)
  - Extract top-level declarations (functions, classes, interfaces, type aliases, variables, exports)
  - Build import bindings and resolve call targets
  - Construct cross-file relations (calls, references, imports)
  - Compute cyclomatic complexity
  - Return `[]ParsedFile`

**Step 2: Write tests**

- Test parsing a simple TypeScript file
- Test parsing a JSX file
- Test import/export extraction
- Test call relation resolution

**Step 3: Verify**

```bash
go test ./internal/kodebase/adapter/... -v -run TestTSAdapter
```

**Step 4: Commit**

```bash
git add -A && git commit -m "feat: implement TypeScript/JavaScript adapter with tree-sitter"
```

---

## Phase 4: Graph & Metrics

### Task 4.1: Implement graph normalizer

**Objective:** Port `normalize-graph.ts` — merge parsed files into a single GraphSnapshot.

**Files:**

- Create: `internal/kodebase/graph/normalize.go`
- Create: `internal/kodebase/graph/normalize_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/normalize-graph.ts` (114 lines)

**Step 1: Implement normalizer**

- `NormalizeGraph(rootPath string, parsedFiles []ParsedFile) *GraphSnapshot`
- Deduplicate files, symbols, external nodes, relations
- Attach symbol IDs to parent files
- Order deterministically

**Step 2: Write tests**

**Step 3: Commit**

```bash
git add -A && git commit -m "feat: implement graph normalizer"
```

---

### Task 4.2: Implement metrics engine

**Objective:** Port `compute-metrics.ts` — compute symbol, file, and directory metrics.

**Files:**

- Create: `internal/kodebase/metrics/compute.go`
- Create: `internal/kodebase/metrics/compute_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/compute-metrics.ts` (480 lines)

**Step 1: Implement metrics computation**

Port all metric calculations:

- **Per-symbol:** blast radius (BFS), centrality (PageRank-like), direct dependents, external reference count, dead export detection, long function, code smells
- **Per-file:** afferent/efferent coupling, instability, orphan/god file detection, circular dependency participation
- **Per-directory:** coupling aggregation
- **Circular dependencies:** DFS cycle detection

**Step 2: Write tests**

- Test blast radius computation
- Test dead export detection
- Test coupling metrics
- Test circular dependency detection

**Step 3: Commit**

```bash
git add -A && git commit -m "feat: implement metrics engine (blast radius, coupling, smells, cycles)"
```

---

## Phase 5: Vault Rendering

### Task 5.1: Implement path utilities

**Objective:** Port `path-utils.ts` — path manipulation and vault path derivation.

**Files:**

- Create: `internal/kodebase/vault/pathutils.go`
- Create: `internal/kodebase/vault/pathutils_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/path-utils.ts` (106 lines)

---

### Task 5.2: Implement text utilities

**Objective:** Port `text-utils.ts` — comment extraction and normalization.

**Files:**

- Create: `internal/kodebase/vault/textutils.go`
- Create: `internal/kodebase/vault/textutils_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/text-utils.ts` (36 lines)

---

### Task 5.3: Implement document renderer

**Objective:** Port `render-documents.ts` — render Obsidian vault documents (raw, wiki, index, base).

**Files:**

- Create: `internal/kodebase/vault/render.go`
- Create: `internal/kodebase/vault/render_wiki.go`
- Create: `internal/kodebase/vault/render_base.go`
- Create: `internal/kodebase/vault/render_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/render-documents.ts` (1692 lines — largest file)

This is the biggest single task. Break into sub-files:

- `render.go` — orchestrator: `RenderKnowledgeBaseDocuments(graph, metrics, topic) []RenderedDocument`
- `render_wiki.go` — 10 wiki concept articles + 3 index pages
- `render_base.go` — 12 Obsidian Base definition files

Each document has YAML frontmatter + Markdown body with `[[wiki-links]]`.

**Step 1: Implement render.go (raw documents)**

- Raw file documents (one per source file)
- Raw symbol documents (one per symbol)
- Raw directory index documents
- Raw language index documents

**Step 2: Implement render_wiki.go (wiki documents)**

- 10 starter wiki concepts: Codebase Overview, Directory Map, Symbol Taxonomy, Dependency Hotspots, Complexity Hotspots, Module Health, Dead Code Report, Code Smells, Circular Dependencies, High-Impact Symbols
- 3 wiki index pages: Dashboard, Concept Index, Source Index

**Step 3: Implement render_base.go (Obsidian Base files)**

- 12 base definition files for Obsidian

**Step 4: Write tests**

- Test raw document generation
- Test wiki link formatting
- Test frontmatter generation

**Step 5: Commit**

```bash
git add -A && git commit -m "feat: implement vault document renderer (raw, wiki, base)"
```

---

### Task 5.4: Implement vault writer

**Objective:** Port `write-vault.ts` — write the full Obsidian vault to disk.

**Files:**

- Create: `internal/kodebase/vault/writer.go`
- Create: `internal/kodebase/vault/writer_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/write-vault.ts` (265 lines)

**Step 1: Implement writer**

- Create topic directory skeleton (raw/, wiki/, outputs/, bases/)
- Write rendered documents in batches
- Manage CLAUDE.md (topic manifest)
- Manage log.md (append-only audit log)
- Remove stale managed wiki concepts

**Step 2: Write tests using t.TempDir()**

**Step 3: Commit**

```bash
git add -A && git commit -m "feat: implement vault writer"
```

---

## Phase 6: Generate Command

### Task 6.1: Wire generate command

**Objective:** Connect scanner, adapters, graph normalizer, metrics, and vault writer into the `generate` CLI command.

**Files:**

- Modify: `internal/cli/generate.go`
- Create: `internal/kodebase/generate.go`
- Create: `internal/kodebase/generate_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/generate-knowledge-base.ts` (114 lines)

**Step 1: Implement generate orchestrator**

```go
func Generate(ctx context.Context, opts GenerateOptions) (*GenerationSummary, error)
```

Pipeline:

1. Scan workspace
2. Select adapters based on languages found
3. Parse files with adapters
4. Normalize graph
5. Compute metrics
6. Render documents
7. Write vault
8. Return summary

**Step 2: Wire into CLI**

Add flags to generate command: --output, --topic, --title, --domain, --include, --exclude, --semantic

**Step 3: Integration test**

Test against a small fixture Go project.

**Step 4: Commit**

```bash
git add -A && git commit -m "feat: wire generate command end-to-end"
```

---

## Phase 7: Vault Reader & Inspect Commands

### Task 7.1: Implement vault reader

**Objective:** Port `vault-reader.ts` — read generated vault back into structured data.

**Files:**

- Create: `internal/kodebase/vault/reader.go`
- Create: `internal/kodebase/vault/reader_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/vault-reader.ts` (207 lines)

---

### Task 7.2: Implement vault query resolver

**Objective:** Port `vault-query.ts` — resolve vault and topic paths.

**Files:**

- Create: `internal/kodebase/vault/query.go`
- Create: `internal/kodebase/vault/query_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/vault-query.ts` (109 lines)

---

### Task 7.3: Implement output formatter

**Objective:** Port `output-formatter.ts` — format tabular data as ASCII table, JSON, or TSV.

**Files:**

- Create: `internal/kodebase/output/formatter.go`
- Create: `internal/kodebase/output/formatter_test.go`

**Reference:** `kodebase/packages/cli/src/knowledge-base/output-formatter.ts` (84 lines)

---

### Task 7.4: Wire inspect command with all subcommands

**Objective:** Port all 9 inspect subcommands: smells, dead-code, complexity, blast-radius, coupling, symbol, file, backlinks, deps, circular-deps.

**Files:**

- Modify: `internal/cli/inspect.go`
- Create: `internal/cli/inspect_smells.go`
- Create: `internal/cli/inspect_deadcode.go`
- Create: `internal/cli/inspect_complexity.go`
- Create: `internal/cli/inspect_blastradius.go`
- Create: `internal/cli/inspect_coupling.go`
- Create: `internal/cli/inspect_symbol.go`
- Create: `internal/cli/inspect_file.go`
- Create: `internal/cli/inspect_backlinks.go`
- Create: `internal/cli/inspect_deps.go`
- Create: `internal/cli/inspect_circulardeps.go`

**Reference:** `kodebase/packages/cli/src/commands/inspect/` (744 lines total)

Each subcommand:

1. Resolves vault/topic via vault-query
2. Reads snapshot via vault-reader
3. Formats output via formatter
4. Supports --format (table|json|tsv), --vault, --topic flags

**Step 1: Implement shared inspect infrastructure**

**Step 2: Implement each subcommand**

**Step 3: Wire into cobra**

**Step 4: Commit**

```bash
git add -A && git commit -m "feat: implement inspect command with 9 subcommands"
```

---

## Phase 8: QMD Integration (Shell)

### Task 8.1: Implement QMD shell client

**Objective:** Port `qmd-client.ts` — call QMD CLI via shell for search and indexing.

**Files:**

- Create: `internal/kodebase/qmd/client.go`
- Create: `internal/kodebase/qmd/client_test.go`

**Reference:** `kodebase/packages/cli/src/integrations/qmd-client.ts` (290 lines)

**Step 1: Implement QMD client**

```go
type QMDClient struct {
    binPath string // path to qmd binary, default "qmd"
}

func (c *QMDClient) Index(ctx context.Context, opts IndexOptions) error
func (c *QMDClient) Search(ctx context.Context, opts SearchOptions) ([]SearchResult, error)
```

Uses `os/exec` to call:

- `qmd collection add <path>` — index
- `qmd collection update <path>` — update
- `qmd search --mode hybrid|lexical|vector <query>` — search

Graceful error when QMD not installed (return `ErrQMDUnavailable`).

**Step 2: Write tests (mock exec)**

**Step 3: Commit**

```bash
git add -A && git commit -m "feat: implement QMD shell client"
```

---

### Task 8.2: Wire search command

**Objective:** Port `search.ts` — hybrid/lexical/vector search via QMD.

**Files:**

- Modify: `internal/cli/search.go`

**Reference:** `kodebase/packages/cli/src/commands/search.ts` (251 lines)

**Step 1: Wire search command**

Add flags: --lex, --vec, --limit, --min-score, --full, --all, --collection

**Step 2: Commit**

```bash
git add -A && git commit -m "feat: wire search command with QMD integration"
```

---

### Task 8.3: Wire index-vault command

**Objective:** Port `index-vault.ts` — create/update QMD collection for a vault topic.

**Files:**

- Modify: `internal/cli/index.go`

**Reference:** `kodebase/packages/cli/src/commands/index-vault.ts` (135 lines)

**Step 1: Wire index command**

Add flags: --embed, --force-embed, --context

**Step 2: Commit**

```bash
git add -A && git commit -m "feat: wire index-vault command with QMD integration"
```

---

## Phase 9: Polish & Integration

### Task 9.1: Add integration test

**Objective:** End-to-end test: generate vault from a fixture project, then inspect it.

**Files:**

- Create: `internal/kodebase/generate_integration_test.go`

**Step 1: Create fixture project in testdata**

Small Go project with:

- Multiple packages
- Circular imports
- Dead exports
- Complex functions

**Step 2: Write integration test**

```go
func TestGenerateAndInspect(t *testing.T) {
    // 1. Generate vault from fixture
    // 2. Verify vault files exist
    // 3. Run inspect commands against vault
    // 4. Verify output
}
```

**Step 3: Commit**

```bash
git add -A && git commit -m "test: add end-to-end integration test"
```

---

### Task 9.2: Update AGENTS.md and config

**Objective:** Update project documentation to reflect kodebase purpose.

**Files:**

- Modify: `AGENTS.md`
- Modify: `CLAUDE.md`
- Modify: `config.example.toml`

**Step 1: Update documentation**

- Project overview: "Kodebase CLI — Turn source code into Obsidian knowledge vaults"
- Package layout table with all new packages
- Commands section with all kodebase commands

**Step 2: Commit**

```bash
git add -A && git commit -m "docs: update AGENTS.md and config for kodebase"
```

---

### Task 9.3: Final verify and push

**Objective:** Run full verification pipeline and push to GitHub.

**Step 1: Run make verify**

```bash
cd ~/dev/projects/kodebase-go && make verify
```

Must pass: fmt + lint + test + build with zero issues.

**Step 2: Fix any issues**

**Step 3: Create GitHub repo and push**

```bash
gh repo create pedronauck/kodebase-go --private --source=. --push
```

---

## Summary

| Phase               | Tasks        | Est. LOC   | Key Dependencies                                                               |
| ------------------- | ------------ | ---------- | ------------------------------------------------------------------------------ |
| 0. Bootstrap        | 1            | ~200       | cobra                                                                          |
| 1. Models           | 1            | ~300       | —                                                                              |
| 2. Scanner          | 1            | ~250       | go-gitignore                                                                   |
| 3. Adapters         | 3            | ~800       | go-tree-sitter, tree-sitter-go, tree-sitter-typescript, tree-sitter-javascript |
| 4. Graph & Metrics  | 2            | ~600       | —                                                                              |
| 5. Vault Rendering  | 4            | ~1500      | —                                                                              |
| 6. Generate Command | 1            | ~200       | —                                                                              |
| 7. Inspect Commands | 4            | ~600       | —                                                                              |
| 8. QMD Integration  | 3            | ~300       | os/exec                                                                        |
| 9. Polish           | 3            | ~100       | —                                                                              |
| **Total**           | **23 tasks** | **~4,850** |                                                                                |

**Parallelizable workstreams** (can run simultaneously via `$team`):

- Stream A: Phases 0-1-2 (scaffold + models + scanner)
- Stream B: Phase 3 (adapters — can start after Phase 1 models)
- Stream C: Phases 4-5 (graph + vault — can start after Phase 1 models)
- Stream D: Phases 7-8 (inspect + QMD — can start after Phase 5 vault)
