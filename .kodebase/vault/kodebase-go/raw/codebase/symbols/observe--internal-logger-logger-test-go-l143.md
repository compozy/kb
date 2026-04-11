---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 147
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger_test.go"
stage: "raw"
start_line: 143
symbol_kind: "method"
symbol_name: "Observe"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Observe"
type: "source"
---

# Codebase Symbol: Observe

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 5
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func (o *testObserver) Observe(_ context.Context, record slog.Record) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]] via `exports` (syntactic)
