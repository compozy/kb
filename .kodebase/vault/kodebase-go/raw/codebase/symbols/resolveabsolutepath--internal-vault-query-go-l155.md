---
blast_radius: 4
centrality: 0.12
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 170
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 155
symbol_kind: "function"
symbol_name: "resolveAbsolutePath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveAbsolutePath"
type: "source"
---

# Codebase Symbol: resolveAbsolutePath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.12
- LOC: 16
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func resolveAbsolutePath(pathValue string) (string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/discovervaultpath--internal-vault-query-go-l28]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/listavailabletopics--internal-vault-query-go-l116]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvevaultquery--internal-vault-query-go-l56]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
