// Package port contém as interfaces (ports) que definem contratos de comunicação
// entre a camada de domínio e a camada de infraestrutura.
// Segue o padrão Ports and Adapters (Hexagonal Architecture).
package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// DREMensalRepository define operações de persistência para DRE Mensal
type DREMensalRepository interface {
	// Create cria um novo DRE mensal
	Create(ctx context.Context, dre *entity.DREMensal) error

	// FindByID busca um DRE por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.DREMensal, error)

	// FindByMesAno busca DRE de um mês específico
	FindByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.DREMensal, error)

	// Update atualiza um DRE existente
	Update(ctx context.Context, dre *entity.DREMensal) error

	// Delete remove um DRE
	Delete(ctx context.Context, tenantID, id string) error

	// List lista DREs com filtros
	List(ctx context.Context, tenantID string, filters DREListFilters) ([]*entity.DREMensal, error)

	// ListByPeriod lista DREs em um período
	ListByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) ([]*entity.DREMensal, error)
}

// DREListFilters filtros para listagem de DREs
type DREListFilters struct {
	Page     int
	PageSize int
	OrderBy  string
}

// FluxoCaixaDiarioRepository define operações para Fluxo de Caixa Diário
type FluxoCaixaDiarioRepository interface {
	// Create cria um novo fluxo de caixa diário
	Create(ctx context.Context, fluxo *entity.FluxoCaixaDiario) error

	// FindByID busca um fluxo por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.FluxoCaixaDiario, error)

	// FindByData busca fluxo de uma data específica
	FindByData(ctx context.Context, tenantID string, data time.Time) (*entity.FluxoCaixaDiario, error)

	// Update atualiza um fluxo existente
	Update(ctx context.Context, fluxo *entity.FluxoCaixaDiario) error

	// Delete remove um fluxo
	Delete(ctx context.Context, tenantID, id string) error

	// ListByDateRange lista fluxos em um período
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.FluxoCaixaDiario, error)
}

// CompensacaoBancariaRepository define operações para Compensações Bancárias
type CompensacaoBancariaRepository interface {
	// Create cria uma nova compensação
	Create(ctx context.Context, comp *entity.CompensacaoBancaria) error

	// FindByID busca uma compensação por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.CompensacaoBancaria, error)

	// FindByReceitaID busca compensação de uma receita
	FindByReceitaID(ctx context.Context, tenantID, receitaID string) (*entity.CompensacaoBancaria, error)

	// Update atualiza uma compensação
	Update(ctx context.Context, comp *entity.CompensacaoBancaria) error

	// Delete remove uma compensação
	Delete(ctx context.Context, tenantID, id string) error

	// List lista compensações com filtros
	List(ctx context.Context, tenantID string, filters CompensacaoListFilters) ([]*entity.CompensacaoBancaria, error)

	// ListByStatus lista compensações por status
	ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusCompensacao) ([]*entity.CompensacaoBancaria, error)

	// ListPendentesCompensacao lista compensações que podem ser marcadas como compensadas (data <= hoje)
	ListPendentesCompensacao(ctx context.Context, tenantID string) ([]*entity.CompensacaoBancaria, error)

	// ListByDateRange lista compensações em um período
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.CompensacaoBancaria, error)
}

// CompensacaoListFilters filtros para listagem de compensações
type CompensacaoListFilters struct {
	Status   *valueobject.StatusCompensacao
	Page     int
	PageSize int
	OrderBy  string
}
