---
blast_radius: 13
centrality: 0.1693
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 337
exported: false
external_reference_count: 4
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 315
symbol_kind: "function"
symbol_name: "createInspectRelationRows"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createInspectRelationRows"
type: "source"
---

# Codebase Symbol: createInspectRelationRows

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 13
- External references: 4
- Centrality: 0.1693
- LOC: 23
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func createInspectRelationRows(relations []vault.VaultRelation) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tobacklinksoutput--internal-cli-inspect-backlinks-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/todependencyoutput--internal-cli-inspect-deps-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
