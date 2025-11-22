package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetMetaBarbeiroUseCase struct {
	repo   port.MetaBarbeiroRepository
	logger *zap.Logger
}

func NewGetMetaBarbeiroUseCase(repo port.MetaBarbeiroRepository, logger *zap.Logger) *GetMetaBarbeiroUseCase {
	return &GetMetaBarbeiroUseCase{repo: repo, logger: logger}
}

func (uc *GetMetaBarbeiroUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.MetaBarbeiro, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	meta, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta barbeiro: %w", err)
	}

	return meta, nil
}
