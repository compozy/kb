---
blast_radius: 2
centrality: 0.1306
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 188
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "magefile.go"
stage: "raw"
start_line: 170
symbol_kind: "function"
symbol_name: "buildLDFlags"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildLDFlags"
type: "source"
---

# Codebase Symbol: buildLDFlags

Source file: [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.1306
- LOC: 19
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func buildLDFlags() string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gitoutput--magefile-go-l190]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildgo--magefile-go-l72]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] via `contains` (syntactic)
