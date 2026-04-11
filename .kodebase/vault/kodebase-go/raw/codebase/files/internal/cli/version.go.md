---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 5
smells:
source_kind: "codebase-file"
source_path: "internal/cli/version.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/version.go"
type: "source"
---

# Codebase File: internal/cli/version.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-version-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newversioncommand--internal-cli-version-go-l11|newVersionCommand (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-version-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newversioncommand--internal-cli-version-go-l11]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/version`

## Backlinks
None
