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
outgoing_relation_count: 43
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/generate/generate.go"
stage: "raw"
symbol_count: 28
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/generate/generate.go"
type: "source"
---

# Codebase File: internal/generate/generate.go

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
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l1|generate (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/scanworkspacefunc--internal-generate-generate-go-l18|scanWorkspaceFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizegraphfunc--internal-generate-generate-go-l19|normalizeGraphFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/computemetricsfunc--internal-generate-generate-go-l20|computeMetricsFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderdocumentsfunc--internal-generate-generate-go-l21|renderDocumentsFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderbasefilesfunc--internal-generate-generate-go-l26|renderBaseFilesFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writevaultfunc--internal-generate-generate-go-l27|writeVaultFunc (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/progressawarelanguageadapter--internal-generate-generate-go-l29|progressAwareLanguageAdapter (interface)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/runner--internal-generate-generate-go-l37|runner (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l51|Generate (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l57|GenerateWithObserver (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/newrunner--internal-generate-generate-go-l61|newRunner (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l82|Generate (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88|GenerateWithObserver (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withdefaults--internal-generate-generate-go-l306|withDefaults (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvepaths--internal-generate-generate-go-l340|resolvePaths (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createtopicmetadata--internal-generate-generate-go-l362|createTopicMetadata (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/workspacelanguages--internal-generate-generate-go-l390|workspaceLanguages (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/selectadapters--internal-generate-generate-go-l407|selectAdapters (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/filesforadapter--internal-generate-generate-go-l423|filesForAdapter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/languagenames--internal-generate-generate-go-l436|languageNames (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/adapternames--internal-generate-generate-go-l445|adapterNames (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/elapsedmillis--internal-generate-generate-go-l454|elapsedMillis (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/emitstagestarted--internal-generate-generate-go-l462|emitStageStarted (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/emitstageprogress--internal-generate-generate-go-l471|emitStageProgress (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/emitstagecompleted--internal-generate-generate-go-l481|emitStageCompleted (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/emitstagefailed--internal-generate-generate-go-l492|emitStageFailed (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/eventfields--internal-generate-generate-go-l503|eventFields (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapternames--internal-generate-generate-go-l445]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computemetricsfunc--internal-generate-generate-go-l20]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtopicmetadata--internal-generate-generate-go-l362]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/elapsedmillis--internal-generate-generate-go-l454]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/emitstagecompleted--internal-generate-generate-go-l481]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/emitstagefailed--internal-generate-generate-go-l492]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/emitstageprogress--internal-generate-generate-go-l471]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/emitstagestarted--internal-generate-generate-go-l462]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/eventfields--internal-generate-generate-go-l503]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filesforadapter--internal-generate-generate-go-l423]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l51]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l82]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l57]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languagenames--internal-generate-generate-go-l436]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrunner--internal-generate-generate-go-l61]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraphfunc--internal-generate-generate-go-l19]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/progressawarelanguageadapter--internal-generate-generate-go-l29]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbasefilesfunc--internal-generate-generate-go-l26]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderdocumentsfunc--internal-generate-generate-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvepaths--internal-generate-generate-go-l340]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/runner--internal-generate-generate-go-l37]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspacefunc--internal-generate-generate-go-l18]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/selectadapters--internal-generate-generate-go-l407]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withdefaults--internal-generate-generate-go-l306]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/workspacelanguages--internal-generate-generate-go-l390]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writevaultfunc--internal-generate-generate-go-l27]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l51]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l82]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l57]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/adapter`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/graph`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/metrics`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/scanner`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `time`

## Backlinks
None
