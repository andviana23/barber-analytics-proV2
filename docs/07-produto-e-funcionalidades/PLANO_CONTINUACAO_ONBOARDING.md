> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🚀 Plano de Continuação - Onboarding Flow

**Data:** 20/11/2025
**Prioridade:** 🔴 CRÍTICA
**Tempo Estimado Total:** 4-6 horas

---

## 📊 Status Atual

### ✅ Completado (80%)

- Database migration aplicada
- Entity & Repository atualizados
- SignupUseCase implementado
- Frontend signup page funcional
- Frontend onboarding page funcional
- AuthContext integrado

### ❌ Pendente (20%)

- **Backend:** CompleteOnboardingUseCase
- **Backend:** TenantHandler
- **Backend:** Dependency Injection
- **Validações:** CNPJ/Email duplicados
- **Testes:** Unit + Integration + E2E

---

## 🎯 Fase 1: Backend Core (CRÍTICO)

**Tempo:** 2 horas
**Objetivo:** Implementar endpoint `POST /tenants/onboarding/complete`

### Task 1.1: Complete Onboarding Use Case

**Arquivo:** `backend/internal/application/usecase/tenant/complete_onboarding_usecase.go`

```bash
# Criar diretório se não existir
mkdir -p backend/internal/application/usecase/tenant

# Criar arquivo
touch backend/internal/application/usecase/tenant/complete_onboarding_usecase.go
```

**Código:**

```go
package tenant

import (
    "context"
    "fmt"
    "time"

    "github.com/andviana23/barber-analytics-backend-v2/internal/domain/entity"
    "github.com/andviana23/barber-analytics-backend-v2/internal/domain/repository"
)

type CompleteOnboardingUseCase struct {
    tenantRepo repository.TenantRepository
}

func NewCompleteOnboardingUseCase(
    tenantRepo repository.TenantRepository,
) *CompleteOnboardingUseCase {
    return &CompleteOnboardingUseCase{
        tenantRepo: tenantRepo,
    }
}

func (uc *CompleteOnboardingUseCase) Execute(ctx context.Context, tenantID string) error {
    // 1. Validar input
    if tenantID == "" {
        return fmt.Errorf("tenant ID is required")
    }

    // 2. Buscar tenant
    tenant, err := uc.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return fmt.Errorf("tenant not found: %w", err)
    }

    // 3. Validar se já completou (idempotência)
    if tenant.OnboardingCompleted {
        return nil // Já completado, retornar sucesso
    }

    // 4. Marcar como completado
    tenant.OnboardingCompleted = true
    tenant.AtualizadoEm = time.Now()

    // 5. Persistir
    if err := uc.tenantRepo.Update(ctx, tenant); err != nil {
        return fmt.Errorf("failed to update tenant: %w", err)
    }

    return nil
}
```

**Validação:**

```bash
# Compilar para verificar erros
cd backend
go build ./internal/application/usecase/tenant/
```

---

### Task 1.2: Tenant Handler

**Arquivo:** `backend/internal/infrastructure/http/handler/tenant_handler.go`

```bash
touch backend/internal/infrastructure/http/handler/tenant_handler.go
```

**Código:**

```go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/andviana23/barber-analytics-backend-v2/internal/application/usecase/tenant"
    "github.com/andviana23/barber-analytics-backend-v2/internal/domain/service"
    httpMiddleware "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/http/middleware"
    "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/http/response"
)

type TenantHandler struct {
    completeOnboardingUC *tenant.CompleteOnboardingUseCase
    jwtService           *service.JWTService
}

func NewTenantHandler(
    completeOnboardingUC *tenant.CompleteOnboardingUseCase,
    jwtService *service.JWTService,
) *TenantHandler {
    return &TenantHandler{
        completeOnboardingUC: completeOnboardingUC,
        jwtService:           jwtService,
    }
}

func (h *TenantHandler) RegisterRoutes(r chi.Router) {
    r.Route("/tenants", func(r chi.Router) {
        // Rotas protegidas
        r.Group(func(r chi.Router) {
            if h.jwtService != nil {
                r.Use(httpMiddleware.ChiAuthMiddleware(h.jwtService))
                r.Use(httpMiddleware.ChiTenantMiddleware())
            }

            r.Post("/onboarding/complete", h.handleCompleteOnboarding)
        })
    })
}

func (h *TenantHandler) handleCompleteOnboarding(w http.ResponseWriter, r *http.Request) {
    // Extrair tenant ID do contexto
    tenantID, err := getTenantIDFromRequest(r)
    if err != nil {
        writeStandardResponse(w, response.Error("FORBIDDEN", "Tenant ID not found in context", err.Error(), ""))
        return
    }

    // Executar use case
    if err := h.completeOnboardingUC.Execute(r.Context(), tenantID); err != nil {
        writeStandardResponse(w, response.Error("INTERNAL_ERROR", "Failed to complete onboarding", err.Error(), ""))
        return
    }

    // Retornar sucesso
    writeStandardResponse(w, response.Success("OK", "Onboarding completed successfully", nil, ""))
}
```

**Validação:**

```bash
go build ./internal/infrastructure/http/handler/
```

---

### Task 1.3: Dependency Injection

**Arquivo:** `backend/cmd/api/main.go`

**Localizar seção de handlers e adicionar:**

```go
// Existing handlers
authHandler := handler.NewAuthHandler(
    loginUC,
    refreshTokenUC,
    createUserUC,
    getMeUC,
    signupUC,
    jwtService,
)

// ✅ ADICIONAR: Tenant Handler
completeOnboardingUC := tenant.NewCompleteOnboardingUseCase(tenantRepo)
tenantHandler := handler.NewTenantHandler(completeOnboardingUC, jwtService)

// Register routes
authHandler.RegisterRoutes(router)
tenantHandler.RegisterRoutes(router) // ✅ ADICIONAR
```

**Validação:**

```bash
# Compilar aplicação completa
make build

# Rodar
make run-backend
```

---

### Task 1.4: Teste Manual

```bash
# Terminal 1: Rodar backend
cd backend
make run-backend

# Terminal 2: Testar signup
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "barberName": "Barbearia Teste",
    "cnpj": "12345678000199",
    "email": "teste@barbearia.com",
    "password": "senha123",
    "name": "João Silva"
  }'

# Copiar token retornado (access_token) e testar onboarding
curl -X POST http://localhost:8080/api/v1/tenants/onboarding/complete \
  -H "Authorization: Bearer {SEU_TOKEN_AQUI}" \
  -H "Content-Type: application/json"

# Verificar no banco
psql $DATABASE_URL -c "SELECT id, nome, onboarding_completed FROM tenants WHERE nome = 'Barbearia Teste';"
```

**Resultado esperado:**

```json
{
  "code": "OK",
  "message": "Onboarding completed successfully",
  "data": null,
  "errors": null,
  "timestamp": "2025-11-20T..."
}
```

---

## 🎯 Fase 2: Validações (IMPORTANTE)

**Tempo:** 1 hora
**Objetivo:** Evitar duplicação de CNPJ/Email

### Task 2.1: Validar Duplicados em SignupUseCase

**Arquivo:** `backend/internal/application/usecase/auth/signup_usecase.go`

**Modificar método `Execute`:**

```go
func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    // 1. Validate Input
    if input.BarberName == "" || input.Email == "" || input.Password == "" || input.Name == "" {
        return nil, errors.New("missing required fields")
    }

    // ✅ ADICIONAR: Validar CNPJ duplicado
    if input.CNPJ != "" {
        existing, _ := uc.tenantRepo.FindByCNPJ(ctx, input.CNPJ)
        if existing != nil {
            return nil, errors.New("CNPJ already registered")
        }
    }

    // ✅ ADICIONAR: Validar Email duplicado
    existingUser, _ := uc.userRepo.FindByEmailAny(ctx, input.Email)
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    // ... resto do código permanece igual
}
```

### Task 2.2: Retornar Erro HTTP Apropriado

**Arquivo:** `backend/internal/infrastructure/http/handler/auth_handler.go`

**Modificar `handleSignup`:**

```go
func (h *AuthHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
    var req auth.SignupInput

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeStandardResponse(w, response.Error("BAD_REQUEST", "Invalid request body", err.Error(), ""))
        return
    }

    resp, err := h.signupUseCase.Execute(r.Context(), req)
    if err != nil {
        // ✅ ADICIONAR: Tratamento de duplicados
        if strings.Contains(err.Error(), "already registered") {
            writeStandardResponse(w, response.Error("CONFLICT", err.Error(), nil, ""))
            return
        }

        writeStandardResponse(w, response.Error("INTERNAL_ERROR", "Failed to signup", err.Error(), ""))
        return
    }

    writeStandardResponse(w, response.Success("CREATED", "Signup successful", resp, ""))
}
```

---

## 🎯 Fase 3: Testes Automatizados (RECOMENDADO)

**Tempo:** 2-3 horas

### Task 3.1: Unit Test - CompleteOnboardingUseCase

**Arquivo:** `backend/internal/application/usecase/tenant/complete_onboarding_usecase_test.go`

```go
package tenant_test

import (
    "context"
    "testing"
    "time"

    "github.com/andviana23/barber-analytics-backend-v2/internal/application/usecase/tenant"
    "github.com/andviana23/barber-analytics-backend-v2/internal/domain/entity"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type mockTenantRepository struct {
    mock.Mock
}

func (m *mockTenantRepository) FindByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
    args := m.Called(ctx, tenantID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *mockTenantRepository) Update(ctx context.Context, tenant *entity.Tenant) error {
    args := m.Called(ctx, tenant)
    return args.Error(0)
}

func TestCompleteOnboardingUseCase_Execute_Success(t *testing.T) {
    // Arrange
    mockRepo := new(mockTenantRepository)
    uc := tenant.NewCompleteOnboardingUseCase(mockRepo)

    testTenant := &entity.Tenant{
        ID:                  "tenant-123",
        Nome:                "Barbearia Teste",
        OnboardingCompleted: false,
        CriadoEm:            time.Now(),
        AtualizadoEm:        time.Now(),
    }

    mockRepo.On("FindByID", mock.Anything, "tenant-123").Return(testTenant, nil)
    mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

    // Act
    err := uc.Execute(context.Background(), "tenant-123")

    // Assert
    assert.NoError(t, err)
    assert.True(t, testTenant.OnboardingCompleted)
    mockRepo.AssertExpectations(t)
}

func TestCompleteOnboardingUseCase_Execute_AlreadyCompleted(t *testing.T) {
    // Arrange
    mockRepo := new(mockTenantRepository)
    uc := tenant.NewCompleteOnboardingUseCase(mockRepo)

    testTenant := &entity.Tenant{
        ID:                  "tenant-123",
        OnboardingCompleted: true, // Já completado
    }

    mockRepo.On("FindByID", mock.Anything, "tenant-123").Return(testTenant, nil)
    // Update NÃO deve ser chamado

    // Act
    err := uc.Execute(context.Background(), "tenant-123")

    // Assert
    assert.NoError(t, err)
    mockRepo.AssertNotCalled(t, "Update")
}
```

**Rodar:**

```bash
go test ./internal/application/usecase/tenant/ -v
```

---

### Task 3.2: Integration Test - Fluxo Completo

**Arquivo:** `backend/tests/integration/onboarding_flow_test.go`

```go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestOnboardingFlow_EndToEnd(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"

    // 1. Signup
    signupPayload := map[string]string{
        "barberName": "Barbearia E2E Test",
        "cnpj":       "99999999000199",
        "email":      "e2e@test.com",
        "password":   "senha123",
        "name":       "E2E User",
    }
    signupBody, _ := json.Marshal(signupPayload)

    resp, err := http.Post(baseURL+"/auth/signup", "application/json", bytes.NewBuffer(signupBody))
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, resp.StatusCode)

    var signupResp map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&signupResp)
    data := signupResp["data"].(map[string]interface{})
    token := data["token"].(string)

    assert.NotEmpty(t, token)

    // 2. Verificar onboarding_completed = false
    req, _ := http.NewRequest("GET", baseURL+"/auth/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err = http.DefaultClient.Do(req)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

    // 3. Complete onboarding
    req, _ = http.NewRequest("POST", baseURL+"/tenants/onboarding/complete", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err = http.DefaultClient.Do(req)
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 4. Verificar onboarding_completed = true
    req, _ = http.NewRequest("GET", baseURL+"/auth/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err = http.DefaultClient.Do(req)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

    // TODO: Parsear response e verificar tenant.onboarding_completed = true
}
```

**Rodar:**

```bash
# Rodar backend em background
make run-backend &

# Rodar integration tests
go test ./tests/integration/ -v

# Matar backend
pkill barber-api
```

---

### Task 3.3: Frontend E2E Test

**Arquivo:** `frontend/e2e/onboarding-flow.spec.ts`

```typescript
import { test, expect } from "@playwright/test";

test.describe("Onboarding Flow", () => {
  const timestamp = Date.now();
  const testEmail = `e2e-${timestamp}@test.com`;
  const testCNPJ = `${timestamp}`.padStart(14, "0");

  test("should complete full signup and onboarding flow", async ({ page }) => {
    // 1. Navigate to signup
    await page.goto("/signup");

    // 2. Fill signup form
    await page.fill(
      '[data-testid="barber-name-input"]',
      `Barbearia E2E ${timestamp}`
    );
    await page.fill('[data-testid="cnpj-input"]', testCNPJ);
    await page.fill('[data-testid="name-input"]', "E2E Test User");
    await page.fill('[data-testid="email-input"]', testEmail);
    await page.fill('[data-testid="password-input"]', "senha123");

    // 3. Submit signup
    await page.click('[data-testid="signup-button"]');

    // 4. Should redirect to onboarding
    await expect(page).toHaveURL("/onboarding", { timeout: 5000 });

    // 5. Verify welcome message
    await expect(
      page.locator("text=Bem-vindo ao Barber Analytics Pro")
    ).toBeVisible();

    // 6. Complete onboarding
    await page.click("text=Começar a Usar");

    // 7. Should redirect to dashboard
    await expect(page).toHaveURL("/dashboard", { timeout: 5000 });

    // 8. Verify dashboard loaded
    await expect(page.locator("text=Dashboard")).toBeVisible();
  });

  test("should prevent duplicate signup with same email", async ({ page }) => {
    // Usar mesmo email do teste anterior (se ainda existir)
    await page.goto("/signup");

    await page.fill('[data-testid="email-input"]', testEmail);
    await page.fill('[data-testid="password-input"]', "senha123");
    await page.fill('[data-testid="barber-name-input"]', "Outra Barbearia");
    await page.fill('[data-testid="cnpj-input"]', "11111111000111");
    await page.fill('[data-testid="name-input"]', "Outro Nome");

    await page.click('[data-testid="signup-button"]');

    // Deve mostrar erro de email já registrado
    await expect(page.locator("text=/email already registered/i")).toBeVisible({
      timeout: 3000,
    });
  });
});
```

**Rodar:**

```bash
cd frontend
npm run test:e2e
```

---

## 📋 Checklist de Execução

### Fase 1: Backend Core (CRÍTICO)

- [ ] Criar `CompleteOnboardingUseCase`
- [ ] Criar `TenantHandler`
- [ ] Registrar routes no `main.go`
- [ ] Compilar sem erros (`make build`)
- [ ] Testar endpoint com curl
- [ ] Verificar banco de dados (onboarding_completed = true)

### Fase 2: Validações

- [ ] Adicionar validação CNPJ duplicado
- [ ] Adicionar validação Email duplicado
- [ ] Retornar HTTP 409 Conflict
- [ ] Testar signup com duplicados

### Fase 3: Testes

- [ ] Unit test: CompleteOnboardingUseCase
- [ ] Integration test: Fluxo completo
- [ ] E2E test: Frontend flow
- [ ] Rodar `make test` sem erros

---

## 🚀 Comandos Rápidos

```bash
# Setup completo
cd backend
mkdir -p internal/application/usecase/tenant
touch internal/application/usecase/tenant/complete_onboarding_usecase.go
touch internal/infrastructure/http/handler/tenant_handler.go

# Build & Test
make build
make test
make run-backend

# Frontend
cd frontend
npm run dev
npm run test:e2e

# Verificar DB
psql $DATABASE_URL -c "SELECT * FROM tenants WHERE onboarding_completed = true;"
```

---

## 📊 Métricas de Sucesso

- ✅ Endpoint `/tenants/onboarding/complete` retorna 200 OK
- ✅ Coluna `onboarding_completed` atualizada para `true` no banco
- ✅ Frontend redireciona para `/dashboard` após completar
- ✅ Testes automatizados passam (unit + integration + e2e)
- ✅ Signups duplicados retornam 409 Conflict

---

**Próxima Ação:** Implementar Fase 1 - Task 1.1 (CompleteOnboardingUseCase)
