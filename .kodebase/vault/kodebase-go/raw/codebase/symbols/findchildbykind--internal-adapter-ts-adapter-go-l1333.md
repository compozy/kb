---
blast_radius: 3
centrality: 0.0543
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 1346
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1333
symbol_kind: "function"
symbol_name: "findChildByKind"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: findChildByKind"
type: "source"
---

# Codebase Symbol: findChildByKind

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0543
- LOC: 14
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func findChildByKind(node *tree_sitter.Node, kind string) *tree_sitter.Node {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
