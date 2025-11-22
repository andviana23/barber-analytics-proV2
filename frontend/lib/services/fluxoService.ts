import { apiClient } from './client';
import {
  FluxoCaixaDiarioResponseSchema,
  DREMensalResponseSchema,
  MesAno,
  ISODate,
} from './schemas';
import { z } from 'zod';

const GetFluxoParamsSchema = z.object({
  data: ISODate.optional(),
  from: ISODate.optional(),
  to: ISODate.optional(),
});

export type GetFluxoParams = z.infer<typeof GetFluxoParamsSchema>;

export async function getFluxoDiario(params: GetFluxoParams = {}) {
  const { data } = GetFluxoParamsSchema.parse(params);
  const query = data ? `?date=${encodeURIComponent(data)}` : '';
  return apiClient.request({
    method: 'GET',
    path: `/financial/cashflow${query}`,
    schema: FluxoCaixaDiarioResponseSchema,
  });
}

export async function getDreMensal(mesAno: string) {
  MesAno.parse(mesAno);
  const query = `?mes_ano=${encodeURIComponent(mesAno)}`;
  return apiClient.request({
    method: 'GET',
    path: `/financial/dre${query}`,
    schema: DREMensalResponseSchema,
  });
}
