# KB Pivot — Task List

**Reference implementation:** `.agents/skills/karpathy-kb/` (Karpathy KB skill — behavioral source for the CLI pivot)

## Tasks

| # | Title | Status | Complexity | Dependencies |
|---|-------|--------|------------|--------------|
| 01 | Implement frontmatter package | pending | medium | — |
| 02 | Extend domain models with KB types | pending | medium | — |
| 03 | Extend config with Firecrawl and OpenRouter | pending | low | — |
| 04 | Implement topic scaffolding | pending | medium | task_01, task_02 |
| 05 | Implement converter registry and simple converters | pending | medium | task_02 |
| 06 | Implement HTML-to-Markdown converter | pending | medium | task_05 |
| 07 | Implement PDF converter | pending | medium | task_05 |
| 08 | Implement Office format converters (DOCX, PPTX, XLSX) | pending | high | task_05 |
| 09 | Implement EPUB and image OCR converters | pending | medium | task_05, task_06 |
| 10 | Implement Firecrawl REST client | pending | medium | task_03 |
| 11 | Implement YouTube transcript extractor | pending | high | task_03 |
| 12 | Implement ingest orchestrator | pending | high | task_01, task_02, task_04, task_05 |
| 13 | Implement lint engine | pending | high | task_01, task_02 |
| 14 | Adapt codebase generate pipeline | pending | high | task_04, task_12 |
| 15 | Rename binary and rewrite CLI root and topic commands | pending | medium | task_04 |
| 16 | Implement CLI ingest commands | pending | high | task_10, task_11, task_12, task_15 |
| 17 | Implement CLI lint command and adapt existing commands | pending | medium | task_13, task_14, task_15 |
| 18 | Integration tests and Makefile update | pending | high | task_16, task_17 |
