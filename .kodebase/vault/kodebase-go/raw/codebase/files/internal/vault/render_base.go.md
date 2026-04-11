---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.6667
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 19
smells:
source_kind: "codebase-file"
source_path: "internal/vault/render_base.go"
stage: "raw"
symbol_count: 13
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/render_base.go"
type: "source"
---

# Codebase File: internal/vault/render_base.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 2
- Instability: 0.6667
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-base-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderbasefiles--internal-vault-render-base-go-l11|RenderBaseFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197|RenderBaseDefinition (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createbasefile--internal-vault-render-base-go-l202|createBaseFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/exprfilter--internal-vault-render-base-go-l209|exprFilter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/andfilter--internal-vault-render-base-go-l213|andFilter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217|baseDefinitionValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/baseviewsvalue--internal-vault-render-base-go-l241|baseViewsValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268|baseFilterValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295|renderYAMLValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderyamlscalar--internal-vault-render-base-go-l364|renderYAMLScalar (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/stringmaptoany--internal-vault-render-base-go-l379|stringMapToAny (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/stringslicetoany--internal-vault-render-base-go-l387|stringSliceToAny (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/andfilter--internal-vault-render-base-go-l213]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basedefinitionvalue--internal-vault-render-base-go-l217]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefiltervalue--internal-vault-render-base-go-l268]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseviewsvalue--internal-vault-render-base-go-l241]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createbasefile--internal-vault-render-base-go-l202]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/exprfilter--internal-vault-render-base-go-l209]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasefiles--internal-vault-render-base-go-l11]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderyamlscalar--internal-vault-render-base-go-l364]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stringmaptoany--internal-vault-render-base-go-l379]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stringslicetoany--internal-vault-render-base-go-l387]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-base-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasedefinition--internal-vault-render-base-go-l197]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasefiles--internal-vault-render-base-go-l11]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
