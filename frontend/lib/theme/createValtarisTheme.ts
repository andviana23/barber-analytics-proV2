import { createTheme, alpha } from '@mui/material/styles';
import { PaletteOptions } from '@mui/material';

export type ThemeMode = 'light' | 'dark';

const lightPalette: PaletteOptions = {
  mode: 'light',
  primary: { main: '#3E5BFF', dark: '#2A42D9' },
  background: { default: '#F5F7FA', paper: '#FFFFFF' },
  text: { primary: '#0E1015', secondary: '#707784' },
  divider: '#E2E5EC',
  success: { main: '#38D69B' },
  error: { main: '#EF4444' },
  warning: { main: '#F4B23E' },
};

const darkPalette: PaletteOptions = {
  mode: 'dark',
  primary: { main: '#3E5BFF', dark: '#2A42D9' },
  background: { default: '#0B0D12', paper: '#12141C' },
  text: { primary: '#FFFFFF', secondary: '#C1C4CE' },
  divider: '#2F343F',
  success: { main: '#3CE9A8' },
  error: { main: '#FF5C5C' },
  warning: { main: '#F6C65C' },
};

export function createValtarisTheme(mode: ThemeMode) {
  const palette = mode === 'dark' ? darkPalette : lightPalette;

  return createTheme({
    palette,
    spacing: 8,
    shape: {
      borderRadius: 10,
    },
    typography: {
      fontFamily: '"Space Grotesk", "Inter", "Segoe UI", system-ui, -apple-system, sans-serif',
      fontWeightRegular: 400,
      fontWeightMedium: 500,
      fontWeightBold: 600,
      h1: { fontWeight: 600, fontSize: '32px', lineHeight: 1.25, letterSpacing: '-0.25px' },
      h2: { fontWeight: 600, fontSize: '28px', lineHeight: 1.3, letterSpacing: '-0.25px' },
      h3: { fontWeight: 600, fontSize: '24px', lineHeight: 1.3 },
      h4: { fontWeight: 600, fontSize: '20px', lineHeight: 1.4 },
      body1: { fontSize: '16px', lineHeight: 1.5 },
      body2: { fontSize: '14px', lineHeight: 1.5 },
      button: { fontWeight: 600, letterSpacing: '0.4px', textTransform: 'none' },
    },
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          ':root': {
            '--valtaris-primary': '#3E5BFF',
            '--valtaris-primary-dark': '#2A42D9',
            '--valtaris-bg': '#F5F7FA',
            '--valtaris-surface': '#FFFFFF',
            '--valtaris-surface-subtle': '#F0F2F6',
            '--valtaris-border': '#E2E5EC',
            '--valtaris-text': '#0E1015',
            '--valtaris-text-muted': '#707784',
            '--valtaris-text-soft': '#A4AAB5',
            '--valtaris-accent-purple': '#8A7CFF',
            '--valtaris-accent-aqua': '#22D3EE',
            '--valtaris-accent-gold': '#D4AF37',
            '--valtaris-success': '#38D69B',
            '--valtaris-danger': '#EF4444',
            '--valtaris-warning': '#F4B23E',
          },
          '.theme-dark': {
            '--valtaris-primary': '#3E5BFF',
            '--valtaris-primary-dark': '#2A42D9',
            '--valtaris-bg': '#0B0D12',
            '--valtaris-surface': '#12141C',
            '--valtaris-surface-subtle': '#1A1D26',
            '--valtaris-border': '#2F343F',
            '--valtaris-text': '#FFFFFF',
            '--valtaris-text-muted': '#C1C4CE',
            '--valtaris-text-soft': '#8B919E',
            '--valtaris-accent-purple': '#9E8FFF',
            '--valtaris-accent-aqua': '#3EE6FF',
            '--valtaris-accent-gold': '#CDA43A',
            '--valtaris-success': '#3CE9A8',
            '--valtaris-danger': '#FF5C5C',
            '--valtaris-warning': '#F6C65C',
          },
          body: {
            backgroundColor: palette.background?.default,
            color: palette.text?.primary,
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 10,
            textTransform: 'none',
            fontWeight: 600,
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            borderRadius: 14,
            border: `1px solid ${alpha(palette.divider || '#E2E5EC', 1)}`,
          },
        },
      },
      MuiDialog: {
        styleOverrides: {
          paper: {
            borderRadius: 16,
          },
        },
      },
    },
  });
}
