---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 371
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 60
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 312
symbol_kind: "function"
symbol_name: "TestToFileLookupOutputIncludesContainedSymbolsAndMetrics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToFileLookupOutputIncludesContainedSymbolsAndMetrics"
type: "source"
---

# Codebase Symbol: TestToFileLookupOutputIncludesContainedSymbolsAndMetrics

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 60
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestToFileLookupOutputIncludesContainedSymbolsAndMetrics(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputstringslice--internal-cli-inspect-test-go-l823]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputvalue--internal-cli-inspect-test-go-l810]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
