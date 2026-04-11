---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 215
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 10
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/models/models.go"
stage: "raw"
start_line: 206
symbol_kind: "struct"
symbol_name: "SymbolMetrics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: SymbolMetrics"
type: "source"
---

# Codebase Symbol: SymbolMetrics

Source file: [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 10
- Dead export: true
- Smells: `dead-export`

## Signature
```text
SymbolMetrics struct {
```

## Documentation
SymbolMetrics stores computed metrics for an individual symbol.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `exports` (syntactic)
