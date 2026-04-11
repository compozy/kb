---
blast_radius: 3
centrality: 0.0589
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 501
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 465
symbol_kind: "function"
symbol_name: "computeCyclomaticComplexity"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: computeCyclomaticComplexity"
type: "source"
---

# Codebase Symbol: computeCyclomaticComplexity

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0589
- LOC: 37
- Dead export: false
- Smells: None

## Signature
```text
func computeCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
