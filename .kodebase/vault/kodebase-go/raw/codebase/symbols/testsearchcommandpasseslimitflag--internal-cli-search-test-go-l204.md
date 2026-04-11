---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 240
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/search_test.go"
stage: "raw"
start_line: 204
symbol_kind: "function"
symbol_name: "TestSearchCommandPassesLimitFlag"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestSearchCommandPassesLimitFlag"
type: "source"
---

# Codebase Symbol: TestSearchCommandPassesLimitFlag

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 37
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestSearchCommandPassesLimitFlag(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]] via `exports` (syntactic)
