---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 123
exported: true
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 38
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 86
symbol_kind: "function"
symbol_name: "Boundaries"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: Boundaries"
type: "source"
---

# Codebase Symbol: Boundaries

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 38
- Dead export: false
- Smells: None

## Signature
```text
func Boundaries() error {
```

## Documentation
Boundaries verifies that package import rules are not violated.
Rules: no package may import cli/.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `exports` (syntactic)
