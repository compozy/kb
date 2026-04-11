---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 139
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 53
outgoing_relation_count: 1
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 87
symbol_kind: "function"
symbol_name: "formatJSON"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: formatJSON"
type: "source"
---

# Codebase Symbol: formatJSON

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 53
- Dead export: false
- Smells: `long-function`

## Signature
```text
func formatJSON(columns []string, data []map[string]any) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formatoutput--internal-output-formatter-go-l36]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
