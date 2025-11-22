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

type ListMetasBarbeiroInput struct {
	TenantID   string
	BarbeiroID *string
	Inicio     valueobject.MesAno
	Fim        valueobject.MesAno
}

type ListMetasBarbeiroUseCase struct {
	repo   port.MetaBarbeiroRepository
	logger *zap.Logger
}

func NewListMetasBarbeiroUseCase(repo port.MetaBarbeiroRepository, logger *zap.Logger) *ListMetasBarbeiroUseCase {
	return &ListMetasBarbeiroUseCase{repo: repo, logger: logger}
}

func (uc *ListMetasBarbeiroUseCase) Execute(ctx context.Context, input ListMetasBarbeiroInput) ([]*entity.MetaBarbeiro, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.BarbeiroID != nil {
		metas, err := uc.repo.ListByBarbeiro(ctx, input.TenantID, *input.BarbeiroID)
		if err != nil {
			return nil, fmt.Errorf("erro ao listar metas barbeiro: %w", err)
		}
		return metas, nil
	}

	metas, err := uc.repo.ListByMesAno(ctx, input.TenantID, input.Inicio)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas barbeiro: %w", err)
	}

	return metas, nil
}
