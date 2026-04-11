---
blast_radius: 1
centrality: 0.0594
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 169
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/search.go"
stage: "raw"
start_line: 149
symbol_kind: "function"
symbol_name: "resolveSearchCollection"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveSearchCollection"
type: "source"
---

# Codebase Symbol: resolveSearchCollection

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0594
- LOC: 21
- Dead export: false
- Smells: None

## Signature
```text
func resolveSearchCollection(options *searchCommandOptions) (string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/runsearchcommand--internal-cli-search-go-l71]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]] via `contains` (syntactic)
