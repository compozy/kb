---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 487
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 476
symbol_kind: "function"
symbol_name: "normalizeSearchMode"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeSearchMode"
type: "source"
---

# Codebase Symbol: normalizeSearchMode

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.137
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeSearchMode(mode SearchMode) (SearchMode, string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/searchcommand--internal-qmd-client-go-l307]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizesearchmoderejectsunsupportedmode--internal-qmd-client-test-go-l476]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
