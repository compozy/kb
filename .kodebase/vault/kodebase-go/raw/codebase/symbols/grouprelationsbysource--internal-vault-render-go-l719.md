---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 725
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 719
symbol_kind: "function"
symbol_name: "groupRelationsBySource"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: groupRelationsBySource"
type: "source"
---

# Codebase Symbol: groupRelationsBySource

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0524
- LOC: 7
- Dead export: false
- Smells: None

## Signature
```text
func groupRelationsBySource(relations []models.RelationEdge) map[string][]models.RelationEdge {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
