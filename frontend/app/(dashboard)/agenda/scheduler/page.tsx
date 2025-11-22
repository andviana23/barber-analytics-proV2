"use client";

import { PageHeader } from '../../../../components/layout/page-header';
import { VCard } from '../../../../components/ui/val-card';

export default function SchedulerPage() {
  return (
    <>
      <PageHeader title="Scheduler" subtitle="DayPilot tematizado com CSS vars VALTARIS" />
      <VCard title="Em construção" subtitle="Use a classe daypilot-valtaris e evite CSS inline">
        <div
          className="daypilot-valtaris"
          style={{
            border: '1px dashed var(--valtaris-border)',
            borderRadius: 12,
            padding: 16,
            color: 'var(--valtaris-text-muted)',
          }}
        >
          Scheduler aqui (aplicar classe e tokens).
        </div>
      </VCard>
    </>
  );
}
