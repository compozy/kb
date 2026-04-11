---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 12
smells:
source_kind: "codebase-file"
source_path: "internal/graph/normalize.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/graph/normalize.go"
type: "source"
---

# Codebase File: internal/graph/normalize.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-go-l1|graph (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/relationkey--internal-graph-normalize-go-l9|relationKey (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17|NormalizeGraph (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/isrenderablefile--internal-graph-normalize-go-l79|isRenderableFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/attachsymbolids--internal-graph-normalize-go-l87|attachSymbolIDs (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/comparerelationedges--internal-graph-normalize-go-l110|compareRelationEdges (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/comparediagnostics--internal-graph-normalize-go-l123|compareDiagnostics (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/comparestrings--internal-graph-normalize-go-l134|compareStrings (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/uniquebykey--internal-graph-normalize-go-l145|uniqueByKey (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/attachsymbolids--internal-graph-normalize-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/comparediagnostics--internal-graph-normalize-go-l123]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/comparerelationedges--internal-graph-normalize-go-l110]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/comparestrings--internal-graph-normalize-go-l134]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graph--internal-graph-normalize-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isrenderablefile--internal-graph-normalize-go-l79]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationkey--internal-graph-normalize-go-l9]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquebykey--internal-graph-normalize-go-l145]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `sort`

## Backlinks
None
