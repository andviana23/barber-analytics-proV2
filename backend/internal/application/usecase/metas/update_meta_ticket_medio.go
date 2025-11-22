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

type UpdateMetaTicketMedioInput struct {
	TenantID        string
	ID              string
	MetaTicketMedio valueobject.Money
}

type UpdateMetaTicketMedioUseCase struct {
	repo   port.MetaTicketMedioRepository
	logger *zap.Logger
}

func NewUpdateMetaTicketMedioUseCase(repo port.MetaTicketMedioRepository, logger *zap.Logger) *UpdateMetaTicketMedioUseCase {
	return &UpdateMetaTicketMedioUseCase{repo: repo, logger: logger}
}

func (uc *UpdateMetaTicketMedioUseCase) Execute(ctx context.Context, input UpdateMetaTicketMedioInput) (*entity.MetaTicketMedio, error) {
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

	// Atualizar valor
	meta.MetaValor = input.MetaTicketMedio
	meta.AtualizadoEm = time.Now()

	if err := uc.repo.Update(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta: %w", err)
	}

	uc.logger.Info("Meta ticket médio atualizada", zap.String("tenant_id", input.TenantID), zap.String("id", input.ID))

	return meta, nil
}
