# 📊 Barber Analytics Pro v2.0 — PRD (Product Requirements Document)

**Status:** ✅ Em Implementação (~75% Concluído)
**Versão:** 4.0 (Atualizado com Design System, RBAC, Audit Logs, Redis, Segurança)
**Data:** 20/11/2025

---

## ⚠️ DOCUMENTO IMPORTANTE

**Este é o PRD original refinado.** Para o **Roadmap completo, profissional e detalhado**, veja:

👉 **[`ROADMAP_COMPLETO_V2.0.md`](./ROADMAP_COMPLETO_V2.0.md)** ← **LEIA PRIMEIRO**

---

## 📁 Documentação Técnica Completa

Toda a documentação técnica está organizada em `/docs/`:

```
docs/
├── ARQUITETURA.md                    # Design: Clean Architecture + DDD
├── ROADMAP_IMPLEMENTACAO_V2.md      # Checklist detalhado 6 fases
├── MODELO_MULTI_TENANT.md           # Multi-tenancy column-based
├── FINANCEIRO.md                    # Domain: Receitas, Despesas, Fluxo
├── ASSINATURAS.md                   # Domain: Assinaturas + Asaas
├── ESTOQUE.md                       # Domain: Produtos (futuro)
├── BANCO_DE_DADOS.md                # Schema, índices, migrations
├── API_REFERENCE.md                 # Endpoints documentados
├── DOMAIN_MODELS.md                 # Go domain entities
├── FLUXO_CRONS.md                   # Scheduled jobs (4x diários)
├── INTEGRACOES_ASAAS.md             # Asaas API integration
├── GUIA_DEV_BACKEND.md              # Go setup + conventions
├── GUIA_DEV_FRONTEND.md             # Next.js setup + patterns
└── GUIA_DEVOPS.md                   # Docker, NGINX, CI/CD
```

---

## 🎯 Visão Geral Refinada

### O Que É Barber Analytics Pro v2.0?

Uma **plataforma SaaS escalável** para gerenciamento de barbearias com:

- ✅ Backend **Go 1.22+** (Clean Architecture + DDD)
- ✅ Frontend **Next.js 16.0.3** (React 19 + App Router)
- ✅ Database **PostgreSQL 14+** gerenciado (Neon)
- ✅ Infraestrutura **Docker + NGINX + CI/CD** profissional
- ✅ Multi-tenant **column-based** (isolamento garantido)
- ✅ Integração **Asaas** para assinaturas/repasse
- ✅ **Crons diários** para sincronização e processamento

### Stack Tecnológica

**Backend:** Go 1.22 + Echo + SQLC + JWT RS256 + PostgreSQL 14+
**Frontend:** Next.js 16.0.3 + React 19 + Tailwind CSS + TanStack Query
**DevOps:** Docker + NGINX + GitHub Actions + Prometheus + Grafana
**Infraestrutura:** VPS Ubuntu 22.04 + Neon (serverless PostgreSQL)

---

## 🆕 Novidades da Versão 4.0 (Nov/2025)

Esta versão consolida tudo que já foi implementado desde a v3.0:

- **Frontend v2 completo + Design System**: Página pública `/design-system-preview`, Storybook 7, Tokens, Dark/Light mode.
- **Correções de UX/SSR & Autenticação**: Refactor do `AppThemeProvider` e auth tokens (`tokens.server.ts` + `tokens.client.ts`).
- **Novos Domínios de Negócio**: Cadastro completo (Clientes, Profissionais, Serviços, Produtos), Lista da vez (Barber Turns).
- **Segurança & Governança**: RBAC (4 roles), Audit Log estruturado, Feature Flags, Rate Limiting.
- **Performance**: Redis caching, Testes de carga (k6).

---

## 📊 As 6 Fases

| Fase | Duração | Foco | Status |
|------|---------|------|--------|
| **0** | 1-3d | Repos, DB, Multi-tenant | ✅ COMPLETA |
| **1** | 3-7d | Docker, NGINX, CI/CD | ✅ COMPLETA |
| **2** | 7-14d | Backend core (auth, financial) | ✅ COMPLETA |
| **3** | 14-28d | Módulos críticos (assinaturas, crons) | ✅ COMPLETA |
| **4** | 14-28d | Frontend Next.js (paralelo) | ✅ COMPLETA |
| **5** | 14-28d | Migração progressiva MVP | 🟡 EM PROGRESSO |
| **6** | 7-14d | Hardening (segurança, observ.) | 🟡 EM PROGRESSO |
| **TOTAL** | 8-12w | MVP 2.0 completo | 🟢 ADIANTADO |

---

## 🔐 Segurança & Multi-Tenancy

### Modelo: Column-Based (tenant_id)

**Por quê?**
- Simplicidade (uma tabela por domínio)
- Escalabilidade até 100k+ tenants
- Sem overhead de schema/database management
- Fácil backup/restore

### Isolamento em 4 Camadas

1. **Auth Layer**: JWT verification
2. **Middleware**: Tenant extraction do token
3. **Repository**: WHERE tenant_id = $1 (obrigatório)
4. **Index**: Composite (tenant_id, date/status)

### Regra Ouro

```go
// ❌ NUNCA esquecer tenant_id
SELECT * FROM receitas WHERE id = $1

// ✅ SEMPRE incluir tenant_id
SELECT * FROM receitas WHERE tenant_id = $1 AND id = $2
```

---

## 💾 Database (Overview)

### Tabelas Principais

```sql
-- Core
tenants (id, nome, cnpj, ativo, plano)
users (tenant_id, email, password_hash, role)
audit_logs (tenant_id, user_id, action, resource, old/new_values)

-- Financial
categorias (tenant_id, nome, tipo)
receitas (tenant_id, descricao, valor, categoria, data, status)
despesas (tenant_id, descricao, valor, categoria, data, status)

-- Subscriptions
planos_assinatura (tenant_id, nome, valor, periodicidade)
assinaturas (tenant_id, plan_id, barbeiro_id, asaas_subscription_id, status)
assinatura_invoices (tenant_id, assinatura_id, valor, status, data_pagamento)

-- Barber Turns (Lista da Vez)
barbers_turn_list (tenant_id, professional_id, current_points, last_turn_at)
barber_turn_history (tenant_id, professional_id, month_year, total_turns)
```

---

## 🔄 Fluxo de Caixa & Automação

### Conceitos

- **Receita**: Entrada de dinheiro (via Asaas ou manual)
- **Despesa**: Saída de dinheiro (comissão, material, etc.)
- **Fluxo**: Receitas - Despesas = Saldo

### 4 Cron Jobs (Diários)

| Horário | Job | Descrição |
|---------|-----|-----------|
| 02:00 | SyncAsaasInvoices | Sincroniza faturas Asaas → Receitas |
| 03:00 | SnapshotFinanceiro | Calcula fluxo do dia, detecta anomalias |
| 04:00 | ProcessarRepasse | Cria comissão para faturas RECEBIDAS |
| 08:00 | Alertas | Verifica anomalias (zero receita, etc.) |

---

## 📋 Checklist Rápido (Task Codes)

### Backend (T-BE-XXX)
- [x] T-BE-001: Go scaffold
- [x] T-BE-002: Config
- [x] T-BE-003: DB + migrations
- [x] T-BE-004: Auth
- [x] T-BE-005-011: Financial CRUD
- [x] T-BE-012: DTOs

### Frontend (T-FE-XXX)
- [x] T-FE-001: Next.js setup
- [x] T-FE-002: API client
- [x] T-FE-003: Auth pages
- [x] T-FE-004-008: Pages (dashboard, receitas, etc.)
- [x] T-FE-009-012: Hooks, forms, components
- [x] T-FE-013-016: UI Components, Tests, Fixes

### Infrastructure (T-INFRA-XXX)
- [x] T-INFRA-001-003: Repos + decisions
- [ ] T-INFRA-004-009: Docker, NGINX, CI/CD (Em progresso)
- [ ] T-INFRA-010-015: Crons, feature flags (Em progresso)

### Quality (T-QA-XXX)
- [x] T-QA-001: Unit tests
- [ ] T-QA-004: Regression tests (Pendente)

### Security (T-SEC-XXX)
- [x] T-SEC-001-004: Rate limiting, audit, RBAC, OWASP

### DevOps (T-OPS-XXX)
- [x] T-OPS-001-003: Prometheus, Grafana, Redis
- [ ] T-OPS-010-011: LGPD, Backup (Pendente)

---

## 📈 Sucesso = Quando...

✅ **Backend**: Auth JWT ✓ | Multi-tenant ✓ | Financial CRUD ✓ | Asaas sync ✓ | Crons ✓ | >80% tests ✓

✅ **Frontend**: Login ✓ | Dashboard ✓ | CRUD receitas ✓ | Assinaturas ✓ | Mobile ✓ | E2E tests ✓

✅ **Infra**: Docker ✓ | NGINX + SSL ⏳ | CI/CD ⏳ | Backup ⏳ | Health checks ✓

✅ **Data**: 100% integridade ✓ | Totais batem ✓ | Feature flags ✓ | Rollout gradual ⏳

✅ **Security**: OWASP ✓ | LGPD ⏳ | Auditoria ✓ | Rate limiting ✓ | Sentry ⏳

---

## 📚 Documentação Completa

- **[ROADMAP_COMPLETO_V2.0.md](./ROADMAP_COMPLETO_V2.0.md)** ← Start here (100% detalhado)
- **[docs/ARQUITETURA.md](./docs/ARQUITETURA.md)** – Design patterns
- **[docs/GUIA_DEV_BACKEND.md](./docs/GUIA_DEV_BACKEND.md)** – Go setup
- **[docs/GUIA_DEV_FRONTEND.md](./docs/GUIA_DEV_FRONTEND.md)** – Next.js setup
- **[docs/GUIA_DEVOPS.md](./docs/GUIA_DEVOPS.md)** – Docker, CI/CD

---

**Documento:** PRD Barber Analytics Pro v2.0
**Status:** ✅ Pronto para Implementação
**Data:** 20/11/2025
**Timeline:** 8-12 semanas

*Documento vivo. Atualizar conforme evolução.*
