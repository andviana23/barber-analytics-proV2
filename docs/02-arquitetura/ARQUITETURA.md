> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🏗️ Arquitetura Barber Analytics Pro v2.0

**Versão:** 2.0  
**Data Criação:** 14/11/2025  
**Status:** Definição e Planejamento  
**Autor:** Arquiteto de Software Sr.

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Princípios Arquiteturais](#princípios-arquiteturais)
3. [Stack Tecnológico](#stack-tecnológico)
4. [Arquitetura em Camadas](#arquitetura-em-camadas)
5. [Estrutura de Diretórios](#estrutura-de-diretórios)
6. [Padrões de Design](#padrões-de-design)
7. [Fluxo de Dados](#fluxo-de-dados)
8. [Multi-Tenancy](#multi-tenancy)
9. [Segurança](#segurança)
10. [Escalabilidade](#escalabilidade)

---

## 🎯 Visão Geral

O Barber Analytics Pro v2.0 é uma plataforma SaaS modular e escalável para gerenciamento completo de barbearias, construída com **Clean Architecture**, **Domain-Driven Design (DDD)** e aderência aos princípios **SOLID**.

### Objetivos Arquiteturais

- ✅ **Independência de Framework**: Lógica de negócio desacoplada de ferramentas
- ✅ **Testabilidade**: Código altamente testável em todos os níveis
- ✅ **Manutenibilidade**: Estrutura clara e padrões consistentes
- ✅ **Escalabilidade**: Suporte a múltiplos tenants e crescimento horizontal
- ✅ **Performance**: Otimizações em queries, cache e processamento assíncrono
- ✅ **Segurança**: Isolamento de dados, auditoria e compliance

---

## 🏛️ Princípios Arquiteturais

### 1. Clean Architecture

```
┌─────────────────────────────────────────┐
│       Presentation Layer (HTTP/UI)      │
├─────────────────────────────────────────┤
│      Application Layer (Use Cases)      │
├─────────────────────────────────────────┤
│    Domain Layer (Business Rules)        │
├─────────────────────────────────────────┤
│  Infrastructure Layer (DB, APIs, etc)   │
└─────────────────────────────────────────┘
```

**Direção de dependências:** Centro (Domain) → Externo (Infrastructure)

### 2. Domain-Driven Design (DDD)

- **Ubiquitous Language**: Linguagem de negócio consistente
- **Bounded Contexts**: Módulos independentes (Financeiro, Assinaturas, Estoque, Lista da Vez)
- **Aggregates**: Entidades relacionadas com raízes claras
- **Value Objects**: Objetos imutáveis sem identidade
- **Repositories**: Abstração de persistência por Aggregate

### 3. SOLID Principles

| Princípio | Aplicação |
|-----------|-----------|
| **S** - SRP | Cada classe tem uma única responsabilidade |
| **O** - OCP | Aberto para extensão, fechado para modificação |
| **L** - LSP | Subtypes são substituíveis por seus tipos base |
| **I** - ISP | Interfaces específicas ao cliente |
| **D** - DIP | Dependências em abstrações, não em implementações |

---

## 🛠️ Stack Tecnológico

### Backend

```yaml
Linguagem: Go 1.24.0 (toolchain go1.24.10)
Framework HTTP: Echo v4 (leve, rápido, middleware-friendly)
ORM/Query Builder: SQLC (type-safe SQL)
Autenticação: JWT (RS256) + Refresh Tokens
Validação: go-playground/validator/v10
Scheduler: Cron em Go (robfig/cron/v3) + systemd para produção
Logger: Zap (structured logging em JSON)
Trace: OpenTelemetry (opcional, futuro)
```

### Banco de Dados

```yaml
Principal: PostgreSQL 14+
Provedor Recomendado: Neon (serverless, backup automático)
Alternativa: Supabase (DB-only mode)
Migrations: golang-migrate/migrate
Backup: Automático (Neon/Supabase) + snapshots periódicos
```

### Frontend (MVP -> V2)

```yaml
MVP 1.0: React 19 + Vite
V2.0 SaaS: Next.js 16.0.3 (App Router) + React 19
Styling: Tailwind CSS 4
State Management: TanStack Query (React Query)
Form Validation: Zod + React Hook Form
UI Components: shadcn/ui
```

### DevOps & Infraestrutura

```yaml
Reverse Proxy: NGINX (SSL/TLS via Certbot)
CI/CD: GitHub Actions
Logs & Monitoring: Grafana + Prometheus
APM: Sentry (para exceções e performance)
Hosting: VPS Ubuntu 22.04 LTS
```

---

## 🏗️ Arquitetura em Camadas

### Backend Go (Clean Architecture)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entrypoint
├── internal/
│   ├── config/                        # Leitura de env
│   ├── domain/                        # Business logic (entities, value objects)
│   ├── application/
│   │   ├── dto/                       # Data Transfer Objects
│   │   ├── mapper/                    # Domain <-> DTO mapping
│   │   └── usecase/                   # Application use cases
│   ├── infrastructure/
│   │   ├── http/                      # HTTP handlers e middlewares
│   │   ├── repository/                # Database repositories
│   │   ├── external/                  # Integrações externas (Asaas, etc)
│   │   └── scheduler/                 # Cron jobs
│   └── ports/                         # Interfaces (abstrações)
├── migrations/                        # SQL migrations
├── tests/                            # Testes integrados
└── go.mod
```

### Camada de Domínio (Domain Layer)

```go
// Entidade - Aggregate Root
type Barbearia struct {
    ID            string
    Nome          string
    CNPJ          string
    Endereco      Endereco           // Value Object
    Barbeiros     []Barbeiro         // Child entities
    Configuracoes Configuracoes      // Value Object
    CriadoEm      time.Time
    AtualizadoEm  time.Time
}

// Entidade - Lista da Vez (Novo Módulo)
type BarbersTurnList struct {
    ID             string
    TenantID       string
    ProfessionalID string
    CurrentPoints  int
    LastTurnAt     time.Time
    IsActive       bool
}

// Value Object - Imutável
type Endereco struct {
    Rua       string
    Numero    int
    Complemento string
    Cidade    string
    UF        string
    CEP       string
}

// Repository Interface (Port)
type BarbeariaRepository interface {
    Save(ctx context.Context, barbearia *Barbearia) error
    FindByID(ctx context.Context, id string) (*Barbearia, error)
    FindByTenantID(ctx context.Context, tenantID string) (*Barbearia, error)
}
```

### Camada de Aplicação (Application Layer)

```go
// Use Case
type CreateReceitaUseCase struct {
    repository domain.ReceitaRepository
    service    domain.CalculoComissaoService
}

func (uc *CreateReceitaUseCase) Execute(ctx context.Context, 
    input CreateReceitaInput) (*CreateReceitaOutput, error) {
    // Validações
    // Lógica de negócio
    // Persistência
    // Retorno
}

// DTO - entrada
type CreateReceitaInput struct {
    TenantID      string    `json:"tenant_id" validate:"required"`
    Descricao     string    `json:"descricao" validate:"required,max=255"`
    Valor         float64   `json:"valor" validate:"required,gt=0"`
    Data          time.Time `json:"data" validate:"required"`
    Categoria     string    `json:"categoria" validate:"required"`
}

// DTO - saída
type CreateReceitaOutput struct {
    ID        string    `json:"id"`
    TenantID  string    `json:"tenant_id"`
    Descricao string    `json:"descricao"`
    Valor     float64   `json:"valor"`
    Status    string    `json:"status"`
}
```

### Camada de Apresentação (HTTP/Delivery Layer)

```go
// Handler
type ReceitaHandler struct {
    createUseCase *application.CreateReceitaUseCase
    listUseCase   *application.ListReceitasUseCase
}

func (h *ReceitaHandler) Create(c echo.Context) error {
    var input application.CreateReceitaInput
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{...})
    }
    
    output, err := h.createUseCase.Execute(c.Request().Context(), input)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{...})
    }
    
    return c.JSON(http.StatusCreated, output)
}
```

### Camada de Infraestrutura (Infrastructure Layer)

```go
// Repository Implementation
type PostgresReceitaRepository struct {
    db *sql.DB
}

func (r *PostgresReceitaRepository) Save(ctx context.Context, 
    receita *domain.Receita) error {
    query := `
        INSERT INTO receitas (id, tenant_id, descricao, valor, data)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := r.db.ExecContext(ctx, query, 
        receita.ID, receita.TenantID, receita.Descricao, 
        receita.Valor, receita.Data)
    return err
}
```

---

## 📂 Estrutura de Diretórios

```
barber-analytics-pro/
│
├── backend/                        # Backend em Go
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   ├── valueobject/
│   │   │   └── service/
│   │   ├── application/
│   │   │   ├── dto/
│   │   │   ├── mapper/
│   │   │   └── usecase/
│   │   ├── infrastructure/
│   │   │   ├── http/
│   │   │   │   ├── handler/
│   │   │   │   ├── middleware/
│   │   │   │   └── route.go
│   │   │   ├── repository/
│   │   │   ├── external/
│   │   │   │   └── asaas/
│   │   │   └── scheduler/
│   │   └── ports/
│   ├── migrations/
│   ├── tests/
│   └── go.mod
│
├── frontend/                       # Frontend Next.js (v2)
│   ├── app/
│   │   ├── (auth)/
│   │   ├── (dashboard)/
│   │   ├── api/
│   │   ├── components/            # Componentes (agora dentro de app/)
│   │   ├── lib/                   # Utils e Hooks (agora dentro de app/)
│   │   ├── public/
│   │   └── package.json
│   │
│   ├── frontend-v1/                # Frontend React/Vite (MVP)
│   ├── src/
│   ├── public/
│   └── package.json
│
├── docs/                              # Documentação
│   ├── ARQUITETURA.md
│   ├── ROADMAP_IMPLEMENTACAO_V2.md
│   ├── MODELO_MULTI_TENANT.md
│   ├── FINANCEIRO.md
│   ├── ASSINATURAS.md
│   ├── ESTOQUE.md
│   ├── BANCO_DE_DADOS.md
│   ├── API_REFERENCE.md
│   ├── DOMAIN_MODELS.md
│   ├── FLUXO_CRONS.md
│   ├── INTEGRACOES_ASAAS.md
│   ├── GUIA_DEV_BACKEND.md
│   ├── GUIA_DEV_FRONTEND.md
│   └── GUIA_DEVOPS.md
│
├── infra/                            # Infraestrutura e DevOps
│   ├── nginx/
│   │   └── nginx.conf
│   └── .github/workflows/
│
├── PRD-BAP-v2.md
└── README.md
```

---

## 🎨 Padrões de Design

### 1. Repository Pattern

Abstração para persistência de dados:

```go
// Port (Interface)
type ReceitaRepository interface {
    Save(ctx context.Context, receita *Receita) error
    FindByID(ctx context.Context, id string) (*Receita, error)
    FindByTenantAndPeriod(ctx context.Context, 
        tenantID string, from, to time.Time) ([]*Receita, error)
}

// Adapter (Implementação)
type PostgresReceitaRepository struct { ... }
```

### 2. Dependency Injection

Injeção de dependências no startup:

```go
func InitializeReceitaHandler(db *sql.DB) *ReceitaHandler {
    repo := repository.NewPostgresReceitaRepository(db)
    createUC := application.NewCreateReceitaUseCase(repo)
    return http.NewReceitaHandler(createUC)
}
```

### 3. DTO (Data Transfer Object)

Separação entre modelo de domínio e dados transmitidos:

```go
// Domain
type Receita struct {
    ID      string
    Valor   float64
    // ...
}

// DTO
type ReceitaResponse struct {
    ID      string  `json:"id"`
    Valor   string  `json:"valor"` // Formatado para JSON
}
```

### 4. Middleware Chain

Middleware para cross-cutting concerns:

```go
app.Use(middleware.Logger())
app.Use(middleware.Recovery())
app.Use(middleware.CORSMiddleware())
app.Use(middleware.AuthMiddleware())
app.Use(middleware.TenantMiddleware())
```

### 5. Service Locator (Opcional)

Para inicialização centralizadas:

```go
type Container struct {
    DB              *sql.DB
    Logger          *zap.Logger
    ReceitaRepo     domain.ReceitaRepository
    DespesaRepo     domain.DespesaRepository
    // ... outros services
}
```

---

## 🔄 Fluxo de Dados

### Fluxo de Requisição HTTP

```
Request HTTP
    ↓
NGINX (Rate Limit, SSL)
    ↓
Echo Router
    ↓
Middleware Chain
  ├── Logger
  ├── Recovery
  ├── Auth (JWT)
  └── Tenant Context
    ↓
Handler (HTTP Layer)
    ├── Bind Request
    ├── Validate Input (Validator)
    └── Call Use Case
    ↓
Use Case (Application Layer)
    ├── Business Logic Validation
    ├── Call Domain Services
    └── Call Repositories
    ↓
Domain Layer
    ├── Business Rules
    ├── Value Object Creation
    └── Entity Validation
    ↓
Repository (Infrastructure)
    └── Database Query (SQLC)
    ↓
Response DTO
    ↓
JSON Response
```

### Fluxo de Processamento Assíncrono (Cron)

```
Scheduler (robfig/cron)
    ↓
Cron Job (ex: Sincronizar Asaas)
    ↓
Use Case (Application Layer)
    ├── Buscar faturas no Asaas
    ├── Mapear para Receitas
    └── Persistir no DB
    ↓
Notificação (opcional)
    └── Log ou Webhook
```

---

## 👥 Multi-Tenancy

### Modelo Selecionado: Column-Based (Tenant per Row)

**Razão**: Simplicidade, segurança, sem complexidade de schema separados.

### Implementação

1. **Coluna tenant_id em todas as tabelas**

```sql
CREATE TABLE receitas (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(10, 2) NOT NULL,
    data DATE NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(id, tenant_id)
);

CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data);
```

2. **Middleware de Tenant**

```go
func TenantMiddleware(c echo.Context) error {
    token := c.Get("user").(*jwt.Token)
    claims := token.Claims.(jwt.MapClaims)
    
    tenantID := claims["tenant_id"].(string)
    c.Set("tenant_id", tenantID)
    
    return c.Next()
}
```

3. **Query Segura**

```go
func (r *PostgresReceitaRepository) FindByTenantAndPeriod(
    ctx context.Context, tenantID string, from, to time.Time) ([]*Receita, error) {
    // Always filter by tenant_id
    query := `
        SELECT id, tenant_id, descricao, valor, data
        FROM receitas
        WHERE tenant_id = $1 AND data BETWEEN $2 AND $3
        ORDER BY data DESC
    `
    return r.db.QueryContext(ctx, query, tenantID, from, to)
}
```

---

## 🔐 Segurança

### Autenticação

- **JWT com RS256** (assimétrico)
- **Refresh Token** com rotação
- **Expiração**: Access Token 15 min, Refresh Token 7 dias

### Autorização

- **Role-Based Access Control (RBAC)**
- **Roles**: Owner, Manager, Employee, Accountant
- **Policies** por contexto (ex: barbeiro vê só suas finanças)

### Isolamento de Dados

- ✅ Sempre filtrar queries por `tenant_id`
- ✅ Validar propriedade de recursos
- ✅ Audit logs em operações sensíveis

### Rate Limiting

- **NGINX**: 100 req/s por IP
- **Aplicação**: 50 req/min por endpoint sensível

### HTTPS/TLS

- **Certificados**: Let's Encrypt + Certbot
- **HSTS**: 1 ano
- **CSP**: Restritivo para frontend

---

## 📈 Escalabilidade

### Banco de Dados

- **Índices estratégicos** em `tenant_id`, datas, status
- **Particionamento** de tabelas largas (receitas, despesas) por ano
- **Connection pooling** via pgBouncer (futuro)
- **Read replicas** no Neon (futuro)

### Backend

- **Stateless API** (escalável horizontalmente)
- **Cache de leitura** (Redis, futuro) para dashboards
- **Bulk operations** com batch inserts
- **Async jobs** fora do request cycle

### Frontend

- **Code splitting** automático no Next.js
- **Image optimization** com next/image
- **CDN** para assets estáticos
- **ISR** (Incremental Static Regeneration) para dashboards

### Monitoramento

- **Prometheus** para métricas
- **Grafana** para dashboards
- **Alertas** para SLA violations
- **Logs centralizados** em Loki ou Datadog

---

## 🔗 Referências

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Echo Framework](https://echo.labstack.com/)
- [Go Best Practices](https://golang.org/doc/effective_go)

---

**Última Atualização:** 14/11/2025  
**Status:** ✅ Aprovado para Implementação
