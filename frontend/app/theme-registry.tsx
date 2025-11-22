"use client";

import { PropsWithChildren, useEffect } from 'react';
import { CacheProvider } from '@emotion/react';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { useServerInsertedHTML } from 'next/navigation';
import createEmotionCache from '../lib/theme/createEmotionCache';
import { createValtarisTheme } from '../lib/theme/createValtarisTheme';
import { useThemeStore } from '../store/theme';

const cache = createEmotionCache();

export function ThemeRegistry({ children }: PropsWithChildren) {
  const mode = useThemeStore((s) => s.mode);
  const setMode = useThemeStore((s) => s.setMode);
  const theme = createValtarisTheme(mode);

  useServerInsertedHTML(() => (
    <style
      data-emotion={`${cache.key} ${Object.keys(cache.inserted).join(' ')}`}
      dangerouslySetInnerHTML={{ __html: Object.values(cache.inserted).join(' ') }}
    />
  ));

  useEffect(() => {
    const saved = typeof window !== 'undefined' ? window.localStorage.getItem('valtaris-theme') : null;
    if (saved === 'dark' || saved === 'light') {
      setMode(saved);
    }
  }, [setMode]);

  useEffect(() => {
    document.body.classList.toggle('theme-dark', mode === 'dark');
  }, [mode]);

  return (
    <CacheProvider value={cache}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </ThemeProvider>
    </CacheProvider>
  );
}
