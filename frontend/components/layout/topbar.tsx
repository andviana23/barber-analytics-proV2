"use client";

import { AppBar, Box, IconButton, Toolbar, Typography, Tooltip } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import LightModeIcon from '@mui/icons-material/LightMode';
import { useThemeStore } from '../../store/theme';

export function Topbar() {
  const mode = useThemeStore((s) => s.mode);
  const toggle = useThemeStore((s) => s.toggle);

  return (
    <AppBar
      position="sticky"
      elevation={0}
      sx={{
        background: 'var(--valtaris-surface)',
        borderBottom: '1px solid var(--valtaris-border)',
        color: 'var(--valtaris-text)',
      }}
    >
      <Toolbar sx={{ display: 'flex', gap: 2 }}>
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Overview
        </Typography>
        <Box sx={{ flexGrow: 1 }} />
        <Tooltip title={mode === 'light' ? 'Ativar Dark' : 'Voltar para Light'}>
          <IconButton
            aria-label="Toggle theme"
            onClick={toggle}
            sx={{
              color: 'var(--valtaris-text)',
              border: '1px solid var(--valtaris-border)',
              borderRadius: '10px',
              '&:hover': { backgroundColor: 'var(--valtaris-surface-subtle)' },
            }}
          >
            {mode === 'light' ? <Brightness4Icon /> : <LightModeIcon />}
          </IconButton>
        </Tooltip>
      </Toolbar>
    </AppBar>
  );
}
