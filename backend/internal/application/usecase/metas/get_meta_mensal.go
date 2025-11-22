package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

// GetMetaMensalUseCase implementa a busca de meta mensal por ID
type GetMetaMensalUseCase struct {
	repo   port.MetaMensalRepository
	logger *zap.Logger
}

// NewGetMetaMensalUseCase cria nova instância do use case
func NewGetMetaMensalUseCase(
repo port.MetaMensalRepository,
logger *zap.Logger,
) *GetMetaMensalUseCase {
	return &GetMetaMensalUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca uma meta mensal por ID
func (uc *GetMetaMensalUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.MetaMensal, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	meta, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta mensal: %w", err)
	}

	return meta, nil
}
