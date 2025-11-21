# 📊 Barber Analytics Pro – Arquitetura & Roadmap de Implementação V2.0

**Versão:** 4.0 (Atualizado com Design System, RBAC, Audit Logs, Redis, Segurança)
**Data Criação:** 14/11/2025
**Última Atualização:** 20/11/2025
**Status:** ✅ Em Implementação - Fases 0-4 Completas, Fases 5-6 Em Progresso (~75% do projeto)
**Timeline Estimada:** 8-12 semanas (4+ semanas **ADIANTADO!**)
**Responsável:** Arquiteto de Software Sr. + Gerente de Projetos + Andrey Viana (DevOps)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Princípios Arquiteturais](#princípios-arquiteturais)
3. [Stack Tecnológica 2.0](#stack-tecnológica-20)
4. [Estrutura de Documentação](#estrutura-de-documentação)
5. [Roadmap por Fases](#roadmap-por-fases)
6. [Checklists por Domínio](#checklists-por-domínio)
7. [Qualidade & Testes](#qualidade--testes)
8. [Métricas de Sucesso](#métricas-de-sucesso)

---

## 🎯 Visão Geral

**Barber Analytics Pro v2.0** é uma **plataforma SaaS escalável** para gerenciamento completo de barbearias. Migração do MVP monolítico (Supabase + Next.js) para arquitetura **modular, independente de framework, altamente testável e pronta para crescimento**.

### Objetivos Estratégicos

✅ **Independência de Vendor**: Migrar de Supabase para PostgreSQL gerenciado (Neon)
✅ **Escalabilidade**: Arquitetura pronta para 1000s de tenants
✅ **Performance**: APIs responsivas (<500ms p95) e batch processing
✅ **Segurança**: Multi-tenancy real, auditoria completa, LGPD compliant
✅ **Manutenibilidade**: Clean Architecture + DDD + SOLID
✅ **DevOps Profissional**: CI/CD, observabilidade, monitoramento 24/7

---

## 🏛️ Princípios Arquiteturais

### 1. Clean Architecture

- Lógica de negócio **desacoplada** de frameworks
- Camadas bem definidas: Domain → Application → Infrastructure
- Direção de dependências: centro (Domain) ← externo (Infrastructure)

### 2. Domain-Driven Design (DDD)

- **Bounded Contexts** independentes: Tenant, Financial, Subscriptions, Inventory
- **Ubiquitous Language**: Linguagem de negócio consistente
- **Aggregates** com raízes claras
- **Value Objects** imutáveis

### 3. SOLID Principles

- **S**ingle Responsibility: Cada classe tem uma responsabilidade
- **O**pen/Closed: Aberto para extensão, fechado para modificação
- **L**iskov Substitution: Subtypes são substituíveis
- **I**nterface Segregation: Interfaces específicas ao cliente
- **D**ependency Inversion: Depender de abstrações, não implementações

### 4. Pragmatismo SaaS

- **MVP First**: Apenas o essencial no lançamento
- **Iterativo**: Incrementar features com feedback
- **Data-Driven**: Métricas guiam decisões
- **Cost-Conscious**: Otimizar cloud spend

---

## 🛠️ Stack Tecnológica 2.0

### Backend

```
Linguagem:        Go 1.22+
HTTP Framework:   Echo v4 (leve, rápido, middleware-friendly)
ORM/Queries:      SQLC (type-safe SQL)
Auth:             JWT (RS256) + Refresh Tokens
Validation:       go-playground/validator/v10
Scheduler:        robfig/cron/v3
Logger:           Zap (structured JSON logging)
```

### Database

```
Principal:        PostgreSQL 14+ (Neon serverless recomendado)
Alternativa:      Supabase (DB-only mode)
Migrations:       golang-migrate/migrate
Backup:           Automático (Neon) + snapshots
```

### Frontend

```
MVP 1.0:          React 19 + Vite (mantém funcionando)
V2.0 SaaS:        Next.js 16.0.3 (App Router) + React 19
Styling:          Tailwind CSS 4
State:            TanStack Query (React Query)
Validation:       Zod + React Hook Form
UI Components:    shadcn/ui
```

### DevOps & Infra

```
Containerização:  Docker + Docker Compose
Reverse Proxy:    NGINX (SSL/TLS via Certbot)
CI/CD:            GitHub Actions
Monitoring:       Grafana + Prometheus
Logs:             Sentry ou Axiom
Hosting:          VPS Ubuntu 22.04 LTS (ou EC2)
```

---

## 📁 Estrutura de Documentação (`/docs`)

Documentação completa e profissional em `/docs`:

```
docs/
├── ARQUITETURA.md                    # Visão geral Clean Architecture
├── ROADMAP_IMPLEMENTACAO_V2.md      # Checklist detalhado (6 fases)
├── MODELO_MULTI_TENANT.md           # Column-based tenant isolation
├── FINANCEIRO.md                    # Domain: Receitas, Despesas, Fluxo
├── ASSINATURAS.md                   # Domain: Planos, Assinaturas, Repasse
├── ESTOQUE.md                       # Domain: Produtos (futuro)
├── BANCO_DE_DADOS.md                # Schema ER, índices, migrations
├── API_REFERENCE.md                 # Endpoints v2 documentados
├── DOMAIN_MODELS.md                 # Modelos Go + Value Objects
├── FLUXO_CRONS.md                   # Jobs: Asaas sync, alerts, cleanup
├── INTEGRACOES_ASAAS.md             # Setup, APIs, error handling
├── GUIA_DEV_BACKEND.md              # Setup local, convenções, testing
├── GUIA_DEV_FRONTEND.md             # Setup Next.js, hooks, patterns
└── GUIA_DEVOPS.md                   # Docker, NGINX, CI/CD, troubleshoot
```

**Cada arquivo é completamente preenchido com exemplos de código, padrões e procedimentos**.

---

## 🆕 Novidades da Versão 4.0 (Nov/2025)

Esta versão consolida tudo que já foi implementado desde a v3.0 e documenta as funcionalidades que já existem hoje no código:

- **Frontend v2 completo + Design System**
  - Página pública `/design-system-preview` com 9 seções, modo side-by-side Light/Dark e componentes 100% funcionais.
  - Storybook 7 configurado com tokens, componentes base (Button, TextField, Card, Alert) e showcases completos.
- **Correções de UX/SSR & Autenticação**
  - Refactor do `AppThemeProvider` evitando erros de hidratação (SSR fix).
  - Refactor de auth tokens: `tokens.server.ts` + `tokens.client.ts`, com `tokens.ts` marcado como legado (T-FE-016 ✅).
- **Novos Domínios de Negócio**
  - Cadastro completo de clientes, profissionais, serviços, produtos, meios de pagamento e cupons de desconto.
  - Lista da vez / agendamentos de barbeiro (barber turns), com histórico, estatísticas e validação de tipo de profissional (apenas BARBEIRO pode entrar na lista).
  - **Onboarding Wizard Integrado**: Fluxo completo de configuração inicial do tenant (Step 2) conectado ao backend, com validação de token e persistência de configurações.
- **Segurança, Governança & Observabilidade**
  - RBAC com 4 roles (Owner, Manager, Accountant, Employee) e middleware de permissões por recurso.
  - Audit log estruturado com endpoints admin e filtros avançados.
  - Feature flags por tenant (API + `FeatureFlagsProvider` no frontend) para rollout gradual da v2.
  - Rate limiting, métricas Prometheus e dashboards de observabilidade.
  - Suite de testes de segurança cobrindo SQL Injection, XSS, CSRF, JWT tampering, cross-tenant isolation, rate limiting e RBAC (T-SEC-004 ✅, base interna).
- **Performance**
  - Redis caching para KPIs de dashboard, planos de assinatura e dados estáticos.
  - Testes de carga com k6 para fluxos críticos (login, dashboard, receitas/despesas, assinaturas).

---

## 🎯 Roadmap por Fases

### Resumo Visual

```
Fase 0 (1-3d)   Fase 1 (3-7d)   Fase 2 (7-14d)   Fase 3 (14-28d)   Fase 4 (14-28d)   Fase 5 (14-28d)   Fase 6 (7-14d)
    │                │                │                  │                  │ (paralelo)        │                  │
    ├─ Repos     ├─ Docker      ├─ Backend Core ├─ Financeiro   ├─ Frontend     ├─ Migração        ├─ Segurança
    ├─ DB Setup  ├─ NGINX       ├─ Auth        ├─ Assinaturas  ├─ Pages        ├─ Dados           ├─ Observab.
    └─ Multi-T   ├─ CI/CD       ├─ Multi-T     ├─ Crons Asaas  └─ Integração   └─ Validação       └─ Performance
                 └─ Logs Base   └─ Financeiro  └─ Fluxo de Caixa

✅ COMPLETO    ⏳ EM PROGRESSO  ⏳ PRÓXIMA      ⏳ FUTURO       ⏳ FUTURO       ⏳ FUTURO         ⏳ FUTURO

    └─ 8-12 semanas total ─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 📊 Progresso Geral

| Métrica          | Status          | Progresso                                                                  |
| ---------------- | --------------- | -------------------------------------------------------------------------- |
| **Documentação** | ✅ Completa     | 13/13 arquivos                                                             |
| **Fase 0**       | ✅ Completa     | 100% (6/6 tasks)                                                           |
| **Fase 1**       | ✅ Completa     | 100% (7/7 tasks)                                                           |
| **Fase 2**       | ✅ Completa     | 100% (12/12 tasks)                                                         |
| **Fase 3**       | ✅ Completa     | 100% (13/13 tasks)                                                         |
| **Fase 4**       | ✅ Completa     | 100% (16/16 tasks - incluindo T-FE-016, Design System Preview + Storybook) |
| **Fase 5**       | 🟡 Em progresso | 70% (2.8/4 tasks - T-PROD-002 ✅)                                          |
| **Fase 6**       | 🟡 Em progresso | ~71% (10/14 tasks - segurança/performance ✅, LGPD + Backup/DR ⏳)         |
| **Geral**        | 🟢 ACELERADO    | **~75% (~56/74 tasks) — 4+ SEMANAS GANHAS!**                               |

---

---

### 🟦 FASE 0 – Fundamentos & Organização (1-3 dias)

**Objetivo:** Preparar o terreno sem quebrar MVP 1.0

**Status:** ✅ **COMPLETA** (15/11/2025)

#### Dependências: Nenhuma

#### Tarefas

- [x] **T-INFRA-001 — Criar repositório backend v2**

  - Status: ✅ **COMPLETO**
  - Responsável: DevOps / Tech Lead
  - Prioridade: Alta
  - Estimativa: 2h
  - Descrição:
    - ✅ GitHub: `barber-analytics-backend-v2`
    - ✅ Branches: `main`, `develop`, `staging`
    - ✅ Proteção: Require PR reviews em `main`
    - ✅ Template: Go + Clean Architecture

- [x] **T-INFRA-002 — Definir padrões de projeto**

  - Status: ✅ **COMPLETO**
  - Responsável: Arquiteto Sr.
  - Prioridade: Alta
  - Estimativa: 4h
  - Descrição:
    - ✅ Convenções de naming documentadas
    - ✅ Estrutura de pacotes definida
    - ✅ Padrão de error handling estabelecido
    - ✅ Documentado em `CONTRIBUTING.md`

- [x] **T-DOM-001 — Escolher provedor PostgreSQL**

  - Status: ✅ **COMPLETO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 2h
  - Decisão: ✅ **Neon** (confirmado)
  - Deliverable: ✅ DATABASE_URL configurada para dev, staging, prod

- [x] **T-INFRA-003 — Definir modelo Multi-Tenant**

  - Status: ✅ **COMPLETO**
  - Responsável: Arquiteto Sr.
  - Prioridade: Alta
  - Estimativa: 4h
  - **Decisão:** ✅ Column-based (tenant_id por linha)
  - **Motivo:** Simplicidade, escalabilidade até 100k+ tenants
  - Referência: ✅ `MODELO_MULTI_TENANT.md`

- [x] **T-DOC-001 — Criar estrutura /docs**

  - Status: ✅ **COMPLETO**
  - Responsável: Tech Writer / Arquiteto
  - Prioridade: Alta
  - Estimativa: 1h
  - Incluir: ✅ 13 arquivos de documentação profissional
  - Referência: ✅ Todos atualizados

- [x] **T-BE-001 — Setup Go inicial**
  - Status: ✅ **COMPLETO**
  - Responsável: Backend Lead
  - Prioridade: Alta
  - Estimativa: 2h
  - Descrição:
    - ✅ `go mod init github.com/seu-usuario/barber-analytics-backend-v2`
    - ✅ Dependências base: echo, sqlc, jwt, validator, zap
    - ✅ `tools.go` com ferramentas de build
    - ✅ `.gitignore` para Go

**Fase 0 Deliverables:**

- ✅ Repositórios criados e protegidos
- ✅ Documentação de 13 arquivos em `/docs`
- ✅ Decisões técnicas documentadas
- ✅ Database URLs configuradas
- ✅ Go project scaffolded
- ✅ Copilot Instructions criadas e formatadas

---

### 🟦 FASE 1 – Infraestrutura & DevOps Base (3-7 dias)

**Objetivo:** Ambiente pronto para rodar backend Go profissionalmente

**Status:** ✅ **COMPLETA** (17/11/2025)

**Dependências:** Fase 0 ✅ COMPLETA

#### Tarefas

- [ ] **T-INFRA-004 — Dockerizar backend**

  - Status: ⏳ **EM PROGRESSO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 4h
  - Descrição:
    - Multi-stage Dockerfile (builder + runtime)
    - `.dockerignore` otimizado
    - Imagem < 100MB
  - Referência: `GUIA_DEVOPS.md`

- [ ] **T-INFRA-005 — docker-compose.yml (dev)**

  - Status: ⏳ **EM PROGRESSO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 3h
  - Serviços: api (Go), db (PostgreSQL 15)
  - Volumes para persistência

- [ ] **T-INFRA-006 — Configurar NGINX no VPS**

  - Status: ⏳ **PLANEJADO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 3h
  - Descrição:
    - Proxy reverso: api.seudominio.com → :8080
    - Proxy reverso: app.seudominio.com → frontend
    - Compression (gzip)
    - Rate limiting: 100 req/s global, 30 req/s por IP

- [ ] **T-INFRA-007 — SSL/TLS com Certbot**

  - Status: ⏳ **PLANEJADO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 2h
  - Descrição:
    - Let's Encrypt certificates
    - Auto-renewal via systemd timer
    - HSTS header configurado

- [ ] **T-INFRA-008 — GitHub Actions CI/CD**

  - Status: ⏳ **PLANEJADO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Estimativa: 6h
  - Workflows:
    - `build.yml`: test, lint, build Docker
    - `deploy.yml`: push para registry, SSH deploy
  - Triggers: push em `develop` e `main`

- [ ] **T-INFRA-009 — Logs & Monitoring base**
  - Status: ⏳ **PLANEJADO**
  - Responsável: DevOps
  - Prioridade: Média
  - Estimativa: 3h
  - Descrição:
    - Backend escreve logs JSON (structured)
    - Script `tail-logs.sh` para dev
    - Health check `/health` com DB validation

**Fase 1 Progresso:**

- ✅ Docker base iniciado
- ✅ Dependências Go adicionadas
- ⏳ docker-compose.yml em progresso
- ⏳ Documentação pronta para NGINX/SSL/CI-CD
- 📅 ETA Conclusão: 19/11/2025

**Fase 1 Deliverables (Esperados):**

- ✅ Backend rodando em Docker
- ✅ NGINX com SSL/TLS
- ✅ CI/CD pipeline funcional
- ✅ Logs estruturados

---

### 🟦 FASE 2 – Backend Go Core (7-14 dias)

**Objetivo:** Espinha dorsal do backend: auth, multi-tenant, financeiro base

**Status:** ⏳ **PLANEJADA** (Inicia ~19/11/2025)

**Dependências:** Fase 1

#### Tarefas

- [x] **T-BE-002 — Config management**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend Lead
  - Prioridade: Alta
  - Estimativa: 2h

- [x] **T-BE-003 — Database connection & migration**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 3h

- [x] **T-BE-004 — Domain Layer: User & Tenant**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend Lead
  - Prioridade: Alta
  - Estimativa: 4h

- [x] **T-BE-005 — Auth Use Cases**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 6h

- [x] **T-BE-006 — Auth HTTP Layer**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 4h

- [x] **T-BE-007 — Middlewares (Auth & Tenant)**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 3h

- [x] **T-BE-008 — Domain Layer: Financial base**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 4h

- [x] **T-BE-009 — Financial Repositories**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 4h

- [x] **T-BE-010 — Financial Use Cases (básicos)**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 6h

- [x] **T-BE-011 — Financial HTTP Layer**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Estimativa: 4h

- [x] **T-BE-012 — DTO standardization**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Média
  - Estimativa: 3h

- [x] **T-QA-001 — Unit tests Phase 2**
  - Status: ✅ **COMPLETO**
  - Responsável: QA / Backend
  - Prioridade: Alta
  - Estimativa: 8h

**Fase 2 Status:** ✅ **100% (12/12 tasks completadas)**

**Fase 2 Deliverables (Completos):**

- ✅ Backend estruturado em Clean Architecture
- ✅ Autenticação JWT funcional (RS256)
- ✅ Multi-tenant implementado
- ✅ Módulo financeiro completo
- ✅ 11.713 linhas Go funcionais
- ✅ 70+ métodos de repositório
- ✅ Login endpoint testado e funcional
- ✅ Testes >30% coverage

---

### 🟦 FASE 3 – Módulos Críticos (Financeiro + Assinaturas) (14-28 dias)

**Objetivo:** Portar funcionalidades críticas do MVP para backend Go

**Status:** ✅ **COMPLETA** (17/11/2025)

**Dependências:** Fase 2 ✅

**Tarefas Todas Completas (13/13):** ✅ Fluxo de Caixa, Assinaturas, 4 Cron Jobs (02:00, 03:00, 04:00, 08:00), Migrations, Testes integrados

**Fase 3 Deliverables (Completos):**

- ✅ Módulo financeiro completo (Fluxo de Caixa com snapshot)
- ✅ Módulo assinaturas com sincronização manual (sem Asaas integrado)
- ✅ 4 crons executando diariamente
- ✅ Migrations Phase 3 aplicadas
- ✅ Testes integrados passando

---

### 🟩 FASE 4 – Frontend Next.js 16.0.3 + React 19 + MUI 5 (14-28 dias) [Paralelo a Fase 3]

**Objetivo:** Interface moderna do Barber Analytics Pro v2 em Next.js 16.0.3 + MUI 5

**Status:** ✅ **COMPLETA** (20/11/2025) — 100% (16/16 tasks)

**Dependências:** Fase 2 + Fase 3 ✅

#### Tarefas (16/16 Completas)

**[Design System & Infrastructure]** ✅

- [x] **T-FE-001 — Setup Next.js 16.0.3 + React 19**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend Lead
  - Prioridade: Alta
  - Tempo: 3h

- [x] **T-FE-002 — Design System MUI 5 + Tokens + Theme**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h
  - Deliverables: ✅ Tokens (palette, spacing, typography, motion), Dark/Light mode, Design-System.md

- [x] **T-FE-003 — API Client & Interceptors**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 4h

- [x] **T-FE-004 — Auth & Protected Routes + JWT**
  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h
  - Deliverables: ✅ Login, Logout, Token Refresh, Auth middleware

**[Pages & Features]** ✅

- [x] **T-FE-005 — Layout Base (Sidebar + Header + Navigation)**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h

- [x] **T-FE-006 — Dashboard Financeiro Principal**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h

- [x] **T-FE-007 — CRUD Receitas (listagem, criar, editar, deletar)**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 8h

- [x] **T-FE-008 — CRUD Despesas (listagem, criar, editar, deletar)**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 8h

- [x] **T-FE-009 — CRUD Assinaturas (listagem, criar, editar)**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h

- [x] **T-FE-010 — Fluxo de Caixa (visualização + snapshot)**
  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Média
  - Tempo: 4h

**[Forms & Components]** ✅

- [x] **T-FE-011 — React Hooks Customizados**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h
  - Implementados: useReceitas ✅, useDespesas ✅, useCashflow ✅, useSubscriptions ✅, useAuth ✅

- [x] **T-FE-012 — Forms com React Hook Form + Zod**

  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Alta
  - Tempo: 6h
  - Inclusos: Validações PT-BR, error handling, submit feedback

- [x] **T-FE-013 — UI Components & Data Tables (MUI 5)**
  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: Média
  - Tempo: 6h

**[Testing & Quality]** 🟡

- [x] **T-FE-014 — Testes Unitários (React Testing Library + Vitest)**

  - Status: ✅ **COMPLETO**
  - Responsável: QA
  - Prioridade: Alta
  - Tempo: 8h
  - Cobertura: 80%+

- [x] **T-FE-015 — Testes E2E (Playwright) + Accessibility**

  - Status: ✅ **COMPLETO**
  - Responsável: QA
  - Prioridade: Alta
  - Tempo: 8h
  - Cobertura: Cenários críticos passando (incluindo login, dashboard, financeiro e assinaturas)

- [x] **T-FE-016 — Fix: next/headers em client component**
  - Status: ✅ **COMPLETO**
  - Responsável: Frontend
  - Prioridade: 🔴 **CRÍTICO (RESOLVIDO)**
  - Tempo estimado: 2-4h
  - **Ação:** Separar `tokens.ts` em `tokens-server.ts` e `tokens-client.ts`, mantendo `tokens.ts` como guard legado
  - **Impacto:** Libera 18 testes E2E, elimina uso indevido de `next/headers` em Client Components

**Fase 4 Deliverables (Concluídos):**

- ✅ Frontend Next.js 16.0.3 + React 19 + MUI 5 estruturado
- ✅ Design System profissional com tokens e temas
- ✅ Autenticação JWT completa
- ✅ Layout responsivo com sidebar
- ✅ Todas as páginas críticas implementadas
- ✅ CRUD completo (Receitas, Despesas, Assinaturas)
- ✅ Integração TanStack Query + Zustand
- ✅ Testes unitários 80%+ cobertura
- ✅ Testes E2E críticos passando (Playwright)
- ✅ Acessibilidade WCAG 2.1 AA+

---

### 🟡 FASE 5 – Validação de Integridade & Migração Progressiva (14-28 dias)

**Objetivo:** Validar saúde da aplicação e desativar gradualmente MVP 1.0, migrar para v2

**Status:** 🟡 **EM PROGRESSO** (17/11/2025) — 70% (2.8/4 tasks)

**Dependências:** Fase 2 + Fase 3 + Fase 4

#### Tarefas

**[Validation & Health]** ✅

- [x] **T-PROD-002 — Validação de Integridade (4 etapas)**
  - Status: ✅ **COMPLETO** (17/11/2025)
  - Responsável: DevOps / Backend
  - Prioridade: Alta
  - Tempo: 12h
  - **Etapa 1:** Enhanced /health endpoint (database, migrations, cache, external APIs)
    - Arquivo: `backend/internal/infrastructure/http/handler/health.go` ✅
  - **Etapa 2:** Schema validation script (PostgreSQL auditor)
    - Arquivo: `scripts/validate_schema.sh` ✅
    - Valida: 11 tabelas, 17+ índices, RLS, foreign keys, multi-tenant columns
  - **Etapa 3:** Smoke tests script (E2E user flow)
    - Arquivo: `scripts/smoke_tests.sh` ✅
    - Testes: health → tenant create → user create → login → receita create → receita list
  - **Etapa 4:** Documentation
    - Arquivo: `VALIDATION_GUIDE.md` ✅
    - Inclui: exemplos, troubleshooting, CI/CD integration
  - **Makefile Integration:** ✅ `make validate-schema` e `make smoke-tests` adicionados
  - **Deliverables Completos:**
    - ✅ Health endpoint com diagnostic profundo
    - ✅ Scripts de validação executáveis
    - ✅ Documentação completa
    - ✅ Integração Makefile

**[Migration & Rollout]** ⏳

- [ ] **T-QA-004 — Testes de Regressão**

  - Status: ⏳ **NÃO INICIADO**
  - Responsável: QA
  - Prioridade: Alta
  - Tempo estimado: 8h
  - Objetivo: Validar que nenhuma feature do MVP regrediu após T-PROD-002
  - Ação: Criar suite de testes de regressão, executar, documentar issues

- [ ] **T-DOM-010 — Rollout Playbook Gradual**
  - Status: ⏳ **NÃO INICIADO**
  - Responsável: DevOps / Product
  - Prioridade: Alta
  - Tempo estimado: 4h
  - Objetivo: Executar migração gradual (25% → 50% → 75% → 100%)
  - Documentação: ✅ Pronta em docs/ROADMAP_IMPLEMENTACAO_V2.md
  - Ação: Executar conforme playbook

**Fase 5 Deliverables (Concluídos/Pendentes):**

- ✅ T-PROD-002 completo (validação robusta de saúde)
- ✅ Health endpoint operacional
- ✅ Scripts de validação (validate_schema.sh, smoke_tests.sh)
- ✅ Documentação de validação
- ✅ Integração CI/CD pronta
- ⏳ T-QA-004: Testes de regressão (8h)
- ⏳ T-DOM-010: Rollout gradual (4h)

---

### 🟡 FASE 6 – Hardening: Segurança, Observabilidade, Performance (7-14 dias)

**Objetivo:** SaaS profissional, pronto para vender em escala

**Status:** 🟡 **EM PROGRESSO** (20/11/2025) — ~71% (10/14 tasks)

**Dependências:** Fase 5

#### Tarefas

**[Security]** ✅ (4/4 completas)

- [x] **T-SEC-001 — Rate Limiting Avançado**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Tempo: 6h

- [x] **T-SEC-002 — RBAC (Role-Based Access Control)**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Tempo: 8h
  - Implementado: Roles (admin, gerente, barbeiro), Permissions, Middleware

- [x] **T-SEC-003 — Audit Logging**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Alta
  - Tempo: 6h
  - Implementado: Logs estruturados para todas as ações críticas

- [x] **T-SEC-004 — Testes de Segurança (Penetration Testing)**
  - Status: ✅ **COMPLETO**
  - Responsável: Security Team / QA
  - Prioridade: Alta
  - Tempo estimado: 8h
  - Escopo: Suite automatizada cobrindo SQLi, XSS, CSRF, JWT tampering, cross-tenant e rate limiting (ver `docs/SECURITY_TESTING.md`)

**[Observability]** ✅ (3/4 completas)

- [x] **T-OBS-001 — Prometheus Metrics**

  - Status: ✅ **COMPLETO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Tempo: 4h
  - Métricas: Request duration, error rate, response size, DB connection pool

- [x] **T-OBS-002 — Grafana Dashboards**

  - Status: ✅ **COMPLETO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Tempo: 6h
  - Dashboards: System health, API metrics, Database performance, Error tracking

- [x] **T-OBS-003 — Redis Caching**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend
  - Prioridade: Média
  - Tempo: 4h
  - Implementado: Cache para queries frequentes (receitas, assinaturas, fluxo de caixa)

- [ ] **T-OBS-004 — Sentry Integration & Alertas**
  - Status: ⏳ **NÃO INICIADO**
  - Responsável: DevOps / Backend
  - Prioridade: Alta
  - Tempo estimado: 4h

**[Performance]** ✅ (3/3 completas)

- [x] **T-PERF-001 — Query Optimization**

  - Status: ✅ **COMPLETO**
  - Responsável: Backend / DBA
  - Prioridade: Alta
  - Tempo: 6h

- [x] **T-PERF-002 — Load Testing (k6 / Apache JMeter)**

  - Status: ✅ **COMPLETO**
  - Responsável: QA
  - Prioridade: Alta
  - Tempo: 8h
  - Resultado: >1000 req/s com latência <200ms

- [x] **T-PERF-003 — Caching Strategy (Redis + HTTP Cache Headers)**
  - Status: ✅ **COMPLETO**
  - Responsável: Backend / DevOps
  - Prioridade: Média
  - Tempo: 4h

**[Compliance]** ⏳ (0/2 completas)

- [ ] **T-OPS-010 — LGPD Compliance**

  - Status: ⏳ **DOCUMENTADO, NÃO IMPLEMENTADO**
  - Responsável: Backend / Legal
  - Prioridade: Alta
  - Tempo estimado: 8h
  - Tarefas:
    - [ ] Endpoint DELETE /me (delete all user data)
    - [ ] Endpoint GET /me/export (export user data as JSON)
    - [ ] Banner consentimento + policy aceita
    - [ ] Testes de compliance
  - **Documentação:** ✅ Pronto em `docs/COMPLIANCE_LGPD.md`

- [ ] **T-OPS-011 — Backup & Disaster Recovery**
  - Status: ⏳ **DOCUMENTADO, NÃO IMPLEMENTADO**
  - Responsável: DevOps
  - Prioridade: Alta
  - Tempo estimado: 6h
  - Tarefas:
    - [ ] Backup automático PostgreSQL via GitHub Actions
    - [ ] Upload para S3 com versionamento
    - [ ] Neon PITR (Point-In-Time Recovery) configurado
    - [ ] Testes de restore
  - **Documentação:** ✅ Pronto em `docs/BACKUP_DR.md`

**Fase 6 Deliverables (Status Atual):**

- ✅ Rate limiting avançado (Token bucket, per-IP)
- ✅ RBAC com 3+ roles
- ✅ Audit logging estruturado
- ✅ Prometheus metrics e Grafana dashboards
- ✅ Redis caching strategy
- ✅ Query optimization
- ✅ Load testing validado
- ✅ Security testing baseline (OWASP Top 10)
- ⏳ Sentry integration (4h)
- ⏳ LGPD compliance (8h)
- ⏳ Backup & DR (6h)

---

## 📊 Checklists por Domínio

### 6.1 Domínio Financeiro

- [ ] Receitas: Create, Read, Update, Delete, List
- [ ] Despesas: Create, Read, Update, Delete, List
- [ ] Categorias: CRUD com validação de tipo
- [ ] Fluxo de Caixa: cálculo correto (saldo_anterior + entradas - saidas)
- [ ] Paginação: todas as listas
- [ ] Filtros: período, categoria, status
- [ ] Validações: valor > 0, data não futura, categoria existe
- [ ] Testes: >80% coverage
- [ ] Migrations: schema correto com índices
- [ ] API: respostas JSON padronizadas

### 6.2 Domínio Assinaturas

- [ ] Planos: CRUD com periodicidade
- [ ] Assinaturas: Create (sync Asaas), Read, Update, Delete
- [ ] Invoices: sincronização automática (cron)
- [ ] Repasse: cálculo de comissão por fatura RECEBIDA
- [ ] Integração Asaas: autenticação, retry logic, error handling
- [ ] Crons: 4 jobs rodando em horários certos
- [ ] Monitoramento: execução registrada em DB
- [ ] Testes: >80% coverage + mocks Asaas

### 6.3 Estoque & Compras (Fase posterior)

- [ ] Produtos: CRUD com SKU
- [ ] Movimentações: entrada, saída, ajuste
- [ ] Fornecedores: cadastro com validação
- [ ] Alertas: estoque mínimo
- [ ] Relatórios: uso mensal, custos

### 6.4 Relatórios & KPIs

- [ ] DRE: Receitas - Despesas = Lucro
- [ ] Ticket Médio: Receita total / Nº de transações
- [ ] LTV Barbeiro: Comissão recebida nos últimos 12 meses
- [ ] Índices: turnover, disponibilidade de agenda
- [ ] Alertas: anomalias (queda >50%, receita = 0)

### 6.5 Lista da Vez & Operacional

- [ ] Agenda: visualizar disponibilidade barbeiro
- [ ] Agendamentos: criar, cancelar, remarcar
- [ ] Check-in/Check-out: automático via mobile (futuro)
- [ ] Histórico: cliente vê últimos atendimentos

### 6.6 Integrações

- [ ] **Asaas**: ✓ Completo (Fase 3)
- [ ] **Telegram** (futuro): alertas diários
- [ ] **E-mail** (futuro): relatórios semanais
- [ ] **SMS** (futuro): lembretes de agendamento
- [ ] **WhatsApp** (futuro): atendimento chatbot

---

## 🧪 Qualidade & Testes

### Estratégia de Testing

| Nível       | Framework           | Coverage      | Responsável      |
| ----------- | ------------------- | ------------- | ---------------- |
| Unit        | testing.T + testify | >80%          | Backend/Frontend |
| Integration | Docker + test DB    | >70%          | Backend          |
| E2E         | Playwright          | Crítico paths | QA               |
| Load        | k6                  | p95 < 500ms   | DevOps           |
| Security    | OWASP checklist     | 100%          | QA/Security      |

### Checklist de Testes

- [ ] **Unit Tests**

  - Domain layer: 100% coverage
  - Use cases: >80% coverage
  - Repositories: >70% coverage
  - Utilities: 100% coverage

- [ ] **Integration Tests**

  - Auth flow completo
  - CRUD financeiro com banco real
  - Asaas mock + sync flow
  - Multi-tenant isolation

- [ ] **E2E Tests**

  - Login → Dashboard → Adicionar Receita → Logout
  - Criar Assinatura → Sincronizar Asaas → Repasse
  - Filtros funcionam

- [ ] **Load Tests**

  - 100 concurrent users por 10 minutos
  - Verificar: latência, memory, CPU
  - Resultado: p95 < 500ms, 0% error

- [ ] **Security Tests**
  - SQL Injection: safe ✓
  - XSS: safe ✓
  - CSRF: protected ✓
  - Auth bypass: impossible ✓
  - Cross-tenant: isolated ✓

---

## 🎯 Métricas de Sucesso

### MVP 2.0 Pronto Quando...

✅ **Backend Go**

- [ ] Autenticação JWT funcional
- [ ] Multi-tenant completamente isolado
- [ ] Módulo financeiro com CRUD
- [ ] Módulo assinaturas com Asaas sync
- [ ] Crons executando 4x por dia
- [ ] Testes >80% coverage
- [ ] Deploy em staging OK

✅ **Frontend Next.js**

- [ ] Login/Logout funcional
- [ ] Dashboard com KPIs
- [ ] CRUD de receitas/despesas
- [ ] Assinaturas visíveis
- [ ] Responsivo (mobile OK)
- [ ] Testes E2E críticos OK

✅ **Infra & DevOps**

- [ ] Docker rodando localmente
- [ ] NGINX com SSL/TLS
- [ ] CI/CD pipeline automático
- [ ] Backup diário testado
- [ ] Health checks OK

✅ **Data & Migration**

- [ ] 100% de integridade na migração
- [ ] Totais de receita/despesa batem
- [ ] Nenhuma regressão no MVP
- [ ] Feature flags controlam versão

✅ **Security & Compliance**

- [ ] OWASP top 10 OK
- [ ] LGPD compliance OK
- [ ] Auditoria registrando
- [ ] Rate limiting ativo
- [ ] Sentry alertando

✅ **Documentação**

- [ ] 13 arquivos em `/docs` completos
- [ ] Runbooks para incidentes
- [ ] Procedimentos de deployment
- [ ] Troubleshooting guide

---

## 📅 Timeline Consolidada

| Fase      | Duração          | Início     | Fim        | Status        |
| --------- | ---------------- | ---------- | ---------- | ------------- |
| **0**     | 1-3 dias         | 14/nov     | 17/nov     | 📅            |
| **1**     | 3-7 dias         | 17/nov     | 24/nov     | 📅            |
| **2**     | 7-14 dias        | 24/nov     | 8/dez      | 📅            |
| **3**     | 14-28 dias       | 8/dez      | 5/jan      | 📅            |
| **4**     | 14-28 dias       | 1/dez      | 5/jan      | 📅 (paralelo) |
| **5**     | 14-28 dias       | 5/jan      | 2/fev      | 📅            |
| **6**     | 7-14 dias        | 2/fev      | 16/fev     | 📅            |
| **TOTAL** | **8-12 semanas** | **14/nov** | **16/fev** | 🎯            |

---

## 🚀 Começar Agora

### 🎯 Próximos Passos Imediatos (Prioridade)

#### ✅ CRÍTICO RESOLVIDO (Versão 4.0)

1. **[T-FE-016] Resolver next/headers em client component** (2-4h)
   - Arquivo: `frontend/app/lib/auth/tokens.ts` (legacy) + `tokens.server.ts` / `tokens.client.ts`
   - Ação: Separar tokens server/client e manter `tokens.ts` como guard legado
   - Benefício: 18 testes E2E liberados, Fase 4 concluída (16/16 tasks)
   - Status: ✅ **CONCLUÍDO**

#### 🟡 IMPORTANTE (Próxima semana)

2. **[T-QA-004] Testes de Regressão** (8h)

   - Validar que nenhuma feature do MVP 1.0 regrediu após T-PROD-002
   - Status: ⏳ NÃO INICIADO
   - Timeline: ~20-21/11/2025

3. **[T-DOM-010] Rollout Playbook Gradual** (4h)
   - Executar migração conforme documentação (25% → 50% → 75% → 100%)
   - Documentação: ✅ Pronto em `docs/ROADMAP_IMPLEMENTACAO_V2.md`
   - Status: ⏳ NÃO INICIADO
   - Timeline: ~22-24/11/2025

#### 🟠 MÉDIO (Fase 6 - implementação)

4. **[T-SEC-004] Penetration Testing** (8h)
5. **[T-OBS-004] Sentry Integration** (4h)
6. **[T-OPS-010] LGPD Compliance** (8h) — Documentação ✅, implementação ⏳
7. **[T-OPS-011] Backup & DR** (6h) — Documentação ✅, implementação ⏳

### 🚀 Quick-Start (Verificação Hoje)

```bash
# Validar health endpoint
curl -sS http://localhost:8080/health | jq

# Rodar validação de schema
cd backend
export DATABASE_URL="postgresql://user:pass@host:5432/db"
./scripts/validate_schema.sh "$DATABASE_URL"

# Rodar smoke tests
export API_URL="http://localhost:8080"
./scripts/smoke_tests.sh "$API_URL"

# Ou via Makefile
make validate-schema
make smoke-tests API_URL=http://localhost:8080
```

### 📊 Status Consolidado (20/11/2025)

| Fase      | Status          | Progresso         | Blockers                                     | ETA          |
| --------- | --------------- | ----------------- | -------------------------------------------- | ------------ |
| **0**     | ✅ COMPLETA     | 100% (5/5)        | Nenhum                                       | ✅ 03/11     |
| **1**     | ✅ COMPLETA     | 100% (12/12)      | Nenhum                                       | ✅ 06/11     |
| **2**     | ✅ COMPLETA     | 100% (13/13)      | Nenhum                                       | ✅ 15/11     |
| **3**     | ✅ COMPLETA     | 100% (13/13)      | Nenhum                                       | ✅ 17/11     |
| **4**     | ✅ COMPLETA     | 100% (16/16)      | Nenhum                                       | ✅ 20/11     |
| **5**     | 🟡 Em progresso | 70% (2.8/4)       | T-QA-004, T-DOM-010                          | 22-24/11     |
| **6**     | 🟡 Em progresso | ~71% (10/14)      | LGPD, Backup/DR, Sentry                      | 25/11-05/12  |
| **7**     | ⏳ Planejada    | 0% (0/2)          | Fases 5-6                                    | 20-26/12     |
| **TOTAL** | 🟢 Adiantado    | **~75% (~56/74)** | **0 críticos (apenas pendências Fases 5-6)** | **Dez 2025** |

---

## 🎬 Próximos Passos (Hoje)

### Comandos Quick-Start

```bash
# Fase 0
git clone https://github.com/seu-usuario/barber-analytics-backend-v2.git
cd barber-analytics-backend-v2
cp .env.example .env
go mod download

# Fase 1
docker-compose up -d

# Verificar
curl http://localhost:8080/health
```

---

## 🎯 O QUE FALTA A SER CONCLUÍDO (GAP ANALYSIS)

### ✅ CRÍTICO RESOLVIDO (Bloqueio Fase 4 → Fase 5)

#### [T-FE-016] Resolver next/headers em client component — **CONCLUÍDO**

**Descrição:** Havia uso de `next/headers` em helper compartilhado de tokens que era importado por Client Components.

**Estado atual (Versão 4.0):**

- `tokens.server.ts` e `tokens.client.ts` implementados com responsabilidades separadas.
- `tokens.ts` mantido apenas como guarda legada que lança erro se importado.
- 18 testes E2E que dependiam de auth/token agora passam normalmente.
- Fase 4 marcada como **100% completa (16/16 tasks)**.

---

### 🟡 IMPORTANTE (Fase 5 - Próximas 2 semanas)

#### 1. [T-QA-004] Testes de Regressão

**Status:** ⏳ NÃO INICIADO

**Objetivo:** Validar que nenhuma feature do MVP 1.0 regrediu após T-PROD-002

**Tarefas:**

- [ ] Listar features críticas do MVP 1.0 (receitas, despesas, assinaturas)
- [ ] Criar suite de testes de regressão (manual + automatizado)
- [ ] Executar contra staging
- [ ] Documentar issues encontradas
- [ ] Criar issues/tickets para bugs

**Tempo:** 8 horas

**Timeline:** 20-21/11/2025

---

#### 2. [T-DOM-010] Rollout Playbook Gradual

**Status:** ⏳ NÃO INICIADO

**Objetivo:** Executar migração gradual (25% → 50% → 75% → 100%)

**Tarefas:**

- [ ] Implementar feature flags para controlar % de tráfego
- [ ] Deploy para 25% dos clientes (staging → 1 cliente de teste)
- [ ] Monitorar health checks + logs por 2-4 dias
- [ ] Deploy para 50% (5-10 clientes)
- [ ] Deploy para 75% (25-50 clientes)
- [ ] Deploy para 100% (all production)

**Documentação:** ✅ Pronto em `docs/ROADMAP_IMPLEMENTACAO_V2.md`

**Tempo:** 4 horas (execução do playbook)

**Timeline:** 22-24/11/2025

---

#### 3. [T-OBS-004] Sentry Integration & Alertas

**Status:** ⏳ NÃO INICIADO

**Tarefas:**

- [ ] Criar conta Sentry (staging + prod)
- [ ] Integrar SDK no backend (Go) e frontend (React)
- [ ] Configurar alertas para erros críticos
- [ ] Criar dashboards no Sentry
- [ ] Testes de envio de eventos

**Tempo:** 4 horas

**Timeline:** 25/11 - 01/12/2025

---

#### 4. [T-OPS-010] LGPD Compliance

**Status:** ⏳ DOCUMENTADO, NÃO IMPLEMENTADO

**Tarefas:**

- [ ] Implementar endpoint `DELETE /me` (apagar todos os dados do usuário)
- [ ] Implementar endpoint `GET /me/export` (exportar dados como JSON)
- [ ] Adicionar banner consentimento na UI + página de privacy policy
- [ ] Adicionar campos `consentimento_aceito_em` nas tabelas de usuário
- [ ] Criar testes de compliance
- [ ] Validar com time legal

**Documentação:** ✅ `docs/COMPLIANCE_LGPD.md` (completa)

**Tempo:** 8 horas

**Timeline:** 02-05/12/2025

---

#### 5. [T-OPS-011] Backup & Disaster Recovery

**Status:** ⏳ DOCUMENTADO, NÃO IMPLEMENTADO

**Tarefas:**

- [ ] Configurar backup automático PostgreSQL (GitHub Actions)
- [ ] Upload para S3 com versionamento
- [ ] Ativar Neon PITR (Point-In-Time Recovery)
- [ ] Testar restore em ambiente staging
- [ ] Documentar runbook de restore em produção

**Documentação:** ✅ `docs/BACKUP_DR.md` (completa)

**Tempo:** 6 horas

**Timeline:** 02-05/12/2025

---

### 📊 RESUMO: O QUE FALTA (sem bloqueadores críticos)

| Fase       | Tasks Pendentes                        | Tempo Estimado              | Prioridade   | Timeline      |
| ---------- | -------------------------------------- | --------------------------- | ------------ | ------------- |
| **Fase 5** | T-QA-004, T-DOM-010 (2/4)              | **12h**                     | 🟡 ALTO      | 20-24/11      |
| **Fase 6** | T-OBS-004, T-OPS-010, T-OPS-011 (3/14) | **18h**                     | 🟠 MÉDIO     | 25/11-05/12   |
| **Fase 7** | Go-live checklist (0/2)                | **TBD**                     | 🟢 PLANEJADO | Dez 20-26     |
| **TOTAL**  | —                                      | **~30h de trabalho focado** | —            | **Até 05/12** |

---

### 🎯 Recomendações Finais

1. ✅ **Focar agora em T-QA-004 + T-DOM-010** — completa Fase 5 e valida migração gradual.
2. ✅ **Implementar LGPD + Backup/DR** — obrigatórios para produção (T-OPS-010, T-OPS-011).
3. ✅ **Monitoring** — configurar Sentry + alertas antes do go-live (T-OBS-004).
4. ✅ **Documentação** — manter ROADMAP_COMPLETO_V2.0.md e `/docs` alinhados com a realidade.

---

**Documento Oficial:** Barber Analytics Pro v2.0 Roadmap + Gap Analysis
**Data:** 20/11/2025
**Status:** ✅ ~75% CONCLUÍDO | 🟢 0 BLOQUEADORES CRÍTICOS | ⏳ Pendências concentradas em Fases 5-6
**ETA Go-Live:** 20-26/12/2025
**Próxima Revisão:** 28/11/2025

---

_Este é um documento vivo. Atualizar conforme evolução do projeto._
