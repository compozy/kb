---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 41
exported: true
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 37
symbol_kind: "function"
symbol_name: "WithObserver"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithObserver"
type: "source"
---

# Codebase Symbol: WithObserver

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0651
- LOC: 5
- Dead export: false
- Smells: None

## Signature
```text
func WithObserver(observer LogObserver) Option {
```

## Documentation
WithObserver receives records after level filtering and before write.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `exports` (syntactic)
