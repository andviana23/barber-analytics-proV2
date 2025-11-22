package user

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type GetUserPreferencesUseCase struct {
	repo   port.UserPreferencesRepository
	logger *zap.Logger
}

func NewGetUserPreferencesUseCase(repo port.UserPreferencesRepository, logger *zap.Logger) *GetUserPreferencesUseCase {
	return &GetUserPreferencesUseCase{repo: repo, logger: logger}
}

func (uc *GetUserPreferencesUseCase) Execute(ctx context.Context, tenantID, userID string) (*entity.UserPreferences, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if userID == "" {
		return nil, domain.ErrInvalidID
	}

	prefs, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar preferências do usuário: %w", err)
	}

	return prefs, nil
}
