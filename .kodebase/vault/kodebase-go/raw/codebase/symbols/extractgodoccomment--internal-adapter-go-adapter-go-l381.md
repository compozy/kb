---
blast_radius: 3
centrality: 0.0589
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 393
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 381
symbol_kind: "function"
symbol_name: "extractGoDocComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractGoDocComment"
type: "source"
---

# Codebase Symbol: extractGoDocComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0589
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func extractGoDocComment(node *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
