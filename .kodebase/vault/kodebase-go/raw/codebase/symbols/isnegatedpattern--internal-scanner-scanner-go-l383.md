---
blast_radius: 5
centrality: 0.1147
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 386
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 383
symbol_kind: "function"
symbol_name: "isNegatedPattern"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isNegatedPattern"
type: "source"
---

# Codebase Symbol: isNegatedPattern

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1147
- LOC: 4
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func isNegatedPattern(pattern string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
