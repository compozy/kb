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
outgoing_relation_count: 9
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_complexity.go"
stage: "raw"
symbol_count: 5
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_complexity.go"
type: "source"
---

# Codebase File: internal/cli/inspect_complexity.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-complexity-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/complexityrow--internal-cli-inspect-complexity-go-l11|complexityRow (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectcomplexitycommand--internal-cli-inspect-complexity-go-l21|newInspectComplexityCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43|toComplexityOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/complexityrowstomaps--internal-cli-inspect-complexity-go-l89|complexityRowsToMaps (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-complexity-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/complexityrow--internal-cli-inspect-complexity-go-l11]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/complexityrowstomaps--internal-cli-inspect-complexity-go-l89]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcomplexitycommand--internal-cli-inspect-complexity-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`

## Backlinks
None
