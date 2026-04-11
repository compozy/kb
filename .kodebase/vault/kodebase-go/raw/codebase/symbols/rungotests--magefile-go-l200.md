---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 208
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 200
symbol_kind: "function"
symbol_name: "runGoTests"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: runGoTests"
type: "source"
---

# Codebase Symbol: runGoTests

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.137
- LOC: 9
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func runGoTests(testArgs ...string) error {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/test--magefile-go-l59]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testintegration--magefile-go-l64]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
