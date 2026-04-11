---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 148
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 52
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 97
symbol_kind: "function"
symbol_name: "TestIndexWithContextAndEmbedRunsExpectedCommands"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestIndexWithContextAndEmbedRunsExpectedCommands"
type: "source"
---

# Codebase Symbol: TestIndexWithContextAndEmbedRunsExpectedCommands

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 52
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestIndexWithContextAndEmbedRunsExpectedCommands(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readinvocationlog--internal-qmd-client-test-go-l537]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefakeqmd--internal-qmd-client-test-go-l490]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `exports` (syntactic)
