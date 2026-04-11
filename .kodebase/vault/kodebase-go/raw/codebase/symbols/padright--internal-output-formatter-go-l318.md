---
blast_radius: 3
centrality: 0.1096
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 325
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 318
symbol_kind: "function"
symbol_name: "padRight"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: padRight"
type: "source"
---

# Codebase Symbol: padRight

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1096
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func padRight(value string, width int) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runecount--internal-output-formatter-go-l327]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formatstringrow--internal-output-formatter-go-l309]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
