---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 66
exported: true
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 64
symbol_kind: "function"
symbol_name: "TestIntegration"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestIntegration"
type: "source"
---

# Codebase Symbol: TestIntegration

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func TestIntegration() error {
```

## Documentation
TestIntegration runs all tests including integration tests.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rungotests--magefile-go-l200]]

## Backlinks
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `exports` (syntactic)
