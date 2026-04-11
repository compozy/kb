---
blast_radius: 13
centrality: 0.0635
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 628
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 620
symbol_kind: "function"
symbol_name: "addToSetMap"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: addToSetMap"
type: "source"
---

# Codebase Symbol: addToSetMap

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.0635
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func addToSetMap[K comparable, V comparable](items map[K]map[V]struct{}, key K, value V) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
