---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 11
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/textutils.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/textutils.go"
type: "source"
---

# Codebase File: internal/vault/textutils.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-textutils-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-vault-textutils-go-l18|NormalizeComment (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-vault-textutils-go-l40|ExtractLeadingComment (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/stripquotes--internal-vault-textutils-go-l58|StripQuotes (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/normalizelinecommentblock--internal-vault-textutils-go-l78|normalizeLineCommentBlock (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isquote--internal-vault-textutils-go-l87|isQuote (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-vault-textutils-go-l40]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isquote--internal-vault-textutils-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-vault-textutils-go-l18]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecommentblock--internal-vault-textutils-go-l78]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-vault-textutils-go-l58]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-textutils-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-vault-textutils-go-l40]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-vault-textutils-go-l18]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-vault-textutils-go-l58]]
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `strings`

## Backlinks
None
