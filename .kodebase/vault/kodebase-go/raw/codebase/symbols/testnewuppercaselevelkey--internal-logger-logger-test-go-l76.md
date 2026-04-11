---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 98
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger_test.go"
stage: "raw"
start_line: 76
symbol_kind: "function"
symbol_name: "TestNewUppercaseLevelKey"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNewUppercaseLevelKey"
type: "source"
---

# Codebase Symbol: TestNewUppercaseLevelKey

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 23
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestNewUppercaseLevelKey(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/new--internal-logger-logger-go-l45]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withwriter--internal-logger-logger-go-l28]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]] via `exports` (syntactic)
