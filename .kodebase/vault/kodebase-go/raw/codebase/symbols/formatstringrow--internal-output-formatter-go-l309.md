---
blast_radius: 2
centrality: 0.0692
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 316
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 309
symbol_kind: "function"
symbol_name: "formatStringRow"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: formatStringRow"
type: "source"
---

# Codebase Symbol: formatStringRow

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0692
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func formatStringRow(values []string, widths []int) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/padright--internal-output-formatter-go-l318]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
