package metas

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// ListMetasMensaisInput define os filtros para listagem
type ListMetasMensaisInput struct {
	TenantID string
	Inicio   valueobject.MesAno
	Fim      valueobject.MesAno
}

// ListMetasMensaisUseCase implementa a listagem de metas mensais
type ListMetasMensaisUseCase struct {
	repo   port.MetaMensalRepository
	logger *zap.Logger
}

// NewListMetasMensaisUseCase cria nova instância do use case
func NewListMetasMensaisUseCase(
	repo port.MetaMensalRepository,
	logger *zap.Logger,
) *ListMetasMensaisUseCase {
	return &ListMetasMensaisUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista metas mensais com filtros
func (uc *ListMetasMensaisUseCase) Execute(ctx context.Context, input ListMetasMensaisInput) ([]*entity.MetaMensal, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	metas, err := uc.repo.ListByPeriod(ctx, input.TenantID, input.Inicio, input.Fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas mensais: %w", err)
	}

	return metas, nil
}
