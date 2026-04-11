---
blast_radius: 2
centrality: 0.1306
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 145
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 134
symbol_kind: "function"
symbol_name: "parseInspectOutputFormat"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseInspectOutputFormat"
type: "source"
---

# Codebase Symbol: parseInspectOutputFormat

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.1306
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseInspectOutputFormat(value string) (output.OutputFormat, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/resolveinspectcontext--internal-cli-inspect-go-l103]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
