> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🎨 Guia de Desenvolvimento - Frontend (Next.js 16.0.3)

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Guia Prático

---

## 📋 Índice

1. [Setup Local](#setup-local)
2. [Estrutura de Projeto](#estrutura-de-projeto)
3. [Convenções](#convenções)
4. [Desenvolvimento](#desenvolvimento)
5. [Styling](#styling)
6. [Estado & Dados](#estado--dados)
7. [Testing](#testing)
8. [Performance](#performance)

---

## 🚀 Setup Local

### Pré-requisitos

```bash
# Node.js
node --version  # Mínimo: 18.17

# Package Manager
npm --version ou yarn --version
```

### Setup

```bash
# 1. Clone/Create
git clone https://github.com/seu-usuario/barber-analytics-frontend.git
cd barber-analytics-frontend

# 2. Instalar dependências
npm install

# 3. Copiar .env
cp .env.example .env.local
# Editar: NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# 4. Rodar dev server
npm run dev

# 5. Abrir navegador
# http://localhost:3000
```

---

## 📁 Estrutura de Projeto

```
frontend/
│
├── app/
│   ├── (auth)/                    # Grupo de layout: Auth (pública)
│   │   ├── layout.tsx
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   │
│   ├── (dashboard)/               # Grupo de layout: Dashboard (privada)
│   │   ├── layout.tsx
│   │   ├── page.tsx              # Dashboard home
│   │   ├── financial/
│   │   │   ├── receitas/
│   │   │   │   ├── page.tsx
│   │   │   │   └── [id]/edit.tsx
│   │   │   ├── despesas/
│   │   │   └── cashflow/
│   │   ├── subscriptions/
│   │   │   ├── page.tsx
│   │   │   └── [id]/
│   │   └── settings/
│   │
│   ├── api/                      # Route handlers (API Backend)
│   │   └── auth/
│   │       └── logout/route.ts
│   │
│   ├── components/               # Componentes (agora dentro de app/)
│   │   ├── ui/                   # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── form.tsx
│   │   │   ├── table.tsx
│   │   │   └── ...
│   │   │
│   │   ├── layout/
│   │   │   ├── navbar.tsx
│   │   │   ├── sidebar.tsx
│   │   │   ├── footer.tsx
│   │   │   └── mobile-nav.tsx
│   │   │
│   │   ├── forms/
│   │   │   ├── receita-form.tsx
│   │   │   ├── despesa-form.tsx
│   │   │   └── login-form.tsx
│   │   │
│   │   ├── dashboard/
│   │   │   ├── kpi-cards.tsx
│   │   │   ├── chart-revenue.tsx
│   │   │   └── recent-activity.tsx
│   │   │
│   │   └── common/
│   │       ├── loading-skeleton.tsx
│   │       ├── empty-state.tsx
│   │       └── error-boundary.tsx
│   │
│   ├── lib/                      # Utils e Hooks (agora dentro de app/)
│   │   ├── api/
│   │   │   ├── client.ts         # Axios instance
│   │   │   └── endpoints.ts
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   ├── useTenant.ts
│   │   │   ├── useReceitas.ts
│   │   │   └── usePagination.ts
│   │   ├── utils/
│   │   │   ├── format.ts         # Formatação (moeda, data)
│   │   │   ├── validation.ts
│   │   │   └── storage.ts
│   │   └── types.ts
│   │
│   ├── globals.css
│   └── layout.tsx               # Root layout
│
├── types/
│   └── index.ts                 # TypeScript types globais
│
├── styles/
│   ├── globals.css
│   └── variables.css            # CSS custom properties
│
├── public/
│   ├── images/
│   ├── icons/
│   └── logo.svg
│
├── package.json
├── tsconfig.json
├── tailwind.config.ts
├── next.config.js
├── .env.example
└── README.md
```

---

## 🎯 Convenções

### Naming

```typescript
// Components: PascalCase
export function ReceitaForm() {}

// Hooks: camelCase, prefixo 'use'
export function useReceitas() {}

// Utilities: camelCase
export const formatCurrency = () => {}

// Types: PascalCase
type Receita = {
    id: string
    valor: number
}

// Constants: UPPER_SNAKE_CASE
const API_BASE_URL = \"http://localhost:8080\"
```

### Imports

```typescript
// ✅ Agrupar imports
import React from \"react\"
import { useState } from \"react\"

import { ReceitaForm } from \"@/components/forms/receita-form\"
import { useReceitas } from \"@/lib/hooks/useReceitas\"
import { formatCurrency } from \"@/lib/utils/format\"
import type { Receita } from \"@/types\"
```

---

## 💻 Desenvolvimento

### Criar Página com Formulário

1. **Criar o tipo**

```typescript
// lib/types.ts
export type Receita = {
    id: string
    descricao: string
    valor: number
    categoria: string
    data: Date
    status: \"CONFIRMADA\" | \"RECEBIDA\" | \"CANCELADA\"
}

export type CreateReceitaInput = Omit<Receita, \"id\" | \"status\">
```

2. **Criar hook de API**

```typescript
// lib/hooks/useReceitas.ts
import { useQuery, useMutation, useQueryClient } from \"@tanstack/react-query\"
import { receitaApi } from \"@/lib/api/endpoints\"

export function useReceitas(tenantId: string, params?: ListParams) {
    return useQuery({
        queryKey: [\"receitas\", tenantId, params],
        queryFn: () => receitaApi.list(tenantId, params),
        staleTime: 5 * 60 * 1000, // 5 minutos
    })
}

export function useCreateReceita() {
    const queryClient = useQueryClient()
    
    return useMutation({
        mutationFn: (data: CreateReceitaInput) => receitaApi.create(data),
        onSuccess: () => {
            // Invalidar cache
            queryClient.invalidateQueries({ queryKey: [\"receitas\"] })
        },
    })
}
```

3. **Criar formulário**

```typescript
// components/forms/receita-form.tsx
import { useForm } from \"react-hook-form\"
import { zodResolver } from \"@hookform/resolvers/zod\"
import { z } from \"zod\"

const schema = z.object({
    descricao: z.string().min(1, \"Obrigatório\").max(255),
    valor: z.string().regex(/^\\d+(\\.\\d{2})?$/, \"Valor inválido\"),
    categoria: z.string().min(1, \"Selecione uma categoria\"),
    data: z.date(),
})

export function ReceitaForm() {
    const { createReceita, isPending } = useCreateReceita()
    const { handleSubmit, control } = useForm({
        resolver: zodResolver(schema),
    })
    
    const onSubmit = async (data: z.infer<typeof schema>) => {
        await createReceita.mutateAsync(data)
    }
    
    return (
        <form onSubmit={handleSubmit(onSubmit)} className=\"space-y-4\">
            {/* Campos do formulário */}
            <button disabled={isPending}>
                {isPending ? \"Salvando...\" : \"Salvar\"}
            </button>
        </form>
    )
}
```

4. **Criar página**

```typescript
// app/(dashboard)/financial/receitas/page.tsx
\"use client\"

import { useState } from \"react\"
import { ReceitaForm } from \"@/components/forms/receita-form\"
import { ReceitaTable } from \"@/components/dashboard/receita-table\"
import { useReceitas } from \"@/lib/hooks/useReceitas\"
import { useTenant } from \"@/lib/hooks/useTenant\"

export default function ReceitasPage() {
    const { tenantId } = useTenant()
    const [params, setParams] = useState({ page: 1, pageSize: 50 })
    
    const { data, isLoading, error } = useReceitas(tenantId, params)
    
    if (isLoading) return <LoadingSkeleton />
    if (error) return <ErrorState error={error} />
    
    return (
        <div className=\"space-y-6\">
            <div className=\"flex justify-between items-center\">
                <h1 className=\"text-3xl font-bold\">Receitas</h1>
                <Dialog>
                    <DialogTrigger>Adicionar</DialogTrigger>
                    <DialogContent>
                        <ReceitaForm />
                    </DialogContent>
                </Dialog>
            </div>
            
            <ReceitaTable data={data?.data || []} />
            
            <Pagination
                current={params.page}
                total={data?.pagination.total}
                onChange={(page) => setParams({ ...params, page })}
            />
        </div>
    )
}
```

---

## 🎨 Styling

### Tailwind CSS

```typescript
// ✅ Usar classes Tailwind
<div className=\"flex items-center justify-between p-4 bg-white rounded-lg shadow\">

// ✅ Componentes reutilizáveis
<Button variant=\"primary\" size=\"lg\">Salvar</Button>

// ✅ Responsive
<div className=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4\">

// ❌ Evitar inline styles
<div style={{ padding: \"16px\" }}>  // ❌
<div className=\"p-4\">             // ✅
```

### Dark Mode (Futuro)

```typescript
// next.config.js
const config = {
    theme: {
        darkMode: \"class\",
        extend: {
            colors: {
                brand: \"#FF6B35\",
            },
        },
    },
}
```

---

## 🔄 Estado & Dados

### TanStack Query Setup

```typescript
// app/providers.tsx
\"use client\"

import { QueryClient, QueryClientProvider } from \"@tanstack/react-query\"
import { ReactNode } from \"react\"

const queryClient = new QueryClient()

export function Providers({ children }: { children: ReactNode }) {
    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    )
}
```

### API Client

```typescript
// lib/api/client.ts
import axios from \"axios\"
import { useAuth } from \"@/lib/hooks/useAuth\"

export const apiClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL,
})

// Interceptor: Adicionar JWT
apiClient.interceptors.request.use((config) => {
    const token = localStorage.getItem(\"access_token\")
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Interceptor: Refresh token em 401
apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 401) {
            // Refresh token logic
        }
        return Promise.reject(error)
    }
)
```

---

## 🧪 Testing

### Unit Tests (Jest)

```typescript
// __tests__/components/receita-form.test.tsx
import { render, screen, fireEvent } from \"@testing-library/react\"
import { ReceitaForm } from \"@/components/forms/receita-form\"

describe(\"ReceitaForm\", () => {
    it(\"renders form\", () => {
        render(<ReceitaForm />)
        expect(screen.getByText(\"Salvar\")).toBeInTheDocument()
    })
    
    it(\"submits form with valid data\", async () => {
        render(<ReceitaForm />)
        fireEvent.change(screen.getByLabelText(\"Descrição\"), {
            target: { value: \"Corte\" },
        })
        fireEvent.click(screen.getByText(\"Salvar\"))
        
        // Assertions
    })
})
```

### E2E Tests (Playwright)

```typescript
// e2e/receitas.spec.ts
import { test, expect } from \"@playwright/test\"

test(\"should create receita\", async ({ page }) => {
    await page.goto(\"/financial/receitas\")
    await page.click(\"button:has-text('Adicionar')\")
    await page.fill(\"input[name='descricao']\", \"Corte\")
    await page.fill(\"input[name='valor']\", \"50.00\")
    await page.click(\"button[type='submit']\")
    
    await expect(page).toContainText(\"Corte\")
})
```

---

## ⚡ Performance

### Code Splitting

```typescript
// ✅ Lazy load componentes pesados
import dynamic from \"next/dynamic\"

const ChartComponent = dynamic(
    () => import(\"@/components/dashboard/chart\"),
    { loading: () => <p>Carregando...</p> }
)
```

### Image Optimization

```typescript
// ✅ Usar next/image
import Image from \"next/image\"

<Image
    src=\"/logo.svg\"
    alt=\"Logo\"
    width={200}
    height={50}
    priority
/>
```

### Memoization

```typescript
// ✅ Usar React.memo para componentes puros
export const ReceitaTable = React.memo(function ReceitaTable({ data }) {
    return <table>{/* ... */}</table>
})
```

---

**Status:** ✅ Guia completo
