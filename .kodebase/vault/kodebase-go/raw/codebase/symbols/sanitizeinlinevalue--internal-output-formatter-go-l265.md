---
blast_radius: 4
centrality: 0.0861
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 289
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 25
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 265
symbol_kind: "function"
symbol_name: "sanitizeInlineValue"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sanitizeInlineValue"
type: "source"
---

# Codebase Symbol: sanitizeInlineValue

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.0861
- LOC: 25
- Dead export: false
- Smells: None

## Signature
```text
func sanitizeInlineValue(value string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
