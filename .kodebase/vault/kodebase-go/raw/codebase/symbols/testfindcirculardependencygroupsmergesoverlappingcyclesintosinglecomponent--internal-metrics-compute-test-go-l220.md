---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 234
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 220
symbol_kind: "function"
symbol_name: "TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent"
type: "source"
---

# Codebase Symbol: TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 15
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findcirculardependencygroups--internal-metrics-compute-go-l375]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `exports` (syntactic)
