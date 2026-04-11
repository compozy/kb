---
blast_radius: 5
centrality: 0.2017
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 808
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 802
symbol_kind: "function"
symbol_name: "mkdirAll"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: mkdirAll"
type: "source"
---

# Codebase Symbol: mkdirAll

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.2017
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func mkdirAll(t *testing.T, path string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createinspecttestvault--internal-cli-inspect-test-go-l641]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/writeinspectmarkdown--internal-cli-inspect-test-go-l792]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
