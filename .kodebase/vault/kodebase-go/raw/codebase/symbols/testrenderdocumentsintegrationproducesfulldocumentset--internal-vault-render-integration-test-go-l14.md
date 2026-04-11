---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 61
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 48
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_integration_test.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "TestRenderDocumentsIntegrationProducesFullDocumentSet"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRenderDocumentsIntegrationProducesFullDocumentSet"
type: "source"
---

# Codebase Symbol: TestRenderDocumentsIntegrationProducesFullDocumentSet

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_integration_test.go|internal/vault/render_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 48
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestRenderDocumentsIntegrationProducesFullDocumentSet(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/finddocument--internal-vault-render-test-go-l240]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefrontmatter--internal-vault-render-test-go-l253]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgraphfixture--internal-vault-render-test-go-l287]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopicfixture--internal-vault-render-test-go-l275]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render_integration_test.go|internal/vault/render_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_integration_test.go|internal/vault/render_integration_test.go]] via `exports` (syntactic)
