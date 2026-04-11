---
blast_radius: 3
centrality: 0.0789
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 342
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 335
symbol_kind: "function"
symbol_name: "normalizePattern"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizePattern"
type: "source"
---

# Codebase Symbol: normalizePattern

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0789
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func normalizePattern(pattern string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/builduserrules--internal-scanner-scanner-go-l308]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
