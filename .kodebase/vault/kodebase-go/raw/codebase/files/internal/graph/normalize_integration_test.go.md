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
source_path: "internal/graph/normalize_integration_test.go"
stage: "raw"
symbol_count: 4
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/graph/normalize_integration_test.go"
type: "source"
---

# Codebase File: internal/graph/normalize_integration_test.go

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
- [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-integration-test-go-l3|graph (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphmergesoverlappingimportsacrossparsedfiles--internal-graph-normalize-integration-test-go-l14|TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/writefixturefile--internal-graph-normalize-integration-test-go-l90|writeFixtureFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/hasrelation--internal-graph-normalize-integration-test-go-l99|hasRelation (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-integration-test-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-graph-normalize-integration-test-go-l99]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphmergesoverlappingimportsacrossparsedfiles--internal-graph-normalize-integration-test-go-l14]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefixturefile--internal-graph-normalize-integration-test-go-l90]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizegraphmergesoverlappingimportsacrossparsedfiles--internal-graph-normalize-integration-test-go-l14]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/adapter`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `testing`

## Backlinks
None
