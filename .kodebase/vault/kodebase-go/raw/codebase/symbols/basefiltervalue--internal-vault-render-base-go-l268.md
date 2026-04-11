---
blast_radius: 7
centrality: 0.0981
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 293
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 268
symbol_kind: "function"
symbol_name: "baseFilterValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: baseFilterValue"
type: "source"
---

# Codebase Symbol: baseFilterValue

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.0981
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func baseFilterValue(filter models.BaseFilter) interface{} {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/baseviewsvalue--internal-vault-render-base-go-l241]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
