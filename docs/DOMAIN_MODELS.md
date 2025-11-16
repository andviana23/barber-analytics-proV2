# 🧬 Domain Models

**Versão:** 2.0  
**Data:** 14/11/2025  
**Status:** Design

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Tenant (Bounded Context Root)](#tenant-bounded-context-root)
3. [User (Bounded Context Auth)](#user-bounded-context-auth)
4. [Receita (Bounded Context Financeiro)](#receita-bounded-context-financeiro)
5. [Despesa (Bounded Context Financeiro)](#despesa-bounded-context-financeiro)
6. [Assinatura (Bounded Context Subscriptions)](#assinatura-bounded-context-subscriptions)
7. [Value Objects](#value-objects)
8. [Enums](#enums)

---

## 🎯 Visão Geral

Modelos de domínio separados por **Bounded Context** (DDD).

```
┌─────────────────────────────────────────┐
│        Root Domain (Sistema)            │
└─────────────────────────────────────────┘
  │
  ├─ Bounded Context: Tenant Management
  │  └─ Agregates: Tenant
  │
  ├─ Bounded Context: Auth & Identity
  │  └─ Agregates: User, Role
  │
  ├─ Bounded Context: Financial
  │  ├─ Agregates: Receita
  │  ├─ Agregates: Despesa
  │  ├─ Agregates: Categoria
  │  └─ Services: FluxoDeCaixa
  │
  ├─ Bounded Context: Subscriptions
  │  ├─ Agregates: PlanoDeassinatura
  │  ├─ Agregates: Assinatura
  │  └─ Agregates: AssinaturaInvoice
  │
  └─ Bounded Context: Inventory (Futuro)
     ├─ Agregates: Produto
     ├─ Agregates: Movimentacao
     └─ Agregates: Fornecedor
```

---

## 👥 Tenant (Bounded Context Root)

**Arquivo:** `internal/domain/tenant/model.go`

```go
package tenant

import (
    \"time\"
    \"github.com/google/uuid\"
)

// Tenant é o Aggregate Root do contexto de tenant
type Tenant struct {
    id          string
    name        string
    cnpj        string
    active      bool
    plan        Plan
    createdAt   time.Time
    updatedAt   time.Time
}

// Value Objects
type Plan string

const (
    PlanFree       Plan = \"FREE\"
    PlanPro        Plan = \"PRO\"
    PlanEnterprise Plan = \"ENTERPRISE\"
)

// Métodos do Agregado
func (t *Tenant) ID() string {
    return t.id
}

func (t *Tenant) IsActive() bool {
    return t.active
}

func (t *Tenant) Deactivate() {
    t.active = false
    t.updatedAt = time.Now()
}

func (t *Tenant) Activate() {
    t.active = true
    t.updatedAt = time.Now()
}

// Factory
func NewTenant(name, cnpj string) (*Tenant, error) {
    if name == \"\" {
        return nil, errors.New(\"name required\")
    }
    
    return &Tenant{
        id:        uuid.NewString(),
        name:      name,
        cnpj:      cnpj,
        active:    true,
        plan:      PlanFree,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }, nil
}

// Repository Interface (Port)
type Repository interface {
    Save(ctx context.Context, tenant *Tenant) error
    FindByID(ctx context.Context, id string) (*Tenant, error)
    FindByName(ctx context.Context, name string) (*Tenant, error)
}
```

---

## 👤 User (Bounded Context Auth)

**Arquivo:** `internal/domain/user/model.go`

```go
package user

import (
    \"time\"
    \"golang.org/x/crypto/bcrypt\"
)

type User struct {
    id           string
    tenantID     string
    email        string
    passwordHash string
    name         string
    role         Role
    active       bool
    lastLogin    *time.Time
    createdAt    time.Time
    updatedAt    time.Time
}

type Role string

const (
    RoleOwner      Role = \"OWNER\"      // Dono da barbearia
    RoleManager    Role = \"MANAGER\"    // Gerente
    RoleAccountant Role = \"ACCOUNTANT\" // Contador/Financeiro
    RoleEmployee   Role = \"EMPLOYEE\"   // Barbeiro/Funcionário
)

// Value Object: Email
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, errors.New(\"invalid email\")
    }
    return &Email{value: value}, nil
}

func (e *Email) String() string {
    return e.value
}

// Factory
func NewUser(tenantID, email, password, name string, role Role) (*User, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    em, err := NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    return &User{
        id:           uuid.NewString(),
        tenantID:     tenantID,
        email:        em.String(),
        passwordHash: string(hash),
        name:         name,
        role:         role,
        active:       true,
        createdAt:    time.Now(),
        updatedAt:    time.Now(),
    }, nil
}

// Verificar senha
func (u *User) VerifyPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
    return err == nil
}

// Repository Interface
type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByTenantAndEmail(ctx context.Context, tenantID, email string) (*User, error)
}
```

---

## 💰 Receita (Bounded Context Financeiro)

**Arquivo:** `internal/domain/financial/receita.go`

```go
package financial

import (
    \"github.com/shopspring/decimal\"
    \"time\"
)

// Receita é o Aggregate Root
type Receita struct {
    id              string
    tenantID        string
    userID          string
    description     string
    amount          decimal.Decimal
    category        *Categoria
    paymentMethod   PaymentMethod
    date            time.Time
    status          ReceiptStatus
    observations    string
    createdAt       time.Time
    updatedAt       time.Time
}

type ReceiptStatus string

const (
    ReceiptConfirmed ReceiptStatus = \"CONFIRMADO\"
    ReceiptReceived  ReceiptStatus = \"RECEBIDO\"
    ReceiptCancelled ReceiptStatus = \"CANCELADO\"
)

// Value Object: Money
type Money struct {
    amount   decimal.Decimal
    currency string // \"BRL\"
}

func NewMoney(value string) (*Money, error) {
    amount, err := decimal.NewFromString(value)
    if err != nil {
        return nil, err
    }
    if amount.LessThanOrEqual(decimal.Zero) {
        return nil, errors.New(\"amount must be positive\")
    }
    return &Money{amount: amount, currency: \"BRL\"}, nil
}

// Factory
func NewReceita(
    tenantID, userID, description, amount string,
    categoryID, paymentMethod string,
    date time.Time) (*Receita, error) {
    
    money, err := NewMoney(amount)
    if err != nil {
        return nil, err
    }
    
    return &Receita{
        id:            uuid.NewString(),
        tenantID:      tenantID,
        userID:        userID,
        description:   description,
        amount:        money.amount,
        paymentMethod: PaymentMethod(paymentMethod),
        date:          date,
        status:        ReceiptConfirmed,
        createdAt:     time.Now(),
        updatedAt:     time.Now(),
    }, nil
}

// Métodos do Agregado
func (r *Receita) Confirm() {
    r.status = ReceiptConfirmed
    r.updatedAt = time.Now()
}

func (r *Receita) Cancel() error {
    if r.status == ReceiptCancelled {
        return errors.New(\"receipt already cancelled\")
    }
    r.status = ReceiptCancelled
    r.updatedAt = time.Now()
    return nil
}

// Repository Interface
type ReceitaRepository interface {
    Save(ctx context.Context, tenantID string, receita *Receita) error
    FindByID(ctx context.Context, tenantID, id string) (*Receita, error)
    FindByTenantAndPeriod(ctx context.Context, tenantID string,
        from, to time.Time, opts FindOptions) ([]*Receita, error)
    Update(ctx context.Context, tenantID string, receita *Receita) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

---

## 💸 Despesa (Bounded Context Financeiro)

Similar a Receita, com adicionais:

```go
type Despesa struct {
    // ... campos similares a Receita
    supplier    string // Fornecedor
    status      ExpenseStatus
}

type ExpenseStatus string

const (
    ExpensePending   ExpenseStatus = \"PENDENTE\"
    ExpensePaid      ExpenseStatus = \"PAGO\"
    ExpenseCancelled ExpenseStatus = \"CANCELADO\"
)

func (d *Despesa) MarkAsPaid() error {
    if d.status == ExpensePaid {
        return errors.New(\"expense already paid\")
    }
    d.status = ExpensePaid
    d.updatedAt = time.Now()
    return nil
}
```

---

## 🎟️ Assinatura (Bounded Context Subscriptions)

**Arquivo:** `internal/domain/subscription/model.go`

```go
package subscription

type Assinatura struct {
    id                  string
    tenantID            string
    planID              string
    barbeiroID          string
    asaasSubscriptionID string
    status              SubscriptionStatus
    startDate           time.Time
    endDate             *time.Time
    nextInvoiceDate     time.Time
    createdAt           time.Time
    updatedAt           time.Time
}

type SubscriptionStatus string

const (
    SubActive     SubscriptionStatus = \"ATIVA\"
    SubCancelled  SubscriptionStatus = \"CANCELADA\"
    SubSuspended  SubscriptionStatus = \"SUSPENSA\"
    SubExpired    SubscriptionStatus = \"EXPIRADA\"
)

// Factory
func NewAssinatura(
    tenantID, planID, barbeiroID, asaasID string,
    startDate time.Time) (*Assinatura, error) {
    
    return &Assinatura{
        id:                  uuid.NewString(),
        tenantID:            tenantID,
        planID:              planID,
        barbeiroID:          barbeiroID,
        asaasSubscriptionID: asaasID,
        status:              SubActive,
        startDate:           startDate,
        nextInvoiceDate:     startDate.AddDate(0, 1, 0),
        createdAt:           time.Now(),
        updatedAt:           time.Now(),
    }, nil
}

func (s *Assinatura) Cancel() error {
    if s.status == SubCancelled {
        return errors.New(\"subscription already cancelled\")
    }
    s.status = SubCancelled
    s.endDate = ptrTime(time.Now())
    s.updatedAt = time.Now()
    return nil
}

// Repository Interface
type Repository interface {
    Save(ctx context.Context, tenantID string, sub *Assinatura) error
    FindByID(ctx context.Context, tenantID, id string) (*Assinatura, error)
    FindByTenant(ctx context.Context, tenantID string,
        status SubscriptionStatus) ([]*Assinatura, error)
}
```

---

## 💎 Value Objects

Value Objects são imutáveis e sem identidade própria.

```go
// Endereco (Value Object)
type Endereco struct {
    Rua          string
    Numero       int
    Complemento  string
    Cidade       string
    UF           string
    CEP          string
}

// Periodo (Value Object)
type Periodo struct {
    DataInicio time.Time
    DataFim    time.Time
}

func (p *Periodo) IsValid() bool {
    return p.DataInicio.Before(p.DataFim)
}

func (p *Periodo) Days() int {
    return int(p.DataFim.Sub(p.DataInicio).Hours() / 24)
}
```

---

## 🏷️ Enums

```go
// Payment Methods
type PaymentMethod string

const (
    PaymentCash       PaymentMethod = \"DINHEIRO\"
    PaymentDebit      PaymentMethod = \"DEBITO\"
    PaymentCredit     PaymentMethod = \"CREDITO\"
    PaymentPix        PaymentMethod = \"PIX\"
    PaymentTransfer   PaymentMethod = \"TRANSFERENCIA\"
)

// Periodicidades
type Periodicidade string

const (
    PerioMensal    Periodicidade = \"MENSAL\"
    PerioTrimestral Periodicidade = \"TRIMESTRAL\"
    PerioAnual     Periodicidade = \"ANUAL\"
)
```

---

**Status:** ✅ Completo
