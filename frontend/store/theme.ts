"use client";

import { create } from 'zustand';

export type ThemeMode = 'light' | 'dark';

type ThemeState = {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
  toggle: () => void;
};

export const useThemeStore = create<ThemeState>((set, get) => ({
  mode: 'light',
  setMode: (mode) => {
    set({ mode });
    if (typeof window !== 'undefined') {
      window.localStorage.setItem('valtaris-theme', mode);
    }
  },
  toggle: () => {
    const next = get().mode === 'light' ? 'dark' : 'light';
    set({ mode: next });
    if (typeof window !== 'undefined') {
      window.localStorage.setItem('valtaris-theme', next);
    }
  },
}));
