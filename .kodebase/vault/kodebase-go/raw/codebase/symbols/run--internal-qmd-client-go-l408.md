---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 440
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 33
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 408
symbol_kind: "method"
symbol_name: "run"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: run"
type: "source"
---

# Codebase Symbol: run

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 33
- Dead export: false
- Smells: None

## Signature
```text
func (client *QMDClient) run(ctx context.Context, spec commandSpec) (string, string, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleandiagnostics--internal-qmd-client-go-l740]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
