import createCache from '@emotion/cache';

// Emotion cache configured for MUI with prepend to ensure styles load first on the page.
export default function createEmotionCache() {
  return createCache({ key: 'mui', prepend: true });
}
