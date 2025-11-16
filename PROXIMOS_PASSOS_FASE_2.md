# 🚀 Próximos Passos: Fase 2 (Backend Core)

**Status:** Fase 1 ✅ Completa  
**Próxima:** Fase 2 - Backend Core  
**Timeline:** 7-14 dias  
**Data Início:** 15/11/2025

---

## 📋 Fase 2: Roadmap Detalhado

### Semana 1: Foundation (Days 1-7)

#### T-BE-001: Config Management (2h)
```go
// internal/config/config.go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Logger   LoggerConfig
    JWT      JWTConfig
    Asaas    AsaasConfig
}

// Ler de .env:
// DATABASE_URL
// PORT
// LOG_LEVEL
// JWT_PRIVATE_KEY_PATH
// JWT_PUBLIC_KEY_PATH
// ASAAS_API_KEY
```

**Arquivo:** `/backend/internal/config/config.go`

#### T-BE-002: Database Connection (1h)
```go
// internal/config/database.go
func NewDatabase(cfg DatabaseConfig) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.URL)
    if err != nil {
        return nil, err
    }
    
    // Pool tuning
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

**Arquivo:** `/backend/internal/config/database.go`

#### T-BE-003: Logger Setup (1h)
```go
// internal/config/logger.go
func NewLogger(level string) *zap.Logger {
    // Structured JSON logging
    // Levels: debug, info, warn, error
    // Output: stdout + file (optional)
}
```

**Arquivo:** `/backend/internal/config/logger.go`

#### T-BE-004: Domain Layer - Entities (3h)

**Estrutura esperada:**
```
internal/domain/
├── tenant/
│   └── tenant.go           (Tenant entity)
├── user/
│   └── user.go             (User entity + roles)
├── financial/
│   ├── receita.go          (Receita entity)
│   ├── despesa.go          (Despesa entity)
│   ├── categoria.go
│   └── money.go            (Value Object)
└── subscription/
    └── subscription.go     (Subscription entity)
```

**Exemplo - Tenant Entity:**
```go
// internal/domain/tenant/tenant.go
package tenant

import "time"

type Tenant struct {
    ID        string
    Nome      string
    CNPJ      string // nullable
    Ativo     bool
    Plano     string // free, pro, enterprise
    CriadoEm  time.Time
    AtualizadoEm time.Time
}

// Validações
func (t *Tenant) Validate() error {
    if t.Nome == "" {
        return errors.New("nome obrigatório")
    }
    // ... mais validações
    return nil
}
```

#### T-BE-005: Domain Layer - Value Objects (2h)

**Exemplos:**
```go
// Money Value Object (imutável)
type Money struct {
    amount    decimal.Decimal
    currency  string
}

// Email Value Object
type Email struct {
    value string
}

// Role Value Object
type Role string
const (
    RoleOwner      Role = "owner"
    RoleManager         = "manager"
    RoleEmployee        = "employee"
    RoleBarbeiro        = "barbeiro"
)
```

#### T-BE-006: Repository Interfaces (2h)

**Exemplo:**
```go
// internal/ports/repository.go
package ports

type TenantRepository interface {
    Save(ctx context.Context, tenant *domain.Tenant) error
    FindByID(ctx context.Context, id string) (*domain.Tenant, error)
    FindByName(ctx context.Context, name string) (*domain.Tenant, error)
}

type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

// ... mais interfaces
```

### Semana 2: Authentication (Days 8-14)

#### T-BE-007: Auth Domain Service (2h)
```go
// internal/domain/user/auth_service.go
type AuthService interface {
    GenerateJWT(user *User) (accessToken, refreshToken string, err error)
    ValidatePassword(hashedPassword, password string) bool
    HashPassword(password string) (string, error)
}
```

#### T-BE-008: Auth Use Cases (4h)
```go
// internal/application/usecase/auth/login.go
type LoginInput struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=6"`
}

type LoginUseCase struct {
    userRepo      ports.UserRepository
    authService   domain.AuthService
    logger        *zap.Logger
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
    // 1. Validar input
    // 2. Buscar usuário por email
    // 3. Validar password
    // 4. Gerar JWT
    // 5. Retornar tokens + user data
}

// Também: RefreshTokenUseCase, CreateUserUseCase
```

#### T-BE-009: Auth Handlers (2h)
```go
// internal/infrastructure/http/handler/auth.go
func (h *AuthHandler) Login(c echo.Context) error {
    var input dto.LoginRequest
    if err := c.Bind(&input); err != nil {
        return c.JSON(400, ErrorResponse{Message: err.Error()})
    }
    
    output, err := h.loginUseCase.Execute(c.Request().Context(), input)
    if err != nil {
        return c.JSON(401, ErrorResponse{Message: "Invalid credentials"})
    }
    
    return c.JSON(200, output)
}
```

#### T-BE-010: Middlewares (2h)
```go
// internal/infrastructure/http/middleware/auth.go
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // 1. Extract token from header
        token := c.Request().Header.Get("Authorization")
        // 2. Validate JWT signature
        // 3. Extract claims (sub, tenant_id, role)
        // 4. Add to context
        // 5. Check expiration
        return next(c)
    }
}

// Também: TenantMiddleware, LoggerMiddleware
```

#### T-BE-011: DTO & Mappers (1h)
```go
// internal/application/dto/request.go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// internal/application/dto/response.go
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`
    User         UserResponse
}

// internal/application/mapper/user_mapper.go
func MapUserToDTO(user *domain.User) UserResponse {
    return UserResponse{
        ID:    user.ID,
        Email: user.Email,
        Nome:  user.Nome,
        Role:  string(user.Role),
    }
}
```

#### T-BE-012: Unit Tests (4h)
```go
// tests/unit/auth_test.go
func TestLoginUseCase_Execute_Success(t *testing.T) {
    // Arrange
    mockUserRepo := &MockUserRepository{}
    uc := usecase.NewLoginUseCase(mockUserRepo, mockAuthService, mockLogger)
    
    input := LoginInput{
        Email:    "user@example.com",
        Password: "password123",
    }
    
    // Act
    output, err := uc.Execute(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.NotEmpty(t, output.AccessToken)
}

func TestLoginUseCase_Execute_InvalidEmail(t *testing.T) {
    // Test error cases
}
```

---

## 📁 Estrutura de Diretórios Esperada

```
backend/
├── cmd/api/
│   └── main.go                     ← Initialize database, router, start server
│
├── internal/
│   ├── config/
│   │   ├── config.go               ← Load from env
│   │   ├── database.go
│   │   ├── logger.go
│   │   └── jwt.go
│   │
│   ├── domain/
│   │   ├── tenant/
│   │   │   └── tenant.go
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   └── auth_service.go
│   │   ├── financial/
│   │   │   ├── receita.go
│   │   │   ├── despesa.go
│   │   │   ├── categoria.go
│   │   │   ├── money.go            ← Value Object
│   │   │   └── calculator.go       ← Service
│   │   └── subscription/
│   │       └── subscription.go
│   │
│   ├── application/
│   │   ├── dto/
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── error.go
│   │   ├── mapper/
│   │   │   ├── user_mapper.go
│   │   │   ├── tenant_mapper.go
│   │   │   └── receita_mapper.go
│   │   └── usecase/
│   │       ├── auth/
│   │       │   ├── login.go
│   │       │   ├── refresh.go
│   │       │   └── create_user.go
│   │       └── tenant/
│   │           └── create_tenant.go
│   │
│   ├── infrastructure/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── auth.go
│   │   │   │   ├── tenant.go
│   │   │   │   └── health.go       ← Já existe!
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── tenant.go
│   │   │   │   ├── logger.go
│   │   │   │   └── recovery.go
│   │   │   └── route.go            ← Register all routes
│   │   │
│   │   ├── repository/
│   │   │   ├── tenant_repository.go
│   │   │   ├── user_repository.go
│   │   │   ├── receita_repository.go
│   │   │   └── despesa_repository.go
│   │   │
│   │   ├── external/
│   │   │   └── asaas/
│   │   │       └── client.go       ← Planejado para Fase 3
│   │   │
│   │   └── scheduler/
│   │       └── scheduler.go        ← Planejado para Fase 3
│   │
│   └── ports/
│       ├── repository.go           ← Interfaces
│       ├── service.go
│       └── handler.go
│
├── tests/
│   ├── unit/
│   │   ├── auth_test.go
│   │   ├── user_test.go
│   │   └── tenant_test.go
│   ├── integration/
│   │   └── auth_integration_test.go
│   └── fixtures/
│       └── seed.go
│
└── migrations/
    └── (já criadas ✅)
```

---

## 🎯 Checklist Fase 2

### Config & Database
- [ ] Config loading from .env
- [ ] Database connection pool
- [ ] Logger setup (Zap)
- [ ] Health check integration

### Domain Layer
- [ ] Tenant entity
- [ ] User entity + roles
- [ ] Receita entity
- [ ] Despesa entity
- [ ] Value objects (Money, Email, Role, etc)
- [ ] Domain services

### Application Layer
- [ ] DTOs (requests/responses)
- [ ] Mappers (domain ↔ DTO)
- [ ] Auth use cases (Login, Refresh, CreateUser)
- [ ] Input validation

### Infrastructure Layer
- [ ] Repository implementations
- [ ] Auth handler
- [ ] Middleware stack
- [ ] Route registration

### Testing
- [ ] Unit tests >80% coverage
- [ ] Mock repositories
- [ ] Integration tests (auth flow)

---

## 🚀 Como Começar Fase 2

### 1. Setup inicial
```bash
cd /home/andrey/projetos/barber-Analytic-proV2/backend

# Criar estrutura de pacotes
mkdir -p internal/{config,domain/{tenant,user,financial,subscription},application/{dto,mapper,usecase/auth},infrastructure/{http/{handler,middleware},repository,external/asaas,scheduler},ports}
mkdir -p tests/{unit,integration,fixtures}
```

### 2. Arquivo main.go
```go
package main

import (
    "log"
    "barber-analytics/internal/config"
    "barber-analytics/internal/infrastructure/http"
)

func main() {
    // 1. Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Setup logger
    logger, _ := config.NewLogger(cfg.Logger.Level)
    
    // 3. Connect database
    db, err := config.NewDatabase(cfg.Database)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // 4. Initialize repositories
    // 5. Initialize use cases
    // 6. Setup routes
    router := http.SetupRoutes(...)
    
    // 7. Start server
    router.Start(":" + cfg.Server.Port)
}
```

### 3. Environment variables (.env)
```env
DATABASE_URL=postgresql://neondb_owner:npg_***@ep-winter-leaf-adhqz08p-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require
PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json

JWT_PRIVATE_KEY_PATH=/opt/barber-api/keys/private.pem
JWT_PUBLIC_KEY_PATH=/opt/barber-api/keys/public.pem
JWT_ISSUER=barber-analytics-pro
JWT_AUDIENCE=barber-analytics-users

ASAAS_API_KEY=asaas_prod_123456
ASAAS_BASE_URL=https://www.asaas.com/api/v3
```

### 4. Começar com T-BE-001
```bash
# File: internal/config/config.go
touch internal/config/config.go
touch internal/config/database.go
touch internal/config/logger.go
touch internal/config/jwt.go

# Implementar LoadConfig() function
```

---

## 📊 Estimativa de Tempo

| Tarefa | Tempo | Crítica |
|--------|-------|---------|
| T-BE-001: Config | 2h | ✅ |
| T-BE-002: Database | 1h | ✅ |
| T-BE-003: Logger | 1h | ✅ |
| T-BE-004: Domain Entities | 3h | ✅ |
| T-BE-005: Value Objects | 2h | ✅ |
| T-BE-006: Repositories (interface) | 2h | ✅ |
| T-BE-007: Auth Service | 2h | ✅ |
| T-BE-008: Use Cases | 4h | ✅ |
| T-BE-009: Handlers | 2h | ✅ |
| T-BE-010: Middlewares | 2h | ✅ |
| T-BE-011: DTOs & Mappers | 1h | ✅ |
| T-BE-012: Unit Tests | 4h | ✅ |
| **TOTAL** | **26h** | **100%** |

**Timeline:** ~7 dias (4h/dia) ou 3-4 dias (full-time)

---

## 📞 Contato & Próximos Passos

**Quando Fase 2 estiver pronta:**
1. ✅ Banco de dados sincronizado com Go entities
2. ✅ Auth funcionando com JWT RS256
3. ✅ Middlewares implementadas
4. ✅ Health check retornando status do banco
5. ✅ >80% coverage em testes

**Então partir para Fase 3:**
- Módulos financeiro (receitas, despesas, fluxo de caixa)
- Integração Asaas
- Cron jobs

---

**🚀 Fase 2 Começa em:** 15/11/2025  
**Estimado em:** 7-14 dias  
**Pronto?** Então vamos para Fase 2! 💪
