---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 361
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 336
symbol_kind: "method"
symbol_name: "indexCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: indexCommand"
type: "source"
---

# Codebase Symbol: indexCommand

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func (client *QMDClient) indexCommand(operation IndexOperation, options IndexOptions) (commandSpec, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
