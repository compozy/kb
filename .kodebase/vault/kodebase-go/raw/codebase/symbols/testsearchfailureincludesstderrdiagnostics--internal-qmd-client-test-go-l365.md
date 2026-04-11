---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 386
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 22
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 365
symbol_kind: "function"
symbol_name: "TestSearchFailureIncludesStderrDiagnostics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestSearchFailureIncludesStderrDiagnostics"
type: "source"
---

# Codebase Symbol: TestSearchFailureIncludesStderrDiagnostics

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 22
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestSearchFailureIncludesStderrDiagnostics(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefakeqmd--internal-qmd-client-test-go-l490]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `exports` (syntactic)
