package metas

import (
"context"
"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

// DeleteMetaMensalUseCase implementa a deleção de meta mensal
type DeleteMetaMensalUseCase struct {
	repo   port.MetaMensalRepository
	logger *zap.Logger
}

// NewDeleteMetaMensalUseCase cria nova instância do use case
func NewDeleteMetaMensalUseCase(
repo port.MetaMensalRepository,
logger *zap.Logger,
) *DeleteMetaMensalUseCase {
	return &DeleteMetaMensalUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute remove uma meta mensal
func (uc *DeleteMetaMensalUseCase) Execute(ctx context.Context, tenantID, id string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if id == "" {
		return domain.ErrInvalidID
	}

	// Verificar se meta existe
	_, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("meta não encontrada: %w", err)
	}

	// Deletar
	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("erro ao deletar meta: %w", err)
	}

	uc.logger.Info("Meta mensal deletada",
zap.String("tenant_id", tenantID),
zap.String("id", id),
)

	return nil
}
