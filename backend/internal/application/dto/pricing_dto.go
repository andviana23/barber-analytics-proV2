package dto

// SaveConfigPrecificacaoRequest representa a requisição para salvar configuração
type SaveConfigPrecificacaoRequest struct {
	MargemDesejada    string `json:"margem_desejada" validate:"required"`
	MarkupAlvo        string `json:"markup_alvo" validate:"required"`
	ImpostoPercentual string `json:"imposto_percentual" validate:"required"`
	ComissaoDefault   string `json:"comissao_default" validate:"required"`
}

// PrecificacaoConfigResponse representa a resposta de configuração
type PrecificacaoConfigResponse struct {
	ID                        string `json:"id"`
	MargemDesejada            string `json:"margem_desejada"`
	MarkupAlvo                string `json:"markup_alvo"`
	ImpostoPercentual         string `json:"imposto_percentual"`
	ComissaoPercentualDefault string `json:"comissao_percentual_default"`
	CriadoEm                  string `json:"criado_em"`
	AtualizadoEm              string `json:"atualizado_em"`
}

// SimularPrecoRequest representa a requisição para simular preço
type SimularPrecoRequest struct {
	ItemID         string                  `json:"item_id" validate:"required,uuid"`
	TipoItem       string                  `json:"tipo_item" validate:"required,oneof=SERVICO PRODUTO"`
	CustoMateriais string                  `json:"custo_materiais" validate:"required"`
	CustoMaoDeObra string                  `json:"custo_mao_de_obra" validate:"required"`
	PrecoAtual     string                  `json:"preco_atual" validate:"required"`
	Parametros     *ParametrosSimulacaoDTO `json:"parametros,omitempty"`
}

// ParametrosSimulacaoDTO representa parâmetros adicionais da simulação
type ParametrosSimulacaoDTO struct {
	TempoMedioMinutos     int     `json:"tempo_medio_minutos,omitempty"`
	QuantidadeMensal      int     `json:"quantidade_mensal,omitempty"`
	CustoPorMinuto        float64 `json:"custo_por_minuto,omitempty"`
	ObservacoesAdicionais string  `json:"observacoes_adicionais,omitempty"`
}

// PrecificacaoSimulacaoResponse representa a resposta de simulação
type PrecificacaoSimulacaoResponse struct {
	ID                  string `json:"id"`
	ItemID              string `json:"item_id"`
	TipoItem            string `json:"tipo_item"`
	CustoMateriais      string `json:"custo_materiais"`
	CustoMaoDeObra      string `json:"custo_mao_de_obra"`
	CustoTotal          string `json:"custo_total"`
	MargemDesejada      string `json:"margem_desejada"`
	ComissaoPercentual  string `json:"comissao_percentual"`
	ImpostoPercentual   string `json:"imposto_percentual"`
	PrecoSugerido       string `json:"preco_sugerido"`
	PrecoAtual          string `json:"preco_atual"`
	DiferencaPercentual string `json:"diferenca_percentual"`
	LucroEstimado       string `json:"lucro_estimado"`
	MargemFinal         string `json:"margem_final"`
	CriadoEm            string `json:"criado_em"`
}
