---
afferent_coupling: 3
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.25
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 10
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_smells.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_smells.go"
type: "source"
---

# Codebase File: internal/cli/inspect_smells.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 3
- Efferent coupling: 1
- Instability: 0.25
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-smells-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/smellrow--internal-cli-inspect-smells-go-l11|smellRow (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectsmellscommand--internal-cli-inspect-smells-go-l19|newInspectSmellsCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37|toSmellOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/includesmellrow--internal-cli-inspect-smells-go-l87|includeSmellRow (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/smellrowstomaps--internal-cli-inspect-smells-go-l104|smellRowsToMaps (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-smells-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/includesmellrow--internal-cli-inspect-smells-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectsmellscommand--internal-cli-inspect-smells-go-l19]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/smellrow--internal-cli-inspect-smells-go-l11]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/smellrowstomaps--internal-cli-inspect-smells-go-l104]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37]]
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
