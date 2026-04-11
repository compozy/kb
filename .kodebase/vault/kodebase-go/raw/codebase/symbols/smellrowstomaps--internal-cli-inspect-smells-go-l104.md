---
blast_radius: 3
centrality: 0.0707
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 116
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_smells.go"
stage: "raw"
start_line: 104
symbol_kind: "function"
symbol_name: "smellRowsToMaps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: smellRowsToMaps"
type: "source"
---

# Codebase Symbol: smellRowsToMaps

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0707
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func smellRowsToMaps(rows []smellRow) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]] via `contains` (syntactic)
