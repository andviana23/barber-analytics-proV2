package dto

// SetMetaMensalRequest representa a requisição para definir meta mensal
type SetMetaMensalRequest struct {
	MesAno          string `json:"mes_ano" validate:"required,len=7"` // YYYY-MM
	MetaFaturamento string `json:"meta_faturamento" validate:"required"`
	Origem          string `json:"origem" validate:"required,oneof=MANUAL AUTOMATICA"`
}

// MetaMensalResponse representa a resposta de meta mensal
type MetaMensalResponse struct {
	ID              string `json:"id"`
	MesAno          string `json:"mes_ano"`
	MetaFaturamento string `json:"meta_faturamento"`
	Origem          string `json:"origem"`
	Status          string `json:"status"`
	Realizado       string `json:"realizado"`
	Percentual      string `json:"percentual"`
	CriadoEm        string `json:"criado_em"`
	AtualizadoEm    string `json:"atualizado_em"`
}

// SetMetaBarbeiroRequest representa a requisição para definir meta de barbeiro
type SetMetaBarbeiroRequest struct {
	BarbeiroID         string `json:"barbeiro_id" validate:"required,uuid"`
	MesAno             string `json:"mes_ano" validate:"required,len=7"`
	MetaServicosGerais string `json:"meta_servicos_gerais" validate:"required"`
	MetaServicosExtras string `json:"meta_servicos_extras" validate:"required"`
	MetaProdutos       string `json:"meta_produtos" validate:"required"`
}

// MetaBarbeiroResponse representa a resposta de meta de barbeiro
type MetaBarbeiroResponse struct {
	ID                       string `json:"id"`
	BarbeiroID               string `json:"barbeiro_id"`
	MesAno                   string `json:"mes_ano"`
	MetaServicosGerais       string `json:"meta_servicos_gerais"`
	MetaServicosExtras       string `json:"meta_servicos_extras"`
	MetaProdutos             string `json:"meta_produtos"`
	RealizadoServicosGerais  string `json:"realizado_servicos_gerais"`
	RealizadoServicosExtras  string `json:"realizado_servicos_extras"`
	RealizadoProdutos        string `json:"realizado_produtos"`
	PercentualServicosGerais string `json:"percentual_servicos_gerais"`
	PercentualServicosExtras string `json:"percentual_servicos_extras"`
	PercentualProdutos       string `json:"percentual_produtos"`
	CriadoEm                 string `json:"criado_em"`
	AtualizadoEm             string `json:"atualizado_em"`
}

// SetMetaTicketRequest representa a requisição para definir meta de ticket médio
type SetMetaTicketRequest struct {
	MesAno     string  `json:"mes_ano" validate:"required,len=7"`
	Tipo       string  `json:"tipo" validate:"required,oneof=GERAL BARBEIRO"`
	BarbeiroID *string `json:"barbeiro_id,omitempty" validate:"omitempty,uuid"`
	MetaValor  string  `json:"meta_valor" validate:"required"`
}

// MetaTicketResponse representa a resposta de meta de ticket médio
type MetaTicketResponse struct {
	ID                   string  `json:"id"`
	MesAno               string  `json:"mes_ano"`
	Tipo                 string  `json:"tipo"`
	BarbeiroID           *string `json:"barbeiro_id,omitempty"`
	MetaValor            string  `json:"meta_valor"`
	TicketMedioRealizado string  `json:"ticket_medio_realizado"`
	Percentual           string  `json:"percentual"`
	CriadoEm             string  `json:"criado_em"`
	AtualizadoEm         string  `json:"atualizado_em"`
}
