---
blast_radius: 4
centrality: 0.3182
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 358
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 322
symbol_kind: "function"
symbol_name: "normalizeFrontmatterValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeFrontmatterValue"
type: "source"
---

# Codebase Symbol: normalizeFrontmatterValue

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.3182
- LOC: 37
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeFrontmatterValue(value interface{}) interface{} {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattermap--internal-vault-reader-go-l314]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalizefrontmattermap--internal-vault-reader-go-l314]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
