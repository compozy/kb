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
outgoing_relation_count: 7
smells:
source_kind: "codebase-file"
source_path: "internal/cli/generate.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/generate.go"
type: "source"
---

# Codebase File: internal/cli/generate.go

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
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newgeneratecommand--internal-cli-generate-go-l15|newGenerateCommand (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newgeneratecommand--internal-cli-generate-go-l15]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `kgenerate (github.com/user/go-devstack/internal/generate)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`

## Backlinks
None
