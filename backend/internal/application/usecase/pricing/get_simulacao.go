package pricing

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetSimulacaoUseCase struct {
	repo   port.PrecificacaoSimulacaoRepository
	logger *zap.Logger
}

func NewGetSimulacaoUseCase(repo port.PrecificacaoSimulacaoRepository, logger *zap.Logger) *GetSimulacaoUseCase {
	return &GetSimulacaoUseCase{repo: repo, logger: logger}
}

func (uc *GetSimulacaoUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.PrecificacaoSimulacao, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	sim, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar simulação: %w", err)
	}

	return sim, nil
}
