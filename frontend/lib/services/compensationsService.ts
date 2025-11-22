import { z } from 'zod';
import { apiClient } from './client';
import { ISODate, SuccessResponseSchema } from './schemas';

// ==================== SCHEMAS ====================

const CompensacaoBancariaResponseSchema = z.object({
  id: z.string(),
  receita_id: z.string(),
  data_transacao: z.string(),
  data_compensacao: z.string(),
  data_compensado: z.string().optional().nullable(),
  valor_bruto: z.string(),
  taxa_percentual: z.string(),
  taxa_fixa: z.string(),
  valor_liquido: z.string(),
  meio_pagamento_id: z.string(),
  d_mais: z.number().int(),
  status: z.enum(['PENDENTE', 'COMPENSADO']),
});

const ListCompensacoesSchema = z.object({
  status: z.enum(['PENDENTE', 'COMPENSADO']).optional(),
  data_inicio: ISODate.optional(),
  data_fim: ISODate.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().max(100).optional(),
});

// ==================== TYPES ====================

export type CompensacaoBancariaResponse = z.infer<
  typeof CompensacaoBancariaResponseSchema
>;
export type ListCompensacoesInput = z.infer<typeof ListCompensacoesSchema>;

// ==================== ENDPOINTS (3 total) ====================

export async function getCompensacao(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/financial/compensations/${id}`,
    schema: CompensacaoBancariaResponseSchema,
  });
}

export async function listCompensacoes(params?: ListCompensacoesInput) {
  const query = ListCompensacoesSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/financial/compensations',
    params: query,
    schema: z.object({
      data: z.array(CompensacaoBancariaResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function deleteCompensacao(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/financial/compensations/${id}`,
    schema: SuccessResponseSchema,
  });
}
