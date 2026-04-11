---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 215
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 213
symbol_kind: "function"
symbol_name: "andFilter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: andFilter"
type: "source"
---

# Codebase Symbol: andFilter

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0615
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func andFilter(conditions ...models.BaseFilter) *models.BaseFilter {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderbasefiles--internal-vault-render-base-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
