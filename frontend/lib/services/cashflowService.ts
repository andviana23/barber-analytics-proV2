import { z } from 'zod';
import { apiClient } from './client';
import { ISODate } from './schemas';

// ==================== SCHEMAS ====================

const FluxoCaixaDiarioResponseSchema = z.object({
  id: z.string(),
  data: z.string(),
  saldo_inicial: z.string(),
  entradas_confirmadas: z.string(),
  entradas_previstas: z.string(),
  saidas_pagas: z.string(),
  saidas_previstas: z.string(),
  saldo_final: z.string(),
  processado_em: z.string(),
});

const ListFluxoCaixaSchema = z.object({
  data_inicio: ISODate.optional(),
  data_fim: ISODate.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().max(100).optional(),
});

// ==================== TYPES ====================

export type FluxoCaixaDiarioResponse = z.infer<
  typeof FluxoCaixaDiarioResponseSchema
>;
export type ListFluxoCaixaInput = z.infer<typeof ListFluxoCaixaSchema>;

// ==================== ENDPOINTS (2 total) ====================

export async function getFluxoCaixa(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/financial/cashflow/${id}`,
    schema: FluxoCaixaDiarioResponseSchema,
  });
}

export async function listFluxoCaixa(params?: ListFluxoCaixaInput) {
  const query = ListFluxoCaixaSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/financial/cashflow',
    params: query,
    schema: z.object({
      data: z.array(FluxoCaixaDiarioResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}
