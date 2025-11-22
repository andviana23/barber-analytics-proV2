# ✅ Checklist Dev — Financeiro

- [x] Domínios alinhados às migrations: status válidos e enums para payables/receivables/compensações ✅
- [x] Repositórios com filtros por tenant, status, vencimento e suporte a paginação/sorting ✅
- [x] Use cases de compensações (receivables) com jobs para marcar D+ ✅
- [ ] Use cases de recorrência (payables) - Pendente (requer definição de regras)
- [x] DRE calcula margens corretas usando categorias e valores ✅
- [x] Fluxo diário grava `fluxo_caixa_diario` e `compensacoes_bancarias` ✅
- [x] Endpoints versionados (`/api/v1/financial/*`) ✅
- [ ] RBAC (owner/manager/accountant) - Pendente (módulo permissões não implementado)
- [ ] Comissões geram PDFs/relatórios - Pendente (baixa prioridade)
- [ ] Dashboard com cache Redis - Pendente (aguardando frontend)
- [x] Métricas Prometheus para cron jobs ✅

**Status Implementação:**

- ✅ **Contas a Pagar:** 6 endpoints (Create, Get, List, Update, Delete, MarcarPagamento)
- ✅ **Contas a Receber:** 6 endpoints (Create, Get, List, Update, Delete, MarcarRecebimento)
- ✅ **Compensações:** 4 endpoints (Create, Get, List, Delete) + Cron Job
- ✅ **Fluxo de Caixa:** 2 endpoints (Get, List) + Cron Job
- ✅ **DRE:** 2 endpoints (Get, List) + Cron Job
- ✅ **Testes:** Unit tests + Integration tests + Smoke tests
