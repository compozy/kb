---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 12
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/inspect_integration_test.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_integration_test.go"
type: "source"
---

# Codebase File: internal/cli/inspect_integration_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
go:build integration

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-integration-test-go-l3|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testinspectcommandsagainstgeneratedfixturevault--internal-cli-inspect-integration-test-go-l18|TestInspectCommandsAgainstGeneratedFixtureVault (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-integration-test-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandsagainstgeneratedfixturevault--internal-cli-inspect-integration-test-go-l18]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandsagainstgeneratedfixturevault--internal-cli-inspect-integration-test-go-l18]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `kgenerate (github.com/user/go-devstack/internal/generate)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `io`
- `imports` (syntactic) -> `log/slog`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `testing`

## Backlinks
None
