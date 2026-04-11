---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 17
domain: "kodebase-go"
end_line: 416
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 1
is_dead_export: false
is_long_function: true
language: "go"
loc: 41
outgoing_relation_count: 0
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 376
symbol_kind: "function"
symbol_name: "frontmatterInt"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: frontmatterInt"
type: "source"
---

# Codebase Symbol: frontmatterInt

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 17
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 41
- Dead export: false
- Smells: `long-function`

## Signature
```text
func frontmatterInt(frontmatter map[string]interface{}, key string) int {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
