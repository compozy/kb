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
outgoing_relation_count: 14
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_circulardeps.go"
type: "source"
---

# Codebase File: internal/cli/inspect_circulardeps.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-circulardeps-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectcirculardepscommand--internal-cli-inspect-circulardeps-go-l12|newInspectCircularDepsCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tocirculardepsoutput--internal-cli-inspect-circulardeps-go-l27|toCircularDepsOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44|circularDependencyRows (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromflags--internal-cli-inspect-circulardeps-go-l56|circularDependencyRowsFromFlags (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74|circularDependencyRowsFromFallback (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/buildinspectimportadjacency--internal-cli-inspect-circulardeps-go-l103|buildInspectImportAdjacency (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145|toCircularDependencyRow (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortcirculardependencyrows--internal-cli-inspect-circulardeps-go-l155|sortCircularDependencyRows (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildinspectimportadjacency--internal-cli-inspect-circulardeps-go-l103]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromflags--internal-cli-inspect-circulardeps-go-l56]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-circulardeps-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcirculardepscommand--internal-cli-inspect-circulardeps-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortcirculardependencyrows--internal-cli-inspect-circulardeps-go-l155]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocirculardepsoutput--internal-cli-inspect-circulardeps-go-l27]]
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/metrics`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
