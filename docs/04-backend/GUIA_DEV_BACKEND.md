> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)
> **Atualizado em: 22/11/2025 - Implementação Completa dos 44 Endpoints** ✅

# 🔧 Guia de Desenvolvimento - Backend (Go)

**Versão:** 2.0
**Data Atualização:** 22/11/2025
**Status:** ✅ 44/44 Endpoints Implementados

---

## 🎉 ATUALIZAÇÃO IMPORTANTE (22/11/2025)

**TODOS OS 44 ENDPOINTS FORAM IMPLEMENTADOS E ESTÃO FUNCIONAIS!**

✅ **Metas** (15 endpoints) - MetaMensal, MetaBarbeiro, MetaTicketMedio
✅ **Precificação** (9 endpoints) - Config + Simulações
✅ **Financeiro** (20 endpoints) - ContaPagar, ContaReceber, Compensação, FluxoCaixa, DRE

Ver detalhes completos em:

- `/Tarefas/01-BLOQUEIOS-BASE/VERTICAL_SLICE_ALL_MODULES.md`
- `/Tarefas/01-BLOQUEIOS-BASE/README.md`

---

## 📋 Índice

1. [Setup Local](#setup-local)
2. [Estrutura de Projeto](#estrutura-de-projeto)
3. [Convenções de Código](#convenções-de-código)
4. [Desenvolvimento](#desenvolvimento)
5. [Testing](#testing)
6. [Deployment](#deployment)

---

## 🚀 Setup Local

### Pré-requisitos

```bash
# Verificar Go
go version  # Mínimo: 1.24

# PostgreSQL
psql --version  # Mínimo: 14

# Ferramentas
brew install golang-migrate
brew install sqlc
```

### Clone e Setup

```bash
# 1. Clone repositório
git clone https://github.com/seu-usuario/barber-analytics-backend-v2.git
cd barber-analytics-backend-v2

# 2. Copy .env
cp .env.example .env
# Editar com DATABASE_URL local

# 3. Instalar dependências
go mod download

# 4. Rodar migrations
migrate -path ./migrations -database $DATABASE_URL up

# 5. Rodar aplicação
go run ./cmd/api/main.go
```

### Health Check

```bash
curl http://localhost:8080/health
# {"status":"ok","timestamp":"2024-11-14T10:30:00Z"}
```

---

## 📁 Estrutura de Projeto

```
backend/
│
├── cmd/
│   └── api/
│       └── main.go                  # Entry point
│
├── internal/                         # Código privado do pacote
│   ├── config/
│   │   └── config.go               # Configuração
│   │
│   ├── domain/                     # Business logic (entities, services)
│   │   ├── tenant/
│   │   ├── user/
│   │   ├── financial/
│   │   │   ├── receita.go
│   │   │   ├── despesa.go
│   │   │   └── error.go
│   │   └── subscription/
│   │
│   ├── application/                # Use cases, DTOs
│   │   ├── dto/
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   ├── mapper/
│   │   │   └── receita_mapper.go
│   │   └── usecase/
│   │       ├── financial/
│   │       │   ├── create_receita.go
│   │       │   └── list_receitas.go
│   │       └── subscription/
│   │
│   ├── infrastructure/             # Implementações concretas
│   │   ├── http/                   # HTTP handlers
│   │   │   ├── handler/
│   │   │   │   ├── receita.go
│   │   │   │   └── subscription.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   └── tenant.go
│   │   │   └── route.go
│   │   │
│   │   ├── repository/             # Database
│   │   │   ├── receita_repository.go
│   │   │   └── subscription_repository.go
│   │   │
│   │   ├── external/               # Integrações externas
│   │   │   └── asaas/
│   │   │       └── client.go
│   │   │
│   │   └── scheduler/              # Cron jobs
│   │       ├── scheduler.go
│   │       └── jobs/
│   │           ├── sync_asaas.go
│   │           └── financial_snapshot.go
│   │
│   └── ports/                      # Interfaces (abstrações)
│       ├── http_handler.go
│       └── repository.go
│
├── migrations/                      # SQL migrations
│   ├── 001_create_tenants.up.sql
│   └── 001_create_tenants.down.sql
│
├── tests/                          # Testes integrados
│   ├── integration/
│   │   └── receita_test.go
│   └── fixtures/
│       └── seed.sql
│
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md
```

---

## 🎯 Convenções de Código

### Naming

```go
// Pacotes: minúsculas, sem underscores
package financial

// Interfaces: PascalCase, suffix -er ou -or
type ReceitaRepository interface {}
type PasswordHasher interface {}

// Structs: PascalCase
type Receita struct {}

// Funções: camelCase (exportadas), PascalCase (exportadas)
func (r *Receita) Cancel() error {}
func NewReceita() *Receita {}

// Constantes: UPPER_SNAKE_CASE
const (
    StatusActive = \"ACTIVE\"
    StatusInactive = \"INACTIVE\"
)

// Variáveis privadas: camelCase
var receitaRepo ReceitaRepository
```

### Error Handling

```go
// ✅ CORRETO: Wrap errors com contexto
if err != nil {
    return fmt.Errorf(\"failed to save receita: %w\", err)
}

// ❌ ERRADO: Ignorar erros
_ = receita.Save()

// ✅ CORRETO: Custom errors
var ErrReceitaNotFound = errors.New(\"receita not found\")

// Usar
if err == ErrReceitaNotFound {
    return c.JSON(404, ErrorResponse{})
}
```

### Contexto

```go
// Sempre passar context como primeiro argumento
func (r *ReceitaRepository) FindByID(ctx context.Context, id string) (*Receita, error) {
    // Respeitar context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
}

// Com timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## 💻 Desenvolvimento

### Criar novo Use Case

1. **Definir DTO**

```go
// internal/application/dto/create_receita.go
type CreateReceitaRequest struct {
    Descricao string `json:\"descricao\" validate:\"required,max=255\"`
    Valor     string `json:\"valor\" validate:\"required,numeric\"`
    Data      time.Time `json:\"data\" validate:\"required\"`
}

type CreateReceitaResponse struct {
    ID    string `json:\"id\"`
    Status string `json:\"status\"`
}
```

2. **Criar Use Case**

```go
// internal/application/usecase/financial/create_receita.go
type CreateReceitaUseCase struct {
    repository domain.ReceitaRepository
    validator  *validator.Validator
}

func NewCreateReceitaUseCase(
    repo domain.ReceitaRepository,
    val *validator.Validator) *CreateReceitaUseCase {
    return &CreateReceitaUseCase{
        repository: repo,
        validator:  val,
    }
}

func (uc *CreateReceitaUseCase) Execute(
    ctx context.Context,
    tenantID string,
    req *dto.CreateReceitaRequest) (*dto.CreateReceitaResponse, error) {

    // Validar
    if err := uc.validator.Struct(req); err != nil {
        return nil, fmt.Errorf(\"validation error: %w\", err)
    }

    // Converter valor
    valor, err := decimal.NewFromString(req.Valor)
    if err != nil {
        return nil, errors.New(\"invalid valor format\")
    }

    // Criar domain entity
    receita := &domain.Receita{...}

    // Persistir
    if err := uc.repository.Save(ctx, tenantID, receita); err != nil {
        return nil, err
    }

    return &dto.CreateReceitaResponse{
        ID:     receita.ID,
        Status: string(receita.Status),
    }, nil
}
```

3. **Criar Handler**

```go
// internal/infrastructure/http/handler/receita.go
type ReceitaHandler struct {
    createUC *application.CreateReceitaUseCase
}

func (h *ReceitaHandler) Create(c echo.Context) error {
    var req dto.CreateReceitaRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, ErrorResponse{Message: err.Error()})
    }

    tenantID := c.Get(\"tenant_id\").(string)

    resp, err := h.createUC.Execute(c.Request().Context(), tenantID, &req)
    if err != nil {
        return c.JSON(500, ErrorResponse{Message: err.Error()})
    }

    return c.JSON(201, resp)
}
```

4. **Registrar Route**

```go
// internal/infrastructure/http/route.go
func setupRoutes(e *echo.Echo, handlers *Handlers) {
    // Grupo protegido
    api := e.Group(\"/api/v2\")
    api.Use(middleware.Auth)
    api.Use(middleware.Tenant)

    // Rotas de receita
    api.POST(\"/financial/receitas\", handlers.Receita.Create)
    api.GET(\"/financial/receitas\", handlers.Receita.List)
}
```

---

## 🧪 Testing

### Unit Tests

```go
// internal/application/usecase/financial/create_receita_test.go
package financial

import (
    \"testing\"
    \"github.com/stretchr/testify/assert\"
)

func TestCreateReceitaUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := &mockReceitaRepository{}
    uc := NewCreateReceitaUseCase(mockRepo, validator.New())

    req := &dto.CreateReceitaRequest{
        Descricao: \"Corte de cabelo\",
        Valor:     \"50.00\",
        Data:      time.Now(),
    }

    // Act
    resp, err := uc.Execute(context.Background(), \"tenant-123\", req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotEmpty(t, resp.ID)
}

func TestCreateReceitaUseCase_Execute_InvalidValue(t *testing.T) {
    // Arrange
    uc := NewCreateReceitaUseCase(mockRepo, validator.New())

    req := &dto.CreateReceitaRequest{
        Valor: \"invalid\",
    }

    // Act
    _, err := uc.Execute(context.Background(), \"tenant-123\", req)

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), \"validation error\")
}
```

### Rodar Testes

```bash
# Todos os testes
make test

# Com coverage
make test-coverage

# Teste específico
go test -run TestCreateReceitaUseCase ./...

# Verbose
go test -v ./...
```

---

## 🚢 Deployment

### Build

```bash
# Build local
make build
```

### Deploy Staging

```bash
# CI/CD automatizado via GitHub Actions
# Manual:
git push origin develop

# GitHub Actions: build → test → deploy
```

### Verificar Deploy

```bash
# Health check
curl https://api-staging.seudominio.com/v2/health

# Verificar logs
ssh ubuntu@vps.com
journalctl -u barber-api -f
```

---

## 📚 Ferramentas Úteis

```bash
# Format código
go fmt ./...

# Lint
golangci-lint run ./...

# Atualizar dependências
go get -u ./...

# Gerar mocks (mockgen)
mockgen -source=internal/domain/receita.go -destination=mocks/mock_receita.go

# Análise estática
go vet ./...

# Benchmark
go test -bench=. -benchmem ./...
```

---

**Status:** ✅ Guia completo
