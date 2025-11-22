# 🔎 Contexto — Precificação

## Estado atual
- Tabelas `precificacao_config` e `precificacao_simulacoes` existem com índices por tenant/item/criado_em.
- Não há domínio/use case/endpoint/hook implementado; simulador é inexistente no frontend.

## Dependências
- Custo médio de produtos/serviços (estoque + financeiro) e comissões automáticas prontos para cálculos.
- Metas/margens definidas para sugerir preços.

## Referência
- `Tarefas/PRECIFICACAO/01-precificacao-simulador.md`
- `docs/10-calculos/*` (margem, markup, ponto de equilíbrio)
