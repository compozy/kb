---
blast_radius: 2
centrality: 0.0661
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 94
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_deadcode.go"
stage: "raw"
start_line: 81
symbol_kind: "function"
symbol_name: "deadCodeRowsToMaps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: deadCodeRowsToMaps"
type: "source"
---

# Codebase Symbol: deadCodeRowsToMaps

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0661
- LOC: 14
- Dead export: false
- Smells: None

## Signature
```text
func deadCodeRowsToMaps(rows []deadCodeRow) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/todeadcodeoutput--internal-cli-inspect-deadcode-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]] via `contains` (syntactic)
