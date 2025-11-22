"use client";

import { Card, CardContent, CardProps, Typography, Box } from '@mui/material';
import { PropsWithChildren } from 'react';

type VCardProps = CardProps &
  PropsWithChildren<{
    title?: string;
    subtitle?: string;
  }>;

export function VCard({ title, subtitle, children, ...rest }: VCardProps) {
  return (
    <Card
      elevation={0}
      {...rest}
      sx={{
        background: 'var(--valtaris-surface)',
        border: '1px solid var(--valtaris-border)',
        borderRadius: '14px',
        boxShadow: '0 12px 30px rgba(0,0,0,0.12)',
        backdropFilter: 'blur(10px)',
        ...(rest.sx || {}),
      }}
    >
      <CardContent sx={{ p: 3, display: 'grid', gap: 1.5 }}>
        {(title || subtitle) && (
          <Box>
            {title && (
              <Typography variant="h6" sx={{ fontWeight: 600, letterSpacing: -0.2 }}>
                {title}
              </Typography>
            )}
            {subtitle && (
              <Typography variant="body2" color="text.secondary">
                {subtitle}
              </Typography>
            )}
          </Box>
        )}
        {children}
      </CardContent>
    </Card>
  );
}
