---
blast_radius: 13
centrality: 0.0635
cyclomatic_complexity: 24
domain: "kodebase-go"
end_line: 366
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 106
outgoing_relation_count: 2
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 261
symbol_kind: "function"
symbol_name: "computeApproxCentrality"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: computeApproxCentrality"
type: "source"
---

# Codebase Symbol: computeApproxCentrality

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 24
- Long function: true
- Blast radius: 13
- External references: 0
- Centrality: 0.0635
- LOC: 106
- Dead export: false
- Smells: `long-function`

## Signature
```text
func computeApproxCentrality(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdependencyrelationtype--internal-metrics-compute-go-l569]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tometricratio--internal-metrics-compute-go-l565]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
