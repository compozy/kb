---
blast_radius: 7
centrality: 0.1241
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 362
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 8
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 355
symbol_kind: "function"
symbol_name: "relation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: relation"
type: "source"
---

# Codebase Symbol: relation

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.1241
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func relation(fromID, toID string, relationType models.RelationType) models.RelationEdge {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
