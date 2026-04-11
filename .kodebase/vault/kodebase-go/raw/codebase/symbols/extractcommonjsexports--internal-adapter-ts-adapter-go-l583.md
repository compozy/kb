---
blast_radius: 2
centrality: 0.0542
cyclomatic_complexity: 22
domain: "kodebase-go"
end_line: 642
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 60
outgoing_relation_count: 3
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 583
symbol_kind: "function"
symbol_name: "extractCommonJSExports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractCommonJSExports"
type: "source"
---

# Codebase Symbol: extractCommonJSExports

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 22
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0542
- LOC: 60
- Dead export: false
- Smells: `long-function`

## Signature
```text
func extractCommonJSExports(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
