---
blast_radius: 9
centrality: 0.0741
cyclomatic_complexity: 12
domain: "kodebase-go"
end_line: 914
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: true
language: "go"
loc: 35
outgoing_relation_count: 2
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 880
symbol_kind: "function"
symbol_name: "resolveTSSymbolName"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveTSSymbolName"
type: "source"
---

# Codebase Symbol: resolveTSSymbolName

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 12
- Long function: true
- Blast radius: 9
- External references: 0
- Centrality: 0.0741
- LOC: 35
- Dead export: false
- Smells: `long-function`

## Signature
```text
func resolveTSSymbolName(node *tree_sitter.Node, source []byte, symbolKind string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
