---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 326
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/models/models.go"
stage: "raw"
start_line: 323
symbol_kind: "struct"
symbol_name: "BaseGroupBy"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: BaseGroupBy"
type: "source"
---

# Codebase Symbol: BaseGroupBy

Source file: [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 4
- Dead export: true
- Smells: `dead-export`

## Signature
```text
BaseGroupBy struct {
```

## Documentation
BaseGroupBy configures the grouping rule for a Base view.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `exports` (syntactic)
