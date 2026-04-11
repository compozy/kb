---
blast_radius: 6
centrality: 0.0981
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 385
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 379
symbol_kind: "function"
symbol_name: "stringMapToAny"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: stringMapToAny"
type: "source"
---

# Codebase Symbol: stringMapToAny

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 6
- External references: 0
- Centrality: 0.0981
- LOC: 7
- Dead export: false
- Smells: None

## Signature
```text
func stringMapToAny(values map[string]string) map[string]interface{} {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/baseviewsvalue--internal-vault-render-base-go-l241]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
