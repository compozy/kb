---
blast_radius: 3
centrality: 0.0941
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 200
exported: true
external_reference_count: 2
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 197
symbol_kind: "function"
symbol_name: "RenderBaseDefinition"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: RenderBaseDefinition"
type: "source"
---

# Codebase Symbol: RenderBaseDefinition

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 3
- External references: 2
- Centrality: 0.0941
- LOC: 4
- Dead export: false
- Smells: None

## Signature
```text
func RenderBaseDefinition(definition models.BaseDefinition) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildwriterequests--internal-vault-writer-go-l166]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/validatebasefile--internal-vault-writer-go-l242]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `exports` (syntactic)
