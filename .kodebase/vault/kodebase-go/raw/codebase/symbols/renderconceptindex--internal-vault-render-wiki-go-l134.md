---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 176
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 43
outgoing_relation_count: 3
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 134
symbol_kind: "function"
symbol_name: "renderConceptIndex"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderConceptIndex"
type: "source"
---

# Codebase Symbol: renderConceptIndex

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0524
- LOC: 43
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func renderConceptIndex(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
