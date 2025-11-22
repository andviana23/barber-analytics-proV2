import { listMovimentacoes } from '@/lib/services/stockService';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

interface Movimentacao {
  id: string;
  estoque_id: string;
  tipo: 'ENTRADA' | 'SAIDA' | 'AJUSTE';
  quantidade: number;
  quantidade_anterior: number;
  quantidade_nova: number;
  custo_unitario: string;
  valor_total: string;
  motivo: string;
  observacoes?: string;
  criado_por: string;
  criado_em: string;
}

interface MovimentacoesFilter {
  estoqueId?: string;
  tipo?: string;
  dataInicio?: string;
  dataFim?: string;
  page?: number;
  pageSize?: number;
}

export function useMovimentacoes(
  filters: MovimentacoesFilter = {},
  options?: Omit<UseQueryOptions<Movimentacao[], Error>, 'queryKey' | 'queryFn'>
) {
  return useQuery<Movimentacao[], Error>({
    queryKey: ['movimentacoes', filters],
    queryFn: () => listMovimentacoes(filters),
    staleTime: 1000 * 60 * 2, // 2 minutos
    ...options,
  });
}
