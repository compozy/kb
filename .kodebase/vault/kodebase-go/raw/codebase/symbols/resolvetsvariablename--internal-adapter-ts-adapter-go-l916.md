---
blast_radius: 10
centrality: 0.1101
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 927
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 916
symbol_kind: "function"
symbol_name: "resolveTSVariableName"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveTSVariableName"
type: "source"
---

# Codebase Symbol: resolveTSVariableName

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 10
- External references: 0
- Centrality: 0.1101
- LOC: 12
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func resolveTSVariableName(node *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
