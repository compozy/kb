---
blast_radius: 24
centrality: 0.2208
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 69
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 29
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/search.go"
stage: "raw"
start_line: 41
symbol_kind: "function"
symbol_name: "newSearchCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newSearchCommand"
type: "source"
---

# Codebase Symbol: newSearchCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 24
- External references: 1
- Centrality: 0.2208
- LOC: 29
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func newSearchCommand() *cobra.Command {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]] via `contains` (syntactic)
