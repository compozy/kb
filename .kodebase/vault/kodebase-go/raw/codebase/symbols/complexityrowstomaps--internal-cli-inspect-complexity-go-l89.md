---
blast_radius: 2
centrality: 0.063
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 103
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_complexity.go"
stage: "raw"
start_line: 89
symbol_kind: "function"
symbol_name: "complexityRowsToMaps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: complexityRowsToMaps"
type: "source"
---

# Codebase Symbol: complexityRowsToMaps

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.063
- LOC: 15
- Dead export: false
- Smells: None

## Signature
```text
func complexityRowsToMaps(rows []complexityRow) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]] via `contains` (syntactic)
