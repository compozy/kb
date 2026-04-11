---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 176
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 162
symbol_kind: "method"
symbol_name: "handleStageCompleted"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: handleStageCompleted"
type: "source"
---

# Codebase Symbol: handleStageCompleted

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 15
- Dead export: false
- Smells: None

## Signature
```text
func (o *generateTextObserver) handleStageCompleted(event kgenerate.Event) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/completedsuffix--internal-cli-generate-output-go-l207]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/humanstagelabel--internal-cli-generate-output-go-l203]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
