---
blast_radius: 15
centrality: 0.1317
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 567
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 565
symbol_kind: "function"
symbol_name: "toMetricRatio"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toMetricRatio"
type: "source"
---

# Codebase Symbol: toMetricRatio

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 15
- External references: 0
- Centrality: 0.1317
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func toMetricRatio(value float64) float64 {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/computeinstability--internal-metrics-compute-go-l556]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
