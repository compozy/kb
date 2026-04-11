---
blast_radius: 5
centrality: 0.0609
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 409
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 397
symbol_kind: "function"
symbol_name: "inspectSectionText"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectSectionText"
type: "source"
---

# Codebase Symbol: inspectSectionText

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 5
- External references: 1
- Centrality: 0.0609
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func inspectSectionText(document vault.VaultDocument, heading string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
