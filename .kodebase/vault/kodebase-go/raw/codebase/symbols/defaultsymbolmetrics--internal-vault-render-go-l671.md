---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 676
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 671
symbol_kind: "function"
symbol_name: "defaultSymbolMetrics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: defaultSymbolMetrics"
type: "source"
---

# Codebase Symbol: defaultSymbolMetrics

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0524
- LOC: 6
- Dead export: false
- Smells: None

## Signature
```text
func defaultSymbolMetrics(symbol models.SymbolNode, metric models.SymbolMetrics) models.SymbolMetrics {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/maxint--internal-vault-render-go-l811]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
