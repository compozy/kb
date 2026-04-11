---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 198
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 33
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 166
symbol_kind: "function"
symbol_name: "buildWriteRequests"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildWriteRequests"
type: "source"
---

# Codebase Symbol: buildWriteRequests

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 33
- Dead export: false
- Smells: None

## Signature
```text
func buildWriteRequests(options WriteVaultOptions) ([]fileWriteRequest, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validatebasefile--internal-vault-writer-go-l242]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validaterendereddocument--internal-vault-writer-go-l200]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
