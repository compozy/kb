---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 32
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
source_path: "internal/generate/events.go"
stage: "raw"
start_line: 32
symbol_kind: "type"
symbol_name: "ObserverFunc"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: ObserverFunc"
type: "source"
---

# Codebase Symbol: ObserverFunc

Source file: [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]]

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
ObserverFunc func(context.Context, Event)
```

## Documentation
ObserverFunc adapts a function into an Observer.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] via `exports` (syntactic)
