# Kodebase Go

**Topic scope:** Generated codebase knowledge topic for `kodebase-go`. This topic stages raw code snapshots in `raw/codebase/` and compiles a starter wiki in `wiki/`.

**Domain:** `kodebase-go`

This file is the schema document for the topic. The `kodebase` CLI manages `raw/codebase/`, `wiki/index/`, and wiki concept pages with `generator: kodebase` frontmatter. Everything else may be extended manually without being overwritten.

## Audit log

See [log.md](log.md) for the append-only record of ingest and compile operations.

## Current wiki articles

- [[kodebase-go/wiki/concepts/Circular Dependencies|Circular Dependencies]]
- [[kodebase-go/wiki/concepts/Code Smells|Code Smells]]
- [[kodebase-go/wiki/concepts/Codebase Overview|Codebase Overview]]
- [[kodebase-go/wiki/concepts/Complexity Hotspots|Complexity Hotspots]]
- [[kodebase-go/wiki/concepts/Dead Code Report|Dead Code Report]]
- [[kodebase-go/wiki/concepts/Dependency Hotspots|Dependency Hotspots]]
- [[kodebase-go/wiki/concepts/Directory Map|Directory Map]]
- [[kodebase-go/wiki/concepts/High-Impact Symbols|High-Impact Symbols]]
- [[kodebase-go/wiki/concepts/Module Health|Module Health]]
- [[kodebase-go/wiki/concepts/Symbol Taxonomy|Symbol Taxonomy]]

## Codebase corpus

- Parsed files: 80
- Parsed symbols: 901
- Raw codebase notes: 1001
- `raw/codebase/files/` - file-level snapshots generated from source files
- `raw/codebase/symbols/` - symbol-level snapshots generated from extracted declarations
- `raw/codebase/indexes/` - generated directory and language inventories
- `bases/` - generated Obsidian Bases views over the raw codebase notes

## Managed starter wiki

- [[kodebase-go/wiki/index/Dashboard|Dashboard]]
- [[kodebase-go/wiki/index/Concept Index|Concept Index]]
- [[kodebase-go/wiki/index/Source Index|Source Index]]

## Research gaps

- Expand the starter wiki into architecture-level articles for the main subsystems.
- Promote repeated query outputs in `outputs/queries/` into first-class wiki articles when they stabilize.
- Add manually curated raw material in `raw/articles/`, `raw/github/`, or `raw/bookmarks/` when source code alone is not enough.
