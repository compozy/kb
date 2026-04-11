---
blast_radius: 7
centrality: 0.2377
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 270
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 8
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 260
symbol_kind: "function"
symbol_name: "testWriteVaultInputs"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testWriteVaultInputs"
type: "source"
---

# Codebase Symbol: testWriteVaultInputs

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 7
- External references: 2
- Centrality: 0.2377
- LOC: 11
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func testWriteVaultInputs(t *testing.T) (models.TopicMetadata, models.GraphSnapshot, []models.RenderedDocument, []models.BaseFile) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgraphfixture--internal-vault-render-test-go-l287]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritabletopicfixture--internal-vault-writer-test-go-l272]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotroundtripswriteroutput--internal-vault-reader-integration-test-go-l12]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultintegrationpersistsfullrenderedvault--internal-vault-writer-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatesclaudemanifestandappendonlylog--internal-vault-writer-test-go-l102]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultrejectsinvalidrendereddocument--internal-vault-writer-test-go-l235]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultreportsprogressforpersistedfiles--internal-vault-writer-test-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
