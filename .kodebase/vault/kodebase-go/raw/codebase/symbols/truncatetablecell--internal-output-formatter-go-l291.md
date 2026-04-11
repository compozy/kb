---
blast_radius: 4
centrality: 0.0861
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 307
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 291
symbol_kind: "function"
symbol_name: "truncateTableCell"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: truncateTableCell"
type: "source"
---

# Codebase Symbol: truncateTableCell

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.0861
- LOC: 17
- Dead export: false
- Smells: None

## Signature
```text
func truncateTableCell(value string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runecount--internal-output-formatter-go-l327]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
