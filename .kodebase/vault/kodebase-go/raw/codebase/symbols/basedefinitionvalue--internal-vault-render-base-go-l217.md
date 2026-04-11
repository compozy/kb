---
blast_radius: 4
centrality: 0.0907
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 239
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 217
symbol_kind: "function"
symbol_name: "baseDefinitionValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: baseDefinitionValue"
type: "source"
---

# Codebase Symbol: baseDefinitionValue

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.0907
- LOC: 23
- Dead export: false
- Smells: None

## Signature
```text
func baseDefinitionValue(definition models.BaseDefinition) map[string]interface{} {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseviewsvalue--internal-vault-render-base-go-l241]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stringmaptoany--internal-vault-render-base-go-l379]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
