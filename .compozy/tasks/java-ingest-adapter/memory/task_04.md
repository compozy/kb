# Task Memory: task_04.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Registrar `adapter.JavaAdapter{}` na orquestração de `internal/generate` sem alterar contratos de estágio/CLI.
- Cobrir seleção determinística em workspace misto e validar help CLI com suporte Java exposto.

## Important Decisions
- A ordem dos adapters no runner foi mantida como `TS -> Go -> Rust -> Java` para preservar determinismo e evitar regressão de seleção nos idiomas já suportados.
- O teste de integração deste task valida fluxo `DryRun` em workspace misto (`go` + `java`) para comprovar detecção e seleção de adapter sem depender da etapa de escrita.

## Learnings
- `supportedCodebaseLanguagesHelp()` já deriva de `models.SupportedLanguageNames()`; a validação explícita de `java` nos testes de help impede regressão silenciosa caso o texto formatado mude.
- A seleção determinística é definida pela ordem da lista de adapters registrada no runner (não pela ordem do scan em runtime).

## Files / Surfaces
- `internal/generate/generate.go`
- `internal/generate/generate_test.go`
- `internal/generate/generate_integration_test.go`
- `internal/cli/generate_test.go`
- `internal/cli/ingest_test.go`

## Errors / Corrections
- Nenhum erro de implementação; ajustes passaram em testes alvo e `make verify`.

## Ready for Next Run
- Task tracking (`task_04.md` e `_tasks.md`) pode ser marcado como concluído com base em evidências de testes alvo + `make verify` verde.
