---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 35
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
symbol_count: 19
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/metrics/compute_test.go"
type: "source"
---

# Codebase File: internal/metrics/compute_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/metrics--internal-metrics-compute-test-go-l1|metrics (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsemptyresultforemptygraph--internal-metrics-compute-test-go-l10|TestComputeMetricsReturnsEmptyResultForEmptyGraph (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricssinglesymbolhaszeroblastradius--internal-metrics-compute-test-go-l32|TestComputeMetricsSingleSymbolHasZeroBlastRadius (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49|TestComputeMetricsCountsBlastRadiusAcrossDependents (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81|TestComputeMetricsCountsFileEfferentCoupling (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103|TestComputeMetricsCountsFileAfferentCoupling (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129|TestComputeMetricsComputesBalancedInstability (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151|TestComputeMetricsFlagsDeadExports (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171|TestComputeMetricsFlagsLongFunctions (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191|TestComputeMetricsDetectsCircularDependencies (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsmergesoverlappingcyclesintosinglecomponent--internal-metrics-compute-test-go-l220|TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsreturnsstablesortedcomponents--internal-metrics-compute-test-go-l236|TestFindCircularDependencyGroupsReturnsStableSortedComponents (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255|TestComputeMetricsReturnsNoCircularDependenciesForAcyclicGraph (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276|TestComputeMetricsAggregatesDirectoryMetrics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/snapshot--internal-metrics-compute-test-go-l305|snapshot (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/filenode--internal-metrics-compute-test-go-l326|fileNode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/symbolnode--internal-metrics-compute-test-go-l336|symbolNode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/relation--internal-metrics-compute-test-go-l355|relation (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/containsstring--internal-metrics-compute-test-go-l364|containsString (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/containsstring--internal-metrics-compute-test-go-l364]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filenode--internal-metrics-compute-test-go-l326]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/metrics--internal-metrics-compute-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relation--internal-metrics-compute-test-go-l355]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/snapshot--internal-metrics-compute-test-go-l305]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolnode--internal-metrics-compute-test-go-l336]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsemptyresultforemptygraph--internal-metrics-compute-test-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricssinglesymbolhaszeroblastradius--internal-metrics-compute-test-go-l32]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsmergesoverlappingcyclesintosinglecomponent--internal-metrics-compute-test-go-l220]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsreturnsstablesortedcomponents--internal-metrics-compute-test-go-l236]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsemptyresultforemptygraph--internal-metrics-compute-test-go-l10]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcomputemetricssinglesymbolhaszeroblastradius--internal-metrics-compute-test-go-l32]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsmergesoverlappingcyclesintosinglecomponent--internal-metrics-compute-test-go-l220]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsreturnsstablesortedcomponents--internal-metrics-compute-test-go-l236]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `testing`

## Backlinks
None
