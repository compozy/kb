---
blast_radius: 13
centrality: 0.0635
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 534
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 519
symbol_kind: "function"
symbol_name: "getNodeFileID"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: getNodeFileID"
type: "source"
---

# Codebase Symbol: getNodeFileID

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.0635
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func getNodeFileID(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-metrics-compute-go-l597]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
