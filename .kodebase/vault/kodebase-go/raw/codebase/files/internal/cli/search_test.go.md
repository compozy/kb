---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 25
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/search_test.go"
stage: "raw"
symbol_count: 10
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/search_test.go"
type: "source"
---

# Codebase File: internal/cli/search_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-search-test-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testsearchcommanddefaultstohybridmode--internal-cli-search-test-go-l14|TestSearchCommandDefaultsToHybridMode (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommanduseslexicalmode--internal-cli-search-test-go-l82|TestSearchCommandUsesLexicalMode (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommandusesvectormode--internal-cli-search-test-go-l120|TestSearchCommandUsesVectorMode (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommanddisplayspathscoreandpreview--internal-cli-search-test-go-l158|TestSearchCommandDisplaysPathScoreAndPreview (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommandpasseslimitflag--internal-cli-search-test-go-l204|TestSearchCommandPassesLimitFlag (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommandhandlesqmdunavailable--internal-cli-search-test-go-l242|TestSearchCommandHandlesQMDUnavailable (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testsearchcommandhelpshowsflags--internal-cli-search-test-go-l281|TestSearchCommandHelpShowsFlags (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/fakesearchclient--internal-cli-search-test-go-l299|fakeSearchClient (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/search--internal-cli-search-test-go-l303|Search (method)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-search-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fakesearchclient--internal-cli-search-test-go-l299]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/search--internal-cli-search-test-go-l303]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanddefaultstohybridmode--internal-cli-search-test-go-l14]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanddisplayspathscoreandpreview--internal-cli-search-test-go-l158]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandhandlesqmdunavailable--internal-cli-search-test-go-l242]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandhelpshowsflags--internal-cli-search-test-go-l281]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandpasseslimitflag--internal-cli-search-test-go-l204]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanduseslexicalmode--internal-cli-search-test-go-l82]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandusesvectormode--internal-cli-search-test-go-l120]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/search--internal-cli-search-test-go-l303]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanddefaultstohybridmode--internal-cli-search-test-go-l14]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanddisplayspathscoreandpreview--internal-cli-search-test-go-l158]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandhandlesqmdunavailable--internal-cli-search-test-go-l242]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandhelpshowsflags--internal-cli-search-test-go-l281]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandpasseslimitflag--internal-cli-search-test-go-l204]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommanduseslexicalmode--internal-cli-search-test-go-l82]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsearchcommandusesvectormode--internal-cli-search-test-go-l120]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/qmd`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
