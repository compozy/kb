---
blast_radius: 4
centrality: 0.0957
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 94
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 40
outgoing_relation_count: 3
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_symbol.go"
stage: "raw"
start_line: 55
symbol_kind: "function"
symbol_name: "toSymbolSummaryOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toSymbolSummaryOutput"
type: "source"
---

# Codebase Symbol: toSymbolSummaryOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.0957
- LOC: 40
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toSymbolSummaryOutput(documents []vault.VaultDocument) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tosymbollookupoutput--internal-cli-inspect-symbol-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]] via `contains` (syntactic)
