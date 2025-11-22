import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  deleteCompensacao,
  getCompensacao,
  listCompensacoes,
  type ListCompensacoesInput,
} from '../services/compensationsService';

// ==================== COMPENSAÇÕES BANCÁRIAS ====================

export function useGetCompensacao(id: string) {
  return useQuery({
    queryKey: ['financial', 'compensations', id],
    queryFn: () => getCompensacao(id),
    enabled: !!id,
  });
}

export function useListCompensacoes(params?: ListCompensacoesInput) {
  return useQuery({
    queryKey: ['financial', 'compensations', 'list', params],
    queryFn: () => listCompensacoes(params),
  });
}

export function useDeleteCompensacao() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteCompensacao,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({
        queryKey: ['financial', 'compensations'],
      });
      queryClient.removeQueries({
        queryKey: ['financial', 'compensations', id],
      });
    },
  });
}
