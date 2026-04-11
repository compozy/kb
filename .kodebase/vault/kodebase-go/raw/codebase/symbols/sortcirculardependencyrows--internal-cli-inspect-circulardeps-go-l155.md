---
blast_radius: 5
centrality: 0.0982
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 161
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 155
symbol_kind: "function"
symbol_name: "sortCircularDependencyRows"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sortCircularDependencyRows"
type: "source"
---

# Codebase Symbol: sortCircularDependencyRows

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0982
- LOC: 7
- Dead export: false
- Smells: None

## Signature
```text
func sortCircularDependencyRows(rows []map[string]any) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
