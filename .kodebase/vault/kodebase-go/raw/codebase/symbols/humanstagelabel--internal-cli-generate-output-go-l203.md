---
blast_radius: 3
centrality: 0.137
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 205
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 203
symbol_kind: "function"
symbol_name: "humanStageLabel"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: humanStageLabel"
type: "source"
---

# Codebase Symbol: humanStageLabel

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.137
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func humanStageLabel(stage string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/handlestagecompleted--internal-cli-generate-output-go-l162]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/handlestagefailed--internal-cli-generate-output-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/handlestagestarted--internal-cli-generate-output-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
