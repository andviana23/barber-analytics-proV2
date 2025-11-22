import { getMetaBarbeiro } from '@/lib/services/metasService';
import { MetaBarbeiroResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type MetaBarbeiroResponse = z.infer<typeof MetaBarbeiroResponseSchema>;

export function useMetasBarbeiro(
  mesAno: string,
  barbeiroId: string,
  options?: Omit<
    UseQueryOptions<MetaBarbeiroResponse, Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<MetaBarbeiroResponse, Error>({
    queryKey: ['metas-barbeiro', mesAno, barbeiroId],
    queryFn: () => getMetaBarbeiro(mesAno, barbeiroId),
    staleTime: 1000 * 60 * 5, // 5 minutos
    enabled: !!barbeiroId,
    ...options,
  });
}
