import { simularPreco } from '@/lib/services/pricingService';
import { PrecificacaoSimulacaoResponseSchema } from '@/lib/services/schemas';
import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { z } from 'zod';

type PrecificacaoSimulacaoResponse = z.infer<
  typeof PrecificacaoSimulacaoResponseSchema
>;

interface SimularPrecoParams {
  itemId: string;
  tipoItem: string;
  custoMateriais: string;
  custoMaoDeObra: string;
  margemDesejada?: string;
  comissaoPercentual?: string;
  impostoPercentual?: string;
}

export function useSimularPreco(
  options?: Omit<
    UseMutationOptions<
      PrecificacaoSimulacaoResponse,
      Error,
      SimularPrecoParams
    >,
    'mutationFn'
  >
) {
  return useMutation<PrecificacaoSimulacaoResponse, Error, SimularPrecoParams>({
    mutationFn: (params) => simularPreco(params),
    ...options,
  });
}
