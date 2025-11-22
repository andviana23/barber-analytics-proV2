import { useMutation, useQuery, useQueryClient } from '@tantml:react-query';
import {
  createPayable,
  deletePayable,
  getPayable,
  listPayables,
  marcarPagamento,
  updatePayable,
  type ListPayablesInput,
  type MarcarPagamentoInput,
  type UpdatePayableInput,
} from '../services/payablesService';

// ==================== CONTAS A PAGAR ====================

export function useCreatePayable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createPayable,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['financial', 'payables'] });
    },
  });
}

export function useGetPayable(id: string) {
  return useQuery({
    queryKey: ['financial', 'payables', id],
    queryFn: () => getPayable(id),
    enabled: !!id,
  });
}

export function useListPayables(params?: ListPayablesInput) {
  return useQuery({
    queryKey: ['financial', 'payables', 'list', params],
    queryFn: () => listPayables(params),
  });
}

export function useUpdatePayable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePayableInput }) =>
      updatePayable(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({
        queryKey: ['financial', 'payables', id],
      });
      queryClient.invalidateQueries({
        queryKey: ['financial', 'payables', 'list'],
      });
    },
  });
}

export function useDeletePayable() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deletePayable,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['financial', 'payables'] });
      queryClient.removeQueries({ queryKey: ['financial', 'payables', id] });
    },
  });
}

export function useMarcarPagamento() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      contaId,
      data,
    }: {
      contaId: string;
      data: MarcarPagamentoInput;
    }) => marcarPagamento(contaId, data),
    onSuccess: (_, { contaId }) => {
      queryClient.invalidateQueries({
        queryKey: ['financial', 'payables', contaId],
      });
      queryClient.invalidateQueries({
        queryKey: ['financial', 'payables', 'list'],
      });
      queryClient.invalidateQueries({ queryKey: ['financial', 'cashflow'] });
    },
  });
}
