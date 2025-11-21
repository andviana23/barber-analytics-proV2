> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📋 Revisão Arquitetural - Onboarding Flow (T-PROD-003)

**Data:** 20/11/2025
**Status:** 🟡 Parcialmente Implementado
**Próxima Fase:** Backend - Complete Onboarding Use Case

---

## 🎯 Objetivo

Implementar fluxo completo de signup e onboarding para novos tenants, desde cadastro inicial até configuração completa com redirecionamento ao dashboard.

---

## ✅ O Que Já Está Implementado

### 1. **Database Layer** ✅ COMPLETO

- **Migration 024:** `onboarding_completed` adicionada à tabela `tenants`
- **Status:** ✅ APLICADO NO NEON (confirmado no arquivo)
- **Entity:** `Tenant` já possui campo `OnboardingCompleted bool`
- **Repository:** `PostgresTenantRepository` já persiste o campo em todas operações

**Validação:**

```sql
-- Verificar coluna existe
SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_name = 'tenants' AND column_name = 'onboarding_completed';
```

---

### 2. **Backend - Signup Use Case** ✅ COMPLETO

**Arquivo:** `backend/internal/application/usecase/auth/signup_usecase.go`

**Implementado:**

- ✅ Validação de inputs (BarberName, CNPJ, Email, Password)
- ✅ Criação de Tenant (com `OnboardingCompleted: false`)
- ✅ Criação de User (role: owner)
- ✅ Geração de JWT token (auto-login)
- ✅ Retorno de `SignupOutput` com TenantID, UserID, Token

**Status:** ✅ FUNCIONANDO

**Pendências:**

- ⚠️ **Transaction Support:** Atualmente faz Save sequencial sem rollback automático
- ⚠️ **CNPJ Validation:** Falta validação de CNPJ duplicado
- ⚠️ **Email Validation:** Falta validação de email duplicado

---

### 3. **Backend - Handler** ✅ COMPLETO

**Arquivo:** `backend/internal/infrastructure/http/handler/auth_handler.go`

**Implementado:**

- ✅ Route: `POST /auth/signup`
- ✅ Handler: `handleSignup()`
- ✅ Parse de `SignupInput`
- ✅ Chamada ao `SignupUseCase`
- ✅ Retorno padronizado (StandardResponse)

**Status:** ✅ FUNCIONANDO

---

### 4. **Frontend - Signup Page** ✅ COMPLETO

**Arquivo:** `frontend/app/(auth)/signup/page.tsx`

**Implementado:**

- ✅ Formulário com validação (React Hook Form + Yup)
- ✅ Campos: barberName, cnpj, name, email, password
- ✅ Error handling e loading states
- ✅ Redirecionamento via `AuthContext.signup()`
- ✅ Design System aplicado (tokens, MUI)

**Status:** ✅ FUNCIONANDO

---

### 5. **Frontend - Onboarding Page** ✅ COMPLETO

**Arquivo:** `frontend/app/onboarding/page.tsx`

**Implementado:**

- ✅ Página de boas-vindas com checklist
- ✅ Botão "Começar a Usar"
- ✅ Chamada à API: `POST /tenants/onboarding/complete`
- ✅ Refetch do user após completar
- ✅ Redirecionamento para `/dashboard`

**Status:** ✅ FUNCIONANDO

---

### 6. **Frontend - AuthContext** ✅ COMPLETO

**Arquivo:** `frontend/app/lib/contexts/AuthContext.tsx`

**Implementado:**

- ✅ Método `signup()` com chamada à API
- ✅ Auto-login após signup (salva tokens)
- ✅ Redirecionamento para `/onboarding` após signup
- ✅ Método `refetchUser()` para atualizar estado local

**Status:** ✅ FUNCIONANDO

---

## ❌ O Que Falta Implementar

### 1. **Backend - Complete Onboarding Use Case** 🔴 CRÍTICO

**Arquivo a criar:** `backend/internal/application/usecase/tenant/complete_onboarding_usecase.go`

**Responsabilidade:**

- Receber `tenantID` do contexto (JWT)
- Buscar tenant no repositório
- Atualizar `OnboardingCompleted = true`
- Persistir no banco
- Retornar sucesso

**Estrutura proposta:**

```go
package tenant

import (
    "context"
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
    // 1. Buscar tenant
    tenant, err := uc.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return fmt.Errorf("tenant not found: %w", err)
    }

    // 2. Validar se já completou (idempotência)
    if tenant.OnboardingCompleted {
        return nil // Já completado, retornar sucesso
    }

    // 3. Marcar como completado
    tenant.OnboardingCompleted = true
    tenant.AtualizadoEm = time.Now()

    // 4. Persistir
    if err := uc.tenantRepo.Update(ctx, tenant); err != nil {
        return fmt.Errorf("failed to update tenant: %w", err)
    }

    return nil
}
```

**Status:** 🔴 NÃO IMPLEMENTADO

---

### 2. **Backend - Tenant Handler** 🔴 CRÍTICO

**Arquivo a criar/modificar:** `backend/internal/infrastructure/http/handler/tenant_handler.go`

**Responsabilidade:**

- Registrar route: `POST /api/v1/tenants/onboarding/complete`
- Middleware: `AuthMiddleware` + `TenantMiddleware`
- Extrair `tenantID` do contexto
- Chamar `CompleteOnboardingUseCase`
- Retornar StandardResponse

**Estrutura proposta:**

```go
package handler

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/andviana23/barber-analytics-backend-v2/internal/application/usecase/tenant"
    "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/http/response"
)

type TenantHandler struct {
    completeOnboardingUC *tenant.CompleteOnboardingUseCase
}

func NewTenantHandler(
    completeOnboardingUC *tenant.CompleteOnboardingUseCase,
) *TenantHandler {
    return &TenantHandler{
        completeOnboardingUC: completeOnboardingUC,
    }
}

func (h *TenantHandler) RegisterRoutes(r chi.Router) {
    r.Route("/tenants", func(r chi.Router) {
        // Rotas protegidas
        r.Group(func(r chi.Router) {
            r.Use(httpMiddleware.ChiAuthMiddleware(h.jwtService))
            r.Use(httpMiddleware.ChiTenantMiddleware())

            r.Post("/onboarding/complete", h.handleCompleteOnboarding)
        })
    })
}

func (h *TenantHandler) handleCompleteOnboarding(w http.ResponseWriter, r *http.Request) {
    // Extrair tenant ID do contexto
    tenantID, err := getTenantIDFromRequest(r)
    if err != nil {
        writeStandardResponse(w, response.Error("FORBIDDEN", "Tenant ID not found", err.Error(), ""))
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

**Status:** 🔴 NÃO IMPLEMENTADO

---

### 3. **Backend - Wire Dependency Injection** 🔴 CRÍTICO

**Arquivo a modificar:** `backend/cmd/api/main.go` (ou arquivo de DI)

**Ações:**

1. Criar `CompleteOnboardingUseCase` no container
2. Criar `TenantHandler` com dependência do use case
3. Registrar routes do `TenantHandler` no router principal

**Exemplo:**

```go
// main.go ou di.go
completeOnboardingUC := tenant.NewCompleteOnboardingUseCase(tenantRepo)
tenantHandler := handler.NewTenantHandler(completeOnboardingUC)

// Registrar routes
tenantHandler.RegisterRoutes(router)
```

**Status:** 🔴 NÃO IMPLEMENTADO

---

### 4. **Frontend - Middleware Enhancement** 🟡 OPCIONAL

**Arquivo:** `frontend/middleware.ts`

**Objetivo:** Verificar `onboardingCompleted` e redirecionar para `/onboarding` se necessário.

**Desafio:** Como obter status de onboarding?

**Opções:**

#### **Opção A: Incluir no JWT Claims** ⭐ RECOMENDADO

```go
// backend/internal/domain/service/jwt_service.go
func (s *JWTService) GenerateAccessToken(userID, tenantID, email, role string, onboardingCompleted bool) (string, error) {
    claims := jwt.MapClaims{
        "user_id":              userID,
        "tenant_id":            tenantID,
        "email":                email,
        "role":                 role,
        "onboarding_completed": onboardingCompleted, // ✅ ADICIONAR
        "exp":                  time.Now().Add(15 * time.Minute).Unix(),
        "iat":                  time.Now().Unix(),
    }
    // ...
}
```

**Prós:**

- ✅ Sem requisição extra ao backend
- ✅ Middleware Next.js pode ler do cookie
- ✅ Performance máxima

**Contras:**

- ⚠️ Precisa atualizar JWT após completar onboarding (fazer logout/login ou refresh)

#### **Opção B: Fetch Separado no Middleware**

```typescript
// frontend/middleware.ts
export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("bap.access_token")?.value;

  if (accessToken && !isPublicRoute(pathname)) {
    // Fetch tenant info
    const tenantInfo = await fetch(`${API_URL}/tenants/me`, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    const data = await tenantInfo.json();

    if (!data.onboarding_completed && pathname !== "/onboarding") {
      return NextResponse.redirect(new URL("/onboarding", request.url));
    }
  }

  return NextResponse.next();
}
```

**Prós:**

- ✅ Sempre atualizado (não depende de token)

**Contras:**

- ❌ Requisição extra a cada navegação (impacto em performance)
- ❌ Pode causar delay no carregamento

**Status:** 🟡 DECISÃO PENDENTE (recomendo Opção A)

---

### 5. **Testes Automatizados** 🟡 RECOMENDADO

#### **Backend Unit Tests**

**Arquivo:** `backend/internal/application/usecase/tenant/complete_onboarding_usecase_test.go`

```go
func TestCompleteOnboardingUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := &mockTenantRepository{}
    uc := NewCompleteOnboardingUseCase(mockRepo)

    // Act
    err := uc.Execute(context.Background(), "tenant-123")

    // Assert
    assert.NoError(t, err)
    assert.True(t, mockRepo.savedTenant.OnboardingCompleted)
}
```

#### **Backend Integration Tests**

**Arquivo:** `backend/tests/integration/onboarding_flow_test.go`

**Fluxo completo:**

1. POST /auth/signup → Verificar tokens retornados
2. GET /auth/me → Verificar `onboarding_completed = false`
3. POST /tenants/onboarding/complete → Sucesso
4. GET /auth/me → Verificar `onboarding_completed = true`

#### **Frontend E2E Tests**

**Arquivo:** `frontend/e2e/onboarding.spec.ts`

```typescript
test("should complete full signup and onboarding flow", async ({ page }) => {
  // 1. Signup
  await page.goto("/signup");
  await page.fill('[data-testid="barber-name-input"]', "Barbearia Teste E2E");
  await page.fill('[data-testid="cnpj-input"]', "12345678000199");
  await page.fill('[data-testid="name-input"]', "João Silva");
  await page.fill('[data-testid="email-input"]', "joao@teste.com");
  await page.fill('[data-testid="password-input"]', "senha123");
  await page.click('[data-testid="signup-button"]');

  // 2. Deve redirecionar para onboarding
  await expect(page).toHaveURL("/onboarding");

  // 3. Completar onboarding
  await page.click("text=Começar a Usar");

  // 4. Deve redirecionar para dashboard
  await expect(page).toHaveURL("/dashboard");
});
```

**Status:** 🟡 NÃO IMPLEMENTADO

---

## 🚨 Issues Identificados

### 1. **Transaction Support** ⚠️ CRÍTICO

**Problema:** `SignupUseCase` faz `Save` sequencial sem rollback:

```go
// ❌ Se Save do User falhar, Tenant fica órfão no banco
if err := uc.tenantRepo.Save(ctx, tenant); err != nil {
    return nil, err
}
// ...
if err := uc.userRepo.Save(ctx, tenant.ID, user); err != nil {
    return nil, err // Tenant já foi salvo!
}
```

**Solução:** Implementar Transaction wrapper no repositório:

```go
type TxManager interface {
    WithTx(ctx context.Context, fn func(context.Context) error) error
}

func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    var tenant *entity.Tenant
    var user *entity.User

    err := uc.txManager.WithTx(ctx, func(txCtx context.Context) error {
        // 1. Create tenant
        tenant = entity.NewTenant(...)
        if err := uc.tenantRepo.Save(txCtx, tenant); err != nil {
            return err
        }

        // 2. Create user
        user = entity.NewUser(...)
        if err := uc.userRepo.Save(txCtx, tenant.ID, user); err != nil {
            return err // Rollback automático
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Gerar token após commit bem-sucedido
    // ...
}
```

**Status:** ⚠️ NÃO IMPLEMENTADO (criar em fase posterior)

---

### 2. **Validação de Duplicados** ⚠️ IMPORTANTE

**Problema:** `SignupUseCase` não valida duplicados de CNPJ/Email.

**Solução:**

```go
func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    // 1. Validar CNPJ duplicado
    if input.CNPJ != "" {
        existing, _ := uc.tenantRepo.FindByCNPJ(ctx, input.CNPJ)
        if existing != nil {
            return nil, errors.New("CNPJ already registered")
        }
    }

    // 2. Validar Email duplicado
    existingUser, _ := uc.userRepo.FindByEmailAny(ctx, input.Email)
    if existingUser != nil {
        return nil, errors.New("Email already registered")
    }

    // ... resto do código
}
```

**Status:** ⚠️ NÃO IMPLEMENTADO

---

### 3. **JWT Claims - Onboarding Status** 🟡 OPCIONAL

**Problema:** Frontend não tem acesso ao status de onboarding sem fetch extra.

**Solução:** Incluir no JWT (ver Opção A acima).

**Status:** 🟡 DECISÃO PENDENTE

---

## 📋 Plano de Continuação

### **Fase 1: Backend - Complete Onboarding** (Prioritário)

**Tempo estimado:** 1-2 horas

**Tarefas:**

1. ✅ Criar `CompleteOnboardingUseCase` em `backend/internal/application/usecase/tenant/`
2. ✅ Criar `TenantHandler` em `backend/internal/infrastructure/http/handler/`
3. ✅ Registrar routes em `main.go` ou DI container
4. ✅ Testar endpoint com curl:
   ```bash
   curl -X POST http://localhost:8080/api/v1/tenants/onboarding/complete \
     -H "Authorization: Bearer {token}"
   ```

---

### **Fase 2: Backend - Validações** (Importante)

**Tempo estimado:** 1 hora

**Tarefas:**

1. ✅ Adicionar validação de CNPJ duplicado em `SignupUseCase`
2. ✅ Adicionar validação de Email duplicado em `SignupUseCase`
3. ✅ Retornar erros HTTP apropriados (409 Conflict)

---

### **Fase 3: Testes** (Recomendado)

**Tempo estimado:** 2-3 horas

**Tarefas:**

1. ✅ Unit test: `CompleteOnboardingUseCase`
2. ✅ Unit test: `SignupUseCase` (duplicados)
3. ✅ Integration test: Fluxo completo signup → onboarding → dashboard
4. ✅ E2E test: Playwright flow completo

---

### **Fase 4: Frontend Middleware** (Opcional)

**Tempo estimado:** 1 hora

**Decisão:** Implementar Opção A (JWT claims) ou Opção B (fetch separado)?

**Tarefas se Opção A:**

1. ✅ Modificar `JWTService.GenerateAccessToken()` para incluir `onboarding_completed`
2. ✅ Atualizar `LoginUseCase` para passar valor correto
3. ✅ Frontend: Decodificar JWT no middleware e verificar claim
4. ✅ Redirecionar para `/onboarding` se `false`

---

### **Fase 5: Transaction Support** (Futuro)

**Tempo estimado:** 3-4 horas

**Tarefas:**

1. ✅ Criar interface `TxManager`
2. ✅ Implementar `PostgresTxManager` com `sql.Tx`
3. ✅ Refatorar `SignupUseCase` para usar transactions
4. ✅ Testar rollback em caso de erro

---

## 🎯 Recomendação Imediata

**Começar por Fase 1 (Backend - Complete Onboarding)**, pois:

1. ✅ É o único bloqueador crítico para fluxo end-to-end funcionar
2. ✅ Frontend já está pronto e aguardando endpoint
3. ✅ Migration já foi aplicada no banco
4. ✅ Implementação é simples e direta (1-2 horas)

**Próximos comandos:**

```bash
# 1. Criar arquivo do use case
touch backend/internal/application/usecase/tenant/complete_onboarding_usecase.go

# 2. Criar arquivo do handler
touch backend/internal/infrastructure/http/handler/tenant_handler.go

# 3. Implementar (código fornecido acima)

# 4. Registrar no main.go

# 5. Testar
make run-backend
curl -X POST http://localhost:8080/api/v1/tenants/onboarding/complete \
  -H "Authorization: Bearer {token_do_signup}"
```

---

## ✅ Checklist de Implementação

### Backend

- [x] Migration 024 aplicada
- [x] Entity Tenant com OnboardingCompleted
- [x] Repository atualizado
- [x] SignupUseCase implementado
- [x] AuthHandler com /signup
- [ ] **CompleteOnboardingUseCase** ← PRÓXIMO
- [ ] **TenantHandler** ← PRÓXIMO
- [ ] **Routes registradas** ← PRÓXIMO
- [ ] Validação de duplicados (CNPJ/Email)
- [ ] Unit tests
- [ ] Integration tests

### Frontend

- [x] Signup page
- [x] Onboarding page
- [x] AuthContext.signup()
- [x] API client configurado
- [ ] Middleware com onboarding check (opcional)
- [ ] E2E tests

---

**Autor:** AI Assistant
**Status:** 📝 Documento de Planejamento
**Próximo Passo:** Implementar Fase 1 (CompleteOnboardingUseCase + TenantHandler)
