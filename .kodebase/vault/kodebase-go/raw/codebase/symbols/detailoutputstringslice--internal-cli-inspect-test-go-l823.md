---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 842
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 823
symbol_kind: "function"
symbol_name: "detailOutputStringSlice"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: detailOutputStringSlice"
type: "source"
---

# Codebase Symbol: detailOutputStringSlice

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 20
- Dead export: false
- Smells: None

## Signature
```text
func detailOutputStringSlice(t *testing.T, output inspectOutput, field string) []string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputvalue--internal-cli-inspect-test-go-l810]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
