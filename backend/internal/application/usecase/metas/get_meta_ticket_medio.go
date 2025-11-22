package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetMetaTicketMedioUseCase struct {
	repo   port.MetaTicketMedioRepository
	logger *zap.Logger
}

func NewGetMetaTicketMedioUseCase(repo port.MetaTicketMedioRepository, logger *zap.Logger) *GetMetaTicketMedioUseCase {
	return &GetMetaTicketMedioUseCase{repo: repo, logger: logger}
}

func (uc *GetMetaTicketMedioUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.MetaTicketMedio, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	meta, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta ticket médio: %w", err)
	}

	return meta, nil
}
