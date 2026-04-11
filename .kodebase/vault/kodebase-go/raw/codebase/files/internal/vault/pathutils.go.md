---
afferent_coupling: 5
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 48
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/vault/pathutils.go"
stage: "raw"
symbol_count: 24
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/pathutils.go"
type: "source"
---

# Codebase File: internal/vault/pathutils.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 5
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-pathutils-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15|ToPosixPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/ispathinside--internal-vault-pathutils-go-l37|IsPathInside (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createfileid--internal-vault-pathutils-go-l77|CreateFileID (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createexternalid--internal-vault-pathutils-go-l82|CreateExternalID (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-vault-pathutils-go-l87|SlugifySegment (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/humanizeslug--internal-vault-pathutils-go-l120|HumanizeSlug (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/derivetopicslug--internal-vault-pathutils-go-l140|DeriveTopicSlug (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/derivetopictitle--internal-vault-pathutils-go-l155|DeriveTopicTitle (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/derivetopicdomain--internal-vault-pathutils-go-l164|DeriveTopicDomain (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-vault-pathutils-go-l169|CreateSymbolID (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181|GetRawFileDocumentPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186|GetRawSymbolDocumentPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193|GetRawDirectoryIndexPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getrawlanguageindexpath--internal-vault-pathutils-go-l203|GetRawLanguageIndexPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208|GetWikiConceptPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213|GetWikiIndexPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/getbasefilepath--internal-vault-pathutils-go-l218|GetBaseFilePath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223|StripMarkdownExtension (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228|ToTopicWikiLink (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/cleancomparablepath--internal-vault-pathutils-go-l237|cleanComparablePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/haswindowsdriveprefix--internal-vault-pathutils-go-l246|hasWindowsDrivePrefix (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizedocumentpathsegment--internal-vault-pathutils-go-l250|normalizeDocumentPathSegment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/splitcomparablepath--internal-vault-pathutils-go-l254|splitComparablePath (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleancomparablepath--internal-vault-pathutils-go-l237]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-vault-pathutils-go-l82]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-vault-pathutils-go-l77]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-vault-pathutils-go-l169]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopicdomain--internal-vault-pathutils-go-l164]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopicslug--internal-vault-pathutils-go-l140]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopictitle--internal-vault-pathutils-go-l155]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getbasefilepath--internal-vault-pathutils-go-l218]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawlanguageindexpath--internal-vault-pathutils-go-l203]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/haswindowsdriveprefix--internal-vault-pathutils-go-l246]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/humanizeslug--internal-vault-pathutils-go-l120]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ispathinside--internal-vault-pathutils-go-l37]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizedocumentpathsegment--internal-vault-pathutils-go-l250]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-vault-pathutils-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/splitcomparablepath--internal-vault-pathutils-go-l254]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-pathutils-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-vault-pathutils-go-l82]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-vault-pathutils-go-l77]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-vault-pathutils-go-l169]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopicdomain--internal-vault-pathutils-go-l164]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopicslug--internal-vault-pathutils-go-l140]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/derivetopictitle--internal-vault-pathutils-go-l155]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getbasefilepath--internal-vault-pathutils-go-l218]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawlanguageindexpath--internal-vault-pathutils-go-l203]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/humanizeslug--internal-vault-pathutils-go-l120]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ispathinside--internal-vault-pathutils-go-l37]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-vault-pathutils-go-l87]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `path`
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `strings`

## Backlinks
None
