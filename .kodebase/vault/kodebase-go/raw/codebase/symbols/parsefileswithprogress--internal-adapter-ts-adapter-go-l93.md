---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 37
domain: "kodebase-go"
end_line: 297
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 205
outgoing_relation_count: 7
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 93
symbol_kind: "method"
symbol_name: "ParseFilesWithProgress"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: ParseFilesWithProgress"
type: "source"
---

# Codebase Symbol: ParseFilesWithProgress

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 37
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 205
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func (adapter TSAdapter) ParseFilesWithProgress(
```

## Documentation
ParseFilesWithProgress parses TS/JS files and reports one progress tick per file.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedexternalnodes--internal-adapter-go-adapter-go-l532]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizeabsolutepath--internal-adapter-ts-adapter-go-l1371]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/pushuniquerelation--internal-adapter-ts-adapter-go-l1357]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationkey--internal-adapter-ts-adapter-go-l1348]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvelocalexports--internal-adapter-ts-adapter-go-l1151]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvereexports--internal-adapter-ts-adapter-go-l1170]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `exports` (syntactic)
