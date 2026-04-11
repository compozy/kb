---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 107
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_integration_test.go"
stage: "raw"
start_line: 99
symbol_kind: "function"
symbol_name: "hasRelation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: hasRelation"
type: "source"
---

# Codebase Symbol: hasRelation

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func hasRelation(relations []models.RelationEdge, fromID string, toID string, relationType models.RelationType) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphmergesoverlappingimportsacrossparsedfiles--internal-graph-normalize-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]] via `contains` (syntactic)
