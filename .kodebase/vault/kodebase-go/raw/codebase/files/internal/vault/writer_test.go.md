---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.3333
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 27
smells:
source_kind: "codebase-file"
source_path: "internal/vault/writer_test.go"
stage: "raw"
symbol_count: 13
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/writer_test.go"
type: "source"
---

# Codebase File: internal/vault/writer_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 1
- Instability: 0.3333
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-writer-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16|TestWriteVaultCreatesTopicSkeletonAndManagedFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatesclaudemanifestandappendonlylog--internal-vault-writer-test-go-l102|TestWriteVaultCreatesClaudeManifestAndAppendOnlyLog (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testwritevaultreportsprogressforpersistedfiles--internal-vault-writer-test-go-l157|TestWriteVaultReportsProgressForPersistedFiles (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187|TestWriteVaultRemovesStaleManagedWikiConceptsOnly (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testwritevaultrejectsinvalidrendereddocument--internal-vault-writer-test-go-l235|TestWriteVaultRejectsInvalidRenderedDocument (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260|testWriteVaultInputs (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testwritabletopicfixture--internal-vault-writer-test-go-l272|testWritableTopicFixture (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/countkinds--internal-vault-writer-test-go-l293|countKinds (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/filteroutdocument--internal-vault-writer-test-go-l310|filterOutDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/readfile--internal-vault-writer-test-go-l321|readFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/assertdirexists--internal-vault-writer-test-go-l332|assertDirExists (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/assertfileexists--internal-vault-writer-test-go-l344|assertFileExists (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertdirexists--internal-vault-writer-test-go-l332]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertfileexists--internal-vault-writer-test-go-l344]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/countkinds--internal-vault-writer-test-go-l293]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filteroutdocument--internal-vault-writer-test-go-l310]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readfile--internal-vault-writer-test-go-l321]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritabletopicfixture--internal-vault-writer-test-go-l272]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultcreatesclaudemanifestandappendonlylog--internal-vault-writer-test-go-l102]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultrejectsinvalidrendereddocument--internal-vault-writer-test-go-l235]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultreportsprogressforpersistedfiles--internal-vault-writer-test-go-l157]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-writer-test-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultcreatesclaudemanifestandappendonlylog--internal-vault-writer-test-go-l102]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultrejectsinvalidrendereddocument--internal-vault-writer-test-go-l235]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultreportsprogressforpersistedfiles--internal-vault-writer-test-go-l157]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/metrics`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `gopkg.in/yaml.v3`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
