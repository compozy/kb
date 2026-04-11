---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 24
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/output/formatter.go"
stage: "raw"
symbol_count: 16
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/output/formatter.go"
type: "source"
---

# Codebase File: internal/output/formatter.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/output--internal-output-formatter-go-l1|output (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/outputformat--internal-output-formatter-go-l17|OutputFormat (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/formatoptions--internal-output-formatter-go-l29|FormatOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/formatoutput--internal-output-formatter-go-l36|FormatOutput (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49|formatTable (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formatjson--internal-output-formatter-go-l87|formatJSON (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formattsv--internal-output-formatter-go-l141|formatTSV (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157|projectStringRows (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177|normalizeCellValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214|normalizeJSONValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/dereferencevalue--internal-output-formatter-go-l248|dereferenceValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sanitizeinlinevalue--internal-output-formatter-go-l265|sanitizeInlineValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/truncatetablecell--internal-output-formatter-go-l291|truncateTableCell (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formatstringrow--internal-output-formatter-go-l309|formatStringRow (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/padright--internal-output-formatter-go-l318|padRight (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/runecount--internal-output-formatter-go-l327|runeCount (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/dereferencevalue--internal-output-formatter-go-l248]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatjson--internal-output-formatter-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatoptions--internal-output-formatter-go-l29]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatoutput--internal-output-formatter-go-l36]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatstringrow--internal-output-formatter-go-l309]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattable--internal-output-formatter-go-l49]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattsv--internal-output-formatter-go-l141]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecellvalue--internal-output-formatter-go-l177]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizejsonvalue--internal-output-formatter-go-l214]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/output--internal-output-formatter-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/outputformat--internal-output-formatter-go-l17]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/padright--internal-output-formatter-go-l318]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/projectstringrows--internal-output-formatter-go-l157]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runecount--internal-output-formatter-go-l327]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sanitizeinlinevalue--internal-output-formatter-go-l265]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/truncatetablecell--internal-output-formatter-go-l291]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatoptions--internal-output-formatter-go-l29]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatoutput--internal-output-formatter-go-l36]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/outputformat--internal-output-formatter-go-l17]]
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `unicode/utf8`

## Backlinks
None
