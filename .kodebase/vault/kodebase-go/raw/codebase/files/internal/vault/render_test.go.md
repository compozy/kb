---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 34
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/vault/render_test.go"
stage: "raw"
symbol_count: 17
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/render_test.go"
type: "source"
---

# Codebase File: internal/vault/render_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-render-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13|TestRenderDocumentsProducesRawWikiAndBaseSurfaces (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawfilefrontmatterandbody--internal-vault-render-test-go-l43|TestRenderDocumentsRawFileFrontmatterAndBody (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawsymbolfrontmatterandsignature--internal-vault-render-test-go-l70|TestRenderDocumentsRawSymbolFrontmatterAndSignature (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdirectoryindexuseswikilinks--internal-vault-render-test-go-l97|TestRenderDocumentsDirectoryIndexUsesWikiLinks (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscodebaseoverviewcontainssummary--internal-vault-render-test-go-l116|TestRenderDocumentsCodebaseOverviewContainsSummary (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdependencyhotspotsliststopfiles--internal-vault-render-test-go-l131|TestRenderDocumentsDependencyHotspotsListsTopFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentscirculardependencieslistsgroups--internal-vault-render-test-go-l142|TestRenderDocumentsCircularDependenciesListsGroups (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdashboardlinkstoallconceptarticles--internal-vault-render-test-go-l153|TestRenderDocumentsDashboardLinksToAllConceptArticles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderbasedefinitionsproducevalidyaml--internal-vault-render-test-go-l178|TestRenderBaseDefinitionsProduceValidYAML (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsbodieshavevalidfrontmatterandkinds--internal-vault-render-test-go-l196|TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsusetopicwikilinksyntax--internal-vault-render-test-go-l221|TestRenderDocumentsUseTopicWikiLinkSyntax (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232|renderFixtureDocuments (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/finddocument--internal-vault-render-test-go-l240|findDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsefrontmatter--internal-vault-render-test-go-l253|parseFrontmatter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testtopicfixture--internal-vault-render-test-go-l275|testTopicFixture (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testgraphfixture--internal-vault-render-test-go-l287|testGraphFixture (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/finddocument--internal-vault-render-test-go-l240]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefrontmatter--internal-vault-render-test-go-l253]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgraphfixture--internal-vault-render-test-go-l287]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderbasedefinitionsproducevalidyaml--internal-vault-render-test-go-l178]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsbodieshavevalidfrontmatterandkinds--internal-vault-render-test-go-l196]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentscirculardependencieslistsgroups--internal-vault-render-test-go-l142]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentscodebaseoverviewcontainssummary--internal-vault-render-test-go-l116]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdashboardlinkstoallconceptarticles--internal-vault-render-test-go-l153]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdependencyhotspotsliststopfiles--internal-vault-render-test-go-l131]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdirectoryindexuseswikilinks--internal-vault-render-test-go-l97]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawfilefrontmatterandbody--internal-vault-render-test-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawsymbolfrontmatterandsignature--internal-vault-render-test-go-l70]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsusetopicwikilinksyntax--internal-vault-render-test-go-l221]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopicfixture--internal-vault-render-test-go-l275]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-render-test-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderbasedefinitionsproducevalidyaml--internal-vault-render-test-go-l178]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsbodieshavevalidfrontmatterandkinds--internal-vault-render-test-go-l196]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentscirculardependencieslistsgroups--internal-vault-render-test-go-l142]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentscodebaseoverviewcontainssummary--internal-vault-render-test-go-l116]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdashboardlinkstoallconceptarticles--internal-vault-render-test-go-l153]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdependencyhotspotsliststopfiles--internal-vault-render-test-go-l131]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsdirectoryindexuseswikilinks--internal-vault-render-test-go-l97]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsproducesrawwikiandbasesurfaces--internal-vault-render-test-go-l13]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawfilefrontmatterandbody--internal-vault-render-test-go-l43]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsrawsymbolfrontmatterandsignature--internal-vault-render-test-go-l70]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrenderdocumentsusetopicwikilinksyntax--internal-vault-render-test-go-l221]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/metrics`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `gopkg.in/yaml.v3`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
