---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 264
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 65
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 200
symbol_kind: "function"
symbol_name: "TestToSymbolLookupOutputReturnsDetailForSingleMatch"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToSymbolLookupOutputReturnsDetailForSingleMatch"
type: "source"
---

# Codebase Symbol: TestToSymbolLookupOutputReturnsDetailForSingleMatch

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 65
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestToSymbolLookupOutputReturnsDetailForSingleMatch(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymbollookupoutput--internal-cli-inspect-symbol-go-l41]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputvalue--internal-cli-inspect-test-go-l810]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
