package user

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type DeleteUserPreferencesUseCase struct {
	repo   port.UserPreferencesRepository
	logger *zap.Logger
}

func NewDeleteUserPreferencesUseCase(repo port.UserPreferencesRepository, logger *zap.Logger) *DeleteUserPreferencesUseCase {
	return &DeleteUserPreferencesUseCase{repo: repo, logger: logger}
}

func (uc *DeleteUserPreferencesUseCase) Execute(ctx context.Context, tenantID, userID string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if userID == "" {
		return domain.ErrInvalidID
	}

	// Verificar se existe
	if _, err := uc.repo.FindByUserID(ctx, userID); err != nil {
		return fmt.Errorf("preferências não encontradas: %w", err)
	}

	// Deletar
	if err := uc.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("erro ao deletar preferências: %w", err)
	}

	uc.logger.Info("Preferências de usuário deletadas",
		zap.String("tenant_id", tenantID),
		zap.String("user_id", userID))

	return nil
}
