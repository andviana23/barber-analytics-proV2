package financial

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

// GenerateDREInput define os dados de entrada para gerar DRE
type GenerateDREInput struct {
	TenantID string
	MesAno   valueobject.MesAno
}

// GenerateDREUseCase implementa a geração de DRE mensal
// Este use case é executado por cron job mensalmente
type GenerateDREUseCase struct {
	dreRepo           port.DREMensalRepository
	contasPagarRepo   port.ContaPagarRepository
	contasReceberRepo port.ContaReceberRepository
	logger            *zap.Logger
}

// NewGenerateDREUseCase cria nova instância do use case
func NewGenerateDREUseCase(
	dreRepo port.DREMensalRepository,
	contasPagarRepo port.ContaPagarRepository,
	contasReceberRepo port.ContaReceberRepository,
	logger *zap.Logger,
) *GenerateDREUseCase {
	return &GenerateDREUseCase{
		dreRepo:           dreRepo,
		contasPagarRepo:   contasPagarRepo,
		contasReceberRepo: contasReceberRepo,
		logger:            logger,
	}
}

// Execute gera ou atualiza o DRE de um mês
func (uc *GenerateDREUseCase) Execute(ctx context.Context, input GenerateDREInput) (*entity.DREMensal, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.MesAno.String() == "" {
		// Usar mês anterior se não informado
		input.MesAno = valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
	}

	// Buscar DRE existente ou criar novo
	dre, err := uc.dreRepo.FindByMesAno(ctx, input.TenantID, input.MesAno)
	if err != nil {
		// Criar novo DRE se não existir
		dre, err = entity.NewDREMensal(input.TenantID, input.MesAno)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar DRE: %w", err)
		}
	}

	// Calcular período do mês
	inicio := input.MesAno.PrimeiroDia()
	fim := input.MesAno.UltimoDia()

	// ===== RECEITAS =====
	// Calcular receitas por subtipo (SERVICO, PRODUTO, PLANO)
	// Nota: Precisa filtrar por subtipo - simplificado aqui
	statusPago := valueobject.StatusContaPago

	totalReceitas, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular receitas: %w", err)
	}

	// TODO: Implementar filtro por subtipo quando disponível no repositório
	// Por enquanto, atribuindo ao serviços (receita principal)
	dre.SetReceitas(totalReceitas, valueobject.Zero(), valueobject.Zero())

	// ===== CUSTOS VARIÁVEIS =====
	// TODO: Buscar comissões e consumo de insumos do período
	// Por enquanto, zerado até implementar módulos relacionados
	dre.SetCustosVariaveis(valueobject.Zero(), valueobject.Zero())

	// ===== DESPESAS =====
	// Calcular despesas fixas e variáveis
	tipoFixo := valueobject.TipoCustoFixo
	tipoVariavel := valueobject.TipoCustoVariavel

	despesasFixas, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular despesas fixas: %w", err)
	}

	// TODO: Filtrar por tipo quando implementado no repositório
	// Por enquanto, considerando tudo como despesa fixa
	_ = tipoFixo
	_ = tipoVariavel

	dre.SetDespesas(despesasFixas, valueobject.Zero())

	// Calcular resultado final
	dre.Calcular()

	// Persistir ou atualizar
	if dre.ProcessadoEm.IsZero() {
		if err := uc.dreRepo.Create(ctx, dre); err != nil {
			return nil, fmt.Errorf("erro ao salvar DRE: %w", err)
		}
	} else {
		if err := uc.dreRepo.Update(ctx, dre); err != nil {
			return nil, fmt.Errorf("erro ao atualizar DRE: %w", err)
		}
	}

	uc.logger.Info("DRE mensal gerado",
		zap.String("tenant_id", input.TenantID),
		zap.String("mes_ano", input.MesAno.String()),
		zap.String("receita_total", dre.ReceitaTotal.String()),
		zap.String("lucro_liquido", dre.LucroLiquido.String()),
	)

	return dre, nil
}

// DefaultMesAnterior retorna o período YYYY-MM do mês anterior ao atual.
func (uc *GenerateDREUseCase) DefaultMesAnterior() valueobject.MesAno {
	return valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
}
