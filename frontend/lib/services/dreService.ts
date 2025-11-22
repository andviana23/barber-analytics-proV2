import { z } from 'zod';
import { apiClient } from './client';
import { DREMensalResponseSchema, MesAno } from './schemas';

// ==================== SCHEMAS ====================

const ListDRESchema = z.object({
  mes_ano_inicio: MesAno.optional(),
  mes_ano_fim: MesAno.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().max(100).optional(),
});

// ==================== TYPES ====================

export type ListDREInput = z.infer<typeof ListDRESchema>;

// ==================== ENDPOINTS (2 total) ====================

export async function getDreMensal(mesAno: string) {
  MesAno.parse(mesAno);
  return apiClient.request({
    method: 'GET',
    path: `/financial/dre/${mesAno}`,
    schema: DREMensalResponseSchema,
  });
}

export async function listDRE(params?: ListDREInput) {
  const query = ListDRESchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/financial/dre',
    params: query,
    schema: z.object({
      data: z.array(DREMensalResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}
