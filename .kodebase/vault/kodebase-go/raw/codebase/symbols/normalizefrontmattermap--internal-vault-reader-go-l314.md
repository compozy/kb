---
blast_radius: 4
centrality: 0.3228
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 320
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 314
symbol_kind: "function"
symbol_name: "normalizeFrontmatterMap"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeFrontmatterMap"
type: "source"
---

# Codebase Symbol: normalizeFrontmatterMap

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.3228
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeFrontmatterMap(values map[string]interface{}) map[string]interface{} {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
