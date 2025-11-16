# 🟦 FASE 2 — Backend Go Core

**Objetivo:** Espinha dorsal do backend: auth, multi-tenant, financeiro base
**Duração:** 7-14 dias
**Dependências:** ✅ Fase 1 completa
**Sprint:** Sprint 2-3

---

## 📊 Progresso Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 2: BACKEND GO CORE                                    │
├─────────────────────────────────────────────────────────────┤
│  Progresso:  ████████████████████░  83% (10/12 concluídas) │
│  Status:     🟡 Em Progresso                                │
│  Prioridade: 🔴 ALTA                                        │
│  Estimativa: 51 horas (36 horas concluídas)                │
│  Sprint:     Sprint 2-3                                     │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Checklist de Tarefas

### ✅ T-BE-002 — Config management
- **Responsável:** Backend Lead
- **Prioridade:** 🔴 Alta
- **Estimativa:** 2h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Arquivo `internal/config/config.go` com validação
- **Horas Gastas:** 2h
- **Detalhes:** Load(), Validate(), IsProduction() com 12 config fields

### ✅ T-BE-003 — Database connection & migration
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 3h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Connection pool + migration `000_create_base_schema.sql`
- **Horas Gastas:** 3h
- **Detalhes:** NewConnection(), Health(), BeginTx(), configurable pool settings

### ✅ T-BE-004 — Domain Layer: User & Tenant
- **Responsável:** Backend Lead
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Entities, Value Objects, Repository interfaces
- **Horas Gastas:** 4h
- **Arquivos Criados:** 13 arquivos
  - Entities: tenant.go, user.go, receita.go, despesa.go, categoria.go
  - Value Objects: email.go, role.go, money.go
  - Repositories: user_repository.go, tenant_repository.go, receita_repository.go, despesa_repository.go, categoria_repository.go
  - Errors: entity/errors.go, valueobject/errors.go

### ✅ T-BE-005 — Auth Use Cases
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Login, RefreshToken, CreateUser use cases
- **Horas Gastas:** 6h
- **Arquivos Criados:** 6 arquivos
  - Services: jwt_service.go (RS256), password_hasher.go (bcrypt)
  - Use Cases: login_usecase.go, refresh_token_usecase.go, create_user_usecase.go
  - DTOs: auth_dto.go

### ✅ T-BE-006 — Auth HTTP Layer
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Endpoints: `/auth/login`, `/auth/refresh`
- **Horas Gastas:** 4h
- **Arquivo Criado:** auth_handler.go (155 linhas)
  - POST /auth/login → AccessToken + RefreshToken
  - POST /auth/refresh → Novo AccessToken
  - POST /auth/users → CreateUser

### ✅ T-BE-007 — Middlewares (Auth & Tenant)
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 3h
- **Sprint:** Sprint 2
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** JWT validation + tenant extraction middlewares
- **Horas Gastas:** 3h
- **Arquivos Criados:** 3 arquivos
  - auth_middleware.go: JWT parsing, validation, claims extraction
  - tenant_middleware.go: Tenant validation + helper functions (GetTenantIDFromContext, GetUserIDFromContext)
  - error_middleware.go: Error handling + panic recovery + CORS

### ✅ T-BE-008 — Domain Layer: Financial base
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 3
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Receita, Despesa, Categoria entities
- **Horas Gastas:** 4h
- **Arquivos Criados:** 3 entidades (receita.go, despesa.go, categoria.go) + Money VO

### ✅ T-BE-009 — Financial Repositories
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 3
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Interfaces + PostgreSQL implementations
- **Horas Gastas:** 4h
- **Arquivos Criados:** 3 repositórios PostgreSQL
  - postgres_receita_repository.go (250+ linhas): Save, FindByID, FindByTenant, FindByTenantAndPeriod, FindByTenantCategoryAndPeriod, FindByTenantStatus, Update, Delete, SumByTenantAndPeriod, Count, CountByStatus
  - postgres_despesa_repository.go (250+ linhas): Mesma interface que Receita
  - postgres_categoria_repository.go (180+ linhas): CRUD com filtro por tipo (RECEITA/DESPESA)

### ✅ T-BE-010 — Financial Use Cases (básicos)
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 6h
- **Sprint:** Sprint 3
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** CreateReceita, ListReceitas, CreateDespesa, ListDespesas
- **Horas Gastas:** 6h
- **Arquivos Criados:** 9 use cases
  - create_receita_usecase.go, list_receitas_usecase.go, update_receita_usecase.go, delete_receita_usecase.go
  - create_despesa_usecase.go, list_despesas_usecase.go, update_despesa_usecase.go, delete_despesa_usecase.go
  - calculate_cashflow_usecase.go

### ✅ T-BE-011 — Financial HTTP Layer
- **Responsável:** Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 3
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Endpoints CRUD para receitas e despesas
- **Horas Gastas:** 4h
- **Arquivos Criados:** 3 handlers
  - receita_handler.go (233 linhas): POST, GET, PUT, DELETE /receitas
  - despesa_handler.go (230 linhas): POST, GET, PUT, DELETE /despesas
  - cashflow_handler.go (82 linhas): GET /cashflow

### 🟡 T-BE-012 — DTO standardization
- **Responsável:** Backend
- **Prioridade:** 🟡 Média
- **Estimativa:** 3h
- **Sprint:** Sprint 3
- **Status:** ✅ **CONCLUÍDO**
- **Deliverable:** Mappers Domain ↔ DTO + response structure
- **Horas Gastas:** 3h
- **Arquivo Criado:** standard_response.go (45 linhas)
  - StandardResponse: Code, Message, Data, Errors, Meta, TraceID, Timestamp
  - Helper functions: Success(), Error()
  - HTTP status mapping com 11 codes (OK, CREATED, BAD_REQUEST, UNAUTHORIZED, FORBIDDEN, NOT_FOUND, CONFLICT, UNPROCESSABLE_ENTITY, RATE_LIMITED, INTERNAL_ERROR, SERVICE_UNAVAILABLE)

### 🔴 T-QA-001 — Unit tests Phase 2
- **Responsável:** QA / Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 8h
- **Sprint:** Sprint 3
- **Status:** ⏳ Não iniciado
- **Deliverable:** Coverage >80% (domain + use cases)

---

## 📈 Métricas de Sucesso

### Fase 2 completa quando:
- [x] ✅ **10 de 12 tasks concluídos (83%)**
- [x] ✅ Backend estruturado em Clean Architecture
- [x] ✅ Autenticação JWT funcional (RS256)
- [x] ✅ Multi-tenant implementado (tenant_id em todos os queries)
- [x] ✅ Módulo financeiro completo (CRUD receitas/despesas + cashflow)
- [x] ✅ HTTP endpoints funcionais (receitas, despesas, cashflow)
- [ ] ✅ Testes com >80% coverage (pendente)
- [x] ✅ Documentação atualizada

---

## 🎯 Deliverables da Fase 2

| # | Deliverable | Status |
|---|-------------|--------|
| 1 | Config management com validação | ✅ CONCLUÍDO |
| 2 | Database connection + migrations | ✅ CONCLUÍDO |
| 3 | Domain Layer (User, Tenant, Financial) | ✅ CONCLUÍDO |
| 4 | Autenticação JWT (RS256) | ✅ CONCLUÍDO |
| 5 | Multi-tenant middleware | ✅ CONCLUÍDO |
| 6 | Financial repository implementations | ✅ CONCLUÍDO |
| 7 | DTO Standardization | ✅ CONCLUÍDO |
| 8 | Financial CRUD Use Cases | ✅ CONCLUÍDO |
| 9 | Financial CRUD Endpoints | ✅ CONCLUÍDO |onclu
---

## 🚀 Próximos Passos

Após completar **100%** da Fase 2:

👉 **Iniciar FASE 3 — Módulos Backend** (`Tarefas/FASE_3_MODULOS_BACKEND.md`)

**Resumo Fase 3:**
- Fluxo de Caixa
- Assinaturas + Asaas Integration
- Cron jobs (4 tarefas diárias)
- Sincronização automática

---

## 📝 Detalhamento Técnico Selecionado

### T-BE-007 — Middleware Multi-Tenant (Exemplo)

```go
// internal/infrastructure/http/middleware/tenant.go
package middleware

import (
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Extrair token do contexto (já validado por AuthMiddleware)
        token := c.Get("user").(*jwt.Token)
        claims := token.Claims.(jwt.MapClaims)

        // Extrair tenant_id do JWT
        tenantID, ok := claims["tenant_id"].(string)
        if !ok || tenantID == "" {
            return echo.NewHTTPError(403, "tenant_id missing in token")
        }

        // Injetar no contexto
        c.Set("tenant_id", tenantID)

        return next(c)
    }
}
```

**Uso:**
```go
// Aplicar em rotas protegidas
protected := e.Group("/api")
protected.Use(AuthMiddleware)
protected.Use(TenantMiddleware)

protected.GET("/receitas", handlers.ListReceitas)
```

---

**Última Atualização:** 14/11/2025
**Status:** 🔴 Não Iniciado (0%)
**Próxima Revisão:** Após completar 50% das tarefas
