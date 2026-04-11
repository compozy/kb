---
blast_radius: 25
centrality: 0.0678
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 74
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 70
symbol_kind: "function"
symbol_name: "bindInspectSharedFlags"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: bindInspectSharedFlags"
type: "source"
---

# Codebase Symbol: bindInspectSharedFlags

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 25
- External references: 0
- Centrality: 0.0678
- LOC: 5
- Dead export: false
- Smells: `high-blast-radius`

## Signature
```text
func bindInspectSharedFlags(flags *pflag.FlagSet, options *inspectSharedOptions) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
