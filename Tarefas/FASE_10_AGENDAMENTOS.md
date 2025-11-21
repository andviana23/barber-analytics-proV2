# 📅 FASE 10 - Módulo de Agendamentos

**Versão:** 1.0
**Data:** 20/11/2025
**Status:** Planejamento
**Prioridade:** Alta (Core Business)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Backend - Core](#backend---core)
3. [Integração DayPilot](#integração-daypilot)
4. [Frontend - UI](#frontend---ui)
5. [Notificações](#notificações)
6. [Testes](#testes)

---

## 🎯 Visão Geral

Implementar sistema completo de agendamentos com:

- ✅ CRUD de agendamentos
- ✅ Integração com DayPilot Scheduler
- ✅ Sincronização tempo real
- ✅ Notificações (SMS/WhatsApp via Asaas)
- ✅ Gestão de disponibilidade de profissionais
- ✅ Regras de negócio (horários, conflitos)

**Dependências:**

- ✅ Módulo Cadastro (clientes, profissionais, serviços)
- ✅ Design System (tokens, tema DayPilot)
- ⏳ Módulo Financeiro (para checkout pós-atendimento)

---

## 🔧 Backend - Core

### T-BE-AGENDA-001: Entidades de Domínio

**Arquivo:** `backend/internal/domain/entity/agendamento.go`

```go
package entity

import (
    "time"
    "github.com/shopspring/decimal"
)

type Agendamento struct {
    ID              string
    TenantID        string
    ClienteID       string
    ProfissionalID  string
    ServicoID       string
    DataHoraInicio  time.Time
    DataHoraFim     time.Time
    Status          AgendamentoStatus
    ValorTotal      decimal.Decimal
    Observacoes     string
    CriadoPor       string
    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

type AgendamentoStatus string

const (
    StatusAgendado    AgendamentoStatus = "AGENDADO"
    StatusConfirmado  AgendamentoStatus = "CONFIRMADO"
    StatusEmAndamento AgendamentoStatus = "EM_ANDAMENTO"
    StatusConcluido   AgendamentoStatus = "CONCLUIDO"
    StatusCancelado   AgendamentoStatus = "CANCELADO"
    StatusNaoCompareceu AgendamentoStatus = "NAO_COMPARECEU"
)

// Validações de domínio
func (a *Agendamento) Validate() error {
    if a.ID == "" {
        return ErrAgendamentoIDRequired
    }
    if a.TenantID == "" {
        return ErrTenantIDRequired
    }
    if a.ClienteID == "" {
        return ErrClienteIDRequired
    }
    if a.ProfissionalID == "" {
        return ErrProfissionalIDRequired
    }
    if a.ServicoID == "" {
        return ErrServicoIDRequired
    }
    if a.DataHoraInicio.IsZero() {
        return ErrDataHoraInicioRequired
    }
    if a.DataHoraFim.IsZero() {
        return ErrDataHoraFimRequired
    }
    if a.DataHoraFim.Before(a.DataHoraInicio) {
        return ErrDataHoraFimAnteriorInicio
    }
    if a.ValorTotal.LessThan(decimal.Zero) {
        return ErrValorTotalInvalido
    }
    return nil
}

// Regras de negócio
func (a *Agendamento) Confirmar() error {
    if a.Status != StatusAgendado {
        return ErrStatusInvalidoParaConfirmacao
    }
    a.Status = StatusConfirmado
    a.AtualizadoEm = time.Now()
    return nil
}

func (a *Agendamento) IniciarAtendimento() error {
    if a.Status != StatusConfirmado && a.Status != StatusAgendado {
        return ErrStatusInvalidoParaIniciar
    }
    a.Status = StatusEmAndamento
    a.AtualizadoEm = time.Now()
    return nil
}

func (a *Agendamento) Concluir() error {
    if a.Status != StatusEmAndamento {
        return ErrStatusInvalidoParaConcluir
    }
    a.Status = StatusConcluido
    a.AtualizadoEm = time.Now()
    return nil
}

func (a *Agendamento) Cancelar(motivo string) error {
    if a.Status == StatusConcluido || a.Status == StatusCancelado {
        return ErrNaoPodeCancelar
    }
    a.Status = StatusCancelado
    a.Observacoes += fmt.Sprintf("\nCANCELADO: %s", motivo)
    a.AtualizadoEm = time.Now()
    return nil
}

func (a *Agendamento) Duracao() time.Duration {
    return a.DataHoraFim.Sub(a.DataHoraInicio)
}

func (a *Agendamento) EstaNoPassado() bool {
    return a.DataHoraInicio.Before(time.Now())
}

func (a *Agendamento) PermiteEdicao() bool {
    return a.Status == StatusAgendado || a.Status == StatusConfirmado
}
```

**Checklist:**

- [ ] Criar entidade Agendamento
- [ ] Implementar métodos de validação
- [ ] Implementar transições de estado
- [ ] Criar testes unitários (>90% coverage)
- [ ] Documentar regras de negócio

---

### T-BE-AGENDA-002: Repository Interface

**Arquivo:** `backend/internal/domain/repository/agendamento_repository.go`

```go
package repository

import (
    "context"
    "time"
    "github.com/andviana23/barber-analytics/internal/domain/entity"
)

type AgendamentoRepository interface {
    // CRUD básico
    Save(ctx context.Context, agendamento *entity.Agendamento) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.Agendamento, error)
    Update(ctx context.Context, agendamento *entity.Agendamento) error
    Delete(ctx context.Context, tenantID, id string) error

    // Queries específicas
    FindByTenantAndPeriod(
        ctx context.Context,
        tenantID string,
        inicio, fim time.Time,
    ) ([]*entity.Agendamento, error)

    FindByProfissionalAndPeriod(
        ctx context.Context,
        tenantID, profissionalID string,
        inicio, fim time.Time,
    ) ([]*entity.Agendamento, error)

    FindByClienteAndPeriod(
        ctx context.Context,
        tenantID, clienteID string,
        inicio, fim time.Time,
    ) ([]*entity.Agendamento, error)

    // Validações de conflito
    VerificarConflito(
        ctx context.Context,
        tenantID, profissionalID string,
        inicio, fim time.Time,
        excludeID string, // Para atualização
    ) (bool, error)

    // Contadores
    CountByProfissionalAndStatus(
        ctx context.Context,
        tenantID, profissionalID string,
        status entity.AgendamentoStatus,
    ) (int64, error)
}
```

**Checklist:**

- [ ] Definir interface do repository
- [ ] Documentar cada método
- [ ] Criar mock para testes

---

### T-BE-AGENDA-003: Use Cases

**Arquivo:** `backend/internal/application/usecase/agendamento/create_agendamento.go`

```go
package agendamento

type CreateAgendamentoUseCase struct {
    agendamentoRepo   repository.AgendamentoRepository
    clienteRepo       repository.ClienteRepository
    profissionalRepo  repository.ProfissionalRepository
    servicoRepo       repository.ServicoRepository
    validator         *validator.Validate
}

type CreateAgendamentoInput struct {
    TenantID        string    `json:"tenant_id" validate:"required,uuid"`
    ClienteID       string    `json:"cliente_id" validate:"required,uuid"`
    ProfissionalID  string    `json:"profissional_id" validate:"required,uuid"`
    ServicoID       string    `json:"servico_id" validate:"required,uuid"`
    DataHoraInicio  time.Time `json:"data_hora_inicio" validate:"required"`
    Observacoes     string    `json:"observacoes" validate:"max=500"`
    CriadoPor       string    `json:"-"` // User ID do contexto
}

type CreateAgendamentoOutput struct {
    ID             string    `json:"id"`
    ClienteID      string    `json:"cliente_id"`
    ClienteNome    string    `json:"cliente_nome"`
    ProfissionalID string    `json:"profissional_id"`
    ProfissionalNome string  `json:"profissional_nome"`
    ServicoID      string    `json:"servico_id"`
    ServicoNome    string    `json:"servico_nome"`
    DataHoraInicio time.Time `json:"data_hora_inicio"`
    DataHoraFim    time.Time `json:"data_hora_fim"`
    Status         string    `json:"status"`
    ValorTotal     string    `json:"valor_total"`
    CriadoEm       time.Time `json:"criado_em"`
}

func (uc *CreateAgendamentoUseCase) Execute(
    ctx context.Context,
    input CreateAgendamentoInput,
) (*CreateAgendamentoOutput, error) {
    // 1. Validar input
    if err := uc.validator.Struct(input); err != nil {
        return nil, ErrInvalidInput
    }

    // 2. Validar que cliente existe e está ativo
    cliente, err := uc.clienteRepo.FindByID(ctx, input.TenantID, input.ClienteID)
    if err != nil {
        return nil, ErrClienteNotFound
    }
    if !cliente.Ativo {
        return nil, ErrClienteInativo
    }

    // 3. Validar que profissional existe e está ativo
    profissional, err := uc.profissionalRepo.FindByID(ctx, input.TenantID, input.ProfissionalID)
    if err != nil {
        return nil, ErrProfissionalNotFound
    }
    if profissional.Status != "ATIVO" {
        return nil, ErrProfissionalInativo
    }

    // 4. Buscar serviço e calcular fim
    servico, err := uc.servicoRepo.FindByID(ctx, input.TenantID, input.ServicoID)
    if err != nil {
        return nil, ErrServicoNotFound
    }
    if !servico.Ativo {
        return nil, ErrServicoInativo
    }

    dataHoraFim := input.DataHoraInicio.Add(time.Duration(servico.Duracao) * time.Minute)

    // 5. Verificar conflito de horário
    conflito, err := uc.agendamentoRepo.VerificarConflito(
        ctx,
        input.TenantID,
        input.ProfissionalID,
        input.DataHoraInicio,
        dataHoraFim,
        "", // Novo agendamento, sem ID para excluir
    )
    if err != nil {
        return nil, err
    }
    if conflito {
        return nil, ErrHorarioConflito
    }

    // 6. Validar horário de trabalho do profissional
    if !uc.validarHorarioTrabalho(profissional, input.DataHoraInicio, dataHoraFim) {
        return nil, ErrForaHorarioTrabalho
    }

    // 7. Criar entidade
    agendamento := &entity.Agendamento{
        ID:             uuid.New().String(),
        TenantID:       input.TenantID,
        ClienteID:      input.ClienteID,
        ProfissionalID: input.ProfissionalID,
        ServicoID:      input.ServicoID,
        DataHoraInicio: input.DataHoraInicio,
        DataHoraFim:    dataHoraFim,
        Status:         entity.StatusAgendado,
        ValorTotal:     servico.Preco,
        Observacoes:    input.Observacoes,
        CriadoPor:      input.CriadoPor,
        CriadoEm:       time.Now(),
        AtualizadoEm:   time.Now(),
    }

    // 8. Validar domínio
    if err := agendamento.Validate(); err != nil {
        return nil, err
    }

    // 9. Persistir
    if err := uc.agendamentoRepo.Save(ctx, agendamento); err != nil {
        return nil, fmt.Errorf("failed to save agendamento: %w", err)
    }

    // 10. Retornar output
    return &CreateAgendamentoOutput{
        ID:               agendamento.ID,
        ClienteID:        agendamento.ClienteID,
        ClienteNome:      cliente.Nome,
        ProfissionalID:   agendamento.ProfissionalID,
        ProfissionalNome: profissional.Nome,
        ServicoID:        agendamento.ServicoID,
        ServicoNome:      servico.Nome,
        DataHoraInicio:   agendamento.DataHoraInicio,
        DataHoraFim:      agendamento.DataHoraFim,
        Status:           string(agendamento.Status),
        ValorTotal:       agendamento.ValorTotal.String(),
        CriadoEm:         agendamento.CriadoEm,
    }, nil
}

func (uc *CreateAgendamentoUseCase) validarHorarioTrabalho(
    profissional *entity.Profissional,
    inicio, fim time.Time,
) bool {
    // TODO: Implementar validação de horário de trabalho
    // Verificar se o horário está dentro do expediente do profissional
    return true
}
```

**Outros Use Cases:**

- `UpdateAgendamentoUseCase`
- `CancelAgendamentoUseCase`
- `ConfirmarAgendamentoUseCase`
- `IniciarAtendimentoUseCase`
- `ConcluirAtendimentoUseCase`
- `ListAgendamentosPorPeriodoUseCase`
- `ListAgendamentosPorProfissionalUseCase`

**Checklist:**

- [ ] Implementar CreateAgendamentoUseCase
- [ ] Implementar UpdateAgendamentoUseCase
- [ ] Implementar CancelAgendamentoUseCase
- [ ] Implementar ConfirmarAgendamentoUseCase
- [ ] Implementar transição de status (iniciar, concluir)
- [ ] Implementar listagem por período
- [ ] Criar testes unitários (>85% coverage)

---

### T-BE-AGENDA-004: Repository PostgreSQL

**Arquivo:** `backend/internal/infrastructure/repository/postgres/agendamento_repository.go`

```go
package postgres

type PostgresAgendamentoRepository struct {
    db *sql.DB
}

func (r *PostgresAgendamentoRepository) Save(
    ctx context.Context,
    agendamento *entity.Agendamento,
) error {
    query := `
        INSERT INTO agendamentos (
            id, tenant_id, cliente_id, profissional_id, servico_id,
            data_hora_inicio, data_hora_fim, status, valor_total,
            observacoes, criado_por, criado_em, atualizado_em
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
    _, err := r.db.ExecContext(
        ctx, query,
        agendamento.ID,
        agendamento.TenantID,
        agendamento.ClienteID,
        agendamento.ProfissionalID,
        agendamento.ServicoID,
        agendamento.DataHoraInicio,
        agendamento.DataHoraFim,
        agendamento.Status,
        agendamento.ValorTotal,
        agendamento.Observacoes,
        agendamento.CriadoPor,
        agendamento.CriadoEm,
        agendamento.AtualizadoEm,
    )
    return err
}

func (r *PostgresAgendamentoRepository) VerificarConflito(
    ctx context.Context,
    tenantID, profissionalID string,
    inicio, fim time.Time,
    excludeID string,
) (bool, error) {
    query := `
        SELECT EXISTS (
            SELECT 1
            FROM agendamentos
            WHERE tenant_id = $1
              AND profissional_id = $2
              AND status NOT IN ('CANCELADO', 'CONCLUIDO', 'NAO_COMPARECEU')
              AND id != $3
              AND (
                  (data_hora_inicio < $5 AND data_hora_fim > $4)
                  OR (data_hora_inicio >= $4 AND data_hora_inicio < $5)
              )
        )
    `
    var exists bool
    err := r.db.QueryRowContext(
        ctx, query,
        tenantID, profissionalID, excludeID, inicio, fim,
    ).Scan(&exists)

    return exists, err
}
```

**Checklist:**

- [ ] Implementar Save
- [ ] Implementar FindByID
- [ ] Implementar Update
- [ ] Implementar Delete (soft delete)
- [ ] Implementar FindByTenantAndPeriod
- [ ] Implementar VerificarConflito
- [ ] Criar índices adequados
- [ ] Testar com dados reais

---

### T-BE-AGENDA-005: Migrations

**Arquivo:** `backend/migrations/XXX_create_agendamentos.up.sql`

```sql
-- Tabela principal de agendamentos
CREATE TABLE agendamentos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    cliente_id UUID NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE RESTRICT,
    servico_id UUID NOT NULL REFERENCES servicos(id) ON DELETE RESTRICT,
    data_hora_inicio TIMESTAMPTZ NOT NULL,
    data_hora_fim TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'AGENDADO'
        CHECK (status IN ('AGENDADO', 'CONFIRMADO', 'EM_ANDAMENTO', 'CONCLUIDO', 'CANCELADO', 'NAO_COMPARECEU')),
    valor_total DECIMAL(10, 2) NOT NULL CHECK (valor_total >= 0),
    observacoes TEXT,
    criado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_period CHECK (data_hora_fim > data_hora_inicio)
);

-- Índices para performance
CREATE INDEX idx_agendamentos_tenant_id ON agendamentos(tenant_id);
CREATE INDEX idx_agendamentos_cliente_id ON agendamentos(tenant_id, cliente_id);
CREATE INDEX idx_agendamentos_profissional_id ON agendamentos(tenant_id, profissional_id);
CREATE INDEX idx_agendamentos_periodo ON agendamentos(tenant_id, data_hora_inicio, data_hora_fim);
CREATE INDEX idx_agendamentos_status ON agendamentos(tenant_id, status);
CREATE INDEX idx_agendamentos_data_status ON agendamentos(tenant_id, data_hora_inicio, status);

-- Índice para detecção de conflitos (CRÍTICO)
CREATE INDEX idx_agendamentos_conflito ON agendamentos(
    tenant_id,
    profissional_id,
    data_hora_inicio,
    data_hora_fim
) WHERE status NOT IN ('CANCELADO', 'CONCLUIDO', 'NAO_COMPARECEU');

-- Trigger para atualizar timestamp
CREATE TRIGGER update_agendamentos_timestamp
    BEFORE UPDATE ON agendamentos
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- Comentários
COMMENT ON TABLE agendamentos IS 'Registro de agendamentos de serviços';
COMMENT ON COLUMN agendamentos.status IS 'AGENDADO: criado | CONFIRMADO: confirmado pelo cliente | EM_ANDAMENTO: atendimento iniciado | CONCLUIDO: finalizado | CANCELADO: cancelado | NAO_COMPARECEU: falta';
COMMENT ON INDEX idx_agendamentos_conflito IS 'Índice parcial para detectar conflitos de horário rapidamente';
```

**Checklist:**

- [ ] Criar tabela agendamentos
- [ ] Criar índices de performance
- [ ] Criar índice de conflito (parcial)
- [ ] Criar triggers de timestamp
- [ ] Criar migration down
- [ ] Testar migration up/down

---

### T-BE-AGENDA-006: HTTP Handlers

**Arquivo:** `backend/internal/infrastructure/http/handler/agendamento_handler.go`

```go
package handler

type AgendamentoHandler struct {
    createUC   *usecase.CreateAgendamentoUseCase
    updateUC   *usecase.UpdateAgendamentoUseCase
    cancelUC   *usecase.CancelAgendamentoUseCase
    confirmarUC *usecase.ConfirmarAgendamentoUseCase
    listUC     *usecase.ListAgendamentosUseCase
}

func (h *AgendamentoHandler) Create(c echo.Context) error {
    var req dto.CreateAgendamentoRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, dto.ErrorResponse{Message: "Invalid request"})
    }

    tenantID := c.Get("tenant_id").(string)
    userID := c.Get("user_id").(string)

    input := usecase.CreateAgendamentoInput{
        TenantID:       tenantID,
        ClienteID:      req.ClienteID,
        ProfissionalID: req.ProfissionalID,
        ServicoID:      req.ServicoID,
        DataHoraInicio: req.DataHoraInicio,
        Observacoes:    req.Observacoes,
        CriadoPor:      userID,
    }

    output, err := h.createUC.Execute(c.Request().Context(), input)
    if err != nil {
        // Map domain errors to HTTP status codes
        switch {
        case errors.Is(err, usecase.ErrHorarioConflito):
            return c.JSON(409, dto.ErrorResponse{
                Code:    "CONFLICT",
                Message: "Horário já ocupado para este profissional",
            })
        case errors.Is(err, usecase.ErrClienteNotFound):
            return c.JSON(404, dto.ErrorResponse{
                Code:    "NOT_FOUND",
                Message: "Cliente não encontrado",
            })
        default:
            return c.JSON(500, dto.ErrorResponse{Message: "Internal error"})
        }
    }

    return c.JSON(201, output)
}

func (h *AgendamentoHandler) List(c echo.Context) error {
    tenantID := c.Get("tenant_id").(string)

    // Parse query params
    dataInicio, _ := time.Parse("2006-01-02", c.QueryParam("data_inicio"))
    dataFim, _ := time.Parse("2006-01-02", c.QueryParam("data_fim"))
    profissionalID := c.QueryParam("profissional_id")
    clienteID := c.QueryParam("cliente_id")
    status := c.QueryParam("status")

    input := usecase.ListAgendamentosInput{
        TenantID:       tenantID,
        DataInicio:     dataInicio,
        DataFim:        dataFim,
        ProfissionalID: profissionalID,
        ClienteID:      clienteID,
        Status:         status,
    }

    output, err := h.listUC.Execute(c.Request().Context(), input)
    if err != nil {
        return c.JSON(500, dto.ErrorResponse{Message: "Internal error"})
    }

    return c.JSON(200, output)
}
```

**Endpoints:**

```
POST   /api/v1/agendamentos
GET    /api/v1/agendamentos
GET    /api/v1/agendamentos/:id
PUT    /api/v1/agendamentos/:id
DELETE /api/v1/agendamentos/:id
POST   /api/v1/agendamentos/:id/confirmar
POST   /api/v1/agendamentos/:id/iniciar
POST   /api/v1/agendamentos/:id/concluir
POST   /api/v1/agendamentos/:id/cancelar
```

**Checklist:**

- [ ] Implementar Create handler
- [ ] Implementar List handler
- [ ] Implementar Get by ID handler
- [ ] Implementar Update handler
- [ ] Implementar Delete handler
- [ ] Implementar Confirmar handler
- [ ] Implementar Iniciar handler
- [ ] Implementar Concluir handler
- [ ] Implementar Cancelar handler
- [ ] Mapear domain errors para HTTP status
- [ ] Adicionar validação de permissões (RBAC)
- [ ] Criar testes de integração

---

## 📅 Integração DayPilot

### T-FE-AGENDA-001: Adapter DayPilot

**Arquivo:** `frontend/app/lib/adapters/daypilotAdapter.ts`

```typescript
import { Agendamento } from "@/app/lib/types/agendamento";
import { DayPilot } from "daypilot-pro-react";
import { tokens } from "@/app/theme/tokens";

/**
 * Mapeia agendamento do backend para evento DayPilot
 */
export function agendamentoToDayPilotEvent(
  agendamento: Agendamento
): DayPilot.EventData {
  return {
    id: agendamento.id,
    start: new DayPilot.Date(agendamento.data_hora_inicio),
    end: new DayPilot.Date(agendamento.data_hora_fim),
    resource: agendamento.profissional_id,
    text: `${agendamento.cliente_nome} - ${agendamento.servico_nome}`,
    backColor: getStatusColor(agendamento.status),
    borderColor: getStatusBorderColor(agendamento.status),
    fontColor: "#FFFFFF",
    tags: {
      agendamento_id: agendamento.id,
      cliente_id: agendamento.cliente_id,
      servico_id: agendamento.servico_id,
      status: agendamento.status,
      valor_total: agendamento.valor_total,
    },
    bubbleHtml: generateBubbleHtml(agendamento),
  };
}

/**
 * Mapeia profissional para resource DayPilot
 */
export function profissionalToDayPilotResource(
  profissional: Profissional
): DayPilot.ResourceData {
  return {
    id: profissional.id,
    name: profissional.nome,
    tags: {
      profissional_id: profissional.id,
      especialidades: profissional.especialidades,
    },
  };
}

/**
 * Cores por status
 */
function getStatusColor(status: string): string {
  const statusColors: Record<string, string> = {
    AGENDADO: tokens.colors.primary[500],
    CONFIRMADO: tokens.colors.success[500],
    EM_ANDAMENTO: tokens.colors.warning[500],
    CONCLUIDO: tokens.colors.neutral[400],
    CANCELADO: tokens.colors.error[500],
    NAO_COMPARECEU: tokens.colors.neutral[300],
  };
  return statusColors[status] || tokens.colors.primary[500];
}

function getStatusBorderColor(status: string): string {
  const statusColors: Record<string, string> = {
    AGENDADO: tokens.colors.primary[700],
    CONFIRMADO: tokens.colors.success[700],
    EM_ANDAMENTO: tokens.colors.warning[700],
    CONCLUIDO: tokens.colors.neutral[600],
    CANCELADO: tokens.colors.error[700],
    NAO_COMPARECEU: tokens.colors.neutral[500],
  };
  return statusColors[status] || tokens.colors.primary[700];
}

/**
 * Gera HTML para bubble (tooltip)
 */
function generateBubbleHtml(agendamento: Agendamento): string {
  return `
    <div style="padding: 12px; font-family: Inter, sans-serif;">
      <div style="font-weight: 600; font-size: 14px; margin-bottom: 8px;">
        ${agendamento.cliente_nome}
      </div>
      <div style="font-size: 13px; color: #64748B; margin-bottom: 4px;">
        <strong>Serviço:</strong> ${agendamento.servico_nome}
      </div>
      <div style="font-size: 13px; color: #64748B; margin-bottom: 4px;">
        <strong>Profissional:</strong> ${agendamento.profissional_nome}
      </div>
      <div style="font-size: 13px; color: #64748B; margin-bottom: 4px;">
        <strong>Horário:</strong> ${formatTime(
          agendamento.data_hora_inicio
        )} - ${formatTime(agendamento.data_hora_fim)}
      </div>
      <div style="font-size: 13px; color: #64748B;">
        <strong>Valor:</strong> R$ ${agendamento.valor_total}
      </div>
      ${
        agendamento.observacoes
          ? `<div style="margin-top: 8px; font-size: 12px; color: #94A3B8; font-style: italic;">
            ${agendamento.observacoes}
          </div>`
          : ""
      }
    </div>
  `;
}

function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString("pt-BR", {
    hour: "2-digit",
    minute: "2-digit",
  });
}
```

**Checklist:**

- [ ] Criar adapter agendamento → DayPilot event
- [ ] Criar adapter profissional → DayPilot resource
- [ ] Implementar cores por status
- [ ] Criar bubble HTML customizado
- [ ] Testar com dados reais

---

### T-FE-AGENDA-002: Componente Scheduler

**Arquivo:** `frontend/app/components/agendamentos/Scheduler.tsx`

```typescript
"use client";

import { useEffect, useRef, useState } from "react";
import { DayPilot, DayPilotScheduler } from "daypilot-pro-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useThemeStore } from "@/app/lib/store/themeStore";
import { getDayPilotTheme } from "@/app/theme/daypilotTheme";
import {
  agendamentoToDayPilotEvent,
  profissionalToDayPilotResource,
} from "@/app/lib/adapters/daypilotAdapter";
import { api } from "@/app/lib/api/client";

interface SchedulerProps {
  tenantId: string;
  dataInicio: Date;
  dataFim: Date;
}

export function Scheduler({ tenantId, dataInicio, dataFim }: SchedulerProps) {
  const schedulerRef = useRef<DayPilot.Scheduler>();
  const { theme } = useThemeStore();
  const queryClient = useQueryClient();

  // Fetch profissionais
  const { data: profissionais = [] } = useQuery({
    queryKey: ["profissionais", tenantId, { ativo: true }],
    queryFn: () => api.profissionais.list(tenantId, { ativo: true }),
  });

  // Fetch agendamentos
  const { data: agendamentos = [], isLoading } = useQuery({
    queryKey: ["agendamentos", tenantId, dataInicio, dataFim],
    queryFn: () =>
      api.agendamentos.list(tenantId, {
        data_inicio: dataInicio.toISOString(),
        data_fim: dataFim.toISOString(),
      }),
    staleTime: 30 * 1000, // 30 segundos
  });

  // Mutation para criar agendamento
  const createMutation = useMutation({
    mutationFn: (data: any) => api.agendamentos.create(tenantId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["agendamentos", tenantId] });
    },
  });

  // Mutation para atualizar agendamento
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: any) =>
      api.agendamentos.update(tenantId, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["agendamentos", tenantId] });
    },
  });

  // Aplicar tema DayPilot
  useEffect(() => {
    const daypilotTheme = getDayPilotTheme(theme);
    const styleEl = document.createElement("style");
    styleEl.innerHTML = daypilotTheme.css;
    document.head.appendChild(styleEl);

    return () => styleEl.remove();
  }, [theme]);

  // Converter dados para DayPilot
  const resources = profissionais.map(profissionalToDayPilotResource);
  const events = agendamentos.map(agendamentoToDayPilotEvent);

  // Event handlers
  const handleTimeRangeSelected = async (
    args: DayPilot.TimeRangeSelectedArgs
  ) => {
    // Abrir modal de criação
    // TODO: Integrar com modal
    console.log("Novo agendamento:", {
      profissional_id: args.resource,
      data_hora_inicio: args.start.toString(),
      data_hora_fim: args.end.toString(),
    });
  };

  const handleEventClick = (args: DayPilot.EventClickArgs) => {
    // Abrir modal de detalhes/edição
    console.log("Agendamento clicado:", args.e.data);
  };

  const handleEventMove = async (args: DayPilot.EventMoveArgs) => {
    // Reagendar (mover horário ou profissional)
    try {
      await updateMutation.mutateAsync({
        id: args.e.id(),
        data: {
          profissional_id: args.newResource,
          data_hora_inicio: args.newStart.toString(),
          data_hora_fim: args.newEnd.toString(),
        },
      });
    } catch (error) {
      args.preventDefault();
      console.error("Erro ao mover agendamento:", error);
    }
  };

  const handleEventResize = async (args: DayPilot.EventResizeArgs) => {
    // Alterar duração
    try {
      await updateMutation.mutateAsync({
        id: args.e.id(),
        data: {
          data_hora_inicio: args.newStart.toString(),
          data_hora_fim: args.newEnd.toString(),
        },
      });
    } catch (error) {
      args.preventDefault();
      console.error("Erro ao redimensionar agendamento:", error);
    }
  };

  return (
    <div style={{ width: "100%", height: "700px" }}>
      <DayPilotScheduler
        {...schedulerRef}
        viewType="Week"
        scale="Day"
        timeHeaders={[
          { groupBy: "Day", format: "dddd, d \\d\\e MMMM" },
          { groupBy: "Hour" },
        ]}
        locale="pt-br"
        cellHeight={60}
        eventHeight={50}
        startDate={DayPilot.Date.today()}
        days={7}
        businessBeginsHour={8}
        businessEndsHour={20}
        showNonBusiness={false}
        resources={resources}
        events={events}
        eventMoveHandling="Update"
        eventResizeHandling="Update"
        timeRangeSelectedHandling="Enabled"
        eventClickHandling="Enabled"
        onTimeRangeSelected={handleTimeRangeSelected}
        onEventClick={handleEventClick}
        onEventMove={handleEventMove}
        onEventResize={handleEventResize}
        rowHeaderColumns={[
          {
            property: "name",
            text: "Profissional",
            width: 200,
          },
        ]}
        bubble={new DayPilot.Bubble()}
        contextMenu={
          new DayPilot.Menu({
            items: [
              {
                text: "Confirmar",
                onClick: (args) => console.log("Confirmar", args),
              },
              {
                text: "Cancelar",
                onClick: (args) => console.log("Cancelar", args),
              },
              { text: "-" },
              {
                text: "Editar",
                onClick: (args) => console.log("Editar", args),
              },
              {
                text: "Excluir",
                onClick: (args) => console.log("Excluir", args),
              },
            ],
          })
        }
      />
    </div>
  );
}
```

**Checklist:**

- [ ] Criar componente Scheduler
- [ ] Integrar com TanStack Query
- [ ] Implementar event handlers (click, move, resize)
- [ ] Aplicar tema dinamicamente
- [ ] Implementar context menu
- [ ] Testar responsividade
- [ ] Adicionar loading states

---

## 🎨 Frontend - UI

### T-FE-AGENDA-003: Modal de Criação

**Arquivo:** `frontend/app/components/agendamentos/CreateAgendamentoModal.tsx`

```typescript
"use client";

import { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Modal } from "@/app/components/ui/Modal";
import { InputField } from "@/app/components/ui/InputField";
import { SelectField } from "@/app/components/ui/SelectField";
import { DateTimePicker } from "@/app/components/ui/DateTimePicker";
import { useCreateAgendamento } from "@/app/lib/hooks/useAgendamentos";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/app/lib/api/client";

const schema = z.object({
  cliente_id: z.string().uuid("Cliente inválido"),
  profissional_id: z.string().uuid("Profissional inválido"),
  servico_id: z.string().uuid("Serviço inválido"),
  data_hora_inicio: z.date(),
  observacoes: z.string().max(500).optional(),
});

type FormData = z.infer<typeof schema>;

interface CreateAgendamentoModalProps {
  open: boolean;
  onClose: () => void;
  tenantId: string;
  initialProfissionalId?: string;
  initialDateTime?: Date;
}

export function CreateAgendamentoModal({
  open,
  onClose,
  tenantId,
  initialProfissionalId,
  initialDateTime,
}: CreateAgendamentoModalProps) {
  const { createAgendamento, isPending } = useCreateAgendamento();

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      profissional_id: initialProfissionalId,
      data_hora_inicio: initialDateTime || new Date(),
    },
  });

  // Fetch clientes
  const { data: clientes = [] } = useQuery({
    queryKey: ["clientes", tenantId, { ativo: true }],
    queryFn: () => api.clientes.list(tenantId, { ativo: true }),
  });

  // Fetch profissionais
  const { data: profissionais = [] } = useQuery({
    queryKey: ["profissionais", tenantId, { ativo: true }],
    queryFn: () => api.profissionais.list(tenantId, { ativo: true }),
  });

  // Fetch serviços
  const servicoIdSelecionado = watch("servico_id");
  const { data: servicos = [] } = useQuery({
    queryKey: ["servicos", tenantId, { ativo: true }],
    queryFn: () => api.servicos.list(tenantId, { ativo: true }),
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createAgendamento.mutateAsync({
        ...data,
        tenant_id: tenantId,
      });
      reset();
      onClose();
    } catch (error) {
      console.error("Erro ao criar agendamento:", error);
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Novo Agendamento"
      subtitle="Agende um horário para o cliente"
      onConfirm={handleSubmit(onSubmit)}
      confirmText="Agendar"
      isLoading={isPending}
      maxWidth="md"
    >
      <form
        id="create-agendamento-form"
        onSubmit={handleSubmit(onSubmit)}
        className="flex flex-col gap-4"
      >
        <SelectField
          name="cliente_id"
          control={control}
          label="Cliente"
          placeholder="Selecione um cliente"
          options={clientes.map((c) => ({
            label: c.nome,
            value: c.id,
          }))}
          error={errors.cliente_id?.message}
          required
        />

        <SelectField
          name="profissional_id"
          control={control}
          label="Profissional"
          placeholder="Selecione um profissional"
          options={profissionais.map((p) => ({
            label: p.nome,
            value: p.id,
          }))}
          error={errors.profissional_id?.message}
          required
        />

        <SelectField
          name="servico_id"
          control={control}
          label="Serviço"
          placeholder="Selecione um serviço"
          options={servicos.map((s) => ({
            label: `${s.nome} - R$ ${s.preco} (${s.duracao} min)`,
            value: s.id,
          }))}
          error={errors.servico_id?.message}
          required
        />

        <Controller
          name="data_hora_inicio"
          control={control}
          render={({ field }) => (
            <DateTimePicker
              label="Data e Hora"
              value={field.value}
              onChange={field.onChange}
              error={errors.data_hora_inicio?.message}
              required
            />
          )}
        />

        <InputField
          name="observacoes"
          control={control}
          label="Observações"
          placeholder="Informações adicionais (opcional)"
          multiline
          rows={3}
          error={errors.observacoes?.message}
        />
      </form>
    </Modal>
  );
}
```

**Checklist:**

- [ ] Criar modal de criação
- [ ] Implementar validação com Zod
- [ ] Integrar com react-hook-form
- [ ] Fetch dinâmico de clientes/profissionais/serviços
- [ ] Implementar DateTimePicker
- [ ] Tratamento de erros
- [ ] Loading states

---

### T-FE-AGENDA-004: Página Principal

**Arquivo:** `frontend/app/(private)/agendamentos/page.tsx`

```typescript
"use client";

import { useState } from "react";
import { Container, Box, Typography, Button } from "@mui/material";
import { Plus } from "lucide-react";
import { Scheduler } from "@/app/components/agendamentos/Scheduler";
import { CreateAgendamentoModal } from "@/app/components/agendamentos/CreateAgendamentoModal";
import { useTenant } from "@/app/lib/hooks/useTenant";

export default function AgendamentosPage() {
  const { tenantId } = useTenant();
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [dataInicio, setDataInicio] = useState(new Date());
  const [dataFim, setDataFim] = useState(
    new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
  );

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h1">Agendamentos</Typography>

        <Button
          variant="contained"
          startIcon={<Plus size={20} />}
          onClick={() => setCreateModalOpen(true)}
        >
          Novo Agendamento
        </Button>
      </Box>

      <Scheduler
        tenantId={tenantId}
        dataInicio={dataInicio}
        dataFim={dataFim}
      />

      <CreateAgendamentoModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        tenantId={tenantId}
      />
    </Container>
  );
}
```

**Checklist:**

- [ ] Criar página principal
- [ ] Integrar Scheduler
- [ ] Adicionar botão de novo agendamento
- [ ] Implementar filtros de período
- [ ] Implementar filtros por profissional/status
- [ ] Adicionar visualizações alternativas (lista, semana, mês)

---

## 📱 Notificações

### T-BE-AGENDA-007: Integração Asaas (SMS/WhatsApp)

**Arquivo:** `backend/internal/infrastructure/external/asaas/notification_service.go`

```go
package asaas

type NotificationService struct {
    client *AsaasClient
}

// Enviar confirmação de agendamento
func (s *NotificationService) SendAgendamentoConfirmacao(
    ctx context.Context,
    agendamento *entity.Agendamento,
    cliente *entity.Cliente,
) error {
    // TODO: Implementar envio via Asaas API
    // POST /notifications
    return nil
}

// Enviar lembrete de agendamento
func (s *NotificationService) SendAgendamentoLembrete(
    ctx context.Context,
    agendamento *entity.Agendamento,
    cliente *entity.Cliente,
) error {
    // TODO: Implementar envio via Asaas API
    // Template: "Lembrete: Você tem agendamento amanhã às XX:XX"
    return nil
}
```

**Checklist:**

- [ ] Pesquisar API de notificações Asaas
- [ ] Implementar envio de confirmação
- [ ] Implementar envio de lembrete (D-1)
- [ ] Criar templates de mensagens
- [ ] Implementar retry e error handling
- [ ] Testar com números reais

---

### T-BE-AGENDA-008: Cron de Lembretes

**Arquivo:** `backend/internal/infrastructure/scheduler/jobs/agendamento_reminder_job.go`

```go
package jobs

type AgendamentoReminderJob struct {
    agendamentoRepo repository.AgendamentoRepository
    clienteRepo     repository.ClienteRepository
    notificationSvc *asaas.NotificationService
    logger          *zap.Logger
}

func (j *AgendamentoReminderJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    j.logger.Info("Starting agendamento reminder job")

    // Buscar agendamentos para amanhã que ainda não foram notificados
    amanha := time.Now().AddDate(0, 0, 1)
    amanhaStart := time.Date(amanha.Year(), amanha.Month(), amanha.Day(), 0, 0, 0, 0, time.Local)
    amanhaEnd := amanhaStart.Add(24 * time.Hour)

    tenants, _ := j.tenantRepo.FindActive(ctx)

    for _, tenant := range tenants {
        agendamentos, _ := j.agendamentoRepo.FindByTenantAndPeriod(
            ctx, tenant.ID, amanhaStart, amanhaEnd,
        )

        for _, agendamento := range agendamentos {
            // Apenas agendados ou confirmados
            if agendamento.Status != entity.StatusAgendado &&
                agendamento.Status != entity.StatusConfirmado {
                continue
            }

            cliente, _ := j.clienteRepo.FindByID(ctx, tenant.ID, agendamento.ClienteID)
            if cliente == nil || cliente.Telefone == "" {
                continue
            }

            // Enviar lembrete
            if err := j.notificationSvc.SendAgendamentoLembrete(ctx, agendamento, cliente); err != nil {
                j.logger.Error("Failed to send reminder",
                    zap.String("agendamento_id", agendamento.ID),
                    zap.Error(err))
            } else {
                j.logger.Info("Reminder sent",
                    zap.String("agendamento_id", agendamento.ID),
                    zap.String("cliente", cliente.Nome))
            }
        }
    }

    j.logger.Info("Agendamento reminder job completed")
}
```

**Schedule:** `0 18 * * *` (18h todo dia)

**Checklist:**

- [ ] Implementar cron job
- [ ] Buscar agendamentos D+1
- [ ] Enviar lembretes
- [ ] Registrar envios para não duplicar
- [ ] Testar em staging

---

## 🧪 Testes

### T-BE-AGENDA-009: Testes Unitários

**Checklist:**

- [ ] Testar entidade Agendamento (validações, transições)
- [ ] Testar CreateAgendamentoUseCase (>85% coverage)
- [ ] Testar VerificarConflito repository
- [ ] Testar validação de horário de trabalho
- [ ] Testar cancelamento com motivo

---

### T-BE-AGENDA-010: Testes de Integração

**Checklist:**

- [ ] Testar criação de agendamento via API
- [ ] Testar detecção de conflito
- [ ] Testar reagendamento (update)
- [ ] Testar transições de status
- [ ] Testar listagem com filtros

---

### T-FE-AGENDA-005: Testes E2E

**Checklist:**

- [ ] Testar criação de agendamento via scheduler
- [ ] Testar drag & drop de agendamento
- [ ] Testar resize de agendamento
- [ ] Testar cancelamento
- [ ] Testar filtros de visualização

---

## 📊 Critérios de Aceitação

### Backend

- ✅ CRUD completo de agendamentos
- ✅ Detecção de conflitos de horário
- ✅ Validação de horário de trabalho
- ✅ Transições de status implementadas
- ✅ >85% test coverage
- ✅ API documentada no Swagger

### Frontend

- ✅ Scheduler DayPilot integrado
- ✅ Criação, edição, cancelamento via modal
- ✅ Drag & drop funcional
- ✅ Filtros de visualização
- ✅ Tema sincronizado (light/dark)
- ✅ Responsivo (desktop/tablet)

### Notificações

- ✅ Confirmação enviada ao criar agendamento
- ✅ Lembrete enviado D-1 às 18h
- ✅ Templates de mensagem customizáveis

---

## 🚀 Ordem de Implementação

1. **Backend Core** (5-7 dias)

   - Entidades, repositories, use cases
   - Migrations
   - Handlers HTTP

2. **Integração DayPilot** (3-4 dias)

   - Adapter
   - Componente Scheduler
   - Event handlers

3. **Frontend UI** (4-5 dias)

   - Modais de criação/edição
   - Página principal
   - Filtros e visualizações

4. **Notificações** (2-3 dias)

   - Integração Asaas
   - Cron de lembretes
   - Templates

5. **Testes & Refinamento** (3-4 dias)
   - Testes unitários
   - Testes de integração
   - Testes E2E
   - Bug fixes

**Total Estimado:** 17-23 dias úteis

---

## 📝 Notas Finais

- **Prioridade:** Alta (módulo core do negócio)
- **Dependências Críticas:** Módulo Cadastro completo
- **Riscos:** Complexidade do DayPilot, detecção de conflitos
- **Mitigação:** Prototipar DayPilot primeiro, usar índice parcial para conflitos

---

**Status:** 📋 Planejamento Completo
**Próxima Fase:** Implementação Backend Core
