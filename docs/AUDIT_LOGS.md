# 📝 Audit Logs - Sistema de Auditoria

**Versão:** 1.0
**Última Atualização:** 15/11/2025
**Status:** ✅ Implementado

---

## 📋 Visão Geral

O sistema de audit logs do Barber Analytics Pro registra **todas as operações de modificação de dados** (CREATE, UPDATE, DELETE) para:

- ✅ **Compliance:** LGPD, SOC2, ISO 27001
- ✅ **Debugging:** Investigar erros e corrupção de dados
- ✅ **Security:** Detectar atividades suspeitas
- ✅ **Accountability:** Rastrear quem fez o quê e quando

---

## 🗄️ Schema da Tabela

### `audit_logs`

```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    action          VARCHAR(50) NOT NULL,              -- CREATE, UPDATE, DELETE, READ
    resource_type   VARCHAR(100),                      -- receita, despesa, assinatura, etc
    resource_name   VARCHAR(100) NOT NULL,             -- /api/v1/receitas, etc
    resource_id     VARCHAR(255),                      -- UUID do recurso afetado
    old_values      JSONB,                             -- Valores anteriores (UPDATE/DELETE)
    new_values      JSONB,                             -- Valores novos (CREATE/UPDATE)
    ip_address      INET,                              -- IP do cliente
    user_agent      TEXT,                              -- User-Agent do cliente
    timestamp       TIMESTAMPTZ DEFAULT NOW(),         -- Data/hora da operação
    deleted_at      TIMESTAMPTZ                        -- Soft delete (retenção 90 dias)
);
```

### Índices

```sql
-- Query por tenant + ordenação por data
CREATE INDEX idx_audit_logs_tenant_timestamp
    ON audit_logs(tenant_id, timestamp DESC);

-- Query por tipo de recurso
CREATE INDEX idx_audit_logs_resource_type
    ON audit_logs(resource_type, tenant_id, timestamp DESC);

-- Query por recurso específico
CREATE INDEX idx_audit_logs_resource
    ON audit_logs(resource_name, resource_id);

-- Query por usuário
CREATE INDEX idx_audit_logs_user
    ON audit_logs(user_id);

-- Query por ação
CREATE INDEX idx_audit_logs_action
    ON audit_logs(action, tenant_id);

-- Soft delete cleanup
CREATE INDEX idx_audit_logs_deleted_at
    ON audit_logs(deleted_at)
    WHERE deleted_at IS NOT NULL;
```

---

## 🎯 Tipos de Ação

| Ação | Descrição | old_values | new_values |
|------|-----------|------------|------------|
| `CREATE` | Criação de recurso | `null` | ✅ Objeto completo |
| `UPDATE` | Atualização de recurso | ✅ Valores antigos | ✅ Valores novos |
| `DELETE` | Exclusão de recurso | ✅ Objeto completo | `null` |
| `READ` | Leitura (opcional) | `null` | `null` |

---

## 🔧 Tipos de Recurso

| Resource Type | Descrição | Exemplo resource_name |
|---------------|-----------|----------------------|
| `receita` | Receitas financeiras | `/api/v1/receitas` |
| `despesa` | Despesas financeiras | `/api/v1/despesas` |
| `assinatura` | Assinaturas (Clube do Trato) | `/api/v1/assinaturas` |
| `produto` | Produtos de estoque | `/api/v1/produtos` |
| `user` | Usuários | `/api/v1/users` |
| `feature_flag` | Feature flags | `/api/v1/admin/feature-flags` |
| `categoria` | Categorias | `/api/v1/categorias` |
| `plano` | Planos de assinatura | `/api/v1/planos` |

---

## 📡 API - Endpoints Admin

### 1. Listar Audit Logs (com filtros)

```http
GET /api/v1/admin/audit-logs?user_id=...&action=...&resource_type=...&date_from=...&date_to=...&limit=50&offset=0
```

**Query Parameters:**
- `user_id` (opcional): Filtrar por usuário específico
- `action` (opcional): Filtrar por ação (`CREATE`, `UPDATE`, `DELETE`)
- `resource_type` (opcional): Filtrar por tipo de recurso (`receita`, `despesa`, etc)
- `resource_id` (opcional): Filtrar por ID específico do recurso
- `date_from` (opcional): Data início (RFC3339, ex: `2025-01-01T00:00:00Z`)
- `date_to` (opcional): Data fim (RFC3339)
- `limit` (opcional): Máximo de resultados (padrão: 50, máximo: 200)
- `offset` (opcional): Paginação (padrão: 0)

**Response:**
```json
{
  "code": "OK",
  "message": "Audit logs recuperados com sucesso",
  "data": {
    "data": [
      {
        "id": "uuid",
        "tenant_id": "uuid",
        "user_id": "uuid",
        "action": "UPDATE",
        "resource_type": "receita",
        "resource_name": "/api/v1/receitas",
        "resource_id": "uuid",
        "old_values": {"valor": 100.00, "descricao": "Venda antiga"},
        "new_values": {"valor": 150.00, "descricao": "Venda atualizada"},
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "timestamp": "2025-11-15T10:30:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "total": 150,
      "limit": 50,
      "offset": 0
    }
  },
  "timestamp": "2025-11-15T10:35:00Z"
}
```

---

### 2. Listar Audit Logs por Usuário

```http
GET /api/v1/admin/audit-logs/user/{user_id}?limit=50&offset=0
```

**Response:** Lista de audit logs do usuário específico

---

### 3. Listar Audit Logs por Recurso

```http
GET /api/v1/admin/audit-logs/resource/{resource_type}/{resource_id}?limit=50&offset=0
```

**Exemplo:**
```http
GET /api/v1/admin/audit-logs/resource/receita/550e8400-e29b-41d4-a716-446655440000
```

**Response:** Histórico completo de mudanças do recurso

---

## 🛠️ Uso no Backend

### Registrando Operações

#### CREATE
```go
err := auditService.RecordCreate(
    ctx,
    tenantID,
    &userID,
    entity.ResourceTypeReceita,
    "/api/v1/receitas",
    receitaID,
    receitaCriada, // Objeto completo
    &ipAddress,
    &userAgent,
)
```

#### UPDATE
```go
err := auditService.RecordUpdate(
    ctx,
    tenantID,
    &userID,
    entity.ResourceTypeReceita,
    "/api/v1/receitas",
    receitaID,
    receitaAntiga,  // Valores antes da mudança
    receitaNova,    // Valores após a mudança
    &ipAddress,
    &userAgent,
)
```

#### DELETE
```go
err := auditService.RecordDelete(
    ctx,
    tenantID,
    &userID,
    entity.ResourceTypeReceita,
    "/api/v1/receitas",
    receitaID,
    receitaDeletada, // Objeto antes de deletar
    &ipAddress,
    &userAgent,
)
```

---

## 🔄 Retenção de Dados (90 dias)

### Política de Retenção

1. **Soft Delete após 90 dias:** Logs são marcados como `deleted_at = NOW()`
2. **Hard Delete após 180 dias:** Logs soft deleted são removidos permanentemente

### Executar Manualmente

```sql
-- Soft delete logs > 90 dias
UPDATE audit_logs
SET deleted_at = NOW()
WHERE deleted_at IS NULL
AND timestamp < NOW() - INTERVAL '90 days';

-- Hard delete logs soft deleted > 90 dias
DELETE FROM audit_logs
WHERE deleted_at IS NOT NULL
AND deleted_at < NOW() - INTERVAL '90 days';
```

### Automação (Cron Job)

**Adicionar ao scheduler:**

```go
// Job: Cleanup audit logs antigos
cronScheduler.AddJob("@daily", func() {
    olderThan := time.Now().AddDate(0, 0, -90)
    count, err := auditLogRepo.SoftDeleteOld(ctx, olderThan)
    if err != nil {
        logger.Error("Failed to soft delete old audit logs", zap.Error(err))
        return
    }
    logger.Info("Soft deleted old audit logs", zap.Int("count", count))
})
```

---

## 📊 Queries Úteis

### 1. Listar últimas 100 ações de um usuário

```sql
SELECT
    action,
    resource_type,
    resource_id,
    timestamp
FROM audit_logs
WHERE tenant_id = 'uuid'
AND user_id = 'uuid'
AND deleted_at IS NULL
ORDER BY timestamp DESC
LIMIT 100;
```

---

### 2. Histórico completo de um recurso (ex: receita)

```sql
SELECT
    user_id,
    action,
    old_values,
    new_values,
    timestamp
FROM audit_logs
WHERE tenant_id = 'uuid'
AND resource_type = 'receita'
AND resource_id = 'uuid'
AND deleted_at IS NULL
ORDER BY timestamp ASC;
```

---

### 3. Detectar exclusões em massa (suspeito)

```sql
SELECT
    user_id,
    COUNT(*) as delete_count,
    MIN(timestamp) as first_delete,
    MAX(timestamp) as last_delete
FROM audit_logs
WHERE tenant_id = 'uuid'
AND action = 'DELETE'
AND timestamp > NOW() - INTERVAL '1 hour'
AND deleted_at IS NULL
GROUP BY user_id
HAVING COUNT(*) > 10
ORDER BY delete_count DESC;
```

---

### 4. Atividade por tipo de recurso (dashboard)

```sql
SELECT
    resource_type,
    action,
    COUNT(*) as count
FROM audit_logs
WHERE tenant_id = 'uuid'
AND timestamp > NOW() - INTERVAL '7 days'
AND deleted_at IS NULL
GROUP BY resource_type, action
ORDER BY count DESC;
```

---

## 🔒 Segurança & Compliance

### LGPD Compliance

✅ **Dados pessoais anonimizados:** Quando usuário é deletado, `user_id` vira `NULL` (ON DELETE SET NULL)
✅ **Retenção limitada:** 90 dias (conforme Art. 15 LGPD)
✅ **Auditoria de acesso:** Registrar `READ` para dados sensíveis (opcional)

### SOC2 Compliance

✅ **Rastreabilidade:** Todos os eventos registrados com timestamp, IP, user agent
✅ **Imutabilidade:** Logs nunca são editados, apenas soft deleted
✅ **Alertas:** Detectar padrões suspeitos (deleções em massa, etc)

---

## ⚠️ Considerações de Performance

### Volume Estimado

| Tenant | Operações/dia | Registros/mês | Storage (90 dias) |
|--------|---------------|---------------|-------------------|
| Pequeno | 500 | 15.000 | ~5 MB |
| Médio | 2.000 | 60.000 | ~20 MB |
| Grande | 10.000 | 300.000 | ~100 MB |

### Otimizações

✅ **Índices estratégicos:** Queries rápidas (<50ms)
✅ **Particionamento (futuro):** Por mês ou tenant
✅ **Archiving (futuro):** Mover logs antigos para S3/Glacier

---

## 🚀 Roadmap Futuro

- [ ] Dashboard visual de auditoria (Grafana)
- [ ] Alertas automáticos (Slack/email) para padrões suspeitos
- [ ] Export para CSV/JSON (compliance)
- [ ] Diff visual (old vs new values)
- [ ] Revert automático (rollback de mudanças)

---

**Última Atualização:** 15/11/2025
**Autor:** Andrey Viana
**Status:** ✅ Produção
