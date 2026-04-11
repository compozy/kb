---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 168
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 143
symbol_kind: "function"
symbol_name: "goFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: goFiles"
type: "source"
---

# Codebase Symbol: goFiles

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func goFiles(root string) ([]string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/fmt--magefile-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
