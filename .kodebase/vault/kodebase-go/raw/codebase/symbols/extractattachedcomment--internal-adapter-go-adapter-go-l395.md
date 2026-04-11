---
blast_radius: 13
centrality: 0.154
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 427
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 33
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 395
symbol_kind: "function"
symbol_name: "extractAttachedComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractAttachedComment"
type: "source"
---

# Codebase Symbol: extractAttachedComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: false
- Blast radius: 13
- External references: 1
- Centrality: 0.154
- LOC: 33
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func extractAttachedComment(node *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-adapter-go-adapter-go-l607]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractgodoccomment--internal-adapter-go-adapter-go-l381]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsdoccomment--internal-adapter-ts-adapter-go-l1004]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
