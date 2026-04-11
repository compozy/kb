---
blast_radius: 13
centrality: 0.0635
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 618
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 606
symbol_kind: "function"
symbol_name: "sortedSetValues"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sortedSetValues"
type: "source"
---

# Codebase Symbol: sortedSetValues

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.0635
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func sortedSetValues[T ~string](items map[T]struct{}) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
