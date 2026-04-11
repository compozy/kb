---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 155
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 141
symbol_kind: "function"
symbol_name: "formatTSV"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: formatTSV"
type: "source"
---

# Codebase Symbol: formatTSV

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 15
- Dead export: false
- Smells: None

## Signature
```text
func formatTSV(columns []string, data []map[string]any) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formatoutput--internal-output-formatter-go-l36]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
