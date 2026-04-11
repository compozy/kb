---
blast_radius: 8
centrality: 0.0625
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 1016
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1004
symbol_kind: "function"
symbol_name: "extractTSDocComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractTSDocComment"
type: "source"
---

# Codebase Symbol: extractTSDocComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 8
- External references: 0
- Centrality: 0.0625
- LOC: 13
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func extractTSDocComment(node *tree_sitter.Node, anchorNode *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
