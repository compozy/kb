---
blast_radius: 5
centrality: 0.0765
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 266
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 241
symbol_kind: "function"
symbol_name: "baseViewsValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: baseViewsValue"
type: "source"
---

# Codebase Symbol: baseViewsValue

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0765
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func baseViewsValue(views []models.BaseView) []interface{} {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stringmaptoany--internal-vault-render-base-go-l379]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stringslicetoany--internal-vault-render-base-go-l387]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
