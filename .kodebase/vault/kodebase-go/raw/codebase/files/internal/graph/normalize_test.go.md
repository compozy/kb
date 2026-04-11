---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 24
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
symbol_count: 14
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/graph/normalize_test.go"
type: "source"
---

# Codebase File: internal/graph/normalize_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-test-go-l1|graph (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphreturnsemptysnapshotfornoparsedfiles--internal-graph-normalize-test-go-l10|TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphpassesthroughsingleparsedfile--internal-graph-normalize-test-go-l22|TestNormalizeGraphPassesThroughSingleParsedFile (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphdeduplicatesfilessymbolsexternalnodesandrelations--internal-graph-normalize-test-go-l58|TestNormalizeGraphDeduplicatesFilesSymbolsExternalNodesAndRelations (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphattachessymbolidstoparentfiles--internal-graph-normalize-test-go-l105|TestNormalizeGraphAttachesSymbolIDsToParentFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132|TestNormalizeGraphSortsCollectionsDeterministically (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphordersdiagnosticsbystagefilepathandmessage--internal-graph-normalize-test-go-l173|TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphomitsdiagnosticonlyfiles--internal-graph-normalize-test-go-l203|TestNormalizeGraphOmitsDiagnosticOnlyFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/assertemptygraphsnapshot--internal-graph-normalize-test-go-l232|assertEmptyGraphSnapshot (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/assertorderedids--internal-graph-normalize-test-go-l252|assertOrderedIDs (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/assertorderedsymbolids--internal-graph-normalize-test-go-l266|assertOrderedSymbolIDs (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/assertorderedexternalids--internal-graph-normalize-test-go-l280|assertOrderedExternalIDs (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsedfilefixture--internal-graph-normalize-test-go-l294|parsedFileFixture (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/capitalize--internal-graph-normalize-test-go-l365|capitalize (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertemptygraphsnapshot--internal-graph-normalize-test-go-l232]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedexternalids--internal-graph-normalize-test-go-l280]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedids--internal-graph-normalize-test-go-l252]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedsymbolids--internal-graph-normalize-test-go-l266]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/capitalize--internal-graph-normalize-test-go-l365]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedfilefixture--internal-graph-normalize-test-go-l294]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphattachessymbolidstoparentfiles--internal-graph-normalize-test-go-l105]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphdeduplicatesfilessymbolsexternalnodesandrelations--internal-graph-normalize-test-go-l58]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphomitsdiagnosticonlyfiles--internal-graph-normalize-test-go-l203]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphordersdiagnosticsbystagefilepathandmessage--internal-graph-normalize-test-go-l173]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphpassesthroughsingleparsedfile--internal-graph-normalize-test-go-l22]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphreturnsemptysnapshotfornoparsedfiles--internal-graph-normalize-test-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphattachessymbolidstoparentfiles--internal-graph-normalize-test-go-l105]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphdeduplicatesfilessymbolsexternalnodesandrelations--internal-graph-normalize-test-go-l58]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphomitsdiagnosticonlyfiles--internal-graph-normalize-test-go-l203]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphordersdiagnosticsbystagefilepathandmessage--internal-graph-normalize-test-go-l173]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphpassesthroughsingleparsedfile--internal-graph-normalize-test-go-l22]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphreturnsemptysnapshotfornoparsedfiles--internal-graph-normalize-test-go-l10]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `testing`

## Backlinks
None
