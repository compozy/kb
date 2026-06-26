# TechSpec: KB Pivot — General-Purpose Knowledge Base CLI

## Executive Summary

This specification describes the transformation of kodebase from a codebase-only Obsidian vault generator into `kb`, a general-purpose knowledge base CLI that implements the structural backbone of the Karpathy KB pattern. The CLI handles the non-LLM phases: topic scaffolding, multi-source ingestion (web URLs, files, YouTube, codebases, bookmarks), structural linting, QMD indexing, and search. LLM-driven compilation and querying remain in the Claude skill layer.

The primary architectural decision is a **converter registry pattern** for document ingestion — each source type (PDF, HTML, DOCX, YouTube, etc.) has a dedicated Go converter behind a common interface, enabling the `ingest` command to support any format without coupling to external runtimes. The primary trade-off is accepting lower conversion quality for some edge-case formats (complex PDFs, PPTX with embedded media) in exchange for a zero-dependency single binary.

The binary is renamed from `kodebase` to `kb`. The existing codebase analysis pipeline (tree-sitter, graph, metrics) is retained as the richest ingest source type, available via `kb ingest codebase`.

## System Architecture

### Component Overview

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Layer (Cobra)                    │
│  topic  │  ingest  │  lint  │  index  │  search │ inspect│
└────┬────┴────┬─────┴───┬───┴────┬────┴────┬────┴───┬────┘
     │         │         │        │         │        │
     v         v         v        v         v        v
┌────────┐ ┌────────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐
│ topic  │ │ ingest │ │ lint │ │ qmd  │ │ qmd  │ │ vault    │
│ mgmt   │ │ engine │ │engine│ │client│ │client│ │ reader + │
└────────┘ └───┬────┘ └──┬───┘ └──────┘ └──────┘ │ metrics  │
               │         │                        └──────────┘
               v         v
        ┌─────────────┐ ┌──────────────┐
        │  Converter   │ │  Vault       │
        │  Registry    │ │  Walker      │
        │  ┌─────────┐ │ │  (frontmatter│
        │  │ PDF     │ │ │   + wikilink │
        │  │ HTML    │ │ │   parsing)   │
        │  │ DOCX    │ │ └──────────────┘
        │  │ YouTube │ │
        │  │ Codebase│ │
        │  │ ...     │ │
        │  └─────────┘ │
        └──────┬───────┘
               │
               v
        ┌──────────────┐     ┌──────────────┐
        │  Vault       │     │  Firecrawl   │
        │  Writer      │     │  REST Client │
        │  (raw/ docs) │     │  (URL ingest)│
        └──────────────┘     └──────────────┘
```

**Data flow for ingestion:**

1. CLI parses command and flags → resolves target topic
2. Ingest engine selects the appropriate converter (by subcommand or file extension)
3. Converter fetches/reads the source and produces Markdown + metadata
4. Ingest engine prepends frontmatter (type, stage, domain, source info, tags)
5. Vault writer saves to `<topic>/raw/<source-type>/<slug>.md`
6. Log appender adds entry to `<topic>/log.md`
7. (Optional) QMD client re-indexes the topic collection

**External system interactions:**

- **Firecrawl REST API** — URL scraping (`ingest url`)
- **YouTube** — Caption and audio extraction via `yt-dlp`; optional STT via configured OpenAI/OpenRouter provider (`ingest youtube`)
- **QMD CLI** — Search indexing and semantic search (via subprocess, same as today)
- **Tesseract** — OCR for images (optional, `ingest file` with image files)
- **OpenAI/OpenRouter API** — Audio transcription for YouTube videos when `--transcribe auto|stt` selects STT

## Implementation Design

### Core Interfaces

```go
// Converter transforms a source into Markdown content.
// Each supported format implements this interface.
type Converter interface {
    // Accepts returns true if this converter handles the given
    // file extension and/or MIME type.
    Accepts(ext string, mimeType string) bool

    // Convert reads from the source and produces Markdown.
    Convert(ctx context.Context, input ConvertInput) (*ConvertResult, error)
}

type ConvertInput struct {
    Reader   io.ReadSeeker // file content (nil for URL-based sources)
    FilePath string        // original file path (empty for URLs)
    URL      string        // source URL (empty for files)
    Options  map[string]any
}

type ConvertResult struct {
    Markdown string         // converted Markdown body
    Title    string         // extracted title (from doc metadata or content)
    Metadata map[string]any // extra metadata (page count, duration, etc.)
}
```

```go
// IngestResult represents a successfully ingested source.
type IngestResult struct {
    Topic      string // topic slug
    SourceType string // "article", "youtube", "codebase", etc.
    FilePath   string // path within the vault where the file was written
    Title      string // extracted or provided title
}
```

```go
// LintIssue represents a single structural problem found in the vault.
type LintIssue struct {
    Kind     LintIssueKind // dead-link, orphan, missing-source, stale, format
    Severity string        // error, warning
    FilePath string        // file containing the issue
    Message  string        // human-readable description
    Target   string        // the problematic link/reference (if applicable)
}
```

### Data Models

**Topic directory structure** (managed by `kb topic new`):

```
<topic-slug>/
├── raw/
│   ├── articles/          # Web scrapes, PDFs, docs
│   ├── bookmarks/         # Bookmark cluster exports
│   ├── codebase/          # Codebase analysis output
│   │   ├── files/         # Per-file markdown with metrics frontmatter
│   │   └── symbols/       # Per-symbol markdown
│   ├── github/            # GitHub README snapshots
│   └── youtube/           # Video transcripts
├── wiki/
│   ├── concepts/          # LLM-compiled articles (not CLI-managed)
│   └── index/             # Dashboard, Concept Index, Source Index
├── outputs/
│   ├── queries/           # Filed-back Q&A results
│   ├── briefings/         # Formatted outputs
│   ├── diagrams/          # Visual artifacts
│   └── reports/           # Lint reports
├── bases/                 # Obsidian Base definition files
├── CLAUDE.md              # Topic schema document
├── AGENTS.md              # Symlink to CLAUDE.md
└── log.md                 # Append-only audit trail
```

**Frontmatter schema** (all ingested files follow `references/frontmatter-schemas.md`):

```yaml
---
title: Source Title
type: source
stage: raw
domain: <topic-domain>
source_kind: article | github-readme | youtube-transcript | codebase-file | codebase-symbol | bookmark-cluster | document
source_url: https://...           # for URL-sourced content
source_path: /path/to/file        # for file-sourced content
scraped: 2026-04-11
tags:
  - <topic-domain>
  - raw
  - <source-kind>
---
```

For `source_kind: youtube-transcript`, `kb ingest youtube` also persists `yt-dlp` video metadata as flat optional frontmatter: `view_count`, `like_count`, `comment_count`, `upload_date`, `duration` seconds, `duration_string`, `channel`, `channel_id`, `uploader_id`, `channel_follower_count`, `categories`, `youtube_tags`, `language`, `live_status`, `was_live`, and `chapter_count`. `youtube_tags` stores YouTube keywords because `tags` remains reserved for KB taxonomy. Chapter payloads are not stored; only `chapter_count` is persisted.

**Config model** (extends existing `internal/config`):

```go
type Config struct {
    App        AppConfig
    Log        LogConfig
    Vault      VaultConfig
    Firecrawl  FirecrawlConfig
    OpenRouter OpenRouterConfig
    STT        STTConfig
    YouTube    YouTubeConfig
}

type VaultConfig struct {
    Root       string
    TopicGlobs []string
}

type FirecrawlConfig struct {
    APIKey string // or FIRECRAWL_API_KEY env
    APIURL string // default: https://api.firecrawl.dev
}

type OpenRouterConfig struct {
    APIKey   string // or OPENROUTER_API_KEY env
    APIURL   string // default: https://openrouter.ai/api
    STTModel string // default: google/gemini-2.5-flash
}

type STTConfig struct {
    Provider      string // openai | openrouter; default: openai
    APIKey        string // OPENAI_API_KEY when provider=openai
    APIURL        string // default: https://api.openai.com
    Model         string // default: gpt-4o-transcribe
    Language      string // default: auto
    Prompt        string
    AudioFormat   string // mp3, mp4, mpeg, mpga, m4a, wav, webm
    ChunkDuration string // default: 10m
    MaxChunkBytes int64
    Concurrency   int
    FFmpegPath    string
}

type YouTubeConfig struct {
    YTDLPPath     string // or YOUTUBE_YT_DLP_PATH env
    Proxy         string // or YOUTUBE_PROXY env
    CookiesFile   string // or YOUTUBE_COOKIES_FILE env
    UserAgent     string // or YOUTUBE_USER_AGENT env
    Transcription string // captions | auto | stt; default: captions
    RetryAttempts int
    RetryBackoff  string
}
```

Provider-aware env behavior keeps OpenAI and OpenRouter credentials separate: `OPENAI_API_KEY` and `OPENAI_API_URL` populate `[stt]` only when `stt.provider = "openai"`; OpenRouter credentials come from `OPENROUTER_API_KEY` and `OPENROUTER_API_URL` when `stt.provider = "openrouter"`. `STT_PROVIDER` and `STT_MODEL` can override the provider/model at runtime. YouTube policy values are `captions`, `auto`, and `stt`; `auto` uses manual captions when available and falls back to STT when only automatic captions or no captions are available.

### API Endpoints

Not applicable — this is a CLI tool. External API interactions:

| Service | Endpoint | Purpose |
|---------|----------|---------|
| Firecrawl | `POST /v2/scrape` | URL → Markdown conversion |
| OpenAI | `POST /v1/audio/transcriptions` | Default audio transcription provider |
| OpenRouter | `POST /api/v1/chat/completions` | Optional compatible audio transcription provider |

## Integration Points

### Firecrawl REST API

- **Purpose**: Convert web URLs to clean Markdown
- **Auth**: Bearer token via `FIRECRAWL_API_KEY` env or config
- **Error handling**: Graceful degradation with actionable error message if API key is missing or request fails
- **Retry**: Exponential backoff with max 3 retries for transient failures (429, 5xx)

### QMD CLI (existing)

- **Purpose**: Semantic/lexical search indexing and querying
- **Integration**: Subprocess via `exec.CommandContext` (unchanged from current implementation)
- **Changes**: Collection naming convention aligns with topic slugs

### STT Providers

- **Purpose**: Speech-to-text for YouTube videos when captions are unavailable or explicitly bypassed
- **Auth**: Bearer token via `OPENAI_API_KEY` for the default `openai` provider, or `OPENROUTER_API_KEY` when `stt.provider = "openrouter"`
- **Flow**: Download audio through `yt-dlp` → segment long files with `ffmpeg` → transcribe chunks → stitch transcript with timestamp offsets
- **Condition**: Invoked by `kb ingest youtube --transcribe stt`, or by `--transcribe auto` when manual captions are unavailable

### Tesseract (optional)

- **Purpose**: OCR for image files during `ingest file`
- **Integration**: Via `gosseract` Go bindings (requires Tesseract system installation)
- **Condition**: Only invoked for image file extensions (.png, .jpg, .jpeg, .tiff, .bmp) when Tesseract is available
- **Graceful degradation**: If Tesseract is not installed, log a warning and skip OCR (extract EXIF metadata only)

## Impact Analysis

| Component | Impact Type | Description and Risk | Required Action |
|-----------|-------------|---------------------|-----------------|
| `cmd/kodebase/` | renamed | Becomes `cmd/kb/`, binary name changes | Rename directory, update Makefile |
| `internal/cli/` | modified | Complete command restructuring: new `topic`, `ingest` parents; `generate` removed as top-level; `inspect` retained | Rewrite command files |
| `internal/generate/` | modified | Becomes the implementation behind `ingest codebase`; output writes to `raw/codebase/` instead of top-level vault | Adapt output paths and options |
| `internal/vault/writer.go` | modified | Extended to support the full topic directory skeleton and generic raw-file writing | Add topic scaffolding and raw-file write functions |
| `internal/vault/render.go` | modified | Codebase rendering writes to `raw/codebase/` subtree; wiki rendering stays for codebase topics | Adjust output paths |
| `internal/config/` | modified | Add `FirecrawlConfig` and `OpenRouterConfig` sections | Extend config struct and TOML parsing |
| `internal/convert/` | new | Converter registry and format-specific converters | Implement from scratch |
| `internal/ingest/` | new | Ingest orchestration: converter selection, frontmatter generation, vault writing, log appending | Implement from scratch |
| `internal/lint/` | new | Vault structural health checker | Implement from scratch |
| `internal/firecrawl/` | new | Firecrawl REST API client | Implement from scratch |
| `internal/youtube/` | new | YouTube caption extraction via `yt-dlp` plus OpenAI/OpenRouter STT providers | Implement from scratch |
| `internal/topic/` | new | Topic scaffolding and management | Implement from scratch |
| `internal/frontmatter/` | new | YAML frontmatter parsing and generation | Implement from scratch |
| `internal/scanner/` | unchanged | Still used by codebase ingest pipeline | No action |
| `internal/adapter/` | unchanged | Still used by codebase ingest pipeline | No action |
| `internal/graph/` | unchanged | Still used by codebase ingest pipeline | No action |
| `internal/metrics/` | unchanged | Still used by codebase ingest pipeline | No action |
| `internal/qmd/` | unchanged | Used by `index` and `search` commands | No action |
| `internal/output/` | unchanged | Used by `lint`, `inspect`, `search` | No action |
| `internal/models/` | modified | Add ingest/topic/lint domain types | Extend with new types |

## Testing Approach

### Unit Tests

- **Converter tests**: Each converter tested with fixture files (small PDF, HTML page, DOCX, XLSX, CSV, JSON). Verify Markdown output structure and metadata extraction.
- **Frontmatter tests**: Round-trip parse/generate for each frontmatter schema variant.
- **Lint tests**: Build fixture vaults with known issues (dead links, orphans, missing sources), verify detection.
- **Topic tests**: Scaffold a topic in `t.TempDir()`, verify directory structure and template content.
- **Firecrawl client tests**: HTTP test server mocking `/v2/scrape` responses.

### Integration Tests

- **Ingest end-to-end**: Ingest a fixture file → verify vault output → verify log entry → verify frontmatter correctness.
- **Codebase ingest**: Run `ingest codebase` against a small Go fixture project → verify `raw/codebase/` output matches current `generate` quality.
- **Lint end-to-end**: Build a vault with mixed healthy and broken content → run lint → verify all issues detected.
- **QMD integration**: Existing QMD integration tests continue to work with isolated `HOME`/`XDG_*`.

## Development Sequencing

### Build Order

1. **`internal/frontmatter/`** — YAML frontmatter parser and generator. Shared by all subsequent components. No dependencies.
2. **`internal/topic/`** — Topic scaffolding (create directory tree, install templates, manage `CLAUDE.md` and `log.md`). Depends on step 1 for frontmatter generation.
3. **`internal/convert/`** — Converter interface, registry, and simple converters (plain text, CSV, JSON, XML, HTML-to-Markdown). Depends on step 1 for metadata types.
4. **`internal/convert/` (complex formats)** — PDF, DOCX, PPTX, XLSX, EPUB converters using ZIP+XML parsing and library integrations. Depends on step 3 for the Converter interface.
5. **`internal/firecrawl/`** — Firecrawl REST API client. No internal dependencies (uses net/http).
6. **`internal/youtube/`** — YouTube caption extractor + STT providers. No internal dependencies.
7. **`internal/ingest/`** — Ingest orchestrator: selects converter, prepends frontmatter, writes to vault, appends log. Depends on steps 1, 2, 3, 4, 5, 6.
8. **`internal/lint/`** — Vault lint engine: walks vault, parses frontmatter, extracts wikilinks, runs checks. Depends on step 1 for frontmatter parsing.
9. **Adapt `internal/generate/`** — Modify generate pipeline to write output to `raw/codebase/` under a topic. Depends on steps 1, 2, 7.
10. **`cmd/kb/` and `internal/cli/`** — Rewrite CLI layer with new command taxonomy. Rename binary. Wire all commands to their implementations. Depends on all previous steps.
11. **Integration tests and Makefile** — End-to-end tests, update build targets, update `make verify`. Depends on step 10.

### Technical Dependencies

- **pdfcpu**: `go get github.com/pdfcpu/pdfcpu`
- **html-to-markdown v2**: `go get github.com/JohannesKaufmann/html-to-markdown/v2`
- **excelize**: `go get github.com/xuri/excelize/v2`
- **gosseract**: `go get github.com/otiai10/gosseract/v2` (requires Tesseract system library)
- **yaml.v3**: `go get gopkg.in/yaml.v3` (for frontmatter parsing)
- **yt-dlp**: required runtime binary for `ingest youtube`
- **ffmpeg**: required runtime binary when STT audio must be segmented
- No infrastructure or external service dependencies beyond QMD, Firecrawl, YouTube, and configured STT providers.

## Monitoring and Observability

- **Structured logging**: Existing `slog` usage extends to all new packages. Each ingest operation logs source type, file path, conversion duration, and output path at info level.
- **Log.md audit trail**: Every CLI operation that modifies the vault appends to `<topic>/log.md` with the standard format: `## [YYYY-MM-DD] <op> | <description>`.
- **Lint reports**: Saved to `outputs/reports/<date>-lint.md` for historical tracking.
- **Conversion errors**: Logged at warn/error level with source path and error detail. Non-fatal conversion failures (e.g., OCR unavailable) produce warnings, not hard failures.

## Technical Considerations

### Key Decisions

- **Converter registry over monolithic ingest**: Each format is independently maintained and testable. New formats require only implementing the `Converter` interface. Trade-off: more files, but each is small and focused.
- **ZIP+XML for Office formats**: DOCX, PPTX, and EPUB are all ZIP archives containing XML. Parsing the XML directly (as markitdown does) avoids heavy Office library dependencies. Trade-off: may miss some formatting details but captures text and structure.
- **YouTube captions-first, STT selectable**: `yt-dlp` is the single YouTube backend for captions and audio. The default `captions` policy avoids STT cost; `auto` uses manual captions when present and falls back to STT when only automatic captions or no captions are available; `stt` forces audio transcription. Trade-off: STT requires provider credentials, costs money, and adds latency, but covers videos without usable captions.
- **Codebase as ingest source, not top-level command**: The `generate` command becomes `ingest codebase`. The full pipeline (scan → parse → graph → metrics → render) is preserved but output goes to `raw/codebase/` within a topic. `inspect` subcommands continue to work against codebase-ingested data. Trade-off: `generate` users must adapt to `ingest codebase`, but the conceptual model is cleaner.

### Known Risks

- **PDF extraction quality**: pdfcpu's text extraction may struggle with complex layouts (multi-column, scanned documents). Mitigation: support a `--ocr` flag that uses Tesseract as a secondary extraction method.
- **DOCX/PPTX coverage**: Rolling our own ZIP+XML parser means we may miss edge cases in complex Office documents. Mitigation: test against a diverse set of real-world documents; accept that 90% coverage is sufficient for knowledge base ingestion.
- **YouTube API changes**: Caption and media extraction depends on YouTube behavior exposed through `yt-dlp`. Mitigation: keep `yt-dlp` updated, support proxy/cookie/user-agent configuration, and use STT when captions are unavailable.
- **Binary size increase**: Adding pdfcpu, excelize, and gosseract bindings will increase the binary. Mitigation: use Go build tags to make OCR (gosseract/Tesseract) opt-in at compile time.

## Architecture Decision Records

- [ADR-001: Topic-Centric CLI Command Taxonomy](adrs/adr-001.md) — Adopt topic/ingest/lint/index/search grouping with ingest subcommands per source type
- [ADR-002: Rename Binary to `kb`](adrs/adr-002.md) — Rename from `kodebase` to `kb` to reflect broader knowledge base focus
- [ADR-003: Native Go Document Conversion with Converter Registry](adrs/adr-003.md) — Implement format conversion natively in Go using a converter registry pattern inspired by markitdown
- [ADR-004: Firecrawl REST API for URL Scraping](adrs/adr-004.md) — Integrate with firecrawl via REST API rather than CLI shell-out
- [ADR-005: Native Go Vault Lint Engine](adrs/adr-005.md) — Rewrite lint logic in Go as a first-class CLI command
