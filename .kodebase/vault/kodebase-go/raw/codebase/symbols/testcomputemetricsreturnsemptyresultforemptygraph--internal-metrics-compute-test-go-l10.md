---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 30
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 10
symbol_kind: "function"
symbol_name: "TestComputeMetricsReturnsEmptyResultForEmptyGraph"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestComputeMetricsReturnsEmptyResultForEmptyGraph"
type: "source"
---

# Codebase Symbol: TestComputeMetricsReturnsEmptyResultForEmptyGraph

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 21
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestComputeMetricsReturnsEmptyResultForEmptyGraph(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `exports` (syntactic)
