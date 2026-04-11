---
blast_radius: 3
centrality: 0.0631
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 298
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 262
symbol_kind: "function"
symbol_name: "renderBacklinkList"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderBacklinkList"
type: "source"
---

# Codebase Symbol: renderBacklinkList

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0631
- LOC: 37
- Dead export: false
- Smells: None

## Signature
```text
func renderBacklinkList(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/linkfornode--internal-vault-render-go-l199]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
