# ✅ Checklist Dev — Estoque

- [ ] Criar tabelas de movimentação/inventário com `tenant_id`, índices por produto/data/tipo.
- [ ] Regras de saldo: não permitir estoque negativo; atualizar custo médio quando aplicável.
- [ ] Endpoints `/stock/*` com RBAC (owner/manager para escrita; employee leitura) e validação de fornecedor/produto do tenant.
- [ ] Integração opcional: flag para criar `contas_a_pagar` em entradas.
- [ ] Job de estoque mínimo com alertas e log em `cron_run_logs`.
- [ ] Relatório Curva ABC calculado server-side com paginação/export.
- [ ] Auditoria registrando entradas/saídas/ajustes.
- [ ] Métricas e logs para movimentações (latência, erros, saldo atualizado).
