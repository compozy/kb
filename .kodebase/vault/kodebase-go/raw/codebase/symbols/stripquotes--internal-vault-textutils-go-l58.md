---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 76
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 1
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/textutils.go"
stage: "raw"
start_line: 58
symbol_kind: "function"
symbol_name: "StripQuotes"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: StripQuotes"
type: "source"
---

# Codebase Symbol: StripQuotes

Source file: [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 19
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func StripQuotes(value string) string {
```

## Documentation
StripQuotes removes a single leading and trailing quote character when present.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isquote--internal-vault-textutils-go-l87]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `exports` (syntactic)
