---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 36
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/reader.go"
stage: "raw"
symbol_count: 20
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/reader.go"
type: "source"
---

# Codebase File: internal/vault/reader.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-reader-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/vaultrelation--internal-vault-reader-go-l17|VaultRelation (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/vaultdocument--internal-vault-reader-go-l24|VaultDocument (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/vaultsnapshot--internal-vault-reader-go-l33|VaultSnapshot (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/readvaultoptions--internal-vault-reader-go-l43|ReadVaultOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/vaultdocumentbucket--internal-vault-reader-go-l47|vaultDocumentBucket (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62|ReadVaultSnapshot (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/extractsection--internal-vault-reader-go-l112|ExtractSection (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/findsymbolsbyname--internal-vault-reader-go-l146|FindSymbolsByName (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createemptysnapshot--internal-vault-reader-go-l181|createEmptySnapshot (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collectmarkdownfiles--internal-vault-reader-go-l192|collectMarkdownFiles (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223|parseVaultDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/classifydocument--internal-vault-reader-go-l251|classifyDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parserelations--internal-vault-reader-go-l264|parseRelations (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsebacklinks--internal-vault-reader-go-l286|parseBacklinks (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortvaultdocuments--internal-vault-reader-go-l308|sortVaultDocuments (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizefrontmattermap--internal-vault-reader-go-l314|normalizeFrontmatterMap (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322|normalizeFrontmatterValue (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/frontmatterstring--internal-vault-reader-go-l360|frontmatterString (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/frontmatterint--internal-vault-reader-go-l376|frontmatterInt (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/classifydocument--internal-vault-reader-go-l251]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectmarkdownfiles--internal-vault-reader-go-l192]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createemptysnapshot--internal-vault-reader-go-l181]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractsection--internal-vault-reader-go-l112]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsymbolsbyname--internal-vault-reader-go-l146]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/frontmatterint--internal-vault-reader-go-l376]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/frontmatterstring--internal-vault-reader-go-l360]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattermap--internal-vault-reader-go-l314]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizefrontmattervalue--internal-vault-reader-go-l322]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsebacklinks--internal-vault-reader-go-l286]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parserelations--internal-vault-reader-go-l264]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readvaultoptions--internal-vault-reader-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortvaultdocuments--internal-vault-reader-go-l308]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-reader-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultdocument--internal-vault-reader-go-l24]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultdocumentbucket--internal-vault-reader-go-l47]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultrelation--internal-vault-reader-go-l17]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultsnapshot--internal-vault-reader-go-l33]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractsection--internal-vault-reader-go-l112]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsymbolsbyname--internal-vault-reader-go-l146]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readvaultoptions--internal-vault-reader-go-l43]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultdocument--internal-vault-reader-go-l24]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultrelation--internal-vault-reader-go-l17]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vaultsnapshot--internal-vault-reader-go-l33]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `gopkg.in/yaml.v3`
- `imports` (syntactic) -> `io/fs`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
