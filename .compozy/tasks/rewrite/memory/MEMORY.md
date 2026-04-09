# Workflow Memory

Keep only durable, cross-task context here. Do not duplicate facts that are obvious from the repository, PRD documents, or git history.

## Current State

## Shared Decisions
- Resolve the Go grammar import path `github.com/tree-sitter/tree-sitter-go/bindings/go` through the root module `github.com/tree-sitter/tree-sitter-go v0.23.4` with a `replace` directive. The standalone `bindings/go` submodule has a broken test dependency on `github.com/tree-sitter/tree-sitter-go`, which causes `go mod tidy` to fail.
- Treat the shorter `internal/...` package layout as active on this branch, including `internal/models`, `internal/adapter`, and `internal/vault`. The `_techspec.md` references to `internal/kodebase/...` are stale historical paths and should not drive new code placement.

## Shared Learnings
- The rewrite task documents for the renderer and final integration still say 12 Obsidian base files, but the authoritative TypeScript reference implementation and upstream tests currently define 11 named base files. Treat this as a task-doc mismatch unless a later product decision introduces a twelfth base.
- Task 10 follows the upstream renderer split: `RenderDocuments` returns markdown documents whose `Body` already includes YAML frontmatter, while `RenderBaseFiles` returns separate `.base` definitions. Downstream writer work should write `document.Body` directly for markdown files and `RenderBaseDefinition(base.Definition)` for base files.
- The upstream `.base` files render as YAML, not JSON. Downstream writer/reader/tests should validate them as YAML documents with `views` arrays.
- Raw vault snapshots do not persist full circular-dependency cycle lists. Read-side features that need the ordered cycles must reconstruct them from the file-level `imports` relations in `snapshot.Files`; the stored `has_circular_dependency` flag only marks participating files.
- The current QMD CLI does not expose the task doc's `qmd search --mode ...` interface. The live shell surface is `qmd query` for hybrid search, `qmd search` for lexical search, `qmd vsearch` for vector search, `qmd collection add <path> --name <collection>`, and `qmd update <collection>`. Search commands honor `--json`; `status`/`update`/`collection list` currently remain human-readable even when `--json` is passed.
- QMD JSON search output stays on stdout while progress/warning text is emitted on stderr. Shell integrations should parse stdout as JSON and reserve stderr for diagnostics, not treat stderr output alone as command failure.
- QMD integration tests need `XDG_CACHE_HOME`, `XDG_CONFIG_HOME`, and `HOME` isolated together to avoid leaking the user's existing collection configuration into temporary test indexes.
- `qmd collection add` is not idempotent for an existing collection name. CLI flows that support re-indexing must detect existing collections first and switch to `qmd update <collection>` instead of retrying `collection add`.
- Empty `qmd status` output is still structurally useful for automation: it reports `Total: 0` / `Vectors: 0` plus a `No collections.` line. Parsers should treat that as a valid zero-collection status rather than an error.

## Open Risks
- Task tracking is inconsistent with the branch state: `task_02.md` is marked completed, but the full `internal/models` package is absent. Task 03 only restores the scanner-facing subset, so future model-dependent tasks may still need the remaining domain types.
- QMD embedding can fail on hosts where sqlite extension loading is unavailable (`sqlite-vec is not available`). Search/index flows should surface stderr cleanly, and integration coverage should avoid assuming vector embedding is always supported.

## Handoffs
