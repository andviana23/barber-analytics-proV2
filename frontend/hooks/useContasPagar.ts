import { listContasPagar } from '@/lib/services/payablesService';
import { ContaPagarResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type ContaPagarResponse = z.infer<typeof ContaPagarResponseSchema>;

interface PayablesFilter {
  status?: string;
  dataInicio?: string;
  dataFim?: string;
  categoria?: string;
  page?: number;
  pageSize?: number;
}

export function useContasPagar(
  filters: PayablesFilter = {},
  options?: Omit<
    UseQueryOptions<ContaPagarResponse[], Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<ContaPagarResponse[], Error>({
    queryKey: ['contas-pagar', filters],
    queryFn: () => listContasPagar(filters),
    staleTime: 1000 * 60 * 2, // 2 minutos
    ...options,
  });
}
