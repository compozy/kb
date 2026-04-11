---
blast_radius: 4
centrality: 0.137
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 263
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader_test.go"
stage: "raw"
start_line: 247
symbol_kind: "function"
symbol_name: "createResolvedVault"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createResolvedVault"
type: "source"
---

# Codebase Symbol: createResolvedVault

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.137
- LOC: 17
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func createResolvedVault(t *testing.T) vault.ResolvedVault {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testfindsymbolsbynameusescaseinsensitivepartialmatch--internal-vault-reader-test-go-l200]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesfrontmatterandclassifiesdocuments--internal-vault-reader-test-go-l12]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesrelationsections--internal-vault-reader-test-go-l81]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotskipsmalformedyamlandwarns--internal-vault-reader-test-go-l137]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] via `contains` (syntactic)
