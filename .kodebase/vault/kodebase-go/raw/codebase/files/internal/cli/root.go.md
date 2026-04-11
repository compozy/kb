---
afferent_coupling: 6
domain: "kodebase-go"
efferent_coupling: 5
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.4545
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 6
smells:
source_kind: "codebase-file"
source_path: "internal/cli/root.go"
stage: "raw"
symbol_count: 3
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/root.go"
type: "source"
---

# Codebase File: internal/cli/root.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 6
- Efferent coupling: 5
- Instability: 0.4545
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-root-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/executecontext--internal-cli-root-go-l10|ExecuteContext (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14|newRootCommand (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-root-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/executecontext--internal-cli-root-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/executecontext--internal-cli-root-go-l10]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `github.com/spf13/cobra`

## Backlinks
None
