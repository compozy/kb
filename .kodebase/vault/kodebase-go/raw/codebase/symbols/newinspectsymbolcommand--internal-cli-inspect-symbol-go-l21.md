---
blast_radius: 25
centrality: 0.0678
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 39
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 0
smells:
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_symbol.go"
stage: "raw"
start_line: 21
symbol_kind: "function"
symbol_name: "newInspectSymbolCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newInspectSymbolCommand"
type: "source"
---

# Codebase Symbol: newInspectSymbolCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 25
- External references: 1
- Centrality: 0.0678
- LOC: 19
- Dead export: false
- Smells: `high-blast-radius`

## Signature
```text
func newInspectSymbolCommand(options *inspectSharedOptions) *cobra.Command {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]] via `contains` (syntactic)
