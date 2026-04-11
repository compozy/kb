---
blast_radius: 4
centrality: 0.191
cyclomatic_complexity: 30
domain: "kodebase-go"
end_line: 686
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: true
language: "go"
loc: 89
outgoing_relation_count: 2
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 598
symbol_kind: "function"
symbol_name: "parseIndexStatus"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseIndexStatus"
type: "source"
---

# Codebase Symbol: parseIndexStatus

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 30
- Long function: true
- Blast radius: 4
- External references: 2
- Centrality: 0.191
- LOC: 89
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func parseIndexStatus(output string) (IndexStatus, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesingleinteger--internal-qmd-client-go-l688]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/status--internal-qmd-client-go-l183]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testparseindexstatusacceptsemptyindex--internal-qmd-client-test-go-l447]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testparseindexstatusparsescollectionsandhealth--internal-qmd-client-test-go-l428]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
