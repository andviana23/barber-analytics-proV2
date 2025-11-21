> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🚩 Feature Flags API - Documentação

## 📋 Visão Geral

Sistema de feature flags por tenant que permite habilitar/desabilitar funcionalidades progressivamente.

**Base URL:** `http://localhost:8080/api/v1` (dev) | `https://api.barberanalytics.com/api/v1` (prod)

---

## 🔓 Endpoints Públicos (Tenant)

### GET /api/v1/feature-flags

Lista todos os feature flags do tenant autenticado.

**Headers:**
```http
X-Tenant-ID: e2e00000-0000-0000-0000-000000000001
```

**Response 200:**
```json
{
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "tenant_id": "e2e00000-0000-0000-0000-000000000001",
      "feature": "use_v2_financial",
      "enabled": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-20T14:25:00Z"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "tenant_id": "e2e00000-0000-0000-0000-000000000001",
      "feature": "use_v2_subscriptions",
      "enabled": false,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "timestamp": "2024-01-20T15:45:30Z"
}
```

**Curl:**
```bash
curl -X GET http://localhost:8080/api/v1/feature-flags \
  -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001"
```

---

## 🔐 Endpoints Admin (Proteção Adicional)

### GET /api/v1/admin/feature-flags

Lista feature flags com filtro opcional por tenant (admin).

**Headers:**
```http
X-Tenant-ID: e2e00000-0000-0000-0000-000000000001
Authorization: Bearer <admin_jwt_token>
```

**Query Parameters:**
- `tenant_id` (opcional): UUID do tenant para filtrar

**Response 200:**
```json
{
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "tenant_id": "e2e00000-0000-0000-0000-000000000001",
      "feature": "use_v2_financial",
      "enabled": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-20T14:25:00Z"
    }
  ],
  "timestamp": "2024-01-20T15:45:30Z"
}
```

**Curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/feature-flags?tenant_id=e2e00000-0000-0000-0000-000000000001" \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### PATCH /api/v1/admin/feature-flags

Habilita ou desabilita um feature flag para um tenant específico.

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <admin_jwt_token>
```

**Request Body:**
```json
{
  "tenant_id": "e2e00000-0000-0000-0000-000000000001",
  "feature": "use_v2_financial",
  "enabled": true
}
```

**Response 200:**
```json
{
  "code": 200,
  "message": "Feature flag updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-20T16:00:00Z"
  },
  "timestamp": "2024-01-20T16:00:00Z"
}
```

**Response 400 (Validação):**
```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": [
    "tenant_id is required",
    "feature must be one of: use_v2_financial, use_v2_subscriptions, use_v2_inventory"
  ],
  "timestamp": "2024-01-20T16:00:00Z"
}
```

**Response 404:**
```json
{
  "code": 404,
  "message": "Feature flag not found",
  "timestamp": "2024-01-20T16:00:00Z"
}
```

**Curl:**
```bash
# Habilitar feature flag
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_jwt_token>" \
  -d '{
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": true
  }'

# Desabilitar feature flag
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_jwt_token>" \
  -d '{
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": false
  }'
```

---

## 🎯 Features Disponíveis

| Feature | Descrição | Impacto |
|---------|-----------|---------|
| `use_v2_financial` | Habilita Financeiro v2 (receitas/despesas/fluxo de caixa) | Endpoints `/api/v1/receitas`, `/api/v1/despesas`, `/api/v1/fluxo-caixa`, `/api/v1/dashboard` |
| `use_v2_subscriptions` | Habilita Clube do Trato v2 (assinaturas + repasses) | Endpoints `/api/v1/assinaturas` |
| `use_v2_inventory` | Habilita Estoque v2 (produtos + movimentações) | Endpoints `/api/v1/produtos`, `/api/v1/estoque` (futuro) |

---

## 🛡️ Middleware de Proteção

Para proteger rotas que dependem de feature flags, use o middleware `FeatureFlagMiddleware`:

### Exemplo de uso (Go):

```go
import (
    httpMiddleware "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/http/middleware"
)

// Proteger endpoint de receitas v2
r.Route("/api/v1/receitas", func(r chi.Router) {
    // Middleware que verifica se use_v2_financial está habilitado
    r.Use(httpMiddleware.FeatureFlagMiddleware(featureFlagCheckUC, "use_v2_financial"))

    r.Get("/", receitaHandler.ListReceitas)
    r.Post("/", receitaHandler.CreateReceita)
})
```

### Response 403 (Feature desabilitada):

```json
{
  "code": 403,
  "message": "Feature 'use_v2_financial' is not enabled for this tenant",
  "timestamp": "2024-01-20T16:00:00Z"
}
```

---

## 🧪 Testando Feature Flags

### 1. Popular seeds iniciais

```bash
cd backend
psql $DATABASE_URL -f scripts/sql/seed_feature_flags.sql
```

### 2. Verificar flags do tenant

```bash
curl -X GET http://localhost:8080/api/v1/feature-flags \
  -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001" | jq
```

### 3. Habilitar feature para tenant beta

```bash
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": true
  }' | jq
```

### 4. Testar endpoint protegido

```bash
# Com flag habilitada -> 200 OK
curl -X GET http://localhost:8080/api/v1/receitas \
  -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001"

# Com flag desabilitada -> 403 Forbidden
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "e2e00000-0000-0000-0000-000000000001", "feature": "use_v2_financial", "enabled": false}'

curl -X GET http://localhost:8080/api/v1/receitas \
  -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001"
# Expected: {"code": 403, "message": "Feature 'use_v2_financial' is not enabled..."}
```

---

## 📊 Workflow de Rollout

### Fase 1: Beta (25% dos tenants)

```bash
# 1. Selecionar tenants beta
BETA_TENANTS=("uuid1" "uuid2" "uuid3")

# 2. Habilitar para cada tenant
for tenant in "${BETA_TENANTS[@]}"; do
    curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
      -H "Content-Type: application/json" \
      -d "{\"tenant_id\": \"$tenant\", \"feature\": \"use_v2_financial\", \"enabled\": true}"
done

# 3. Monitorar por 1 semana
# - Sentry: Erros 500
# - Logs: Divergências de cálculo
# - Feedback: Usuários reportam problemas?
```

### Fase 2: Expansão (50% → 75% → 100%)

Repetir processo acima aumentando gradualmente.

### Rollback (se necessário)

```bash
# Desabilitar para tenant específico
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": false
  }'
```

---

## 🔍 Validação de Integridade

Após habilitar flag para um tenant, validar:

### 1. Dashboard retorna dados corretos

```bash
curl -H "X-Tenant-ID: <uuid>" http://localhost:8080/api/v1/dashboard | jq
```

**Esperado:** Totais de receitas/despesas/saldo corretos.

### 2. Receitas v2 retornam dados

```bash
curl -H "X-Tenant-ID: <uuid>" http://localhost:8080/api/v1/receitas | jq
```

### 3. Comparar com MVP (se dual-read ativo)

Ver seção de Dual-Read no frontend (`T-FE-013`).

---

## 🆘 Troubleshooting

### Erro: "Feature flag not found"

**Causa:** Feature flag não foi criado para o tenant.

**Solução:**
```bash
# Executar seeds
psql $DATABASE_URL -f backend/scripts/sql/seed_feature_flags.sql
```

### Erro: "MISSING_TENANT"

**Causa:** Header `X-Tenant-ID` ausente ou inválido.

**Solução:** Adicionar header correto.

### Erro 403: "Feature not enabled"

**Causa:** Flag está desabilitada (`enabled = false`).

**Solução:** Habilitar via PATCH `/admin/feature-flags`.

---

**Última atualização:** 15/11/2025
**Versão:** 1.0.0
**Autor:** Barber Analytics Pro Team
