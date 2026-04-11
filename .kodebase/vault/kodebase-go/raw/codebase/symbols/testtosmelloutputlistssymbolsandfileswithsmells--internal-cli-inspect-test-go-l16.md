---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 53
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 38
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 16
symbol_kind: "function"
symbol_name: "TestToSmellOutputListsSymbolsAndFilesWithSmells"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToSmellOutputListsSymbolsAndFilesWithSmells"
type: "source"
---

# Codebase Symbol: TestToSmellOutputListsSymbolsAndFilesWithSmells

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 38
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestToSmellOutputListsSymbolsAndFilesWithSmells(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
