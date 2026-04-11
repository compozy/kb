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
outgoing_relation_count: 6
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/reader_integration_test.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/reader_integration_test.go"
type: "source"
---

# Codebase File: internal/vault/reader_integration_test.go

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
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-reader-integration-test-go-l3|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotroundtripswriteroutput--internal-vault-reader-integration-test-go-l12|TestReadVaultSnapshotRoundTripsWriterOutput (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotroundtripswriteroutput--internal-vault-reader-integration-test-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-reader-integration-test-go-l3]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotroundtripswriteroutput--internal-vault-reader-integration-test-go-l12]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `testing`

## Backlinks
None
