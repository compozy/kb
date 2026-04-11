---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 93
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 22
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/query_test.go"
stage: "raw"
start_line: 72
symbol_kind: "function"
symbol_name: "TestResolveVaultQueryAutoResolvesSingleTopic"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestResolveVaultQueryAutoResolvesSingleTopic"
type: "source"
---

# Codebase Symbol: TestResolveVaultQueryAutoResolvesSingleTopic

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 22
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestResolveVaultQueryAutoResolvesSingleTopic(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-vault-query-test-go-l206]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetopicmarker--internal-vault-query-test-go-l214]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `exports` (syntactic)
