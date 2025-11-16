# 🎟️ Módulo de Assinaturas & Clube

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Design Finalizado

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Entidades de Domínio](#entidades-de-domínio)
3. [Casos de Uso](#casos-de-uso)
4. [Integração Asaas](#integração-asaas)
5. [Schema do Banco](#schema-do-banco)
6. [Fluxo Repasse Barbeiro](#fluxo-repasse-barbeiro)

---

## 🎯 Visão Geral

Módulo de **gerenciamento de planos de assinatura** para barbeiros. Permite criar planos (ex: 4 cortes/mês), vender para barbeiros, integrar com Asaas para pagamentos recorrentes, e automatizar repasses pós-vencimento.

---

## 🏛️ Entidades de Domínio

### PlanoDeassinatura (Aggregate Root)

```go
type PlanoDeassinatura struct {
    ID              string
    TenantID        string
    Nome            string
    Descricao       string
    Valor           decimal.Decimal
    Periodicidade   Periodicidade  // MENSAL, TRIMESTRAL, ANUAL
    QuantidadeServicos int          // Ex: 4 cortes
    Ativa           bool
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

type Periodicidade string

const (
    PeriodMensal    Periodicidade = "MENSAL"
    PeriodTrimestral Periodicidade = "TRIMESTRAL"
    PeriodAnual     Periodicidade = "ANUAL"
)
```

### Assinatura (Aggregate Root)

```go
type Assinatura struct {
    ID                  string
    TenantID            string
    PlanID              string
    BarbeiroID          string
    AsaasSubscriptionID string  // ID da sub no Asaas
    Status              SubscriptionStatus
    DataInicio          time.Time
    DataFim             *time.Time
    ProximaFaturaData   time.Time
    CriadoEm            time.Time
    AtualizadoEm        time.Time
}

type SubscriptionStatus string

const (
    SubActive     SubscriptionStatus = "ATIVA"
    SubCancelled  SubscriptionStatus = "CANCELADA"
    SubSuspended  SubscriptionStatus = "SUSPENSA"
    SubExpired    SubscriptionStatus = "EXPIRADA"
)
```

### AssinaturaInvoice (Entity)

```go
type AssinaturaInvoice struct {
    ID              string
    TenantID        string
    AssinaturaID    string
    AsaasInvoiceID  string  // ID da fatura no Asaas
    Valor           decimal.Decimal
    Status          InvoiceStatus
    DataVencimento  time.Time
    DataPagamento   *time.Time
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

type InvoiceStatus string

const (
    InvoicePending    InvoiceStatus = "PENDENTE"
    InvoiceConfirmed  InvoiceStatus = "CONFIRMADA"
    InvoiceReceived   InvoiceStatus = "RECEBIDA"
    InvoiceCancelled  InvoiceStatus = "CANCELADA"
    InvoiceRefunded   InvoiceStatus = "REEMBOLSADA"
)
```

---

## 💼 Casos de Uso

### CreateAssinaturaUseCase

```go
type CreateAssinaturaUseCase struct {
    repository     domain.AssinaturaRepository
    asaasClient    *external.AsaasClient
    validator      domain.AssinaturaValidator
}

type CreateAssinaturaInput struct {
    TenantID    string `json:"tenant_id" validate:"required"`
    PlanID      string `json:"plan_id" validate:"required"`
    BarbeiroID  string `json:"barbeiro_id" validate:"required"`
    DataInicio  time.Time `json:"data_inicio" validate:"required"`
}

func (uc *CreateAssinaturaUseCase) Execute(
    ctx context.Context, input CreateAssinaturaInput) (*CreateAssinaturaOutput, error) {
    
    // 1. Validar input
    if err := uc.validator.Validate(input); err != nil {
        return nil, err
    }
    
    // 2. Buscar plano
    plan, err := uc.repository.FindPlanByID(ctx, input.TenantID, input.PlanID)
    if err != nil {
        return nil, ErrPlanNotFound
    }
    
    // 3. Criar no Asaas
    asaasResp, err := uc.asaasClient.CreateSubscription(&external.CreateSubscriptionReq{
        CustomerID: input.BarbeiroID, // ou customer externo
        BillingType: "RECURRING",
        Value: plan.Valor,
        Cycle: "MONTHLY", // TODO: map periodicidade
        NextDueDate: input.DataInicio,
    })
    if err != nil {
        return nil, ErrAsaasCreationFailed
    }
    
    // 4. Criar entidade de domínio
    assinatura := &domain.Assinatura{
        ID:                  uuid.NewString(),
        TenantID:            input.TenantID,
        PlanID:              input.PlanID,
        BarbeiroID:          input.BarbeiroID,
        AsaasSubscriptionID: asaasResp.ID,
        Status:              domain.SubActive,
        DataInicio:          input.DataInicio,
        ProximaFaturaData:   input.DataInicio.AddDate(0, 1, 0),
        CriadoEm:            time.Now(),
        AtualizadoEm:        time.Now(),
    }
    
    // 5. Persistir
    if err := uc.repository.SaveAssinatura(ctx, input.TenantID, assinatura); err != nil {
        // TODO: rollback no Asaas
        return nil, err
    }
    
    return &CreateAssinaturaOutput{
        ID:     assinatura.ID,
        Status: string(assinatura.Status),
    }, nil
}
```

### SincronizarFaturasAsaasUseCase

```go
type SincronizarFaturasAsaasUseCase struct {
    invoiceRepo  domain.AssinaturaInvoiceRepository
    asaasClient  *external.AsaasClient
}

func (uc *SincronizarFaturasAsaasUseCase) Execute(ctx context.Context, tenantID string) error {
    // 1. Buscar todas assinaturas ativas do tenant
    assinaturas, err := uc.repository.FindByTenant(ctx, tenantID, domain.SubActive)
    if err != nil {
        return err
    }
    
    // 2. Para cada assinatura, buscar faturas no Asaas
    for _, sub := range assinaturas {
        invoices, err := uc.asaasClient.ListInvoices(&external.ListInvoicesReq{
            SubscriptionID: sub.AsaasSubscriptionID,
            Status: []string{"PENDING", "CONFIRMED", "RECEIVED"},
        })
        if err != nil {
            log.Error("Failed to sync invoices", "subscription_id", sub.ID, "err", err)
            continue
        }
        
        // 3. Sincronizar cada fatura
        for _, asaasInvoice := range invoices {
            // Verificar se já existe
            existing, _ := uc.invoiceRepo.FindByAsaasID(ctx, tenantID, asaasInvoice.ID)
            
            if existing == nil {
                // Criar nova
                invoice := &domain.AssinaturaInvoice{
                    ID:              uuid.NewString(),
                    TenantID:        tenantID,
                    AssinaturaID:    sub.ID,
                    AsaasInvoiceID:  asaasInvoice.ID,
                    Valor:           decimal.NewFromFloat(asaasInvoice.Value),
                    Status:          mapAsaasStatus(asaasInvoice.Status),
                    DataVencimento:  asaasInvoice.DueDate,
                    DataPagamento:   asaasInvoice.PaidDate,
                    CriadoEm:        time.Now(),
                    AtualizadoEm:    time.Now(),
                }
                uc.invoiceRepo.Save(ctx, tenantID, invoice)
            } else {
                // Atualizar status
                existing.Status = mapAsaasStatus(asaasInvoice.Status)
                existing.DataPagamento = asaasInvoice.PaidDate
                existing.AtualizadoEm = time.Now()
                uc.invoiceRepo.Update(ctx, tenantID, existing)
            }
        }
    }
    
    return nil
}
```

---

## 🔗 Integração Asaas

### AsaasClient

```go
type AsaasClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

// Criar assinatura
func (c *AsaasClient) CreateSubscription(req *CreateSubscriptionReq) (*CreateSubscriptionResp, error) {
    // POST https://api.asaas.com/v3/subscriptions
    // Headers: Authorization: Bearer {apiKey}
    // Body: {"customerEmail": ..., "value": ..., "nextDueDate": ..., ...}
}

// Listar faturas de uma assinatura
func (c *AsaasClient) ListInvoices(req *ListInvoicesReq) ([]*Invoice, error) {
    // GET https://api.asaas.com/v3/invoices?subscriptionId={id}
    // Retorna array de invoices
}

// Cancelar assinatura
func (c *AsaasClient) CancelSubscription(subscriptionID string) error {
    // POST https://api.asaas.com/v3/subscriptions/{id}/cancel
}
```

### Error Handling

```go
// Possíveis erros Asaas:
// - 401: API key inválida
// - 422: Dados inválidos
// - 429: Rate limit
// - 5xx: Servidor Asaas

func handleAsaasError(resp *http.Response) error {
    if resp.StatusCode == 401 {
        return ErrAsaasUnauthorized
    }
    if resp.StatusCode == 422 {
        return ErrAsaasInvalidData
    }
    if resp.StatusCode == 429 {
        // Implementar retry com backoff
        return ErrAsaasRateLimit
    }
    return ErrAsaasServerError
}
```

---

## 🗄️ Schema do Banco

### Tabela: planos_assinatura

```sql
CREATE TABLE planos_assinatura (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    valor DECIMAL(10, 2) NOT NULL CHECK (valor > 0),
    periodicidade VARCHAR(50) NOT NULL,
    quantidade_servicos INT DEFAULT 0,
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, nome)
);
```

### Tabela: assinaturas

```sql
CREATE TABLE assinaturas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES planos_assinatura(id) ON DELETE RESTRICT,
    barbeiro_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    asaas_subscription_id VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ATIVA',
    data_inicio DATE NOT NULL,
    data_fim DATE,
    proxima_fatura_data DATE NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_assinaturas_tenant ON assinaturas(tenant_id);
CREATE INDEX idx_assinaturas_status ON assinaturas(tenant_id, status);
CREATE INDEX idx_assinaturas_barbeiro ON assinaturas(barbeiro_id);
```

### Tabela: assinatura_invoices

```sql
CREATE TABLE assinatura_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    assinatura_id UUID NOT NULL REFERENCES assinaturas(id) ON DELETE CASCADE,
    asaas_invoice_id VARCHAR(255) UNIQUE NOT NULL,
    valor DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDENTE',
    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_invoices_tenant ON assinatura_invoices(tenant_id);
CREATE INDEX idx_invoices_assinatura ON assinatura_invoices(assinatura_id);
CREATE INDEX idx_invoices_status ON assinatura_invoices(status, data_vencimento);
```

---

## 💳 Fluxo Repasse Barbeiro

### Contexto

Barbeiro tem assinatura ativa. Quando a fatura é **PAGA/RECEBIDA**, a comissão é repassada ao barbeiro no dia combinado.

### Fluxo

```
Fatura Criada (PENDENTE)
    ↓ (via Cron diário: sincronizar com Asaas)
Fatura Confirmada (CONFIRMADA)
    ↓ (quando cliente paga)
Fatura Recebida (RECEBIDA)
    ↓ (Cron: ProcessarRepasseUseCase)
Criar Receita no financeiro (comissão para barbeiro)
    ↓
Gerar Despesa (comissão saída do caixa)
    ↓
Contabilizar em Fluxo de Caixa
```

### ProcessarRepasseUseCase

```go
type ProcessarRepasseUseCase struct {
    invoiceRepo   domain.AssinaturaInvoiceRepository
    receitaRepo   domain.ReceitaRepository
    despesaRepo   domain.DespesaRepository
    configRepo    domain.ConfiguracaoRepository
}

func (uc *ProcessarRepasseUseCase) Execute(ctx context.Context, tenantID string) error {
    // 1. Buscar faturas pagas que ainda não foram repassadas
    invoices, err := uc.invoiceRepo.FindUnprocessed(ctx, tenantID, domain.InvoiceReceived)
    if err != nil {
        return err
    }
    
    for _, invoice := range invoices {
        // 2. Buscar assinatura
        assinatura, _ := uc.assinaturaRepository.FindByID(ctx, tenantID, invoice.AssinaturaID)
        
        // 3. Buscar plano
        plan, _ := uc.planRepository.FindByID(ctx, tenantID, assinatura.PlanID)
        
        // 4. Buscar configuração de repasse (% ou valor fixo)
        config, _ := uc.configRepo.FindByTenant(ctx, tenantID)
        
        // 5. Calcular repasse
        // Ex: 70% para barbeiro, 30% para barbearia
        percentualBarbeiro := config.PercentualRepasseBarbeiro // 0.7
        valorRepasse := invoice.Valor.Mul(decimal.NewFromFloat(percentualBarbeiro))
        
        // 6. Criar receita (entrada de assinatura)
        receita := &domain.Receita{
            ID:          uuid.NewString(),
            TenantID:    tenantID,
            UsuarioID:   "", // System
            Descricao:   fmt.Sprintf("Assinatura %s - Barbeiro %s", plan.Nome, assinatura.BarbeiroID),
            Valor:       invoice.Valor,
            Status:      domain.ReceiptReceived,
            Data:        time.Now(),
            CriadoEm:    time.Now(),
            AtualizadoEm: time.Now(),
        }
        uc.receitaRepo.Save(ctx, tenantID, receita)
        
        // 7. Criar despesa (comissão para barbeiro)
        despesa := &domain.Despesa{
            ID:          uuid.NewString(),
            TenantID:    tenantID,
            UsuarioID:   assinatura.BarbeiroID,
            Descricao:   fmt.Sprintf("Comissão - Plano %s", plan.Nome),
            Valor:       valorRepasse,
            Status:      domain.ExpensePaid,
            Data:        time.Now(), // ou data_pagamento configurada
            CriadoEm:    time.Now(),
            AtualizadoEm: time.Now(),
        }
        uc.despesaRepo.Save(ctx, tenantID, despesa)
        
        // 8. Marcar fatura como processada
        invoice.Processada = true
        uc.invoiceRepo.Update(ctx, tenantID, invoice)
    }
    
    return nil
}
```

### Regra de Negócio

- RN-SUB-001: Repasse só ocorre se fatura está "RECEBIDA"
- RN-SUB-002: Repasse não pode ocorrer duas vezes para mesma fatura
- RN-SUB-003: Percentual de repasse é configurável por tenant
- RN-SUB-004: Repasse pode ter data de execução adiada

---

**Status:** ✅ Design finalizado
