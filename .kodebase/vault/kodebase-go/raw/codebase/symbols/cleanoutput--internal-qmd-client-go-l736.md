---
blast_radius: 9
centrality: 0.2655
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 738
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 736
symbol_kind: "function"
symbol_name: "cleanOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: cleanOutput"
type: "source"
---

# Codebase Symbol: cleanOutput

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 9
- External references: 0
- Centrality: 0.2655
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func cleanOutput(output string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleandiagnostics--internal-qmd-client-go-l740]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parseupdateresult--internal-qmd-client-go-l502]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
