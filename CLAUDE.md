# Kodebase Go Port — AGENTS.md

## HIGH PRIORITY

- **YOU MUST** use `gpt-5.4` model with `xhigh` reasoning effort
- **YOU MUST** run `make verify` (fmt + lint + test + build) before completing ANY task
- **IF YOU DON'T CHECK SKILLS** your task will be invalidated
- **NEVER** use workarounds — use `no-workarounds` + `systematic-debugging` skills for bugs
- **NEVER** add dependencies by hand in `go.mod` — always use `go get`
- **NEVER** run `git restore`, `git checkout`, `git reset`, `git clean`, `git rm` without explicit permission

## PROJECT OVERVIEW

Kodebase CLI — A Go port of [kodebase](https://github.com/compozy/kodebase).
Turns source code repositories into Karpathy-style Obsidian knowledge vaults with rich code metrics.

**Reference TypeScript source:** `/root/kodebase` (~6,700 LOC, 24 files)

## IMPLEMENTATION PLAN

Read the full plan at: `docs/plans/implementation-plan.md`

## PACKAGE LAYOUT

| Path | Responsibility |
|------|---------------|
| `cmd/kodebase` | CLI entry point |
| `internal/cli` | Cobra command definitions |
| `internal/kodebase/models` | Domain types and interfaces |
| `internal/kodebase/scanner` | Workspace file discovery |
| `internal/kodebase/adapter` | Tree-sitter language adapters (Go, TS/JS) |
| `internal/kodebase/graph` | Graph normalization |
| `internal/kodebase/metrics` | Metrics engine |
| `internal/kodebase/vault` | Vault rendering, writing, reading |
| `internal/kodebase/qmd` | QMD shell client |
| `internal/kodebase/output` | Output formatting |
| `internal/config` | Config loading |
| `internal/logger` | Structured logging |
| `internal/version` | Build metadata |

## BUILD COMMANDS

```bash
make verify    # fmt -> lint -> test -> build (BLOCKING GATE)
make fmt       # gofmt
make lint      # golangci-lint (zero issues)
make test      # go test ./... -race
make build     # go build ./...
make deps      # go mod tidy
```

## CODING STYLE

- Go 1.24, cobra for CLI, slog for logging
- `fmt.Errorf("context: %w", err)` for error wrapping
- `errors.Is()`/`errors.As()` for error matching
- `context.Context` as first arg on all functions crossing runtime boundaries
- Table-driven tests with `t.Run`, `t.Parallel()`, `t.TempDir()`
- No `panic()`/`log.Fatal()` in production paths
- Functional options for complex constructors
- Compile-time interface checks: `var _ Interface = (*Type)(nil)`
- Always reference the TypeScript source in `/root/kodebase` when implementing

## ARCHITECTURE

Single-binary Go CLI. Pipeline: Scan → Parse (tree-sitter adapters) → Normalize Graph → Compute Metrics → Render Documents → Write Vault → Read/Inspect.

- `LanguageAdapter` interface: `Supports(lang) bool`, `ParseFiles(files, rootPath) ([]ParsedFile, error)`
- Adapters: GoAdapter (tree-sitter-go), TSAdapter (tree-sitter-typescript + tree-sitter-javascript)
- QMD integration via shell calls (`os/exec`), graceful fallback when not installed
