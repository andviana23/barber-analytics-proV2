# 🔎 Contexto — Módulo Financeiro

## Estado atual
- Banco de dados contém todas as tabelas e índices necessários (payables, receivables, compensações, fluxo diário, DRE).
- Nenhuma entidade/repositório/use case/handler foi criado para essas tabelas (vide `Tarefas/CONCLUIR`).
- Frontend não possui páginas/hooks para os fluxos avançados (payables/receivables, DRE, fluxo compensado, comissões).

## Dependências
- Base técnica + LGPD/Backup concluídos (pacotes 01 e 02).
- Módulo Financeiro desbloqueia Estoque (consumo automático) e Metas (progressão financeira).

## Referências
- `Tarefas/FINANCEIRO/*.md`
- `docs/10-calculos/*` (ticket médio, ponto de equilíbrio, margem, fluxo compensado)
- `DATABASE_MIGRATIONS_COMPLETED.md`
