---
blast_radius: 3
centrality: 0.0641
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 88
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
source_path: "internal/cli/inspect_coupling.go"
stage: "raw"
start_line: 75
symbol_kind: "function"
symbol_name: "couplingRowsToMaps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: couplingRowsToMaps"
type: "source"
---

# Codebase Symbol: couplingRowsToMaps

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0641
- LOC: 14
- Dead export: false
- Smells: None

## Signature
```text
func couplingRowsToMaps(rows []couplingRow) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]] via `contains` (syntactic)
