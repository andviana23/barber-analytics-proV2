import { listContasReceber } from '@/lib/services/receivablesService';
import { ContaReceberResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type ContaReceberResponse = z.infer<typeof ContaReceberResponseSchema>;

interface ReceivablesFilter {
  status?: string;
  origem?: string;
  dataInicio?: string;
  dataFim?: string;
  page?: number;
  pageSize?: number;
}

export function useContasReceber(
  filters: ReceivablesFilter = {},
  options?: Omit<
    UseQueryOptions<ContaReceberResponse[], Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<ContaReceberResponse[], Error>({
    queryKey: ['contas-receber', filters],
    queryFn: () => listContasReceber(filters),
    staleTime: 1000 * 60 * 2, // 2 minutos
    ...options,
  });
}
