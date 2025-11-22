package pricing

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type DeleteSimulacaoUseCase struct {
	repo   port.PrecificacaoSimulacaoRepository
	logger *zap.Logger
}

func NewDeleteSimulacaoUseCase(repo port.PrecificacaoSimulacaoRepository, logger *zap.Logger) *DeleteSimulacaoUseCase {
	return &DeleteSimulacaoUseCase{repo: repo, logger: logger}
}

func (uc *DeleteSimulacaoUseCase) Execute(ctx context.Context, tenantID, id string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if id == "" {
		return domain.ErrInvalidID
	}

	if _, err := uc.repo.FindByID(ctx, tenantID, id); err != nil {
		return fmt.Errorf("simulação não encontrada: %w", err)
	}

	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("erro ao deletar simulação: %w", err)
	}

	uc.logger.Info("Simulação de precificação deletada", zap.String("tenant_id", tenantID), zap.String("id", id))

	return nil
}
