---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.5
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 29
smells:
source_kind: "codebase-file"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
symbol_count: 14
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/ts_adapter_test.go"
type: "source"
---

# Codebase File: internal/adapter/ts_adapter_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 1
- Instability: 0.5
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-test-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testtsadaptersupportstslikelanguages--internal-adapter-ts-adapter-test-go-l13|TestTSAdapterSupportsTSLikeLanguages (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29|TestTSAdapterParsesSimpleTypeScriptFile (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78|TestTSAdapterParsesTSXComponent (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99|TestTSAdapterParsesJavaScriptRequireImports (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140|TestTSAdapterExtractsClassAndMethodSymbols (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177|TestTSAdapterBuildsImportBindingsForDefaultNamedAndNamespaceImports (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236|TestTSAdapterHandlesNamedAndStarReExports (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadaptercomputescyclomaticcomplexity--internal-adapter-ts-adapter-test-go-l285|TestTSAdapterComputesCyclomaticComplexity (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtsadapterproducesdiagnosticsforparseerrors--internal-adapter-ts-adapter-test-go-l306|TestTSAdapterProducesDiagnosticsForParseErrors (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsesingletslikefile--internal-adapter-ts-adapter-test-go-l329|parseSingleTSLikeFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsetslikesources--internal-adapter-ts-adapter-test-go-l340|parseTSLikeSources (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/mustfindparsedfile--internal-adapter-ts-adapter-test-go-l377|mustFindParsedFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/languageforpath--internal-adapter-ts-adapter-test-go-l390|languageForPath (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languageforpath--internal-adapter-ts-adapter-test-go-l390]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindparsedfile--internal-adapter-ts-adapter-test-go-l377]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesingletslikefile--internal-adapter-ts-adapter-test-go-l329]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsetslikesources--internal-adapter-ts-adapter-test-go-l340]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadaptercomputescyclomaticcomplexity--internal-adapter-ts-adapter-test-go-l285]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterproducesdiagnosticsforparseerrors--internal-adapter-ts-adapter-test-go-l306]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadaptersupportstslikelanguages--internal-adapter-ts-adapter-test-go-l13]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadaptercomputescyclomaticcomplexity--internal-adapter-ts-adapter-test-go-l285]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterproducesdiagnosticsforparseerrors--internal-adapter-ts-adapter-test-go-l306]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadaptersupportstslikelanguages--internal-adapter-ts-adapter-test-go-l13]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
