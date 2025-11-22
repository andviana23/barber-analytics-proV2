import { getFluxoCaixaDiario } from '@/lib/services/fluxoService';
import { FluxoCaixaDiarioResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type FluxoCaixaDiarioResponse = z.infer<typeof FluxoCaixaDiarioResponseSchema>;

interface DateRangeFilter {
  dataInicio: string;
  dataFim: string;
}

export function useFluxoCaixaCompensado(
  dateRange: DateRangeFilter,
  options?: Omit<
    UseQueryOptions<FluxoCaixaDiarioResponse[], Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<FluxoCaixaDiarioResponse[], Error>({
    queryKey: [
      'fluxo-caixa-compensado',
      dateRange.dataInicio,
      dateRange.dataFim,
    ],
    queryFn: () => getFluxoCaixaDiario(dateRange.dataInicio, dateRange.dataFim),
    staleTime: 1000 * 60 * 3, // 3 minutos
    ...options,
  });
}
