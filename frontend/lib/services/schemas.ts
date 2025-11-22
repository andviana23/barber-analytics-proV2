import { z } from 'zod';

export const MoneyString = z.string().regex(/^-?\d+(\.\d+)?$/, 'valor monetário inválido');
export const UUID = z.string().uuid();
export const ISODate = z.string().regex(/^\d{4}-\d{2}-\d{2}/, 'data inválida');
export const MesAno = z.string().regex(/^\d{4}-\d{2}$/, 'mes_ano inválido');

export const ErrorResponseSchema = z.object({
  error: z.string(),
  message: z.string(),
  details: z.record(z.unknown()).optional(),
});

export const SuccessResponseSchema = z.object({
  success: z.boolean(),
  message: z.string(),
  data: z.unknown().optional(),
});

export const ContaPagarResponseSchema = z.object({
  id: UUID,
  descricao: z.string(),
  categoria_id: UUID,
  fornecedor: z.string(),
  valor: MoneyString,
  tipo: z.string(),
  recorrente: z.boolean(),
  periodicidade: z.string().optional(),
  data_vencimento: ISODate,
  data_pagamento: z.string().nullable().optional(),
  status: z.string(),
  comprovante_url: z.string().nullable().optional(),
  pix_code: z.string().nullable().optional(),
  observacoes: z.string().nullable().optional(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});

export const ContaReceberResponseSchema = z.object({
  id: UUID,
  origem: z.string(),
  assinatura_id: UUID.nullable().optional(),
  descricao_origem: z.string(),
  valor: MoneyString,
  valor_pago: MoneyString,
  valor_aberto: MoneyString,
  data_vencimento: ISODate,
  data_recebimento: z.string().nullable().optional(),
  status: z.string(),
  observacoes: z.string().nullable().optional(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});

export const FluxoCaixaDiarioResponseSchema = z.object({
  id: UUID,
  data: ISODate,
  saldo_inicial: MoneyString,
  entradas_confirmadas: MoneyString,
  entradas_previstas: MoneyString,
  saidas_pagas: MoneyString,
  saidas_previstas: MoneyString,
  saldo_final: MoneyString,
  processado_em: z.string(),
});

export const CompensacaoBancariaResponseSchema = z.object({
  id: UUID,
  receita_id: UUID,
  data_transacao: ISODate,
  data_compensacao: ISODate,
  data_compensado: z.string().nullable().optional(),
  valor_bruto: MoneyString,
  taxa_percentual: MoneyString,
  taxa_fixa: MoneyString,
  valor_liquido: MoneyString,
  meio_pagamento_id: UUID,
  d_mais: z.number().int(),
  status: z.string(),
});

export const DREMensalResponseSchema = z.object({
  id: UUID,
  mes_ano: MesAno,
  receita_servicos: MoneyString,
  receita_produtos: MoneyString,
  receita_planos: MoneyString,
  receita_total: MoneyString,
  custo_comissoes: MoneyString,
  custo_insumos: MoneyString,
  custo_variavel_total: MoneyString,
  despesa_fixa: MoneyString,
  despesa_variavel: MoneyString,
  despesa_total: MoneyString,
  resultado_bruto: MoneyString,
  resultado_operacional: MoneyString,
  margem_bruta: MoneyString,
  margem_operacional: MoneyString,
  lucro_liquido: MoneyString,
  processado_em: z.string(),
});

export const PrecificacaoConfigResponseSchema = z.object({
  id: UUID,
  margem_desejada: z.string(),
  markup_alvo: z.string(),
  imposto_percentual: z.string(),
  comissao_percentual_default: z.string(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});

export const PrecificacaoSimulacaoResponseSchema = z.object({
  id: UUID,
  item_id: UUID,
  tipo_item: z.string(),
  custo_materiais: z.string(),
  custo_mao_de_obra: z.string(),
  custo_total: z.string(),
  margem_desejada: z.string(),
  comissao_percentual: z.string(),
  imposto_percentual: z.string(),
  preco_sugerido: z.string(),
  preco_atual: z.string(),
  diferenca_percentual: z.string(),
  lucro_estimado: z.string(),
  margem_final: z.string(),
  criado_em: z.string(),
});

export const MetaMensalResponseSchema = z.object({
  id: UUID,
  mes_ano: MesAno,
  meta_faturamento: MoneyString,
  origem: z.string(),
  status: z.string(),
  realizado: MoneyString,
  percentual: z.string(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});

export const MetaBarbeiroResponseSchema = z.object({
  id: UUID,
  barbeiro_id: UUID,
  mes_ano: MesAno,
  meta_servicos_gerais: MoneyString,
  meta_servicos_extras: MoneyString,
  meta_produtos: MoneyString,
  realizado_servicos_gerais: MoneyString,
  realizado_servicos_extras: MoneyString,
  realizado_produtos: MoneyString,
  percentual_servicos_gerais: z.string(),
  percentual_servicos_extras: z.string(),
  percentual_produtos: z.string(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});

export const MetaTicketResponseSchema = z.object({
  id: UUID,
  mes_ano: MesAno,
  tipo: z.string(),
  barbeiro_id: UUID.nullable().optional(),
  meta_valor: MoneyString,
  ticket_medio_realizado: MoneyString,
  percentual: z.string(),
  criado_em: z.string(),
  atualizado_em: z.string(),
});
