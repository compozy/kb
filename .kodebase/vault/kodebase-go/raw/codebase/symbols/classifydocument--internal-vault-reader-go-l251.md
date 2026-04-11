---
blast_radius: 1
centrality: 0.0579
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 262
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 251
symbol_kind: "function"
symbol_name: "classifyDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: classifyDocument"
type: "source"
---

# Codebase Symbol: classifyDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0579
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func classifyDocument(frontmatter map[string]interface{}) vaultDocumentBucket {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/frontmatterstring--internal-vault-reader-go-l360]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
