# MEMORY — instagram-ingest

- **Goal:** Implementar `kb ingest instagram <url>` (reel/post-vídeo/IGTV por URL única),
  reusando o pipeline yt-dlp+STT do YouTube via extração de um núcleo compartilhado.
  Success = `make verify` verde + feature funcional + testes.

- **Constraints/Assumptions:**
  - `make verify` (fmt→lint→test→build→boundaries) é gate inegociável; lint zero findings.
  - Não quebrar YouTube nem `kb ingest channel`.
  - yt-dlp já instalado (2026.06.09); `instagram:user` BROKEN ⇒ sem bulk de perfil.
  - Firecrawl recusa Instagram (descartado). Plano aprovado em `.claude/plans/ok-fa-a-um-plano-lazy-squirrel.md`.

- **Key decisions:**
  - Núcleo novo `internal/mediadl` (yt-dlp backend + STT + Extractor genérico).
  - `youtube`/`instagram` = shims finos; `youtube` mantém type-aliases p/ CLI+channel.
  - `mediadl.Extractor.Extract(ctx, ParsedURL, opts)` — parsing fica no shim.
  - Config: seção `[instagram]` dedicada; default transcription `auto`.
  - Corpo do doc: `## Caption` + `## Transcript`; degrada p/ caption-only se STT vazio.
  - SourceKind `instagram-video`, raw dir `instagram`.
  - mediadl exporta: ParsedURL, BackendConfig (com Platform p/ msg de erro), PlaylistEntry,
    NewExtractor, ListPlaylistEntries, IsNetworkBlocked, IsContextError, IsTranscriptUnavailable.

- **State: COMPLETE — `make verify` green (788 tests, boundaries OK), integration test passes.**
  - **Done:** Fase 1 `internal/mediadl` (núcleo + testes movidos); Fase 2 `youtube` slim shim +
    aliases + channel via ListPlaylistEntries; Fase 3 `internal/instagram`; config `[instagram]`;
    models SourceKindInstagramVideo; ingest raw dir; CLI `ingest_instagram.go` + register + tests;
    config.example.toml + CLAUDE.md docs; integration test (fake yt-dlp, no network).
  - **Now:** Nothing — feature shipped. Awaiting user review (no commit made; on main, would branch).
  - **Next (out of scope / possible follow-ups):** INSTAGRAM_* env overrides for parity with
    YOUTUBE_*; real-network smoke test; bulk-profile ingest if yt-dlp `instagram:user` is fixed.

- **Open questions:** nenhuma (todas resolvidas no grilling).

- **Working set:**
  - internal/youtube/{youtube,ytdlp,transcription,openai,openrouter,channel}.go (+ _test)
  - internal/{config,models,ingest,cli}/...
  - config.example.toml, CLAUDE.md
