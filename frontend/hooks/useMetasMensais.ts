import { getMetaMensal } from '@/lib/services/metasService';
import { MetaMensalResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type MetaMensalResponse = z.infer<typeof MetaMensalResponseSchema>;

export function useMetasMensais(
  mesAno: string,
  options?: Omit<
    UseQueryOptions<MetaMensalResponse, Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<MetaMensalResponse, Error>({
    queryKey: ['metas-mensais', mesAno],
    queryFn: () => getMetaMensal(mesAno),
    staleTime: 1000 * 60 * 5, // 5 minutos
    ...options,
  });
}
