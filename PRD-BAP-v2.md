# 📊 Barber Analytics Pro v2.0 — PRD (Product Requirements Document)

**Status:** ✅ Documento Refinado e Organizado  
**Versão:** 2.0  
**Data:** 14/11/2025

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
- ✅ Frontend **Next.js 15** (React 19 + App Router)
- ✅ Database **PostgreSQL 14+** gerenciado (Neon)
- ✅ Infraestrutura **Docker + NGINX + CI/CD** profissional
- ✅ Multi-tenant **column-based** (isolamento garantido)
- ✅ Integração **Asaas** para assinaturas/repasse
- ✅ **Crons diários** para sincronização e processamento

### Stack Tecnológica

**Backend:** Go 1.22 + Echo + SQLC + JWT RS256 + PostgreSQL 14+  
**Frontend:** Next.js 15 + React 19 + Tailwind CSS + TanStack Query  
**DevOps:** Docker + NGINX + GitHub Actions + Prometheus + Grafana  
**Infraestrutura:** VPS Ubuntu 22.04 + Neon (serverless PostgreSQL)

---

## 📊 As 6 Fases

| Fase | Duração | Foco | Status |
|------|---------|------|--------|
| **0** | 1-3d | Repos, DB, Multi-tenant | �� |
| **1** | 3-7d | Docker, NGINX, CI/CD | 📅 |
| **2** | 7-14d | Backend core (auth, financial) | 📅 |
| **3** | 14-28d | Módulos críticos (assinaturas, crons) | 📅 |
| **4** | 14-28d | Frontend Next.js (paralelo) | 📅 |
| **5** | 14-28d | Migração progressiva MVP | 📅 |
| **6** | 7-14d | Hardening (segurança, observ.) | 📅 |
| **TOTAL** | 8-12w | MVP 2.0 completo | 🎯 |

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
```

---

## �� Fluxo de Caixa & Automação

### Conceitos

- **Receita**: Entrada de dinheiro (via Asaas ou manual)
- **Despesa**: Saída de dinheiro (comissão, material, etc.)
- **Fluxo**: Receitas - Despesas = Saldo

### 4 Cron Jobs (Diários)

| Horário | Job | Descrição |
|---------|-----|-----------|
| 02:00   | SyncAsaasInvoices | Sincroniza faturas Asaas → Receitas |
| 03:00   | SnapshotFinanceiro | Calcula fluxo do dia, detecta anomalias |
| 04:00   | ProcessarRepasse | Cria comissão para faturas RECEBIDAS |
| 08:00   | Alertas | Verifica anomalias (zero receita, etc.) |

### Repasse Barbeiro (Exemplo)

**Barbeiro tem 70% de comissão**

1. Fatura Asaas: R$ 100 RECEBIDA
2. Cron cria Receita: R$ 100 (entrada)
3. Cron cria Despesa: R$ 30 (comissão)
4. Barbeiro recebe: R$ 70 (líquido)

---

## 🔗 Integração Asaas

**O Quê é?** Gateway de pagamento para assinaturas recorrentes.

**Por quê?** Facilita:
- Criar assinaturas para barbeiros
- Sincronizar faturas (daily)
- Processar repassos (automático)

**APIs Utilizadas:**
```http
POST   /subscriptions          # Criar assinatura
GET    /invoices?subscription  # Listar faturas
DELETE /subscriptions/{id}     # Cancelar
```

**Error Handling:**
- `401`: API key inválida
- `422`: Validação falhou
- `429`: Rate limit (retry com backoff)
- `5xx`: Retry automático (exponential)

---

## 📋 Checklist Rápido (Task Codes)

### Backend (T-BE-XXX)
- [ ] T-BE-001: Go scaffold
- [ ] T-BE-002: Config
- [ ] T-BE-003: DB + migrations
- [ ] T-BE-004: Auth
- [ ] T-BE-005-011: Financial CRUD
- [ ] T-BE-012: DTOs

### Frontend (T-FE-XXX)
- [ ] T-FE-001: Next.js setup
- [ ] T-FE-002: API client
- [ ] T-FE-003: Auth pages
- [ ] T-FE-004-008: Pages (dashboard, receitas, etc.)
- [ ] T-FE-009-012: Hooks, forms, components

### Infrastructure (T-INFRA-XXX)
- [ ] T-INFRA-001-003: Repos + decisions
- [ ] T-INFRA-004-009: Docker, NGINX, CI/CD
- [ ] T-INFRA-010-015: Crons, feature flags

### Quality (T-QA-XXX)
- [ ] T-QA-001-004: Unit, integration, E2E, regression

### Security (T-SEC-XXX)
- [ ] T-SEC-001-004: Rate limiting, audit, RBAC, OWASP

### DevOps (T-OPS-XXX)
- [ ] T-OPS-001-005: Prometheus, Grafana, Sentry, alerts, backup

---

## 📈 Sucesso = Quando...

✅ **Backend**: Auth JWT ✓ | Multi-tenant ✓ | Financial CRUD ✓ | Asaas sync ✓ | Crons ✓ | >80% tests ✓

✅ **Frontend**: Login ✓ | Dashboard ✓ | CRUD receitas ✓ | Assinaturas ✓ | Mobile ✓ | E2E tests ✓

✅ **Infra**: Docker ✓ | NGINX + SSL ✓ | CI/CD ✓ | Backup ✓ | Health checks ✓

✅ **Data**: 100% integridade ✓ | Totais batem ✓ | Feature flags ✓ | Rollout gradual ✓

✅ **Security**: OWASP ✓ | LGPD ✓ | Auditoria ✓ | Rate limiting ✓ | Sentry ✓

---

## 🚀 Começar Agora

### Passo 1: Ler Documentação (30 min)
```bash
cat ROADMAP_COMPLETO_V2.0.md
cd docs/ && cat ARQUITETURA.md GUIA_DEV_BACKEND.md
```

### Passo 2: Fase 0 Setup (1-3 dias)
```bash
# Backend repo
git init barber-analytics-backend-v2
cd barber-analytics-backend-v2
go mod init github.com/seu-usuario/barber-analytics-backend-v2
go get github.com/labstack/echo/v4 github.com/golang-jwt/jwt/v5 github.com/lib/pq
```

### Passo 3: Fase 1 Docker (3-7 dias)
```bash
# Dockerfile
docker build -t barber-api:latest .
docker-compose up -d
curl http://localhost:8080/health
```

---

## �� Notas Importantes

⚠️ **Multi-tenant**: Sempre `tenant_id` em queries. PR review.

⚠️ **Migrations**: Versionadas no git. Testar rollback.

⚠️ **Secrets**: GitHub Secrets. NUNCA `.env` real commited.

⚠️ **Backup**: Antes de migração em produção.

⚠️ **Dependencies**: Weekly updates (security patches).

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
**Data:** 14/11/2025  
**Timeline:** 8-12 semanas

*Documento vivo. Atualizar conforme evolução.*
