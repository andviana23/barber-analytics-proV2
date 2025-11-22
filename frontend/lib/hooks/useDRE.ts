import { useQuery } from '@tanstack/react-query';
import {
  getDreMensal,
  listDRE,
  type ListDREInput,
} from '../services/dreService';

// ==================== DRE ====================

export function useGetDRE(mesAno: string) {
  return useQuery({
    queryKey: ['financial', 'dre', mesAno],
    queryFn: () => getDreMensal(mesAno),
    enabled: !!mesAno,
  });
}

export function useListDRE(params?: ListDREInput) {
  return useQuery({
    queryKey: ['financial', 'dre', 'list', params],
    queryFn: () => listDRE(params),
  });
}
