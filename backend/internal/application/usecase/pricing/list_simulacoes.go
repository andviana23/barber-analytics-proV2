package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type ListSimulacoesInput struct {
	TenantID string
	ItemID   *string
	TipoItem *string
}

type ListSimulacoesUseCase struct {
	repo   port.PrecificacaoSimulacaoRepository
	logger *zap.Logger
}

func NewListSimulacoesUseCase(repo port.PrecificacaoSimulacaoRepository, logger *zap.Logger) *ListSimulacoesUseCase {
	return &ListSimulacoesUseCase{repo: repo, logger: logger}
}

func (uc *ListSimulacoesUseCase) Execute(ctx context.Context, input ListSimulacoesInput) ([]*entity.PrecificacaoSimulacao, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	var sims []*entity.PrecificacaoSimulacao
	var err error

	filters := port.SimulacaoListFilters{}
	if input.TipoItem != nil {
		filters.TipoItem = input.TipoItem
	}

	if input.ItemID != nil {
		tipoItem := ""
		if input.TipoItem != nil {
			tipoItem = *input.TipoItem
		}
		sims, err = uc.repo.ListByItem(ctx, input.TenantID, *input.ItemID, tipoItem, filters)
	} else {
		sims, err = uc.repo.List(ctx, input.TenantID, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao listar simulações: %w", err)
	}

	return sims, nil
}
