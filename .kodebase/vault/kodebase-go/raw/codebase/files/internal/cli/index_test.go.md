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
outgoing_relation_count: 22
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/index_test.go"
stage: "raw"
symbol_count: 8
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/index_test.go"
type: "source"
---

# Codebase File: internal/cli/index_test.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-index-test-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testindexcommandresolvestopicpathbeforecallingqmd--internal-cli-index-test-go-l15|TestIndexCommandResolvesTopicPathBeforeCallingQMD (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexcommandupdatesexistingcollection--internal-cli-index-test-go-l100|TestIndexCommandUpdatesExistingCollection (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexcommandhandlesqmdunavailable--internal-cli-index-test-go-l147|TestIndexCommandHandlesQMDUnavailable (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testindexcommandhelpshowsflags--internal-cli-index-test-go-l190|TestIndexCommandHelpShowsFlags (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/fakeindexclient--internal-cli-index-test-go-l208|fakeIndexClient (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/status--internal-cli-index-test-go-l213|Status (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/index--internal-cli-index-test-go-l220|Index (method)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-index-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fakeindexclient--internal-cli-index-test-go-l208]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/index--internal-cli-index-test-go-l220]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/status--internal-cli-index-test-go-l213]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandhandlesqmdunavailable--internal-cli-index-test-go-l147]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandhelpshowsflags--internal-cli-index-test-go-l190]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandresolvestopicpathbeforecallingqmd--internal-cli-index-test-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandupdatesexistingcollection--internal-cli-index-test-go-l100]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/index--internal-cli-index-test-go-l220]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/status--internal-cli-index-test-go-l213]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandhandlesqmdunavailable--internal-cli-index-test-go-l147]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandhelpshowsflags--internal-cli-index-test-go-l190]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandresolvestopicpathbeforecallingqmd--internal-cli-index-test-go-l15]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testindexcommandupdatesexistingcollection--internal-cli-index-test-go-l100]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/qmd`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
