---
blast_radius: 5
centrality: 0.0861
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 212
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 36
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 177
symbol_kind: "function"
symbol_name: "normalizeCellValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeCellValue"
type: "source"
---

# Codebase Symbol: normalizeCellValue

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0861
- LOC: 36
- Dead export: false
- Smells: None

## Signature
```text
func normalizeCellValue(value any) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/dereferencevalue--internal-output-formatter-go-l248]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
