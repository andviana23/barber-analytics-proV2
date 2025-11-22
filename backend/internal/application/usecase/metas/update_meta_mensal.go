package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
"go.uber.org/zap"
)

// UpdateMetaMensalInput define os dados para atualização
type UpdateMetaMensalInput struct {
	TenantID        string
	ID              string
	MetaFaturamento valueobject.Money
}

// UpdateMetaMensalUseCase implementa a atualização de meta mensal
type UpdateMetaMensalUseCase struct {
	repo   port.MetaMensalRepository
	logger *zap.Logger
}

// NewUpdateMetaMensalUseCase cria nova instância do use case
func NewUpdateMetaMensalUseCase(
repo port.MetaMensalRepository,
logger *zap.Logger,
) *UpdateMetaMensalUseCase {
	return &UpdateMetaMensalUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute atualiza uma meta mensal existente
func (uc *UpdateMetaMensalUseCase) Execute(ctx context.Context, input UpdateMetaMensalInput) (*entity.MetaMensal, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	// Buscar meta existente
	meta, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta: %w", err)
	}

	// Atualizar valor
	if err := meta.AtualizarMeta(input.MetaFaturamento); err != nil {
		return nil, fmt.Errorf("erro ao atualizar meta: %w", err)
	}

	// Salvar
	if err := uc.repo.Update(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta: %w", err)
	}

	uc.logger.Info("Meta mensal atualizada",
zap.String("tenant_id", input.TenantID),
zap.String("id", input.ID),
zap.String("nova_meta", input.MetaFaturamento.String()),
	)

	return meta, nil
}
