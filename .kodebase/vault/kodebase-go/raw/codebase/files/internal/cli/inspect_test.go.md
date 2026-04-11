---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 11
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.9167
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 57
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
symbol_count: 28
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_test.go"
type: "source"
---

# Codebase File: internal/cli/inspect_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 11
- Instability: 0.9167
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-test-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputlistssymbolsandfileswithsmells--internal-cli-inspect-test-go-l16|TestToSmellOutputListsSymbolsAndFilesWithSmells (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtodeadcodeoutputlistsdeadexportsandorphanfiles--internal-cli-inspect-test-go-l55|TestToDeadCodeOutputListsDeadExportsAndOrphanFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocomplexityoutputsortsbydescendingcomplexity--internal-cli-inspect-test-go-l90|TestToComplexityOutputSortsByDescendingComplexity (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputsortsbydescendingblastradius--internal-cli-inspect-test-go-l132|TestToBlastRadiusOutputSortsByDescendingBlastRadius (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputsortsbyinstability--internal-cli-inspect-test-go-l166|TestToCouplingOutputSortsByInstability (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdetailforsinglematch--internal-cli-inspect-test-go-l200|TestToSymbolLookupOutputReturnsDetailForSingleMatch (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnssummaryformultiplematches--internal-cli-inspect-test-go-l266|TestToSymbolLookupOutputReturnsSummaryForMultipleMatches (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdescriptiveerrorforunknownname--internal-cli-inspect-test-go-l300|TestToSymbolLookupOutputReturnsDescriptiveErrorForUnknownName (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312|TestToFileLookupOutputIncludesContainedSymbolsAndMetrics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputreturnsdescriptiveerrorforunknownpath--internal-cli-inspect-test-go-l373|TestToFileLookupOutputReturnsDescriptiveErrorForUnknownPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtobacklinksoutputlistsincomingreferencesforsymbol--internal-cli-inspect-test-go-l385|TestToBacklinksOutputListsIncomingReferencesForSymbol (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtodependencyoutputlistsoutgoingdependenciesforsymbol--internal-cli-inspect-test-go-l414|TestToDependencyOutputListsOutgoingDependenciesForSymbol (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtodependencyoutputsupportsexactfilepathlookup--internal-cli-inspect-test-go-l442|TestToDependencyOutputSupportsExactFilePathLookup (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputlistsfileswithcirculardependencyflag--internal-cli-inspect-test-go-l467|TestToCircularDepsOutputListsFilesWithCircularDependencyFlag (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputfallsbacktosccdetectionforlegacyvaults--internal-cli-inspect-test-go-l508|TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputshowsmessagewhennocycles--internal-cli-inspect-test-go-l538|TestToCircularDepsOutputShowsMessageWhenNoCycles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectcommandjsonformatproducesvalidjson--internal-cli-inspect-test-go-l550|TestInspectCommandJSONFormatProducesValidJSON (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectcommandtsvformatproducesheaderandrows--internal-cli-inspect-test-go-l572|TestInspectCommandTSVFormatProducesHeaderAndRows (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectcommandreturnsdescriptiveerrorformissingvault--internal-cli-inspect-test-go-l594|TestInspectCommandReturnsDescriptiveErrorForMissingVault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectsubcommandsrespondtohelp--internal-cli-inspect-test-go-l611|TestInspectSubcommandsRespondToHelp (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637|testVaultDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createinspecttestvault--internal-cli-inspect-test-go-l641|createInspectTestVault (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writeinspectmarkdown--internal-cli-inspect-test-go-l792|writeInspectMarkdown (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/mkdirall--internal-cli-inspect-test-go-l802|mkdirAll (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/detailoutputvalue--internal-cli-inspect-test-go-l810|detailOutputValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/detailoutputstringslice--internal-cli-inspect-test-go-l823|detailOutputStringSlice (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testfiledocumentforcycle--internal-cli-inspect-test-go-l844|testFileDocumentForCycle (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspecttestvault--internal-cli-inspect-test-go-l641]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputstringslice--internal-cli-inspect-test-go-l823]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/detailoutputvalue--internal-cli-inspect-test-go-l810]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-cli-inspect-test-go-l802]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfiledocumentforcycle--internal-cli-inspect-test-go-l844]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandjsonformatproducesvalidjson--internal-cli-inspect-test-go-l550]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandreturnsdescriptiveerrorformissingvault--internal-cli-inspect-test-go-l594]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandtsvformatproducesheaderandrows--internal-cli-inspect-test-go-l572]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectsubcommandsrespondtohelp--internal-cli-inspect-test-go-l611]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtobacklinksoutputlistsincomingreferencesforsymbol--internal-cli-inspect-test-go-l385]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputsortsbydescendingblastradius--internal-cli-inspect-test-go-l132]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputfallsbacktosccdetectionforlegacyvaults--internal-cli-inspect-test-go-l508]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputlistsfileswithcirculardependencyflag--internal-cli-inspect-test-go-l467]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputshowsmessagewhennocycles--internal-cli-inspect-test-go-l538]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocomplexityoutputsortsbydescendingcomplexity--internal-cli-inspect-test-go-l90]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocouplingoutputsortsbyinstability--internal-cli-inspect-test-go-l166]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodeadcodeoutputlistsdeadexportsandorphanfiles--internal-cli-inspect-test-go-l55]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodependencyoutputlistsoutgoingdependenciesforsymbol--internal-cli-inspect-test-go-l414]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodependencyoutputsupportsexactfilepathlookup--internal-cli-inspect-test-go-l442]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputreturnsdescriptiveerrorforunknownpath--internal-cli-inspect-test-go-l373]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosmelloutputlistssymbolsandfileswithsmells--internal-cli-inspect-test-go-l16]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdescriptiveerrorforunknownname--internal-cli-inspect-test-go-l300]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdetailforsinglematch--internal-cli-inspect-test-go-l200]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnssummaryformultiplematches--internal-cli-inspect-test-go-l266]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writeinspectmarkdown--internal-cli-inspect-test-go-l792]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandjsonformatproducesvalidjson--internal-cli-inspect-test-go-l550]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandreturnsdescriptiveerrorformissingvault--internal-cli-inspect-test-go-l594]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandtsvformatproducesheaderandrows--internal-cli-inspect-test-go-l572]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectsubcommandsrespondtohelp--internal-cli-inspect-test-go-l611]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtobacklinksoutputlistsincomingreferencesforsymbol--internal-cli-inspect-test-go-l385]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputsortsbydescendingblastradius--internal-cli-inspect-test-go-l132]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputfallsbacktosccdetectionforlegacyvaults--internal-cli-inspect-test-go-l508]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputlistsfileswithcirculardependencyflag--internal-cli-inspect-test-go-l467]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputshowsmessagewhennocycles--internal-cli-inspect-test-go-l538]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocomplexityoutputsortsbydescendingcomplexity--internal-cli-inspect-test-go-l90]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocouplingoutputsortsbyinstability--internal-cli-inspect-test-go-l166]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodeadcodeoutputlistsdeadexportsandorphanfiles--internal-cli-inspect-test-go-l55]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodependencyoutputlistsoutgoingdependenciesforsymbol--internal-cli-inspect-test-go-l414]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtodependencyoutputsupportsexactfilepathlookup--internal-cli-inspect-test-go-l442]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputreturnsdescriptiveerrorforunknownpath--internal-cli-inspect-test-go-l373]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosmelloutputlistssymbolsandfileswithsmells--internal-cli-inspect-test-go-l16]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdescriptiveerrorforunknownname--internal-cli-inspect-test-go-l300]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnsdetailforsinglematch--internal-cli-inspect-test-go-l200]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosymbollookupoutputreturnssummaryformultiplematches--internal-cli-inspect-test-go-l266]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
