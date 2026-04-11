---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 19
smells:
source_kind: "codebase-file"
source_path: "internal/cli/search.go"
stage: "raw"
symbol_count: 10
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/search.go"
type: "source"
---

# Codebase File: internal/cli/search.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-search-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/searchcommandclient--internal-cli-search-go-l17|searchCommandClient (interface)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/searchcommandoptions--internal-cli-search-go-l21|searchCommandOptions (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newsearchcommand--internal-cli-search-go-l41|newSearchCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/runsearchcommand--internal-cli-search-go-l71|runSearchCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvesearchmode--internal-cli-search-go-l134|resolveSearchMode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvesearchcollection--internal-cli-search-go-l149|resolveSearchCollection (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsesearchoutputformat--internal-cli-search-go-l171|parseSearchOutputFormat (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/searchresultstorows--internal-cli-search-go-l184|searchResultsToRows (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/wrapqmdcommanderror--internal-cli-search-go-l196|wrapQMDCommandError (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-search-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newsearchcommand--internal-cli-search-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesearchoutputformat--internal-cli-search-go-l171]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvesearchcollection--internal-cli-search-go-l149]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvesearchmode--internal-cli-search-go-l134]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runsearchcommand--internal-cli-search-go-l71]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchcommandclient--internal-cli-search-go-l17]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchcommandoptions--internal-cli-search-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchresultstorows--internal-cli-search-go-l184]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/wrapqmdcommanderror--internal-cli-search-go-l196]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/output`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/qmd`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `strings`

## Backlinks
None
