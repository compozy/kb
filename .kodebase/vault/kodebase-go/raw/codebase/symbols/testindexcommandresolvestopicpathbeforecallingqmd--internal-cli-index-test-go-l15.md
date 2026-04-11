---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 98
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 84
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/index_test.go"
stage: "raw"
start_line: 15
symbol_kind: "function"
symbol_name: "TestIndexCommandResolvesTopicPathBeforeCallingQMD"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestIndexCommandResolvesTopicPathBeforeCallingQMD"
type: "source"
---

# Codebase Symbol: TestIndexCommandResolvesTopicPathBeforeCallingQMD

Source file: [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 84
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestIndexCommandResolvesTopicPathBeforeCallingQMD(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `exports` (syntactic)
