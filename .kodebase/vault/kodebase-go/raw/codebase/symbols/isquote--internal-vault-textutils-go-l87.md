---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 89
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/textutils.go"
stage: "raw"
start_line: 87
symbol_kind: "function"
symbol_name: "isQuote"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isQuote"
type: "source"
---

# Codebase Symbol: isQuote

Source file: [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func isQuote(char byte) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/stripquotes--internal-vault-textutils-go-l58]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `contains` (syntactic)
