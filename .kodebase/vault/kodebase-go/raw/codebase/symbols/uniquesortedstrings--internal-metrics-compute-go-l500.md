---
blast_radius: 17
centrality: 0.1902
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 517
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 18
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 500
symbol_kind: "function"
symbol_name: "uniqueSortedStrings"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: uniqueSortedStrings"
type: "source"
---

# Codebase Symbol: uniqueSortedStrings

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 17
- External references: 0
- Centrality: 0.1902
- LOC: 18
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func uniqueSortedStrings(values []string) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/findcirculardependencygroups--internal-metrics-compute-go-l375]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
