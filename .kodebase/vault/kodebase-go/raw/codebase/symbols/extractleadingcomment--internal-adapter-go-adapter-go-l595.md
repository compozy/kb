---
blast_radius: 4
centrality: 0.0607
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 605
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 595
symbol_kind: "function"
symbol_name: "extractLeadingComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractLeadingComment"
type: "source"
---

# Codebase Symbol: extractLeadingComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 4
- External references: 1
- Centrality: 0.0607
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func extractLeadingComment(sourceText string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-adapter-go-adapter-go-l607]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecomment--internal-adapter-go-adapter-go-l631]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
