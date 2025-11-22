import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createReceivable,
  deleteReceivable,
  getReceivable,
  listReceivables,
  marcarRecebimento,
  updateReceivable,
  type ListReceivablesInput,
  type MarcarRecebimentoInput,
  type UpdateReceivableInput,
} from '../services/receivablesService';

// ==================== CONTAS A RECEBER ====================

export function useCreateReceivable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createReceivable,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['financial', 'receivables'] });
    },
  });
}

export function useGetReceivable(id: string) {
  return useQuery({
    queryKey: ['financial', 'receivables', id],
    queryFn: () => getReceivable(id),
    enabled: !!id,
  });
}

export function useListReceivables(params?: ListReceivablesInput) {
  return useQuery({
    queryKey: ['financial', 'receivables', 'list', params],
    queryFn: () => listReceivables(params),
  });
}

export function useUpdateReceivable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateReceivableInput }) =>
      updateReceivable(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({
        queryKey: ['financial', 'receivables', id],
      });
      queryClient.invalidateQueries({
        queryKey: ['financial', 'receivables', 'list'],
      });
    },
  });
}

export function useDeleteReceivable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteReceivable,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['financial', 'receivables'] });
      queryClient.removeQueries({ queryKey: ['financial', 'receivables', id] });
    },
  });
}

export function useMarcarRecebimento() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      contaId,
      data,
    }: {
      contaId: string;
      data: MarcarRecebimentoInput;
    }) => marcarRecebimento(contaId, data),
    onSuccess: (_, { contaId }) => {
      queryClient.invalidateQueries({
        queryKey: ['financial', 'receivables', contaId],
      });
      queryClient.invalidateQueries({
        queryKey: ['financial', 'receivables', 'list'],
      });
      queryClient.invalidateQueries({ queryKey: ['financial', 'cashflow'] });
      queryClient.invalidateQueries({
        queryKey: ['financial', 'compensations'],
      });
    },
  });
}
