---
blast_radius: 3
centrality: 0.1246
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 175
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 3
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/output/formatter.go"
stage: "raw"
start_line: 157
symbol_kind: "function"
symbol_name: "projectStringRows"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: projectStringRows"
type: "source"
---

# Codebase Symbol: projectStringRows

Source file: [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1246
- LOC: 19
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func projectStringRows(columns []string, data []map[string]any, truncate bool) [][]string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sanitizeinlinevalue--internal-output-formatter-go-l265]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/truncatetablecell--internal-output-formatter-go-l291]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattsv--internal-output-formatter-go-l141]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] via `contains` (syntactic)
