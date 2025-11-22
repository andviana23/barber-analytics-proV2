"use client";

import { Skeleton, SkeletonProps, Stack } from '@mui/material';
import { PropsWithChildren } from 'react';

type VSkeletonProps = SkeletonProps &
  PropsWithChildren<{
    lines?: number;
  }>;

export function VSkeleton({ lines = 3, ...rest }: VSkeletonProps) {
  return (
    <Stack spacing={1}>
      {Array.from({ length: lines }).map((_, idx) => (
        <Skeleton
          key={idx}
          variant="rectangular"
          height={12}
          sx={{ borderRadius: '8px', backgroundColor: 'var(--valtaris-surface-subtle)' }}
          {...rest}
        />
      ))}
    </Stack>
  );
}
