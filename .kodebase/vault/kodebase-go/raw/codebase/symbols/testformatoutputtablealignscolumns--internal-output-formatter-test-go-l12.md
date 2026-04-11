---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 67
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 56
outgoing_relation_count: 0
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter_test.go"
stage: "raw"
start_line: 12
symbol_kind: "function"
symbol_name: "TestFormatOutputTableAlignsColumns"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestFormatOutputTableAlignsColumns"
type: "source"
---

# Codebase Symbol: TestFormatOutputTableAlignsColumns

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter_test.go|internal/output/formatter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 56
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestFormatOutputTableAlignsColumns(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/output/formatter_test.go|internal/output/formatter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter_test.go|internal/output/formatter_test.go]] via `exports` (syntactic)
