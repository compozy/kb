---
blast_radius: 3
centrality: 0.0543
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 1283
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 4
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1265
symbol_kind: "function"
symbol_name: "declarationExportNames"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: declarationExportNames"
type: "source"
---

# Codebase Symbol: declarationExportNames

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0543
- LOC: 19
- Dead export: false
- Smells: None

## Signature
```text
func declarationExportNames(node *tree_sitter.Node, source []byte) []string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gettssymbolkind--internal-adapter-ts-adapter-go-l859]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
