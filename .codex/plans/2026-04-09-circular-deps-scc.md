# Escalar Detecção de Dependências Circulares

## Summary
- Substituir a enumeração exponencial de ciclos simples por SCCs determinísticos e lineares no estágio `metrics`.
- Preservar `HasCircularDependency` por arquivo e passar `CircularDependencies` a representar grupos cíclicos.
- Alinhar `inspect circular-deps` à referência TypeScript: listar arquivos participantes, com fallback SCC para vaults antigos.

## Key Changes
- Centralizar a lógica de SCC em um helper reutilizável.
- Atualizar `metrics` para usar SCCs e marcar arquivos participantes.
- Atualizar wiki/article para descrever grupos cíclicos em vez de loops simples.
- Atualizar `inspect circular-deps` para retornar arquivos participantes.
- Remover a duplicação do algoritmo atual entre `metrics` e `inspect`.

## Public Interfaces
- `models.MetricsResult.CircularDependencies` passa a significar grupos cíclicos/SCCs.
- `kodebase inspect circular-deps` deixa de emitir `{cycle, files}` e passa a emitir arquivos participantes.

## Test Plan
- Caso simples `A -> B -> C -> A` continua retornando um único grupo.
- Casos com ciclos sobrepostos no mesmo SCC retornam um grupo único.
- Múltiplos SCCs independentes mantêm ordenação determinística.
- `inspect circular-deps` passa a esperar arquivos participantes.
- Fallback para vault antigo sem `has_circular_dependency`.
- Verificação final com `make verify` e reprodução em repositório grande.

## Assumptions
- Saída aprovada: `generate` reporta SCCs; `inspect circular-deps` lista arquivos participantes.
- Ordenação lexicográfica estável para grupos e arquivos.
- Self-loop de um único arquivo permanece fora do escopo salvo exigência explícita de teste.
