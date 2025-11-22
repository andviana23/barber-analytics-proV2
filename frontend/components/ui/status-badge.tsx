"use client";

import { Chip, ChipProps } from '@mui/material';

type Status = 'success' | 'warning' | 'danger';

const variants: Record<Status, { bg: string; color: string; label: string }> = {
  success: { bg: 'rgba(56,214,155,0.14)', color: '#38D69B', label: 'Sucesso' },
  warning: { bg: 'rgba(244,178,62,0.16)', color: '#F4B23E', label: 'Atenção' },
  danger: { bg: 'rgba(239,68,68,0.16)', color: '#EF4444', label: 'Alerta' },
};

type StatusBadgeProps = Omit<ChipProps, 'label' | 'color'> & { status: Status; label?: string };

export function StatusBadge({ status, label, ...rest }: StatusBadgeProps) {
  const variant = variants[status];
  return (
    <Chip
      {...rest}
      label={label ?? variant.label}
      size="small"
      sx={{
        bgcolor: variant.bg,
        color: variant.color,
        borderRadius: '999px',
        fontWeight: 600,
        letterSpacing: 0.2,
      }}
    />
  );
}
