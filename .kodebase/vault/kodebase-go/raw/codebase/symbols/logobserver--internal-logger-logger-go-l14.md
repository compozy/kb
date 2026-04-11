---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 16
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
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 14
symbol_kind: "interface"
symbol_name: "LogObserver"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "interface"
title: "Codebase Symbol: LogObserver"
type: "source"
---

# Codebase Symbol: LogObserver

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

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
LogObserver interface {
```

## Documentation
LogObserver receives a clone of each log record that passes the handler's
level filter.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `exports` (syntactic)
