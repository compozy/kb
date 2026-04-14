# Correção Estrutural: Lint de Raw Codebase + Suporte Rust

## Resumo

- Corrigir o `dead-link` falso no lint pela causa raiz: o lint vai parar de interpretar `[[...]]` literais vindos do código em seções descritivas de `raw/codebase/*`.
- Corrigir `openfang` pela causa raiz: Rust passa a ser linguagem suportada de primeira classe na pipeline `scan -> parse -> normalize -> render -> write`.
- Não entram workarounds: sem suppress de lint, sem escaping artificial do texto derivado do código, sem “file-only mode” para Rust, sem adapter regex paralelo.

## Mudanças de Implementação

- `internal/lint` deixa de usar uma extração global de wikilinks para todo markdown e passa a calcular `links` por tipo de documento, guiado por `frontmatter.source_kind`.
- Para `source_kind=codebase-file`, o grafo de links do lint considera apenas `## Symbols`, `## Outgoing Relations` e `## Backlinks`.
- Para `source_kind=codebase-symbol`, o grafo considera apenas a linha `Source file: [[...]]`, `## Outgoing Relations` e `## Backlinks`.
- Para `source_kind=codebase-directory-index` e `codebase-language-index`, o grafo considera apenas `## Files` e `## Symbols`.
- Para todos os outros documentos, o comportamento atual permanece: extrair wikilinks do corpo inteiro.
- A implementação deve ficar em `internal/lint` como um helper explícito de regiões válidas para link graph; pode reutilizar `vault.ExtractSection`, mas a regra de quais regiões contam pertence ao lint.
- O conteúdo de `Module Notes` e `Documentation` continua sendo renderizado exatamente como hoje; o fix é no consumo pelo lint, não na renderização.

- `internal/models` ganha `LangRust`, `SupportedLanguages()` passa a incluir `rust`, e os pontos de ajuda/erro do CLI passam a derivar a lista suportada de uma única fonte compartilhada.
- `internal/scanner` passa a reconhecer `.rs`.
- `go.mod` adiciona o binding oficial `github.com/tree-sitter/tree-sitter-rust/bindings/go`; `internal/adapter/treesitter.go` expõe `rustLanguage()`.
- `internal/adapter` ganha `RustAdapter`, registrado no runner padrão de `internal/generate`.
- O `RustAdapter` extrai, no mínimo, estes símbolos: `module`, `function`, `method`, `struct`, `enum`, `trait`, `typeAlias`, `const`, `static`, `union`, `macro`.
- O adapter extrai doc comments de arquivo via `//!` / `/*! ... */` e doc comments de símbolo via `///` / `/** ... */`; `signature` vem do trecho sintático da declaração.
- `exported=true` segue a visibilidade Rust real: `pub` no item, e `pub use` gera `RelExports` para símbolos internos resolvidos.
- O adapter constrói um índice de resolução de módulos/crates antes do parse final, baseado em `Cargo.toml` e caminhos convencionais (`src/lib.rs`, `src/main.rs`, `mod.rs`, `foo.rs`, `src/bin/*.rs`).
- O adapter resolve `mod`, `use crate::`, `self::`, `super::`, aliases e imports agrupados para arquivos/símbolos internos quando houver mapeamento; o que não resolver vira `ExternalNode`.
- As relações emitidas para Rust serão: `RelContains`, `RelImports`, `RelExports`, `RelCalls` e `RelReferences`; todas com `ConfidenceSyntactic` nesta entrega.
- `RelCalls` cobre chamadas diretas resolvíveis por binding/import/local symbol index. Não entra um modo “semântico” novo nesta entrega.

## Interfaces Públicas

- A lista de linguagens suportadas visível ao usuário passa a incluir `rust` em:
  - `kb ingest codebase --help`
  - `kb generate --help`
  - mensagens de erro de scan sem arquivos suportados
  - `GenerationSummary.detectedLanguages`
- Nenhum novo flag ou arquivo de configuração é adicionado.
- O layout de vault e o ownership model recém-estabilizado permanecem iguais: `raw/codebase/*`, `wiki/codebase/*`, bridges em `wiki/index/*`, e bloco gerenciado em `CLAUDE.md`.

## Testes e Aceite

- Unitários de lint:
  - `codebase-file` com `Module Notes` contendo `[[literal]]` não gera `dead-link`, mas links em `Symbols`/`Outgoing Relations`/`Backlinks` continuam válidos.
  - `codebase-symbol` com `Documentation` contendo `[[literal]]` não gera `dead-link`, mas a linha `Source file: [[...]]` e as seções de relações continuam sendo lidas.
  - `codebase-directory-index` e `codebase-language-index` continuam contribuindo links via `Files`/`Symbols`.
- Unitários do adapter Rust:
  - suporta `LangRust`
  - extrai docs, signatures e visibilidade
  - cobre `struct/enum/trait/type alias/const/static/union/macro`
  - cobre `impl` methods
  - cobre `mod`, `use`, `pub use`, alias e grouped imports
  - emite diagnósticos em arquivo com erro sintático
- Integração:
  - fixture Rust com `Cargo.toml`, módulos aninhados e `pub use` gera `filesScanned > 0`, `filesParsed > 0`, `symbolsExtracted > 0`, docs raw e relações internas.
  - help/tests de CLI passam a mostrar `rust`.
  - `make verify` continua sendo gate obrigatório.
- Aceite real em `~/dev/knowledge` após implementação:
  - `kb ingest codebase .resources/goclaw --topic goclaw`
  - `kb ingest codebase .resources/openclaw --topic openclaw`
  - `kb ingest codebase .resources/openfang --topic openfang`
  - `kb lint goclaw`, `kb lint openclaw` e `kb lint openfang` sem `dead-link` falso derivado de `raw/codebase/*`
  - `openfang` deixa de parecer um repo JS pequeno: o scan precisa detectar Rust explicitamente e processar os `.rs` do workspace, não só os 26 arquivos JS observados hoje
  - `kb topic info/list` continuam funcionando e os bridges/indexes atuais não regridem

## Assumptions

- O fix do lint é estrutural e localizado ao link graph; não altera renderização, snapshots raw nem semântica dos documentos manuais.
- O suporte Rust desta entrega busca paridade sintática máxima dentro da arquitetura atual, sem introduzir um novo modo semântico separado.
- Macros passam a ser símbolos de primeira classe quando declaradas; invocações entram apenas como relações quando forem resolvíveis de forma sintática.
