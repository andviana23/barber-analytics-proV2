# VALTARIS v 1.0 — Guia de Desenvolvimento Frontend

Guia prático para trabalhar no frontend do VALTARIS usando o Design System “Cyber Luxury”. Light é padrão; Dark é opt-in via store.

## Stack Oficial (resumo)
- Next.js 16 (App Router), React 19, TypeScript.
- MUI 5 + Emotion (ThemeRegistry + `create-emotion-cache`).
- CSS Variables `--valtaris-*` para light/dark.
- Zustand (tema/prefs/layout), React Hook Form + Zod, TanStack Query.
- Testes: Jest + Testing Library; Playwright para E2E.

## Setup Local
```bash
node --version   # >= 18.17
npm install      # dentro de frontend/
cp .env.example .env.local
# configure NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
npm run dev
```
Abrir `http://localhost:3000`.

## Estrutura de Projeto (essencial)
```
frontend/
  app/
    layout.tsx              # Root; envolve ThemeRegistry
    theme-registry.tsx      # Emotion + MUI + CSS vars
    (dashboard)/...         # Telas privadas
    (auth)/...              # Telas públicas
    components/
      ui/                   # Wrappers DS: VButton, VCard, VModal, etc.
      forms/                # RHF + Zod wrappers
      layout/               # Sidebar, Topbar, DashboardLayout
    lib/                    # api client, hooks, utils, theme
    store/                  # Zustand stores (tema, layout)
  docs/03-frontend          # Documentação do DS
```

## Como criar uma nova tela alinhada ao DS
1) Garantir que o arquivo tem `"use client"` se usar hooks.  
2) A página já está sob `ThemeRegistry` em `app/layout.tsx`.  
3) Usar componentes do DS (`VButton`, `VCard`, `VModal`, wrappers de formulário).  
4) Estilizar via `sx`/styled e `var(--valtaris-*)`, nunca hex direto.  
5) Formulários: RHF + Zod + wrappers; mensagens de erro em `var(--valtaris-danger)`.  
6) Dados: TanStack Query com chaves que incluem `tenantId`; loading com skeleton e não com `setTimeout`.  
7) Agendadores/tabelas externas: classe `daypilot-valtaris` + CSS vars, sem CSS inline stringificado.  
8) Acessibilidade: focus ring visível (`accent.aqua`), hit area 40px+.

## Convenções rápidas
- Components: PascalCase (`ReceitaForm`). Hooks: `useX`. Types: PascalCase. Constantes: UPPER_SNAKE_CASE.
- Imports agrupados; paths absolutos com `@/`.
- Evite `any`; tipar mutations/queries com generics do TanStack Query.

## Styling e Tema
- MUI é o layer principal; Tailwind é opcional como utilitário, se já existir.
- Use `sx` com tokens: `sx={{ color: "var(--valtaris-text)" }}`.
- Modais e overlays devem usar vidro (`backdrop-filter: blur(...)`), borda metálica e sombras definidas em `01-FOUNDATIONS.md`.

## Estado & Dados
- TanStack Query: instanciar `QueryClient` no provider do App Router. StaleTime padrão 5min para listas.
- API client: axios com interceptors de auth; base URL via env.
- Zustand: persistir tema (`valtaris-theme`) ao trocar light/dark.

## Testing
- Unit: Jest + Testing Library para componentes e hooks.
- E2E: Playwright para fluxos críticos (login, receba/pagar, scheduler).
- Snapshot só para componentes puros; evite para páginas dinâmicas.

## Performance
- Code splitting com `next/dynamic` para gráficos pesados.
- `next/image` para assets; respeitar `priority` apenas em hero.
- Memorize tabelas grandes (`React.memo`) e derive memos com `useMemo`.

## Pontes com a Documentação do DS
- Tokens e materiais: `01-FOUNDATIONS.md`.
- ThemeRegistry/Theme Store: `02-ARCHITECTURE.md`.
- Componentes: `03-COMPONENTS.md`.
- Padrões de formulários e acessibilidade: `04-PATTERNS.md`.
- Componentes críticos e layouts: `COMPONENTES_CRITICOS.md`.
