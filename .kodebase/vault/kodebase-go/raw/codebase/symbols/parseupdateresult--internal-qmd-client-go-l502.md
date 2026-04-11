---
blast_radius: 2
centrality: 0.1047
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 553
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: true
language: "go"
loc: 52
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 502
symbol_kind: "function"
symbol_name: "parseUpdateResult"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseUpdateResult"
type: "source"
---

# Codebase Symbol: parseUpdateResult

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 2
- External references: 1
- Centrality: 0.1047
- LOC: 52
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func parseUpdateResult(output string, operation IndexOperation) (UpdateResult, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testparseupdateresultparsesaddandupdatesummaries--internal-qmd-client-test-go-l388]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
