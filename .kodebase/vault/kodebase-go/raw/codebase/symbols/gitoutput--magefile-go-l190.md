---
blast_radius: 3
centrality: 0.1618
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 198
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
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
start_line: 190
symbol_kind: "function"
symbol_name: "gitOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: gitOutput"
type: "source"
---

# Codebase Symbol: gitOutput

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1618
- LOC: 9
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func gitOutput(args ...string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildldflags--magefile-go-l170]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
