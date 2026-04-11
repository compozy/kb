---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 474
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 469
symbol_kind: "method"
symbol_name: "resolveSnippet"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: resolveSnippet"
type: "source"
---

# Codebase Symbol: resolveSnippet

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: false
- Smells: None

## Signature
```text
func (payload searchResultPayload) resolveSnippet(full bool) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/firstnonempty--internal-qmd-client-go-l746]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
