---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 49
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/qmd/client_test.go"
stage: "raw"
symbol_count: 23
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/qmd/client_test.go"
type: "source"
---

# Codebase File: internal/qmd/client_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-test-go-l1|qmd (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testsearchreturnserrqmdunavailableformissingbinary--internal-qmd-client-test-go-l14|TestSearchReturnsErrQMDUnavailableForMissingBinary (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27|TestIndexAddConstructsExpectedArguments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexupdateconstructsexpectedarguments--internal-qmd-client-test-go-l64|TestIndexUpdateConstructsExpectedArguments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexwithcontextandembedrunsexpectedcommands--internal-qmd-client-test-go-l97|TestIndexWithContextAndEmbedRunsExpectedCommands (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexrejectsinvalidinputs--internal-qmd-client-test-go-l150|TestIndexRejectsInvalidInputs (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchhybridmodeusesquerycommand--internal-qmd-client-test-go-l183|TestSearchHybridModeUsesQueryCommand (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchallomitslimitandacceptsmodealias--internal-qmd-client-test-go-l216|TestSearchAllOmitsLimitAndAcceptsModeAlias (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchpasseslimitminscoreandfullflags--internal-qmd-client-test-go-l246|TestSearchPassesLimitMinScoreAndFullFlags (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchparsesjsonandnormalizesresults--internal-qmd-client-test-go-l283|TestSearchParsesJSONAndNormalizesResults (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchfullusesbodywhenpresent--internal-qmd-client-test-go-l318|TestSearchFullUsesBodyWhenPresent (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcontextcancellationstopsrunningcommand--internal-qmd-client-test-go-l343|TestSearchContextCancellationStopsRunningCommand (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchfailureincludesstderrdiagnostics--internal-qmd-client-test-go-l365|TestSearchFailureIncludesStderrDiagnostics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparseupdateresultparsesaddandupdatesummaries--internal-qmd-client-test-go-l388|TestParseUpdateResultParsesAddAndUpdateSummaries (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparseembedresultparsessuccessandnowork--internal-qmd-client-test-go-l408|TestParseEmbedResultParsesSuccessAndNoWork (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparseindexstatusparsescollectionsandhealth--internal-qmd-client-test-go-l428|TestParseIndexStatusParsesCollectionsAndHealth (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparseindexstatusacceptsemptyindex--internal-qmd-client-test-go-l447|TestParseIndexStatusAcceptsEmptyIndex (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparsehumandurationmillisecondsparsesmultipleunits--internal-qmd-client-test-go-l463|TestParseHumanDurationMillisecondsParsesMultipleUnits (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnormalizesearchmoderejectsunsupportedmode--internal-qmd-client-test-go-l476|TestNormalizeSearchModeRejectsUnsupportedMode (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/fakeqmdoptions--internal-qmd-client-test-go-l484|fakeQMDOptions (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writefakeqmd--internal-qmd-client-test-go-l490|writeFakeQMD (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/readinvocationlog--internal-qmd-client-test-go-l537|readInvocationLog (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/shellquote--internal-qmd-client-test-go-l572|shellQuote (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fakeqmdoptions--internal-qmd-client-test-go-l484]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readinvocationlog--internal-qmd-client-test-go-l537]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/shellquote--internal-qmd-client-test-go-l572]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexrejectsinvalidinputs--internal-qmd-client-test-go-l150]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexupdateconstructsexpectedarguments--internal-qmd-client-test-go-l64]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexwithcontextandembedrunsexpectedcommands--internal-qmd-client-test-go-l97]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizesearchmoderejectsunsupportedmode--internal-qmd-client-test-go-l476]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseembedresultparsessuccessandnowork--internal-qmd-client-test-go-l408]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparsehumandurationmillisecondsparsesmultipleunits--internal-qmd-client-test-go-l463]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseindexstatusacceptsemptyindex--internal-qmd-client-test-go-l447]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseindexstatusparsescollectionsandhealth--internal-qmd-client-test-go-l428]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseupdateresultparsesaddandupdatesummaries--internal-qmd-client-test-go-l388]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchallomitslimitandacceptsmodealias--internal-qmd-client-test-go-l216]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcontextcancellationstopsrunningcommand--internal-qmd-client-test-go-l343]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchfailureincludesstderrdiagnostics--internal-qmd-client-test-go-l365]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchfullusesbodywhenpresent--internal-qmd-client-test-go-l318]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchhybridmodeusesquerycommand--internal-qmd-client-test-go-l183]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchparsesjsonandnormalizesresults--internal-qmd-client-test-go-l283]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchpasseslimitminscoreandfullflags--internal-qmd-client-test-go-l246]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchreturnserrqmdunavailableformissingbinary--internal-qmd-client-test-go-l14]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefakeqmd--internal-qmd-client-test-go-l490]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexrejectsinvalidinputs--internal-qmd-client-test-go-l150]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexupdateconstructsexpectedarguments--internal-qmd-client-test-go-l64]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexwithcontextandembedrunsexpectedcommands--internal-qmd-client-test-go-l97]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizesearchmoderejectsunsupportedmode--internal-qmd-client-test-go-l476]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseembedresultparsessuccessandnowork--internal-qmd-client-test-go-l408]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparsehumandurationmillisecondsparsesmultipleunits--internal-qmd-client-test-go-l463]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseindexstatusacceptsemptyindex--internal-qmd-client-test-go-l447]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseindexstatusparsescollectionsandhealth--internal-qmd-client-test-go-l428]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseupdateresultparsesaddandupdatesummaries--internal-qmd-client-test-go-l388]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchallomitslimitandacceptsmodealias--internal-qmd-client-test-go-l216]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcontextcancellationstopsrunningcommand--internal-qmd-client-test-go-l343]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchfailureincludesstderrdiagnostics--internal-qmd-client-test-go-l365]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchfullusesbodywhenpresent--internal-qmd-client-test-go-l318]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchhybridmodeusesquerycommand--internal-qmd-client-test-go-l183]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchparsesjsonandnormalizesresults--internal-qmd-client-test-go-l283]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchpasseslimitminscoreandfullflags--internal-qmd-client-test-go-l246]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchreturnserrqmdunavailableformissingbinary--internal-qmd-client-test-go-l14]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`
- `imports` (syntactic) -> `time`

## Backlinks
None
