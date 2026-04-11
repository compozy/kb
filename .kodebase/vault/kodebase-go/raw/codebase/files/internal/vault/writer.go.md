---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 43
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/writer.go"
stage: "raw"
symbol_count: 28
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/writer.go"
type: "source"
---

# Codebase File: internal/vault/writer.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 2
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-writer-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writevaultoptions--internal-vault-writer-go-l25|WriteVaultOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/writevaultresult--internal-vault-writer-go-l34|WriteVaultResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/writeprogress--internal-vault-writer-go-l41|WriteProgress (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/filewriterequest--internal-vault-writer-go-l47|fileWriteRequest (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53|WriteVault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/validatetopic--internal-vault-writer-go-l113|validateTopic (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensuretopicskeleton--internal-vault-writer-go-l126|ensureTopicSkeleton (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resetmanagedsubtrees--internal-vault-writer-go-l149|resetManagedSubtrees (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/buildwriterequests--internal-vault-writer-go-l166|buildWriteRequests (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/validaterendereddocument--internal-vault-writer-go-l200|validateRenderedDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/validatebasefile--internal-vault-writer-go-l242|validateBaseFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/expecteddocumentplacement--internal-vault-writer-go-l257|expectedDocumentPlacement (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/cleantopicrelativepath--internal-vault-writer-go-l270|cleanTopicRelativePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensuredirectories--internal-vault-writer-go-l285|ensureDirectories (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writefilesinbatches--internal-vault-writer-go-l306|writeFilesInBatches (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writeprogressreporter--internal-vault-writer-go-l333|writeProgressReporter (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newwriteprogressreporter--internal-vault-writer-go-l339|newWriteProgressReporter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/report--internal-vault-writer-go-l346|Report (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359|buildTopicClaude (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensureagentssymlink--internal-vault-writer-go-l442|ensureAgentsSymlink (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensuretopicgitkeeps--internal-vault-writer-go-l454|ensureTopicGitkeeps (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/ensuregitkeep--internal-vault-writer-go-l473|ensureGitkeep (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/removemanagedwikiconcepts--internal-vault-writer-go-l489|removeManagedWikiConcepts (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/hasmanagedgenerator--internal-vault-writer-go-l520|hasManagedGenerator (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/appendlog--internal-vault-writer-go-l529|appendLog (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/countwrittendocuments--internal-vault-writer-go-l581|countWrittenDocuments (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writetextfile--internal-vault-writer-go-l598|writeTextFile (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/appendlog--internal-vault-writer-go-l529]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildwriterequests--internal-vault-writer-go-l166]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleantopicrelativepath--internal-vault-writer-go-l270]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/countwrittendocuments--internal-vault-writer-go-l581]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensureagentssymlink--internal-vault-writer-go-l442]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuredirectories--internal-vault-writer-go-l285]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuregitkeep--internal-vault-writer-go-l473]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuretopicgitkeeps--internal-vault-writer-go-l454]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuretopicskeleton--internal-vault-writer-go-l126]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/expecteddocumentplacement--internal-vault-writer-go-l257]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filewriterequest--internal-vault-writer-go-l47]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasmanagedgenerator--internal-vault-writer-go-l520]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newwriteprogressreporter--internal-vault-writer-go-l339]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/removemanagedwikiconcepts--internal-vault-writer-go-l489]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/report--internal-vault-writer-go-l346]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resetmanagedsubtrees--internal-vault-writer-go-l149]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validatebasefile--internal-vault-writer-go-l242]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validaterendereddocument--internal-vault-writer-go-l200]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validatetopic--internal-vault-writer-go-l113]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-writer-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefilesinbatches--internal-vault-writer-go-l306]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writeprogress--internal-vault-writer-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writeprogressreporter--internal-vault-writer-go-l333]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetextfile--internal-vault-writer-go-l598]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevaultoptions--internal-vault-writer-go-l25]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevaultresult--internal-vault-writer-go-l34]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/report--internal-vault-writer-go-l346]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writeprogress--internal-vault-writer-go-l41]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevaultoptions--internal-vault-writer-go-l25]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevaultresult--internal-vault-writer-go-l34]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
