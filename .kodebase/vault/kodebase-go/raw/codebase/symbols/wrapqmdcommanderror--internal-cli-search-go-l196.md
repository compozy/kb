---
blast_radius: 2
centrality: 0.0738
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 206
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/search.go"
stage: "raw"
start_line: 196
symbol_kind: "function"
symbol_name: "wrapQMDCommandError"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: wrapQMDCommandError"
type: "source"
---

# Codebase Symbol: wrapQMDCommandError

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.0738
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func wrapQMDCommandError(command string, err error) error {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/runindexcommand--internal-cli-index-go-l63]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/runsearchcommand--internal-cli-search-go-l71]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]] via `contains` (syntactic)
