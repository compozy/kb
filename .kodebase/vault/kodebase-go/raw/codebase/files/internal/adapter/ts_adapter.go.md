---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 55
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
symbol_count: 44
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/ts_adapter.go"
type: "source"
---

# Codebase File: internal/adapter/ts_adapter.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 2
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tsadapter--internal-adapter-ts-adapter-go-l32|TSAdapter (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/tscalltarget--internal-adapter-ts-adapter-go-l34|tsCallTarget (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tsimportbinding--internal-adapter-ts-adapter-go-l40|tsImportBinding (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tslocalexport--internal-adapter-ts-adapter-go-l47|tsLocalExport (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tsreexport--internal-adapter-ts-adapter-go-l52|tsReExport (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tssymbolmatch--internal-adapter-ts-adapter-go-l58|tsSymbolMatch (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsedtsfile--internal-adapter-ts-adapter-go-l65|parsedTSFile (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-ts-adapter-go-l78|Supports (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-ts-adapter-go-l88|ParseFiles (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93|ParseFilesWithProgress (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299|parseTSFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/selecttslanguage--internal-adapter-ts-adapter-go-l378|selectTSLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391|extractTSImports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484|extractTSExports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583|extractCommonJSExports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644|extractRequireBindings (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766|createTSSymbolMatch (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782|extractClassMethodSymbols (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805|extractVariableSymbols (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826|createTSSymbol (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/gettssymbolkind--internal-adapter-ts-adapter-go-l859|getTSSymbolKind (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880|resolveTSSymbolName (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916|resolveTSVariableName (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formattssignature--internal-adapter-ts-adapter-go-l929|formatTSSignature (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formattsreturntype--internal-adapter-ts-adapter-go-l968|formatTSReturnType (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formattsvariabletypesuffix--internal-adapter-ts-adapter-go-l986|formatTSVariableTypeSuffix (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extracttsdoccomment--internal-adapter-ts-adapter-go-l1004|extractTSDocComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collecttscalltargets--internal-adapter-ts-adapter-go-l1018|collectTSCallTargets (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/computetscyclomaticcomplexity--internal-adapter-ts-adapter-go-l1069|computeTSCyclomaticComplexity (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/istsexecutableboundary--internal-adapter-ts-adapter-go-l1107|isTSExecutableBoundary (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/applylocalexportstate--internal-adapter-ts-adapter-go-l1116|applyLocalExportState (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvelocalexports--internal-adapter-ts-adapter-go-l1151|resolveLocalExports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvereexports--internal-adapter-ts-adapter-go-l1170|resolveReExports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolverelativeimportfile--internal-adapter-ts-adapter-go-l1207|resolveRelativeImportFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createtsparsediagnostic--internal-adapter-ts-adapter-go-l1253|createTSParseDiagnostic (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265|declarationExportNames (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isdefaultexport--internal-adapter-ts-adapter-go-l1285|isDefaultExport (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290|matchCommonJSExportTarget (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/findchildbykind--internal-adapter-ts-adapter-go-l1333|findChildByKind (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/relationkey--internal-adapter-ts-adapter-go-l1348|relationKey (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/pushuniquerelation--internal-adapter-ts-adapter-go-l1357|pushUniqueRelation (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizeabsolutepath--internal-adapter-ts-adapter-go-l1371|normalizeAbsolutePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/clonestringset--internal-adapter-ts-adapter-go-l1380|cloneStringSet (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/applylocalexportstate--internal-adapter-ts-adapter-go-l1116]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/clonestringset--internal-adapter-ts-adapter-go-l1380]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collecttscalltargets--internal-adapter-ts-adapter-go-l1018]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computetscyclomaticcomplexity--internal-adapter-ts-adapter-go-l1069]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtsparsediagnostic--internal-adapter-ts-adapter-go-l1253]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsdoccomment--internal-adapter-ts-adapter-go-l1004]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findchildbykind--internal-adapter-ts-adapter-go-l1333]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsreturntype--internal-adapter-ts-adapter-go-l968]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattssignature--internal-adapter-ts-adapter-go-l929]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsvariabletypesuffix--internal-adapter-ts-adapter-go-l986]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gettssymbolkind--internal-adapter-ts-adapter-go-l859]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdefaultexport--internal-adapter-ts-adapter-go-l1285]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/istsexecutableboundary--internal-adapter-ts-adapter-go-l1107]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizeabsolutepath--internal-adapter-ts-adapter-go-l1371]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedtsfile--internal-adapter-ts-adapter-go-l65]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-ts-adapter-go-l88]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/pushuniquerelation--internal-adapter-ts-adapter-go-l1357]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationkey--internal-adapter-ts-adapter-go-l1348]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvelocalexports--internal-adapter-ts-adapter-go-l1151]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvereexports--internal-adapter-ts-adapter-go-l1170]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolverelativeimportfile--internal-adapter-ts-adapter-go-l1207]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/selecttslanguage--internal-adapter-ts-adapter-go-l378]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-ts-adapter-go-l78]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsadapter--internal-adapter-ts-adapter-go-l32]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tscalltarget--internal-adapter-ts-adapter-go-l34]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsimportbinding--internal-adapter-ts-adapter-go-l40]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tslocalexport--internal-adapter-ts-adapter-go-l47]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsreexport--internal-adapter-ts-adapter-go-l52]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tssymbolmatch--internal-adapter-ts-adapter-go-l58]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-ts-adapter-go-l88]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-ts-adapter-go-l78]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsadapter--internal-adapter-ts-adapter-go-l32]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `tree_sitter (github.com/tree-sitter/go-tree-sitter)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
