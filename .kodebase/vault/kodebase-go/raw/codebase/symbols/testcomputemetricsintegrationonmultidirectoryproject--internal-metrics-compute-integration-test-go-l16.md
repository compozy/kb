---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 91
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 76
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_integration_test.go"
stage: "raw"
start_line: 16
symbol_kind: "function"
symbol_name: "TestComputeMetricsIntegrationOnMultiDirectoryProject"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestComputeMetricsIntegrationOnMultiDirectoryProject"
type: "source"
---

# Codebase Symbol: TestComputeMetricsIntegrationOnMultiDirectoryProject

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 76
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestComputeMetricsIntegrationOnMultiDirectoryProject(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsymbolbynameandfile--internal-metrics-compute-integration-test-go-l122]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetsworkspace--internal-metrics-compute-integration-test-go-l93]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]] via `exports` (syntactic)
