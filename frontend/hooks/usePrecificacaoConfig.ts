import { getPrecificacaoConfig } from '@/lib/services/pricingService';
import { PrecificacaoConfigResponseSchema } from '@/lib/services/schemas';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { z } from 'zod';

type PrecificacaoConfigResponse = z.infer<
  typeof PrecificacaoConfigResponseSchema
>;

export function usePrecificacaoConfig(
  options?: Omit<
    UseQueryOptions<PrecificacaoConfigResponse, Error>,
    'queryKey' | 'queryFn'
  >
) {
  return useQuery<PrecificacaoConfigResponse, Error>({
    queryKey: ['precificacao-config'],
    queryFn: () => getPrecificacaoConfig(),
    staleTime: 1000 * 60 * 10, // 10 minutos (config muda raramente)
    ...options,
  });
}
