# 🟦 FASE 4 — Frontend 2.0 (Next.js)

**Objetivo:** Frontend Next.js apontando para novo backend Go
**Duração:** 14-28 dias
**Dependências:** ✅ Fase 2 completa (APIs básicas disponíveis)
**Sprint:** Sprint 4-6 (PARALELO à Fase 3)

---

## 📊 Progresso Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 4: FRONTEND 2.0 (NEXT.JS)                             │
├─────────────────────────────────────────────────────────────┤
│  Progresso:  ██████████████████████████  96.875% (15.5/16)  │
│  Status:     🟡 Em Andamento (T-QA-003 75% completo)        │
│  Prioridade: 🔴 ALTA                                        │
│  Estimativa: 67 horas (62h concluídas + 6h T-QA-003)        │
│  Sprint:     Sprint 4-6 (paralelo Fase 3)                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎨 Alinhamento com o Designer System

- Use `docs/Designer-System.md` v2.0.0 como fonte única de tokens, componentes e diretrizes (cores, espaçamentos, tipografia Inter + JetBrains Mono, sombras, bordas).
- Stack oficial: Next.js 15 (App Router) + React 19 + MUI 5 + TanStack Query + DayPilot Scheduler (Gantt futuro), com TanStack encarregando fetch/cache.
- Tema compartilhado: `useThemeStore` (Zustand com persistência), `ThemeToggle`, `useThemeSyncDom`, variáveis CSS (`data-theme` + `color-scheme`) e detecção de `prefers-color-scheme` preservando o modo light como padrão.
- Tokens devem nortear todos os estilos (`app/theme/tokens.ts` + `components`), nunca hardcode de cores/spacing, e o uso do Tailwind deve ficar limitado a utilitários ou evitar misturar com `sx`.
- Scheduler DayPilot e demais componentes visuais devem consumir os tokens (background, bordas, texto, eventos) e responder aos breakpoints (Mobile → Day, Tablet → Week, Desktop → WorkWeek).
- Acessibilidade WCAG 2.1 AA é obrigatória: contraste ≥ 4.5:1, foco visível 2px, labels/ARIA, navegação por teclado e respeito a `prefers-reduced-motion`.

## ✅ Checklist de Tarefas

### ✅ T-FE-001 — Setup Next.js v2
- **Responsável:** Frontend Lead
- **Prioridade:** 🔴 Alta
- **Estimativa:** 3h
- **Sprint:** Sprint 4
- **Status:** ✅ Concluído
- **Deliverable:** Next.js 15 + TypeScript + App Router + MUI 5 core com tokens

#### Critérios de Aceitação
- [x] `create-next-app` executado com TypeScript (Next.js 16)
- [x] `@next/font` (Inter + JetBrains Mono) configurado
- [x] `app/theme/tokens.ts` expõe cores, espaçamentos e sombras do Designer System
- [x] `app/layout.tsx` envolve `ThemeProvider` + `CssBaseline` + `QueryClientProvider`
- [x] `.env.local.example` criado
- [x] Config: Image optimization, compression
- [x] App Router estruturado
- [x] Testado: `pnpm run dev` → http://localhost:3000

#### Arquivos Criados
- `app/fonts.ts` - Configuração Inter + JetBrains Mono
- `app/theme/tokens.ts` - Tokens completos do Design System
- `app/theme/index.ts` - Tema MUI com tokens integrados
- `app/providers.tsx` - QueryClientProvider + ThemeProvider
- `app/layout.tsx` - Root layout atualizado
- `.env.local.example` + `.env.local`
- `next.config.ts` - Otimizações configuradas

---

### ✅ T-FE-002 — API Client & Interceptors
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 4
- **Status:** ✅ Concluído
- **Deliverable:** Axios instance com retry + JWT interceptor

#### Critérios de Aceitação
- [x] Axios instance configurada
- [x] Base URL: `NEXT_PUBLIC_API_URL`
- [x] Request interceptor (add Authorization header)
- [x] Response interceptor (401 handling → refresh token)
- [x] Retry logic (3 tentativas para 5xx com backoff exponencial)

#### Arquivos Criados
- `app/lib/api/client.ts` - API client completo com interceptors
- `app/lib/hooks/useApi.ts` - Hook para integração React Query
- `app/lib/utils/storage.ts` - Utilities para localStorage
- `app/lib/types/auth.ts` - Tipos TypeScript para auth

#### Funcionalidades Implementadas
- Request interceptor adiciona `Authorization: Bearer {token}`
- Response interceptor detecta 401 e faz refresh automático do token
- Fila de requests durante refresh para evitar race conditions
- Retry com axios-retry: 3 tentativas, backoff exponencial
- Logout automático quando refresh falha
- Helpers tipados: `api.get<T>()`, `api.post<T>()`, etc.

---

### ✅ T-FE-003 — Auth & Protected Routes
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 4
- **Status:** ✅ Concluído
- **Deliverable:** Login page + middleware de proteção

#### Critérios de Aceitação
- [x] Login page criada: `app/(auth)/login/page.tsx`
- [x] Token storage (cookies com helper client/server)
- [x] Middleware para rotas protegidas
- [x] Logout funcional
- [x] Redirect após login: `/dashboard`

#### Arquivos Criados
- `app/lib/types/auth.ts` - Tipos completos (User, Tenant, LoginCredentials, etc.)
- `app/lib/auth/tokens.ts` - Helper para cookies (server + client side)
- `app/lib/contexts/AuthContext.tsx` - Context com login/logout/user state
- `app/lib/hooks/useLogout.ts` - Hook para logout
- `app/(auth)/login/page.tsx` - Login page completo com React Hook Form + Yup
- `app/providers.tsx` - Atualizado com AuthProvider

---

### ✅ T-FE-004 — Layout & Navigation
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 4
- **Status:** ✅ Concluído
- **Deliverable:** Layout principal com sidebar + topbar

#### Critérios de Aceitação
- [x] Root layout com providers (TanStack Query + Auth)
- [x] Dashboard layout (sidebar, topbar)
- [x] Tenant info exibido no topbar
- [x] Navegação: Dashboard, Receitas, Despesas, Assinaturas, Fluxo Caixa
- [x] Responsivo (mobile + desktop)

#### Arquivos Criados
- `app/components/layout/Sidebar.tsx` - Menu lateral com navegação e collapse
- `app/components/layout/Topbar.tsx` - Barra superior com tenant info e menu usuário
- `app/components/layout/DashboardLayout.tsx` - Layout wrapper
- `app/(private)/layout.tsx` - Layout para rotas privadas
- `app/(private)/dashboard/page.tsx` - Dashboard com KPI cards
- `app/(private)/receitas/page.tsx` - Página de receitas
- `app/(private)/despesas/page.tsx` - Página de despesas
- `app/(private)/assinaturas/page.tsx` - Página de assinaturas
- `app/(private)/fluxo-caixa/page.tsx` - Página de fluxo de caixa

---

### ✅ T-FE-005 — Dashboard page
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 5
- **Status:** ✅ Concluído
- **Deliverable:** Dashboard com KPIs + charts

#### Critérios de Aceitação
- [x] KPI cards: Receita, Despesa, Saldo, Assinantes
- [x] Charts (Recharts):
  - [x] Receita mensal (bar chart)
  - [x] Despesas por categoria (pie chart)
- [x] Recent activity feed
- [x] Período selecionável (mês atual, último mês, etc)

#### Arquivos Criados
- `app/lib/utils/format.ts` - Formatação de currency, datas, números
- `app/lib/types/financial.ts` - Tipos completos para entidades financeiras
- `app/lib/hooks/useDashboard.ts` - Hooks para métricas e dados do dashboard
- `app/components/dashboard/DashboardKPI.tsx` - Card KPI reutilizável
- `app/components/dashboard/RevenueBarChart.tsx` - Gráfico de barras (receitas)
- `app/components/dashboard/ExpensePieChart.tsx` - Gráfico de pizza (despesas)
- `app/components/dashboard/RecentActivity.tsx` - Feed de atividades
- `app/components/dashboard/PeriodSelector.tsx` - Seletor de período

---

### ✅ T-FE-006 — Receitas & Despesas pages
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 8h
- **Sprint:** Sprint 5
- **Status:** ✅ Concluído
- **Deliverable:** Páginas CRUD para receitas e despesas

#### Critérios de Aceitação
- [x] Tabelas com paginação (10/25/50 itens por página)
- [x] Filtros: período, categoria, status, busca
- [x] Modal criar/editar com validação Zod
- [x] Delete com confirmação (ConfirmDialog)
- [x] Formatação: R$ 1.000,00
- [x] Loading states e error handling
- [x] Toast notifications (notistack)

#### Arquivos Criados
- `app/components/financial/FinancialTable.tsx` - Tabela reutilizável com paginação
- `app/components/financial/FinancialFilters.tsx` - Filtros reutilizáveis
- `app/components/financial/FinancialModal.tsx` - Modal com React Hook Form + Zod
- `app/components/common/ConfirmDialog.tsx` - Diálogo de confirmação
- `app/(private)/receitas/page.tsx` - Página completa de Receitas
- `app/(private)/despesas/page.tsx` - Página completa de Despesas
- `app/lib/hooks/useReceitas.ts` - CRUD hooks com notificações
- `app/lib/hooks/useDespesas.ts` - CRUD hooks com notificações
- `app/lib/hooks/useCategorias.ts` - Hook para categorias
- `app/providers.tsx` - Integração SnackbarProvider (notistack)

---

### ✅ T-FE-007 — Assinaturas page
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 5
- **Status:** ✅ Concluído
- **Deliverable:** Página de assinaturas

#### Critérios de Aceitação
- [x] Lista de planos de assinatura
- [x] Lista de assinantes (por tenant)
- [x] Status badge (ativo, cancelado, pendente)
- [x] Ação: Cancelar assinatura (com confirmação)
- [x] Filtro: status, plano

#### Arquivos Criados
- `app/lib/types/subscription.ts` - Tipos para planos e assinaturas
- `app/lib/hooks/usePlans.ts` - Hooks para planos
- `app/lib/hooks/useSubscriptions.ts` - Hooks CRUD para assinaturas
- `app/components/subscriptions/SubscriptionStatusBadge.tsx` - Badge de status
- `app/components/subscriptions/SubscriptionFilters.tsx` - Filtros reutilizáveis
- `app/components/subscriptions/SubscriptionsTable.tsx` - Tabela com paginação
- `app/components/subscriptions/PlansGrid.tsx` - Grid de cards de planos
- `app/(private)/assinaturas/page.tsx` - Página completa com stats

---

### ✅ T-FE-008 — Fluxo de Caixa page
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 4h
- **Sprint:** Sprint 5
- **Status:** ✅ Concluído
- **Deliverable:** Página de fluxo de caixa

#### Critérios de Aceitação
- [x] Seletor de período (início/fim)
- [x] Tabela com saldo diário:
  - [x] Data, Entradas, Saídas, Saldo
- [x] Card resumo: Saldo inicial, Total entradas, Total saídas, Saldo final
- [x] Gráfico de linha (saldo ao longo do tempo)

#### Arquivos Criados
- `app/lib/types/cashflow.ts` - Tipos para fluxo de caixa
- `app/lib/hooks/useCashflow.ts` - Hook para buscar dados de fluxo
- `app/components/cashflow/CashflowSummary.tsx` - Cards de resumo
- `app/components/cashflow/CashflowChart.tsx` - Gráfico de linha (Recharts)
- `app/components/cashflow/CashflowTable.tsx` - Tabela detalhamento diário
- `app/(private)/fluxo-caixa/page.tsx` - Página completa com seletor período

---

### ✅ T-FE-009 — React Hooks customizados
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 5
- **Status:** ✅ Concluído
- **Deliverable:** Hooks reutilizáveis

#### Critérios de Aceitação
- [x] `useAuth()` - estado auth + login/logout (AuthContext)
- [x] `useTenant()` - tenant atual (via AuthContext + user.tenant)
- [x] `useReceitas()` - lista com paginação
- [x] `useDespesas()` - lista com paginação
- [x] `useAssinaturas()` - lista (useSubscriptions)
- [x] `useFluxoCaixa()` - cálculo de fluxo (useCashflow)

#### Hooks Implementados
- `useAuth()` - Exportado de AuthContext, expõe {user, isAuthenticated, login, logout, isLoading}
- `useReceitas()` - Lista com filtros, CRUD completo + toast notifications
- `useDespesas()` - Lista com filtros, CRUD completo + toast notifications
- `useCategorias(type)` - Lista categorias por tipo (receita/despesa)
- `useSubscriptions(filters)` - Lista assinaturas com paginação e filtros
- `usePlans()` - Lista planos de assinatura
- `useCashflow(filters)` - Busca dados de fluxo de caixa por período
- `useDashboard()` - Métricas e dados do dashboard
- `useLogout()` - Wrapper para logout com invalidação de queries

---

### ✅ T-FE-010 — Forms com React Hook Form + Zod
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído
- **Deliverable:** Forms validados

#### Critérios de Aceitação
- [x] `ReceitaForm` com validação Zod (via FinancialModal)
- [x] `DespesaForm` com validação Zod (via FinancialModal)
- [x] `LoginForm` com validação Zod (Yup - compatível)
- [x] Validação em tempo real (mode: 'onChange')
- [x] Mensagens de erro customizadas
- [x] Loading states e disabled buttons

#### Schemas Criados
- `app/lib/schemas/auth.schema.ts` - loginSchema
- `app/lib/schemas/financial.schema.ts` - receitaSchema, despesaSchema
- `app/lib/schemas/index.ts` - Exports centralizados

#### Componentes de Formulário
- `FinancialModal` - Form unificado para receitas/despesas com React Hook Form + Zod
  - Validação em tempo real
  - NumericFormat para moeda (R$ 1.000,00)
  - Controller para campos MUI
  - Error handling com helperText
  - Loading states e botões disabled
- `LoginForm` (login/page.tsx) - Form com validação Yup
  - Campos email/password
  - Validação em tempo real
  - Loading spinner
  - Error messages

#### Validações Implementadas
- Descrição: min 3 chars, max 255 chars
- Valor: número positivo, min 0.01
- Data: formato YYYY-MM-DD obrigatório
- Categoria: UUID obrigatório
- Status: enum validado (pending/received para receita, pending/paid para despesa)
- Email: formato válido obrigatório
- Senha: min 6 caracteres

---

### ✅ T-FE-011 — Componentes do Design System (MUI)
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 4h
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído
- **Deliverable:** Biblioteca base de componentes MUI consumindo tokens e estados acessíveis

#### Critérios de Aceitação
- [x] Button (variants: primary, secondary, outline, ghost) com estados `loading` e `disabled`
- [x] Input (text, number, date) com `FormHelperText`, labels vinculados e mensagens de erro aria-live
- [x] Table (paginação, sorting, zebra, focus-visible) utilizando tokens para backgrounds e borders
- [x] Dialog/Modal (DialogTitle, DialogContent, DialogActions) consumindo tokens e respeitando `prefers-reduced-motion`
- [x] Select/Dropdown com `aria-expanded` e foco visível
- [x] Badge (status colors) + ícones com contraste WCAG 2.1 AA
- [x] Todos os componentes usam cores, espaçamentos, bordas e sombras de `app/theme/tokens.ts` e expõem props para `aria`/focus states

#### Arquivos Criados
- `app/components/ui/Button.tsx` - LoadingButton com 4 variants (primary, secondary, outline, ghost), estados loading/disabled, focus visível 2px
- `app/components/ui/InputField.tsx` - TextField wrapper com label vinculado, FormHelperText, aria-invalid, aria-describedby, contraste WCAG AA
- `app/components/ui/DataTable.tsx` - Table com paginação (10/25/50), sorting opcional, zebra stripes, loading skeletons, empty state, focus visível
- `app/components/ui/Dialog.tsx` - DialogRoot, DialogTitle, DialogContent, DialogActions com tokens, prefers-reduced-motion, responsivo (fullScreen mobile)
- `app/components/ui/SelectField.tsx` - Select wrapper com label, aria-expanded, focus visível, options com disabled
- `app/components/ui/StatusBadge.tsx` - Chip com mapeamento de cores (ATIVO/verde, PENDENTE/laranja, CANCELADO/vermelho, INATIVO/cinza), ícones MUI, contraste AA
- `app/components/ui/index.ts` - Exports centralizados de todos os componentes

---

### ✅ T-FE-012 — Formatting & Utils
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 3h
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído
- **Deliverable:** Funções utilitárias

#### Critérios de Aceitação
- [x] `formatCurrency(1000)` → "R$ 1.000,00"
- [x] `formatDate()` → "14/11/2025"
- [x] `formatPercent()` → "12.5%"
- [x] `storage` (localStorage wrapper com tipos)
- [x] `cn()` (className utility)

#### Arquivos Criados/Atualizados
- `app/lib/utils/format.ts` - Já existia com formatCurrency, formatDate, formatDateTime, formatRelativeDate, formatPercent, formatNumber, parseCurrency, truncate, formatCNPJ, formatCPF, formatPhone (todos usando Intl.NumberFormat e date-fns com locale pt-BR)
- `app/lib/utils/cn.ts` - Utility para combinar classNames condicionalmente (similar ao clsx), suporta strings, objetos condicionais, arrays
- `app/lib/utils/storage.ts` - Atualizado com métodos genéricos tipados: get<T>(), set<T>(), remove(), clear(), has() + métodos específicos para tokens (getAccessToken, setTokens, clearAll), SSR-safe
- `app/lib/utils/index.ts` - Exports centralizados de todas as utilities

---

### ✅ T-FE-013 — Theme System & Dark Mode
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 3h
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído
- **Deliverable:** Theme store + Toggle + CSS vars alinhados ao Designer System

#### Critérios de Aceitação
- [x] `lib/store/themeStore.ts` com Zustand + persistência em `localStorage`, mode default `light`
- [x] `components/ThemeToggle.tsx` usando MUI `IconButton`, `Brightness4/Brightness7` e `aria-label`
- [x] `hooks/useSystemTheme.ts` aplicando `prefers-color-scheme` apenas quando não houver tema salvo
- [x] `hooks/useThemeSyncDom.ts` definindo `data-theme` e `color-scheme` no `document.documentElement`
- [x] `styles/theme-variables.css` com variáveis CSS (light/dark) derivadas dos tokens
- [x] Tema sincronizado com os tokens (`tokens.colorPrimary`, `tokens.bgPaper`, etc.) para garantir 4.5:1 de contraste

#### Arquivos Criados/Atualizados
- `app/lib/store/themeStore.ts` - Zustand store com persist (localStorage chave 'bap-theme'), mode padrão 'light', setMode(), toggleMode()
- `app/lib/hooks/useSystemTheme.ts` - Detecta prefers-color-scheme, lazy initialization, observa mudanças em tempo real, SSR-safe
- `app/lib/hooks/useThemeSyncDom.ts` - Sincroniza mode com DOM (data-theme, color-scheme), cleanup automático, SSR-safe
- `app/styles/theme-variables.css` - CSS vars para light/dark mode, contraste WCAG AA >= 4.5:1, prefers-reduced-motion, 160 linhas
- `app/components/common/ThemeToggle.tsx` - IconButton com Brightness4/7, Tooltip, aria-label acessível, rotação 180° no hover, focus visível
- `app/components/providers/AppThemeProvider.tsx` - Wrapper que integra store + MUI ThemeProvider + DOM sync + system preference, hidratação automática
- `app/providers.tsx` - Atualizado para usar AppThemeProvider (substitui ThemeProvider + CssBaseline direto)
- `app/components/layout/Topbar.tsx` - ThemeToggle integrado ao lado do menu de usuário
- `app/globals.css` - Import de theme-variables.css
- `app/theme/index.ts` - Já tinha createAppTheme(mode) para gerar tema MUI dinâmico

#### Features Implementadas
- **Persistência**: Tema salvo em localStorage, sobrevive a reloads
- **System Preference**: Auto-detecta prefers-color-scheme na primeira carga se não houver tema salvo
- **DOM Sync**: data-theme e color-scheme aplicados automaticamente
- **MUI Integration**: Tema MUI dinâmico baseado em mode (light/dark)
- **CSS Variables**: 40+ variáveis CSS para light e dark, contraste validado
- **Acessibilidade**: Focus visível 2px, aria-label, prefers-reduced-motion
- **Animações**: Rotação 180° no hover do toggle, transições suaves (300ms)

---

### ✅ T-FE-014 — DayPilot Scheduler Integration
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 4h
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído
- **Deliverable:** Scheduler DayPilot responsivo e tematizado

#### Critérios de Aceitação
- [x] `components/Scheduler.tsx` usa `DayPilot.Scheduler`, `useQuery`/`useMutation` (TanStack Query) e `useThemeStore`
- [x] `getDayPilotTheme(mode)` injeta CSS usando tokens (`bgDefault`, `textPrimary`, `borderDefault`, `colorPrimary`, etc.)
- [x] `components/ResponsiveScheduler` usa `useMediaQuery` para ajustar viewType/cellHeight (Mobile: Day, Tablet: Week, Desktop: WorkWeek)
- [x] Eventos e resources revalidam via `useQueryClient` após mutations
- [x] DayPilot grid/cells header/selected/event states respeitam tokens de sombras/bordas e respondem a modo escuro

#### Implementação Realizada
**Arquivos criados:**
- ✅ `app/lib/types/scheduler.ts` — Tipos TypeScript (SchedulerResource, SchedulerEvent, SchedulerConfig, SchedulerCallbacks)
- ✅ `app/lib/daypilot/theme.ts` — Utility getDayPilotTheme(mode) + applyDayPilotTheme() com injeção CSS
- ✅ `app/components/scheduler/Scheduler.tsx` — Componente base com DayPilot.Scheduler + TanStack Query
- ✅ `app/components/scheduler/ResponsiveScheduler.tsx` — Wrapper responsivo com useMediaQuery
- ✅ `app/components/scheduler/index.ts` — Barrel exports
- ✅ `app/lib/hooks/useSchedulerEvents.ts` — Hook opcional para integração Query
- ✅ `app/(private)/agendamentos/page.tsx` — Página demo com eventos de exemplo

**Features implementadas:**
- 🎨 Tema light/dark dinâmico via useThemeStore
- 📱 Responsivo: Day (mobile), Week (tablet), WorkWeek (desktop)
- 🔄 Revalidação automática com TanStack Query após mutations
- 🎯 Drag & drop de eventos (move + resize)
- ➕ Criação por seleção de range
- 🖱️ Click em eventos
- 💅 Estilização 100% via tokens (cores, sombras, bordas, tipografia)
- ⚡ Custom CSS para hover effects, transitions, scrollbar
- 📦 TypeScript strict mode compliant

**Pacotes instalados:**
- @daypilot/daypilot-lite-react v4.8.1 (free version)

---

### ✅ T-FE-015 — Acessibilidade & Auditoria de Componentes
- **Responsável:** Frontend
- **Prioridade:** 🔴 Alta (elevado de Média)
- **Estimativa:** 5h (100% completo)
- **Sprint:** Sprint 6
- **Status:** ✅ Concluído e Validado
- **Deliverable:** Componentes acessíveis WCAG 2.1 AA + testes a11y + documentação

#### Critérios de Aceitação
- [x] `components/design-system/AccessibleInput.tsx` com `label`, `aria-describedby`, `aria-invalid`, `role="alert"`, `aria-live="polite"`
- [x] Foco 2px outline usando `tokens.focus` aplicado em Button/Input/Modal via `:focus-visible`
- [x] Contraste text.primary/secondary ≥ 4.5:1 validado em light/dark mode
- [x] `prefers-reduced-motion: reduce` aplicado globalmente em CSS + hook React
- [x] `jest-axe` instalado com testes para Button (8), AccessibleInput (10), Modal (7) — **25/25 passed, zero violations**
- [x] Seção "Acessibilidade" documentada em `docs/Designer-System.md` (7 seções completas: contraste, foco, ARIA, keyboard nav, reduced motion, componente acessível, testes)

#### Resultados da Validação

**Testes Automatizados:**
```bash
$ pnpm test:a11y

Test Suites: 3 passed, 3 total
Tests:       25 passed, 25 total
Time:        3.568 s

✅ Button.a11y.test.tsx: 8 tests passed (default, loading, disabled, 4 variants, with icon)
✅ AccessibleInput.a11y.test.tsx: 10 tests passed (5 axe + 5 ARIA)
✅ Modal.a11y.test.tsx: 7 tests passed (3 axe + 3 ARIA + 1 keyboard)
```

**Arquivos Criados/Modificados:**
- ✅ `frontend/app/components/design-system/AccessibleInput.tsx` (68 linhas)
- ✅ `frontend/app/theme/tokens.ts` (+8 linhas: focus token)
- ✅ `frontend/app/theme/index.ts` (+20 linhas: MuiButton/InputBase/Dialog focus)
- ✅ `frontend/app/styles/theme-variables.css` (+17 linhas: @media prefers-reduced-motion)
- ✅ `frontend/app/lib/hooks/usePrefersReducedMotion.ts` (60 linhas)
- ✅ `frontend/jest.setup.ts` (68 linhas)
- ✅ `frontend/jest.config.ts` (42 linhas)
- ✅ `frontend/__tests__/accessibility/Button.a11y.test.tsx` (68 linhas)
- ✅ `frontend/__tests__/accessibility/AccessibleInput.a11y.test.tsx` (129 linhas)
- ✅ `frontend/__tests__/accessibility/Modal.a11y.test.tsx` (117 linhas)
- ✅ `frontend/package.json` (+6 scripts de teste)
- ✅ `docs/Designer-System.md` (seção "Acessibilidade" expandida: 500+ linhas)

**Dependências Instaladas:**
- ts-node@10.9.2
- jest@30.2.0
- jest-axe@10.0.0
- @testing-library/react@16.3.0
- @testing-library/jest-dom@6.9.1
- @testing-library/user-event@14.6.1
- jest-environment-jsdom@30.2.0
- @types/jest@30.0.0
- @types/jest-axe@3.5.9

**Garantias WCAG 2.1 AA:**
| Critério | Status | Evidência |
|----------|--------|-----------|
| Contraste ≥4.5:1 | ✅ | Validado na tabela de cores (Designer-System.md) |
| Foco visível 2px | ✅ | tokens.focus aplicado em MUI theme |
| ARIA completo | ✅ | AccessibleInput com aria-invalid/describedby/live |
| Keyboard nav | ✅ | Tab order lógico, Enter/Space/Escape |
| Reduced motion | ✅ | CSS @media + usePrefersReducedMotion hook |
| Zero violations | ✅ | 25/25 testes jest-axe passando |

**Próximas Ações Recomendadas:**
1. ⏳ Validação manual com NVDA/VoiceOver (screen reader)
2. ⏳ Chrome DevTools Lighthouse audit (target: 100% accessibility score)
3. ⏳ Instalar axe DevTools extension — escanear páginas principais
4. ⏳ Testar keyboard navigation em fluxos principais (login, CRUD)

---

### 🟡 T-QA-003 — Frontend tests
- **Responsável:** QA
- **Prioridade:** 🔴 Alta (elevado de Média)
- **Estimativa:** 6h → 8h (ajustado)
- **Sprint:** Sprint 6
- **Status:** ⏳ Em Andamento (75% completo - 6/8h)
- **Deliverable:** Testes componentes (Jest/RTL) + E2E (Playwright)

#### Critérios de Aceitação
- [x] Testes componentes (Jest + React Testing Library)
  - [x] Button: estados default/loading/disabled + eventos onClick (30 testes)
  - [x] AccessibleInput: label linking, aria attributes, mensagens erro (18 testes)
  - [x] Modal/Dialog: open/close, focus trap, keyboard (Escape) (19 testes)
- [x] E2E login flow (Playwright) - **IMPLEMENTADO mas aguardando correção de build**
  - [x] Login → Dashboard → Logout (3 cenários)
  - [x] Validar redirecionamentos e presença de KPIs
- [x] E2E CRUD receita (Playwright) - **IMPLEMENTADO mas aguardando correção de build**
  - [x] Criar → Verificar tabela → Editar → Deletar (6 cenários)
  - [x] Usar selectors estáveis (data-testid)
- [x] Scripts configurados: `test:unit`, `test:e2e`, `test:coverage`

#### Resultados Alcançados

**Testes Unitários (100% completo):**
```bash
$ pnpm test:unit

Test Suites: 3 passed, 3 total
Tests:       67 passed, 67 total
Time:        ~5s

✅ Button.test.tsx: 30 testes (rendering, click, loading, disabled, variants, fullWidth, props, a11y)
✅ AccessibleInput.test.tsx: 18 testes (rendering, helper text, ARIA, user interaction, props)
✅ Modal.test.tsx: 19 testes (rendering, interaction, accessibility, content)
```

**Testes E2E (Implementados, aguardando correção de build):**
- ✅ `e2e/login.spec.ts` - 3 cenários (login success, invalid credentials, logout)
- ✅ `e2e/receitas.spec.ts` - 6 cenários (create, edit, delete, filter, paginate, search)
- ✅ `e2e/fixtures/auth.ts` - Setup de autenticação reutilizável

**Bloqueio Atual:**
- ❌ Next.js build error: `next/headers` sendo importado em componente client (`app/lib/auth/tokens.ts`)
- 🔧 **Ação necessária:** Separar lógica server/client em tokens.ts ou usar alternativa client-side
- ⏸️ Testes E2E não podem rodar até resolver erro de build

**Arquivos Criados:**
- `frontend/e2e/login.spec.ts` (66 linhas) - Testes login/logout
- `frontend/e2e/receitas.spec.ts` (151 linhas) - Testes CRUD completo
- `frontend/e2e/fixtures/auth.ts` (45 linhas) - Helper autenticação
- `frontend/__tests__/unit/Button.test.tsx` (197 linhas) - 30 testes unitários
- `frontend/__tests__/unit/AccessibleInput.test.tsx` (143 linhas) - 18 testes unitários
- `frontend/__tests__/unit/Modal.test.tsx` (161 linhas) - 19 testes unitários

#### Plano de Execução Detalhado

**Fase 1: Component Tests - Button (1.5h)**

1.1. Criar `__tests__/unit/Button.test.tsx`:
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '@/components/ui/Button';

describe('Button Component', () => {
  it('should render with children', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should call onClick when clicked', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click</Button>);
    fireEvent.click(screen.getByText('Click'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should not call onClick when disabled', () => {
    const handleClick = jest.fn();
    render(<Button disabled onClick={handleClick}>Disabled</Button>);
    fireEvent.click(screen.getByText('Disabled'));
    expect(handleClick).not.toHaveBeenCalled();
  });

  it('should show loading spinner when loading', () => {
    render(<Button loading>Loading</Button>);
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('should be disabled when loading', () => {
    render(<Button loading>Loading</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('should apply variant styles', () => {
    const { rerender } = render(<Button variant="primary">Primary</Button>);
    expect(screen.getByRole('button')).toHaveClass('MuiButton-contained');

    rerender(<Button variant="outline">Outline</Button>);
    expect(screen.getByRole('button')).toHaveClass('MuiButton-outlined');
  });

  it('should support fullWidth prop', () => {
    render(<Button fullWidth>Full Width</Button>);
    expect(screen.getByRole('button')).toHaveStyle({ width: '100%' });
  });
});
```

**Fase 2: Component Tests - AccessibleInput (1.5h)**

2.1. Criar `__tests__/unit/AccessibleInput.test.tsx`:
```typescript
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AccessibleInput } from '@/components/design-system/AccessibleInput';

describe('AccessibleInput Component', () => {
  it('should render with label', () => {
    render(<AccessibleInput id="test" label="Email" />);
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
  });

  it('should link label to input via htmlFor', () => {
    render(<AccessibleInput id="email-input" label="Email" />);
    const label = screen.getByText('Email');
    const input = screen.getByLabelText('Email');
    expect(label).toHaveAttribute('for', 'email-input');
    expect(input).toHaveAttribute('id', 'email-input');
  });

  it('should display helper text', () => {
    render(
      <AccessibleInput
        id="test"
        label="Email"
        helperText="Enter your email"
      />
    );
    expect(screen.getByText('Enter your email')).toBeInTheDocument();
  });

  it('should have aria-invalid when error', () => {
    render(
      <AccessibleInput id="test" label="Email" error />
    );
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true');
  });

  it('should link helper text with aria-describedby', () => {
    render(
      <AccessibleInput
        id="email"
        label="Email"
        error
        helperText="Invalid email format"
      />
    );
    const input = screen.getByLabelText('Email');
    const helperText = screen.getByText('Invalid email format');
    expect(input).toHaveAttribute('aria-describedby', helperText.id);
  });

  it('should have role="alert" on error helper text', () => {
    render(
      <AccessibleInput
        id="test"
        label="Email"
        error
        helperText="Error message"
      />
    );
    const helperText = screen.getByText('Error message');
    expect(helperText).toHaveAttribute('role', 'alert');
  });

  it('should have aria-live="polite" on error helper text', () => {
    render(
      <AccessibleInput
        id="test"
        label="Email"
        error
        helperText="Error message"
      />
    );
    const helperText = screen.getByText('Error message');
    expect(helperText).toHaveAttribute('aria-live', 'polite');
  });

  it('should accept user input', async () => {
    const user = userEvent.setup();
    render(<AccessibleInput id="test" label="Name" />);
    const input = screen.getByLabelText('Name') as HTMLInputElement;
    await user.type(input, 'John Doe');
    expect(input.value).toBe('John Doe');
  });

  it('should show focus outline on focus', async () => {
    const user = userEvent.setup();
    render(<AccessibleInput id="test" label="Email" />);
    const input = screen.getByLabelText('Email');
    await user.click(input);
    expect(input).toHaveFocus();
  });
});
```

**Fase 3: Component Tests - Modal (1h)**

3.1. Criar `__tests__/unit/Modal.test.tsx`:
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@/components/ui/Dialog';
import { Button } from '@/components/ui/Button';

describe('Modal Component', () => {
  it('should not render when closed', () => {
    render(
      <Dialog open={false} onClose={() => {}}>
        <DialogTitle>Test Modal</DialogTitle>
      </Dialog>
    );
    expect(screen.queryByText('Test Modal')).not.toBeInTheDocument();
  });

  it('should render when open', () => {
    render(
      <Dialog open onClose={() => {}}>
        <DialogTitle>Test Modal</DialogTitle>
        <DialogContent>Modal content</DialogContent>
      </Dialog>
    );
    expect(screen.getByText('Test Modal')).toBeInTheDocument();
    expect(screen.getByText('Modal content')).toBeInTheDocument();
  });

  it('should call onClose when clicking backdrop', () => {
    const handleClose = jest.fn();
    render(
      <Dialog open onClose={handleClose}>
        <DialogTitle>Test</DialogTitle>
      </Dialog>
    );
    const backdrop = screen.getByRole('presentation').firstChild;
    if (backdrop) fireEvent.click(backdrop as Element);
    expect(handleClose).toHaveBeenCalled();
  });

  it('should call onClose when pressing Escape', () => {
    const handleClose = jest.fn();
    render(
      <Dialog open onClose={handleClose}>
        <DialogTitle>Test</DialogTitle>
      </Dialog>
    );
    fireEvent.keyDown(screen.getByRole('dialog'), { key: 'Escape' });
    expect(handleClose).toHaveBeenCalled();
  });

  it('should focus first focusable element', () => {
    render(
      <Dialog open onClose={() => {}}>
        <DialogTitle>Test</DialogTitle>
        <DialogContent>
          <Button>First Button</Button>
          <Button>Second Button</Button>
        </DialogContent>
      </Dialog>
    );
    expect(screen.getByText('First Button')).toHaveFocus();
  });

  it('should have aria-labelledby', () => {
    render(
      <Dialog open onClose={() => {}}>
        <DialogTitle id="dialog-title">Test Modal</DialogTitle>
        <DialogContent>Content</DialogContent>
      </Dialog>
    );
    const dialog = screen.getByRole('dialog');
    expect(dialog).toHaveAttribute('aria-labelledby', 'dialog-title');
  });

  it('should render actions', () => {
    render(
      <Dialog open onClose={() => {}}>
        <DialogTitle>Test</DialogTitle>
        <DialogContent>Content</DialogContent>
        <DialogActions>
          <Button>Cancel</Button>
          <Button>Confirm</Button>
        </DialogActions>
      </Dialog>
    );
    expect(screen.getByText('Cancel')).toBeInTheDocument();
    expect(screen.getByText('Confirm')).toBeInTheDocument();
  });
});
```

**Fase 4: Playwright Setup + E2E Login (2h)**

4.1. Instalar Playwright:
```bash
cd frontend
pnpm add -D @playwright/test
npx playwright install
```

4.2. Criar `playwright.config.ts`:
```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'pnpm dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

4.3. Criar `.env.test` (credenciais de teste):
```env
TEST_USER_EMAIL=teste@barberanalytics.com
TEST_USER_PASSWORD=senha123
NEXT_PUBLIC_API_URL=http://localhost:8080
```

4.4. Criar `e2e/auth.setup.ts` (fixture de autenticação):
```typescript
import { test as setup, expect } from '@playwright/test';

const authFile = 'playwright/.auth/user.json';

setup('authenticate', async ({ page }) => {
  await page.goto('/login');
  await page.getByLabel('Email').fill(process.env.TEST_USER_EMAIL!);
  await page.getByLabel('Senha').fill(process.env.TEST_USER_PASSWORD!);
  await page.getByRole('button', { name: /entrar/i }).click();

  await page.waitForURL('/dashboard');
  expect(page.url()).toContain('/dashboard');

  await page.context().storageState({ path: authFile });
});
```

4.5. Criar `e2e/login.spec.ts`:
```typescript
import { test, expect } from '@playwright/test';

test.describe('Login Flow', () => {
  test('should login successfully and redirect to dashboard', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel('Email').fill(process.env.TEST_USER_EMAIL!);
    await page.getByLabel('Senha').fill(process.env.TEST_USER_PASSWORD!);
    await page.getByRole('button', { name: /entrar/i }).click();

    await page.waitForURL('/dashboard');
    expect(page.url()).toContain('/dashboard');

    // Verificar KPIs
    await expect(page.getByTestId('kpi-receita')).toBeVisible();
    await expect(page.getByTestId('kpi-despesa')).toBeVisible();
    await expect(page.getByTestId('kpi-saldo')).toBeVisible();
  });

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel('Email').fill('invalid@test.com');
    await page.getByLabel('Senha').fill('wrongpassword');
    await page.getByRole('button', { name: /entrar/i }).click();

    await expect(page.getByText(/credenciais inválidas/i)).toBeVisible();
  });

  test('should logout successfully', async ({ page }) => {
    await page.goto('/dashboard');

    await page.getByRole('button', { name: /menu/i }).click();
    await page.getByRole('menuitem', { name: /sair/i }).click();

    await page.waitForURL('/login');
    expect(page.url()).toContain('/login');
  });
});
```

**Fase 5: E2E CRUD Receita (2h)**

5.1. Criar `e2e/receitas.spec.ts`:
```typescript
import { test, expect } from '@playwright/test';

test.describe('Receitas CRUD', () => {
  test.use({ storageState: 'playwright/.auth/user.json' });

  test('should create, edit, and delete receita', async ({ page }) => {
    await page.goto('/receitas');

    // Criar
    await page.getByRole('button', { name: /nova receita/i }).click();
    await page.getByLabel('Descrição').fill('Corte de cabelo teste E2E');
    await page.getByLabel('Valor').fill('50.00');
    await page.getByLabel('Data').fill('2025-11-15');
    await page.getByLabel('Categoria').selectOption({ label: 'Serviços' });
    await page.getByRole('button', { name: /salvar/i }).click();

    // Verificar toast
    await expect(page.getByText(/receita criada com sucesso/i)).toBeVisible();

    // Verificar tabela
    await expect(page.getByText('Corte de cabelo teste E2E')).toBeVisible();

    // Editar (opcional - se implementado)
    const row = page.getByText('Corte de cabelo teste E2E').locator('..');
    await row.getByRole('button', { name: /editar/i }).click();
    await page.getByLabel('Valor').fill('55.00');
    await page.getByRole('button', { name: /salvar/i }).click();

    await expect(page.getByText(/receita atualizada/i)).toBeVisible();

    // Deletar
    await row.getByRole('button', { name: /excluir/i }).click();
    await page.getByRole('button', { name: /confirmar/i }).click();

    await expect(page.getByText(/receita excluída/i)).toBeVisible();
    await expect(page.getByText('Corte de cabelo teste E2E')).not.toBeVisible();
  });

  test('should filter receitas by period', async ({ page }) => {
    await page.goto('/receitas');

    await page.getByLabel('Data inicial').fill('2025-11-01');
    await page.getByLabel('Data final').fill('2025-11-30');
    await page.getByRole('button', { name: /filtrar/i }).click();

    await page.waitForResponse(response =>
      response.url().includes('/receitas') && response.status() === 200
    );

    // Verificar se tabela foi atualizada
    await expect(page.getByRole('table')).toBeVisible();
  });

  test('should paginate receitas', async ({ page }) => {
    await page.goto('/receitas');

    const firstRowText = await page.locator('tbody tr:first-child td:first-child').textContent();

    await page.getByRole('button', { name: /próxima página/i }).click();

    const newFirstRowText = await page.locator('tbody tr:first-child td:first-child').textContent();
    expect(newFirstRowText).not.toBe(firstRowText);
  });
});
```

**Fase 6: Scripts & CI (0.5h)**

6.1. Atualizar `package.json`:
```json
{
  "scripts": {
    "test:unit": "jest --testPathPattern='__tests__/unit'",
    "test:a11y": "jest --testPathPattern='__tests__/accessibility'",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:debug": "playwright test --debug",
    "test:coverage": "jest --coverage",
    "test": "pnpm test:unit && pnpm test:a11y"
  }
}
```

6.2. Criar `.github/workflows/frontend-tests.yml` (CI):
```yaml
name: Frontend Tests

on:
  push:
    branches: [develop, main]
    paths:
      - 'frontend/**'
  pull_request:
    branches: [develop, main]
    paths:
      - 'frontend/**'

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: pnpm/action-setup@v2
        with:
          version: 8
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
          cache: 'pnpm'
          cache-dependency-path: frontend/pnpm-lock.yaml

      - name: Install dependencies
        working-directory: frontend
        run: pnpm install

      - name: Run unit tests
        working-directory: frontend
        run: pnpm test:unit

      - name: Run a11y tests
        working-directory: frontend
        run: pnpm test:a11y

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: pnpm/action-setup@v2
        with:
          version: 8
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
          cache: 'pnpm'

      - name: Install dependencies
        working-directory: frontend
        run: pnpm install

      - name: Install Playwright browsers
        working-directory: frontend
        run: npx playwright install --with-deps

      - name: Run E2E tests
        working-directory: frontend
        run: pnpm test:e2e
        env:
          TEST_USER_EMAIL: ${{ secrets.TEST_USER_EMAIL }}
          TEST_USER_PASSWORD: ${{ secrets.TEST_USER_PASSWORD }}

      - uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: frontend/playwright-report/
          retention-days: 30
```

**Validação Final:**
- [x] Jest: `pnpm test:unit` - Button (30), AccessibleInput (18), Modal (19) = **67 testes passando**
- [x] Axe: `pnpm test:a11y` - Zero violações (25 testes de acessibilidade)
- [ ] E2E: `pnpm test:e2e` - Login (3) + CRUD receita (6) = **9 testes implementados**
  - ⚠️ **Bloqueado:** Erro de build Next.js (`next/headers` em client component)
  - 🔧 **Ação necessária:** Refatorar `app/lib/auth/tokens.ts` para separar lógica server/client
- [x] Coverage: `pnpm test:coverage` - Configurado
- [x] CI: GitHub Actions workflow criado (`.github/workflows/frontend-tests.yml`)
- [ ] Documentação atualizada no README

**Próximos Passos:**

1. **Resolver erro de build (CRÍTICO):**
   ```typescript
   // Problema: app/lib/auth/tokens.ts importa next/headers em contexto client
   // Solução: Criar tokens-server.ts e tokens-client.ts separados
   // - tokens-server.ts: usa next/headers (cookies)
   // - tokens-client.ts: usa localStorage
   // - AuthContext deve usar tokens-client.ts
   ```

2. **Executar testes E2E:**
   ```bash
   cd frontend
   pnpm dev  # Iniciar servidor Next.js
   pnpm test:e2e  # Executar testes Playwright
   ```

3. **Validar cobertura de testes:**
   ```bash
   pnpm test:coverage
   # Target: >80% lines, >70% branches
   ```

4. **Adicionar data-testid em componentes:**
   - DashboardKPI cards: `data-testid="kpi-receita"`, `data-testid="kpi-despesa"`, etc.
   - FinancialTable rows: `data-testid="receita-row-{id}"`
   - Buttons de ação: `data-testid="btn-edit"`, `data-testid="btn-delete"`

5. **Documentar no README.md:**
   - Seção "Running Tests" com comandos disponíveis
   - Troubleshooting para erros comuns
   - Guia para adicionar novos testes

---

## 📈 Métricas de Sucesso

### Fase 4 completa quando:
- [ ] ✅ Todas as 16 tasks (frontend + QA) concluídas (100%)
- [ ] ✅ Frontend Next.js + MUI 5 estruturado com tokens do Designer System
- [ ] ✅ DayPilot Scheduler responsivo e tematizado funcionando
- [ ] ✅ Tema dark/light com ThemeStore, ThemeToggle, CSS vars e preferência do SO
- [ ] ✅ Acessibilidade WCAG 2.1 AA aplicada (contraste, focus, aria, reduz motion)
- [ ] ✅ Integração com backend Go funcionando
- [ ] ✅ Responsividade testada (mobile + desktop)
- [ ] ✅ Deploy em staging realizado
- [ ] ✅ Testes E2E passando

---

## 🎯 Deliverables da Fase 4

| # | Deliverable | Status |
|---|-------------|--------|
| 1 | Next.js 15 configurado com MUI + tokens | ⏳ Pendente |
| 2 | API Client com interceptors | ⏳ Pendente |
| 3 | Auth + Protected routes | ⏳ Pendente |
| 4 | Dashboard com KPIs + charts | ⏳ Pendente |
| 5 | Páginas CRUD (receitas, despesas) | ⏳ Pendente |
| 6 | Página Assinaturas | ⏳ Pendente |
| 7 | Hooks customizados | ⏳ Pendente |
| 8 | Forms validados (Zod) | ⏳ Pendente |
| 9 | Componentes do Design System (MUI + acessibilidade) | ⏳ Pendente |
| 10 | Theme + Dark Mode (ThemeStore + CSS vars) | ✅ Concluído |
| 11 | DayPilot Scheduler responsivo e tematizado | ✅ Concluído |
| 12 | Acessibilidade WCAG 2.1 AA aplicada | ⏳ Pendente |
| 13 | Testes E2E funcionando | ⏳ Pendente |

---

## 🚀 Próximos Passos

Após completar **100%** da Fase 4:

👉 **Iniciar FASE 5 — Migração** (`Tarefas/FASE_5_MIGRACAO.md`)

**Resumo Fase 5:**
- Feature flags (beta mode)
- Migração progressiva de dados MVP → v2
- Testes de regressão
- Desativação gradual do MVP 1.0

---

## 📝 Detalhamento Técnico Selecionado

### T-FE-002 — API Client (Exemplo)

```typescript
// lib/api/client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
});

// Request interceptor: add JWT
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: handle 401
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Tentar refresh token
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const { data } = await axios.post('/auth/refresh', { refreshToken });
          localStorage.setItem('access_token', data.access_token);
          // Retry original request
          return apiClient(error.config);
        } catch {
          // Logout
          localStorage.clear();
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

### T-FE-013 — Theme Store & Toggle (Exemplo)

```tsx
// lib/store/themeStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

type ThemeMode = 'light' | 'dark';

interface ThemeStore {
  theme: ThemeMode;
  toggleTheme: () => void;
  setTheme: (mode: ThemeMode) => void;
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set) => ({
      theme: 'light',
      toggleTheme: () => set((state) => ({
        theme: state.theme === 'light' ? 'dark' : 'light',
      })),
      setTheme: (mode: ThemeMode) => set({ theme: mode }),
    }),
    {
      name: 'bap-theme-storage',
      storage: typeof window !== 'undefined' ? localStorage : undefined,
    }
  )
);
```

```tsx
// components/ThemeToggle.tsx
'use client';
import { IconButton } from '@mui/material';
import { Brightness4, Brightness7 } from '@mui/icons-material';
import { useThemeStore } from '@/lib/store/themeStore';

export function ThemeToggle() {
  const { theme, toggleTheme } = useThemeStore();
  const isDark = theme === 'dark';

  return (
    <IconButton
      onClick={toggleTheme}
      title={isDark ? 'Passar para light mode' : 'Passar para dark mode'}
      aria-label="theme-toggle"
    >
      {isDark ? <Brightness7 /> : <Brightness4 />}
    </IconButton>
  );
}
```

---

**Última Atualização:** 14/11/2025
**Status:** 🔴 Não Iniciado (0%)
**Próxima Revisão:** Após completar 50% das tarefas
