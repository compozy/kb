---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 331
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 306
symbol_kind: "function"
symbol_name: "writeFilesInBatches"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeFilesInBatches"
type: "source"
---

# Codebase Symbol: writeFilesInBatches

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func writeFilesInBatches(ctx context.Context, files []fileWriteRequest, report func(string)) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetextfile--internal-vault-writer-go-l598]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
