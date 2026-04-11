---
blast_radius: 9
centrality: 0.0685
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 984
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 2
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 968
symbol_kind: "function"
symbol_name: "formatTSReturnType"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: formatTSReturnType"
type: "source"
---

# Codebase Symbol: formatTSReturnType

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 9
- External references: 0
- Centrality: 0.0685
- LOC: 17
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func formatTSReturnType(node *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formattssignature--internal-adapter-ts-adapter-go-l929]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
