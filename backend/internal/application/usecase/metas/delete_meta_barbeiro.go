package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type DeleteMetaBarbeiroUseCase struct {
	repo   port.MetaBarbeiroRepository
	logger *zap.Logger
}

func NewDeleteMetaBarbeiroUseCase(repo port.MetaBarbeiroRepository, logger *zap.Logger) *DeleteMetaBarbeiroUseCase {
	return &DeleteMetaBarbeiroUseCase{repo: repo, logger: logger}
}

func (uc *DeleteMetaBarbeiroUseCase) Execute(ctx context.Context, tenantID, id string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if id == "" {
		return domain.ErrInvalidID
	}

	if _, err := uc.repo.FindByID(ctx, tenantID, id); err != nil {
		return fmt.Errorf("meta não encontrada: %w", err)
	}

	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("erro ao deletar meta: %w", err)
	}

	uc.logger.Info("Meta barbeiro deletada", zap.String("tenant_id", tenantID), zap.String("id", id))

	return nil
}
