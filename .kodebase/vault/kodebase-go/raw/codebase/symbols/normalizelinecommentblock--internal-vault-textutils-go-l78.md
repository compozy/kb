---
blast_radius: 2
centrality: 0.1338
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 85
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/textutils.go"
stage: "raw"
start_line: 78
symbol_kind: "function"
symbol_name: "normalizeLineCommentBlock"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeLineCommentBlock"
type: "source"
---

# Codebase Symbol: normalizeLineCommentBlock

Source file: [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.1338
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeLineCommentBlock(rawComment string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-vault-textutils-go-l40]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-vault-textutils-go-l18]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] via `contains` (syntactic)
