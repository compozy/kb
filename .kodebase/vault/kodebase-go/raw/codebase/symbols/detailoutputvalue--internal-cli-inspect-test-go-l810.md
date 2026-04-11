---
blast_radius: 3
centrality: 0.1421
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 821
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 810
symbol_kind: "function"
symbol_name: "detailOutputValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: detailOutputValue"
type: "source"
---

# Codebase Symbol: detailOutputValue

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1421
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func detailOutputValue(t *testing.T, output inspectOutput, field string) any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/detailoutputstringslice--internal-cli-inspect-test-go-l823]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdetailforsinglematch--internal-cli-inspect-test-go-l200]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
