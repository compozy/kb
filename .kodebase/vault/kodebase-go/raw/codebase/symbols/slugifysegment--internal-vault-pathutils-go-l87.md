---
blast_radius: 15
centrality: 0.1548
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 117
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 5
is_dead_export: true
is_long_function: false
language: "go"
loc: 31
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 87
symbol_kind: "function"
symbol_name: "SlugifySegment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: SlugifySegment"
type: "source"
---

# Codebase Symbol: SlugifySegment

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 15
- External references: 0
- Centrality: 0.1548
- LOC: 31
- Dead export: true
- Smells: `bottleneck`, `dead-export`

## Signature
```text
func SlugifySegment(value string) string {
```

## Documentation
SlugifySegment converts a free-form segment into a filesystem-friendly slug.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-vault-pathutils-go-l169]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/derivetopicslug--internal-vault-pathutils-go-l140]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
