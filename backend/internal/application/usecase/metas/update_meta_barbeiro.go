package metas

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

type UpdateMetaBarbeiroInput struct {
	TenantID        string
	ID              string
	MetaFaturamento valueobject.Money
	MetaServicos    int
}

type UpdateMetaBarbeiroUseCase struct {
	repo   port.MetaBarbeiroRepository
	logger *zap.Logger
}

func NewUpdateMetaBarbeiroUseCase(repo port.MetaBarbeiroRepository, logger *zap.Logger) *UpdateMetaBarbeiroUseCase {
	return &UpdateMetaBarbeiroUseCase{repo: repo, logger: logger}
}

func (uc *UpdateMetaBarbeiroUseCase) Execute(ctx context.Context, input UpdateMetaBarbeiroInput) (*entity.MetaBarbeiro, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	meta, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta: %w", err)
	}

	// Atualizar valores
	meta.MetaServicosGerais = input.MetaFaturamento
	meta.MetaServicosExtras = valueobject.Zero()
	meta.MetaProdutos = valueobject.Zero()
	meta.AtualizadoEm = time.Now()

	if err := uc.repo.Update(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta: %w", err)
	}

	uc.logger.Info("Meta barbeiro atualizada", zap.String("tenant_id", input.TenantID), zap.String("id", input.ID))

	return meta, nil
}
