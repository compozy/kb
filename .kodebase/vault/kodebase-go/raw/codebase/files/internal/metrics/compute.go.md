---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 27
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/metrics/compute.go"
stage: "raw"
symbol_count: 20
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/metrics/compute.go"
type: "source"
---

# Codebase File: internal/metrics/compute.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/metrics--internal-metrics-compute-go-l1|metrics (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23|ComputeMetrics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261|computeApproxCentrality (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/detectcirculardependencygroups--internal-metrics-compute-go-l368|detectCircularDependencyGroups (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/findcirculardependencygroups--internal-metrics-compute-go-l375|FindCircularDependencyGroups (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/buildfileimportadjacency--internal-metrics-compute-go-l470|buildFileImportAdjacency (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/uniquesortedstrings--internal-metrics-compute-go-l500|uniqueSortedStrings (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/getnodefileid--internal-metrics-compute-go-l519|getNodeFileID (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/getdirectorypath--internal-metrics-compute-go-l536|getDirectoryPath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isentrypointfile--internal-metrics-compute-go-l540|isEntryPointFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/computeinstability--internal-metrics-compute-go-l556|computeInstability (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tometricratio--internal-metrics-compute-go-l565|toMetricRatio (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isdependencyrelationtype--internal-metrics-compute-go-l569|isDependencyRelationType (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isfunctionlike--internal-metrics-compute-go-l574|isFunctionLike (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/conditionalsmell--internal-metrics-compute-go-l578|conditionalSmell (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collectsortedsmells--internal-metrics-compute-go-l585|collectSortedSmells (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createfileid--internal-metrics-compute-go-l597|createFileID (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/containskey--internal-metrics-compute-go-l601|containsKey (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortedsetvalues--internal-metrics-compute-go-l606|sortedSetValues (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/addtosetmap--internal-metrics-compute-go-l620|addToSetMap (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/addtosetmap--internal-metrics-compute-go-l620]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildfileimportadjacency--internal-metrics-compute-go-l470]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectsortedsmells--internal-metrics-compute-go-l585]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computeinstability--internal-metrics-compute-go-l556]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/conditionalsmell--internal-metrics-compute-go-l578]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/containskey--internal-metrics-compute-go-l601]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-metrics-compute-go-l597]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detectcirculardependencygroups--internal-metrics-compute-go-l368]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findcirculardependencygroups--internal-metrics-compute-go-l375]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getdirectorypath--internal-metrics-compute-go-l536]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getnodefileid--internal-metrics-compute-go-l519]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdependencyrelationtype--internal-metrics-compute-go-l569]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isentrypointfile--internal-metrics-compute-go-l540]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isfunctionlike--internal-metrics-compute-go-l574]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/metrics--internal-metrics-compute-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedsetvalues--internal-metrics-compute-go-l606]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tometricratio--internal-metrics-compute-go-l565]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquesortedstrings--internal-metrics-compute-go-l500]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findcirculardependencygroups--internal-metrics-compute-go-l375]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `math`
- `imports` (syntactic) -> `path`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
