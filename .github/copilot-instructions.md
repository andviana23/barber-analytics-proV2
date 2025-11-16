# 🧠 GitHub Copilot — Barber Analytics Pro

**Instruções Oficiais**

| Informação | Valor |
|-----------|-------|
| **Versão** | 4.0 |
| **Autor** | Andrey Viana |
| **Objetivo** | Transformar o Copilot em um desenvolvedor sênior do Barber Analytics Pro, seguindo fielmente a arquitetura, guias e design system oficiais |

---

## 🌐 0. Como usar estas instruções

Sempre que o Copilot gerar, alterar ou sugerir código, ele **deve seguir**:

### Documentação Principal

- `ARQUITETURA.md`
- `GUIA_DEV_BACKEND.md`
- `GUIA_DEV_FRONTEND.md`
- `GUIA_DEVOPS.md`
- `Designer-System.md`
- `BANCO_DE_DADOS.md`
- `FINANCEIRO.md`
- `ESTOQUE.md`
- `ASSINATURAS.md`
- `FLUXO_CRONS.md`
- `API_REFERENCE.md`

### Regra de Ouro

> 🎯 Se houver conflito entre qualquer coisa e estes documentos, **os documentos do projeto são a VERDADE absoluta**.
>
> Quando em dúvida: pergunte, descreva as opções e NUNCA invente padrão novo.

---

## 🎯 1. Regras Fundamentais (Críticas e Obrigatórias)

### 🔥 1.1 NUNCA escrever SQL fora do @pgsql

**Regra ABSOLUTA:**

Qualquer alteração de schema, tabela, índice, view, função, RLS, trigger ou migração:

- ✅ **DEVE** ser feita usando macros como `@pgsql_query`, `@pgsql_modify`, `@pgsql_connect`, `@pgsql_db_context` (padrão do projeto)
- ❌ **NÃO** pode ser gerada como `.sql` solto, nem comando direto em terminal, nem código SQL embedado em Go

**Quando precisar mexer em banco:**

- Gerar **apenas** o bloco de instruções @pgsql (sem executar nada por conta própria)
- Referir-se a `BANCO_DE_DADOS.md` para nomes de tabelas, colunas e convenções

---

### 🧾 1.2 NÃO criar arquivos `.md` sem o Andrey pedir

- ❌ Nada de criar documentação nova automaticamente
- ❌ Nada de gerar guias, readmes, checklists, changelogs sem pedido explícito
- ✅ Se Andrey pedir documentação, seguir tom e estrutura dos arquivos atuais (`GUIA_DEV_*`, `Designer-System.md`, etc.)

---

### 🧱 1.3 Seguir a Arquitetura Oficial SEM alterar

Backend e frontend já possuem arquitetura definida. O Copilot:

- ❌ **NÃO** pode criar nova camada, novo padrão ou "atalhos"
- ❌ **NÃO** pode jogar regra de negócio em lugar inadequado (handler, componente React, etc.)
- ✅ **DEVE** respeitar a divisão entre:
  - **Domain** — Lógica de negócio pura
  - **Application** — Casos de uso e orquestração
  - **Infrastructure** — Implementação técnica
  - **HTTP / UI** — Camada de apresentação
  - **Jobs / Crons** — Tarefas agendadas

Sempre que sugerir algo, deve estar claramente encaixado em uma dessas camadas, conforme `ARQUITETURA.md` e `GUIA_DEV_BACKEND.md`.

---

### 🎨 1.4 Seguir o Design System oficial (Designer-System.md)

Para frontend, é **LEI**:

#### Stack Base

- **MUI 5** + Design Tokens gerenciados no theme
- Tokens expostos via `theme` + constantes reutilizáveis (DayPilot, testes, CSS utilities)
- Dark/Light mode controlados por tokens
- Acessibilidade padrão WCAG 2.1 AA+

#### O Copilot **NÃO** pode:

- ❌ Usar Tailwind ou classes soltas tipo `bg-white`, `text-gray-700`, `#123456` direto
- ❌ Criar componente fora do sistema (ex: `<button>` cru com estilos manuais)
- ❌ Ignorar tokens de cor, spacing, tipografia e motion

#### O Copilot **DEVE**:

- ✅ Usar componentes e padrões documentados em `Designer-System.md`
- ✅ Consultar tokens de palette, spacing, radius, shadows, typography e motion
- ✅ Integrar DayPilot usando os tokens (cores, backgrounds, estados light/dark)

---

### 📦 1.5 Padrão de retorno `{ data, error }`

#### Backend (Go)

Usecases e handlers retornam estruturas no padrão:

```go
type Result[T any] struct {
    Data  *T
    Error error
}
```

#### Frontend (TypeScript)

Hooks e services retornam:

```typescript
type Result<T> = {
  data: T | null;
  error: Error | null;
};
```

**Regra:** Nunca retornar valores soltos em fluxos principais. Sempre encapsular em `{ data, error }` (ou equivalente detalhado do projeto).

---

### 🧱 1.6 Erros sempre com contexto

#### Backend (Go)

Usar logger padrão (ex.: zerolog) com mensagem de contexto:

```go
log.Error().Err(err).Msg("falha ao criar receita")
```

#### Frontend (TypeScript)

Mensagens claras para o usuário, em **PT-BR**:

```typescript
toast.error("Não foi possível salvar a receita. Tente novamente.");
```

---

### 🧬 1.7 Multi-tenant é obrigatório

#### Regras para qualquer acesso a dados:

- ✅ Sempre considerar `tenant_id` (no domínio, nas queries e no filtro de RLS)
- ✅ Nunca retornar dados de um tenant para outro
- ✅ Handlers devem obter `tenant_id` via middleware/ctx e repassar ao usecase
- ✅ Tabelas devem conter colunas obrigatórias (ver seção Banco)

---

## 🏗️ 2. Arquitetura Oficial (Visão Geral)

### 🧬 2.1 Backend (Go)

#### Estrutura Base

```
cmd/api/main.go                           → Entrypoint
internal/config                           → Configuração
internal/domain                           → Regras de negócio puras
internal/application                      → DTOs, usecases, mappers
internal/infrastructure
  ├── http                               → Handlers, middlewares, rotas
  ├── repository                         → Acesso a dados
  ├── external                           → Integrações (Asaas, etc.)
  ├── scheduler                          → Crons/jobs
  └── database                           → Conexão e abstrações
migrations                                → Histórico (ligado ao @pgsql)
tests                                     → Testes unitários/integrados
```

#### Princípios

- ✅ Dependência sempre aponta para o domínio
- ✅ Domain não conhece infra nem HTTP

---

### 🎨 2.2 Frontend (Next.js 15 + React 19 + MUI 5)

#### Estrutura Base

```
frontend/app
  ├── (auth)                            → Login/logout, reset de senha
  ├── (private)                         → Dashboards, financeiro, estoque, assinaturas
  └── ...

frontend/components
  └── design-system                     → Componentes visuais reutilizáveis

frontend/lib
  ├── api                               → Client HTTP, interceptors, base URLs
  ├── hooks                             → Hooks de domínio (useReceitas, useAssinaturas, etc.)
  └── store                             → Zustand/TanStack Query para estado global

frontend/theme                           → Tokens.ts, providers.tsx, theme-variables.css
```

#### Princípios

- ✅ App Router é a fonte de verdade das rotas
- ✅ Sem fetch direto em componentes de página (exceto loaders específicos)

---

## 🔧 3. Backend (Go) — Regras Detalhadas

### 3.1 Domain Layer

**Responsabilidades:**
- Somente lógica de negócio pura
- Entities representam o modelo de domínio (Receita, Despesa, Assinatura, Produto, etc.)
- Value Objects para conceitos imutáveis (Money, StatusAssinatura, Categoria, etc.)

**Restrições:**
- ❌ Não conhecer banco, HTTP, JSON, headers, contexto web

**Quando criar ou alterar regras:**
- Ver se isso é realmente regra de negócio → colocar em `internal/domain`
- Expor comportamento em métodos de entidade/serviço
- Não misturar validação de transporte (HTTP) com regra de domínio

---

### 3.2 Application Layer (Usecases + DTOs + Mappers)

**Responsabilidades:**
- Usecases recebem DTOs de entrada (já validados pelo handler) + context
- Usam repositórios via interfaces definidas no domínio (`internal/domain/repository`)
- Retornam DTOs de saída ou o tipo genérico `{ Data, Error }`
- Orquestram fluxo, não implementam infra

**Sempre que criar feature:**
1. Criar DTO de entrada/saída
2. Criar UseCase em `internal/application/usecase/<bounded-context>`
3. Registrar o UseCase na injeção de dependências
4. Só depois criar Handler HTTP chamando esse UseCase

---

### 3.3 Infrastructure Layer

#### Repositórios

- ✅ Vivem em `internal/infrastructure/repository`
- ✅ Implementam interfaces definidas em `internal/domain/repository`
- ✅ Usam abstrações de DB do projeto (ex.: `database.Connection`)
- ✅ Nunca contornam RLS ou multi-tenant
- ✅ Sempre incluem `tenant_id` nos filtros

#### HTTP (Handlers + Middlewares)

**Handlers devem fazer:**
- Bind/validação do request
- Extração de `tenant_id` / `user_id`
- Chamada a UseCases
- Montagem de resposta HTTP

**Middlewares devem fazer:**
- Autenticação (JWT)
- Resolução de tenant
- Logging de request/response

#### Scheduler / Crons

- ✅ Jobs em `internal/infrastructure/scheduler`
- ✅ Cada job chama usecases específicos (nunca lógica solta)
- ✅ Sempre com:
  - Idempotência (registrar execuções)
  - Circuit breaker
  - Logs estruturados

---

## 🎨 4. Frontend — Regras Detalhadas

### 4.1 Stack e Responsabilidades

- **Next.js 15** App Router
- **React 19**
- **MUI 5** com theming profundo (tokens)
- **TanStack Query** para dados remotos
- **Zustand** (ou store custom) para estado global crítico (auth, tenant, tema)
- **React Hook Form + Zod** para formulários

---

### 4.2 Design System (MUI + Tokens)

**Princípios:**
- ✅ Tokens expostos via theme (palette, spacing, typography, radius, shadow, motion)
- ✅ Componentes DS devem importar tokens em vez de hardcode (theme.palette.primary.main, etc.)

**Todos os componentes novos devem:**
- Usar DS existente como base (Button, Input, etc.)
- Ser consistentes em padding, radius e tipografia
- Estar em `components/design-system` (quando reutilizáveis)

**Exemplo — Wrapper de botão:**

```typescript
import { Button as MuiButton, ButtonProps } from '@mui/material';

export function Button(props: ButtonProps) {
  return <MuiButton variant="contained" color="primary" {...props} />;
}
```

**Exemplo — Uso de tokens via sx:**

```typescript
<Box
  sx={(theme) => ({
    padding: theme.spacing(2),
    borderRadius: theme.shape.borderRadius,
    backgroundColor: theme.palette.background.paper,
  })}
/>
```

---

### 4.3 Dados e Hooks

**Padrão de hooks:**

```typescript
function useReceitas(filters) {
  const { data, error, isLoading } = useQuery({
    queryKey: ['receitas', filters],
    queryFn: () => api.receitas.list(filters),
  });

  return { data, error, isLoading };
}
```

**Regra:** Criar hooks em `frontend/hooks` (ex.: `useReceitas`, `useDespesas`, `useAssinaturas`)

---

### 4.4 Forms

- ✅ React Hook Form + Zod para validações
- ✅ Mensagens em PT-BR
- ✅ Inputs DS integrados com register e error

---

### 4.5 DayPilot Scheduler Integration

- ✅ Configuração de cores e temas vem dos tokens
- ✅ Nada de inline color hardcoded
- ✅ Respeitar dark/light mode conforme `Designer-System.md`

---

## 🗄️ 5. Banco de Dados (Neon/Postgres)

### 5.1 Acesso

- ✅ **Sempre via macros/infra do projeto** (`@pgsql_*`)
- ❌ Nunca construir conexão manual via `sql.Open` sem seguir padrão

---

### 5.2 Colunas Obrigatórias

Cada tabela relevante deve ter, no mínimo:

```sql
tenant_id          UUID         NOT NULL
criado_em          TIMESTAMPTZ  DEFAULT NOW()
atualizado_em      TIMESTAMPTZ  DEFAULT NOW()
ativo              BOOLEAN      DEFAULT true      (quando aplicável)
```

---

### 5.3 RLS

- ✅ Todas as tabelas com dados sensíveis devem ter RLS habilitado
- ✅ Policies sempre baseadas em `tenant_id` ligado ao JWT

---

### 5.4 IDs

- ❌ Nunca usar `SERIAL`
- ✅ IDs preferencialmente `uuid DEFAULT gen_random_uuid()`

---

## 💰 6. Módulo Financeiro

(Ver detalhes em `FINANCEIRO.md`)

### Escopo

- ✅ Bounded context próprio (Financeiro)
- ✅ Lida com: Receitas, Despesas, Categorias, Centros de custo, DRE, Fluxo de caixa

### Regras para o Copilot

- ✅ Campos financeiros usam Value Object Money
- ✅ Não misturar lógica de assinatura com financeiro (cada módulo no seu contexto)
- ✅ DRE e fluxo de caixa devem respeitar diferença entre data de competência e data de pagamento
- ✅ Hooks e endpoints devem refletir isso claramente

---

## 📦 7. Módulo de Estoque

(Ver `ESTOQUE.md`)

### Entidades

- Produto
- Movimentação
- Compra
- Fornecedor

### Funcionalidades

- Controle de estoque (entradas/saídas)
- Integração com financeiro (compra gera despesa)
- Possibilidade de parcelas (forma de pagamento)

### Regras

- ✅ Movimentação de estoque sempre em função de um evento (compra, venda, ajuste)
- ✅ Integração com financeiro deve ser feita por usecases específicos (não dentro do repository)
- ✅ Formular lógica de estoque respeitando multi-tenant

---

## 🔁 8. Assinaturas (Clube do Trato)

(Ver `ASSINATURAS.md`)

### Fluxo

- Cliente paga assinatura
- Valor cai pré-datado na barbearia
- Repasse ao barbeiro é pós recebimento efetivo (ex.: atende em novembro, recebe comissão em janeiro)

### Regras

**Sempre diferenciar:**
- Status no Asaas (pago, pendente, cancelado)
- Status interno (confirmado para repasse, repasse feito, etc)

**Implementação:**
- ✅ Usecases de assinatura não podem usar Asaas diretamente
- ✅ Usar adapters em `infrastructure/external/asaas`
- ✅ Qualquer lógica de projeção de receita deve ser feita no módulo correto (nunca em handler ou componente)

---

## ⏱️ 9. Fluxo de Crons / Jobs

(Ver `FLUXO_CRONS.md`)

### O que é agendado

- Relatórios diários de financeiro
- Atualização de status de assinatura
- Consolidação de dashboards
- Health-check

### Regras

**Implementação:**
- ✅ Jobs nunca implementam regra de domínio diretamente
- ✅ Sempre chamam usecases

**Características:**
- ✅ Idempotentes (registrar execuções)
- ✅ Monitorados (logs + possíveis integrações com Prometheus/Grafana/Sentry)

**Copilot deve ajudar a:**
- Centralizar a orquestração
- Manter consistência de nomes
- Garantir que falhas sejam logadas com contexto

---

## ⚙️ 10. DevOps / CI/CD / Infra

(Ver `GUIA_DEVOPS.md`)

### Arquitetura

- Deploy automatizado via GitHub Actions
- Backend em VPS (PM2, NGINX reverse proxy)
- Frontend em Vercel (ou similar)
- Banco em Neon (Postgres fully-managed)

### Quando o Copilot mexer em `.github/workflows/*.yml`

- ✅ Não quebrar jobs existentes
- ✅ Seguir nomes e padrões do pipeline atual
- ✅ Respeitar variáveis de ambiente já definidas nos secrets

### Configurações de log/monitoramento

- ✅ Seguir instruções de Sentry/Prometheus/Grafana do guia
- ❌ Nunca expor secrets ou tokens em log

---

## 🔐 11. Segurança

- ✅ JWT RS256 para autenticação
- ✅ Refresh tokens com rotação
- ✅ Multi-tenant obrigatório em todas as queries
- ✅ RLS ativo nas tabelas sensíveis
- ❌ Nunca registrar senhas ou dados sensíveis em logs
- ✅ Inputs de usuário devem ser validados no backend (sem confiar só no frontend)

---

## 🧪 12. Testes

### Backend (Go)

- ✅ `go test ./...` deve passar

**Testes unitários para:**
- Entities
- Value Objects
- Usecases

**Testes de integração para:**
- Repositórios
- Endpoints HTTP críticos

### Frontend (TypeScript)

- ✅ React Testing Library + Vitest (ou Jest, conforme guia)
- ✅ Testes de acessibilidade com jest-axe/a11y

**Testes E2E para fluxos críticos:**
- Login
- Cadastro/edição de receita/despesa
- Criação/edição de assinatura
- Movimentação de estoque

---

## 🧩 13. Convenções de Nomenclatura

### Backend (Go)

| Tipo | Padrão | Exemplo |
|------|--------|---------|
| Usecases | PascalCase com "UseCase" | `CreateReceitaUseCase` |
| Repositories | PascalCase com "Repository" | `ReceitaRepository` |
| DTOs | PascalCase com "Input"/"Output" | `CreateReceitaInput` |
| Entidades | PascalCase | `Receita`, `Despesa` |

### Frontend (TypeScript)

| Tipo | Padrão | Exemplo |
|------|--------|---------|
| Hooks | camelCase com "use" | `useReceitas` |
| Services | camelCase com "Service" | `receitasService` |
| Páginas | PascalCase com "Page" | `ReceitasPage` |
| Componentes DS | PascalCase | `Button`, `Card`, `DataTable` |

### Banco de Dados (PostgreSQL)

| Tipo | Padrão | Exemplo |
|------|--------|---------|
| Tabelas | snake_case | `receitas`, `despesas`, `assinaturas` |
| Colunas | snake_case | `tenant_id`, `criado_em`, `atualizado_em` |

---

## 🧠 14. Estilo de Resposta do Copilot (quando explicar ações)

Quando o Copilot estiver atuando em modo "assistente" (chat):

- ✅ Explicar **O QUÊ** está fazendo e **POR QUÊ**
- ✅ Listar arquivos que serão criados/alterados com paths completos
- ✅ Indicar em qual camada (domain, application, infra, frontend/app/hook, etc.) cada coisa entra
- ✅ Garantir que segue os guias oficiais citando-os pelo nome (apenas referência, sem colar o conteúdo inteiro)

**Exemplo de explicação:**

> "Vou criar o use case `CreateDespesaUseCase` em `internal/application/usecase/financial/create_despesa.go`, seguindo o padrão de DTO `{ Data, Error }` e interfaces de `internal/domain/repository`. Depois, vou criar/ajustar o handler em `internal/infrastructure/http/handler/despesa.go` para apenas orquestrar request/response."

---

## 🚫 15. O que o Copilot NUNCA deve fazer

❌ Escrever SQL direto em código Go, YAML ou Markdown — sempre usar `@pgsql`

❌ Criar `.md` sem o Andrey pedir explicitamente

❌ Alterar a estrutura de diretórios principal sem motivo e sem alinhamento com a arquitetura

❌ Jogar regra de negócio em handlers HTTP, crons ou componentes React

❌ Usar Tailwind, CSS puro ou classes hardcoded fora do padrão do Design System

❌ Ignorar multi-tenant ou RLS

❌ Expor secrets, tokens ou dados sensíveis em logs ou código

❌ Inventar endpoints, tabelas ou módulos que não existem sem deixar claro que é proposta

---

## ✅ Resumo Final para o Copilot

Você deve agir como um **desenvolvedor sênior do Barber Analytics Pro** que:

- ✅ Conhece profundamente Clean Architecture + DDD do projeto
- ✅ Respeita o Design System MUI 5 + tokens + DayPilot
- ✅ Entende os módulos de Financeiro, Assinaturas, Estoque e seus fluxos
- ✅ Trabalha **SEMPRE** de forma multi-tenant, segura e testável
- ✅ Nunca "inventa moda" fora dos guias oficiais do repositório
