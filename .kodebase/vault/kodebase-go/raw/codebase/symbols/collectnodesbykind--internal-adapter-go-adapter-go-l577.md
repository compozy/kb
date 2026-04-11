---
blast_radius: 9
centrality: 0.0978
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 586
exported: false
external_reference_count: 3
has_smells: false
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 577
symbol_kind: "function"
symbol_name: "collectNodesByKind"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: collectNodesByKind"
type: "source"
---

# Codebase Symbol: collectNodesByKind

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 9
- External references: 3
- Centrality: 0.0978
- LOC: 10
- Dead export: false
- Smells: None

## Signature
```text
func collectNodesByKind(node *tree_sitter.Node, targetKind string) []tree_sitter.Node {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
