---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 130
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 1
is_dead_export: false
is_long_function: true
language: "go"
loc: 68
outgoing_relation_count: 3
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/index.go"
stage: "raw"
start_line: 63
symbol_kind: "function"
symbol_name: "runIndexCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: runIndexCommand"
type: "source"
---

# Codebase Symbol: runIndexCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 68
- Dead export: false
- Smells: `long-function`

## Signature
```text
func runIndexCommand(cmd *cobra.Command, options *indexCommandOptions) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/chooseindexoperation--internal-cli-index-go-l151]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findcollectionstatus--internal-cli-index-go-l161]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/wrapqmdcommanderror--internal-cli-search-go-l196]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]] via `contains` (syntactic)
