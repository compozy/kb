---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 753
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 746
symbol_kind: "function"
symbol_name: "firstNonEmpty"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: firstNonEmpty"
type: "source"
---

# Codebase Symbol: firstNonEmpty

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.137
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func firstNonEmpty(values ...string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/normalize--internal-qmd-client-go-l459]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvesnippet--internal-qmd-client-go-l469]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
