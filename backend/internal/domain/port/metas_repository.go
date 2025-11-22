package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// MetaMensalRepository define operações para Metas Mensais
type MetaMensalRepository interface {
	// Create cria uma nova meta mensal
	Create(ctx context.Context, meta *entity.MetaMensal) error

	// FindByID busca uma meta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.MetaMensal, error)

	// FindByMesAno busca meta de um mês específico
	FindByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.MetaMensal, error)

	// Update atualiza uma meta existente
	Update(ctx context.Context, meta *entity.MetaMensal) error

	// Delete remove uma meta
	Delete(ctx context.Context, tenantID, id string) error

	// ListAtivas lista metas ativas
	ListAtivas(ctx context.Context, tenantID string) ([]*entity.MetaMensal, error)

	// ListByPeriod lista metas em um período
	ListByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) ([]*entity.MetaMensal, error)
}

// MetaBarbeiroRepository define operações para Metas de Barbeiro
type MetaBarbeiroRepository interface {
	// Create cria uma nova meta de barbeiro
	Create(ctx context.Context, meta *entity.MetaBarbeiro) error

	// FindByID busca uma meta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.MetaBarbeiro, error)

	// FindByBarbeiroMesAno busca meta de um barbeiro em um mês
	FindByBarbeiroMesAno(ctx context.Context, tenantID, barbeiroID string, mesAno valueobject.MesAno) (*entity.MetaBarbeiro, error)

	// Update atualiza uma meta existente
	Update(ctx context.Context, meta *entity.MetaBarbeiro) error

	// Delete remove uma meta
	Delete(ctx context.Context, tenantID, id string) error

	// ListByBarbeiro lista todas as metas de um barbeiro
	ListByBarbeiro(ctx context.Context, tenantID, barbeiroID string) ([]*entity.MetaBarbeiro, error)

	// ListByMesAno lista todas as metas de barbeiros de um mês
	ListByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) ([]*entity.MetaBarbeiro, error)
}

// MetaTicketMedioRepository define operações para Metas de Ticket Médio
type MetaTicketMedioRepository interface {
	// Create cria uma nova meta de ticket médio
	Create(ctx context.Context, meta *entity.MetaTicketMedio) error

	// FindByID busca uma meta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.MetaTicketMedio, error)

	// FindGeralByMesAno busca meta geral de um mês
	FindGeralByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.MetaTicketMedio, error)

	// FindBarbeiroByMesAno busca meta de um barbeiro em um mês
	FindBarbeiroByMesAno(ctx context.Context, tenantID, barbeiroID string, mesAno valueobject.MesAno) (*entity.MetaTicketMedio, error)

	// Update atualiza uma meta existente
	Update(ctx context.Context, meta *entity.MetaTicketMedio) error

	// Delete remove uma meta
	Delete(ctx context.Context, tenantID, id string) error

	// ListByMesAno lista todas as metas de um mês (geral + barbeiros)
	ListByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) ([]*entity.MetaTicketMedio, error)

	// ListByBarbeiro lista metas de um barbeiro
	ListByBarbeiro(ctx context.Context, tenantID, barbeiroID string) ([]*entity.MetaTicketMedio, error)
}
