---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 169
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 5
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 151
symbol_kind: "function"
symbol_name: "TestComputeMetricsFlagsDeadExports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestComputeMetricsFlagsDeadExports"
type: "source"
---

# Codebase Symbol: TestComputeMetricsFlagsDeadExports

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 19
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestComputeMetricsFlagsDeadExports(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/containsstring--internal-metrics-compute-test-go-l364]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filenode--internal-metrics-compute-test-go-l326]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/snapshot--internal-metrics-compute-test-go-l305]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolnode--internal-metrics-compute-test-go-l336]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `exports` (syntactic)
