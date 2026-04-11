---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 214
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 32
outgoing_relation_count: 4
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 183
symbol_kind: "function"
symbol_name: "TestSearchHybridModeUsesQueryCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestSearchHybridModeUsesQueryCommand"
type: "source"
---

# Codebase Symbol: TestSearchHybridModeUsesQueryCommand

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 32
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestSearchHybridModeUsesQueryCommand(t *testing.T) {
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
