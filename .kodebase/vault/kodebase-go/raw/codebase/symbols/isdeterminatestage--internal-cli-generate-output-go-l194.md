---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 201
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 194
symbol_kind: "function"
symbol_name: "isDeterminateStage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isDeterminateStage"
type: "source"
---

# Codebase Symbol: isDeterminateStage

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func isDeterminateStage(event kgenerate.Event) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/handlestagestarted--internal-cli-generate-output-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
