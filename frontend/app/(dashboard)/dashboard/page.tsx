"use client";

import { VCard } from '../../../components/ui/val-card';
import { PageHeader } from '../../../components/layout/page-header';
import { VTable } from '../../../components/ui/val-table';
import { StatusBadge } from '../../../components/ui/status-badge';
import { VSkeleton } from '../../../components/ui/val-skeleton';

const sampleRows = [
  { id: '1', name: 'Receita Serviços', value: 'R$ 12.300', status: <StatusBadge status="success" /> },
  { id: '2', name: 'Despesas Variáveis', value: 'R$ 4.100', status: <StatusBadge status="warning" /> },
  { id: '3', name: 'Inadimplência', value: 'R$ 900', status: <StatusBadge status="danger" /> },
];

export default function DashboardPage() {
  return (
    <>
      <PageHeader title="Dashboard" subtitle="Resumo financeiro em tempo real" />
      <VCard title="Visão Geral">
        <VSkeleton lines={1} />
      </VCard>
      <VCard title="Indicadores">
        <VTable
          columns={[
            { id: 'name', label: 'Indicador' },
            { id: 'value', label: 'Valor', align: 'right' },
            { id: 'status', label: 'Status', align: 'right' },
          ]}
          rows={sampleRows}
        />
      </VCard>
    </>
  );
}
