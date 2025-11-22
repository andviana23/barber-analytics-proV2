import { registrarEntrada } from '@/lib/services/stockService';
import {
  useMutation,
  UseMutationOptions,
  useQueryClient,
} from '@tantml:react-query';

interface RegistrarEntradaInput {
  estoqueId: string;
  quantidade: number;
  custoUnitario: string;
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

export function useRegistrarEntrada(
  options?: Omit<
    UseMutationOptions<MovimentacaoResponse, Error, RegistrarEntradaInput>,
    'mutationFn'
  >
) {
  const queryClient = useQueryClient();

  return useMutation<MovimentacaoResponse, Error, RegistrarEntradaInput>({
    mutationFn: (input) => registrarEntrada(input),
    onSuccess: (data) => {
      // Invalida cache de estoque e movimentações
      queryClient.invalidateQueries({ queryKey: ['estoque'] });
      queryClient.invalidateQueries({ queryKey: ['movimentacoes'] });
      queryClient.invalidateQueries({ queryKey: ['estoque', data.estoque_id] });
    },
    ...options,
  });
}
