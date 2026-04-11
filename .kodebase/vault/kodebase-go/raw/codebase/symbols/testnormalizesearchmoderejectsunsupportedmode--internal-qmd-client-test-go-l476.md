---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 482
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 476
symbol_kind: "function"
symbol_name: "TestNormalizeSearchModeRejectsUnsupportedMode"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNormalizeSearchModeRejectsUnsupportedMode"
type: "source"
---

# Codebase Symbol: TestNormalizeSearchModeRejectsUnsupportedMode

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 7
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestNormalizeSearchModeRejectsUnsupportedMode(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizesearchmode--internal-qmd-client-go-l476]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `exports` (syntactic)
