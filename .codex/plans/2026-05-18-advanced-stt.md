# Implementar STT Avancado Greenfield para `kb ingest youtube`

## Summary
- Direcao atual aceita pelo usuario: greenfield, sem back-compat para o caminho antigo.
- `yt-dlp` e o unico backend de YouTube para metadata, captions e audio; se ele falhar, a ingestao falha com erro acionavel.
- A CLI expõe apenas `--transcribe captions|auto|stt`; a flag antiga foi removida.
- STT usa providers plugaveis (`openai` default e `openrouter` opcional), chunking com `ffmpeg`, e proveniencia no frontmatter.

## API, Config e CLI
- Adicionar `youtube.TranscriptionPolicy` com valores `captions`, `auto`, `stt`.
  - `captions`: usa captions do YouTube sem custo de STT.
  - `auto`: usa caption manual quando existir; se houver apenas caption automatica ou nenhuma caption, usa STT.
  - `stt`: ignora captions e transcreve audio como fonte primaria.
- Adicionar flag `kb ingest youtube --transcribe captions|auto|stt`; resolver default por config `[youtube].transcription`, caindo para `captions`.
- Remover a flag antiga de STT; comandos com a flag removida devem falhar como flag desconhecida.
- Adicionar `[stt]` em `Config`, `config.example.toml` e env overrides:
  - `provider = "openai"` com suporte a `openai|openrouter`.
  - `model = "gpt-4o-transcribe"` para OpenAI; OpenRouter usa `openrouter.stt_model` quando selecionado.
  - `api_url`, `api_key`, `language = "auto"`, `prompt = ""`, `audio_format = "mp3"`, `chunk_duration = "10m"`, `max_chunk_bytes = 24000000`, `concurrency = 2`, `ffmpeg_path = "ffmpeg"`.
  - Env novos: `OPENAI_API_KEY`, `OPENAI_API_URL`, `STT_PROVIDER`, `STT_MODEL`; envs `OPENROUTER_API_KEY` e `OPENROUTER_API_URL` configuram provider `openrouter`.
  - Env OpenAI deve ser provider-aware e nao vazar para OpenRouter quando `STT_PROVIDER=openrouter`.

## Implementation Changes
- Refatorar `internal/youtube` para separar responsabilidades:
  - `ytDLPBackend.loadInfo` obtem metadata.
  - `ytDLPBackend.extractCaptionsFromInfo` aceita se captions automaticas sao permitidas.
  - `ytDLPBackend.downloadAudio` baixa audio com `yt-dlp --extract-audio --audio-format <format>`, reaproveitando `baseArgs()` para proxy, cookies, user-agent e retries.
- Ajustar o orquestrador `Extractor.Extract`:
  - `captions` chama apenas captions.
  - `auto` busca captions manuais; se nao houver manual, chama STT via audio do `yt-dlp`.
  - `stt` baixa audio e transcreve sem tentar captions.
  - Nao ha fallback legado de YouTube nem caminho alternativo de audio.
- Implementar providers:
  - `OpenAITranscriber` usa `POST /v1/audio/transcriptions` com multipart form data, `model`, `response_format=json`, `language` quando nao for `auto`, e `prompt` quando configurado.
  - `OpenRouterClient` preserva o modo compativel via chat completions e consome `prompt`/`language` configurados.
- Implementar chunking:
  - Baixar um arquivo temporario de audio via `yt-dlp`.
  - Se o arquivo exceder `max_chunk_bytes` ou a duracao configurada exigir divisao, segmentar com `ffmpeg -f segment` em blocos de `chunk_duration`.
  - Transcrever chunks com concorrencia limitada por `stt.concurrency`, preservar ordem e gerar markdown com headings por offset (`## 00:00`, `## 10:00`, etc.).
  - Se `ffmpeg` estiver ausente ou um chunk continuar acima do limite, retornar erro explicito e acionavel.
- Estender ingest para proveniencia sem acoplar YouTube ao writer:
  - Adicionar `ExtraFrontmatter map[string]any` em `ingest.Options`, validando que nao sobrescreve chaves base como `title`, `type`, `stage`, `source_kind`, `source_url`.
  - `ingest youtube` passa `transcript_source`, `transcription_policy`, `transcript_language`, `caption_kind` quando aplicavel, e `stt_provider`/`stt_model` quando a fonte for STT.

## Test Plan
- `ytDLPBackend.downloadAudio` usa `--extract-audio`, `--audio-format`, `--proxy`, `--cookies`, `--user-agent`, retries e retorna arquivo/formato gerado pelo fake `yt-dlp`.
- `Extractor` usa audio do `yt-dlp` para STT quando captions manuais estao indisponiveis.
- `Extractor` falha quando `yt-dlp` esta indisponivel, sem fallback alternativo.
- `captions`, `auto` e `stt` escolhem corretamente entre manual captions, automatic captions e STT.
- `OpenAITranscriber` envia multipart correto para `/v1/audio/transcriptions`, parseia `{ "text": ... }`, reporta erro de API e valida chave ausente com mensagem util.
- Chunking preserva ordem, offsets e limite de concorrencia usando fake `ffmpeg` e fake transcriber.
- CLI aceita `--transcribe`, rejeita a flag removida como desconhecida, e repassa a policy para o extractor.
- Config carrega `[stt]`, aplica defaults/env overrides provider-aware e rejeita provider/policy/concurrency/formato invalidos.
- Ingest grava os campos de proveniencia no frontmatter e rejeita `ExtraFrontmatter` que tente sobrescrever chaves reservadas.

## Verification
- `rtk go test ./internal/youtube ./internal/config ./internal/cli ./internal/ingest`
- `rtk make verify`
- Sem chamadas reais a OpenAI/OpenRouter nos testes; usar `httptest` e scripts fake para `yt-dlp`/`ffmpeg`.

## Assumptions
- Provider `local` fica fora deste rollout; a nova interface deixa a extensao simples sem aumentar a superficie agora.
- Diarizacao nao entra neste rollout; `response_format=json` e texto simples sao suficientes para KB raw docs.
- `ffmpeg` e `yt-dlp` sao dependencias externas para STT avancado; erros devem dizer exatamente qual binario/config falta.

