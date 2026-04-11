---
afferent_coupling: 12
domain: "kodebase-go"
efferent_coupling: 10
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.4545
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 32
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/cli/inspect.go"
stage: "raw"
symbol_count: 23
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect.go"
type: "source"
---

# Codebase File: internal/cli/inspect.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 12
- Efferent coupling: 10
- Instability: 0.4545
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectsharedoptions--internal-cli-inspect-go-l17|inspectSharedOptions (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectoutput--internal-cli-inspect-go-l23|inspectOutput (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectcontext--internal-cli-inspect-go-l28|inspectContext (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41|newInspectCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/bindinspectsharedflags--internal-cli-inspect-go-l70|bindInspectSharedFlags (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/runinspectcommand--internal-cli-inspect-go-l76|runInspectCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolveinspectcontext--internal-cli-inspect-go-l103|resolveInspectContext (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parseinspectoutputformat--internal-cli-inspect-go-l134|parseInspectOutputFormat (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isfunctionlikedocument--internal-cli-inspect-go-l147|isFunctionLikeDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156|inspectFrontmatterString (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170|inspectFrontmatterStringArray (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193|inspectFrontmatterBool (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213|inspectFrontmatterInt (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254|inspectFrontmatterFloat (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectdetailentry--internal-cli-inspect-go-l295|inspectDetailEntry (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createinspectdetailoutput--internal-cli-inspect-go-l300|createInspectDetailOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createinspectrelationrows--internal-cli-inspect-go-l315|createInspectRelationRows (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/findinspectfilebysourcepath--internal-cli-inspect-go-l339|findInspectFileBySourcePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/findsingleinspectsymbolmatch--internal-cli-inspect-go-l350|findSingleInspectSymbolMatch (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolveinspectentity--internal-cli-inspect-go-l377|resolveInspectEntity (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectsectiontext--internal-cli-inspect-go-l397|inspectSectionText (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectsymbolsforfile--internal-cli-inspect-go-l411|inspectSymbolsForFile (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/bindinspectsharedflags--internal-cli-inspect-go-l70]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectdetailoutput--internal-cli-inspect-go-l300]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectrelationrows--internal-cli-inspect-go-l315]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findinspectfilebysourcepath--internal-cli-inspect-go-l339]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsingleinspectsymbolmatch--internal-cli-inspect-go-l350]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectcontext--internal-cli-inspect-go-l28]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectdetailentry--internal-cli-inspect-go-l295]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectoutput--internal-cli-inspect-go-l23]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectsectiontext--internal-cli-inspect-go-l397]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectsharedoptions--internal-cli-inspect-go-l17]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectsymbolsforfile--internal-cli-inspect-go-l411]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isfunctionlikedocument--internal-cli-inspect-go-l147]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseinspectoutputformat--internal-cli-inspect-go-l134]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveinspectcontext--internal-cli-inspect-go-l103]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveinspectentity--internal-cli-inspect-go-l377]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runinspectcommand--internal-cli-inspect-go-l76]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/spf13/pflag`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/output`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
