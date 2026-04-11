---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 55
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/textutils.go"
stage: "raw"
start_line: 40
symbol_kind: "function"
symbol_name: "ExtractLeadingComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ExtractLeadingComment"
type: "source"
---

# Codebase Symbol: ExtractLeadingComment

Source file: [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 16
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func ExtractLeadingComment(sourceText string) string {
```

## Documentation
ExtractLeadingComment returns the first leading block or line comment from source text.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-vault-textutils-go-l18]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecommentblock--internal-vault-textutils-go-l78]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `exports` (syntactic)
