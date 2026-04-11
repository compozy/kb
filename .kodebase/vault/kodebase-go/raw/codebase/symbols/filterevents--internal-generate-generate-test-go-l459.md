---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 467
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate_test.go"
stage: "raw"
start_line: 459
symbol_kind: "function"
symbol_name: "filterEvents"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: filterEvents"
type: "source"
---

# Codebase Symbol: filterEvents

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func filterEvents(events []Event, kind EventKind, stage string) []Event {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrunnergenerateemitsparseandwriteprogressevents--internal-generate-generate-test-go-l352]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `contains` (syntactic)
