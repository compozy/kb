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
outgoing_relation_count: 10
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/qmd/client_integration_test.go"
stage: "raw"
symbol_count: 4
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/qmd/client_integration_test.go"
type: "source"
---

# Codebase File: internal/qmd/client_integration_test.go

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
go:build integration

## Symbols
- [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-integration-test-go-l3|qmd (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13|TestQMDClientIndexesAndSearchesTempVault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/writemarkdownfile--internal-qmd-client-integration-test-go-l65|writeMarkdownFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sanitizetestidentifier--internal-qmd-client-integration-test-go-l74|sanitizeTestIdentifier (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-integration-test-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sanitizetestidentifier--internal-qmd-client-integration-test-go-l74]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writemarkdownfile--internal-qmd-client-integration-test-go-l65]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `os/exec`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
