---
blast_radius: 23
centrality: 1
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 30
exported: false
external_reference_count: 22
has_smells: true
incoming_relation_count: 24
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 5
smells:
  - "bottleneck"
  - "feature-envy"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/root.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "newRootCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newRootCommand"
type: "source"
---

# Codebase Symbol: newRootCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/root.go|internal/cli/root.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 23
- External references: 22
- Centrality: 1
- LOC: 17
- Dead export: false
- Smells: `bottleneck`, `feature-envy`, `high-blast-radius`

## Signature
```text
func newRootCommand() *cobra.Command {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newgeneratecommand--internal-cli-generate-go-l15]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newindexcommand--internal-cli-index-go-l37]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcommand--internal-cli-inspect-go-l41]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newsearchcommand--internal-cli-search-go-l41]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newversioncommand--internal-cli-version-go-l11]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandpassesflagsandprintsjsonsummary--internal-cli-generate-test-go-l15]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidlogformat--internal-cli-generate-test-go-l182]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidprogressmode--internal-cli-generate-test-go-l167]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritesjsoneventswhenrequested--internal-cli-generate-test-go-l111]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritestexteventstostderrbydefault--internal-cli-generate-test-go-l80]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexcommandhandlesqmdunavailable--internal-cli-index-test-go-l147]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexcommandhelpshowsflags--internal-cli-index-test-go-l190]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexcommandresolvestopicpathbeforecallingqmd--internal-cli-index-test-go-l15]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexcommandupdatesexistingcollection--internal-cli-index-test-go-l100]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectparenthelplistsallsubcommands--internal-cli-inspect-helpers-test-go-l260]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testversioncommandprintsbuildversion--internal-cli-inspect-helpers-test-go-l281]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectcommandjsonformatproducesvalidjson--internal-cli-inspect-test-go-l550]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectcommandreturnsdescriptiveerrorformissingvault--internal-cli-inspect-test-go-l594]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectcommandtsvformatproducesheaderandrows--internal-cli-inspect-test-go-l572]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/executecontext--internal-cli-root-go-l10]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommandreturnsresultsagainstindexedvault--internal-cli-search-index-integration-test-go-l21]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommanddefaultstohybridmode--internal-cli-search-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommanddisplayspathscoreandpreview--internal-cli-search-test-go-l158]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommandhandlesqmdunavailable--internal-cli-search-test-go-l242]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommandhelpshowsflags--internal-cli-search-test-go-l281]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommandpasseslimitflag--internal-cli-search-test-go-l204]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommanduseslexicalmode--internal-cli-search-test-go-l82]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcommandusesvectormode--internal-cli-search-test-go-l120]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/root.go|internal/cli/root.go]] via `contains` (syntactic)
