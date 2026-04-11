---
blast_radius: 2
centrality: 0.0677
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 549
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 18
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 532
symbol_kind: "function"
symbol_name: "sortedExternalNodes"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sortedExternalNodes"
type: "source"
---

# Codebase Symbol: sortedExternalNodes

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.0677
- LOC: 18
- Dead export: false
- Smells: None

## Signature
```text
func sortedExternalNodes(externalNodes map[string]models.ExternalNode) []models.ExternalNode {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
