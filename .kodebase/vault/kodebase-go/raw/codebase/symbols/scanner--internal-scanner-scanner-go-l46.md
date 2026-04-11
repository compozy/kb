---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 48
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 46
symbol_kind: "struct"
symbol_name: "Scanner"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: Scanner"
type: "source"
---

# Codebase Symbol: Scanner

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 3
- Dead export: true
- Smells: `dead-export`

## Signature
```text
Scanner struct {
```

## Documentation
Scanner discovers supported source files inside a repository workspace.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
