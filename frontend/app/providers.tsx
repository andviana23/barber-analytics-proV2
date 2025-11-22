"use client";

import { PropsWithChildren, useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeRegistry } from './theme-registry';

export function Providers({ children }: PropsWithChildren) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <ThemeRegistry>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </ThemeRegistry>
  );
}
