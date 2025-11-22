"use client";

import { Button, ButtonProps, alpha, useTheme } from '@mui/material';

type Kind = 'primary' | 'secondary' | 'ghost';

type VButtonProps = ButtonProps & {
  kind?: Kind;
};

export function VButton({ kind = 'primary', ...props }: VButtonProps) {
  const theme = useTheme();
  const map = {
    primary: {
      bg: 'var(--valtaris-primary)',
      hover: 'var(--valtaris-primary-dark)',
      text: '#FFFFFF',
      border: 'transparent',
    },
    secondary: {
      bg: 'var(--valtaris-surface-subtle)',
      hover: 'var(--valtaris-surface)',
      text: 'var(--valtaris-text)',
      border: 'var(--valtaris-border)',
    },
    ghost: {
      bg: 'transparent',
      hover: alpha(theme.palette.primary.main, 0.08),
      text: 'var(--valtaris-text)',
      border: alpha('#8A7CFF', 0.3),
    },
  }[kind];

  return (
    <Button
      {...props}
      variant="contained"
      sx={{
        px: 2.5,
        py: 1.25,
        minHeight: 40,
        borderRadius: '10px',
        textTransform: 'none',
        fontWeight: 600,
        letterSpacing: 0.4,
        backgroundColor: map.bg,
        color: map.text,
        border: `1px solid ${map.border}`,
        boxShadow: '0 10px 30px rgba(62, 91, 255, 0.12)',
        '&:hover': {
          backgroundColor: map.hover,
          boxShadow: '0 14px 36px rgba(34, 211, 238, 0.18)',
        },
        '&:focus-visible': {
          outline: '2px solid var(--valtaris-accent-aqua)',
          outlineOffset: 2,
        },
        '&.Mui-disabled': {
          opacity: 0.48,
          boxShadow: 'none',
        },
      }}
    />
  );
}
