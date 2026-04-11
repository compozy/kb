---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 151
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 145
symbol_kind: "struct"
symbol_name: "QMDClient"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: QMDClient"
type: "source"
---

# Codebase Symbol: QMDClient

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 7
- Dead export: true
- Smells: `dead-export`

## Signature
```text
QMDClient struct {
```

## Documentation
QMDClient executes QMD shell commands with context-aware process management.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
