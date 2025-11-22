package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type UpdatePrecificacaoConfigInput struct {
	TenantID string
	Config   *entity.PrecificacaoConfig
}

type UpdatePrecificacaoConfigUseCase struct {
	repo   port.PrecificacaoConfigRepository
	logger *zap.Logger
}

func NewUpdatePrecificacaoConfigUseCase(repo port.PrecificacaoConfigRepository, logger *zap.Logger) *UpdatePrecificacaoConfigUseCase {
	return &UpdatePrecificacaoConfigUseCase{repo: repo, logger: logger}
}

func (uc *UpdatePrecificacaoConfigUseCase) Execute(ctx context.Context, input UpdatePrecificacaoConfigInput) (*entity.PrecificacaoConfig, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Verificar se existe
	_, err := uc.repo.FindByTenantID(ctx, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("configuração não encontrada: %w", err)
	}

	// Atualizar
	if err := uc.repo.Update(ctx, input.Config); err != nil {
		return nil, fmt.Errorf("erro ao atualizar configuração: %w", err)
	}

	uc.logger.Info("Configuração de precificação atualizada", zap.String("tenant_id", input.TenantID))

	return input.Config, nil
}
