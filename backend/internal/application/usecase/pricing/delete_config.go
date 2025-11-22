package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type DeletePrecificacaoConfigUseCase struct {
	repo   port.PrecificacaoConfigRepository
	logger *zap.Logger
}

func NewDeletePrecificacaoConfigUseCase(repo port.PrecificacaoConfigRepository, logger *zap.Logger) *DeletePrecificacaoConfigUseCase {
	return &DeletePrecificacaoConfigUseCase{repo: repo, logger: logger}
}

func (uc *DeletePrecificacaoConfigUseCase) Execute(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}

	if _, err := uc.repo.FindByTenantID(ctx, tenantID); err != nil {
		return fmt.Errorf("configuração não encontrada: %w", err)
	}

	if err := uc.repo.Delete(ctx, tenantID); err != nil {
		return fmt.Errorf("erro ao deletar configuração: %w", err)
	}

	uc.logger.Info("Configuração de precificação deletada", zap.String("tenant_id", tenantID))

	return nil
}
