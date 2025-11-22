"use client";

import { Dialog, DialogTitle, DialogContent, IconButton, DialogProps } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { PropsWithChildren } from 'react';

type VModalProps = DialogProps &
  PropsWithChildren<{
    title?: string;
    onClose?: () => void;
  }>;

export function VModal({ title, onClose, children, ...rest }: VModalProps) {
  return (
    <Dialog
      {...rest}
      onClose={onClose}
      slotProps={{
        backdrop: {
          sx: {
            backdropFilter: 'blur(16px)',
            backgroundColor: 'rgba(11,13,18,0.65)',
          },
        },
        paper: {
          sx: {
            borderRadius: '16px',
            border: '1px solid rgba(255,255,255,0.08)',
            background: 'var(--valtaris-surface)',
            boxShadow: '0 24px 60px rgba(0,0,0,0.5)',
            minWidth: 420,
          },
        },
      }}
    >
      {title && (
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {title}
          <IconButton
            aria-label="Fechar"
            onClick={onClose}
            sx={{
              ml: 'auto',
              color: 'var(--valtaris-text)',
              '&:hover': { color: 'var(--valtaris-primary)' },
            }}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        </DialogTitle>
      )}
      <DialogContent sx={{ pb: 3 }}>{children}</DialogContent>
    </Dialog>
  );
}
