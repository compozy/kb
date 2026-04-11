---
blast_radius: 2
centrality: 0.0542
cyclomatic_complexity: 17
domain: "kodebase-go"
end_line: 482
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 92
outgoing_relation_count: 6
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 391
symbol_kind: "function"
symbol_name: "extractTSImports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractTSImports"
type: "source"
---

# Codebase Symbol: extractTSImports

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 17
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0542
- LOC: 92
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func extractTSImports(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolverelativeimportfile--internal-adapter-ts-adapter-go-l1207]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
