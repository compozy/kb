---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 493
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate_test.go"
stage: "raw"
start_line: 478
symbol_kind: "function"
symbol_name: "testClock"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testClock"
type: "source"
---

# Codebase Symbol: testClock

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func testClock(instants ...time.Time) func() time.Time {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrunnergeneratecallspipelinestagesinorder--internal-generate-generate-test-go-l58]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `contains` (syntactic)
