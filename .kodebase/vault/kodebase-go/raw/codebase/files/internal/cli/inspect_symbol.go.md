---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.3333
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 11
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_symbol.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_symbol.go"
type: "source"
---

# Codebase File: internal/cli/inspect_symbol.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 1
- Instability: 0.3333
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-symbol-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/inspectsymbolsummaryrow--internal-cli-inspect-symbol-go-l12|inspectSymbolSummaryRow (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectsymbolcommand--internal-cli-inspect-symbol-go-l21|newInspectSymbolCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tosymbollookupoutput--internal-cli-inspect-symbol-go-l41|toSymbolLookupOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tosymbolsummaryoutput--internal-cli-inspect-symbol-go-l55|toSymbolSummaryOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96|toSymbolDetailOutput (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-symbol-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectsymbolsummaryrow--internal-cli-inspect-symbol-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectsymbolcommand--internal-cli-inspect-symbol-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymbollookupoutput--internal-cli-inspect-symbol-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosymbolsummaryoutput--internal-cli-inspect-symbol-go-l55]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
