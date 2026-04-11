---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 69
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/models/models.go"
stage: "raw"
symbol_count: 35
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/models/models.go"
type: "source"
---

# Codebase File: internal/models/models.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/models--internal-models-models-go-l1|models (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/supportedlanguage--internal-models-models-go-l4|SupportedLanguage (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/supportedlanguages--internal-models-models-go-l20|SupportedLanguages (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/relationtype--internal-models-models-go-l25|RelationType (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/relationconfidence--internal-models-models-go-l43|RelationConfidence (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/diagnosticseverity--internal-models-models-go-l53|DiagnosticSeverity (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/diagnosticstage--internal-models-models-go-l63|DiagnosticStage (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/documentkind--internal-models-models-go-l79|DocumentKind (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/managedarea--internal-models-models-go-l91|ManagedArea (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/baseviewtype--internal-models-models-go-l103|BaseViewType (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/structureddiagnostic--internal-models-models-go-l115|StructuredDiagnostic (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/graphfile--internal-models-models-go-l126|GraphFile (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/symbolnode--internal-models-models-go-l136|SymbolNode (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/externalnode--internal-models-models-go-l152|ExternalNode (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/relationedge--internal-models-models-go-l160|RelationEdge (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsedfile--internal-models-models-go-l168|ParsedFile (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/graphsnapshot--internal-models-models-go-l177|GraphSnapshot (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scannedsourcefile--internal-models-models-go-l187|ScannedSourceFile (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scannedworkspace--internal-models-models-go-l194|ScannedWorkspace (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/languageadapter--internal-models-models-go-l200|LanguageAdapter (interface)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/symbolmetrics--internal-models-models-go-l206|SymbolMetrics (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/filemetrics--internal-models-models-go-l218|FileMetrics (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/directorymetrics--internal-models-models-go-l230|DirectoryMetrics (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/metricsresult--internal-models-models-go-l237|MetricsResult (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/rendereddocument--internal-models-models-go-l247|RenderedDocument (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/topicmetadata--internal-models-models-go-l256|TopicMetadata (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generateoptions--internal-models-models-go-l267|GenerateOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generationtimings--internal-models-models-go-l279|GenerationTimings (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/generationsummary--internal-models-models-go-l291|GenerationSummary (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/basefilter--internal-models-models-go-l310|BaseFilter (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/baseproperty--internal-models-models-go-l318|BaseProperty (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/basegroupby--internal-models-models-go-l323|BaseGroupBy (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/baseview--internal-models-models-go-l329|BaseView (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/basedefinition--internal-models-models-go-l339|BaseDefinition (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/basefile--internal-models-models-go-l347|BaseFile (struct)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basedefinition--internal-models-models-go-l339]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefile--internal-models-models-go-l347]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefilter--internal-models-models-go-l310]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basegroupby--internal-models-models-go-l323]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseproperty--internal-models-models-go-l318]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseview--internal-models-models-go-l329]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseviewtype--internal-models-models-go-l103]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/diagnosticseverity--internal-models-models-go-l53]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/diagnosticstage--internal-models-models-go-l63]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/directorymetrics--internal-models-models-go-l230]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/documentkind--internal-models-models-go-l79]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/externalnode--internal-models-models-go-l152]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filemetrics--internal-models-models-go-l218]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generateoptions--internal-models-models-go-l267]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generationsummary--internal-models-models-go-l291]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generationtimings--internal-models-models-go-l279]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graphfile--internal-models-models-go-l126]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graphsnapshot--internal-models-models-go-l177]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languageadapter--internal-models-models-go-l200]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/managedarea--internal-models-models-go-l91]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/metricsresult--internal-models-models-go-l237]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/models--internal-models-models-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedfile--internal-models-models-go-l168]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationconfidence--internal-models-models-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationedge--internal-models-models-go-l160]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationtype--internal-models-models-go-l25]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendereddocument--internal-models-models-go-l247]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedsourcefile--internal-models-models-go-l187]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedworkspace--internal-models-models-go-l194]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/structureddiagnostic--internal-models-models-go-l115]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguage--internal-models-models-go-l4]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguages--internal-models-models-go-l20]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolmetrics--internal-models-models-go-l206]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolnode--internal-models-models-go-l136]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/topicmetadata--internal-models-models-go-l256]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basedefinition--internal-models-models-go-l339]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefile--internal-models-models-go-l347]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basefilter--internal-models-models-go-l310]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/basegroupby--internal-models-models-go-l323]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseproperty--internal-models-models-go-l318]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseview--internal-models-models-go-l329]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/baseviewtype--internal-models-models-go-l103]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/diagnosticseverity--internal-models-models-go-l53]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/diagnosticstage--internal-models-models-go-l63]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/directorymetrics--internal-models-models-go-l230]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/documentkind--internal-models-models-go-l79]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/externalnode--internal-models-models-go-l152]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filemetrics--internal-models-models-go-l218]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generateoptions--internal-models-models-go-l267]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generationsummary--internal-models-models-go-l291]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generationtimings--internal-models-models-go-l279]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graphfile--internal-models-models-go-l126]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/graphsnapshot--internal-models-models-go-l177]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languageadapter--internal-models-models-go-l200]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/managedarea--internal-models-models-go-l91]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/metricsresult--internal-models-models-go-l237]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedfile--internal-models-models-go-l168]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationconfidence--internal-models-models-go-l43]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationedge--internal-models-models-go-l160]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/relationtype--internal-models-models-go-l25]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendereddocument--internal-models-models-go-l247]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedsourcefile--internal-models-models-go-l187]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedworkspace--internal-models-models-go-l194]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/structureddiagnostic--internal-models-models-go-l115]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguage--internal-models-models-go-l4]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguages--internal-models-models-go-l20]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolmetrics--internal-models-models-go-l206]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/symbolnode--internal-models-models-go-l136]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/topicmetadata--internal-models-models-go-l256]]

## Backlinks
None
