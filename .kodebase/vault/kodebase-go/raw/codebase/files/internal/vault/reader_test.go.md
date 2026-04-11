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
outgoing_relation_count: 20
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/reader_test.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/reader_test.go"
type: "source"
---

# Codebase File: internal/vault/reader_test.go

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
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-reader-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesfrontmatterandclassifiesdocuments--internal-vault-reader-test-go-l12|TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesrelationsections--internal-vault-reader-test-go-l81|TestReadVaultSnapshotParsesRelationSections (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotskipsmalformedyamlandwarns--internal-vault-reader-test-go-l137|TestReadVaultSnapshotSkipsMalformedYAMLAndWarns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsheadingbody--internal-vault-reader-test-go-l171|TestExtractSectionReturnsHeadingBody (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsemptystringwhenmissing--internal-vault-reader-test-go-l192|TestExtractSectionReturnsEmptyStringWhenMissing (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testfindsymbolsbynameusescaseinsensitivepartialmatch--internal-vault-reader-test-go-l200|TestFindSymbolsByNameUsesCaseInsensitivePartialMatch (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createresolvedvault--internal-vault-reader-test-go-l247|createResolvedVault (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writemarkdowndocument--internal-vault-reader-test-go-l265|writeMarkdownDocument (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createresolvedvault--internal-vault-reader-test-go-l247]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsemptystringwhenmissing--internal-vault-reader-test-go-l192]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsheadingbody--internal-vault-reader-test-go-l171]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindsymbolsbynameusescaseinsensitivepartialmatch--internal-vault-reader-test-go-l200]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesfrontmatterandclassifiesdocuments--internal-vault-reader-test-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesrelationsections--internal-vault-reader-test-go-l81]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotskipsmalformedyamlandwarns--internal-vault-reader-test-go-l137]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-reader-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writemarkdowndocument--internal-vault-reader-test-go-l265]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsemptystringwhenmissing--internal-vault-reader-test-go-l192]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractsectionreturnsheadingbody--internal-vault-reader-test-go-l171]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfindsymbolsbynameusescaseinsensitivepartialmatch--internal-vault-reader-test-go-l200]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesfrontmatterandclassifiesdocuments--internal-vault-reader-test-go-l12]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotparsesrelationsections--internal-vault-reader-test-go-l81]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testreadvaultsnapshotskipsmalformedyamlandwarns--internal-vault-reader-test-go-l137]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
