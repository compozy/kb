---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 47
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 3
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 36
symbol_kind: "function"
symbol_name: "FormatOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: FormatOutput"
type: "source"
---

# Codebase Symbol: FormatOutput

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 12
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func FormatOutput(options FormatOptions) string {
```

## Documentation
FormatOutput renders tabular data using the requested output format.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatjson--internal-output-formatter-go-l87]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsv--internal-output-formatter-go-l141]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `exports` (syntactic)
