---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 130
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 41
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 90
symbol_kind: "function"
symbol_name: "TestToComplexityOutputSortsByDescendingComplexity"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToComplexityOutputSortsByDescendingComplexity"
type: "source"
---

# Codebase Symbol: TestToComplexityOutputSortsByDescendingComplexity

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 41
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestToComplexityOutputSortsByDescendingComplexity(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
