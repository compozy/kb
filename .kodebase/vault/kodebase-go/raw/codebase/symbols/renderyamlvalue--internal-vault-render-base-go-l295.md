---
blast_radius: 5
centrality: 0.0907
cyclomatic_complexity: 23
domain: "kodebase-go"
end_line: 362
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: true
language: "go"
loc: 68
outgoing_relation_count: 3
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 295
symbol_kind: "function"
symbol_name: "renderYAMLValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderYAMLValue"
type: "source"
---

# Codebase Symbol: renderYAMLValue

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 23
- Long function: true
- Blast radius: 5
- External references: 0
- Centrality: 0.0907
- LOC: 68
- Dead export: false
- Smells: `long-function`

## Signature
```text
func renderYAMLValue(value interface{}, indent int) []string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderyamlscalar--internal-vault-render-base-go-l364]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
