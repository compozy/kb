---
blast_radius: 3
centrality: 0.0543
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 1288
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1285
symbol_kind: "function"
symbol_name: "isDefaultExport"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isDefaultExport"
type: "source"
---

# Codebase Symbol: isDefaultExport

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0543
- LOC: 4
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func isDefaultExport(node *tree_sitter.Node, source []byte) bool {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
