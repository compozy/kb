---
blast_radius: 25
centrality: 0.0678
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 30
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_deadcode.go"
stage: "raw"
start_line: 19
symbol_kind: "function"
symbol_name: "newInspectDeadCodeCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newInspectDeadCodeCommand"
type: "source"
---

# Codebase Symbol: newInspectDeadCodeCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 25
- External references: 1
- Centrality: 0.0678
- LOC: 12
- Dead export: false
- Smells: `high-blast-radius`

## Signature
```text
func newInspectDeadCodeCommand(options *inspectSharedOptions) *cobra.Command {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]] via `contains` (syntactic)
