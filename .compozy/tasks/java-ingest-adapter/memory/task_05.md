# Task Memory: task_05.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implementar resolucao profunda Java (deep-first) com fallback sintatico automatico, diagnostico estruturado de fallback e ordenacao deterministica, incluindo cobertura de testes unitarios e de integracao para os cenarios exigidos da task_05.

## Important Decisions
- Introduzida abstracao interna de resolver no `java_adapter` com duas estrategias: `javaDeepResolver` (semantic) e `javaSyntacticResolver` (fallback), sem churn de API entre pacotes.
- O fallback passou a ser acionado somente para alvos que falharam na resolucao profunda, evitando duplicacao semantic+syntactic para o mesmo vinculo quando deep ja resolveu.
- O diagnostico `JAVA_RESOLUTION_FALLBACK` foi agregado por arquivo com `severity=warning` e `stage=parse`, incluindo detalhes ordenados deterministicamente por alvo/razao.

## Learnings
- O baseline de cobertura de `internal/adapter` estava abaixo da meta da task; foi necessario ampliar testes de helper e ramos internos do Java adapter para atingir `80.0%` sem expandir para mudancas fora do escopo Java.
- A ordenacao explicita de diagnosticos e relacoes e essencial para manter igualdade entre execucoes repetidas em fixtures de integracao.

## Files / Surfaces
- `internal/adapter/java_adapter.go`
- `internal/adapter/java_adapter_test.go`
- `internal/adapter/java_adapter_integration_test.go`
- `.compozy/tasks/java-ingest-adapter/task_05.md`
- `.compozy/tasks/java-ingest-adapter/_tasks.md`
- `.compozy/tasks/java-ingest-adapter/memory/task_05.md`
- `.compozy/tasks/java-ingest-adapter/memory/MEMORY.md`

## Errors / Corrections
- Durante patch incremental do `java_adapter.go`, um artefato literal `*** End Patch` ficou no arquivo e foi removido na sequencia.
- A primeira versao do fallback aplicava resolucao sintatica para todos os alvos; corrigido para aplicar somente aos alvos unresolved da fase deep.

## Ready for Next Run
- Task pronta para revisao manual com diff preparado (auto-commit desabilitado).
- Verificacoes executadas com sucesso: testes Java unit/integration, cobertura `internal/adapter` em `80.0%` e `make verify` completo.
