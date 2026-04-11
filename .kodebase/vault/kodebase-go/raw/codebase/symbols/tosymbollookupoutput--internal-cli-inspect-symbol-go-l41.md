---
blast_radius: 3
centrality: 0.1586
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 53
exported: false
external_reference_count: 3
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 3
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_symbol.go"
stage: "raw"
start_line: 41
symbol_kind: "function"
symbol_name: "toSymbolLookupOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toSymbolLookupOutput"
type: "source"
---

# Codebase Symbol: toSymbolLookupOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 3
- External references: 3
- Centrality: 0.1586
- LOC: 13
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func toSymbolLookupOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsingleinspectsymbolmatch--internal-cli-inspect-go-l350]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymbolsummaryoutput--internal-cli-inspect-symbol-go-l55]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdescriptiveerrorforunknownname--internal-cli-inspect-test-go-l300]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdetailforsinglematch--internal-cli-inspect-test-go-l200]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnssummaryformultiplematches--internal-cli-inspect-test-go-l266]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]] via `contains` (syntactic)
