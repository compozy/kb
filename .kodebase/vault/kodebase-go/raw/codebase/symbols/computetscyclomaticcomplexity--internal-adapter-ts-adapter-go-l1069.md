---
blast_radius: 8
centrality: 0.0625
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 1105
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1069
symbol_kind: "function"
symbol_name: "computeTSCyclomaticComplexity"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: computeTSCyclomaticComplexity"
type: "source"
---

# Codebase Symbol: computeTSCyclomaticComplexity

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 8
- External references: 0
- Centrality: 0.0625
- LOC: 37
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func computeTSCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
