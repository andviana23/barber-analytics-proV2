package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// ContaPagarRepository define operações para Contas a Pagar
type ContaPagarRepository interface {
	// Create cria uma nova conta a pagar
	Create(ctx context.Context, conta *entity.ContaPagar) error

	// FindByID busca uma conta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.ContaPagar, error)

	// Update atualiza uma conta existente
	Update(ctx context.Context, conta *entity.ContaPagar) error

	// Delete remove uma conta
	Delete(ctx context.Context, tenantID, id string) error

	// List lista contas com filtros
	List(ctx context.Context, tenantID string, filters ContaPagarListFilters) ([]*entity.ContaPagar, error)

	// ListByStatus lista contas por status
	ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaPagar, error)

	// ListVencendoEm lista contas que vencem em até N dias
	ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaPagar, error)

	// ListAtrasadas lista contas atrasadas
	ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaPagar, error)

	// ListByDateRange lista contas em um período (data vencimento)
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaPagar, error)

	// SumByPeriod soma valores de contas em um período
	SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error)

	// SumByCategoria soma valores por categoria
	SumByCategoria(ctx context.Context, tenantID, categoriaID string, inicio, fim time.Time) (valueobject.Money, error)
}

// ContaPagarListFilters filtros para listagem de contas a pagar
type ContaPagarListFilters struct {
	Status      *valueobject.StatusConta
	Tipo        *valueobject.TipoCusto
	CategoriaID *string
	Fornecedor  *string
	Recorrente  *bool
	Page        int
	PageSize    int
	OrderBy     string
}

// ContaReceberRepository define operações para Contas a Receber
type ContaReceberRepository interface {
	// Create cria uma nova conta a receber
	Create(ctx context.Context, conta *entity.ContaReceber) error

	// FindByID busca uma conta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.ContaReceber, error)

	// Update atualiza uma conta existente
	Update(ctx context.Context, conta *entity.ContaReceber) error

	// Delete remove uma conta
	Delete(ctx context.Context, tenantID, id string) error

	// List lista contas com filtros
	List(ctx context.Context, tenantID string, filters ContaReceberListFilters) ([]*entity.ContaReceber, error)

	// ListByStatus lista contas por status
	ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaReceber, error)

	// ListByAssinatura lista contas de uma assinatura
	ListByAssinatura(ctx context.Context, tenantID, assinaturaID string) ([]*entity.ContaReceber, error)

	// ListVencendoEm lista contas que vencem em até N dias
	ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaReceber, error)

	// ListAtrasadas lista contas atrasadas
	ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaReceber, error)

	// ListByDateRange lista contas em um período (data vencimento)
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaReceber, error)

	// SumByPeriod soma valores de contas em um período
	SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error)

	// SumByOrigem soma valores por origem
	SumByOrigem(ctx context.Context, tenantID, origem string, inicio, fim time.Time) (valueobject.Money, error)
}

// ContaReceberListFilters filtros para listagem de contas a receber
type ContaReceberListFilters struct {
	Status       *valueobject.StatusConta
	Origem       *string
	AssinaturaID *string
	Page         int
	PageSize     int
	OrderBy      string
}
