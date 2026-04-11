---
blast_radius: 5
centrality: 0.1612
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 121
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 108
symbol_kind: "function"
symbol_name: "parseLevel"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseLevel"
type: "source"
---

# Codebase Symbol: parseLevel

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1612
- LOC: 14
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseLevel(level string) (slog.Level, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/new--internal-logger-logger-go-l45]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
