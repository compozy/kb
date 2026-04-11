---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 136
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 4
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 116
symbol_kind: "function"
symbol_name: "ListAvailableTopics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ListAvailableTopics"
type: "source"
---

# Codebase Symbol: ListAvailableTopics

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 21
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func ListAvailableTopics(options VaultQueryOptions) ([]string, error) {
```

## Documentation
ListAvailableTopics returns the marker-backed topic directories in deterministic order.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuredirectory--internal-vault-query-go-l172]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/listtopicdirectories--internal-vault-query-go-l179]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveabsolutepath--internal-vault-query-go-l155]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvevaultpath--internal-vault-query-go-l138]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `exports` (syntactic)
