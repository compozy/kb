# Task Memory: task_11.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Consolidar a suíte de regressão da Fase 2 para Java cobrindo nested types, wildcard imports, política de ambiguidade, cenário enterprise com metadata multi-módulo, E2E de CLI/lint e validação de budget de performance com evidência reproduzível.

## Important Decisions
- Adicionado um cenário integrado `TestJavaAdapterPhase2EnterpriseScenarioRegression` para validar no mesmo caso: ambiguidade resolvida por metadata de módulo, nested class call, wildcard resolution e fallback previsível para wildcard ausente.
- O fixture Java de CLI foi evoluído para layout enterprise de 3 módulos (`shared-a`, `shared-b`, `app`) com `app/build.gradle` apontando para `shared-b`, mantendo asserts estáveis de resumo, artefatos e lint.
- A validação de `internal/generate` recebeu cenário de integração próprio (`TestGenerateIntegrationBuildsVaultFromJavaPhase2Workspace`) e o teste de budget passou a registrar baseline/java/overhead via `t.Logf` para evidência direta no output.

## Learnings
- A combinação de import ambíguo explícito + wildcard + nested type se mantém determinística quando a dependência de módulo (`app -> shared-b`) está disponível no fixture.
- Para evidência de benchmark em regressão contínua, registrar delta no teste de budget reduz ambiguidade de interpretação quando o teste passa.
- Cobertura do pacote `internal/adapter` permanece acima do limite exigido após os novos cenários (`80.7%` com `go test -tags integration ./internal/adapter -cover`).

## Files / Surfaces
- `internal/adapter/java_adapter_integration_test.go`
- `internal/cli/workflow_test_helpers_test.go`
- `internal/cli/workflow_integration_test.go`
- `internal/generate/generate_integration_test.go`
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md`

## Errors / Corrections
- Nenhuma falha bloqueante após as mudanças; testes focados passaram na primeira execução.
- `make verify` passou completo após atualização dos cenários.

## Ready for Next Run
- Evidência executada nesta task:
  - `go test -tags integration ./internal/adapter -run "TestJavaAdapterPhase2EnterpriseScenarioRegression"`
  - `go test ./internal/cli -run "TestWriteJavaMultiModuleCodebaseFixtureCreatesDeterministicLayout"`
  - `go test -tags integration ./internal/cli -run "TestCLIIntegrationScaffoldIngestJavaWorkspaceCodebase"`
  - `go test -tags integration ./internal/generate -run "TestGenerateIntegrationBuildsVaultFromJavaPhase2Workspace|TestGenerateIntegrationJavaIngestPerformanceBudget" -v`
  - `go test -tags integration ./internal/adapter -cover` (`80.7%`)
  - `make verify` (PASS)
