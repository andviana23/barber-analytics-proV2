import { registrarSaida } from '@/lib/services/stockService';
import {
  useMutation,
  UseMutationOptions,
  useQueryClient,
} from '@tanstack/react-query';

interface RegistrarSaidaInput {
  estoqueId: string;
  quantidade: number;
  motivo: string;
  observacoes?: string;
}

interface MovimentacaoResponse {
  id: string;
  estoque_id: string;
  tipo: string;
  quantidade: number;
  quantidade_anterior: number;
  quantidade_nova: number;
  custo_unitario: string;
  valor_total: string;
  motivo: string;
  observacoes?: string;
  criado_em: string;
}

export function useRegistrarSaida(
  options?: Omit<
    UseMutationOptions<MovimentacaoResponse, Error, RegistrarSaidaInput>,
    'mutationFn'
  >
) {
  const queryClient = useQueryClient();

  return useMutation<MovimentacaoResponse, Error, RegistrarSaidaInput>({
    mutationFn: (input) => registrarSaida(input),
    onSuccess: (data) => {
      // Invalida cache de estoque e movimentações
      queryClient.invalidateQueries({ queryKey: ['estoque'] });
      queryClient.invalidateQueries({ queryKey: ['movimentacoes'] });
      queryClient.invalidateQueries({ queryKey: ['estoque', data.estoque_id] });
    },
    ...options,
  });
}
