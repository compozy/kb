---
blast_radius: 7
centrality: 0.2356
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 329
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 327
symbol_kind: "function"
symbol_name: "runeCount"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: runeCount"
type: "source"
---

# Codebase Symbol: runeCount

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.2356
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func runeCount(value string) int {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/padright--internal-output-formatter-go-l318]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/truncatetablecell--internal-output-formatter-go-l291]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
