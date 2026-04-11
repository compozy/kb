---
blast_radius: 2
centrality: 0.0573
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 463
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 35
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 429
symbol_kind: "function"
symbol_name: "extractCallTargetNames"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractCallTargetNames"
type: "source"
---

# Codebase Symbol: extractCallTargetNames

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0573
- LOC: 35
- Dead export: false
- Smells: None

## Signature
```text
func extractCallTargetNames(node *tree_sitter.Node, source []byte) []string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
