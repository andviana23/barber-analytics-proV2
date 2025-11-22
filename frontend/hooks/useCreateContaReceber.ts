import { createContaReceber } from '@/lib/services/receivablesService';
import { ContaReceberResponseSchema } from '@/lib/services/schemas';
import {
  useMutation,
  UseMutationOptions,
  useQueryClient,
} from '@tanstack/react-query';
import { z } from 'zod';

type ContaReceberResponse = z.infer<typeof ContaReceberResponseSchema>;

interface CreateContaReceberInput {
  origem: 'SERVICO' | 'PRODUTO' | 'PLANO';
  assinaturaId?: string;
  descricaoOrigem: string;
  valor: string;
  dataVencimento: string;
  observacoes?: string;
}

export function useCreateContaReceber(
  options?: Omit<
    UseMutationOptions<ContaReceberResponse, Error, CreateContaReceberInput>,
    'mutationFn'
  >
) {
  const queryClient = useQueryClient();

  return useMutation<ContaReceberResponse, Error, CreateContaReceberInput>({
    mutationFn: (input) => createContaReceber(input),
    onSuccess: () => {
      // Invalida cache de contas a receber
      queryClient.invalidateQueries({ queryKey: ['contas-receber'] });
      queryClient.invalidateQueries({ queryKey: ['fluxo-caixa-compensado'] });
    },
    ...options,
  });
}
