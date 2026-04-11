---
blast_radius: 4
centrality: 0.109
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 800
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 792
symbol_kind: "function"
symbol_name: "writeInspectMarkdown"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeInspectMarkdown"
type: "source"
---

# Codebase Symbol: writeInspectMarkdown

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.109
- LOC: 9
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func writeInspectMarkdown(t *testing.T, topicPath, relativePath, content string) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-cli-inspect-test-go-l802]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createinspecttestvault--internal-cli-inspect-test-go-l641]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
