# 📌 Backlog — Financeiro

## 🔴 Obrigatórios
1. [ ] **T-FIN-001 — Contas a Pagar** — ref. `Tarefas/FINANCEIRO/03-contas-a-pagar.md`
   - Implementar domínios/repos/use cases + endpoints `/financial/payables` (CRUD, recorrência, notificações D-5/D-1/D0) usando `contas_a_pagar`.
   - Upload de comprovante seguro; status `ABERTO/PAGO/ATRASADO` com transições validadas.
2. [ ] **T-FIN-002 — Contas a Receber** — ref. `Tarefas/FINANCEIRO/04-contas-a-receber.md`
   - Modelar `contas_a_receber` (origem assinatura/serviço/outro), sync manual com Asaas, conciliação e inadimplência.
   - Endpoints `/financial/receivables` + notificações de atraso.
3. [ ] **T-FIN-003 — Fluxo de Caixa Compensado** — ref. `Tarefas/FINANCEIRO/07-fluxo-caixa-compensado.md`
   - Use cases para gerar `fluxo_caixa_diario` e `compensacoes_bancarias` (D+ configurável em `meios_pagamento.d_mais`).
   - Endpoint `/financial/cashflow/compensado` com projeções D+N e compensações.
4. [ ] **T-FIN-004 — Comissões Automáticas** — ref. `Tarefas/FINANCEIRO/05-comissoes-automaticas.md`
   - Engine de cálculo (fixo/percentual/degrau) sobre faturas recebidas; geração de PDFs/relatórios.
   - Integração com `barber_commissions` e dashboard.
5. [ ] **T-FIN-005 — DRE Completo** — ref. `Tarefas/FINANCEIRO/02-dre.md` e `06-dre-completo.md`
   - Agregação mensal em `dre_mensal` usando `categorias.tipo_custo` e `receitas.subtipo`.
   - Endpoints de comparação M/M e exportação PDF.
6. [ ] **T-FIN-006 — Dashboard Financeiro** — ref. `Tarefas/FINANCEIRO/01-dashboard-financeiro.md`
   - Endpoint agregado + UI (metas, PE, fluxo, DRE) com cache Redis e invalidation.

## 🧭 Dependências cruzadas
- Fluxo compensado depende de payables/receivables + `meios_pagamento.d_mais`.
- DRE usa dados de payables/receivables + categorias com `tipo_custo` e `receitas.subtipo`.
- Dashboard consome resultados de T-FIN-001..005; executar por último.
