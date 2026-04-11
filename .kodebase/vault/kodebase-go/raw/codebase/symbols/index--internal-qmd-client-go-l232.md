---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 15
domain: "kodebase-go"
end_line: 305
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 74
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 232
symbol_kind: "method"
symbol_name: "Index"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Index"
type: "source"
---

# Codebase Symbol: Index

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 15
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 74
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func (client *QMDClient) Index(ctx context.Context, options IndexOptions) (IndexResult, error) {
```

## Documentation
Index creates or updates a QMD collection and returns the structured summary.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizeindexoperation--internal-qmd-client-go-l489]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseupdateresult--internal-qmd-client-go-l502]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
