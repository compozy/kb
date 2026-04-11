---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 132
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 1
is_dead_export: false
is_long_function: true
language: "go"
loc: 62
outgoing_relation_count: 5
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/search.go"
stage: "raw"
start_line: 71
symbol_kind: "function"
symbol_name: "runSearchCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: runSearchCommand"
type: "source"
---

# Codebase Symbol: runSearchCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 62
- Dead export: false
- Smells: `long-function`

## Signature
```text
func runSearchCommand(cmd *cobra.Command, query string, options *searchCommandOptions) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesearchoutputformat--internal-cli-search-go-l171]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvesearchcollection--internal-cli-search-go-l149]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvesearchmode--internal-cli-search-go-l134]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchresultstorows--internal-cli-search-go-l184]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/wrapqmdcommanderror--internal-cli-search-go-l196]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]] via `contains` (syntactic)
