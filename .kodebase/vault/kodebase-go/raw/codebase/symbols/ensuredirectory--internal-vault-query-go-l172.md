---
blast_radius: 2
centrality: 0.0723
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 177
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 172
symbol_kind: "function"
symbol_name: "ensureDirectory"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ensureDirectory"
type: "source"
---

# Codebase Symbol: ensureDirectory

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0723
- LOC: 6
- Dead export: false
- Smells: None

## Signature
```text
func ensureDirectory(pathToCheck, label string) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdirectorypath--internal-vault-query-go-l210]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/listavailabletopics--internal-vault-query-go-l116]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvevaultquery--internal-vault-query-go-l56]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
