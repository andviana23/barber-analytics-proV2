import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  deleteMetaBarbeiro,
  deleteMetaMensal,
  deleteMetaTicket,
  getMetaBarbeiro,
  getMetaMensal,
  getMetaTicket,
  listMetasBarbeiro,
  listMetasMensais,
  listMetasTicket,
  setMetaBarbeiro,
  setMetaMensal,
  setMetaTicket,
  updateMetaBarbeiro,
  updateMetaMensal,
  updateMetaTicket,
  type ListMetasBarbeiroInput,
  type ListMetasMensaisInput,
  type ListMetasTicketInput,
  type UpdateMetaBarbeiroInput,
  type UpdateMetaMensalInput,
  type UpdateMetaTicketInput,
} from '../services/metasService';

// ==================== META MENSAL ====================

export function useSetMetaMensal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setMetaMensal,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'monthly'] });
    },
  });
}

export function useGetMetaMensal(id: string) {
  return useQuery({
    queryKey: ['metas', 'monthly', id],
    queryFn: () => getMetaMensal(id),
    enabled: !!id,
  });
}

export function useListMetasMensais(params?: ListMetasMensaisInput) {
  return useQuery({
    queryKey: ['metas', 'monthly', 'list', params],
    queryFn: () => listMetasMensais(params),
  });
}

export function useUpdateMetaMensal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMetaMensalInput }) =>
      updateMetaMensal(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'monthly', id] });
      queryClient.invalidateQueries({ queryKey: ['metas', 'monthly', 'list'] });
    },
  });
}

export function useDeleteMetaMensal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteMetaMensal,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'monthly'] });
      queryClient.removeQueries({ queryKey: ['metas', 'monthly', id] });
    },
  });
}

// ==================== META BARBEIRO ====================

export function useSetMetaBarbeiro() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setMetaBarbeiro,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'barbers'] });
    },
  });
}

export function useGetMetaBarbeiro(id: string) {
  return useQuery({
    queryKey: ['metas', 'barbers', id],
    queryFn: () => getMetaBarbeiro(id),
    enabled: !!id,
  });
}

export function useListMetasBarbeiro(params?: ListMetasBarbeiroInput) {
  return useQuery({
    queryKey: ['metas', 'barbers', 'list', params],
    queryFn: () => listMetasBarbeiro(params),
  });
}

export function useUpdateMetaBarbeiro() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMetaBarbeiroInput }) =>
      updateMetaBarbeiro(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'barbers', id] });
      queryClient.invalidateQueries({ queryKey: ['metas', 'barbers', 'list'] });
    },
  });
}

export function useDeleteMetaBarbeiro() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteMetaBarbeiro,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'barbers'] });
      queryClient.removeQueries({ queryKey: ['metas', 'barbers', id] });
    },
  });
}

// ==================== META TICKET MÉDIO ====================

export function useSetMetaTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setMetaTicket,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'ticket'] });
    },
  });
}

export function useGetMetaTicket(id: string) {
  return useQuery({
    queryKey: ['metas', 'ticket', id],
    queryFn: () => getMetaTicket(id),
    enabled: !!id,
  });
}

export function useListMetasTicket(params?: ListMetasTicketInput) {
  return useQuery({
    queryKey: ['metas', 'ticket', 'list', params],
    queryFn: () => listMetasTicket(params),
  });
}

export function useUpdateMetaTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMetaTicketInput }) =>
      updateMetaTicket(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'ticket', id] });
      queryClient.invalidateQueries({ queryKey: ['metas', 'ticket', 'list'] });
    },
  });
}

export function useDeleteMetaTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteMetaTicket,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['metas', 'ticket'] });
      queryClient.removeQueries({ queryKey: ['metas', 'ticket', id] });
    },
  });
}
