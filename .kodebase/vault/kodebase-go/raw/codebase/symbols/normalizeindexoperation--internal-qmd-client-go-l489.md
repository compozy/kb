---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 500
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 489
symbol_kind: "function"
symbol_name: "normalizeIndexOperation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeIndexOperation"
type: "source"
---

# Codebase Symbol: normalizeIndexOperation

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0615
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func normalizeIndexOperation(operation IndexOperation) (IndexOperation, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
