---
blast_radius: 5
centrality: 0.0893
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 377
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 364
symbol_kind: "function"
symbol_name: "renderYAMLScalar"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderYAMLScalar"
type: "source"
---

# Codebase Symbol: renderYAMLScalar

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0893
- LOC: 14
- Dead export: false
- Smells: None

## Signature
```text
func renderYAMLScalar(value interface{}) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
