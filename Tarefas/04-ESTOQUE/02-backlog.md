# 📌 Backlog — Estoque

## 🔴 Obrigatórios
1. [ ] **T-EST-001 — Entrada de Estoque** — ref. `01-entrada.md`
   - Tabela de movimentações com tipo `ENTRADA`, vínculos com fornecedor e custos; opção de gerar `contas_a_pagar`.
   - Endpoint `/stock/entries` + UI de entrada.
2. [ ] **T-EST-002 — Saída de Estoque** — ref. `02-saida.md`
   - Movimentação `SAIDA` com motivos e validação de saldo; bloqueio de estoque negativo.
   - Endpoint `/stock/outputs` + UI.
3. [ ] **T-EST-003 — Consumo Automático por Serviço** — ref. `03-consumo-automatico.md`
   - Ficha técnica por serviço; baixa automática em atendimento/assinatura.
   - Integração com financeiro opcional (custo por serviço).
4. [ ] **T-EST-004 — Inventário** — ref. `04-inventario.md`
   - Contagem física, divergências e ajustes de saldo.
   - Auditoria e relatórios de ajustes.
5. [ ] **T-EST-005 — Estoque Mínimo & Alertas** — ref. `06-estoque-minimo.md`
   - Configuração de `estoque_minimo` por item; job de alerta e sugestão de compra.
6. [ ] **T-EST-006 — Curva ABC** — ref. `05-curva-abc.md`
   - Relatório Pareto (A/B/C) baseado em consumo e valor; exportação.

## 🧭 Dependências cruzadas
- Entradas podem gerar payables → requer módulo financeiro ativo.
- Consumo automático depende de serviços configurados e operações de saída/entrada.
- Alertas de estoque mínimo dependem de inventário atualizado.
