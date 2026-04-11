---
blast_radius: 8
centrality: 0.0711
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 313
exported: false
external_reference_count: 2
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 300
symbol_kind: "function"
symbol_name: "createInspectDetailOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createInspectDetailOutput"
type: "source"
---

# Codebase Symbol: createInspectDetailOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 8
- External references: 2
- Centrality: 0.0711
- LOC: 14
- Dead export: false
- Smells: None

## Signature
```text
func createInspectDetailOutput(entries ...inspectDetailEntry) inspectOutput {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
