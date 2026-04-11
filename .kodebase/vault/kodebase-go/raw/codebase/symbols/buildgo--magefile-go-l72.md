---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 82
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 72
symbol_kind: "function"
symbol_name: "buildGo"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildGo"
type: "source"
---

# Codebase Symbol: buildGo

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func buildGo() error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildldflags--magefile-go-l170]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/build--magefile-go-l68]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
