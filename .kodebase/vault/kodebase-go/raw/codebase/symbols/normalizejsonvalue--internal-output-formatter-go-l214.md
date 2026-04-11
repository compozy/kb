---
blast_radius: 7
centrality: 0.1427
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 246
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 33
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 214
symbol_kind: "function"
symbol_name: "normalizeJSONValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeJSONValue"
type: "source"
---

# Codebase Symbol: normalizeJSONValue

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.1427
- LOC: 33
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeJSONValue(value any) any {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/dereferencevalue--internal-output-formatter-go-l248]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formatjson--internal-output-formatter-go-l87]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
