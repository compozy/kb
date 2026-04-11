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
outgoing_relation_count: 6
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/scanner/scanner_integration_test.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/scanner/scanner_integration_test.go"
type: "source"
---

# Codebase File: internal/scanner/scanner_integration_test.go

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
- [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-integration-test-go-l3|scanner (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11|TestScanWorkspaceIntegrationNestedProject (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-integration-test-go-l3]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11]]
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `testing`

## Backlinks
None
