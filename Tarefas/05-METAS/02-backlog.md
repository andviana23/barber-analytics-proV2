# 📌 Backlog — Metas

## 🔴 Obrigatórios
1. [ ] **T-MET-001 — Meta Geral do Mês** — ref. `01-meta-geral-mes.md`
   - CRUD metas em `metas_mensais` (UNIQUE tenant+mes_ano), status, origem (manual/automática), alertas de progresso.
   - Endpoint `/metas/mensal` + cards no dashboard.
2. [ ] **T-MET-002 — Meta por Barbeiro** — ref. `02-meta-por-barbeiro.md`
   - `metas_barbeiro`: metas individuais com ranking e progressão.
   - Endpoint `/metas/barbeiros` + páginas de ranking.
3. [ ] **T-MET-003 — Meta de Ticket Médio** — ref. `03-meta-ticket-medio.md`
   - `metas_ticket_medio`: metas globais e por barbeiro; cálculo usando receitas/atendimentos.
   - Endpoint `/metas/ticket-medio` + gráficos.
4. [ ] **T-MET-004 — Metas Automáticas** — ref. `04-metas-automaticas.md`
   - Engine para sugerir metas com base em faturamento mínimo, margem e histórico.
   - Jobs para gerar metas no início do mês; integração com dashboard.

## 🧭 Dependências cruzadas
- Usa dados de DRE/Fluxo/Receitas/Despesas; executar após Financeiro completo.
- Metas automáticas dependem de metas base (1-3) e cálculos de margem/PE.
