---
blast_radius: 8
centrality: 0.0625
cyclomatic_complexity: 13
domain: "kodebase-go"
end_line: 966
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 38
outgoing_relation_count: 3
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 929
symbol_kind: "function"
symbol_name: "formatTSSignature"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: formatTSSignature"
type: "source"
---

# Codebase Symbol: formatTSSignature

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 13
- Long function: true
- Blast radius: 8
- External references: 0
- Centrality: 0.0625
- LOC: 38
- Dead export: false
- Smells: `long-function`

## Signature
```text
func formatTSSignature(node *tree_sitter.Node, source []byte, symbolKind, name string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsreturntype--internal-adapter-ts-adapter-go-l968]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsvariabletypesuffix--internal-adapter-ts-adapter-go-l986]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
