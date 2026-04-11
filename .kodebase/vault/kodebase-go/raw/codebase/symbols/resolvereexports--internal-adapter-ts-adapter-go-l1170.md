---
blast_radius: 1
centrality: 0.0569
cyclomatic_complexity: 12
domain: "kodebase-go"
end_line: 1205
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 36
outgoing_relation_count: 0
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1170
symbol_kind: "function"
symbol_name: "resolveReExports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveReExports"
type: "source"
---

# Codebase Symbol: resolveReExports

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 12
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0569
- LOC: 36
- Dead export: false
- Smells: `long-function`

## Signature
```text
func resolveReExports(parsedEntries []parsedTSFile, exportedSymbolsByFilePath map[string]map[string]string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
