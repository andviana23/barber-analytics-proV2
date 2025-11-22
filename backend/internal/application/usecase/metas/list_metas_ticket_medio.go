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

type ListMetasTicketMedioInput struct {
	TenantID string
	Inicio   valueobject.MesAno
	Fim      valueobject.MesAno
}

type ListMetasTicketMedioUseCase struct {
	repo   port.MetaTicketMedioRepository
	logger *zap.Logger
}

func NewListMetasTicketMedioUseCase(repo port.MetaTicketMedioRepository, logger *zap.Logger) *ListMetasTicketMedioUseCase {
	return &ListMetasTicketMedioUseCase{repo: repo, logger: logger}
}

func (uc *ListMetasTicketMedioUseCase) Execute(ctx context.Context, input ListMetasTicketMedioInput) ([]*entity.MetaTicketMedio, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	metas, err := uc.repo.ListByMesAno(ctx, input.TenantID, input.Inicio)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas ticket médio: %w", err)
	}

	return metas, nil
}
