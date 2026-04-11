---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 58
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 1
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 58
symbol_kind: "type"
symbol_name: "IndexOperation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: IndexOperation"
type: "source"
---

# Codebase Symbol: IndexOperation

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`type`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 1
- Dead export: true
- Smells: `dead-export`

## Signature
```text
IndexOperation string
```

## Documentation
IndexOperation selects whether Index creates a collection or updates one.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
