---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 184
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/models/models.go"
stage: "raw"
start_line: 177
symbol_kind: "struct"
symbol_name: "GraphSnapshot"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: GraphSnapshot"
type: "source"
---

# Codebase Symbol: GraphSnapshot

Source file: [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 8
- Dead export: true
- Smells: `dead-export`

## Signature
```text
GraphSnapshot struct {
```

## Documentation
GraphSnapshot is the merged graph output across all parsed files.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `exports` (syntactic)
