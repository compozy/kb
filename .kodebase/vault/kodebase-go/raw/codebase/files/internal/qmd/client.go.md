---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 70
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/qmd/client.go"
stage: "raw"
symbol_count: 42
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/qmd/client.go"
type: "source"
---

# Codebase File: internal/qmd/client.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
Package qmd provides a shell-backed client for the QMD CLI.

## Symbols
- [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-go-l2|qmd (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/lookpathfunc--internal-qmd-client-go-l42|lookPathFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/commandcontextfunc--internal-qmd-client-go-l43|commandContextFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/searchmode--internal-qmd-client-go-l46|SearchMode (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/indexoperation--internal-qmd-client-go-l58|IndexOperation (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/searchoptions--internal-qmd-client-go-l68|SearchOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/searchresult--internal-qmd-client-go-l79|SearchResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/indexoptions--internal-qmd-client-go-l88|IndexOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/updateresult--internal-qmd-client-go-l98|UpdateResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/embedresult--internal-qmd-client-go-l108|EmbedResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/collectioninfo--internal-qmd-client-go-l116|CollectionInfo (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/indexstatus--internal-qmd-client-go-l125|IndexStatus (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/indexresult--internal-qmd-client-go-l134|IndexResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/clientoption--internal-qmd-client-go-l142|ClientOption (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/qmdclient--internal-qmd-client-go-l145|QMDClient (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154|NewClient (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169|WithBinaryPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withindexname--internal-qmd-client-go-l176|WithIndexName (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/status--internal-qmd-client-go-l183|Status (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/search--internal-qmd-client-go-l198|Search (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232|Index (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/searchcommand--internal-qmd-client-go-l307|searchCommand (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/indexcommand--internal-qmd-client-go-l336|indexCommand (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/embedcommand--internal-qmd-client-go-l363|embedCommand (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/statuscommand--internal-qmd-client-go-l375|statusCommand (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/baseargs--internal-qmd-client-go-l382|baseArgs (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvebinary--internal-qmd-client-go-l391|resolveBinary (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/run--internal-qmd-client-go-l408|run (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/commandspec--internal-qmd-client-go-l442|commandSpec (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/searchresultpayload--internal-qmd-client-go-l447|searchResultPayload (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalize--internal-qmd-client-go-l459|normalize (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvesnippet--internal-qmd-client-go-l469|resolveSnippet (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizesearchmode--internal-qmd-client-go-l476|normalizeSearchMode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizeindexoperation--internal-qmd-client-go-l489|normalizeIndexOperation (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parseupdateresult--internal-qmd-client-go-l502|parseUpdateResult (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555|parseEmbedResult (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598|parseIndexStatus (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsesingleinteger--internal-qmd-client-go-l688|parseSingleInteger (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsehumandurationmilliseconds--internal-qmd-client-go-l701|parseHumanDurationMilliseconds (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736|cleanOutput (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/cleandiagnostics--internal-qmd-client-go-l740|cleanDiagnostics (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/firstnonempty--internal-qmd-client-go-l746|firstNonEmpty (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseargs--internal-qmd-client-go-l382]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleandiagnostics--internal-qmd-client-go-l740]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/clientoption--internal-qmd-client-go-l142]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectioninfo--internal-qmd-client-go-l116]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/commandcontextfunc--internal-qmd-client-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/commandspec--internal-qmd-client-go-l442]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/embedcommand--internal-qmd-client-go-l363]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/embedresult--internal-qmd-client-go-l108]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/firstnonempty--internal-qmd-client-go-l746]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexcommand--internal-qmd-client-go-l336]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexoperation--internal-qmd-client-go-l58]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexoptions--internal-qmd-client-go-l88]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexresult--internal-qmd-client-go-l134]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexstatus--internal-qmd-client-go-l125]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/lookpathfunc--internal-qmd-client-go-l42]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalize--internal-qmd-client-go-l459]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizeindexoperation--internal-qmd-client-go-l489]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizesearchmode--internal-qmd-client-go-l476]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsehumandurationmilliseconds--internal-qmd-client-go-l701]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesingleinteger--internal-qmd-client-go-l688]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseupdateresult--internal-qmd-client-go-l502]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/qmd--internal-qmd-client-go-l2]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/qmdclient--internal-qmd-client-go-l145]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvebinary--internal-qmd-client-go-l391]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvesnippet--internal-qmd-client-go-l469]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/run--internal-qmd-client-go-l408]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/search--internal-qmd-client-go-l198]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchcommand--internal-qmd-client-go-l307]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchmode--internal-qmd-client-go-l46]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchoptions--internal-qmd-client-go-l68]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchresult--internal-qmd-client-go-l79]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchresultpayload--internal-qmd-client-go-l447]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/status--internal-qmd-client-go-l183]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/statuscommand--internal-qmd-client-go-l375]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/updateresult--internal-qmd-client-go-l98]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withindexname--internal-qmd-client-go-l176]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/clientoption--internal-qmd-client-go-l142]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectioninfo--internal-qmd-client-go-l116]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/embedresult--internal-qmd-client-go-l108]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexoperation--internal-qmd-client-go-l58]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexoptions--internal-qmd-client-go-l88]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexresult--internal-qmd-client-go-l134]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/indexstatus--internal-qmd-client-go-l125]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/qmdclient--internal-qmd-client-go-l145]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/search--internal-qmd-client-go-l198]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchmode--internal-qmd-client-go-l46]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchoptions--internal-qmd-client-go-l68]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/searchresult--internal-qmd-client-go-l79]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/status--internal-qmd-client-go-l183]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/updateresult--internal-qmd-client-go-l98]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withindexname--internal-qmd-client-go-l176]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `os/exec`
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
