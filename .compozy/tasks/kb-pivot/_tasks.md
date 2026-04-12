# KB Pivot — Task List

**Reference implementation:** `.agents/skills/karpathy-kb/` (Karpathy KB skill — behavioral source for the CLI pivot)

## Tasks

| # | Title | Status | Complexity | Dependencies |
|---|-------|--------|------------|--------------|
| 01 | Implement frontmatter package | completed | medium | — |
| 02 | Extend domain models with KB types | completed | medium | — |
| 03 | Extend config with Firecrawl and OpenRouter | completed | low | — |
| 04 | Implement topic scaffolding | completed | medium | task_01, task_02 |
| 05 | Implement converter registry and simple converters | completed | medium | task_02 |
| 06 | Implement HTML-to-Markdown converter | completed | medium | task_05 |
| 07 | Implement PDF converter | completed | medium | task_05 |
| 08 | Implement Office format converters (DOCX, PPTX, XLSX) | completed | high | task_05 |
| 09 | Implement EPUB and image OCR converters | completed | medium | task_05, task_06 |
| 10 | Implement Firecrawl REST client | completed | medium | task_03 |
| 11 | Implement YouTube transcript extractor | completed | high | task_03 |
| 12 | Implement ingest orchestrator | completed | high | task_01, task_02, task_04, task_05 |
| 13 | Implement lint engine | completed | high | task_01, task_02 |
| 14 | Adapt codebase generate pipeline | completed | high | task_04, task_12 |
| 15 | Rename binary and rewrite CLI root and topic commands | completed | medium | task_04 |
| 16 | Implement CLI ingest commands | completed | high | task_10, task_11, task_12, task_15 |
| 17 | Implement CLI lint command and adapt existing commands | completed | medium | task_13, task_14, task_15 |
| 18 | Integration tests and Makefile update | completed | high | task_16, task_17 |
