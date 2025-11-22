import { z } from 'zod';
import { apiClient } from './client';
import {
  MoneyString,
  PrecificacaoConfigResponseSchema,
  PrecificacaoSimulacaoResponseSchema,
  SuccessResponseSchema,
  UUID,
} from './schemas';

// ==================== SCHEMAS ====================

// Config
const SaveConfigSchema = z.object({
  margem_desejada: z.string(),
  markup_alvo: z.string(),
  imposto_percentual: z.string(),
  comissao_default: z.string(),
});

const UpdateConfigSchema = SaveConfigSchema.partial();

// Simulação
const SimularPrecoSchema = z.object({
  item_id: UUID,
  tipo_item: z.enum(['SERVICO', 'PRODUTO']),
  custo_materiais: MoneyString,
  custo_mao_de_obra: MoneyString,
  preco_atual: MoneyString,
  parametros: z
    .object({
      tempo_medio_minutos: z.number().int().optional(),
      quantidade_mensal: z.number().int().optional(),
      custo_por_minuto: z.number().optional(),
      observacoes_adicionais: z.string().optional(),
    })
    .optional(),
});

const SaveSimulacaoSchema = z.object({
  item_id: UUID,
  tipo_item: z.enum(['SERVICO', 'PRODUTO']),
  custo_materiais: MoneyString,
  custo_mao_de_obra: MoneyString,
  preco_atual: MoneyString,
  preco_sugerido: MoneyString,
  margem_calculada: z.string(),
  markup_calculado: z.string(),
  break_even: MoneyString,
  observacoes: z.string().optional(),
});

const ListSimulacoesSchema = z.object({
  tipo_item: z.enum(['SERVICO', 'PRODUTO']).optional(),
  page: z.number().int().positive().optional(),
  page_size: z.number().int().positive().optional(),
});

// ==================== TYPES ====================

export type SaveConfigInput = z.infer<typeof SaveConfigSchema>;
export type UpdateConfigInput = z.infer<typeof UpdateConfigSchema>;

export type SimularPrecoInput = z.infer<typeof SimularPrecoSchema>;
export type SaveSimulacaoInput = z.infer<typeof SaveSimulacaoSchema>;
export type ListSimulacoesInput = z.infer<typeof ListSimulacoesSchema>;

// ==================== CONFIG (4 endpoints) ====================

export async function saveConfig(input: SaveConfigInput) {
  const payload = SaveConfigSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/pricing/config',
    body: payload,
    schema: PrecificacaoConfigResponseSchema,
  });
}

export async function getConfig() {
  return apiClient.request({
    method: 'GET',
    path: '/pricing/config',
    schema: PrecificacaoConfigResponseSchema,
  });
}

export async function updateConfig(input: UpdateConfigInput) {
  const payload = UpdateConfigSchema.parse(input);
  return apiClient.request({
    method: 'PUT',
    path: '/pricing/config',
    body: payload,
    schema: PrecificacaoConfigResponseSchema,
  });
}

export async function deleteConfig() {
  return apiClient.request({
    method: 'DELETE',
    path: '/pricing/config',
    schema: SuccessResponseSchema,
  });
}

// ==================== SIMULAÇÃO (5 endpoints) ====================

export async function simularPreco(input: SimularPrecoInput) {
  const payload = SimularPrecoSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/pricing/simulate',
    body: payload,
    schema: PrecificacaoSimulacaoResponseSchema,
  });
}

export async function saveSimulacao(input: SaveSimulacaoInput) {
  const payload = SaveSimulacaoSchema.parse(input);
  return apiClient.request({
    method: 'POST',
    path: '/pricing/simulations',
    body: payload,
    schema: PrecificacaoSimulacaoResponseSchema,
  });
}

export async function getSimulacao(id: string) {
  return apiClient.request({
    method: 'GET',
    path: `/pricing/simulations/${id}`,
    schema: PrecificacaoSimulacaoResponseSchema,
  });
}

export async function listSimulacoes(params?: ListSimulacoesInput) {
  const query = ListSimulacoesSchema.parse(params || {});
  return apiClient.request({
    method: 'GET',
    path: '/pricing/simulations',
    params: query,
    schema: z.object({
      data: z.array(PrecificacaoSimulacaoResponseSchema),
      page: z.number(),
      page_size: z.number(),
      total: z.number(),
    }),
  });
}

export async function deleteSimulacao(id: string) {
  return apiClient.request({
    method: 'DELETE',
    path: `/pricing/simulations/${id}`,
    schema: SuccessResponseSchema,
  });
}
