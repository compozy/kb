---
blast_radius: 12
centrality: 0.18
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 285
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 275
symbol_kind: "function"
symbol_name: "testTopicFixture"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testTopicFixture"
type: "source"
---

# Codebase Symbol: testTopicFixture

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 12
- External references: 1
- Centrality: 0.18
- LOC: 11
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func testTopicFixture() models.TopicMetadata {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsintegrationproducesfulldocumentset--internal-vault-render-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
