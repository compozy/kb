---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 160
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 155
symbol_kind: "method"
symbol_name: "handleStageProgress"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: handleStageProgress"
type: "source"
---

# Codebase Symbol: handleStageProgress

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: false
- Smells: None

## Signature
```text
func (o *generateTextObserver) handleStageProgress(event kgenerate.Event) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
