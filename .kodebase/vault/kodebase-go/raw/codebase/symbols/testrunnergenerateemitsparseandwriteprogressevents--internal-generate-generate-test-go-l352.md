---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 457
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 106
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate_test.go"
stage: "raw"
start_line: 352
symbol_kind: "function"
symbol_name: "TestRunnerGenerateEmitsParseAndWriteProgressEvents"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRunnerGenerateEmitsParseAndWriteProgressEvents"
type: "source"
---

# Codebase Symbol: TestRunnerGenerateEmitsParseAndWriteProgressEvents

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 106
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestRunnerGenerateEmitsParseAndWriteProgressEvents(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filterevents--internal-generate-generate-test-go-l459]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/firstevent--internal-generate-generate-test-go-l469]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `exports` (syntactic)
