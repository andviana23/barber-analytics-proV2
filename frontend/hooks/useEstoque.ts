import { listEstoque } from '@/lib/services/stockService';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

interface EstoqueItem {
  id: string;
  nome: string;
  codigo: string;
  categoria: string;
  quantidade_atual: number;
  quantidade_minima: number;
  unidade_medida: string;
  custo_unitario: string;
  valor_total: string;
  status: string;
  criado_em: string;
  atualizado_em: string;
}

interface EstoqueFilter {
  categoria?: string;
  status?: string;
  abaixoMinimo?: boolean;
  page?: number;
  pageSize?: number;
}

export function useEstoque(
  filters: EstoqueFilter = {},
  options?: Omit<UseQueryOptions<EstoqueItem[], Error>, 'queryKey' | 'queryFn'>
) {
  return useQuery<EstoqueItem[], Error>({
    queryKey: ['estoque', filters],
    queryFn: () => listEstoque(filters),
    staleTime: 1000 * 60 * 3, // 3 minutos
    ...options,
  });
}
