> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 💰 Módulo Financeiro - Domain & Implementação

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Design Finalizado

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Entidades de Domínio](#entidades-de-domínio)
3. [Casos de Uso](#casos-de-uso)
4. [Regras de Negócio](#regras-de-negócio)
5. [Schema do Banco](#schema-do-banco)
6. [Fluxo de Dados](#fluxo-de-dados)

---

## 🎯 Visão Geral

Módulo responsável por **gerenciar receitas, despesas, categorias e fluxo de caixa** da barbearia. Suporta múltiplas formas de pagamento e integração com sistema de assinaturas.

### Agregados Principais

```
Receita (Aggregate Root)
├── ID
├── TenantID
├── Descricao
├── Valor
├── Categoria (Value Object)
├── MetodoPagamento (Value Object)
├── Data
└── CriadoEm

Despesa (Aggregate Root)
├── ID
├── TenantID
├── Descricao
├── Valor
├── Categoria (Value Object)
├── Fornecedor
├── MetodoPagamento (Value Object)
├── Data
└── CriadoEm

Categoria (Entity)
├── ID
├── TenantID
├── Nome
├── Tipo (RECEITA ou DESPESA)
└── Ativa

FluxoDeCaixa (Value Object / Read Model)
├── Periodo
├── SaldoInicial
├── Entradas
├── Saidas
├── SaldoFinal
```

---

## 🏛️ Entidades de Domínio

### 1. Receita

```go
package domain

import (
    "time"
    "github.com/shopspring/decimal"
)

type Receita struct {
    ID              string
    TenantID        string
    UsuarioID       string
    Descricao       string
    Valor           decimal.Decimal // Usar Decimal, não float64
    Categoria       *Categoria
    MetodoPagamento MetodoPagamento
    Data            time.Time
    Status          ReceiptStatus // CONFIRMADO, RECEBIDO, CANCELADO
    Observacoes     string
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

type ReceiptStatus string

const (
    ReceiptConfirmed ReceiptStatus = "CONFIRMADO"
    ReceiptReceived  ReceiptStatus = "RECEBIDO"
    ReceiptCancelled ReceiptStatus = "CANCELADO"
)

// Validações de domínio
func (r *Receita) Validate() error {
    if r.ID == "" {
        return ErrReceitaIDRequired
    }
    if r.TenantID == "" {
        return ErrTenantIDRequired
    }
    if r.Valor.LessThanOrEqual(decimal.Zero) {
        return ErrReceitaValorInvalido
    }
    if r.Data.IsZero() {
        return ErrReceitaDataRequired
    }
    return nil
}

// Business logic
func (r *Receita) Confirmar() error {
    if r.Status != ReceiptConfirmed {
        r.Status = ReceiptConfirmed
        r.AtualizadoEm = time.Now()
    }
    return nil
}

func (r *Receita) Cancelar() error {
    if r.Status == ReceiptCancelled {
        return ErrReceitaJaCancelada
    }
    r.Status = ReceiptCancelled
    r.AtualizadoEm = time.Now()
    return nil
}
```

### 2. Despesa

```go
type Despesa struct {
    ID              string
    TenantID        string
    UsuarioID       string
    Descricao       string
    Valor           decimal.Decimal
    Categoria       *Categoria
    Fornecedor      string
    MetodoPagamento MetodoPagamento
    Data            time.Time
    Status          ExpenseStatus
    Observacoes     string
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

type ExpenseStatus string

const (
    ExpensePending  ExpenseStatus = "PENDENTE"
    ExpensePaid     ExpenseStatus = "PAGO"
    ExpenseCancelled ExpenseStatus = "CANCELADO"
)

func (d *Despesa) Validate() error {
    if d.ID == "" {
        return ErrDespesaIDRequired
    }
    if d.TenantID == "" {
        return ErrTenantIDRequired
    }
    if d.Valor.LessThanOrEqual(decimal.Zero) {
        return ErrDespesaValorInvalido
    }
    if d.Data.IsZero() {
        return ErrDespesaDataRequired
    }
    return nil
}

func (d *Despesa) Marcar() error {
    if d.Status != ExpensePending {
        return ErrDespesaJaPaga
    }
    d.Status = ExpensePaid
    d.AtualizadoEm = time.Now()
    return nil
}
```

### 3. Categoria (Value Object)

```go
type Categoria struct {
    ID          string
    TenantID    string
    Nome        string
    Tipo        TipoCategoria
    Cor         string // #RRGGBB
    Ativa       bool
    CriadoEm    time.Time
}

type TipoCategoria string

const (
    TipoCategoriaReceita TipoCategoria = "RECEITA"
    TipoCategoriaDespesa TipoCategoria = "DESPESA"
)

func (c *Categoria) Validate() error {
    if c.Nome == "" {
        return ErrCategoriaNomeRequired
    }
    if c.Tipo != TipoCategoriaReceita && c.Tipo != TipoCategoriaDespesa {
        return ErrCategoriaTipoInvalido
    }
    return nil
}
```

### 4. MetodoPagamento (Value Object)

```go
type MetodoPagamento string

const (
    MetodoDinheiro    MetodoPagamento = "DINHEIRO"
    MetodoDebito      MetodoPagamento = "DEBITO"
    MetodoCredito     MetodoPagamento = "CREDITO"
    MetodoPix         MetodoPagamento = "PIX"
    MetodoTransferencia MetodoPagamento = "TRANSFERENCIA"
)

func (m MetodoPagamento) IsValid() bool {
    switch m {
    case MetodoDinheiro, MetodoDebito, MetodoCredito, MetodoPix, MetodoTransferencia:
        return true
    default:
        return false
    }
}
```

### 5. FluxoDeCaixa (Read Model / DTO)

```go
type FluxoDeCaixa struct {
    TenantID      string
    Periodo       Periodo
    SaldoInicial  decimal.Decimal
    Entradas      decimal.Decimal // Receitas confirmadas
    Saidas        decimal.Decimal  // Despesas pagas
    SaldoFinal    decimal.Decimal
}

type Periodo struct {
    DataInicio time.Time
    DataFim    time.Time
}

func (f *FluxoDeCaixa) Calcular() error {
    if f.DataInicio.After(f.DataFim) {
        return ErrPeriodoInvalido
    }
    f.SaldoFinal = f.SaldoInicial.Add(f.Entradas).Sub(f.Saidas)
    return nil
}
```

---

## 💼 Casos de Uso

### 1. CreateReceitaUseCase

```go
type CreateReceitaUseCase struct {
    repository domain.ReceitaRepository
    validator  domain.ReceitaValidator
}

type CreateReceitaInput struct {
    TenantID        string    `json:"tenant_id" validate:"required"`
    Descricao       string    `json:"descricao" validate:"required,max=255"`
    Valor           string    `json:"valor" validate:"required,numeric"` // String para precisão
    CategoriaID     string    `json:"categoria_id" validate:"required"`
    MetodoPagamento string    `json:"metodo_pagamento" validate:"required"`
    Data            time.Time `json:"data" validate:"required"`
    Observacoes     string    `json:"observacoes" validate:"max=500"`
}

type CreateReceitaOutput struct {
    ID              string    `json:"id"`
    TenantID        string    `json:"tenant_id"`
    Descricao       string    `json:"descricao"`
    Valor           string    `json:"valor"`
    CategoriaID     string    `json:"categoria_id"`
    MetodoPagamento string    `json:"metodo_pagamento"`
    Status          string    `json:"status"`
    CriadoEm        time.Time `json:"criado_em"`
}

func (uc *CreateReceitaUseCase) Execute(
    ctx context.Context, input CreateReceitaInput) (*CreateReceitaOutput, error) {
    
    // 1. Validar input
    if err := uc.validator.Validate(input); err != nil {
        return nil, ErrInvalidInput
    }
    
    // 2. Converter valor para Decimal
    valor, err := decimal.NewFromString(input.Valor)
    if err != nil || valor.LessThanOrEqual(decimal.Zero) {
        return nil, ErrReceitaValorInvalido
    }
    
    // 3. Criar entidade de domínio
    receita := &domain.Receita{
        ID:              uuid.NewString(),
        TenantID:        input.TenantID,
        Descricao:       input.Descricao,
        Valor:           valor,
        CategoriaID:     input.CategoriaID,
        MetodoPagamento: domain.MetodoPagamento(input.MetodoPagamento),
        Data:            input.Data,
        Observacoes:     input.Observacoes,
        Status:          domain.ReceiptConfirmed,
        CriadoEm:        time.Now(),
        AtualizadoEm:    time.Now(),
    }
    
    // 4. Validar regras de negócio
    if err := receita.Validate(); err != nil {
        return nil, err
    }
    
    // 5. Persistir
    if err := uc.repository.Save(ctx, input.TenantID, receita); err != nil {
        return nil, err
    }
    
    // 6. Retornar output
    return &CreateReceitaOutput{
        ID:              receita.ID,
        TenantID:        receita.TenantID,
        Descricao:       receita.Descricao,
        Valor:           receita.Valor.String(),
        CategoriaID:     input.CategoriaID,
        MetodoPagamento: string(receita.MetodoPagamento),
        Status:          string(receita.Status),
        CriadoEm:        receita.CriadoEm,
    }, nil
}
```

### 2. ListReceitasUseCase

```go
type ListReceitasUseCase struct {
    repository domain.ReceitaRepository
}

type ListReceitasInput struct {
    TenantID   string
    DataInicio time.Time `json:"data_inicio"`
    DataFim    time.Time `json:"data_fim"`
    CategoriaID string   `json:"categoria_id"` // Opcional
    Page       int       `json:"page" default:"1"`
    PageSize   int       `json:"page_size" default:"50"`
}

type ListReceitasOutput struct {
    Receitas   []ReceitaDTO  `json:"receitas"`
    Total      int64         `json:"total"`
    Page       int           `json:"page"`
    PageSize   int           `json:"page_size"`
}

func (uc *ListReceitasUseCase) Execute(
    ctx context.Context, input ListReceitasInput) (*ListReceitasOutput, error) {
    
    // Validações
    if input.TenantID == "" {
        return nil, ErrTenantIDRequired
    }
    
    if input.DataInicio.After(input.DataFim) {
        return nil, ErrPeriodoInvalido
    }
    
    // Limitar page size
    if input.PageSize > 100 {
        input.PageSize = 100
    }
    
    // Query com paginação
    receitas, total, err := uc.repository.FindByTenantAndPeriod(
        ctx, input.TenantID, input.DataInicio, input.DataFim,
        input.CategoriaID, input.Page, input.PageSize,
    )
    
    if err != nil {
        return nil, err
    }
    
    // Mapear para DTO
    receitasDTO := make([]ReceitaDTO, len(receitas))
    for i, r := range receitas {
        receitasDTO[i] = mapReceitaToDTO(r)
    }
    
    return &ListReceitasOutput{
        Receitas: receitasDTO,
        Total:    total,
        Page:     input.Page,
        PageSize: input.PageSize,
    }, nil
}
```

### 3. GetFluxoDeCaixaUseCase

```go
type GetFluxoDeCaixaUseCase struct {
    receitaRepo domain.ReceitaRepository
    despesaRepo domain.DespesaRepository
}

type GetFluxoDeCaixaInput struct {
    TenantID   string
    DataInicio time.Time
    DataFim    time.Time
}

func (uc *GetFluxoDeCaixaUseCase) Execute(
    ctx context.Context, input GetFluxoDeCaixaInput) (*domain.FluxoDeCaixa, error) {
    
    // 1. Calcular entradas (receitas confirmadas)
    entradas, err := uc.receitaRepo.SumByTenantAndPeriod(
        ctx, input.TenantID, input.DataInicio, input.DataFim,
        domain.ReceiptConfirmed,
    )
    if err != nil {
        return nil, err
    }
    
    // 2. Calcular saídas (despesas pagas)
    saidas, err := uc.despesaRepo.SumByTenantAndPeriod(
        ctx, input.TenantID, input.DataInicio, input.DataFim,
        domain.ExpensePaid,
    )
    if err != nil {
        return nil, err
    }
    
    // 3. Buscar saldo do dia anterior
    saldoAnterior := decimal.Zero // TODO: implementar
    
    // 4. Calcular saldo final
    fluxo := &domain.FluxoDeCaixa{
        TenantID:     input.TenantID,
        Periodo:      domain.Periodo{input.DataInicio, input.DataFim},
        SaldoInicial: saldoAnterior,
        Entradas:     entradas,
        Saidas:       saidas,
    }
    
    if err := fluxo.Calcular(); err != nil {
        return nil, err
    }
    
    return fluxo, nil
}
```

---

## 📐 Regras de Negócio

| Regra | Descrição | Validação |
|-------|-----------|-----------|
| **RN-FIN-001** | Receita deve ter valor > 0 | `valor > 0` |
| **RN-FIN-002** | Despesa deve ter valor > 0 | `valor > 0` |
| **RN-FIN-003** | Receita/Despesa devem ter categoria válida | `categoria_id NOT NULL` |
| **RN-FIN-004** | Receita cancelada não pode ser reativada | `status = CANCELADO` é irreversível |
| **RN-FIN-005** | Receita não pode ter data futura (default: hoje) | `data <= TODAY()` |
| **RN-FIN-006** | Cada categoria pertence a um único tenant | `categoria.tenant_id = receita.tenant_id` |
| **RN-FIN-007** | Fluxo de caixa só conta receitas "RECEBIDO" | `status = "RECEBIDO"` |
| **RN-FIN-008** | Fluxo de caixa só conta despesas "PAGO" | `status = "PAGO"` |

---

## 🗄️ Schema do Banco

### Tabela: receitas

```sql
CREATE TABLE receitas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(15, 2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'CONFIRMADO',
    observacoes TEXT,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data DESC);
CREATE INDEX idx_receitas_tenant_categoria ON receitas(tenant_id, categoria_id);
CREATE INDEX idx_receitas_tenant_status ON receitas(tenant_id, status);
CREATE INDEX idx_receitas_data ON receitas(data);
```

### Tabela: despesas

```sql
CREATE TABLE despesas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(15, 2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    fornecedor VARCHAR(255),
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'PENDENTE',
    observacoes TEXT,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_despesas_tenant_id ON despesas(tenant_id);
CREATE INDEX idx_despesas_tenant_data ON despesas(tenant_id, data DESC);
CREATE INDEX idx_despesas_tenant_categoria ON despesas(tenant_id, categoria_id);
CREATE INDEX idx_despesas_tenant_status ON despesas(tenant_id, status);
```

### Tabela: categorias

```sql
CREATE TABLE categorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('RECEITA', 'DESPESA')),
    cor VARCHAR(7) DEFAULT '#000000',
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, nome)
);

CREATE INDEX idx_categorias_tenant_id ON categorias(tenant_id);
CREATE INDEX idx_categorias_tenant_tipo ON categorias(tenant_id, tipo);
```

---

## 🔄 Fluxo de Dados

### Criar Receita

```
POST /financial/receitas
{
  "descricao": "Corte de cabelo",
  "valor": "50.00",
  "categoria_id": "cat-123",
  "metodo_pagamento": "PIX",
  "data": "2024-11-14"
}
        ↓
Handler valida request
        ↓
CreateReceitaUseCase.Execute()
  - Converter valor para Decimal
  - Criar entidade Receita
  - Validar regras de negócio
  - Chamar repository.Save()
        ↓
Repository (PostgreSQL)
  - INSERT com tenant_id obrigatório
  - Retornar receita criada
        ↓
Response 201 Created
{
  "id": "uuid-123",
  "status": "CONFIRMADO",
  "criado_em": "2024-11-14T10:30:00Z"
}
```

### Listar Receitas com Filtros

```
GET /financial/receitas?from=2024-11-01&to=2024-11-30&categoria_id=cat-123
        ↓
Handler parseia query params
        ↓
ListReceitasUseCase.Execute()
  - Validar período
  - Chamar repository.FindByTenantAndPeriod()
  - Mapear para DTOs
        ↓
Repository (PostgreSQL)
  SELECT * FROM receitas
  WHERE tenant_id = $1
    AND data BETWEEN $2 AND $3
    AND categoria_id = $4
  ORDER BY data DESC
  LIMIT 50 OFFSET 0
        ↓
Response 200 OK
{
  "receitas": [...],
  "total": 125,
  "page": 1,
  "page_size": 50
}
```

---

**Status:** ✅ Design finalizado  
**Próxima fase:** Implementação em Go
