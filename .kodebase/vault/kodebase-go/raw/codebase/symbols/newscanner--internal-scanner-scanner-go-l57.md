---
blast_radius: 13
centrality: 0.2634
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 66
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 57
symbol_kind: "function"
symbol_name: "NewScanner"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: NewScanner"
type: "source"
---

# Codebase Symbol: NewScanner

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.2634
- LOC: 10
- Dead export: true
- Smells: `bottleneck`, `dead-export`

## Signature
```text
func NewScanner(opts ...Option) *Scanner {
```

## Documentation
NewScanner constructs a scanner using functional options.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
