package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// PrecificacaoConfigRepository define operações para Configuração de Precificação
type PrecificacaoConfigRepository interface {
	// Create cria uma nova configuração
	Create(ctx context.Context, config *entity.PrecificacaoConfig) error

	// FindByTenantID busca a configuração do tenant (deve ser única por tenant)
	FindByTenantID(ctx context.Context, tenantID string) (*entity.PrecificacaoConfig, error)

	// Update atualiza a configuração existente
	Update(ctx context.Context, config *entity.PrecificacaoConfig) error

	// Delete remove a configuração
	Delete(ctx context.Context, tenantID string) error
}

// PrecificacaoSimulacaoRepository define operações para Simulações de Precificação
type PrecificacaoSimulacaoRepository interface {
	// Create cria uma nova simulação
	Create(ctx context.Context, simulacao *entity.PrecificacaoSimulacao) error

	// FindByID busca uma simulação por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.PrecificacaoSimulacao, error)

	// Update atualiza uma simulação existente
	Update(ctx context.Context, simulacao *entity.PrecificacaoSimulacao) error

	// Delete remove uma simulação
	Delete(ctx context.Context, tenantID, id string) error

	// ListByItem lista simulações de um item específico
	ListByItem(ctx context.Context, tenantID, itemID, tipoItem string, filters SimulacaoListFilters) ([]*entity.PrecificacaoSimulacao, error)

	// List lista todas as simulações com filtros
	List(ctx context.Context, tenantID string, filters SimulacaoListFilters) ([]*entity.PrecificacaoSimulacao, error)

	// GetLatestByItem busca a simulação mais recente de um item
	GetLatestByItem(ctx context.Context, tenantID, itemID, tipoItem string) (*entity.PrecificacaoSimulacao, error)
}

// SimulacaoListFilters filtros para listagem de simulações
type SimulacaoListFilters struct {
	TipoItem *string // SERVICO ou PRODUTO
	Page     int
	PageSize int
	OrderBy  string
}
