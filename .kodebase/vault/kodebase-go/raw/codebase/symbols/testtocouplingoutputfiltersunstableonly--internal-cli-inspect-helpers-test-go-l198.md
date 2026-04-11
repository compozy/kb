---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 212
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_helpers_test.go"
stage: "raw"
start_line: 198
symbol_kind: "function"
symbol_name: "TestToCouplingOutputFiltersUnstableOnly"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToCouplingOutputFiltersUnstableOnly"
type: "source"
---

# Codebase Symbol: TestToCouplingOutputFiltersUnstableOnly

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 15
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestToCouplingOutputFiltersUnstableOnly(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] via `exports` (syntactic)
