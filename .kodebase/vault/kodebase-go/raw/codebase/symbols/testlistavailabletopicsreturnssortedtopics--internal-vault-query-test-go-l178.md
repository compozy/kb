---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 204
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 27
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/query_test.go"
stage: "raw"
start_line: 178
symbol_kind: "function"
symbol_name: "TestListAvailableTopicsReturnsSortedTopics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestListAvailableTopicsReturnsSortedTopics"
type: "source"
---

# Codebase Symbol: TestListAvailableTopicsReturnsSortedTopics

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 27
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestListAvailableTopicsReturnsSortedTopics(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-vault-query-test-go-l206]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetopicmarker--internal-vault-query-test-go-l214]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `exports` (syntactic)
