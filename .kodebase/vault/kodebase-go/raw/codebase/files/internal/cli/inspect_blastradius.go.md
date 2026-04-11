---
afferent_coupling: 3
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.25
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 9
smells:
source_kind: "codebase-file"
source_path: "internal/cli/inspect_blastradius.go"
stage: "raw"
symbol_count: 5
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_blastradius.go"
type: "source"
---

# Codebase File: internal/cli/inspect_blastradius.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 3
- Efferent coupling: 1
- Instability: 0.25
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-blastradius-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/blastradiusrow--internal-cli-inspect-blastradius-go-l11|blastRadiusRow (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newinspectblastradiuscommand--internal-cli-inspect-blastradius-go-l20|newInspectBlastRadiusCommand (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47|toBlastRadiusOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/blastradiusrowstomaps--internal-cli-inspect-blastradius-go-l92|blastRadiusRowsToMaps (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/blastradiusrow--internal-cli-inspect-blastradius-go-l11]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/blastradiusrowstomaps--internal-cli-inspect-blastradius-go-l92]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-blastradius-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectblastradiuscommand--internal-cli-inspect-blastradius-go-l20]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/spf13/cobra`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `sort`

## Backlinks
None
