"use client";

import { Box } from '@mui/material';
import { PropsWithChildren } from 'react';
import { Sidebar } from './sidebar';
import { Topbar } from './topbar';

export function DashboardShell({ children }: PropsWithChildren) {
  return (
    <Box sx={{ display: 'grid', gridTemplateColumns: '260px 1fr', minHeight: '100vh', background: 'var(--valtaris-bg)' }}>
      <Sidebar />
      <Box component="main" sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
        <Topbar />
        <Box sx={{ p: 3, display: 'grid', gap: 3, flexGrow: 1, background: 'var(--valtaris-bg)' }}>{children}</Box>
      </Box>
    </Box>
  );
}
