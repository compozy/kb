---
blast_radius: 1
centrality: 0.0579
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 249
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 27
outgoing_relation_count: 4
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 223
symbol_kind: "function"
symbol_name: "parseVaultDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseVaultDocument"
type: "source"
---

# Codebase Symbol: parseVaultDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0579
- LOC: 27
- Dead export: false
- Smells: None

## Signature
```text
func parseVaultDocument(markdown, relativePath string, warn func(string)) (VaultDocument, bool) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractsection--internal-vault-reader-go-l112]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattermap--internal-vault-reader-go-l314]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsebacklinks--internal-vault-reader-go-l286]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parserelations--internal-vault-reader-go-l264]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
