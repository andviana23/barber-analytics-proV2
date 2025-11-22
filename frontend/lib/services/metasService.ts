import { z } from 'zod';
import { apiClient } from './client';
import {
  MesAno,
  MetaBarbeiroResponseSchema,
  MetaMensalResponseSchema,
  MetaTicketResponseSchema,
  MoneyString,
  SuccessResponseSchema,
  UUID,
} from './schemas';

// ==================== SCHEMAS ====================

// Meta Mensal
const SetMetaMensalSchema = z.object({
  mes_ano: MesAno,
  meta_faturamento: MoneyString,
  origem: z.enum(['MANUAL', 'AUTOMATICA']),
});

const UpdateMetaMensalSchema = SetMetaMensalSchema.partial().required({
  mes_ano: true,
});

const ListMetasMensaisSchema = z.object({
  mes_ano_inicio: MesAno.optional(),
  mes_ano_fim: MesAno.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().optional(),
});

// Meta Barbeiro
const SetMetaBarbeiroSchema = z.object({
  barbeiro_id: UUID,
  mes_ano: MesAno,
  meta_servicos_gerais: MoneyString,
  meta_servicos_extras: MoneyString,
  meta_produtos: MoneyString,
});

const UpdateMetaBarbeiroSchema = SetMetaBarbeiroSchema.partial().required({
  mes_ano: true,
});

const ListMetasBarbeiroSchema = z.object({
  barbeiro_id: UUID.optional(),
  mes_ano: MesAno.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().optional(),
});

// Meta Ticket Médio
const SetMetaTicketSchema = z.object({
  mes_ano: MesAno,
  tipo: z.enum(['GERAL', 'BARBEIRO']),
  barbeiro_id: UUID.optional().nullable(),
  meta_valor: MoneyString,
});

const UpdateMetaTicketSchema = SetMetaTicketSchema.partial().required({
  mes_ano: true,
});

const ListMetasTicketSchema = z.object({
  tipo: z.enum(['GERAL', 'BARBEIRO']).optional(),
  barbeiro_id: UUID.optional(),
  mes_ano: MesAno.optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().optional(),
});

// ==================== TYPES ====================

export type SetMetaMensalInput = z.infer<typeof SetMetaMensalSchema>;
export type UpdateMetaMensalInput = z.infer<typeof UpdateMetaMensalSchema>;
export type ListMetasMensaisInput = z.infer<typeof ListMetasMensaisSchema>;

export type SetMetaBarbeiroInput = z.infer<typeof SetMetaBarbeiroSchema>;
export type UpdateMetaBarbeiroInput = z.infer<typeof UpdateMetaBarbeiroSchema>;
export type ListMetasBarbeiroInput = z.infer<typeof ListMetasBarbeiroSchema>;

export type SetMetaTicketInput = z.infer<typeof SetMetaTicketSchema>;
export type UpdateMetaTicketInput = z.infer<typeof UpdateMetaTicketSchema>;
export type ListMetasTicketInput = z.infer<typeof ListMetasTicketSchema>;

// ==================== META MENSAL (5 endpoints) ====================

export async function setMetaMensal(input: SetMetaMensalInput) {
  const payload = SetMetaMensalSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/metas/monthly',
    body: payload,
    schema: MetaMensalResponseSchema,
  });
}

export async function getMetaMensal(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/metas/monthly/${id}`,
    schema: MetaMensalResponseSchema,
  });
}

export async function listMetasMensais(params?: ListMetasMensaisInput) {
  const query = ListMetasMensaisSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/metas/monthly',
    params: query,
    schema: z.object({
      data: z.array(MetaMensalResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function updateMetaMensal(
  id: string,
  input: UpdateMetaMensalInput
) {
  const payload = UpdateMetaMensalSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: `/metas/monthly/${id}`,
    body: payload,
    schema: MetaMensalResponseSchema,
  });
}

export async function deleteMetaMensal(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/metas/monthly/${id}`,
    schema: SuccessResponseSchema,
  });
}

// ==================== META BARBEIRO (5 endpoints) ====================

export async function setMetaBarbeiro(input: SetMetaBarbeiroInput) {
  const payload = SetMetaBarbeiroSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/metas/barbers',
    body: payload,
    schema: MetaBarbeiroResponseSchema,
  });
}

export async function getMetaBarbeiro(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/metas/barbers/${id}`,
    schema: MetaBarbeiroResponseSchema,
  });
}

export async function listMetasBarbeiro(params?: ListMetasBarbeiroInput) {
  const query = ListMetasBarbeiroSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/metas/barbers',
    params: query,
    schema: z.object({
      data: z.array(MetaBarbeiroResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function updateMetaBarbeiro(
  id: string,
  input: UpdateMetaBarbeiroInput
) {
  const payload = UpdateMetaBarbeiroSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: `/metas/barbers/${id}`,
    body: payload,
    schema: MetaBarbeiroResponseSchema,
  });
}

export async function deleteMetaBarbeiro(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/metas/barbers/${id}`,
    schema: SuccessResponseSchema,
  });
}

// ==================== META TICKET MÉDIO (5 endpoints) ====================

export async function setMetaTicket(input: SetMetaTicketInput) {
  const payload = SetMetaTicketSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/metas/ticket',
    body: payload,
    schema: MetaTicketResponseSchema,
  });
}

export async function getMetaTicket(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/metas/ticket/${id}`,
    schema: MetaTicketResponseSchema,
  });
}

export async function listMetasTicket(params?: ListMetasTicketInput) {
  const query = ListMetasTicketSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/metas/ticket',
    params: query,
    schema: z.object({
      data: z.array(MetaTicketResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function updateMetaTicket(
  id: string,
  input: UpdateMetaTicketInput
) {
  const payload = UpdateMetaTicketSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: `/metas/ticket/${id}`,
    body: payload,
    schema: MetaTicketResponseSchema,
  });
}

export async function deleteMetaTicket(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/metas/ticket/${id}`,
    schema: SuccessResponseSchema,
  });
}
