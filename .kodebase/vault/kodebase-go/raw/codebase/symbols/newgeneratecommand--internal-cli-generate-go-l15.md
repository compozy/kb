---
blast_radius: 24
centrality: 0.2208
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 66
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 52
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/generate.go"
stage: "raw"
start_line: 15
symbol_kind: "function"
symbol_name: "newGenerateCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newGenerateCommand"
type: "source"
---

# Codebase Symbol: newGenerateCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate.go|internal/cli/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: true
- Blast radius: 24
- External references: 1
- Centrality: 0.2208
- LOC: 52
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`, `long-function`

## Signature
```text
func newGenerateCommand() *cobra.Command {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate.go|internal/cli/generate.go]] via `contains` (syntactic)
