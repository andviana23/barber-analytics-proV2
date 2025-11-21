> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 👥 Modelo Multi-Tenant

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Aprovado

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Modelo Selecionado](#modelo-selecionado)
3. [Schema do Banco](#schema-do-banco)
4. [Isolamento de Dados](#isolamento-de-dados)
5. [Implementação](#implementação)
6. [Segurança](#segurança)
7. [Performance](#performance)
8. [Escalabilidade](#escalabilidade)

---

## 🎯 Visão Geral

Multi-tenancy é o modelo onde **múltiplas barbearias** (tenants) compartilham a mesma infraestrutura de software, com isolamento completo de dados.

### Três Modelos Possíveis

| Modelo | Isolamento | Escala | Complexidade | Custos |
|--------|-----------|--------|-------------|----|
| **Column-based** | Por linha | Média | Baixa | Baixo |
| **Schema-based** | Por schema | Alta | Média | Médio |
| **Database-based** | DB separado | Muito alta | Alta | Alto |

**Escolha:** **Column-based** ✅

**Razão:** Simplicidade, escalabilidade para 1000s de tenants, custo baixo, fácil de gerenciar.

---

## 🏗️ Modelo Selecionado: Column-Based

### Conceito

Cada linha no banco tem uma coluna `tenant_id` que identifica a qual barbearia pertence.

```
Tabela: receitas
┌─────┬───────────┬────────────┬────────┐
│ id  │ tenant_id │ descricao  │ valor  │
├─────┼───────────┼────────────┼────────┤
│ 1   │ tenant-A  │ Corte      │ 50.00  │
│ 2   │ tenant-A  │ Barba      │ 30.00  │
│ 3   │ tenant-B  │ Corte      │ 60.00  │
│ 4   │ tenant-B  │ Barba      │ 35.00  │
└─────┴───────────┴────────────┴────────┘

Tenant A vê apenas linhas com tenant_id = 'tenant-A'
Tenant B vê apenas linhas com tenant_id = 'tenant-B'
```

---

## 🗄️ Schema do Banco

### Tabela Core: Tenants

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    cnpj VARCHAR(14) UNIQUE,
    ativo BOOLEAN DEFAULT true,
    plano VARCHAR(50) DEFAULT 'free', -- free, pro, enterprise
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tenants_ativo ON tenants(ativo);
CREATE INDEX idx_tenants_cnpj ON tenants(cnpj);
```

### Tabela Core: Usuários com Tenant

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee', -- owner, manager, accountant, employee
    ativo BOOLEAN DEFAULT true,
    ultimo_login TIMESTAMP,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, email) -- Mesmo email em tenants diferentes é OK
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_tenant_email ON users(tenant_id, email);
```

### Exemplo: Tabelas com Tenant

```sql
CREATE TABLE receitas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(10, 2) NOT NULL,
    categoria VARCHAR(100),
    data DATE NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

-- Índices críticos
CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data DESC);
CREATE INDEX idx_receitas_tenant_categoria ON receitas(tenant_id, categoria);
```

### Convenção: Todo `INSERT` e `UPDATE` precisa de `tenant_id`

```sql
-- ✅ CORRETO
INSERT INTO receitas (id, tenant_id, descricao, valor, data)
VALUES (gen_random_uuid(), 'tenant-123', 'Corte', 50.00, NOW());

-- ❌ PERIGOSO (sem filtro tenant)
SELECT * FROM receitas WHERE data = '2024-11-14';

-- ✅ CORRETO (com filtro tenant)
SELECT * FROM receitas WHERE tenant_id = $1 AND data = $2;
```

---

## 🔐 Isolamento de Dados

### Princípio: Defense in Depth

**Múltiplas camadas de proteção** contra vazamento de dados.

#### Camada 1: Autenticação JWT (RS256)

```go
// internal/infrastructure/http/middleware/auth.go
package middleware

import (
    "crypto/rsa"
    "fmt"
    "net/http"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
    publicKey *rsa.PublicKey
}

func NewAuthMiddleware(publicKey *rsa.PublicKey) *AuthMiddleware {
    return &AuthMiddleware{publicKey: publicKey}
}

func (m *AuthMiddleware) ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Extrair Bearer token
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "missing Authorization header",
            })
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Verificar signature com public key
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return m.publicKey, nil
        })
        
        if err != nil || !token.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "invalid token",
            })
        }
        
        // Extrair claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "invalid claims",
            })
        }
        
        // Validar campos obrigatórios
        userID, ok := claims["sub"].(string)
        if !ok || userID == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "missing user_id in token",
            })
        }
        
        tenantID, ok := claims["tenant_id"].(string)
        if !ok || tenantID == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "missing tenant_id in token",
            })
        }
        
        // Injetar no contexto
        c.Set("user_id", userID)
        c.Set("tenant_id", tenantID)
        c.Set("role", claims["role"])
        
        return next(c)
    }
}
```

#### Camada 2: Middleware de Tenant Validation

```go
// internal/infrastructure/http/middleware/tenant.go
package middleware

import (
    "context"
    "net/http"
    
    "github.com/labstack/echo/v4"
    "github.com/andviana23/barber-analytics-backend-v2/internal/domain"
)

type TenantMiddleware struct {
    tenantRepo domain.TenantRepository
    logger     *zap.Logger
}

func NewTenantMiddleware(
    tenantRepo domain.TenantRepository,
    logger *zap.Logger,
) *TenantMiddleware {
    return &TenantMiddleware{
        tenantRepo: tenantRepo,
        logger:     logger,
    }
}

func (m *TenantMiddleware) ValidateTenant(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Extrair tenant_id do contexto (injetado pelo AuthMiddleware)
        tenantID, ok := c.Get("tenant_id").(string)
        if !ok || tenantID == "" {
            m.logger.Warn("missing tenant_id in context")
            return c.JSON(http.StatusForbidden, map[string]string{
                "error": "tenant_id required",
            })
        }
        
        // Validar UUID format
        if err := validateUUID(tenantID); err != nil {
            m.logger.Warn("invalid tenant_id format",
                zap.String("tenant_id", tenantID),
                zap.Error(err))
            return c.JSON(http.StatusBadRequest, map[string]string{
                "error": "invalid tenant_id format",
            })
        }
        
        // Verificar se tenant existe e está ativo
        ctx := c.Request().Context()
        tenant, err := m.tenantRepo.FindByID(ctx, tenantID)
        if err != nil {
            m.logger.Error("tenant not found",
                zap.String("tenant_id", tenantID),
                zap.Error(err))
            return c.JSON(http.StatusForbidden, map[string]string{
                "error": "tenant not found or inactive",
            })
        }
        
        if !tenant.Ativo {
            m.logger.Warn("inactive tenant attempted access",
                zap.String("tenant_id", tenantID),
                zap.String("tenant_nome", tenant.Nome))
            return c.JSON(http.StatusForbidden, map[string]string{
                "error": "tenant inactive",
            })
        }
        
        // Injetar tenant completo no contexto (se necessário)
        c.Set("tenant", tenant)
        
        m.logger.Debug("tenant validated",
            zap.String("tenant_id", tenantID),
            zap.String("tenant_nome", tenant.Nome))
        
        return next(c)
    }
}

func validateUUID(id string) error {
    _, err := uuid.Parse(id)
    return err
}
```

#### Camada 3: Repository Layer

```go
// ✅ SEMPRE filtra por tenant_id
func (r *PostgresReceitaRepository) FindByID(
    ctx context.Context, id, tenantID string) (*Receita, error) {
    
    query := `
        SELECT id, tenant_id, descricao, valor, data
        FROM receitas
        WHERE id = $1 AND tenant_id = $2  -- Dupla verificação
    `
    row := r.db.QueryRowContext(ctx, query, id, tenantID)
    // ...
}

// ❌ NUNCA aceita requisição sem tenantID
func (r *PostgresReceitaRepository) ListAll(ctx context.Context) ([]*Receita, error) {
    // ❌ BUG: Retorna receitas de TODOS os tenants!
    query := `SELECT * FROM receitas`
    // ...
}
```

#### Camada 4: HTTP Handler

```go
func (h *ReceitaHandler) GetByID(c echo.Context) error {
    id := c.Param("id")
    tenantID := c.Get("tenant_id").(string) // Obrigatório do middleware
    
    receita, err := h.repo.FindByID(c.Request().Context(), id, tenantID)
    if err != nil {
        return c.JSON(http.StatusNotFound, "Receita não encontrada")
    }
    
    return c.JSON(http.StatusOK, receita)
}
```

---

## 🛠️ Implementação

### Padrão: Tenant-Aware Repository

```go
// Interface
type ReceitaRepository interface {
    Save(ctx context.Context, tenantID string, receita *Receita) error
    FindByID(ctx context.Context, tenantID, id string) (*Receita, error)
    FindByTenant(ctx context.Context, tenantID string, opts FindOptions) ([]*Receita, error)
}

// Implementação
type PostgresReceitaRepository struct {
    db *sql.DB
}

func (r *PostgresReceitaRepository) FindByID(
    ctx context.Context, tenantID, id string) (*Receita, error) {
    
    // 1. Validar inputs
    if tenantID == "" || id == "" {
        return nil, errors.New("tenant_id and id required")
    }
    
    // 2. Query com tenant_id obrigatório
    query := `
        SELECT id, tenant_id, descricao, valor, data, criado_em
        FROM receitas
        WHERE id = $1 AND tenant_id = $2
        LIMIT 1
    `
    
    row := r.db.QueryRowContext(ctx, query, id, tenantID)
    
    var receita Receita
    err := row.Scan(&receita.ID, &receita.TenantID, &receita.Descricao, 
                     &receita.Valor, &receita.Data, &receita.CriadoEm)
    
    if err == sql.ErrNoRows {
        return nil, errors.New("not found")
    }
    
    return &receita, nil
}
```

### Padrão: Context Injection

```go
// Use case recebe tenantID do contexto
type GetReceitasUseCase struct {
    repo ReceitaRepository
}

func (uc *GetReceitasUseCase) Execute(
    ctx context.Context, tenantID string) ([]*Receita, error) {
    
    // Validação básica
    if tenantID == "" {
        return nil, errors.New("tenant_id required")
    }
    
    return uc.repo.FindByTenant(ctx, tenantID, FindOptions{})
}
```

---

## �️ RLS (Row Level Security) Policies

### PostgreSQL RLS para Dupla Proteção

```sql
-- Habilitar RLS em todas as tabelas com tenant_id
ALTER TABLE receitas ENABLE ROW LEVEL SECURITY;
ALTER TABLE despesas ENABLE ROW LEVEL SECURITY;
ALTER TABLE assinaturas ENABLE ROW LEVEL SECURITY;
ALTER TABLE assinatura_invoices ENABLE ROW LEVEL SECURITY;

-- Policy: View Own Unit (SELECT)
CREATE POLICY "view_own_unit"
ON receitas
FOR SELECT USING (
  tenant_id IN (
    SELECT DISTINCT tenant_id 
    FROM users 
    WHERE id = current_setting('app.user_id', true)::uuid
  )
);

-- Policy: Insert Own Unit (INSERT)
CREATE POLICY "insert_own_unit"
ON receitas
FOR INSERT WITH CHECK (
  tenant_id IN (
    SELECT DISTINCT tenant_id 
    FROM users 
    WHERE id = current_setting('app.user_id', true)::uuid
  )
);

-- Policy: Update Own Unit (UPDATE)
CREATE POLICY "update_own_unit"
ON receitas
FOR UPDATE USING (
  tenant_id IN (
    SELECT DISTINCT tenant_id 
    FROM users 
    WHERE id = current_setting('app.user_id', true)::uuid
  )
);

-- Policy: Delete Own Unit (DELETE)
CREATE POLICY "delete_own_unit"
ON receitas
FOR DELETE USING (
  tenant_id IN (
    SELECT DISTINCT tenant_id 
    FROM users 
    WHERE id = current_setting('app.user_id', true)::uuid
  )
);
```

### Ativar RLS no Connection Pool

```go
// internal/infrastructure/repository/postgres.go
package repository

import (
    "context"
    "database/sql"
    "fmt"
)

type PostgresRepository struct {
    db *sql.DB
}

// SetTenantContext configura o user_id para RLS policies
func (r *PostgresRepository) SetTenantContext(ctx context.Context, userID string) error {
    query := "SET LOCAL app.user_id = $1"
    _, err := r.db.ExecContext(ctx, query, userID)
    if err != nil {
        return fmt.Errorf("failed to set tenant context: %w", err)
    }
    return nil
}

// WithTransaction executa operação em transação com RLS
func (r *PostgresRepository) WithTransaction(
    ctx context.Context,
    userID string,
    fn func(*sql.Tx) error,
) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // Configurar RLS context
    _, err = tx.ExecContext(ctx, "SET LOCAL app.user_id = $1", userID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to set RLS context: %w", err)
    }
    
    // Executar operação
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    // Commit
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

---

## �🔒 Segurança

### Ataques Possíveis e Mitigação

| Ataque | Descrição | Mitigação |
|--------|-----------|-----------|
| **Tenant Enumeration** | Descobrir IDs de outros tenants | Rate limiting, logs de auditoria |
| **JWT Tampering** | Alterar tenant_id no token | RS256 signature verification |
| **SQL Injection** | Inserir SQL malicioso | Prepared statements (SQLC) |
| **Mass Assignment** | Alterar tenant_id na request | DTO validation, bloquear tenant_id no input |
| **Cross-Tenant Leak** | Esquecer de filtrar tenant_id | Code review, testes de segurança |

### Checklist de Segurança

- [ ] **JWT Verification**
  ```go
  // Usar RS256 (assimétrico)
  // Private key: server (sign)
  // Public key: publicado (verify)
  ```

- [ ] **Prepared Statements**
  ```go
  // ✅ CORRETO
  query := "SELECT * FROM receitas WHERE tenant_id = $1 AND id = $2"
  row := db.QueryContext(ctx, query, tenantID, id)
  
  // ❌ PERIGOSO
  query := fmt.Sprintf("SELECT * FROM receitas WHERE tenant_id = '%s'", tenantID)
  ```

- [ ] **Input Validation**
  ```go
  // Validar UUIDs
  _, err := uuid.Parse(tenantID)
  if err != nil {
    return errors.New("invalid tenant_id format")
  }
  ```

- [ ] **Auditoria**
  ```go
  // Registrar acessos a dados sensíveis
  type AuditLog struct {
      TenantID  string
      UserID    string
      Action    string // READ, WRITE, DELETE
      Resource  string
      Timestamp time.Time
  }
  ```

- [ ] **Teste de Segurança**
  ```go
  func TestCrossTenantAccess(t *testing.T) {
      // Usuário de tenant-A tenta acessar dados de tenant-B
      // Deve retornar 404 ou 403
  }
  ```

---

## ⚡ Performance

### Índices Estratégicos

```sql
-- Índices para queries comuns
CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data DESC);
CREATE INDEX idx_receitas_tenant_categoria ON receitas(tenant_id, categoria);

-- Para listagens com paginação
CREATE INDEX idx_receitas_tenant_criado ON receitas(tenant_id, criado_em DESC);

-- Para filtros de status
CREATE INDEX idx_receitas_tenant_status ON receitas(tenant_id, status);
```

### Query Optimization

```go
// ❌ Lento: Sem índices, table scan
SELECT * FROM receitas WHERE data = '2024-11-14'

// ✅ Rápido: Com índice composto
SELECT * FROM receitas WHERE tenant_id = $1 AND data = $2

// Resultado: 1ms vs 500ms em tabela grande
```

### Caching de Leitura (Futuro)

```go
// Com Redis para dashboards que acessam muitas vezes
type CachedReceitaRepository struct {
    db    *sql.DB
    cache *redis.Client
}

func (r *CachedReceitaRepository) GetSummary(ctx context.Context, tenantID string) (*Summary, error) {
    // Verificar cache
    key := fmt.Sprintf("summary:%s", tenantID)
    cached, err := r.cache.Get(ctx, key).Result()
    if err == nil {
        var summary Summary
        json.Unmarshal([]byte(cached), &summary)
        return &summary, nil
    }
    
    // Calcular se não em cache
    summary := r.calculateSummary(ctx, tenantID)
    
    // Cachear por 1 hora
    r.cache.Set(ctx, key, json.Marshal(summary), time.Hour)
    
    return summary, nil
}
```

---

## 📈 Escalabilidade

### Limites do Modelo Column-Based

```
✅ Funciona bem até: 10,000s de tenants
✅ Com particionamento: 100,000s de tenants
⚠️ Depois disso: Considerar sharding ou database-per-tenant
```

### Particionamento de Tabelas (Futuro)

```sql
-- Particionar por tenant_id para melhor performance
CREATE TABLE receitas_partitioned (
    id UUID,
    tenant_id UUID,
    data DATE,
    valor DECIMAL,
    ...
) PARTITION BY HASH (tenant_id);

CREATE TABLE receitas_part_0 PARTITION OF receitas_partitioned
    FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE receitas_part_1 PARTITION OF receitas_partitioned
    FOR VALUES WITH (MODULUS 4, REMAINDER 1);
-- ... etc
```

### Sharding (Futuro Distante)

```
Se um dia chegar a 100,000s de tenants:

Backend Go → Router (qual shard?)
              ↓
         Shard 1: Tenants A-J
         Shard 2: Tenants K-S
         Shard 3: Tenants T-Z
         
Cada shard tem seu próprio PostgreSQL
```

---

## 🧪 Testes

### Teste de Isolamento

```go
// tests/integration/tenant_isolation_test.go
package integration

import (
    "context"
    "testing"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTenantIsolation(t *testing.T) {
    ctx := context.Background()
    
    // Setup: Criar dois tenants distintos
    tenant1ID := uuid.NewString()
    tenant2ID := uuid.NewString()
    
    tenant1, err := tenantRepo.Create(ctx, &domain.Tenant{
        ID:   tenant1ID,
        Nome: "Barbearia Alpha",
        CNPJ: "12345678000190",
        Ativo: true,
    })
    require.NoError(t, err)
    
    tenant2, err := tenantRepo.Create(ctx, &domain.Tenant{
        ID:   tenant2ID,
        Nome: "Barbearia Beta",
        CNPJ: "98765432000100",
        Ativo: true,
    })
    require.NoError(t, err)
    
    // Inserir dados em tenant1
    receita1, err := receitaRepo.Save(ctx, tenant1ID, &domain.Receita{
        ID:        uuid.NewString(),
        Descricao: "Corte Alpha",
        Valor:     50.00,
        Data:      time.Now(),
    })
    require.NoError(t, err)
    
    // Inserir dados em tenant2
    receita2, err := receitaRepo.Save(ctx, tenant2ID, &domain.Receita{
        ID:        uuid.NewString(),
        Descricao: "Corte Beta",
        Valor:     60.00,
        Data:      time.Now(),
    })
    require.NoError(t, err)
    
    // Test 1: Tenant1 busca APENAS seus dados
    receitas1, err := receitaRepo.FindByTenant(ctx, tenant1ID, domain.FindOptions{})
    require.NoError(t, err)
    assert.Len(t, receitas1, 1)
    assert.Equal(t, "Corte Alpha", receitas1[0].Descricao)
    assert.Equal(t, 50.00, receitas1[0].Valor)
    assert.Equal(t, tenant1ID, receitas1[0].TenantID)
    
    // Test 2: Tenant2 busca APENAS seus dados
    receitas2, err := receitaRepo.FindByTenant(ctx, tenant2ID, domain.FindOptions{})
    require.NoError(t, err)
    assert.Len(t, receitas2, 1)
    assert.Equal(t, "Corte Beta", receitas2[0].Descricao)
    assert.Equal(t, 60.00, receitas2[0].Valor)
    assert.Equal(t, tenant2ID, receitas2[0].TenantID)
    
    // Test 3: Tenant1 tenta acessar dado de Tenant2 (cross-tenant)
    _, err = receitaRepo.FindByID(ctx, tenant1ID, receita2.ID)
    assert.Error(t, err)
    assert.Equal(t, "not found", err.Error())
    
    // Test 4: Tenant2 tenta acessar dado de Tenant1 (cross-tenant)
    _, err = receitaRepo.FindByID(ctx, tenant2ID, receita1.ID)
    assert.Error(t, err)
    assert.Equal(t, "not found", err.Error())
}

func TestTenantMiddlewareIsolation(t *testing.T) {
    // Setup HTTP test server
    e := echo.New()
    
    // Middleware stack
    e.Use(middleware.AuthMiddleware(...))
    e.Use(middleware.TenantMiddleware(...))
    
    // Test endpoint
    e.GET("/api/v2/receitas/:id", handlers.GetReceitaByID)
    
    // Tenant1 JWT
    tenant1Token := generateJWT("user-123", "tenant-abc", "owner")
    
    // Tenant2 JWT
    tenant2Token := generateJWT("user-456", "tenant-xyz", "owner")
    
    // Create receita for tenant-abc
    receitaID := createReceita(t, "tenant-abc", "Corte", 50.00)
    
    // Test 1: Tenant1 acessa sua própria receita
    req1 := httptest.NewRequest(http.MethodGet, "/api/v2/receitas/"+receitaID, nil)
    req1.Header.Set("Authorization", "Bearer "+tenant1Token)
    rec1 := httptest.NewRecorder()
    e.ServeHTTP(rec1, req1)
    
    assert.Equal(t, http.StatusOK, rec1.Code)
    
    // Test 2: Tenant2 tenta acessar receita de Tenant1
    req2 := httptest.NewRequest(http.MethodGet, "/api/v2/receitas/"+receitaID, nil)
    req2.Header.Set("Authorization", "Bearer "+tenant2Token)
    rec2 := httptest.NewRecorder()
    e.ServeHTTP(rec2, req2)
    
    assert.Equal(t, http.StatusNotFound, rec2.Code)
}

func TestRLSPolicyEnforcement(t *testing.T) {
    // Test com RLS policies habilitadas
    ctx := context.Background()
    
    // Setup connection com user_id configurado
    db := setupDBWithRLS(t)
    defer db.Close()
    
    tenant1ID := uuid.NewString()
    user1ID := uuid.NewString()
    
    // Criar tenant e usuário
    createTenant(t, db, tenant1ID)
    createUser(t, db, user1ID, tenant1ID)
    
    // Configurar RLS context
    _, err := db.ExecContext(ctx, "SET LOCAL app.user_id = $1", user1ID)
    require.NoError(t, err)
    
    // Inserir receita
    receitaID := uuid.NewString()
    _, err = db.ExecContext(ctx,
        "INSERT INTO receitas (id, tenant_id, descricao, valor, data) VALUES ($1, $2, $3, $4, $5)",
        receitaID, tenant1ID, "Corte", 50.00, time.Now())
    require.NoError(t, err)
    
    // Tentar acessar com user de outro tenant (deve falhar)
    user2ID := uuid.NewString()
    tenant2ID := uuid.NewString()
    createTenant(t, db, tenant2ID)
    createUser(t, db, user2ID, tenant2ID)
    
    _, err = db.ExecContext(ctx, "SET LOCAL app.user_id = $1", user2ID)
    require.NoError(t, err)
    
    // Query deve retornar vazio (RLS bloqueia)
    var count int
    err = db.QueryRowContext(ctx,
        "SELECT COUNT(*) FROM receitas WHERE id = $1",
        receitaID).Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 0, count) // RLS policy bloqueou acesso
}
```

---

## 📝 Checklist de Implementação

### Backend (Go)

- [ ] **Schema**: Coluna `tenant_id UUID NOT NULL` em TODAS as tabelas
- [ ] **Índices**: Compostos `(tenant_id, field)` para queries comuns
- [ ] **Migrations**: Todas as tabelas com FK para `tenants(id) ON DELETE CASCADE`
- [ ] **JWT**: Claims incluem `tenant_id` (RS256 signature)
- [ ] **Middleware Auth**: `AuthMiddleware.ValidateJWT()` extrai `tenant_id` e injeta no context
- [ ] **Middleware Tenant**: `TenantMiddleware.ValidateTenant()` verifica se tenant está ativo
- [ ] **Repository Pattern**: TODOS os métodos recebem `tenantID string` como primeiro parâmetro
- [ ] **Query Safety**: TODAS as queries incluem `WHERE tenant_id = $1`
- [ ] **RLS Policies**: PostgreSQL RLS habilitado em todas as tabelas
- [ ] **Context Injection**: `SetTenantContext()` configurado para RLS policies
- [ ] **Error Handling**: 404 quando cross-tenant access (não 403 para evitar enumeration)
- [ ] **Audit Logs**: Registrar `tenant_id` + `user_id` em todas as operações críticas

### Testing

- [ ] **Unit Tests**: Repository methods validam `tenant_id` != ""
- [ ] **Integration Tests**: `TestTenantIsolation()` verifica isolamento completo
- [ ] **Security Tests**: `TestCrossTenantAccess()` valida bloqueio cross-tenant
- [ ] **RLS Tests**: `TestRLSPolicyEnforcement()` valida PostgreSQL RLS
- [ ] **HTTP Tests**: `TestTenantMiddlewareIsolation()` valida middleware stack

### Security

- [ ] **SQL Injection**: Prepared statements ($1, $2) em TODAS as queries
- [ ] **JWT Verification**: RS256 signature com public key verification
- [ ] **UUID Validation**: `uuid.Parse()` antes de usar tenant_id
- [ ] **Input Validation**: DTO validation com `go-playground/validator`
- [ ] **Rate Limiting**: 30 req/min por tenant em endpoints sensíveis
- [ ] **Audit Logs**: `audit_logs` tabela registrando WHO, WHAT, WHEN, WHERE

### Documentation

- [x] **MODELO_MULTI_TENANT.md**: Documentação completa (este arquivo)
- [ ] **API_REFERENCE.md**: Atualizar com header `X-Tenant-ID` (se necessário)
- [ ] **GUIA_DEV_BACKEND.md**: Adicionar seção sobre multi-tenancy patterns
- [ ] **Runbook**: Procedimento de resposta a incidente de vazamento de dados

### Monitoring

- [ ] **Metrics**: Prometheus counter `tenant_requests_total{tenant_id, endpoint}`
- [ ] **Alerts**: Alert se tenant acessa > 1000 req/min (possível ataque)
- [ ] **Logs**: Structured logging com `tenant_id` em TODOS os logs
- [ ] **Traces**: OpenTelemetry span tags com `tenant.id` attribute

---

## 🚨 Troubleshooting

### Problema: Receita não encontrada mas existe no banco

**Causa:** Esqueceu de filtrar por `tenant_id` na query

**Solução:**
```sql
-- ❌ Errado
SELECT * FROM receitas WHERE id = $1

-- ✅ Correto
SELECT * FROM receitas WHERE id = $1 AND tenant_id = $2
```

### Problema: JWT não contém tenant_id

**Causa:** Token gerado antes de adicionar claim `tenant_id`

**Solução:**
```go
// Regerar token com claim obrigatório
claims := jwt.MapClaims{
    "sub": userID,
    "tenant_id": tenantID, // ← Obrigatório
    "role": role,
    "exp": time.Now().Add(15 * time.Minute).Unix(),
}
```

### Problema: RLS policies não estão funcionando

**Causa:** `app.user_id` não está sendo configurado na conexão

**Solução:**
```go
// Antes de cada query em transação
tx.ExecContext(ctx, "SET LOCAL app.user_id = $1", userID)
```

### Problema: Cross-tenant access bypass

**Causa:** Handler não está usando middleware `TenantMiddleware`

**Solução:**
```go
// Adicionar middleware na rota
api := e.Group("/api/v2")
api.Use(middleware.AuthMiddleware(...))
api.Use(middleware.TenantMiddleware(...)) // ← Obrigatório
```

---

## 📚 Referências

- [Multi-tenancy Best Practices - PostgreSQL Wiki](https://wiki.postgresql.org/wiki/Multi-tenancy)
- [Row Level Security - PostgreSQL Docs](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [JWT Best Practices - RFC 8725](https://datatracker.ietf.org/doc/html/rfc8725)
- [OWASP Top 10 - Broken Access Control](https://owasp.org/Top10/A01_2021-Broken_Access_Control/)
- [Neon Multi-Tenant Architecture](https://neon.tech/docs/guides/multi-tenant-architecture)

---

**Status:** ✅ Documentação Completa  
**Última atualização:** 14/11/2025  
**Responsável:** @arquiteto  
**Completa T-INFRA-003** (Define Multi-Tenant model)
