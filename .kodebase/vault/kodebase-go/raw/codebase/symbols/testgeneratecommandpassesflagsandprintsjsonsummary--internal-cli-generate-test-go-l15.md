---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 78
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 64
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_test.go"
stage: "raw"
start_line: 15
symbol_kind: "function"
symbol_name: "TestGenerateCommandPassesFlagsAndPrintsJSONSummary"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestGenerateCommandPassesFlagsAndPrintsJSONSummary"
type: "source"
---

# Codebase Symbol: TestGenerateCommandPassesFlagsAndPrintsJSONSummary

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 64
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestGenerateCommandPassesFlagsAndPrintsJSONSummary(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]] via `exports` (syntactic)
