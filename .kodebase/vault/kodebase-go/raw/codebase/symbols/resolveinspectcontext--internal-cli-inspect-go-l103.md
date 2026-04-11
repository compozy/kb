---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 132
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 30
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 103
symbol_kind: "function"
symbol_name: "resolveInspectContext"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveInspectContext"
type: "source"
---

# Codebase Symbol: resolveInspectContext

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 30
- Dead export: false
- Smells: None

## Signature
```text
func resolveInspectContext(options *inspectSharedOptions) (inspectContext, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseinspectoutputformat--internal-cli-inspect-go-l134]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/runinspectcommand--internal-cli-inspect-go-l76]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
