---
afferent_coupling: 3
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
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
symbol_count: 17
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/go_adapter_test.go"
type: "source"
---

# Codebase File: internal/adapter/go_adapter_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 3
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-go-adapter-test-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testgoadaptersupportsonlygo--internal-adapter-go-adapter-test-go-l13|TestGoAdapterSupportsOnlyGo (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29|TestGoAdapterParsesSimpleGoFile (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadaptersetsexportedflags--internal-adapter-go-adapter-test-go-l88|TestGoAdapterSetsExportedFlags (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractstypedeclarations--internal-adapter-go-adapter-test-go-l108|TestGoAdapterExtractsTypeDeclarations (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsmethodsignaturewithreceiver--internal-adapter-go-adapter-test-go-l136|TestGoAdapterExtractsMethodSignatureWithReceiver (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157|TestGoAdapterExtractsImportRelations (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196|TestGoAdapterExtractsCallRelationsForDirectIdentifiers (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadaptercomputescyclomaticcomplexity--internal-adapter-go-adapter-test-go-l218|TestGoAdapterComputesCyclomaticComplexity (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterproducesdiagnosticsforparseerrors--internal-adapter-go-adapter-test-go-l260|TestGoAdapterProducesDiagnosticsForParseErrors (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsdoccomments--internal-adapter-go-adapter-test-go-l287|TestGoAdapterExtractsDocComments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgoadaptermoduledocusesonlyleadingcomment--internal-adapter-go-adapter-test-go-l312|TestGoAdapterModuleDocUsesOnlyLeadingComment (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsesinglegofile--internal-adapter-go-adapter-test-go-l331|parseSingleGoFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsegosources--internal-adapter-go-adapter-test-go-l342|parseGoSources (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379|mustFindSymbol (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/mustfindexternalnode--internal-adapter-go-adapter-test-go-l392|mustFindExternalNode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/hasrelation--internal-adapter-go-adapter-test-go-l405|hasRelation (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-go-adapter-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-adapter-go-adapter-test-go-l405]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindexternalnode--internal-adapter-go-adapter-test-go-l392]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegosources--internal-adapter-go-adapter-test-go-l342]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesinglegofile--internal-adapter-go-adapter-test-go-l331]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptercomputescyclomaticcomplexity--internal-adapter-go-adapter-test-go-l218]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsdoccomments--internal-adapter-go-adapter-test-go-l287]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsmethodsignaturewithreceiver--internal-adapter-go-adapter-test-go-l136]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractstypedeclarations--internal-adapter-go-adapter-test-go-l108]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptermoduledocusesonlyleadingcomment--internal-adapter-go-adapter-test-go-l312]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterproducesdiagnosticsforparseerrors--internal-adapter-go-adapter-test-go-l260]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptersetsexportedflags--internal-adapter-go-adapter-test-go-l88]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptersupportsonlygo--internal-adapter-go-adapter-test-go-l13]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptercomputescyclomaticcomplexity--internal-adapter-go-adapter-test-go-l218]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsdoccomments--internal-adapter-go-adapter-test-go-l287]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractsmethodsignaturewithreceiver--internal-adapter-go-adapter-test-go-l136]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterextractstypedeclarations--internal-adapter-go-adapter-test-go-l108]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptermoduledocusesonlyleadingcomment--internal-adapter-go-adapter-test-go-l312]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadapterproducesdiagnosticsforparseerrors--internal-adapter-go-adapter-test-go-l260]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptersetsexportedflags--internal-adapter-go-adapter-test-go-l88]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgoadaptersupportsonlygo--internal-adapter-go-adapter-test-go-l13]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
