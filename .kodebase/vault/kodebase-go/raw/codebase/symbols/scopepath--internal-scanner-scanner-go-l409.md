---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 420
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 409
symbol_kind: "function"
symbol_name: "scopePath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: scopePath"
type: "source"
---

# Codebase Symbol: scopePath

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func scopePath(relativeDirectory string, relativePath string) (string, bool) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/isignored--internal-scanner-scanner-go-l388]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
