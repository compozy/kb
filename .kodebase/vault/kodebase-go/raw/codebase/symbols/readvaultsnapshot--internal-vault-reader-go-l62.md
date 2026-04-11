---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 109
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 48
outgoing_relation_count: 6
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 62
symbol_kind: "function"
symbol_name: "ReadVaultSnapshot"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ReadVaultSnapshot"
type: "source"
---

# Codebase Symbol: ReadVaultSnapshot

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 48
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func ReadVaultSnapshot(resolvedVault ResolvedVault, options ReadVaultOptions) (VaultSnapshot, error) {
```

## Documentation
ReadVaultSnapshot walks a topic directory and parses every managed markdown file.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/classifydocument--internal-vault-reader-go-l251]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectmarkdownfiles--internal-vault-reader-go-l192]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createemptysnapshot--internal-vault-reader-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortvaultdocuments--internal-vault-reader-go-l308]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `exports` (syntactic)
