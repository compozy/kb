---
blast_radius: 6
centrality: 0.1746
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 371
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 365
symbol_kind: "function"
symbol_name: "capitalize"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: capitalize"
type: "source"
---

# Codebase Symbol: capitalize

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 6
- External references: 0
- Centrality: 0.1746
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func capitalize(value string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsedfilefixture--internal-graph-normalize-test-go-l294]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
