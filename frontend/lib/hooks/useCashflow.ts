import { useQuery } from '@tanstack/react-query';
import {
  getFluxoCaixa,
  listFluxoCaixa,
  type ListFluxoCaixaInput,
} from '../services/cashflowService';

// ==================== FLUXO DE CAIXA ====================

export function useGetFluxoCaixa(id: string) {
  return useQuery({
    queryKey: ['financial', 'cashflow', id],
    queryFn: () => getFluxoCaixa(id),
    enabled: !!id,
  });
}

export function useListFluxoCaixa(params?: ListFluxoCaixaInput) {
  return useQuery({
    queryKey: ['financial', 'cashflow', 'list', params],
    queryFn: () => listFluxoCaixa(params),
  });
}
