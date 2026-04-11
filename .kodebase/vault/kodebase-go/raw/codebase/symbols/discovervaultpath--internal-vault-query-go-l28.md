---
blast_radius: 3
centrality: 0.1122
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 53
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 2
smells:
  - "bottleneck"
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 28
symbol_kind: "function"
symbol_name: "DiscoverVaultPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: DiscoverVaultPath"
type: "source"
---

# Codebase Symbol: DiscoverVaultPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1122
- LOC: 26
- Dead export: true
- Smells: `bottleneck`, `dead-export`

## Signature
```text
func DiscoverVaultPath(cwd string) (string, error) {
```

## Documentation
DiscoverVaultPath walks up from cwd until it finds a `.kodebase/vault` directory.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdirectorypath--internal-vault-query-go-l210]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveabsolutepath--internal-vault-query-go-l155]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/resolvevaultpath--internal-vault-query-go-l138]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `exports` (syntactic)
