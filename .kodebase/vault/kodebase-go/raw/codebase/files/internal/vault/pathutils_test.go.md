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
source_path: "internal/vault/pathutils_test.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/pathutils_test.go"
type: "source"
---

# Codebase File: internal/vault/pathutils_test.go

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
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-pathutils-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testcreatefileiddeterministic--internal-vault-pathutils-test-go-l10|TestCreateFileIDDeterministic (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testcreatesymboliddeterministicandunique--internal-vault-pathutils-test-go-l25|TestCreateSymbolIDDeterministicAndUnique (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtoposixpath--internal-vault-pathutils-test-go-l52|TestToPosixPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testispathinside--internal-vault-pathutils-test-go-l79|TestIsPathInside (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testdocumentpathderivationhelpers--internal-vault-pathutils-test-go-l107|TestDocumentPathDerivationHelpers (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtopichelpers--internal-vault-pathutils-test-go-l153|TestTopicHelpers (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtopicwikilinkhelpers--internal-vault-pathutils-test-go-l181|TestTopicWikiLinkHelpers (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testpathhelpershandleemptyinputs--internal-vault-pathutils-test-go-l197|TestPathHelpersHandleEmptyInputs (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcreatefileiddeterministic--internal-vault-pathutils-test-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcreatesymboliddeterministicandunique--internal-vault-pathutils-test-go-l25]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdocumentpathderivationhelpers--internal-vault-pathutils-test-go-l107]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testispathinside--internal-vault-pathutils-test-go-l79]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testpathhelpershandleemptyinputs--internal-vault-pathutils-test-go-l197]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopichelpers--internal-vault-pathutils-test-go-l153]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopicwikilinkhelpers--internal-vault-pathutils-test-go-l181]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoposixpath--internal-vault-pathutils-test-go-l52]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-pathutils-test-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcreatefileiddeterministic--internal-vault-pathutils-test-go-l10]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testcreatesymboliddeterministicandunique--internal-vault-pathutils-test-go-l25]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdocumentpathderivationhelpers--internal-vault-pathutils-test-go-l107]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testispathinside--internal-vault-pathutils-test-go-l79]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testpathhelpershandleemptyinputs--internal-vault-pathutils-test-go-l197]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopichelpers--internal-vault-pathutils-test-go-l153]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtopicwikilinkhelpers--internal-vault-pathutils-test-go-l181]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoposixpath--internal-vault-pathutils-test-go-l52]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `testing`

## Backlinks
None
