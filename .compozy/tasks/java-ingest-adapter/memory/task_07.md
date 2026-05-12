# Task Memory: task_07.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implementar resolução consistente para tipos Java aninhados (`Outer.Inner`) sem regressão em resolução top-level.
- Entregar cobertura unitária e de integração para nested/inner types com saída determinística.

## Important Decisions
- Tipos aninhados passam a ser modelados com nome qualificado no símbolo (ex.: `Outer.Inner`) para codificar ownership no próprio contrato existente sem alterar `models.SymbolNode`.
- Extração de declarações Java passa a ser recursiva dentro do corpo de tipos (`class/interface/enum/record`) para incluir nested declarations e métodos internos.
- Resolução de qualifier preserva cadeia completa (`Outer.Inner`) em vez de truncar para o último segmento.
- Resolvedor deep agora considera candidatos qualificados por import/local/package com expansão de prefixo (`Outer` + `Inner`) para resolver referências nested.

## Learnings
- A indexação local por nome simples só é segura quando o símbolo resolve para um único FQN no arquivo; em caso ambíguo, manter apenas nomes qualificados evita resolução incorreta.
- O benchmark de cobertura para `internal/adapter` depende de execução com `-tags integration` para manter o pacote acima do limiar de 80%.

## Files / Surfaces
- `internal/adapter/java_adapter.go`
- `internal/adapter/java_adapter_test.go`
- `internal/adapter/java_adapter_integration_test.go`
- `.compozy/tasks/java-ingest-adapter/memory/task_07.md`
- `.compozy/tasks/java-ingest-adapter/memory/MEMORY.md`
- `.compozy/tasks/java-ingest-adapter/task_07.md`
- `.compozy/tasks/java-ingest-adapter/_tasks.md`

## Errors / Corrections
- Teste `TestResolveJavaMethodInvocationFallbackParsing` falhou após preservar qualifier completo; assert foi corrigido para `com.example.Helper`.
- Cobertura unitária inicial ficou abaixo do alvo; adicionados testes para branches de helpers nested e resolução de call target, elevando `go test ./internal/adapter -cover` para 80.0%.

## Ready for Next Run
- Validar task seguinte considerando continuidade da resolução qualificada para wildcard imports e cenários ambíguos.
- Não foi criado commit automático (`--auto-commit=false`); diff permanece para revisão manual.
