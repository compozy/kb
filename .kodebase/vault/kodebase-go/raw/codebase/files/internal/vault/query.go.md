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
outgoing_relation_count: 21
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/query.go"
stage: "raw"
symbol_count: 11
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/query.go"
type: "source"
---

# Codebase File: internal/vault/query.go

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
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-query-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvedvault--internal-vault-query-go-l14|ResolvedVault (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/vaultqueryoptions--internal-vault-query-go-l21|VaultQueryOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/discovervaultpath--internal-vault-query-go-l28|DiscoverVaultPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/resolvevaultquery--internal-vault-query-go-l56|ResolveVaultQuery (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/listavailabletopics--internal-vault-query-go-l116|ListAvailableTopics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/resolvevaultpath--internal-vault-query-go-l138|resolveVaultPath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolveabsolutepath--internal-vault-query-go-l155|resolveAbsolutePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensuredirectory--internal-vault-query-go-l172|ensureDirectory (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/listtopicdirectories--internal-vault-query-go-l179|listTopicDirectories (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isdirectorypath--internal-vault-query-go-l210|isDirectoryPath (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/discovervaultpath--internal-vault-query-go-l28]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuredirectory--internal-vault-query-go-l172]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdirectorypath--internal-vault-query-go-l210]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/listavailabletopics--internal-vault-query-go-l116]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/listtopicdirectories--internal-vault-query-go-l179]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveabsolutepath--internal-vault-query-go-l155]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvedvault--internal-vault-query-go-l14]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvevaultpath--internal-vault-query-go-l138]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvevaultquery--internal-vault-query-go-l56]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-query-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultqueryoptions--internal-vault-query-go-l21]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/discovervaultpath--internal-vault-query-go-l28]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/listavailabletopics--internal-vault-query-go-l116]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvedvault--internal-vault-query-go-l14]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvevaultquery--internal-vault-query-go-l56]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultqueryoptions--internal-vault-query-go-l21]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
