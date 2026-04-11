---
blast_radius: 9
centrality: 0.2449
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 238
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 232
symbol_kind: "function"
symbol_name: "renderFixtureDocuments"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderFixtureDocuments"
type: "source"
---

# Codebase Symbol: renderFixtureDocuments

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 9
- External references: 0
- Centrality: 0.2449
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func renderFixtureDocuments(t *testing.T) []models.RenderedDocument {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgraphfixture--internal-vault-render-test-go-l287]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopicfixture--internal-vault-render-test-go-l275]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsbodieshavevalidfrontmatterandkinds--internal-vault-render-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscirculardependencieslistsgroups--internal-vault-render-test-go-l142]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscodebaseoverviewcontainssummary--internal-vault-render-test-go-l116]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdashboardlinkstoallconceptarticles--internal-vault-render-test-go-l153]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdependencyhotspotsliststopfiles--internal-vault-render-test-go-l131]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdirectoryindexuseswikilinks--internal-vault-render-test-go-l97]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawfilefrontmatterandbody--internal-vault-render-test-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawsymbolfrontmatterandsignature--internal-vault-render-test-go-l70]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsusetopicwikilinksyntax--internal-vault-render-test-go-l221]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
