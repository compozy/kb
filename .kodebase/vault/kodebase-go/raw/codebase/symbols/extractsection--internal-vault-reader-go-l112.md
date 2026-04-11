---
blast_radius: 2
centrality: 0.0631
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 143
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 32
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 112
symbol_kind: "function"
symbol_name: "ExtractSection"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ExtractSection"
type: "source"
---

# Codebase Symbol: ExtractSection

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0631
- LOC: 32
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func ExtractSection(body, heading string) string {
```

## Documentation
ExtractSection returns the markdown content under the named level-two heading.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `exports` (syntactic)
