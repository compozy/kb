---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 412
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 28
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 385
symbol_kind: "function"
symbol_name: "TestToBacklinksOutputListsIncomingReferencesForSymbol"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToBacklinksOutputListsIncomingReferencesForSymbol"
type: "source"
---

# Codebase Symbol: TestToBacklinksOutputListsIncomingReferencesForSymbol

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 28
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestToBacklinksOutputListsIncomingReferencesForSymbol(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tobacklinksoutput--internal-cli-inspect-backlinks-go-l31]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
