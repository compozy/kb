---
blast_radius: 18
centrality: 0.2727
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 252
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 250
symbol_kind: "function"
symbol_name: "normalizeDocumentPathSegment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: normalizeDocumentPathSegment"
type: "source"
---

# Codebase Symbol: normalizeDocumentPathSegment

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 18
- External references: 0
- Centrality: 0.2727
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func normalizeDocumentPathSegment(value string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
