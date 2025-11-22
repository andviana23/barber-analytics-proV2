# 📌 Backlog — Precificação

## 🔴 Obrigatório
1. [ ] **T-PREC-001 — Simulador de Precificação** — ref. `01-precificacao-simulador.md`
   - CRUD de `precificacao_config` (margem desejada, markup, impostos, comissionamento padrão) — UNIQUE por tenant.
   - Endpoint `/pricing/simulations` para calcular preço sugerido usando custos, impostos, comissões e meta de margem.
   - Persistir histórico em `precificacao_simulacoes`; exportar resultado via API pública (quando habilitado).
   - UI de simulador com presets e salvar configurações.

## 🧭 Dependências cruzadas
- Usa custos de estoque/serviços e comissões; depende de Financeiro + Estoque.
- Metas fornecem margem alvo; usar como default quando existente.
