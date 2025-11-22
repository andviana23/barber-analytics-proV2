# VALTARIS v 1.0 — Componentes Críticos

Camada mínima que garante consistência “Cyber Luxury” e uso dos tokens. Light é padrão; todos respeitam CSS vars e MUI palette.

- **ThemeRegistry / ThemeProvider**  
  - Local: `app/theme-registry.tsx` (ver `02-ARCHITECTURE.md`).  
  - Responsável por injetar Emotion + MUI + CSS vars (`--valtaris-*`) e aplicar `.theme-dark`.

- **Layout Base (DashboardLayout)**  
  - Inclui Sidebar, Topbar, conteúdo com `max-width` opcional e grid de 24–32px.  
  - Usa `VCard` para seções e `VButton` para CTAs principais.

- **Formulários-chave (RHF + Zod + MUI)**  
  - FinancialModal, SubscriptionForms, cadastros gerais.  
  - Sempre com wrappers (`FormInput`, `FormSelect`) e estados de erro em `var(--valtaris-danger)`.  
  - Ver `04-PATTERNS.md` para standard de validação.

- **Hooks de Dados (TanStack Query)**  
  - `useReceitas`, `useDespesas`, `useSubscriptions`, `useDashboard` etc.  
  - Cache keys devem incluir `tenantId`; staleTime mínimo 5min para listas.  
  - Erros renderizados com `StatusBadge`/empty states padrão.

- **Lista da Vez / Scheduler (DayPilot)**  
  - Classe `daypilot-valtaris` + CSS vars (ver `02-ARCHITECTURE.md`).  
  - Nunca usar estilos inline stringificados; tema reage ao light/dark automaticamente.

- **Design System base**  
  - `VButton`, `VTextField/FormInput`, `VCard`, `VModal`, `StatusBadge`, `VTable`.  
  - Sempre referenciar tokens (`sx`/CSS vars), sem novas cores ou sombras fora do contrato.
