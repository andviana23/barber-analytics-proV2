package mapper

import (
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
)

// ToPrecificacaoConfigResponse converte entidade para DTO Response
func ToPrecificacaoConfigResponse(config *entity.PrecificacaoConfig) dto.PrecificacaoConfigResponse {
	return dto.PrecificacaoConfigResponse{
		ID:                        config.ID,
		MargemDesejada:            config.MargemDesejada.String(),
		MarkupAlvo:                config.MarkupAlvo.String(),
		ImpostoPercentual:         config.ImpostoPercentual.String(),
		ComissaoPercentualDefault: config.ComissaoPercentualDefault.String(),
		CriadoEm:                  config.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
		AtualizadoEm:              config.AtualizadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToPrecificacaoSimulacaoResponse converte entidade para DTO Response
func ToPrecificacaoSimulacaoResponse(sim *entity.PrecificacaoSimulacao) dto.PrecificacaoSimulacaoResponse {
	return dto.PrecificacaoSimulacaoResponse{
		ID:                  sim.ID,
		ItemID:              sim.ItemID,
		TipoItem:            sim.TipoItem,
		CustoMateriais:      sim.CustoMateriais.String(),
		CustoMaoDeObra:      sim.CustoMaoDeObra.String(),
		CustoTotal:          sim.CustoTotal.String(),
		MargemDesejada:      sim.MargemDesejada.String(),
		ComissaoPercentual:  sim.ComissaoPercentual.String(),
		ImpostoPercentual:   sim.ImpostoPercentual.String(),
		PrecoSugerido:       sim.PrecoSugerido.String(),
		PrecoAtual:          sim.PrecoAtual.String(),
		DiferencaPercentual: sim.DiferencaPercentual.String(),
		LucroEstimado:       sim.LucroEstimado.String(),
		MargemFinal:         sim.MargemFinal.String(),
		CriadoEm:            sim.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// FromSaveConfigPrecificacaoRequest converte DTO Request para parâmetros do use case
func FromSaveConfigPrecificacaoRequest(req dto.SaveConfigPrecificacaoRequest) (
	margemDesejada, impostoPercentual, comissaoDefault valueobject.Percentage,
	markupAlvo decimal.Decimal,
	err error,
) {
	// Parse margem desejada
	margemDecimal, err := decimal.NewFromString(req.MargemDesejada)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, fmt.Errorf("margem_desejada inválido: %w", err)
	}
	margemDesejada, err = valueobject.NewPercentage(margemDecimal)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, err
	}

	// Parse markup alvo
	markupAlvo, err = decimal.NewFromString(req.MarkupAlvo)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, fmt.Errorf("markup_alvo inválido: %w", err)
	}

	// Parse impostos
	impostosDecimal, err := decimal.NewFromString(req.ImpostoPercentual)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, fmt.Errorf("imposto_percentual inválido: %w", err)
	}
	impostoPercentual, err = valueobject.NewPercentage(impostosDecimal)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, err
	}

	// Parse comissão default
	comissaoDecimal, err := decimal.NewFromString(req.ComissaoDefault)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, fmt.Errorf("comissao_default inválido: %w", err)
	}
	comissaoDefault, err = valueobject.NewPercentage(comissaoDecimal)
	if err != nil {
		return valueobject.Percentage{}, valueobject.Percentage{}, valueobject.Percentage{}, decimal.Zero, err
	}

	return margemDesejada, impostoPercentual, comissaoDefault, markupAlvo, nil
}

// FromSimularPrecoRequest converte DTO Request para parâmetros do use case
func FromSimularPrecoRequest(req dto.SimularPrecoRequest) (
	custoMateriais, custoMaoDeObra, precoAtual valueobject.Money,
	parametros *entity.ParametrosSimulacao,
	err error,
) {
	// Parse custos
	custoMateriaisDecimal, err := decimal.NewFromString(req.CustoMateriais)
	if err != nil {
		return valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, nil, fmt.Errorf("custo_materiais inválido: %w", err)
	}
	custoMateriais = valueobject.NewMoneyFromDecimal(custoMateriaisDecimal)

	custoMaoObraDecimal, err := decimal.NewFromString(req.CustoMaoDeObra)
	if err != nil {
		return valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, nil, fmt.Errorf("custo_mao_de_obra inválido: %w", err)
	}
	custoMaoDeObra = valueobject.NewMoneyFromDecimal(custoMaoObraDecimal)

	precoAtualDecimal, err := decimal.NewFromString(req.PrecoAtual)
	if err != nil {
		return valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, nil, fmt.Errorf("preco_atual inválido: %w", err)
	}
	precoAtual = valueobject.NewMoneyFromDecimal(precoAtualDecimal)

	// Converter parâmetros opcionais de DTO para entidade
	if req.Parametros != nil {
		parametros = &entity.ParametrosSimulacao{
			TempoMedioMinutos:     req.Parametros.TempoMedioMinutos,
			QuantidadeMensal:      req.Parametros.QuantidadeMensal,
			CustoPorMinuto:        req.Parametros.CustoPorMinuto,
			ObservacoesAdicionais: req.Parametros.ObservacoesAdicionais,
		}
	}

	return custoMateriais, custoMaoDeObra, precoAtual, parametros, nil
}
