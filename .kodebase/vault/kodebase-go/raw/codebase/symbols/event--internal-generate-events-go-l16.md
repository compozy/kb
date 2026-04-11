---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 24
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/generate/events.go"
stage: "raw"
start_line: 16
symbol_kind: "struct"
symbol_name: "Event"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: Event"
type: "source"
---

# Codebase Symbol: Event

Source file: [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 9
- Dead export: true
- Smells: `dead-export`

## Signature
```text
Event struct {
```

## Documentation
Event reports structured progress from the generate pipeline.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `exports` (syntactic)
