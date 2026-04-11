---
blast_radius: 3
centrality: 0.0667
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 105
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_blastradius.go"
stage: "raw"
start_line: 92
symbol_kind: "function"
symbol_name: "blastRadiusRowsToMaps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: blastRadiusRowsToMaps"
type: "source"
---

# Codebase Symbol: blastRadiusRowsToMaps

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0667
- LOC: 14
- Dead export: false
- Smells: None

## Signature
```text
func blastRadiusRowsToMaps(rows []blastRadiusRow) []map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]] via `contains` (syntactic)
