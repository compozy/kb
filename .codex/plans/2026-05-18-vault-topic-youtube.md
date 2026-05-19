# Correção Greenfield do Modelo de Vault, Tópicos e YouTube

## Summary

Corrigir o `kb` para tratar a raiz do repo como vault Obsidian quando houver `kb.toml`, resolver tópicos por path relativo, unificar a validade de tópico em torno de `CLAUDE.md`, usar `topic.yaml` como fonte estruturada de metadata, manter `raw/youtube` como layout canônico e melhorar a extração do YouTube com proxy/cookies/retry/erros acionáveis.

A implementação pode fazer breaking changes porque o projeto está em alpha. O objetivo é remover decisões antigas inconsistentes, não preservar compatibilidade por alias permanente.

## Key Changes

- Adicionar suporte a `kb.toml` na raiz do vault:
  - Exemplo de contrato:
    ```toml
    [vault]
    root = "."
    topic_globs = ["*", "harness/*", "social-media/*"]
    ```
  - `vault.DiscoverVaultPath` passa a subir diretórios procurando `kb.toml` primeiro; se encontrado, resolve `[vault].root` relativo ao diretório do config.
  - Manter `.kb/vault` apenas como fallback legado enquanto testes existentes forem atualizados; `kb.toml` tem precedência.

- Unificar resolução de tópicos:
  - Criar uma API compartilhada, por exemplo `topic.Resolve(vaultPath, topicRef)` ou equivalente, usada por `topic.Info`, `topic.List`, `vault.ResolveVaultQuery`, ingest, lint/search/index/inspect.
  - Tópico válido passa a significar: diretório existente com `CLAUDE.md` arquivo.
  - `--topic harness/goclaw` é o identificador canônico para tópicos aninhados.
  - `topic.List` e `vault.ListAvailableTopics` varrem os `topic_globs` do `kb.toml`, retornando slugs relativos como `go-best-practices` e `harness/goclaw`.
  - Não resolver por leaf-name implícito (`goclaw`) para evitar ambiguidade.

- Tornar skeleton autocurável:
  - `topic.Info` não deve rejeitar tópico por ausência de subpastas, `log.md` ou symlink `AGENTS.md`.
  - `ingest.Ingest` deve chamar `topic.Resolve` e `topic.EnsureCurrentSkeleton` antes de ler metadata/log.
  - `EnsureCurrentSkeleton` deve criar diretórios faltantes, `log.md` se ausente e `AGENTS.md` quando ausente; não deve sobrescrever `CLAUDE.md` existente nem substituir `AGENTS.md` manual sem necessidade.
  - Remover a regra rígida `AGENTS.md` precisa ser symlink exatamente para `"CLAUDE.md"`.

- Usar `topic.yaml` como fonte de metadata:
  - `topic.yaml` vira fonte primária para `slug`, `title` e `domain`.
  - `CLAUDE.md` fica como fallback apenas quando `topic.yaml` não existir ou estiver incompleto.
  - Regex de H1/domain em `CLAUDE.md` deve ser mantida só como fallback legado, não como fonte principal.
  - `topic.New` deve escrever `topic.yaml` junto com `CLAUDE.md`, `AGENTS.md`, `log.md` e skeleton.

- Manter `raw/youtube` como layout canônico:
  - Novos ingests de `SourceKindYouTubeTranscript` continuam escrevendo em `raw/youtube`.
  - `lint` valida schema de transcripts em `raw/youtube`.
  - `raw/transcripts` não vira alias permanente para evitar dois layouts válidos.
  - Adicionar migração explícita, por exemplo `kb migrate transcripts --topic <topic>`, que move arquivos de `raw/transcripts` para `raw/youtube`, preserva conteúdo/frontmatter quando válido, atualiza `log.md` e falha em conflito de nomes sem sobrescrever.

- Melhorar YouTube com suporte real de rede:
  - Expandir config com `[youtube]`:
    ```toml
    [youtube]
    proxy = ""
    cookies_file = ""
    user_agent = ""
    yt_dlp_path = "yt-dlp"
    transcription = "captions"
    retry_attempts = 3
    retry_backoff = "1s"
    ```
  - Adicionar env overrides para valores operacionais: `YOUTUBE_YT_DLP_PATH`, `YOUTUBE_PROXY`, `YOUTUBE_COOKIES_FILE`, `YOUTUBE_USER_AGENT`.
  - Construir `youtube.Extractor` com `YoutubeConfig`, `STTConfig` e `OpenRouterConfig`.
  - Usar `yt-dlp` como backend obrigatório para metadata, captions e audio; proxy/cookies/user-agent/retries são repassados como argumentos do processo.
  - Implementar retry/backoff por argumentos do `yt-dlp`; não repetir erro sem legenda ou URL inválida.
  - Classificar erros:
    - captions sem track: `transcript_unavailable`
    - HTTP 400/403/429 em captions/stream: erro acionável indicando provável bloqueio de IP e sugerindo `[youtube].proxy` ou cookies
    - audio bloqueado durante STT: `audio_unavailable` com causa de rede e ação sobre proxy/cookies/user-agent
  - Não implementar `po_token` nesta rodada; não expor campo sem suporte real.

## Public Interfaces / Types

- `config.Config` ganha `VaultConfig` e `YoutubeConfig`.
- `config.example.toml` documenta `[vault]` e `[youtube]`.
- `internal/vault` passa a carregar/usar configuração de descoberta de vault e globs de tópico.
- `internal/topic` passa a expor uma única resolução tolerante de tópico, usada por todos os comandos.
- CLI:
  - `--topic` aceita path relativo, como `harness/goclaw`.
  - Novo comando de migração explícita para `raw/transcripts -> raw/youtube`.
  - `kb ingest youtube` usa config `[youtube]` automaticamente; `--transcribe captions|auto|stt` controla a política de transcrição.

## Test Plan

- Config:
  - carrega `kb.toml` com `[vault] root/topic_globs`;
  - rejeita keys desconhecidas;
  - aplica env overrides de YouTube;
  - valida proxy/cookies/user-agent/retry defaults.

- Vault/topic:
  - descobre vault pela raiz com `kb.toml`;
  - preserva fallback `.kb/vault`;
  - lista tópicos planos e aninhados por glob;
  - resolve `--topic harness/goclaw`;
  - não resolve leaf ambíguo implicitamente;
  - `topic.Info` aceita tópico parcial com `CLAUDE.md`;
  - `EnsureCurrentSkeleton` autocura diretórios/log sem sobrescrever marker existente;
  - `topic.yaml` vence `CLAUDE.md` para title/domain.

- Ingest/lint/layout:
  - `ingest youtube` escreve em `raw/youtube`;
  - ingest funciona em tópico parcial porque cura antes de ler metadata;
  - lint valida `raw/youtube`;
  - migração move `raw/transcripts/*.md` para `raw/youtube/*.md`, registra log e falha em conflito.

- YouTube:
  - extractor resolve e chama `yt-dlp` com argumentos de rede configurados;
  - proxy configurado é repassado ao `yt-dlp`;
  - cookies configurados são repassados ao `yt-dlp`;
  - retries acontecem para status transitórios e não acontecem para transcript disabled;
  - erro 400/403/429 retorna mensagem acionável sobre proxy/cookies/IP;
  - STT roda conforme `--transcribe auto|stt` e sempre usa audio baixado pelo `yt-dlp`.

- Gate:
  - Rodar `rtk make verify` como verificação final obrigatória.

## Assumptions

- `kb.toml` é o marcador/config público escolhido para vault na raiz.
- `raw/youtube` é o layout canônico; `raw/transcripts` é legado a ser migrado explicitamente.
- Tópicos aninhados são identificados por path relativo completo.
- Como o projeto está em alpha, não vamos manter aliases silenciosos que perpetuem dois modelos válidos.
- `po_token` fica fora porque o adapter atual não oferece suporte direto e expor config sem aplicação real seria uma falsa correção.
