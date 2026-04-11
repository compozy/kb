---
blast_radius: 8
centrality: 0.1518
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 291
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 272
symbol_kind: "function"
symbol_name: "testWritableTopicFixture"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testWritableTopicFixture"
type: "source"
---

# Codebase Symbol: testWritableTopicFixture

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 8
- External references: 0
- Centrality: 0.1518
- LOC: 20
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func testWritableTopicFixture(t *testing.T) models.TopicMetadata {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
