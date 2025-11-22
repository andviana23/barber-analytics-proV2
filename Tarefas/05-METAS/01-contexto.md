# 🔎 Contexto — Metas

## Estado atual
- Tabelas `metas_mensais`, `metas_barbeiro`, `metas_ticket_medio` já existem (migrations 032-034) com índices por tenant/mes_ano/barbeiro.
- Nenhum domínio/use case/endpoint foi criado; frontend não possui páginas/hooks de metas.

## Dependências
- Requer DRE/Fluxo para calcular progresso e alertas.
- Metas automáticas usam faturamento mínimo/margem (docs de cálculos) e custos/receitas reais.

## Referências
- `Tarefas/METAS/01-meta-geral-mes.md`
- `Tarefas/METAS/02-meta-por-barbeiro.md`
- `Tarefas/METAS/03-meta-ticket-medio.md`
- `Tarefas/METAS/04-metas-automaticas.md`
