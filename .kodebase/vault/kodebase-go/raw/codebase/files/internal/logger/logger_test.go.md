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
outgoing_relation_count: 20
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/logger/logger_test.go"
stage: "raw"
symbol_count: 8
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/logger/logger_test.go"
type: "source"
---

# Codebase File: internal/logger/logger_test.go

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
- [[kodebase-go/raw/codebase/symbols/logger--internal-logger-logger-test-go-l1|logger (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testnew--internal-logger-logger-test-go-l12|TestNew (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnewwritesjson--internal-logger-logger-test-go-l53|TestNewWritesJSON (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnewuppercaselevelkey--internal-logger-logger-test-go-l76|TestNewUppercaseLevelKey (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100|TestNewWithObserver (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnewlevelfiltering--internal-logger-logger-test-go-l122|TestNewLevelFiltering (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testobserver--internal-logger-logger-test-go-l138|testObserver (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/observe--internal-logger-logger-test-go-l143|Observe (method)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logger--internal-logger-logger-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observe--internal-logger-logger-test-go-l143]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnew--internal-logger-logger-test-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewlevelfiltering--internal-logger-logger-test-go-l122]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewuppercaselevelkey--internal-logger-logger-test-go-l76]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewwritesjson--internal-logger-logger-test-go-l53]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testobserver--internal-logger-logger-test-go-l138]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observe--internal-logger-logger-test-go-l143]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnew--internal-logger-logger-test-go-l12]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewlevelfiltering--internal-logger-logger-test-go-l122]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewuppercaselevelkey--internal-logger-logger-test-go-l76]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewwritesjson--internal-logger-logger-test-go-l53]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `log/slog`
- `imports` (syntactic) -> `sync`
- `imports` (syntactic) -> `testing`

## Backlinks
None
