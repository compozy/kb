---
blast_radius: 21
centrality: 0.3243
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 407
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: true
language: "go"
loc: 121
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 287
symbol_kind: "function"
symbol_name: "testGraphFixture"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testGraphFixture"
type: "source"
---

# Codebase Symbol: testGraphFixture

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: true
- Blast radius: 21
- External references: 2
- Centrality: 0.3243
- LOC: 121
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`, `long-function`

## Signature
```text
func testGraphFixture() models.GraphSnapshot {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsintegrationproducesfulldocumentset--internal-vault-render-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderbasedefinitionsproducevalidyaml--internal-vault-render-test-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
