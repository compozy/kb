---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: true
language: "go"
outgoing_relation_count: 30
smells:
  - "god-file"
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/generate_output.go"
stage: "raw"
symbol_count: 18
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/generate_output.go"
type: "source"
---

# Codebase File: internal/cli/generate_output.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`, `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-output-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/generateprogressmode--internal-cli-generate-output-go-l18|generateProgressMode (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/generatelogformat--internal-cli-generate-output-go-l26|generateLogFormat (type)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsegenerateprogressmode--internal-cli-generate-output-go-l33|parseGenerateProgressMode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsegeneratelogformat--internal-cli-generate-output-go-l46|parseGenerateLogFormat (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newgenerateobserver--internal-cli-generate-output-go-l57|newGenerateObserver (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isterminalwriter--internal-cli-generate-output-go-l80|isTerminalWriter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/generatejsonobserver--internal-cli-generate-output-go-l89|generateJSONObserver (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l94|ObserveGenerateEvent (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generatetextobserver--internal-cli-generate-output-go-l104|generateTextObserver (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l112|ObserveGenerateEvent (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/handlestagestarted--internal-cli-generate-output-go-l132|handleStageStarted (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/handlestageprogress--internal-cli-generate-output-go-l155|handleStageProgress (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/handlestagecompleted--internal-cli-generate-output-go-l162|handleStageCompleted (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/handlestagefailed--internal-cli-generate-output-go-l178|handleStageFailed (method)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isdeterminatestage--internal-cli-generate-output-go-l194|isDeterminateStage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/humanstagelabel--internal-cli-generate-output-go-l203|humanStageLabel (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/completedsuffix--internal-cli-generate-output-go-l207|completedSuffix (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-output-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/completedsuffix--internal-cli-generate-output-go-l207]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatejsonobserver--internal-cli-generate-output-go-l89]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatelogformat--internal-cli-generate-output-go-l26]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generateprogressmode--internal-cli-generate-output-go-l18]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatetextobserver--internal-cli-generate-output-go-l104]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handlestagecompleted--internal-cli-generate-output-go-l162]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handlestagefailed--internal-cli-generate-output-go-l178]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handlestageprogress--internal-cli-generate-output-go-l155]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handlestagestarted--internal-cli-generate-output-go-l132]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/humanstagelabel--internal-cli-generate-output-go-l203]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdeterminatestage--internal-cli-generate-output-go-l194]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isterminalwriter--internal-cli-generate-output-go-l80]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newgenerateobserver--internal-cli-generate-output-go-l57]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l112]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l94]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegeneratelogformat--internal-cli-generate-output-go-l46]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegenerateprogressmode--internal-cli-generate-output-go-l33]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l112]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-cli-generate-output-go-l94]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/schollz/progressbar/v3`
- `imports` (syntactic) -> `kgenerate (github.com/user/go-devstack/internal/generate)`
- `imports` (syntactic) -> `golang.org/x/term`
- `imports` (syntactic) -> `io`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `sync`
- `imports` (syntactic) -> `time`

## Backlinks
None
