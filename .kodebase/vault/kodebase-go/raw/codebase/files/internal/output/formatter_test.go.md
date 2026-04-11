---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 22
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/output/formatter_test.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/output/formatter_test.go"
type: "source"
---

# Codebase File: internal/output/formatter_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/output-test--internal-output-formatter-test-go-l1|output_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testformatoutputtablealignscolumns--internal-output-formatter-test-go-l12|TestFormatOutputTableAlignsColumns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputtablehandlesvariablewidthsandtruncation--internal-output-formatter-test-go-l69|TestFormatOutputTableHandlesVariableWidthsAndTruncation (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputjsonproducesvalidprojectedobjects--internal-output-formatter-test-go-l97|TestFormatOutputJSONProducesValidProjectedObjects (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputtsvrendersheaderandsanitizescells--internal-output-formatter-test-go-l144|TestFormatOutputTSVRendersHeaderAndSanitizesCells (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputemptydata--internal-output-formatter-test-go-l161|TestFormatOutputEmptyData (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputsinglerowinallformats--internal-output-formatter-test-go-l192|TestFormatOutputSingleRowInAllFormats (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputspecialcharactersremainvalidjson--internal-output-formatter-test-go-l229|TestFormatOutputSpecialCharactersRemainValidJSON (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testformatoutputdefaultsunsupportedformatstotable--internal-output-formatter-test-go-l262|TestFormatOutputDefaultsUnsupportedFormatsToTable (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/output-test--internal-output-formatter-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputdefaultsunsupportedformatstotable--internal-output-formatter-test-go-l262]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputemptydata--internal-output-formatter-test-go-l161]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputjsonproducesvalidprojectedobjects--internal-output-formatter-test-go-l97]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputsinglerowinallformats--internal-output-formatter-test-go-l192]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputspecialcharactersremainvalidjson--internal-output-formatter-test-go-l229]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtablealignscolumns--internal-output-formatter-test-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtablehandlesvariablewidthsandtruncation--internal-output-formatter-test-go-l69]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtsvrendersheaderandsanitizescells--internal-output-formatter-test-go-l144]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputdefaultsunsupportedformatstotable--internal-output-formatter-test-go-l262]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputemptydata--internal-output-formatter-test-go-l161]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputjsonproducesvalidprojectedobjects--internal-output-formatter-test-go-l97]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputsinglerowinallformats--internal-output-formatter-test-go-l192]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputspecialcharactersremainvalidjson--internal-output-formatter-test-go-l229]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtablealignscolumns--internal-output-formatter-test-go-l12]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtablehandlesvariablewidthsandtruncation--internal-output-formatter-test-go-l69]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testformatoutputtsvrendersheaderandsanitizescells--internal-output-formatter-test-go-l144]]
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/output`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
