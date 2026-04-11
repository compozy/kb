---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 303
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 28
outgoing_relation_count: 4
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 276
symbol_kind: "function"
symbol_name: "TestComputeMetricsAggregatesDirectoryMetrics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestComputeMetricsAggregatesDirectoryMetrics"
type: "source"
---

# Codebase Symbol: TestComputeMetricsAggregatesDirectoryMetrics

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 28
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestComputeMetricsAggregatesDirectoryMetrics(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filenode--internal-metrics-compute-test-go-l326]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relation--internal-metrics-compute-test-go-l355]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/snapshot--internal-metrics-compute-test-go-l305]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `exports` (syntactic)
