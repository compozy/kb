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
outgoing_relation_count: 8
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_deadcode.go"
stage: "raw"
symbol_count: 5
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_deadcode.go"
type: "source"
---

# Codebase File: internal/cli/inspect_deadcode.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-deadcode-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/deadcoderow--internal-cli-inspect-deadcode-go-l10|deadCodeRow (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectdeadcodecommand--internal-cli-inspect-deadcode-go-l19|newInspectDeadCodeCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/todeadcodeoutput--internal-cli-inspect-deadcode-go-l32|toDeadCodeOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/deadcoderowstomaps--internal-cli-inspect-deadcode-go-l81|deadCodeRowsToMaps (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-deadcode-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/deadcoderow--internal-cli-inspect-deadcode-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/deadcoderowstomaps--internal-cli-inspect-deadcode-go-l81]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectdeadcodecommand--internal-cli-inspect-deadcode-go-l19]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/todeadcodeoutput--internal-cli-inspect-deadcode-go-l32]]
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`

## Backlinks
None
