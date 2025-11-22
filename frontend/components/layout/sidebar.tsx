"use client";

import { Box, Divider, List, ListItemButton, ListItemIcon, ListItemText, Typography } from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';
import EventIcon from '@mui/icons-material/Event';
import PeopleIcon from '@mui/icons-material/People';
import Link from 'next/link';

const navItems = [
  { label: 'Dashboard', href: '/dashboard', icon: <DashboardIcon /> },
  { label: 'Financeiro', href: '/financeiro', icon: <AccountBalanceWalletIcon /> },
  { label: 'Assinaturas', href: '/assinaturas', icon: <PeopleIcon /> },
  { label: 'Agenda', href: '/agenda', icon: <EventIcon /> },
];

export function Sidebar() {
  return (
    <Box
      component="aside"
      sx={{
        width: 260,
        background: 'var(--valtaris-surface)',
        borderRight: '1px solid var(--valtaris-border)',
        minHeight: '100vh',
        position: 'sticky',
        top: 0,
        p: 2,
        display: 'flex',
        flexDirection: 'column',
        gap: 2,
      }}
    >
      <Typography variant="h6" sx={{ fontWeight: 700, letterSpacing: 0.5 }}>
        VALTARIS
      </Typography>
      <Divider />
      <List dense disablePadding>
        {navItems.map((item) => (
          <ListItemButton
            key={item.href}
            component={Link}
            href={item.href}
            sx={{
              borderRadius: '10px',
              mb: 0.5,
              '&:hover': { backgroundColor: 'var(--valtaris-surface-subtle)' },
            }}
          >
            <ListItemIcon sx={{ color: 'var(--valtaris-text-muted)', minWidth: 36 }}>{item.icon}</ListItemIcon>
            <ListItemText primary={item.label} />
          </ListItemButton>
        ))}
      </List>
    </Box>
  );
}
