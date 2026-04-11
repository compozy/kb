---
blast_radius: 3
centrality: 0.1432
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 374
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 360
symbol_kind: "function"
symbol_name: "frontmatterString"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: frontmatterString"
type: "source"
---

# Codebase Symbol: frontmatterString

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1432
- LOC: 15
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func frontmatterString(frontmatter map[string]interface{}, key string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/classifydocument--internal-vault-reader-go-l251]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/findsymbolsbyname--internal-vault-reader-go-l146]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
