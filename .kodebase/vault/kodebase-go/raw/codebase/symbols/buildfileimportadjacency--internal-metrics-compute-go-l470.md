---
blast_radius: 14
centrality: 0.0777
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 498
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 29
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 470
symbol_kind: "function"
symbol_name: "buildFileImportAdjacency"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildFileImportAdjacency"
type: "source"
---

# Codebase Symbol: buildFileImportAdjacency

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 14
- External references: 0
- Centrality: 0.0777
- LOC: 29
- Dead export: false
- Smells: None

## Signature
```text
func buildFileImportAdjacency(files []models.GraphFile, relations []models.RelationEdge) map[string][]string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/detectcirculardependencygroups--internal-metrics-compute-go-l368]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
