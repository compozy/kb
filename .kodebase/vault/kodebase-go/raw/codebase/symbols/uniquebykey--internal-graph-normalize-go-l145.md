---
blast_radius: 9
centrality: 0.1145
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 160
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize.go"
stage: "raw"
start_line: 145
symbol_kind: "function"
symbol_name: "uniqueByKey"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: uniqueByKey"
type: "source"
---

# Codebase Symbol: uniqueByKey

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 9
- External references: 0
- Centrality: 0.1145
- LOC: 16
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func uniqueByKey[T any, K comparable](items []T, keyFn func(T) K) []T {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]] via `contains` (syntactic)
