---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 137
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 18
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 120
symbol_kind: "function"
symbol_name: "HumanizeSlug"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: HumanizeSlug"
type: "source"
---

# Codebase Symbol: HumanizeSlug

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 18
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func HumanizeSlug(value string) string {
```

## Documentation
HumanizeSlug converts a hyphenated slug into a title-cased label.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/derivetopictitle--internal-vault-pathutils-go-l155]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
