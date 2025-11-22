import { getDreMensal } from '@/lib/services/dreService';
import { DREMensalResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type DREMensalResponse = z.infer<typeof DREMensalResponseSchema>;

export function useDRE(
  mesAno: string,
  options?: Omit<
    UseQueryOptions<DREMensalResponse, Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<DREMensalResponse, Error>({
    queryKey: ['dre', mesAno],
    queryFn: () => getDreMensal(mesAno),
    staleTime: 1000 * 60 * 5, // 5 minutos
    ...options,
  });
}
