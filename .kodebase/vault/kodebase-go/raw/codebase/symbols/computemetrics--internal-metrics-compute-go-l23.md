---
blast_radius: 12
centrality: 0.2097
cyclomatic_complexity: 52
domain: "kodebase-go"
end_line: 259
exported: true
external_reference_count: 12
has_smells: true
incoming_relation_count: 14
is_dead_export: false
is_long_function: true
language: "go"
loc: 237
outgoing_relation_count: 14
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 23
symbol_kind: "function"
symbol_name: "ComputeMetrics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ComputeMetrics"
type: "source"
---

# Codebase Symbol: ComputeMetrics

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 52
- Long function: true
- Blast radius: 12
- External references: 12
- Centrality: 0.2097
- LOC: 237
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func ComputeMetrics(graph models.GraphSnapshot) models.MetricsResult {
```

## Documentation
ComputeMetrics derives symbol, file, directory, and circular dependency metrics for a graph snapshot.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/addtosetmap--internal-metrics-compute-go-l620]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectsortedsmells--internal-metrics-compute-go-l585]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computeinstability--internal-metrics-compute-go-l556]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/conditionalsmell--internal-metrics-compute-go-l578]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/containskey--internal-metrics-compute-go-l601]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-metrics-compute-go-l597]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detectcirculardependencygroups--internal-metrics-compute-go-l368]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getdirectorypath--internal-metrics-compute-go-l536]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getnodefileid--internal-metrics-compute-go-l519]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdependencyrelationtype--internal-metrics-compute-go-l569]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isentrypointfile--internal-metrics-compute-go-l540]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isfunctionlike--internal-metrics-compute-go-l574]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedsetvalues--internal-metrics-compute-go-l606]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsintegrationonmultidirectoryproject--internal-metrics-compute-integration-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsemptyresultforemptygraph--internal-metrics-compute-test-go-l10]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricssinglesymbolhaszeroblastradius--internal-metrics-compute-test-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `exports` (syntactic)
