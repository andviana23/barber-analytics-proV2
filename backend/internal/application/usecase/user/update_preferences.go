package user

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type UpdateUserPreferencesInput struct {
	TenantID    string
	UserID      string
	Preferences *entity.UserPreferences
}

type UpdateUserPreferencesUseCase struct {
	repo   port.UserPreferencesRepository
	logger *zap.Logger
}

func NewUpdateUserPreferencesUseCase(repo port.UserPreferencesRepository, logger *zap.Logger) *UpdateUserPreferencesUseCase {
	return &UpdateUserPreferencesUseCase{repo: repo, logger: logger}
}

func (uc *UpdateUserPreferencesUseCase) Execute(ctx context.Context, input UpdateUserPreferencesInput) (*entity.UserPreferences, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.UserID == "" {
		return nil, domain.ErrInvalidID
	}

	// Verificar se existe
	_, err := uc.repo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("preferências não encontradas: %w", err)
	}

	// Atualizar
	if err := uc.repo.Update(ctx, input.Preferences); err != nil {
		return nil, fmt.Errorf("erro ao atualizar preferências: %w", err)
	}

	uc.logger.Info("Preferências de usuário atualizadas",
		zap.String("tenant_id", input.TenantID),
		zap.String("user_id", input.UserID))

	return input.Preferences, nil
}
