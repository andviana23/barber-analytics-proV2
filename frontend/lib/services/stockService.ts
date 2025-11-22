import { apiClient } from './client';

// Placeholder para futuras operações de estoque.
// Mantido para cumprir o contrato de serviços, sem implementar endpoints inexistentes.

export async function getEstoqueResumo() {
  return apiClient.request<{ message: string }>({
    method: 'GET',
    path: '/stock/summary',
  });
}

export async function registrarEntrada(payload: Record<string, unknown>) {
  return apiClient.request<{ message: string }>({
    method: 'POST',
    path: '/stock/entries',
    body: payload,
  });
}
