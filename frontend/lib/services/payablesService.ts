import { z } from 'zod';
import { apiClient } from './client';
import {
  ContaPagarResponseSchema,
  ErrorResponseSchema,
  ISODate,
  MoneyString,
  SuccessResponseSchema,
  UUID,
} from './schemas';

// ==================== SCHEMAS ====================

const CreatePayableSchema = z.object({
  descricao: z.string().min(3),
  categoria_id: UUID,
  fornecedor: z.string().min(1),
  valor: MoneyString,
  tipo: z.enum(['FIXO', 'VARIAVEL']),
  data_vencimento: ISODate,
  recorrente: z.boolean().optional(),
  periodicidade: z.string().optional(),
  pix_code: z.string().optional(),
  observacoes: z.string().optional(),
});

const UpdatePayableSchema = CreatePayableSchema.partial();

const ListPayablesSchema = z.object({
  categoria_id: UUID.optional(),
  status: z.enum(['PENDENTE', 'PAGO', 'VENCIDO', 'CANCELADO']).optional(),
  tipo: z.enum(['FIXO', 'VARIAVEL']).optional(),
  data_vencimento_inicio: ISODate.optional(),
  data_vencimento_fim: ISODate.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().optional(),
});

const MarcarPagamentoSchema = z.object({
  data_pagamento: ISODate,
  comprovante_url: z.string().optional(),
});

// ==================== TYPES ====================

export type CreatePayableInput = z.infer<typeof CreatePayableSchema>;
export type UpdatePayableInput = z.infer<typeof UpdatePayableSchema>;
export type ListPayablesInput = z.infer<typeof ListPayablesSchema>;
export type MarcarPagamentoInput = z.infer<typeof MarcarPagamentoSchema>;

// ==================== ENDPOINTS (6 total) ====================

export async function createPayable(input: CreatePayableInput) {
  const payload = CreatePayableSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/financial/payables',
    body: payload,
    schema: ContaPagarResponseSchema,
  });
}

export async function getPayable(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/financial/payables/${id}`,
    schema: ContaPagarResponseSchema,
  });
}

export async function listPayables(params?: ListPayablesInput) {
  const query = ListPayablesSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/financial/payables',
    params: query,
    schema: z.object({
      data: z.array(ContaPagarResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function updatePayable(id: string, input: UpdatePayableInput) {
  const payload = UpdatePayableSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: `/financial/payables/${id}`,
    body: payload,
    schema: ContaPagarResponseSchema,
  });
}

export async function deletePayable(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/financial/payables/${id}`,
    schema: SuccessResponseSchema,
  });
}

export async function marcarPagamento(
  contaId: string,
  input: MarcarPagamentoInput
) {
  const payload = MarcarPagamentoSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: `/financial/payables/${contaId}/payment`,
    body: payload,
    schema: SuccessResponseSchema,
  });
}

export const PayablesErrorSchema = ErrorResponseSchema;
