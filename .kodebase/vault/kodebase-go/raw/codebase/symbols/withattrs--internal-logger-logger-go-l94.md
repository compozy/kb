---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 99
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 94
symbol_kind: "method"
symbol_name: "WithAttrs"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: WithAttrs"
type: "source"
---

# Codebase Symbol: WithAttrs

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func (h *observedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `exports` (syntactic)
