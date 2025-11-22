package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type GetPrecificacaoConfigUseCase struct {
	repo   port.PrecificacaoConfigRepository
	logger *zap.Logger
}

func NewGetPrecificacaoConfigUseCase(repo port.PrecificacaoConfigRepository, logger *zap.Logger) *GetPrecificacaoConfigUseCase {
	return &GetPrecificacaoConfigUseCase{repo: repo, logger: logger}
}

func (uc *GetPrecificacaoConfigUseCase) Execute(ctx context.Context, tenantID string) (*entity.PrecificacaoConfig, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	config, err := uc.repo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configuração de precificação: %w", err)
	}

	return config, nil
}
