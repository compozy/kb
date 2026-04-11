---
blast_radius: 5
centrality: 0.2361
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 78
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 18
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 61
symbol_kind: "function"
symbol_name: "newRunner"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newRunner"
type: "source"
---

# Codebase Symbol: newRunner

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 5
- External references: 1
- Centrality: 0.2361
- LOC: 18
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func newRunner() runner {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l57]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgenerateintegrationbuildsvaultfromfixturerepository--internal-generate-generate-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
