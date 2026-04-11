---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 31
smells:
source_kind: "codebase-file"
source_path: "magefile.go"
stage: "raw"
symbol_count: 14
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: magefile.go"
type: "source"
---

# Codebase File: magefile.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: true
- Circular dependency: false
- Smells: None

## Module Notes
go:build mage

## Symbols
- [[kodebase-go/raw/codebase/symbols/main--magefile-go-l3|main (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/deps--magefile-go-l28|Deps (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/fmt--magefile-go-l32|Fmt (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/lint--magefile-go-l44|Lint (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/test--magefile-go-l59|Test (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testintegration--magefile-go-l64|TestIntegration (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/build--magefile-go-l68|Build (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/buildgo--magefile-go-l72|buildGo (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/boundaries--magefile-go-l86|Boundaries (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/verify--magefile-go-l125|Verify (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/gofiles--magefile-go-l143|goFiles (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/buildldflags--magefile-go-l170|buildLDFlags (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/gitoutput--magefile-go-l190|gitOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rungotests--magefile-go-l200|runGoTests (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/boundaries--magefile-go-l86]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/build--magefile-go-l68]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildgo--magefile-go-l72]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildldflags--magefile-go-l170]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/deps--magefile-go-l28]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fmt--magefile-go-l32]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gitoutput--magefile-go-l190]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gofiles--magefile-go-l143]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/lint--magefile-go-l44]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/main--magefile-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rungotests--magefile-go-l200]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/test--magefile-go-l59]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testintegration--magefile-go-l64]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/verify--magefile-go-l125]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/boundaries--magefile-go-l86]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/build--magefile-go-l68]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/deps--magefile-go-l28]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fmt--magefile-go-l32]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/lint--magefile-go-l44]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/test--magefile-go-l59]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testintegration--magefile-go-l64]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/verify--magefile-go-l125]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/magefile/mage/sh`
- `imports` (syntactic) -> `io/fs`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `os/exec`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `time`

## Backlinks
None
