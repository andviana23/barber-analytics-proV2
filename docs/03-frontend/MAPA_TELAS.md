# VALTARIS v 1.0 — Mapa de Telas

Roteiro das telas principais e como o Design System se encaixa. Light é padrão; todas usam tokens do tema.

- **Login / Onboarding**  
  - Layout público `(auth)`; usa `VCard` + `VButton` primário; formulários com RHF+Zod.

- **Dashboard (KPIs)**  
  - Home de `(dashboard)`; cards com vidro (`VCard`), badges de status e gráficos tematizados.

- **Financeiro: Receitas / Despesas / Fluxo de Caixa**  
  - Tabelas `VTable`, filtros com `FormInput`, modais `VModal`.  
  - Agregações e filtros sempre com cache TanStack Query incluindo `tenantId`.

- **Financeiro: Contas a Pagar/Receber**  
  - Formulários com estados de erro em `var(--valtaris-danger)`; badges de status oficiais.

- **Assinaturas (planos, assinantes)**  
  - Forms e tabelas; usar `StatusBadge` para ciclo de cobrança; modais com vidro.

- **Lista da Vez / Agendamentos**  
  - DayPilot com classe `daypilot-valtaris` e CSS vars; não usar estilos inline.

- **Cadastros (clientes, profissionais, serviços, produtos, meios, cupons)**  
  - Reutilizar wrappers de formulário; tabelas com `compact mode` opcional.

- **Design System Preview / Storybook**  
  - Página demo para tokens + componentes; referenciar `01-FOUNDATIONS.md` e `03-COMPONENTS.md`.

Para rotas/URLs detalhadas, ver `frontend/app/(dashboard)` e `API_REFERENCE.md`.
