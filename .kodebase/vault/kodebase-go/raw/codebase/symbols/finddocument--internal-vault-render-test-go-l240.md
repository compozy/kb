---
blast_radius: 10
centrality: 0.2485
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 251
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 240
symbol_kind: "function"
symbol_name: "findDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: findDocument"
type: "source"
---

# Codebase Symbol: findDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 10
- External references: 1
- Centrality: 0.2485
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func findDocument(t *testing.T, documents []models.RenderedDocument, relativePath string) models.RenderedDocument {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsintegrationproducesfulldocumentset--internal-vault-render-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscirculardependencieslistsgroups--internal-vault-render-test-go-l142]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscodebaseoverviewcontainssummary--internal-vault-render-test-go-l116]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdashboardlinkstoallconceptarticles--internal-vault-render-test-go-l153]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdependencyhotspotsliststopfiles--internal-vault-render-test-go-l131]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdirectoryindexuseswikilinks--internal-vault-render-test-go-l97]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawfilefrontmatterandbody--internal-vault-render-test-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawsymbolfrontmatterandsignature--internal-vault-render-test-go-l70]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsusetopicwikilinksyntax--internal-vault-render-test-go-l221]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
