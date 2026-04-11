---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 153
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 22
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 132
symbol_kind: "method"
symbol_name: "handleStageStarted"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: handleStageStarted"
type: "source"
---

# Codebase Symbol: handleStageStarted

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 22
- Dead export: false
- Smells: None

## Signature
```text
func (o *generateTextObserver) handleStageStarted(event kgenerate.Event) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/humanstagelabel--internal-cli-generate-output-go-l203]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdeterminatestage--internal-cli-generate-output-go-l194]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
