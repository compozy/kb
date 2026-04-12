# Contributing to kb

Thanks for your interest in contributing. Whether it's a bug report, a new language adapter, a file converter, or a documentation fix, every contribution helps.

---

## Prerequisites

- [Go](https://go.dev) >= 1.24
- [Git](https://git-scm.com)

Optional (auto-downloaded by the build system if missing):

- [golangci-lint](https://golangci-lint.run) v2
- [gotestsum](https://github.com/gotestyourself/gotestsum)

---

## Getting Started

```bash
git clone https://github.com/pedronauck/kodebase-go.git
cd kodebase-go
make verify   # Verify everything passes before making changes
```

### Project Structure

```text
cmd/
  kb/
    main.go                       # Program entrypoint
internal/
  cli/                            # Cobra command tree and command adapters
    root.go
    topic.go                      # kb topic {new,list,info}
    ingest.go                     # kb ingest parent command + shared helpers
    ingest_url.go                 # kb ingest url
    ingest_file.go                # kb ingest file
    ingest_youtube.go             # kb ingest youtube
    ingest_codebase.go            # kb ingest codebase
    ingest_bookmarks.go           # kb ingest bookmarks
    lint.go                       # kb lint
    generate.go                   # Hidden compatibility alias
    inspect.go                    # Router for inspect subcommands
    inspect_*.go                  # Inspect subcommand implementations
    search.go
    index.go
    version.go
  topic/                          # Topic scaffolding and metadata
  ingest/                         # Ingest orchestration and vault writes
  convert/                        # Converter registry and format-specific converters
  firecrawl/                      # Firecrawl REST API client
  youtube/                        # YouTube transcript extraction + STT fallback
  frontmatter/                    # YAML frontmatter parsing and generation
  lint/                           # KB structural lint engine
  generate/                       # Codebase-to-vault pipeline orchestration
  models/                         # Domain types, snapshots, and interfaces
  scanner/                        # Workspace discovery and ignore filtering
  adapter/                        # Tree-sitter parsing adapters
    go_adapter.go
    ts_adapter.go
    treesitter.go
  graph/                          # Graph normalization
  metrics/                        # File, symbol, and directory metrics
  vault/                          # Rendering, writing, reading, query helpers
  qmd/                            # QMD shell client integration
  output/                         # Table / JSON / TSV output rendering
  config/                         # TOML config and env-backed secrets
  logger/                         # Structured slog logger
  version/                        # Build metadata
magefile.go                       # Mage build tasks (wrapped by Makefile)
```

---

## Development Workflow

| Command                 | Description                               |
| ----------------------- | ----------------------------------------- |
| `make fmt`              | Format all Go files with gofmt            |
| `make lint`             | Run golangci-lint v2 with zero tolerance  |
| `make test`             | Unit tests with race detector             |
| `make test-integration` | Unit + integration tests                  |
| `make build`            | Build binary to `bin/kb` with ldflags     |
| `make verify`           | fmt -> lint -> test -> build -> boundaries|
| `make deps`             | Run `go mod tidy`                         |

**`make verify` must pass before submitting a PR.** The CI pipeline runs the same command.

---

## Code Style

- **File naming:** `snake_case.go` for all Go files
- **Exports:** Capitalize public symbols; keep internal logic unexported
- **Formatting:** `gofmt` (standard Go formatting, enforced by `make fmt`)
- **Linting:** golangci-lint v2 with zero warnings -- warnings are treated as errors
- **Imports:** Group in order: stdlib, third-party, internal
- **Dependencies:** Use `go get` for dependency changes, never hand-edit `go.mod`
- **CLI commands:** Use [Cobra](https://github.com/spf13/cobra). Keep commands thin -- delegate to packages like `internal/topic`, `internal/ingest`, `internal/lint`, `internal/generate`, `internal/vault`, and `internal/qmd`

---

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org) enforced by CI.

```text
feat(cli): add python language adapter
fix(cli): correct cyclomatic complexity for method receivers
refactor(cli): extract shared metric computation
test(cli): add coverage for blast-radius edge cases
docs(repo): update README with search command examples
chore(repo): update dependencies
```

**Valid types:** `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `build`, `ci`

**Valid scopes:** `cli`, `repo`, `docs`, `test`, `ci`

PR titles must follow the same format (enforced by CI).

---

## Testing

- **Framework:** Standard `testing` package with [gotestsum](https://github.com/gotestyourself/gotestsum) runner
- **Run tests:** `make test` or `make verify`
- **Test location:** Co-located with source (`*_test.go` next to the file under test)
- **File pattern:** `<module>_test.go`

### Conventions

- Table-driven tests with `t.Run` subtests
- `t.TempDir()` for filesystem isolation
- Focused test helpers with `t.Helper()`
- Integration tests use the `//go:build integration` build tag and run via `make test-integration`
- Race detector is always enabled (`-race`)

---

## Adding a New Language Adapter

This is one of the most impactful contributions you can make. `kb` uses a clean adapter interface that makes adding new languages straightforward.

### Step 1: Create the Adapter File

```text
internal/adapter/<language>_adapter.go
```

### Step 2: Implement the LanguageAdapter Interface

The interface is defined in `internal/models/models.go`:

```go
type LanguageAdapter interface {
    Supports(lang SupportedLanguage) bool
    ParseFiles(files []ScannedSourceFile, rootPath string) ([]ParsedFile, error)
}
```

Your `ParseFiles` method must return, for each file:

- **`File`** -- A `GraphFile` node with the file's path, language, and symbol IDs
- **`Symbols`** -- Slice of `SymbolNode` entries (functions, classes, interfaces, types, variables, methods)
- **`ExternalNodes`** -- External module references
- **`Relations`** -- Edges connecting files and symbols (`imports`, `exports`, `calls`, `references`, `declares`, `contains`)
- **`Diagnostics`** -- Any parse warnings or errors

### Step 3: Register Language Extensions

Add your language's file extensions to the `SupportedLanguage` constants in `internal/models/models.go` and extend `SupportedLanguages()`.

### Step 4: Add Tree-sitter Bindings

Add the tree-sitter grammar dependency and wire it in `internal/adapter/treesitter.go`.

### Step 5: Wire the Adapter

Add adapter instantiation and file dispatch in `internal/generate/generate.go`, following the pattern of the existing Go and TypeScript adapters.

### Step 6: Write Tests

Add `<language>_adapter_test.go` in `internal/adapter/` with fixture files in `internal/adapter/testdata/`.

### Step 7: Update Documentation

Add the language to the "Supported Languages" table in `README.md`.

### Reference Implementations

- **`go_adapter.go`** -- Tree-sitter-based Go parser. Relation confidence: `syntactic`.
- **`ts_adapter.go`** -- Tree-sitter-based TypeScript/JavaScript parser. Relation confidence: `syntactic`.

---

## Adding a New File Converter

The converter registry (`internal/convert`) makes adding new file format support straightforward.

### Step 1: Create the Converter File

```text
internal/convert/<format>.go
```

### Step 2: Implement the Converter Interface

The interface is defined in `internal/models/models.go`:

```go
type Converter interface {
    Extensions() []string
    Convert(ctx context.Context, input ConvertInput) (*ConvertResult, error)
}
```

### Step 3: Register the Converter

Add your converter to `NewRegistry()` in `internal/convert/registry.go`.

### Step 4: Write Tests

Add `<format>_test.go` in `internal/convert/` with fixture files in `internal/convert/testdata/`.

### Reference Implementations

- **`pdf.go`** -- PDF text extraction via pdfcpu
- **`docx.go`** -- DOCX XML extraction
- **`html.go`** -- HTML-to-markdown conversion

---

## Pull Request Process

1. Fork the repository and create a feature branch from `main`
2. Make your changes following the code style guidelines above
3. Run `make verify` and ensure it passes with zero warnings
4. Write a clear PR description explaining what changed and why
5. PR titles must follow Conventional Commits format with a scope
6. Keep PRs focused -- one feature or fix per PR

---

## Reporting Issues

- **Bugs:** Use the [bug report template](https://github.com/pedronauck/kodebase-go/issues/new?template=bug-report.yml)
- **Features:** Use the [feature request template](https://github.com/pedronauck/kodebase-go/issues/new?template=feature-request.yml)

Include the command you ran, the output you got, and the output you expected. `kb version` output and Go version help us reproduce faster.

---

## Code of Conduct

Be respectful and constructive. We follow the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).
