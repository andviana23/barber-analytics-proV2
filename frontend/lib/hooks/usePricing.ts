import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  deleteConfig,
  deleteSimulacao,
  getConfig,
  getSimulacao,
  listSimulacoes,
  saveConfig,
  saveSimulacao,
  simularPreco,
  updateConfig,
  type ListSimulacoesInput,
} from '../services/pricingService';

// ==================== CONFIG ====================

export function useSaveConfig() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: saveConfig,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricing', 'config'] });
    },
  });
}

export function useGetConfig() {
  return useQuery({
    queryKey: ['pricing', 'config'],
    queryFn: getConfig,
  });
}

export function useUpdateConfig() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateConfig,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricing', 'config'] });
    },
  });
}

export function useDeleteConfig() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteConfig,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricing', 'config'] });
    },
  });
}

// ==================== SIMULAÇÃO ====================

export function useSimularPreco() {
  return useMutation({
    mutationFn: simularPreco,
  });
}

export function useSaveSimulacao() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: saveSimulacao,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricing', 'simulations'] });
    },
  });
}

export function useGetSimulacao(id: string) {
  return useQuery({
    queryKey: ['pricing', 'simulations', id],
    queryFn: () => getSimulacao(id),
    enabled: !!id,
  });
}

export function useListSimulacoes(params?: ListSimulacoesInput) {
  return useQuery({
    queryKey: ['pricing', 'simulations', 'list', params],
    queryFn: () => listSimulacoes(params),
  });
}

export function useDeleteSimulacao() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteSimulacao,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['pricing', 'simulations'] });
      queryClient.removeQueries({ queryKey: ['pricing', 'simulations', id] });
    },
  });
}
