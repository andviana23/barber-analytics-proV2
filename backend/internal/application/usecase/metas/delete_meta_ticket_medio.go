package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type DeleteMetaTicketMedioUseCase struct {
	repo   port.MetaTicketMedioRepository
	logger *zap.Logger
}

func NewDeleteMetaTicketMedioUseCase(repo port.MetaTicketMedioRepository, logger *zap.Logger) *DeleteMetaTicketMedioUseCase {
	return &DeleteMetaTicketMedioUseCase{repo: repo, logger: logger}
}

func (uc *DeleteMetaTicketMedioUseCase) Execute(ctx context.Context, tenantID, id string) error {
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

	uc.logger.Info("Meta ticket médio deletada", zap.String("tenant_id", tenantID), zap.String("id", id))

	return nil
}
