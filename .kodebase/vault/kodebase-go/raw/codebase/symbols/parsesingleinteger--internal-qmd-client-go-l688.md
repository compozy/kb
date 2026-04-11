---
blast_radius: 5
centrality: 0.1319
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 699
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
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
start_line: 688
symbol_kind: "function"
symbol_name: "parseSingleInteger"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseSingleInteger"
type: "source"
---

# Codebase Symbol: parseSingleInteger

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1319
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseSingleInteger(pattern *regexp.Regexp, input string) (int, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
