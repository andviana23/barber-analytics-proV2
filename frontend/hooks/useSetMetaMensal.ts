import { setMetaMensal } from '@/lib/services/metasService';
import { MetaMensalResponseSchema } from '@/lib/services/schemas';
import {
  useMutation,
  UseMutationOptions,
  useQueryClient,
} from '@tanstack/react-query';
import { z } from 'zod';

type MetaMensalResponse = z.infer<typeof MetaMensalResponseSchema>;

interface SetMetaMensalInput {
  mesAno: string;
  metaFaturamento: string;
  origem: 'MANUAL' | 'AUTOMATICA';
}

export function useSetMetaMensal(
  options?: Omit<
    UseMutationOptions<MetaMensalResponse, Error, SetMetaMensalInput>,
    'mutationFn'
  >
) {
  const queryClient = useQueryClient();

  return useMutation<MetaMensalResponse, Error, SetMetaMensalInput>({
    mutationFn: (input) => setMetaMensal(input),
    onSuccess: (data) => {
      // Invalida cache de metas mensais
      queryClient.invalidateQueries({
        queryKey: ['metas-mensais', data.mes_ano],
      });
      queryClient.invalidateQueries({ queryKey: ['metas-mensais'] });
    },
    ...options,
  });
}
