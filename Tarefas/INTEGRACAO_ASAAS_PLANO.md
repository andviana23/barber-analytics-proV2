# 📋 Plano de Implementação - Integração API Asaas

**Versão:** 1.0
**Data:** 17/11/2025
**Status:** 📝 Planejamento
**Responsável:** Equipe Backend
**Estimativa Total:** 27-39 horas (5-7 dias)

---

## 📊 Índice

1. [Contexto e Objetivos](#contexto-e-objetivos)
2. [Informações a Capturar](#informações-a-capturar)
3. [Fases de Implementação](#fases-de-implementação)
4. [Checklist Completo](#checklist-completo)
5. [Estimativas de Tempo](#estimativas-de-tempo)
6. [Próximos Passos](#próximos-passos)

---

## 🎯 Contexto e Objetivos

### Status Atual do Sistema

**✅ Implementado:**
- Backend Go com Clean Architecture (72% completo)
- Sistema de assinaturas **MANUAL** implementado
- Entidades de domínio prontas (`Assinatura`, `AssinaturaInvoice`)
- Documentação completa em `/docs/INTEGRACOES_ASAAS.md`
- 4 Cron Jobs rodando (02:00, 03:00, 04:00, 08:00)

**⚠️ Não Implementado:**
- Integração com API Asaas
- Cliente HTTP para Asaas
- Sincronização automática de faturas
- Webhooks Asaas

### Objetivo da Integração

Integrar o sistema com a **API do Asaas** para:

1. ✅ Criar assinaturas automaticamente no Asaas
2. ✅ Sincronizar status de assinaturas (ATIVA, ATRASADA, AGUARDANDO)
3. ✅ Sincronizar faturas e pagamentos
4. ✅ Capturar informações de clientes
5. ✅ Automatizar repasse de comissões

---

## 📋 Informações a Capturar do Asaas

### 1. Status de Pagamento da Assinatura
- **Campo Asaas:** `payment.status`
- **Valores:** `PENDING`, `CONFIRMED`, `RECEIVED`
- **Mapeamento:**
  - `CONFIRMED` → Pagamento confirmado ✅
  - `RECEIVED` → Dinheiro recebido na conta ✅
  - `PENDING` → Aguardando pagamento ⏳

### 2. Status da Assinatura
- **Campo Asaas:** `subscription.status`
- **Valores:** `ACTIVE`, `OVERDUE`, `AWAITING_PAYMENT`, `CANCELLED`, `EXPIRED`
- **Mapeamento:**
  - `ACTIVE` → Assinatura ativa ✅
  - `OVERDUE` → Assinatura atrasada ⚠️
  - `AWAITING_PAYMENT` → Aguardando pagamento ⏳

### 3. Data de Pagamento da Assinatura
- **Campo Asaas:** `invoice.paymentDate`
- **Formato:** `YYYY-MM-DD` ou `YYYY-MM-DDTHH:MM:SS`
- **Armazenar em:** `assinatura_invoices.data_pagamento`

### 4. Data de Vencimento da Assinatura
- **Campo Asaas:** `invoice.dueDate` e `invoice.originalDueDate`
- **Armazenar em:** `assinatura_invoices.data_vencimento`

### 5. Data de Previsão de Recebimento
- **Campo Asaas:** `invoice.estimatedPaymentDate`
- **Armazenar em:** `assinatura_invoices.data_previsao_recebimento` (novo campo)

### 6. Nome do Cliente
- **Campo Asaas:** `customer.name`
- **Buscar via:** `GET /customers/{customer_id}`
- **Exibir em:** Frontend (lista de assinaturas)

---

## 🚀 Fases de Implementação

---

### **FASE 1: Infraestrutura Base (4-6 horas)**

#### 📦 Tarefa 1.1 - Criar Cliente HTTP Asaas

**Arquivo:** `backend/internal/infrastructure/external/asaas/client.go`

**Estrutura:**
```go
package asaas

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "go.uber.org/zap"
)

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
    logger     *zap.Logger
}

func NewClient(apiKey, baseURL string, logger *zap.Logger) *Client {
    return &Client{
        apiKey:  apiKey,
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        logger: logger,
    }
}
```

**Métodos a implementar:**

1. **CreateSubscription** - Criar assinatura no Asaas
```go
func (c *Client) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*SubscriptionResponse, error)
```

2. **GetSubscription** - Buscar assinatura por ID
```go
func (c *Client) GetSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResponse, error)
```

3. **ListInvoices** - Listar faturas de uma assinatura
```go
func (c *Client) ListInvoices(ctx context.Context, subscriptionID string) (*InvoiceListResponse, error)
```

4. **GetInvoice** - Buscar detalhes de uma fatura específica
```go
func (c *Client) GetInvoice(ctx context.Context, invoiceID string) (*InvoiceResponse, error)
```

5. **CancelSubscription** - Cancelar assinatura no Asaas
```go
func (c *Client) CancelSubscription(ctx context.Context, subscriptionID string) error
```

6. **GetCustomer** - Buscar informações do cliente
```go
func (c *Client) GetCustomer(ctx context.Context, customerID string) (*CustomerResponse, error)
```

**Features obrigatórias:**
- ✅ Autenticação via `Authorization: Bearer {apiKey}`
- ✅ Retry logic com exponential backoff
- ✅ Timeout configurável (30s)
- ✅ Error handling robusto
- ✅ Logging estruturado (request_id, method, status_code)

**Exemplo de retry logic:**
```go
func (c *Client) doRequestWithRetry(ctx context.Context, req *http.Request, maxRetries int) (*http.Response, error) {
    var lastErr error
    backoff := 1 * time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := c.httpClient.Do(req)

        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        if err != nil {
            lastErr = err
        } else {
            lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
            resp.Body.Close()
        }

        if attempt < maxRetries-1 {
            c.logger.Warn("retry request",
                zap.Int("attempt", attempt+1),
                zap.Duration("backoff", backoff),
                zap.Error(lastErr),
            )
            time.Sleep(backoff)
            backoff *= 2 // Exponential backoff
        }
    }

    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

---

#### 📦 Tarefa 1.2 - Criar DTOs da API Asaas

**Arquivo:** `backend/internal/infrastructure/external/asaas/dto.go`

```go
package asaas

import "time"

// ========== REQUEST DTOs ==========

type CreateSubscriptionRequest struct {
    Customer    string  `json:"customer"`              // ID do cliente no Asaas
    BillingType string  `json:"billingType"`           // CREDIT_CARD, BOLETO, PIX
    Value       float64 `json:"value"`                 // Valor da assinatura
    NextDueDate string  `json:"nextDueDate"`           // YYYY-MM-DD
    Cycle       string  `json:"cycle"`                 // MONTHLY, QUARTERLY, YEARLY
    Description string  `json:"description,omitempty"` // Descrição opcional
}

// ========== RESPONSE DTOs ==========

type SubscriptionResponse struct {
    Object      string  `json:"object"`
    ID          string  `json:"id"`
    Status      string  `json:"status"` // ACTIVE, OVERDUE, AWAITING_PAYMENT, CANCELLED, EXPIRED
    Customer    string  `json:"customer"`
    Value       float64 `json:"value"`
    NextDueDate string  `json:"nextDueDate"`
    Cycle       string  `json:"cycle"`
    DateCreated string  `json:"dateCreated"`
}

type InvoiceResponse struct {
    Object               string   `json:"object"`
    ID                   string   `json:"id"`
    Subscription         string   `json:"subscription"`
    Status               string   `json:"status"` // PENDING, CONFIRMED, RECEIVED
    Value                float64  `json:"value"`
    DueDate              string   `json:"dueDate"`
    OriginalDueDate      string   `json:"originalDueDate"`
    PaymentDate          *string  `json:"paymentDate"`
    ClientPaymentDate    *string  `json:"clientPaymentDate"`
    EstimatedPaymentDate *string  `json:"estimatedPaymentDate"`
    InvoiceURL           string   `json:"invoiceUrl"`
    DateCreated          string   `json:"dateCreated"`
}

type InvoiceListResponse struct {
    Object  string            `json:"object"`
    HasMore bool              `json:"hasMore"`
    Data    []InvoiceResponse `json:"data"`
}

type CustomerResponse struct {
    Object  string `json:"object"`
    ID      string `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    CpfCnpj string `json:"cpfCnpj"`
}

type ErrorResponse struct {
    Errors []struct {
        Code        string `json:"code"`
        Description string `json:"description"`
    } `json:"errors"`
}
```

---

#### 📦 Tarefa 1.3 - Configuração de Ambiente

**Arquivo:** `backend/.env`

Adicionar variáveis:
```bash
# ========== Asaas Integration ==========
ASAAS_API_KEY=sk_live_xxxxxxxxxxxxx
ASAAS_BASE_URL=https://api.asaas.com/v3
ASAAS_REQUEST_TIMEOUT=30s
FEATURE_ASAAS_INTEGRATION=true
```

**Arquivo:** `backend/internal/config/config.go`

Adicionar campos:
```go
type Config struct {
    // ... campos existentes ...

    // Asaas Integration
    AsaasAPIKey             string
    AsaasBaseURL            string
    AsaasTimeout            time.Duration
    AsaasIntegrationEnabled bool
}

func LoadConfig() (*Config, error) {
    // ... código existente ...

    cfg.AsaasAPIKey = os.Getenv("ASAAS_API_KEY")
    cfg.AsaasBaseURL = getEnvOrDefault("ASAAS_BASE_URL", "https://api.asaas.com/v3")

    timeoutStr := getEnvOrDefault("ASAAS_REQUEST_TIMEOUT", "30s")
    cfg.AsaasTimeout, _ = time.ParseDuration(timeoutStr)

    cfg.AsaasIntegrationEnabled = getEnvOrDefault("FEATURE_ASAAS_INTEGRATION", "false") == "true"

    return cfg, nil
}
```

---

### **FASE 2: Atualização de Entidades (2-3 horas)**

#### 📦 Tarefa 2.1 - Adicionar Campos em Assinatura

**Arquivo:** `backend/internal/domain/entity/assinatura.go`

Adicionar campos:
```go
type Assinatura struct {
    // ... campos existentes ...

    // Campos Asaas
    asaasStatus           *string    // ACTIVE, OVERDUE, AWAITING_PAYMENT
    ultimaSincronizacao   *time.Time // Última sincronização com Asaas
}
```

Adicionar getters:
```go
func (a *Assinatura) AsaasStatus() *string     { return a.asaasStatus }
func (a *Assinatura) UltimaSincronizacao() *time.Time { return a.ultimaSincronizacao }
```

Adicionar método:
```go
// AtualizarStatusAsaas atualiza o status sincronizado do Asaas
func (a *Assinatura) AtualizarStatusAsaas(asaasStatus string) {
    a.asaasStatus = &asaasStatus
    now := time.Now()
    a.ultimaSincronizacao = &now
    a.updatedAt = now
}
```

Atualizar `ReconstructAssinatura`:
```go
func ReconstructAssinatura(
    // ... parâmetros existentes ...
    asaasStatus *string,
    ultimaSincronizacao *time.Time,
) *Assinatura {
    return &Assinatura{
        // ... campos existentes ...
        asaasStatus:         asaasStatus,
        ultimaSincronizacao: ultimaSincronizacao,
    }
}
```

---

#### 📦 Tarefa 2.2 - Adicionar Campos em AssinaturaInvoice

**Arquivo:** `backend/internal/domain/entity/assinatura_invoice.go`

Adicionar campos:
```go
type AssinaturaInvoice struct {
    // ... campos existentes ...

    // Campos Asaas
    asaasStatus             *string    // PENDING, CONFIRMED, RECEIVED
    dataPrevisaoRecebimento *time.Time // estimatedPaymentDate
    clientPaymentDate       *time.Time // Data que cliente pagou
    invoiceURL              string     // URL da fatura no Asaas
    processada              bool       // Se já foi processada para repasse
}
```

Adicionar getters:
```go
func (i *AssinaturaInvoice) AsaasStatus() *string              { return i.asaasStatus }
func (i *AssinaturaInvoice) DataPrevisaoRecebimento() *time.Time { return i.dataPrevisaoRecebimento }
func (i *AssinaturaInvoice) ClientPaymentDate() *time.Time     { return i.clientPaymentDate }
func (i *AssinaturaInvoice) InvoiceURL() string                { return i.invoiceURL }
func (i *AssinaturaInvoice) Processada() bool                  { return i.processada }
```

Adicionar métodos:
```go
// AtualizarDadosAsaas atualiza informações sincronizadas do Asaas
func (i *AssinaturaInvoice) AtualizarDadosAsaas(
    asaasStatus string,
    invoiceURL string,
    estimatedPaymentDate *time.Time,
    clientPaymentDate *time.Time,
) {
    i.asaasStatus = &asaasStatus
    i.invoiceURL = invoiceURL
    i.dataPrevisaoRecebimento = estimatedPaymentDate
    i.clientPaymentDate = clientPaymentDate
    i.updatedAt = time.Now()
}

// MarcarComoProcessada marca a invoice como processada para repasse
func (i *AssinaturaInvoice) MarcarComoProcessada() error {
    if i.status != InvoicePago {
        return errors.New("apenas invoices pagas podem ser marcadas como processadas")
    }
    if i.processada {
        return errors.New("invoice já foi processada")
    }

    i.processada = true
    i.updatedAt = time.Now()
    return nil
}
```

Atualizar `ReconstructAssinaturaInvoice`:
```go
func ReconstructAssinaturaInvoice(
    // ... parâmetros existentes ...
    asaasStatus *string,
    dataPrevisaoRecebimento *time.Time,
    clientPaymentDate *time.Time,
    invoiceURL string,
    processada bool,
) *AssinaturaInvoice {
    return &AssinaturaInvoice{
        // ... campos existentes ...
        asaasStatus:             asaasStatus,
        dataPrevisaoRecebimento: dataPrevisaoRecebimento,
        clientPaymentDate:       clientPaymentDate,
        invoiceURL:              invoiceURL,
        processada:              processada,
    }
}
```

---

### **FASE 3: Casos de Uso (6-8 horas)**

#### 📦 Tarefa 3.1 - Criar Assinatura com Asaas

**Arquivo:** `backend/internal/application/usecase/subscription/create_assinatura_asaas_usecase.go`

```go
package subscription

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "go.uber.org/zap"

    "backend/internal/domain/entity"
    "backend/internal/domain/repository"
    "backend/internal/infrastructure/external/asaas"
)

type CreateAssinaturaAsaasInput struct {
    TenantID    uuid.UUID `json:"tenant_id" validate:"required"`
    PlanID      uuid.UUID `json:"plan_id" validate:"required"`
    BarbeiroID  uuid.UUID `json:"barbeiro_id" validate:"required"`
    CustomerID  string    `json:"customer_id" validate:"required"` // ID do cliente no Asaas
    BillingType string    `json:"billing_type" validate:"required"` // CREDIT_CARD, BOLETO, PIX
    DataInicio  time.Time `json:"data_inicio" validate:"required"`
}

type CreateAssinaturaAsaasOutput struct {
    ID                  uuid.UUID `json:"id"`
    AsaasSubscriptionID string    `json:"asaas_subscription_id"`
    Status              string    `json:"status"`
    AsaasStatus         string    `json:"asaas_status"`
    ProximaFaturaData   time.Time `json:"proxima_fatura_data"`
}

type CreateAssinaturaAsaasUseCase struct {
    assinaturaRepo repository.AssinaturaRepository
    planoRepo      repository.PlanoAssinaturaRepository
    asaasClient    *asaas.Client
    logger         *zap.Logger
}

func NewCreateAssinaturaAsaasUseCase(
    assinaturaRepo repository.AssinaturaRepository,
    planoRepo repository.PlanoAssinaturaRepository,
    asaasClient *asaas.Client,
    logger *zap.Logger,
) *CreateAssinaturaAsaasUseCase {
    return &CreateAssinaturaAsaasUseCase{
        assinaturaRepo: assinaturaRepo,
        planoRepo:      planoRepo,
        asaasClient:    asaasClient,
        logger:         logger,
    }
}

func (uc *CreateAssinaturaAsaasUseCase) Execute(
    ctx context.Context,
    input CreateAssinaturaAsaasInput,
) (*CreateAssinaturaAsaasOutput, error) {

    // 1. Buscar plano
    plano, err := uc.planoRepo.FindByID(ctx, input.TenantID, input.PlanID)
    if err != nil {
        return nil, fmt.Errorf("plano não encontrado: %w", err)
    }

    // 2. Criar assinatura no Asaas
    asaasReq := &asaas.CreateSubscriptionRequest{
        Customer:    input.CustomerID,
        BillingType: input.BillingType,
        Value:       plano.Valor().InexactFloat64(),
        NextDueDate: input.DataInicio.Format("2006-01-02"),
        Cycle:       mapPeriodicidadeToCycle(plano.Periodicidade()),
        Description: fmt.Sprintf("Plano %s - Barbeiro", plano.Nome()),
    }

    asaasResp, err := uc.asaasClient.CreateSubscription(ctx, asaasReq)
    if err != nil {
        uc.logger.Error("falha ao criar assinatura no Asaas",
            zap.String("tenant_id", input.TenantID.String()),
            zap.Error(err),
        )
        return nil, fmt.Errorf("erro ao criar assinatura no Asaas: %w", err)
    }

    // 3. Criar assinatura no banco
    assinatura, err := entity.NewAssinaturaComAsaas(
        input.TenantID,
        input.PlanID,
        input.BarbeiroID,
        asaasResp.ID,
        input.DataInicio,
        input.DataInicio, // proxima_fatura_data = data_inicio
    )
    if err != nil {
        // Tentar cancelar no Asaas (rollback)
        _ = uc.asaasClient.CancelSubscription(ctx, asaasResp.ID)
        return nil, fmt.Errorf("erro ao criar entidade assinatura: %w", err)
    }

    // Atualizar status Asaas
    assinatura.AtualizarStatusAsaas(asaasResp.Status)

    // 4. Persistir no banco
    if err := uc.assinaturaRepo.Create(ctx, input.TenantID, assinatura); err != nil {
        // Tentar cancelar no Asaas (rollback)
        _ = uc.asaasClient.CancelSubscription(ctx, asaasResp.ID)
        return nil, fmt.Errorf("erro ao salvar assinatura: %w", err)
    }

    uc.logger.Info("assinatura criada com sucesso",
        zap.String("assinatura_id", assinatura.ID().String()),
        zap.String("asaas_subscription_id", asaasResp.ID),
    )

    return &CreateAssinaturaAsaasOutput{
        ID:                  assinatura.ID(),
        AsaasSubscriptionID: asaasResp.ID,
        Status:              string(assinatura.Status()),
        AsaasStatus:         asaasResp.Status,
        ProximaFaturaData:   assinatura.ProximaFaturaData(),
    }, nil
}

func mapPeriodicidadeToCycle(periodicidade entity.Periodicidade) string {
    switch periodicidade {
    case entity.PeriodMensal:
        return "MONTHLY"
    case entity.PeriodTrimestral:
        return "QUARTERLY"
    case entity.PeriodAnual:
        return "YEARLY"
    default:
        return "MONTHLY"
    }
}
```

---

#### 📦 Tarefa 3.2 - Sincronizar Assinaturas do Asaas

**Arquivo:** `backend/internal/application/usecase/subscription/sync_assinaturas_asaas_usecase.go`

```go
package subscription

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "go.uber.org/zap"

    "backend/internal/domain/repository"
    "backend/internal/infrastructure/external/asaas"
)

type SyncAssinaturasAsaasUseCase struct {
    assinaturaRepo repository.AssinaturaRepository
    asaasClient    *asaas.Client
    logger         *zap.Logger
}

func NewSyncAssinaturasAsaasUseCase(
    assinaturaRepo repository.AssinaturaRepository,
    asaasClient *asaas.Client,
    logger *zap.Logger,
) *SyncAssinaturasAsaasUseCase {
    return &SyncAssinaturasAsaasUseCase{
        assinaturaRepo: assinaturaRepo,
        asaasClient:    asaasClient,
        logger:         logger,
    }
}

func (uc *SyncAssinaturasAsaasUseCase) Execute(ctx context.Context, tenantID uuid.UUID) error {
    uc.logger.Info("iniciando sincronização de assinaturas", zap.String("tenant_id", tenantID.String()))

    // 1. Buscar todas assinaturas do tenant que têm asaas_subscription_id
    assinaturas, err := uc.assinaturaRepo.FindByTenant(ctx, tenantID)
    if err != nil {
        return fmt.Errorf("erro ao buscar assinaturas: %w", err)
    }

    syncCount := 0
    errorCount := 0

    // 2. Para cada assinatura com Asaas
    for _, assinatura := range assinaturas {
        if assinatura.AsaasSubscriptionID() == nil {
            continue // Pular assinaturas manuais
        }

        asaasSubID := *assinatura.AsaasSubscriptionID()

        // 3. Buscar dados no Asaas
        asaasResp, err := uc.asaasClient.GetSubscription(ctx, asaasSubID)
        if err != nil {
            uc.logger.Error("erro ao buscar assinatura no Asaas",
                zap.String("assinatura_id", assinatura.ID().String()),
                zap.String("asaas_subscription_id", asaasSubID),
                zap.Error(err),
            )
            errorCount++
            continue
        }

        // 4. Atualizar status e data
        assinatura.AtualizarStatusAsaas(asaasResp.Status)

        if nextDueDate, err := time.Parse("2006-01-02", asaasResp.NextDueDate); err == nil {
            _ = assinatura.AtualizarProximaFatura(nextDueDate)
        }

        // 5. Persistir atualização
        if err := uc.assinaturaRepo.Update(ctx, tenantID, assinatura); err != nil {
            uc.logger.Error("erro ao atualizar assinatura",
                zap.String("assinatura_id", assinatura.ID().String()),
                zap.Error(err),
            )
            errorCount++
            continue
        }

        syncCount++
    }

    uc.logger.Info("sincronização de assinaturas concluída",
        zap.String("tenant_id", tenantID.String()),
        zap.Int("synced", syncCount),
        zap.Int("errors", errorCount),
    )

    return nil
}
```

---

#### 📦 Tarefa 3.3 - Sincronizar Faturas do Asaas

**Arquivo:** `backend/internal/application/usecase/subscription/sync_invoices_asaas_usecase.go`

```go
package subscription

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "go.uber.org/zap"

    "backend/internal/domain/entity"
    "backend/internal/domain/repository"
    "backend/internal/infrastructure/external/asaas"
)

type SyncInvoicesAsaasUseCase struct {
    assinaturaRepo repository.AssinaturaRepository
    invoiceRepo    repository.AssinaturaInvoiceRepository
    asaasClient    *asaas.Client
    logger         *zap.Logger
}

func NewSyncInvoicesAsaasUseCase(
    assinaturaRepo repository.AssinaturaRepository,
    invoiceRepo repository.AssinaturaInvoiceRepository,
    asaasClient *asaas.Client,
    logger *zap.Logger,
) *SyncInvoicesAsaasUseCase {
    return &SyncInvoicesAsaasUseCase{
        assinaturaRepo: assinaturaRepo,
        invoiceRepo:    invoiceRepo,
        asaasClient:    asaasClient,
        logger:         logger,
    }
}

func (uc *SyncInvoicesAsaasUseCase) Execute(ctx context.Context, tenantID uuid.UUID) error {
    uc.logger.Info("iniciando sincronização de faturas", zap.String("tenant_id", tenantID.String()))

    // 1. Buscar assinaturas com Asaas
    assinaturas, err := uc.assinaturaRepo.FindByTenant(ctx, tenantID)
    if err != nil {
        return fmt.Errorf("erro ao buscar assinaturas: %w", err)
    }

    syncCount := 0
    createCount := 0
    updateCount := 0
    errorCount := 0

    // 2. Para cada assinatura
    for _, assinatura := range assinaturas {
        if assinatura.AsaasSubscriptionID() == nil {
            continue // Pular manuais
        }

        asaasSubID := *assinatura.AsaasSubscriptionID()

        // 3. Buscar faturas no Asaas
        invoicesResp, err := uc.asaasClient.ListInvoices(ctx, asaasSubID)
        if err != nil {
            uc.logger.Error("erro ao buscar faturas no Asaas",
                zap.String("assinatura_id", assinatura.ID().String()),
                zap.Error(err),
            )
            errorCount++
            continue
        }

        // 4. Processar cada fatura
        for _, asaasInvoice := range invoicesResp.Data {
            // Verificar se já existe no DB
            existing, err := uc.invoiceRepo.FindByAsaasInvoiceID(ctx, tenantID, asaasInvoice.ID)

            if err != nil || existing == nil {
                // Criar nova fatura
                if err := uc.createInvoiceFromAsaas(ctx, tenantID, assinatura.ID(), &asaasInvoice); err != nil {
                    uc.logger.Error("erro ao criar fatura",
                        zap.String("asaas_invoice_id", asaasInvoice.ID),
                        zap.Error(err),
                    )
                    errorCount++
                    continue
                }
                createCount++
            } else {
                // Atualizar fatura existente
                if err := uc.updateInvoiceFromAsaas(ctx, tenantID, existing, &asaasInvoice); err != nil {
                    uc.logger.Error("erro ao atualizar fatura",
                        zap.String("invoice_id", existing.ID().String()),
                        zap.Error(err),
                    )
                    errorCount++
                    continue
                }
                updateCount++
            }
            syncCount++
        }
    }

    uc.logger.Info("sincronização de faturas concluída",
        zap.String("tenant_id", tenantID.String()),
        zap.Int("total_synced", syncCount),
        zap.Int("created", createCount),
        zap.Int("updated", updateCount),
        zap.Int("errors", errorCount),
    )

    return nil
}

func (uc *SyncInvoicesAsaasUseCase) createInvoiceFromAsaas(
    ctx context.Context,
    tenantID, assinaturaID uuid.UUID,
    asaasInvoice *asaas.InvoiceResponse,
) error {
    // Parse datas
    dueDate, _ := time.Parse("2006-01-02", asaasInvoice.DueDate)

    // Competência = mês da fatura
    competenciaInicio := dueDate.AddDate(0, -1, 0)
    competenciaFim := dueDate.AddDate(0, 0, -1)

    // Criar invoice
    invoice, err := entity.NewAssinaturaInvoiceComAsaas(
        tenantID,
        assinaturaID,
        asaasInvoice.ID,
        decimal.NewFromFloat(asaasInvoice.Value),
        dueDate,
        competenciaInicio,
        competenciaFim,
    )
    if err != nil {
        return err
    }

    // Atualizar dados do Asaas
    estimatedDate := parseAsaasDate(asaasInvoice.EstimatedPaymentDate)
    clientDate := parseAsaasDate(asaasInvoice.ClientPaymentDate)

    invoice.AtualizarDadosAsaas(
        asaasInvoice.Status,
        asaasInvoice.InvoiceURL,
        estimatedDate,
        clientDate,
    )

    // Mapear status
    invoice = mapAsaasStatusToInvoice(invoice, asaasInvoice.Status)

    // Registrar pagamento se foi pago
    if asaasInvoice.PaymentDate != nil {
        paymentDate, _ := time.Parse("2006-01-02", *asaasInvoice.PaymentDate)
        _ = invoice.RegistrarPagamento(paymentDate)
    }

    // Salvar
    return uc.invoiceRepo.Create(ctx, tenantID, invoice)
}

func (uc *SyncInvoicesAsaasUseCase) updateInvoiceFromAsaas(
    ctx context.Context,
    tenantID uuid.UUID,
    invoice *entity.AssinaturaInvoice,
    asaasInvoice *asaas.InvoiceResponse,
) error {
    // Atualizar dados do Asaas
    estimatedDate := parseAsaasDate(asaasInvoice.EstimatedPaymentDate)
    clientDate := parseAsaasDate(asaasInvoice.ClientPaymentDate)

    invoice.AtualizarDadosAsaas(
        asaasInvoice.Status,
        asaasInvoice.InvoiceURL,
        estimatedDate,
        clientDate,
    )

    // Mapear status
    invoice = mapAsaasStatusToInvoice(invoice, asaasInvoice.Status)

    // Registrar pagamento se foi pago e ainda não tinha sido registrado
    if asaasInvoice.PaymentDate != nil && invoice.DataPagamento() == nil {
        paymentDate, _ := time.Parse("2006-01-02", *asaasInvoice.PaymentDate)
        _ = invoice.RegistrarPagamento(paymentDate)
    }

    // Atualizar
    return uc.invoiceRepo.Update(ctx, tenantID, invoice)
}

func parseAsaasDate(dateStr *string) *time.Time {
    if dateStr == nil || *dateStr == "" {
        return nil
    }

    t, err := time.Parse("2006-01-02", *dateStr)
    if err != nil {
        return nil
    }
    return &t
}

func mapAsaasStatusToInvoice(invoice *entity.AssinaturaInvoice, asaasStatus string) *entity.AssinaturaInvoice {
    // Mapear status Asaas → Status interno
    switch asaasStatus {
    case "PENDING":
        // Mantém como PENDENTE
    case "CONFIRMED":
        // Manter como PENDENTE mas atualizar asaas_status
    case "RECEIVED":
        // Status PAGO (já tratado em RegistrarPagamento)
    }
    return invoice
}
```

---

#### 📦 Tarefa 3.4 - Buscar Informações do Cliente

**Arquivo:** `backend/internal/application/usecase/subscription/get_customer_info_usecase.go`

```go
package subscription

import (
    "context"
    "fmt"

    "go.uber.org/zap"

    "backend/internal/infrastructure/external/asaas"
)

type GetCustomerInfoInput struct {
    CustomerID string `json:"customer_id" validate:"required"`
}

type GetCustomerInfoOutput struct {
    CustomerID string `json:"customer_id"`
    Nome       string `json:"nome"`
    Email      string `json:"email"`
    CpfCnpj    string `json:"cpf_cnpj"`
}

type GetCustomerInfoUseCase struct {
    asaasClient *asaas.Client
    logger      *zap.Logger
}

func NewGetCustomerInfoUseCase(
    asaasClient *asaas.Client,
    logger *zap.Logger,
) *GetCustomerInfoUseCase {
    return &GetCustomerInfoUseCase{
        asaasClient: asaasClient,
        logger:      logger,
    }
}

func (uc *GetCustomerInfoUseCase) Execute(
    ctx context.Context,
    input GetCustomerInfoInput,
) (*GetCustomerInfoOutput, error) {

    // Buscar cliente no Asaas
    customer, err := uc.asaasClient.GetCustomer(ctx, input.CustomerID)
    if err != nil {
        uc.logger.Error("erro ao buscar cliente no Asaas",
            zap.String("customer_id", input.CustomerID),
            zap.Error(err),
        )
        return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
    }

    return &GetCustomerInfoOutput{
        CustomerID: customer.ID,
        Nome:       customer.Name,
        Email:      customer.Email,
        CpfCnpj:    customer.CpfCnpj,
    }, nil
}
```

---

### **FASE 4: Cron Jobs (3-4 horas)**

#### 📦 Tarefa 4.1 - Job de Sincronização de Assinaturas

**Arquivo:** `backend/internal/infrastructure/scheduler/jobs/sync_assinaturas_asaas_job.go`

```go
package jobs

import (
    "context"
    "time"

    "go.uber.org/zap"

    "backend/internal/application/usecase/subscription"
    "backend/internal/domain/repository"
)

type SyncAssinaturasAsaasJob struct {
    tenantRepo repository.TenantRepository
    usecase    *subscription.SyncAssinaturasAsaasUseCase
    logger     *zap.Logger
}

func NewSyncAssinaturasAsaasJob(
    tenantRepo repository.TenantRepository,
    usecase *subscription.SyncAssinaturasAsaasUseCase,
    logger *zap.Logger,
) *SyncAssinaturasAsaasJob {
    return &SyncAssinaturasAsaasJob{
        tenantRepo: tenantRepo,
        usecase:    usecase,
        logger:     logger,
    }
}

func (j *SyncAssinaturasAsaasJob) Name() string {
    return "SyncAssinaturasAsaasJob"
}

func (j *SyncAssinaturasAsaasJob) Schedule() string {
    return "0 */30 * * * *" // A cada 30 minutos
}

func (j *SyncAssinaturasAsaasJob) Execute(ctx context.Context) error {
    startTime := time.Now()
    j.logger.Info("iniciando job de sincronização de assinaturas Asaas")

    // Buscar todos os tenants
    tenants, err := j.tenantRepo.FindAll(ctx)
    if err != nil {
        j.logger.Error("erro ao buscar tenants", zap.Error(err))
        return err
    }

    successCount := 0
    errorCount := 0

    // Sincronizar cada tenant
    for _, tenant := range tenants {
        if err := j.usecase.Execute(ctx, tenant.ID()); err != nil {
            j.logger.Error("erro ao sincronizar assinaturas do tenant",
                zap.String("tenant_id", tenant.ID().String()),
                zap.Error(err),
            )
            errorCount++
            continue
        }
        successCount++
    }

    duration := time.Since(startTime)
    j.logger.Info("job de sincronização de assinaturas concluído",
        zap.Int("tenants_synced", successCount),
        zap.Int("errors", errorCount),
        zap.Duration("duration", duration),
    )

    return nil
}
```

---

#### 📦 Tarefa 4.2 - Job de Sincronização de Faturas

**Arquivo:** `backend/internal/infrastructure/scheduler/jobs/sync_invoices_asaas_job.go`

```go
package jobs

import (
    "context"
    "time"

    "go.uber.org/zap"

    "backend/internal/application/usecase/subscription"
    "backend/internal/domain/repository"
)

type SyncInvoicesAsaasJob struct {
    tenantRepo repository.TenantRepository
    usecase    *subscription.SyncInvoicesAsaasUseCase
    logger     *zap.Logger
}

func NewSyncInvoicesAsaasJob(
    tenantRepo repository.TenantRepository,
    usecase *subscription.SyncInvoicesAsaasUseCase,
    logger *zap.Logger,
) *SyncInvoicesAsaasJob {
    return &SyncInvoicesAsaasJob{
        tenantRepo: tenantRepo,
        usecase:    usecase,
        logger:     logger,
    }
}

func (j *SyncInvoicesAsaasJob) Name() string {
    return "SyncInvoicesAsaasJob"
}

func (j *SyncInvoicesAsaasJob) Schedule() string {
    return "0 0 2 * * *" // Diariamente às 02:00
}

func (j *SyncInvoicesAsaasJob) Execute(ctx context.Context) error {
    startTime := time.Now()
    j.logger.Info("iniciando job de sincronização de faturas Asaas")

    // Buscar todos os tenants
    tenants, err := j.tenantRepo.FindAll(ctx)
    if err != nil {
        j.logger.Error("erro ao buscar tenants", zap.Error(err))
        return err
    }

    successCount := 0
    errorCount := 0

    // Sincronizar cada tenant
    for _, tenant := range tenants {
        if err := j.usecase.Execute(ctx, tenant.ID()); err != nil {
            j.logger.Error("erro ao sincronizar faturas do tenant",
                zap.String("tenant_id", tenant.ID().String()),
                zap.Error(err),
            )
            errorCount++
            continue
        }
        successCount++
    }

    duration := time.Since(startTime)
    j.logger.Info("job de sincronização de faturas concluído",
        zap.Int("tenants_synced", successCount),
        zap.Int("errors", errorCount),
        zap.Duration("duration", duration),
    )

    return nil
}
```

---

#### 📦 Tarefa 4.3 - Registrar Jobs no Scheduler

**Arquivo:** `backend/internal/infrastructure/scheduler/setup.go`

Atualizar função `SetupJobs`:

```go
func SetupJobs(
    logger *zap.Logger,
    // ... repositórios existentes ...
    asaasClient *asaas.Client, // ADICIONAR
    tenantRepo repository.TenantRepository, // ADICIONAR
) []Job {

    // ... jobs existentes ...

    // ADICIONAR: Use cases Asaas
    syncAssinaturasUC := subscription.NewSyncAssinaturasAsaasUseCase(
        assinaturaRepo,
        asaasClient,
        logger,
    )

    syncInvoicesUC := subscription.NewSyncInvoicesAsaasUseCase(
        assinaturaRepo,
        invoiceRepo,
        asaasClient,
        logger,
    )

    return []Job{
        // Jobs existentes
        NewSubscriptionValidationJob(logger, assinaturaRepo, invoiceRepo, planoRepo),
        NewFinancialSnapshotJob(logger, receitaRepo, despesaRepo, snapshotRepo),
        NewCommissionProcessingJob(logger, invoiceRepo, assinaturaRepo, receitaRepo, despesaRepo),
        NewAlertsJob(logger, invoiceRepo, assinaturaRepo),

        // NOVOS JOBS ASAAS
        jobs.NewSyncAssinaturasAsaasJob(tenantRepo, syncAssinaturasUC, logger),
        jobs.NewSyncInvoicesAsaasJob(tenantRepo, syncInvoicesUC, logger),
    }
}
```

---

### **FASE 5: Endpoints da API (3-4 horas)**

#### 📦 Tarefa 5.1 - Endpoint: Criar Assinatura com Asaas

**Arquivo:** `backend/internal/infrastructure/http/handler/subscription_asaas_handler.go`

```go
package handler

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "go.uber.org/zap"

    "backend/internal/application/usecase/subscription"
    "backend/internal/infrastructure/http/response"
)

type SubscriptionAsaasHandler struct {
    createAsaasUC      *subscription.CreateAssinaturaAsaasUseCase
    getCustomerInfoUC  *subscription.GetCustomerInfoUseCase
    syncAssinaturasUC  *subscription.SyncAssinaturasAsaasUseCase
    syncInvoicesUC     *subscription.SyncInvoicesAsaasUseCase
    logger             *zap.Logger
}

func NewSubscriptionAsaasHandler(
    createAsaasUC *subscription.CreateAssinaturaAsaasUseCase,
    getCustomerInfoUC *subscription.GetCustomerInfoUseCase,
    syncAssinaturasUC *subscription.SyncAssinaturasAsaasUseCase,
    syncInvoicesUC *subscription.SyncInvoicesAsaasUseCase,
    logger *zap.Logger,
) *SubscriptionAsaasHandler {
    return &SubscriptionAsaasHandler{
        createAsaasUC:     createAsaasUC,
        getCustomerInfoUC: getCustomerInfoUC,
        syncAssinaturasUC: syncAssinaturasUC,
        syncInvoicesUC:    syncInvoicesUC,
        logger:            logger,
    }
}

// POST /api/v1/subscriptions/asaas
func (h *SubscriptionAsaasHandler) CreateAssinaturaAsaas(c echo.Context) error {
    var input subscription.CreateAssinaturaAsaasInput

    if err := c.Bind(&input); err != nil {
        return response.Error(c, http.StatusBadRequest, "Dados inválidos", err)
    }

    // Pegar tenant_id do contexto
    tenantID := getTenantIDFromContext(c)
    input.TenantID = tenantID

    output, err := h.createAsaasUC.Execute(c.Request().Context(), input)
    if err != nil {
        h.logger.Error("erro ao criar assinatura Asaas", zap.Error(err))
        return response.Error(c, http.StatusInternalServerError, "Erro ao criar assinatura", err)
    }

    return response.Success(c, http.StatusCreated, "Assinatura criada com sucesso", output)
}

// GET /api/v1/subscriptions/customer/:customer_id
func (h *SubscriptionAsaasHandler) GetCustomerInfo(c echo.Context) error {
    customerID := c.Param("customer_id")

    input := subscription.GetCustomerInfoInput{
        CustomerID: customerID,
    }

    output, err := h.getCustomerInfoUC.Execute(c.Request().Context(), input)
    if err != nil {
        h.logger.Error("erro ao buscar cliente", zap.Error(err))
        return response.Error(c, http.StatusInternalServerError, "Erro ao buscar cliente", err)
    }

    return response.Success(c, http.StatusOK, "Cliente encontrado", output)
}

// POST /api/v1/subscriptions/sync
func (h *SubscriptionAsaasHandler) ForceSync(c echo.Context) error {
    var req struct {
        SyncType string `json:"sync_type"` // "assinaturas" ou "invoices"
    }

    if err := c.Bind(&req); err != nil {
        return response.Error(c, http.StatusBadRequest, "Dados inválidos", err)
    }

    tenantID := getTenantIDFromContext(c)
    ctx := c.Request().Context()

    var err error
    switch req.SyncType {
    case "assinaturas":
        err = h.syncAssinaturasUC.Execute(ctx, tenantID)
    case "invoices":
        err = h.syncInvoicesUC.Execute(ctx, tenantID)
    default:
        return response.Error(c, http.StatusBadRequest, "sync_type inválido (use 'assinaturas' ou 'invoices')", nil)
    }

    if err != nil {
        h.logger.Error("erro ao sincronizar", zap.String("sync_type", req.SyncType), zap.Error(err))
        return response.Error(c, http.StatusInternalServerError, "Erro ao sincronizar", err)
    }

    return response.Success(c, http.StatusOK, "Sincronização concluída", map[string]string{
        "sync_type": req.SyncType,
    })
}
```

---

#### 📦 Tarefa 5.2 - Registrar Rotas

**Arquivo:** `backend/internal/infrastructure/http/routes/routes.go`

Adicionar rotas:

```go
// Assinaturas Asaas (requer autenticação)
subscriptionAsaasHandler := handler.NewSubscriptionAsaasHandler(
    createAsaasUC,
    getCustomerInfoUC,
    syncAssinaturasUC,
    syncInvoicesUC,
    logger,
)

asaasGroup := api.Group("/subscriptions")
asaasGroup.Use(authMiddleware)
asaasGroup.Use(tenantMiddleware)

asaasGroup.POST("/asaas", subscriptionAsaasHandler.CreateAssinaturaAsaas)
asaasGroup.GET("/customer/:customer_id", subscriptionAsaasHandler.GetCustomerInfo)
asaasGroup.POST("/sync", subscriptionAsaasHandler.ForceSync)
```

---

### **FASE 6: Migrations (1-2 horas)**

#### 📦 Tarefa 6.1 - Migration: Adicionar Colunas Asaas

**Arquivo:** `backend/migrations/014_add_asaas_fields.up.sql`

```sql
-- ========== Adicionar campos Asaas em assinaturas ==========
ALTER TABLE assinaturas
    ADD COLUMN IF NOT EXISTS asaas_status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS ultima_sincronizacao TIMESTAMPTZ;

-- Índice para buscar por status Asaas
CREATE INDEX IF NOT EXISTS idx_assinaturas_asaas_status
    ON assinaturas(tenant_id, asaas_status)
    WHERE asaas_status IS NOT NULL;

-- ========== Adicionar campos Asaas em assinatura_invoices ==========
ALTER TABLE assinatura_invoices
    ADD COLUMN IF NOT EXISTS asaas_status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS data_previsao_recebimento DATE,
    ADD COLUMN IF NOT EXISTS client_payment_date DATE,
    ADD COLUMN IF NOT EXISTS invoice_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS processada BOOLEAN DEFAULT false;

-- Índice para buscar faturas não processadas
CREATE INDEX IF NOT EXISTS idx_invoices_processada
    ON assinatura_invoices(tenant_id, processada, status)
    WHERE processada = false;

-- Índice para buscar por status Asaas
CREATE INDEX IF NOT EXISTS idx_invoices_asaas_status
    ON assinatura_invoices(tenant_id, asaas_status)
    WHERE asaas_status IS NOT NULL;

-- Comentários
COMMENT ON COLUMN assinaturas.asaas_status IS 'Status sincronizado do Asaas: ACTIVE, OVERDUE, AWAITING_PAYMENT';
COMMENT ON COLUMN assinaturas.ultima_sincronizacao IS 'Data/hora da última sincronização com Asaas';
COMMENT ON COLUMN assinatura_invoices.asaas_status IS 'Status sincronizado do Asaas: PENDING, CONFIRMED, RECEIVED';
COMMENT ON COLUMN assinatura_invoices.data_previsao_recebimento IS 'Data prevista de recebimento (estimatedPaymentDate)';
COMMENT ON COLUMN assinatura_invoices.client_payment_date IS 'Data que o cliente realizou o pagamento';
COMMENT ON COLUMN assinatura_invoices.invoice_url IS 'URL da fatura no Asaas';
COMMENT ON COLUMN assinatura_invoices.processada IS 'Indica se a fatura já foi processada para repasse de comissão';
```

**Arquivo:** `backend/migrations/014_add_asaas_fields.down.sql`

```sql
-- Remover índices
DROP INDEX IF EXISTS idx_assinaturas_asaas_status;
DROP INDEX IF EXISTS idx_invoices_processada;
DROP INDEX IF EXISTS idx_invoices_asaas_status;

-- Remover colunas de assinaturas
ALTER TABLE assinaturas
    DROP COLUMN IF EXISTS asaas_status,
    DROP COLUMN IF EXISTS ultima_sincronizacao;

-- Remover colunas de assinatura_invoices
ALTER TABLE assinatura_invoices
    DROP COLUMN IF EXISTS asaas_status,
    DROP COLUMN IF EXISTS data_previsao_recebimento,
    DROP COLUMN IF EXISTS client_payment_date,
    DROP COLUMN IF EXISTS invoice_url,
    DROP COLUMN IF EXISTS processada;
```

---

### **FASE 7: Testes (4-6 horas)**

#### 📦 Tarefa 7.1 - Testes Unitários do Client

**Arquivo:** `backend/internal/infrastructure/external/asaas/client_test.go`

```go
package asaas

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
)

func TestCreateSubscription_Success(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "/subscriptions", r.URL.Path)
        assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{
            "object": "subscription",
            "id": "sub_123",
            "status": "ACTIVE",
            "value": 49.90
        }`))
    }))
    defer server.Close()

    // Client
    client := NewClient("test_api_key", server.URL, zap.NewNop())

    // Execute
    resp, err := client.CreateSubscription(context.Background(), &CreateSubscriptionRequest{
        Customer:    "cus_123",
        BillingType: "CREDIT_CARD",
        Value:       49.90,
        NextDueDate: "2025-12-01",
        Cycle:       "MONTHLY",
    })

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "sub_123", resp.ID)
    assert.Equal(t, "ACTIVE", resp.Status)
}

func TestCreateSubscription_Unauthorized(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte(`{"errors": [{"code": "unauthorized", "description": "API key inválida"}]}`))
    }))
    defer server.Close()

    // Client
    client := NewClient("invalid_key", server.URL, zap.NewNop())

    // Execute
    _, err := client.CreateSubscription(context.Background(), &CreateSubscriptionRequest{})

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unauthorized")
}

func TestListInvoices_Success(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "GET", r.Method)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{
            "object": "list",
            "hasMore": false,
            "data": [
                {
                    "object": "invoice",
                    "id": "inv_123",
                    "status": "CONFIRMED",
                    "value": 49.90
                }
            ]
        }`))
    }))
    defer server.Close()

    // Client
    client := NewClient("test_api_key", server.URL, zap.NewNop())

    // Execute
    resp, err := client.ListInvoices(context.Background(), "sub_123")

    // Assert
    assert.NoError(t, err)
    assert.Len(t, resp.Data, 1)
    assert.Equal(t, "inv_123", resp.Data[0].ID)
}

func TestRetryLogic_TransientError(t *testing.T) {
    attemptCount := 0

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attemptCount++
        if attemptCount < 3 {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"object": "subscription", "id": "sub_123"}`))
    }))
    defer server.Close()

    client := NewClient("test_api_key", server.URL, zap.NewNop())

    resp, err := client.CreateSubscription(context.Background(), &CreateSubscriptionRequest{})

    assert.NoError(t, err)
    assert.Equal(t, "sub_123", resp.ID)
    assert.Equal(t, 3, attemptCount) // Deve ter tentado 3 vezes
}
```

---

#### 📦 Tarefa 7.2 - Testes de Integração

**Arquivo:** `backend/tests/integration/asaas_sync_test.go`

```go
package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSyncAssinaturas_Integration(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()

    // Criar tenant, plano, assinatura
    tenant := createTestTenant(t, db)
    plano := createTestPlano(t, db, tenant.ID())
    assinatura := createTestAssinatura(t, db, tenant.ID(), plano.ID())

    // Mock Asaas client
    mockAsaas := &MockAsaasClient{
        GetSubscriptionFunc: func(ctx context.Context, id string) (*asaas.SubscriptionResponse, error) {
            return &asaas.SubscriptionResponse{
                ID:     id,
                Status: "ACTIVE",
            }, nil
        },
    }

    // Use case
    syncUC := subscription.NewSyncAssinaturasAsaasUseCase(
        assinaturaRepo,
        mockAsaas,
        logger,
    )

    // Execute
    err := syncUC.Execute(context.Background(), tenant.ID())
    require.NoError(t, err)

    // Verify
    updated, err := assinaturaRepo.FindByID(context.Background(), tenant.ID(), assinatura.ID())
    require.NoError(t, err)
    assert.NotNil(t, updated.AsaasStatus())
    assert.Equal(t, "ACTIVE", *updated.AsaasStatus())
}

func TestSyncInvoices_Idempotency(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()

    tenant := createTestTenant(t, db)
    assinatura := createTestAssinatura(t, db, tenant.ID(), uuid.New())

    mockAsaas := &MockAsaasClient{
        ListInvoicesFunc: func(ctx context.Context, subID string) (*asaas.InvoiceListResponse, error) {
            return &asaas.InvoiceListResponse{
                Data: []asaas.InvoiceResponse{
                    {
                        ID:     "inv_123",
                        Status: "CONFIRMED",
                        Value:  49.90,
                    },
                },
            }, nil
        },
    }

    syncUC := subscription.NewSyncInvoicesAsaasUseCase(
        assinaturaRepo,
        invoiceRepo,
        mockAsaas,
        logger,
    )

    // Execute 2 vezes
    err1 := syncUC.Execute(context.Background(), tenant.ID())
    err2 := syncUC.Execute(context.Background(), tenant.ID())

    require.NoError(t, err1)
    require.NoError(t, err2)

    // Verify: Deve ter criado apenas 1 fatura
    invoices, err := invoiceRepo.FindByAssinatura(context.Background(), tenant.ID(), assinatura.ID())
    require.NoError(t, err)
    assert.Len(t, invoices, 1)
}
```

---

### **FASE 8: Documentação & Monitoramento (2-3 horas)**

#### 📦 Tarefa 8.1 - Atualizar Documentação da API

**Arquivo:** `backend/docs/API_REFERENCE.md`

Adicionar seção:

```markdown
## Assinaturas Asaas

### POST /api/v1/subscriptions/asaas

Cria uma assinatura integrada com Asaas.

**Autenticação:** Bearer Token
**Permissões:** `create_assinatura`

**Request:**
```json
{
  "plan_id": "uuid-do-plano",
  "barbeiro_id": "uuid-do-barbeiro",
  "customer_id": "cus_000000000000000001",
  "billing_type": "CREDIT_CARD",
  "data_inicio": "2025-12-01"
}
```

**Response 201:**
```json
{
  "code": 201,
  "message": "Assinatura criada com sucesso",
  "data": {
    "id": "uuid",
    "asaas_subscription_id": "sub_000000000000000001",
    "status": "ATIVA",
    "asaas_status": "ACTIVE",
    "proxima_fatura_data": "2025-12-01"
  }
}
```

**Erros:**
- `400` - Dados inválidos
- `500` - Erro ao criar assinatura no Asaas

---

### GET /api/v1/subscriptions/customer/:customer_id

Busca informações de um cliente no Asaas.

**Response 200:**
```json
{
  "code": 200,
  "data": {
    "customer_id": "cus_000000000000000001",
    "nome": "João Silva",
    "email": "joao@example.com",
    "cpf_cnpj": "12345678901"
  }
}
```

---

### POST /api/v1/subscriptions/sync

Força sincronização manual com Asaas.

**Request:**
```json
{
  "sync_type": "assinaturas"
}
```

**Response 200:**
```json
{
  "code": 200,
  "message": "Sincronização concluída",
  "data": {
    "sync_type": "assinaturas"
  }
}
```
```

---

#### 📦 Tarefa 8.2 - Métricas Prometheus

**Arquivo:** `backend/internal/infrastructure/external/asaas/metrics.go`

```go
package asaas

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    AsaasRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "asaas_requests_total",
            Help: "Total de requests para API Asaas",
        },
        []string{"method", "status_code"},
    )

    AsaasRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "asaas_request_duration_seconds",
            Help:    "Duração das requests para Asaas",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method"},
    )

    AsaasErrorsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "asaas_errors_total",
            Help: "Total de erros na integração Asaas",
        },
        []string{"error_type"},
    )
)
```

Usar nas funções do client:

```go
func (c *Client) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
    start := time.Now()
    defer func() {
        AsaasRequestDuration.WithLabelValues("CreateSubscription").Observe(time.Since(start).Seconds())
    }()

    // ... lógica ...

    AsaasRequestsTotal.WithLabelValues("CreateSubscription", fmt.Sprintf("%d", resp.StatusCode)).Inc()

    if err != nil {
        AsaasErrorsTotal.WithLabelValues("api_error").Inc()
        return nil, err
    }

    return result, nil
}
```

---

#### 📦 Tarefa 8.3 - Alertas Grafana

**Arquivo:** `backend/docs/observability/grafana/asaas_alerts.json`

```json
{
  "alerts": [
    {
      "name": "Taxa de Erro Asaas Alta",
      "condition": "rate(asaas_errors_total[5m]) > 0.05",
      "severity": "warning",
      "description": "Taxa de erro na API Asaas está acima de 5%"
    },
    {
      "name": "Latência Asaas Alta",
      "condition": "histogram_quantile(0.95, asaas_request_duration_seconds) > 5",
      "severity": "warning",
      "description": "Latência P95 da API Asaas está acima de 5 segundos"
    },
    {
      "name": "Asaas Indisponível",
      "condition": "up{job=\"asaas\"} == 0 for 10m",
      "severity": "critical",
      "description": "API Asaas está indisponível por mais de 10 minutos"
    },
    {
      "name": "Falha na Sincronização",
      "condition": "increase(asaas_sync_failures_total[1h]) > 3",
      "severity": "critical",
      "description": "Mais de 3 falhas consecutivas na sincronização"
    }
  ]
}
```

---

## ✅ Checklist Completo de Implementação

### **FASE 1 - Infraestrutura Base (4-6h)**
- [ ] 1.1 - Criar `AsaasClient` com métodos básicos
- [ ] 1.2 - Implementar retry logic com exponential backoff
- [ ] 1.3 - Criar todos os DTOs (Request/Response)
- [ ] 1.4 - Adicionar configuração de ambiente (.env)
- [ ] 1.5 - Atualizar `config.go` com campos Asaas
- [ ] 1.6 - Testes unitários do client (5+ casos)

### **FASE 2 - Entidades (2-3h)**
- [ ] 2.1 - Adicionar campos Asaas em `Assinatura`
- [ ] 2.2 - Adicionar getters e setters
- [ ] 2.3 - Adicionar campos Asaas em `AssinaturaInvoice`
- [ ] 2.4 - Atualizar métodos `Reconstruct`
- [ ] 2.5 - Testes unitários das entidades

### **FASE 3 - Casos de Uso (6-8h)**
- [ ] 3.1 - `CreateAssinaturaAsaasUseCase`
- [ ] 3.2 - `SyncAssinaturasAsaasUseCase`
- [ ] 3.3 - `SyncInvoicesAsaasUseCase`
- [ ] 3.4 - `GetCustomerInfoUseCase`
- [ ] 3.5 - Testes de cada use case (mock Asaas)
- [ ] 3.6 - Tratamento de erros e rollback

### **FASE 4 - Cron Jobs (3-4h)**
- [ ] 4.1 - `SyncAssinaturasAsaasJob` (30 min)
- [ ] 4.2 - `SyncInvoicesAsaasJob` (diário 02:00)
- [ ] 4.3 - Registrar jobs no scheduler
- [ ] 4.4 - Testar execução manual dos jobs
- [ ] 4.5 - Logs estruturados

### **FASE 5 - Endpoints API (3-4h)**
- [ ] 5.1 - Handler `CreateAssinaturaAsaas`
- [ ] 5.2 - Handler `GetCustomerInfo`
- [ ] 5.3 - Handler `ForceSync`
- [ ] 5.4 - Registrar rotas
- [ ] 5.5 - Middlewares (auth, tenant)
- [ ] 5.6 - Testes de integração dos endpoints

### **FASE 6 - Migrations (1-2h)**
- [ ] 6.1 - Migration 014 (UP)
- [ ] 6.2 - Migration 014 (DOWN)
- [ ] 6.3 - Testar aplicação da migration
- [ ] 6.4 - Testar rollback
- [ ] 6.5 - Validar índices criados

### **FASE 7 - Testes (4-6h)**
- [ ] 7.1 - Testes unitários (client)
- [ ] 7.2 - Testes unitários (use cases)
- [ ] 7.3 - Testes de integração (DB + mock Asaas)
- [ ] 7.4 - Testes com Asaas Sandbox
- [ ] 7.5 - Testes de idempotência
- [ ] 7.6 - Testes de tratamento de erros
- [ ] 7.7 - Cobertura >80%

### **FASE 8 - Documentação & Monitoramento (2-3h)**
- [ ] 8.1 - Atualizar `API_REFERENCE.md`
- [ ] 8.2 - Atualizar `INTEGRACOES_ASAAS.md`
- [ ] 8.3 - Adicionar métricas Prometheus
- [ ] 8.4 - Configurar alertas Grafana
- [ ] 8.5 - Criar runbook de troubleshooting

### **FASE 9 - Deploy & Validação (2-3h)**
- [ ] 9.1 - Deploy em staging
- [ ] 9.2 - Executar smoke tests
- [ ] 9.3 - Criar assinatura de teste
- [ ] 9.4 - Forçar sincronização
- [ ] 9.5 - Validar logs e métricas
- [ ] 9.6 - Ativar feature flag `FEATURE_ASAAS_INTEGRATION=true`
- [ ] 9.7 - Deploy em produção (gradual)

---

## ⏱️ Estimativa de Tempo

| Fase | Tempo Estimado | Prioridade | Complexidade |
|------|----------------|------------|--------------|
| **FASE 1** - Infraestrutura | 4-6 horas | 🔴 ALTA | Média |
| **FASE 2** - Entidades | 2-3 horas | 🔴 ALTA | Baixa |
| **FASE 3** - Casos de Uso | 6-8 horas | 🔴 ALTA | Alta |
| **FASE 4** - Cron Jobs | 3-4 horas | 🟡 MÉDIA | Média |
| **FASE 5** - Endpoints | 3-4 horas | 🔴 ALTA | Média |
| **FASE 6** - Migrations | 1-2 horas | 🔴 ALTA | Baixa |
| **FASE 7** - Testes | 4-6 horas | 🟡 MÉDIA | Média |
| **FASE 8** - Documentação | 2-3 horas | 🟢 BAIXA | Baixa |
| **FASE 9** - Deploy | 2-3 horas | 🔴 ALTA | Média |
| **TOTAL** | **27-39 horas** | - | **~5-7 dias** |

---

## 🎯 Próximos Passos

### **Opção A - Implementação Completa (Recomendado)**
**Tempo:** 5-7 dias
**Resultado:** Integração completa com Asaas

**Passos:**
1. Criar conta no Asaas (sandbox + produção)
2. Gerar API Keys
3. Seguir fases 1-9 em ordem
4. Validar cada fase antes de prosseguir

---

### **Opção B - MVP Rápido**
**Tempo:** 2-3 dias
**Resultado:** Apenas sincronização de faturas

**Passos:**
1. **FASE 1** - Client básico (sem retry)
2. **FASE 3.3** - Apenas `SyncInvoicesAsaasUseCase`
3. **FASE 5.2** - Endpoint listar faturas
4. **FASE 6** - Migration
5. Testar em sandbox

---

### **Opção C - Continuar Manual**
**Tempo:** 0 dias
**Resultado:** Sistema continua funcionando 100% manual

**Motivo:** Sistema já funciona perfeitamente sem Asaas

---

## 🚨 Riscos e Mitigações

| Risco | Impacto | Probabilidade | Mitigação |
|-------|---------|---------------|-----------|
| **API Asaas indisponível** | Alto | Média | Retry logic + cache local + fallback manual |
| **Rate limiting Asaas** | Médio | Média | Exponential backoff + respeitar limites (100 req/min) |
| **Dados inconsistentes** | Alto | Baixa | Sincronização idempotente + audit logs |
| **Chave API vazada** | Crítico | Baixa | Nunca commitar .env + rotar chaves mensalmente |
| **Timeout em sync** | Médio | Média | Timeout configurável + jobs em background |
| **Erro no rollback** | Alto | Baixa | Testar rollback em staging + logs detalhados |

---

## 📚 Referências

- [Documentação Asaas](https://docs.asaas.com)
- [API Asaas - Subscriptions](https://docs.asaas.com/reference/crear-una-suscripcin)
- [API Asaas - Invoices](https://docs.asaas.com/reference/obtener-facturas)
- [Documentação Interna - INTEGRACOES_ASAAS.md](/docs/INTEGRACOES_ASAAS.md)
- [Documentação Interna - ASSINATURAS.md](/docs/ASSINATURAS.md)

---

**Status:** 📝 Planejamento Completo
**Próxima Ação:** Aguardando decisão do time (Opção A, B ou C)
**Data:** 17/11/2025
