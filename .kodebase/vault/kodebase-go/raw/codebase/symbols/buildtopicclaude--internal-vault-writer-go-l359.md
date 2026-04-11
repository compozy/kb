---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 440
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 82
outgoing_relation_count: 3
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 359
symbol_kind: "function"
symbol_name: "buildTopicClaude"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildTopicClaude"
type: "source"
---

# Codebase Symbol: buildTopicClaude

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 82
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func buildTopicClaude(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
