---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 12
domain: "kodebase-go"
end_line: 155
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 54
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 102
symbol_kind: "function"
symbol_name: "TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog"
type: "source"
---

# Codebase Symbol: TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 12
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 54
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readfile--internal-vault-writer-test-go-l321]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `exports` (syntactic)
