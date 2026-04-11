---
blast_radius: 16
centrality: 0.1973
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 646
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 631
symbol_kind: "function"
symbol_name: "normalizeLineComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeLineComment"
type: "source"
---

# Codebase Symbol: normalizeLineComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 16
- External references: 0
- Centrality: 0.1973
- LOC: 16
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeLineComment(rawComment string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-adapter-go-adapter-go-l607]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
