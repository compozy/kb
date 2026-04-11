---
blast_radius: 16
centrality: 0.164
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 468
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: true
language: "go"
loc: 94
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 375
symbol_kind: "function"
symbol_name: "FindCircularDependencyGroups"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: FindCircularDependencyGroups"
type: "source"
---

# Codebase Symbol: FindCircularDependencyGroups

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 16
- External references: 2
- Centrality: 0.164
- LOC: 94
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func FindCircularDependencyGroups(adjacency map[string][]string) [][]string {
```

## Documentation
FindCircularDependencyGroups returns deterministic strongly connected file
groups for the provided adjacency list. Groups with a single node are
ignored, so self-referential files are not treated as circular components.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquesortedstrings--internal-metrics-compute-go-l500]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/detectcirculardependencygroups--internal-metrics-compute-go-l368]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsmergesoverlappingcyclesintosinglecomponent--internal-metrics-compute-test-go-l220]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testfindcirculardependencygroupsreturnsstablesortedcomponents--internal-metrics-compute-test-go-l236]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `exports` (syntactic)
