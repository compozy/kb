---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 37
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 1
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/textutils.go"
stage: "raw"
start_line: 18
symbol_kind: "function"
symbol_name: "NormalizeComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: NormalizeComment"
type: "source"
---

# Codebase Symbol: NormalizeComment

Source file: [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 20
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func NormalizeComment(rawComment string) string {
```

## Documentation
NormalizeComment strips Go/TS comment delimiters while preserving the comment text.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecommentblock--internal-vault-textutils-go-l78]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-vault-textutils-go-l40]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `exports` (syntactic)
