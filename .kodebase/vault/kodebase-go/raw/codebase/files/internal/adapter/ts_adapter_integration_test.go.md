---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 5
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/adapter/ts_adapter_integration_test.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/ts_adapter_integration_test.go"
type: "source"
---

# Codebase File: internal/adapter/ts_adapter_integration_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 2
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
go:build integration

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-integration-test-go-l3|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11|TestTSAdapterIntegrationParsesMultiFileProject (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-ts-adapter-integration-test-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `testing`

## Backlinks
None
