---
blast_radius: 4
centrality: 0.1384
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 734
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 34
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 701
symbol_kind: "function"
symbol_name: "parseHumanDurationMilliseconds"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseHumanDurationMilliseconds"
type: "source"
---

# Codebase Symbol: parseHumanDurationMilliseconds

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 4
- External references: 1
- Centrality: 0.1384
- LOC: 34
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseHumanDurationMilliseconds(input string) (int, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testparsehumandurationmillisecondsparsesmultipleunits--internal-qmd-client-test-go-l463]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
