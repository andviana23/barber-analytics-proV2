import { z } from 'zod';
import { apiClient } from './client';
import {
  ContaReceberResponseSchema,
  ErrorResponseSchema,
  ISODate,
  MoneyString,
  SuccessResponseSchema,
  UUID,
} from './schemas';

// ==================== SCHEMAS ====================

const CreateReceivableSchema = z.object({
  origem: z.string().min(1),
  assinatura_id: UUID.optional().nullable(),
  descricao_origem: z.string().min(1),
  valor: MoneyString,
  data_vencimento: ISODate,
  observacoes: z.string().optional(),
});

const UpdateReceivableSchema = z.object({
  origem: z.string().min(1).optional(),
  assinatura_id: UUID.optional().nullable(),
  descricao_origem: z.string().min(1).optional(),
  valor: MoneyString.optional(),
  data_vencimento: ISODate.optional(),
  observacoes: z.string().optional(),
});

const ListReceivablesSchema = z.object({
  status: z.enum(['PENDENTE', 'PARCIAL', 'RECEBIDO', 'CANCELADO']).optional(),
  origem: z.string().optional(),
  data_inicio: ISODate.optional(),
  data_fim: ISODate.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().max(100).optional(),
});

const MarcarRecebimentoSchema = z.object({
  valor_pago: MoneyString,
  data_recebimento: ISODate,
});

// ==================== TYPES ====================

export type CreateReceivableInput = z.infer<typeof CreateReceivableSchema>;
export type UpdateReceivableInput = z.infer<typeof UpdateReceivableSchema>;
export type ListReceivablesInput = z.infer<typeof ListReceivablesSchema>;
export type MarcarRecebimentoInput = z.infer<typeof MarcarRecebimentoSchema>;

// ==================== ENDPOINTS (6 total) ====================

export async function createReceivable(input: CreateReceivableInput) {
  const payload = CreateReceivableSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/financial/receivables',
    body: payload,
    schema: ContaReceberResponseSchema,
  });
}

export async function getReceivable(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/financial/receivables/${id}`,
    schema: ContaReceberResponseSchema,
  });
}

export async function listReceivables(params?: ListReceivablesInput) {
  const query = ListReceivablesSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/financial/receivables',
    params: query,
    schema: z.object({
      data: z.array(ContaReceberResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function updateReceivable(
  id: string,
  input: UpdateReceivableInput
) {
  const payload = UpdateReceivableSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: `/financial/receivables/${id}`,
    body: payload,
    schema: ContaReceberResponseSchema,
  });
}

export async function deleteReceivable(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/financial/receivables/${id}`,
    schema: SuccessResponseSchema,
  });
}

export async function marcarRecebimento(
  contaId: string,
  input: MarcarRecebimentoInput
) {
  const payload = MarcarRecebimentoSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: `/financial/receivables/${contaId}/receipt`,
    body: payload,
    schema: SuccessResponseSchema,
  });
}

export const ReceivablesErrorSchema = ErrorResponseSchema;
