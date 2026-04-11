---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 113
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 58
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 56
symbol_kind: "function"
symbol_name: "ResolveVaultQuery"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ResolveVaultQuery"
type: "source"
---

# Codebase Symbol: ResolveVaultQuery

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 58
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func ResolveVaultQuery(options VaultQueryOptions) (ResolvedVault, error) {
```

## Documentation
ResolveVaultQuery resolves the target vault and topic from the provided options.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuredirectory--internal-vault-query-go-l172]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/listtopicdirectories--internal-vault-query-go-l179]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveabsolutepath--internal-vault-query-go-l155]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvevaultpath--internal-vault-query-go-l138]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `exports` (syntactic)
