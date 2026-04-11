---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 132
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize.go"
stage: "raw"
start_line: 123
symbol_kind: "function"
symbol_name: "compareDiagnostics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: compareDiagnostics"
type: "source"
---

# Codebase Symbol: compareDiagnostics

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 10
- Dead export: false
- Smells: None

## Signature
```text
func compareDiagnostics(left models.StructuredDiagnostic, right models.StructuredDiagnostic) int {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/comparestrings--internal-graph-normalize-go-l134]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]] via `contains` (syntactic)
