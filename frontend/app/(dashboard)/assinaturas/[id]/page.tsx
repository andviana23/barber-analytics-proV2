"use client";

import { useParams } from 'next/navigation';
import { PageHeader } from '../../../../components/layout/page-header';
import { VCard } from '../../../../components/ui/val-card';

export default function AssinaturaDetalhePage() {
  const params = useParams();
  const id = Array.isArray(params?.id) ? params?.id[0] : params?.id;

  return (
    <>
      <PageHeader title="Detalhe da Assinatura" subtitle={`Assinatura ${id ?? ''}`} />
      <VCard title="Em construção" subtitle="Use o DS para billing, histórico e ações" />
    </>
  );
}
