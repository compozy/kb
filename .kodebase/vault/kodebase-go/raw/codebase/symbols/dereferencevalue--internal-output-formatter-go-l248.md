---
blast_radius: 7
centrality: 0.2087
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 263
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 248
symbol_kind: "function"
symbol_name: "dereferenceValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: dereferenceValue"
type: "source"
---

# Codebase Symbol: dereferenceValue

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.2087
- LOC: 16
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func dereferenceValue(value any) any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
