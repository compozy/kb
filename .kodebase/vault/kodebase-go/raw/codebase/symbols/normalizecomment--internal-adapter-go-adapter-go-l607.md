---
blast_radius: 15
centrality: 0.142
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 629
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 607
symbol_kind: "function"
symbol_name: "normalizeComment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeComment"
type: "source"
---

# Codebase Symbol: normalizeComment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 15
- External references: 0
- Centrality: 0.142
- LOC: 23
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeComment(rawComment string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecomment--internal-adapter-go-adapter-go-l631]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
