---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 170
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 113
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate_test.go"
stage: "raw"
start_line: 58
symbol_kind: "function"
symbol_name: "TestRunnerGenerateCallsPipelineStagesInOrder"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRunnerGenerateCallsPipelineStagesInOrder"
type: "source"
---

# Codebase Symbol: TestRunnerGenerateCallsPipelineStagesInOrder

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 113
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestRunnerGenerateCallsPipelineStagesInOrder(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testclock--internal-generate-generate-test-go-l478]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] via `exports` (syntactic)
