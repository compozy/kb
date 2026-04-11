# Workflow Memory

Keep only durable, cross-task context here. Do not duplicate facts that are obvious from the repository, PRD documents, or git history.

## Current State

## Shared Decisions

- For KB pivot work, follow the `Converter.Convert(ctx, input ConvertInput)` signature from `task_02.md` and `_techspec.md`; `adrs/adr-003.md` still shows an older signature and should not drive new implementations.
- Dot-prefixed Karpathy KB assets under `.agents/skills/karpathy-kb/assets/` are embedded through the repo-root package in `karpathy_assets.go` using `//go:embed all:.agents/...`; downstream `internal/...` packages should import that FS instead of trying to embed parent-directory assets directly.
- `internal/ingest.Ingest` returns `IngestResult.FilePath` as a vault-root-relative path that includes the topic slug (for example `topic/raw/articles/doc.md`); downstream CLI work can print it directly without recomputing the topic prefix.
- `internal/lint` resolves wikilinks against topic-relative paths, topic-prefixed paths (for example `topic/wiki/...`), file stems, and frontmatter `title` values so both generated KB path links and hand-authored title links resolve inside a topic.
- CLI commands now inherit a single root `--vault` flag from `internal/cli/root.go`; future CLI work should read that value via the shared helpers in `internal/cli/vault_flag.go` instead of defining duplicate per-command `--vault` flags.

## Shared Learnings

- `github.com/JohannesKaufmann/html-to-markdown/v2` top-level helpers only wire the base and commonmark plugins; KB HTML conversion must build a custom converter with the table plugin enabled or HTML tables will not render as Markdown tables.
- `github.com/pdfcpu/pdfcpu` exposes PDF metadata and decoded per-page content streams, but not a ready-made plain-text extraction API; KB PDF conversion needs a small parser on top of `pdfcpu.ExtractPageContent`.
- `pdfcpu` lazily initializes shared config-directory globals in `NewDefaultConfiguration`, which trips the race detector under concurrent first use; if KB conversion only needs built-in defaults, call `pdfapi.DisableConfigDir()` before concurrent PDF work.
- If a task runs `go get` before the new import exists in source, the dependency can remain marked `// indirect`; run `go mod tidy` before closeout so `go.mod` reflects the real direct dependency set.
- `github.com/kkdai/youtube/v2@v2.10.6` upgrades the repo `go` directive to `1.26`; use `v2.10.5` on this branch if you need transcript APIs without moving the module past `go 1.24.0`.
- `gosseract` OCR builds need the native Tesseract and Leptonica development headers at compile time; without them the build fails early on `leptonica/allheaders.h` before any OCR-tagged tests can run.
- Firecrawl `/v2/scrape` returns markdown in `data.markdown`; callers should read page title from `data.metadata.title` and prefer `data.metadata.sourceURL`, then `data.metadata.url`, as the canonical source URL.
- `internal/qmd` tests that fake the CLI with generated shell scripts can intermittently hit `ETXTBSY` under parallel load if they exec the transient script path directly; use a fake client `commandContext` that runs `/bin/sh <script> ...` instead of direct script execution.

## Open Risks

- `adrs/adr-003.md` is stale relative to the task spec and tech spec, so future converter-related tasks can drift unless they check the newer sources first.

## Handoffs
