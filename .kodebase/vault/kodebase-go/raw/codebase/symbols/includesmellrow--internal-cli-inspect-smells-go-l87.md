---
blast_radius: 3
centrality: 0.0707
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 102
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_smells.go"
stage: "raw"
start_line: 87
symbol_kind: "function"
symbol_name: "includeSmellRow"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: includeSmellRow"
type: "source"
---

# Codebase Symbol: includeSmellRow

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0707
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func includeSmellRow(smells []string, filter string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]] via `contains` (syntactic)
