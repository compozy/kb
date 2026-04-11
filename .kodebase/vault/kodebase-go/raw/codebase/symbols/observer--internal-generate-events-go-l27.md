---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 29
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/generate/events.go"
stage: "raw"
start_line: 27
symbol_kind: "interface"
symbol_name: "Observer"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "interface"
title: "Codebase Symbol: Observer"
type: "source"
---

# Codebase Symbol: Observer

Source file: [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]]

## Kind
`interface`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 3
- Dead export: true
- Smells: `dead-export`

## Signature
```text
Observer interface {
```

## Documentation
Observer receives structured generation events.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `exports` (syntactic)
