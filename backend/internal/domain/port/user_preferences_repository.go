package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// UserPreferencesRepository define operações para Preferências do Usuário (LGPD)
type UserPreferencesRepository interface {
	// Create cria novas preferências de usuário
	Create(ctx context.Context, preferences *entity.UserPreferences) error

	// FindByUserID busca preferências de um usuário
	FindByUserID(ctx context.Context, userID string) (*entity.UserPreferences, error)

	// Update atualiza preferências existentes
	Update(ctx context.Context, preferences *entity.UserPreferences) error

	// Delete remove preferências de um usuário
	Delete(ctx context.Context, userID string) error
}
