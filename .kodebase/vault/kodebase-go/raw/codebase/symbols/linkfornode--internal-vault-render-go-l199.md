---
blast_radius: 5
centrality: 0.1581
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 219
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 199
symbol_kind: "function"
symbol_name: "linkForNode"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: linkForNode"
type: "source"
---

# Codebase Symbol: linkForNode

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1581
- LOC: 21
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func linkForNode(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderbacklinklist--internal-vault-render-go-l262]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrelationlist--internal-vault-render-go-l221]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
