import { createContaPagar } from '@/lib/services/payablesService';
import { ContaPagarResponseSchema } from '@/lib/services/schemas';
import {
  useMutation,
  UseMutationOptions,
  useQueryClient,
} from '@tanstack/react-query';
import { z } from 'zod';

type ContaPagarResponse = z.infer<typeof ContaPagarResponseSchema>;

interface CreateContaPagarInput {
  descricao: string;
  categoriaId: string;
  fornecedor: string;
  valor: string;
  tipo: 'FIXO' | 'VARIAVEL';
  recorrente: boolean;
  periodicidade?: string;
  dataVencimento: string;
  observacoes?: string;
}

export function useCreateContaPagar(
  options?: Omit<
    UseMutationOptions<ContaPagarResponse, Error, CreateContaPagarInput>,
    'mutationFn'
  >
) {
  const queryClient = useQueryClient();

  return useMutation<ContaPagarResponse, Error, CreateContaPagarInput>({
    mutationFn: (input) => createContaPagar(input),
    onSuccess: () => {
      // Invalida cache de contas a pagar
      queryClient.invalidateQueries({ queryKey: ['contas-pagar'] });
      queryClient.invalidateQueries({ queryKey: ['fluxo-caixa-compensado'] });
    },
    ...options,
  });
}
