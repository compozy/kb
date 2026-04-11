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
outgoing_relation_count: 17
smells:
source_kind: "codebase-file"
source_path: "internal/cli/index.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/index.go"
type: "source"
---

# Codebase File: internal/cli/index.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-index-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/indexcommandclient--internal-cli-index-go-l16|indexCommandClient (interface)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/indexcommandoptions--internal-cli-index-go-l21|indexCommandOptions (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newindexcommand--internal-cli-index-go-l37|newIndexCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/runindexcommand--internal-cli-index-go-l63|runIndexCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/indexresultpayload--internal-cli-index-go-l132|indexResultPayload (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/indexstatuspayload--internal-cli-index-go-l144|indexStatusPayload (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/chooseindexoperation--internal-cli-index-go-l151|chooseIndexOperation (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/findcollectionstatus--internal-cli-index-go-l161|findCollectionStatus (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/chooseindexoperation--internal-cli-index-go-l151]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-index-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findcollectionstatus--internal-cli-index-go-l161]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexcommandclient--internal-cli-index-go-l16]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexcommandoptions--internal-cli-index-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexresultpayload--internal-cli-index-go-l132]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexstatuspayload--internal-cli-index-go-l144]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newindexcommand--internal-cli-index-go-l37]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runindexcommand--internal-cli-index-go-l63]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/qmd`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `strings`

## Backlinks
None
