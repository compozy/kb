# Portabilidade de `make verify` no Windows — Design

## Objetivo

Fazer `make verify` passar nativamente no Windows sem exigir WSL, Developer
Mode ou privilégios de criação de symlink, preservando o comportamento atual em
Linux, macOS e WSL.

## Causas-raiz

1. Testes com caminhos fictícios POSIX esperam a string original após a
   normalização legítima de `filepath` no Windows.
2. A criação de tópicos exige um symlink `AGENTS.md -> CLAUDE.md`; usuários
   comuns do Windows podem não possuir o privilégio necessário.
3. Fakes de `yt-dlp`, `ffmpeg` e QMD dependem de scripts POSIX e de `/bin/sh`.

## Desenho aprovado

- Normalizar expectativas de caminho nos testes com `filepath.Clean` ou
  caminhos de `t.TempDir()`.
- Manter symlink em plataformas Unix; quando o Windows negar especificamente o
  privilégio de symlink, gravar uma cópia regular de `CLAUDE.md` em `AGENTS.md`.
- Substituir fakes de shell por helpers de processo Go multiplataforma nos
  testes afetados.
- Adicionar testes de regressão para o fallback de `AGENTS.md` e os helpers;
  não pular testes e não adicionar requisitos de ambiente.

## Validação

Executar os pacotes afetados e depois `make verify` em Windows. A suite deve
continuar exigindo symlink em Unix e aceitar o fallback apenas no erro de
privilégio do Windows.
