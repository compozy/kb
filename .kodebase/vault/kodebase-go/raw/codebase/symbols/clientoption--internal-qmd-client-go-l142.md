---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 142
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
start_line: 142
symbol_kind: "type"
symbol_name: "ClientOption"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: ClientOption"
type: "source"
---

# Codebase Symbol: ClientOption

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
ClientOption func(*QMDClient)
```

## Documentation
ClientOption mutates a QMDClient configuration.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
