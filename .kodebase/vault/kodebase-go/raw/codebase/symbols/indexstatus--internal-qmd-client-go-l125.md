---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 130
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 125
symbol_kind: "struct"
symbol_name: "IndexStatus"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: IndexStatus"
type: "source"
---

# Codebase Symbol: IndexStatus

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: true
- Smells: `dead-export`

## Signature
```text
IndexStatus struct {
```

## Documentation
IndexStatus mirrors QMD's status view for collection health.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
