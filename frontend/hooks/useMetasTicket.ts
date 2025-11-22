import { getMetaTicket } from '@/lib/services/metasService';
import { MetaTicketResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type MetaTicketResponse = z.infer<typeof MetaTicketResponseSchema>;

export function useMetasTicket(
  mesAno: string,
  options?: Omit<
    UseQueryOptions<MetaTicketResponse[], Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<MetaTicketResponse[], Error>({
    queryKey: ['metas-ticket', mesAno],
    queryFn: () => getMetaTicket(mesAno),
    staleTime: 1000 * 60 * 5, // 5 minutos
    ...options,
  });
}
