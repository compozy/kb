<div align="center">

# kb

[English](README.md) | [Português (Brasil)](README.pt-BR.md)

### Crie e mantenha bases de conhecimento organizadas por tópico.

[![Licença MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/pedronauck/kodebase-go/ci.yaml?branch=main&label=CI)](https://github.com/pedronauck/kodebase-go/actions)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)](https://go.dev/)

[Instalação](#instalação) &#8226; [Em ação](#veja-em-ação) &#8226; [Recursos](#recursos) &#8226; [Comandos](#comandos) &#8226; [Contribuição](#contribuindo)

</div>

---

`kb` é uma CLI Go de binário único para criar e manter bases de conhecimento organizadas por tópico, no padrão [Karpathy KB](https://github.com/karpathy). Ela cuida do fluxo sem LLM: criação de tópicos, ingestão de várias fontes (URLs, arquivos, YouTube, código e favoritos), lint estrutural, análise de código, indexação/busca com QMD e comandos de inspeção. A compilação com LLM fica na sua camada de agentes.

Sem SaaS. Sem nuvem. Apenas Markdown.

---

## Instalação

> [!NOTE]
> O `kb` funciona ainda melhor com sua skill complementar em [`skills/`](skills/). Instale-a com `npx skills add https://github.com/compozy/kb --skill kb`.

### Homebrew

```bash
brew install compozy/kb/kb
```

### npm

```bash
npm install -g @compozy/kb
```

### Go

```bash
go install github.com/compozy/kb/cmd/kb@latest
```

### Compilar do código-fonte

```bash
git clone https://github.com/compozy/kb.git
cd kb
make build
# o binário fica em bin/kb
```

**Opcional** — para busca semântica:

```bash
npm install -g @tobilu/qmd
```

> [!NOTE]
> **Requisitos:** Go >= 1.24 é necessário apenas para compilar a partir do código-fonte. O binário distribuído do `kb` não tem dependência obrigatória em tempo de execução; instale as ferramentas abaixo apenas para as capacidades que utilizar.

| Capacidade | Dependência ou configuração adicional |
| --- | --- |
| `ingest url` (padrão `http-local`) | Nenhuma. O `kb` busca URLs HTTPS públicas diretamente. |
| `ingest url --provider browser` / `--render` | Um executável Chrome ou Chromium configurado em `[browser].command`. |
| `ingest url --provider firecrawl` | Uma chave da API do [Firecrawl](https://firecrawl.dev) (`FIRECRAWL_API_KEY`). |
| `ingest youtube` com legendas | [yt-dlp](https://github.com/yt-dlp/yt-dlp) no `PATH` ou configurado em `youtube.yt_dlp_path`. |
| `ingest youtube --transcribe auto\|stt` | `yt-dlp` para baixar áudio e `ffmpeg` para conversão e segmentação de áudios longos, além do provedor STT escolhido. |
| STT OpenAI / OpenRouter | A chave de API do provedor (`OPENAI_API_KEY` ou `OPENROUTER_API_KEY`). |
| STT local `whispercpp` | `whisper-server` do whisper.cpp e um modelo GGML compatível; configure seus caminhos em `[whispercpp]`. Não requer chave de API nem upload de áudio. |
| `search` / `index` | [QMD](https://github.com/tobilu/qmd), instalado separadamente. |

<details>
<summary><strong>O que ele acessa</strong></summary>

- **Cria arquivos** em `.kb/vault/` no repositório alvo (ou em um `--vault` personalizado).
- **Lê** arquivos-fonte no repositório alvo (nunca os modifica).
- **Chamadas de rede** — `ingest url` busca a URL HTTPS diretamente por padrão; os providers opcionais Firecrawl e browser fazem chamadas de rede/processo adicionais. `ingest youtube` acessa o YouTube pelo `yt-dlp`; `ingest youtube --transcribe auto|stt` também usa o provedor STT configurado. Todos os outros comandos são totalmente locais.
- **Sem telemetria** — nada é enviado em segundo plano.
- **Desinstalação:** remova o binário `kb` do seu `PATH` e apague o diretório `.kb/`.

</details>

---

## Veja em ação

Crie um tópico e faça ingestão de conteúdo de várias fontes:

```bash
# cria um novo tópico
$ kb topic new rust-lang "Rust Language" programming
{
  "slug": "rust-lang",
  "title": "Rust Language",
  "domain": "programming"
}

# cria um tópico em bundle OKF
$ kb topic new rust-catalog "Rust Catalog" programming --mode okf

# ingere um artigo da web
$ kb ingest url https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html --topic rust-lang

# ingere um PDF local
$ kb ingest file ./rust-reference.pdf --topic rust-lang

# ingere a transcrição de um vídeo do YouTube
$ kb ingest youtube https://www.youtube.com/watch?v=... --topic rust-lang

# ingere um snapshot completo de um código-fonte
# a primeira execução cria automaticamente ./my-rust-project/.kb/vault/rust-lang
$ kb ingest codebase ./my-rust-project --topic rust-lang --progress never

# verifica problemas estruturais no tópico
$ kb lint rust-lang

# promove um documento wiki compilado ao bundle OKF
$ kb promote .kb/vault/rust-lang/wiki/concepts/Ownership.md --to rust-catalog --type Concept

# verifica o bundle OKF
$ kb okf check rust-catalog --strict
```

Analise snapshots de código pelo terminal:

```bash
$ kb inspect complexity --top 5

 symbol_name       | cyclomatic_complexity | loc | source_path
 computeMetrics    | 12                    | 89  | src/compute-metrics.ts
 parseTypeScript   | 9                     | 67  | src/adapters/typescript.ts
 normalizeGraph    | 8                     | 45  | src/normalize-graph.ts
 renderDocuments   | 7                     | 112 | src/render-documents.ts
 scanWorkspace     | 6                     | 54  | src/scan-workspace.ts
```

```bash
$ kb inspect dead-code

 kind   | name           | source_path                | reason
 symbol | oldHelper      | src/utils.ts               | dead-export
 file   | unused-cfg.ts  | src/config/unused-cfg.ts   | orphan-file
```

Busque linguagem natural no seu vault (requer QMD):

```bash
$ kb search "error handling patterns" --limit 3

 title                  | score | path
 Error Handling Guide   | 0.89  | wiki/concepts/Error Handling.md
 src/error-boundary.ts  | 0.74  | raw/codebase/files/src/error-boundary.ts.md
 handleError            | 0.61  | raw/codebase/symbols/handleError--src-utils-l42.md
```

---

## Recursos

**Bases de conhecimento por tópico** — o `kb` organiza conhecimento em tópicos, cada um com seus próprios diretórios `raw/`, `wiki/`, `outputs/` e `bases/`. Use `kb topic new` para criar manualmente ou deixe `kb ingest codebase` inicializar um tópico na primeira execução.

**Ingestão de várias fontes** — ingira artigos da web diretamente por HTTPS (com renderização por browser ou Firecrawl opcionais), arquivos locais (PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, texto simples e imagens com OCR), legendas ou transcrições STT do YouTube, bases de código e grupos de favoritos. Cada fonte passa por um registro de conversores que normaliza o conteúdo em Markdown com frontmatter.

**Análise de código** — aponte `kb ingest codebase` para um repositório e ele gera uma camada de vault [Obsidian](https://obsidian.md) com símbolos, arquivos e relações de dependência mapeados em notas Markdown interligadas. Ele calcula complexidade ciclomática, raio de impacto, acoplamento, instabilidade e código morto; depois compila artigos wiki e visualizações interativas do [Base](https://obsidian.md/blog/bases/).

**Lint estrutural** — `kb lint` verifica frontmatter ausente, wikilinks quebrados, arquivos órfãos e outros problemas estruturais. Os relatórios podem ser salvos em Markdown em `outputs/reports/`.

**10 subcomandos de inspeção** — consulte seu código como um banco de dados pelo terminal. Classifique funções por complexidade, encontre exports mortos, siga cadeias de dependências e detecte imports circulares. Três formatos de saída: tabela, JSON e TSV; sem dependências externas.

**Saída nativa do Obsidian** — todas as notas geradas usam wikilinks, frontmatter YAML e referências compatíveis com backlinks. As visualizações Base fornecem tabelas filtráveis e ordenáveis no Obsidian: explorador de símbolos, pontos de complexidade, zona de perigo, saúde de módulos e muito mais.

**Busca semântica** — indexe seu vault com QMD para busca híbrida, lexical ou vetorial em toda a documentação. Útil para onboarding, revisão de arquitetura ou para alimentar LLMs com contexto.

**Saída amigável a IA** — cada vault contém `CLAUDE.md` e `AGENTS.md` como documentos de esquema. A saída Markdown estruturada foi feita para consumo direto por LLMs: metadados em frontmatter, estrutura consistente e referências explícitas.

**Binário único** — sem dependências em tempo de execução. Um único binário `kb` faz tudo. Construído em Go para inicialização rápida e baixo uso de memória.

---

## Por que kb

A maioria das ferramentas de análise de código entrega um dashboard. O `kb` entrega uma base de conhecimento.

**SonarQube e CodeClimate** dizem que seu código tem problemas. O `kb` diz quais são, onde se conectam e entrega um espaço de trabalho estruturado para raciocinar sobre eles. A saída é Markdown seu, não um dashboard SaaS que desaparece quando você cancela a assinatura.

**Sourcegraph** é excelente para busca de código entre repositórios. O `kb` serve para entender profundamente um único repositório — sua arquitetura, padrões de acoplamento e superfície de risco — e construir um artefato de conhecimento persistente ao redor dele.

**Obsidian e Notion** são ótimos para anotações. O `kb` automatiza a criação da estrutura e a ingestão para que você comece com organização, não com uma página em branco; depois amplie o vault com suas notas e análises.

A diferença-chave: a saída do `kb` se acumula. Uma análise do SonarQube do mês passado é dado obsoleto. Um vault do `kb` do mês passado é uma base de conhecimento que você vem cultivando há 30 dias.

---

## Comandos

### `kb topic`

Crie e gerencie tópicos da base de conhecimento.

```bash
kb topic new <slug> <title> <domain> [--mode wiki|okf]   # Cria um tópico
kb topic list                           # Lista todos os tópicos do vault
kb topic info <slug>                    # Mostra metadados de um tópico
```

`--mode wiki` é o scaffold padrão. `--mode okf` cria um bundle Open Knowledge File na raiz, com `index.md`, `log.md` e orientações de autoria OKF.

### `kb promote`

Promova um documento wiki compilado para um tópico OKF sem alterar o documento-fonte.

```bash
kb promote <wiki-doc> --to <okf-topic> --type <type> [--description <text>]
```

### `kb okf`

Verifique a conformidade de um tópico OKF com o bundle.

```bash
kb okf check <topic> [--strict] [--format table|json|tsv]
```

### `kb ingest`

Ingira material-fonte em um tópico. `url`, `file`, `youtube` e `bookmarks` exigem um tópico existente; `codebase` pode criá-lo na primeira execução.

```bash
kb ingest url <url> --topic <slug>                # Busca uma URL HTTPS pública localmente
kb ingest url <url> --provider browser --topic <slug> # Renderiza uma página JavaScript com Chromium
kb ingest url <url> --provider firecrawl --topic <slug> # Usa Firecrawl quando configurado
kb ingest url <url> --render --topic <slug>       # Atalho para --provider browser
kb ingest file <path> --topic <slug>              # Converte e ingere um arquivo local
kb ingest youtube <url> --topic <slug> [--transcribe captions|auto|stt] [--sub-langs orig,pt]
kb ingest codebase <path> --topic <slug>          # Analisa código e cria o tópico se necessário
kb ingest bookmarks <path> --topic <slug>         # Ingere um arquivo Markdown de grupo de favoritos
```

**Formatos aceitos** por `ingest file`: PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, texto simples (`.txt`, `.md`) e imagens (PNG, JPG, TIFF, BMP, GIF — com OCR opcional via Tesseract).

| Subcomando de ingestão | `--topic` | Flags adicionais |
| --- | --- | --- |
| `url` | obrigatório | `--provider http-local\|browser\|firecrawl`, `--render` |
| `file` | obrigatório | -- |
| `youtube` | obrigatório | `--transcribe captions\|auto\|stt` (padrão: `captions`), `--sub-langs`/`--lang` (idiomas de legenda; padrão: `orig`) |
| `codebase` | obrigatório | `--vault`, `--output` (apelido obsoleto), `--title`, `--domain`, `--include`, `--exclude`, `--semantic`, `--progress`, `--log-format` |
| `bookmarks` | obrigatório | -- |

#### Providers de ingestão de URL

`kb ingest url` usa o provider local `http-local` por padrão; por isso, artigos HTTPS públicos comuns não requerem API ou credenciais de terceiros. Ele segue redirecionamentos apenas para hosts HTTPS públicos, rejeita localhost e IPs privados ou link-local (inclusive endereços resolvidos por DNS), desabilita proxies e limita a resposta a 20 MiB. O conteúdo baixado passa pelo registro normal de conversores; portanto, HTML e formatos de documento aceitos são ingeridos como arquivos locais. O frontmatter do resultado registra a URL final, horário da busca, tipo de conteúdo e hash SHA-256.

Use `--provider browser` (ou `--render`) para páginas renderizadas em JavaScript após configurar `[browser].command`. Use `--provider firecrawl` quando preferir o serviço hospedado de extração do Firecrawl; ele requer `FIRECRAWL_API_KEY`. Todos os providers aceitam somente URLs HTTPS públicas.

### `kb lint`

Verifique problemas estruturais de um tópico.

```bash
kb lint [<slug>] [--format table|json|tsv] [--save] [--topic <slug>]
```

`--save` grava um relatório Markdown em `outputs/reports/<date>-lint.md`.

### `kb inspect`

Consulte dados do vault de código usando frontmatter e métricas extraídas.

```bash
kb inspect <subcommand> [options]
```

| Categoria | Subcomando | Descrição |
| --- | --- | --- |
| Métricas | `smells` | Lista símbolos e arquivos com code smells detectados |
| Métricas | `dead-code` | Lista exports mortos e arquivos órfãos |
| Métricas | `complexity` | Classifica funções por complexidade ciclomática |
| Métricas | `blast-radius` | Classifica símbolos por raio de impacto (dependentes transitivos) |
| Métricas | `coupling` | Classifica arquivos por instabilidade (eferente/aferente) |
| Grafo | `backlinks` | Mostra quem referencia um símbolo informado |
| Grafo | `deps` | Mostra relações de saída de um arquivo ou símbolo |
| Grafo | `circular-deps` | Lista arquivos em ciclos de dependência circular |
| Consulta | `symbol` | Busca aproximada e visão detalhada de um símbolo |
| Consulta | `file` | Consulta exata de arquivo pelo caminho-fonte |

Flags compartilhadas: `--format table|json|tsv` (padrão: `table`), `--vault <path>`, `--topic <slug>`.

### `kb search`

Busca semântica nos documentos do vault. Requer [QMD](https://github.com/tobilu/qmd).

```bash
kb search <query> [options]
```

| Flag | Padrão | Descrição |
| --- | --- | --- |
| `--lex` | false | Usa somente busca por palavra-chave BM25 |
| `--vec` | false | Usa somente busca por similaridade vetorial |
| `--limit` | 10 | Máximo de resultados retornados |
| `--full` | false | Mostra o documento completo em vez de trecho |
| `--min-score` | -- | Limiar mínimo de similaridade |
| `--all` | false | Retorna todos os resultados acima do limiar |
| `--collection` | -- | Nome explícito da coleção QMD |
| `--format` | table | Formato: `table`, `json` ou `tsv` |

O modo padrão é híbrido (BM25 + similaridade vetorial). Use `--lex` ou `--vec` para restringir. Em hosts nos quais o QMD não executa busca vetorial por ausência de `sqlite-vec`, o modo híbrido padrão recua automaticamente para busca lexical.

### `kb index`

Crie ou atualize uma coleção QMD para indexação semântica. Requer [QMD](https://github.com/tobilu/qmd).

```bash
kb index [options]
```

| Flag | Padrão | Descrição |
| --- | --- | --- |
| `--vault` | -- | Raiz do vault |
| `--topic` | -- | Slug do tópico no vault |
| `--name` | -- | Sobrescreve o nome derivado da coleção QMD |
| `--embed` | true | Executa embeddings após sincronizar arquivos |
| `--context` | -- | Anexa contexto para melhorar relevância da busca |
| `--force-embed` | false | Força novo embedding de todos os documentos |

Se o QMD estiver instalado sem suporte a `sqlite-vec`, `kb index` ainda sincroniza a coleção e informa `embedStatus: "skipped_unavailable"` na saída JSON para que a busca lexical continue disponível. Use `--force-embed` apenas quando o embedding vetorial for necessário no host atual.

### `kb version`

Imprime metadados da versão de build.

---

## O que é gerado

```text
.kb/vault/
  <topic-slug>/
    CLAUDE.md                    # Documento de esquema para LLMs
    AGENTS.md                    # Referência do projeto para agentes
    log.md                       # Log de operações somente de acréscimo
    raw/
      codebase/                  # Snapshot de código gerado por máquina
        files/                   #   Uma nota Markdown por arquivo-fonte
        symbols/                 #   Uma nota Markdown por símbolo extraído
        indexes/
          directories/           #   Inventários por diretório
          languages/             #   Inventários por linguagem
      <ingested-sources>.md      # Artigos, transcrições e documentos ingeridos
    wiki/
      concepts/                  # Artigos wiki sintetizados (tópicos de código)
        Codebase Overview.md
        Directory Map.md
        Symbol Taxonomy.md
        Dependency Hotspots.md
        Complexity Hotspots.md
        Module Health.md
        Dead Code Report.md
        Code Smells.md
        Circular Dependencies.md
        High-Impact Symbols.md
      index/
        Dashboard.md             # Página inicial
        Concept Index.md         # Lista de artigos
        Source Index.md          # Índice reverso
    outputs/                     # Seus resultados de análise (preservados)
      reports/                   # Relatórios de lint (com --save)
    bases/                       # Visualizações Obsidian Base
```

- **`raw/`** — snapshots de fontes gerados por máquina e documentos ingeridos. Snapshots de código são atualizados em toda execução; documentos ingeridos são somente de acréscimo.
- **`wiki/`** — artigos iniciais sintetizados a partir de métricas de código. Áreas gerenciadas são atualizadas; suas adições são preservadas.
- **`outputs/`** — lugar para seus briefings, consultas, diagramas e relatórios. Nunca é alterado pelo `kb`.
- **`bases/`** — arquivos `.base` do Obsidian Base para visualizações interativas de tabela/cartão/lista de métricas.

---

## Linguagens aceitas (análise de código)

| Linguagem | Extensões | Parser | Confiança da relação |
| --- | --- | --- | --- |
| TypeScript | `.ts`, `.tsx` | `tree-sitter` | sintática |
| JavaScript | `.js`, `.jsx` | `tree-sitter` | sintática |
| Go | `.go` | `tree-sitter` | sintática |

Quer adicionar uma linguagem? Veja [CONTRIBUTING.md](CONTRIBUTING.md#adding-a-new-language-adapter).

---

## Formatos de arquivo aceitos (ingestão)

| Formato | Extensões | Observações |
| --- | --- | --- |
| PDF | `.pdf` | Extração nativa de texto por pdfcpu |
| DOCX | `.docx` | Extração baseada em XML |
| XLSX | `.xlsx` | Conversão de planilha em tabela Markdown |
| PPTX | `.pptx` | Extração de texto dos slides |
| EPUB | `.epub` | Extração de capítulos com HTML para Markdown |
| HTML | `.html`, `.htm` | Conversão HTML para Markdown |
| CSV | `.csv` | Conversão de tabela |
| JSON | `.json` | Bloco de código formatado |
| XML | `.xml` | Bloco de código formatado |
| Texto simples | `.txt`, `.md` | Sem conversão |
| Imagens | `.png`, `.jpg`, `.tiff`, `.bmp`, `.gif` | OCR opcional via Tesseract |

---

## Relações rastreadas

| Relação | Descrição |
| --- | --- |
| `imports` | Um arquivo importa outro arquivo ou módulo externo |
| `exports` | Um arquivo exporta um símbolo |
| `calls` | Uma função ou método chama outro símbolo |
| `references` | O código referencia um símbolo |
| `declares` | Um arquivo declara um símbolo |
| `contains` | Um arquivo contém estruturalmente um símbolo |

---

## Code smells detectados

| Smell | Escopo | Condição |
| --- | --- | --- |
| `dead-export` | símbolo | Exportado, mas nunca referenciado fora do arquivo |
| `long-function` | símbolo | Função com > 50 LOC ou complexidade ciclomática > 10 |
| `high-blast-radius` | símbolo | Mais de 20 dependentes transitivos |
| `bottleneck` | símbolo | Centralidade de intermediação > 0,1 |
| `feature-envy` | símbolo | Mais referências a símbolos de outro arquivo que aos próprios |
| `god-file` | arquivo | Mais de 15 símbolos ou acoplamento eferente > 10 |
| `orphan-file` | arquivo | Acoplamento aferente zero e não é ponto de entrada |

---

## Caminhos excluídos

O scanner:

- Sempre ignora `.git`, `.hg`, `.svn`, links simbólicos e a própria raiz de vault configurada.
- Aplica ignores padrão para `vendor/`, `.turbo/`, `.next/`, `node_modules/`, `dist/`, `build/` e `coverage/`.
- Respeita arquivos `.gitignore` encontrados da raiz do scan para baixo, inclusive aninhados.
- Aplica padrões `--exclude` depois das regras de ignore do repositório.
- Aplica padrões `--include` por último, como reinclusões explícitas.

`--include` e `--exclude` usam padrões no estilo `.gitignore`, avaliados contra caminhos relativos à raiz do scan.

---

## Configuração

O `kb` é configurado principalmente por flags de CLI. A configuração opcional de runtime é carregada de um arquivo TOML e de variáveis de ambiente.

| Variável | Origem | Descrição |
| --- | --- | --- |
| `APP_CONFIG` | env | Caminho para o arquivo de configuração TOML |
| `FIRECRAWL_API_KEY` | env / TOML | Chave Firecrawl para `ingest url --provider firecrawl` |
| `FIRECRAWL_API_URL` | env / TOML | Endpoint Firecrawl |
| `browser.command` | TOML | Executável Chrome/Chromium para `ingest url --provider browser` ou `--render` |
| `OPENAI_API_KEY` | env / TOML | Chave OpenAI STT do provider `openai` padrão |
| `OPENAI_API_URL` | env / TOML | URL-base STT compatível com OpenAI |
| `STT_PROVIDER` | env / TOML | Provider STT: `openai`, `openrouter` ou `whispercpp` local |
| `STT_MODEL` | env / TOML | Sobrescrita de modelo STT, por exemplo `gpt-4o-transcribe` |
| `OPENROUTER_API_KEY` | env / TOML | Chave OpenRouter quando `stt.provider = "openrouter"` |
| `OPENROUTER_API_URL` | env / TOML | Endpoint OpenRouter |
| `YOUTUBE_YT_DLP_PATH` | env / TOML | Caminho do executável `yt-dlp` para `ingest youtube` |
| `YOUTUBE_PROXY` | env / TOML | URL de proxy passada ao `yt-dlp` |
| `YOUTUBE_COOKIES_FILE` | env / TOML | Arquivo Netscape de cookies passado ao `yt-dlp` |
| `YOUTUBE_USER_AGENT` | env / TOML | User-Agent passado ao `yt-dlp` |
| `YOUTUBE_CAPTION_LANGUAGES` | env / TOML | Lista de preferência de idiomas de legenda separada por vírgulas, por exemplo `orig,pt,en` |

`ingest youtube` aceita três políticas de transcrição: `captions` usa apenas legendas do YouTube, `stt` força transcrição de áudio e `auto` usa legendas manuais quando disponíveis e recorre a STT quando existem somente legendas automáticas ou não há legendas. STT grava frontmatter de proveniência como `transcript_source`, `stt_provider` e `stt_model`.

A seleção de legendas usa por padrão `caption_languages = ["orig"]`, para que vídeos não ingleses utilizem a faixa de legenda nativa/original em vez da tradução automática em inglês do YouTube. `orig` é resolvido por vídeo a partir dos metadados do yt-dlp, com faixas automáticas `<lang>-orig` preferidas quando não houver legenda manual no idioma original. Sobrescreva por comando com `--sub-langs` ou `--lang`; valores da CLI têm prioridade sobre `YOUTUBE_CAPTION_LANGUAGES`, que tem prioridade sobre `[youtube].caption_languages`. Legendas automáticas traduzidas ficam desabilitadas, salvo se `[youtube].allow_translated_captions = true`, pois usam o endpoint de tradução do YouTube e sofrem mais rate limiting.

Ingestões do YouTube também persistem metadados do vídeo do `yt-dlp` no frontmatter da transcrição para ordenação e filtragem: contagens de engajamento (`view_count`, `like_count`, `comment_count`), contexto de publicação (`upload_date`, `duration`, `duration_string`, `channel`, `channel_id`, `uploader_id`, `channel_follower_count`) e classificação (`categories`, `youtube_tags`, `language`, `live_status`, `was_live`, `chapter_count`). Campos escalares ocultos ou ausentes recebem `null`; `categories` e `youtube_tags` recebem listas vazias quando ausentes; `chapter_count` é `0` quando não há metadados de capítulos. Palavras-chave do vídeo usam `youtube_tags`, porque `tags` continua sendo o campo de taxonomia do KB; apenas a contagem de capítulos é armazenada, não seu conteúdo completo.

STT tem custo e latência reais. Prefira `captions` para ingestão normal e use `auto` ou `stt` quando as legendas do YouTube não estiverem disponíveis ou forem inadequadas.

### Transcrição local com whisper.cpp

Defina `stt.provider` como `whispercpp` para executar a transcrição localmente. Quando STT for necessário, o `kb` reutiliza um servidor já escutando no endereço configurado ou inicia o próprio `whisper-server`; não há wrapper PowerShell, chave de API ou upload externo de áudio.

```toml
[stt]
provider = "whispercpp"
model = "whisper.cpp-small"

[whispercpp]
server_path = "C:/Users/me/.local/whisper.cpp/cuda/Release/whisper-server.exe"
model_path = "C:/Users/me/.local/whisper.cpp/models/ggml-small.bin"
host = "127.0.0.1"
port = 8188
startup_timeout = "30s"
```

A instalação local ainda requer `whisper-server` e um modelo GGML compatível. O `kb` inicia o servidor com `--convert` e a rota compatível `/v1/audio/transcriptions`, depois o mantém em execução para ingestões posteriores.

Veja [`config.example.toml`](config.example.toml) para o esquema TOML completo.

---

## Desenvolvimento

**Pré-requisito:** [Go](https://go.dev) >= 1.24

```bash
git clone https://github.com/pedronauck/kodebase-go.git
cd kodebase-go
make verify    # format + lint + test + build + boundaries
```

| Comando | Descrição |
| --- | --- |
| `make fmt` | Formata todos os arquivos Go com gofmt |
| `make lint` | Executa golangci-lint com tolerância zero |
| `make test` | Testes unitários com detector de corrida |
| `make test-integration` | Testes unitários + de integração |
| `make build` | Compila o binário em `bin/kb` |
| `make verify` | fmt -> lint -> test -> build -> boundaries |
| `make deps` | Executa `go mod tidy` |

Consulte [CONTRIBUTING.md](CONTRIBUTING.md) para estilo de código, requisitos de teste e como adicionar um adaptador de linguagem.

---

## Contribuindo

O `kb` usa licença MIT e é desenvolvido abertamente. Contribuições de todos os tipos são bem-vindas:

- **Adaptadores de linguagem** — adicione suporte a Python, Rust, Java ou qualquer linguagem com gramática tree-sitter.
- **Conversores de arquivo** — adicione suporte a novos formatos no registro de conversores.
- **Novos detectores de code smell** — o motor de métricas foi feito para ser estendido.
- **Modelos de artigo wiki** — artigos iniciais melhores geram vaults melhores.
- **Bugs e sugestões** — [abra uma issue](https://github.com/pedronauck/kodebase-go/issues), nós lemos todas.

Consulte [CONTRIBUTING.md](CONTRIBUTING.md) para a preparação do ambiente e as diretrizes.

---

## Licença

MIT — veja [LICENSE](LICENSE).
