# VALTARIS v1.0 — Frontend Architecture (Design System)

> Esta documentação explica **como o Design System do VALTARIS é ligado ao código**: ThemeRegistry, MUI, Emotion, Zustand (tema) e CSS vars `--valtaris-*`.

---

## 1. Stack Oficial do Frontend VALTARIS

### 1.1 Core

- Next.js 16 (App Router)
- React 19
- TypeScript

### 1.2 Design System & UI

- **MUI 5** (base de componentes)
- Emotion (estilos + cache)
- CSS Variables `--valtaris-*` para cores, surfaces e materiais
- Tema Light/Dark sincronizado com:
  - MUI Theme
  - CSS vars globais

### 1.3 Estado, Formulários, Dados

- Zustand (Theme Store, preferências, layout)
- React Hook Form + Zod (forms)
- TanStack Query (dados remotos)
- DayPilot (agenda/scheduler) estilizado via CSS vars

---

## 2. ThemeRegistry (MUI + Emotion)

O **ThemeRegistry** é o ponto central que integra:

- Cache do Emotion (para SSR e estilo correto)
- ThemeProvider do MUI
- Injeção de CSS vars `--valtaris-*` para Light/Dark

Responsabilidades:

1. Criar o **MUI Theme** correto com base no tema atual (`light` ou `dark`).
2. Aplicar `<ThemeProvider>` em volta da aplicação.
3. Injetar estilos globais necessários (ex.: reset, fontes, CSS vars).
4. Garantir que o **SSR do Next.js 16** respeite o estilo (via `useServerInsertedHTML` ou padrão equivalente recomendado).

---

## 3. Theme Store (Zustand)

A Theme Store (Zustand) controla:

- `theme`: `"light"` | `"dark"`
- Preferências do usuário (ex.: seguir sistema operacional ou não)
- Ações:
  - `setTheme("light" | "dark")`
  - `toggleTheme()`

Fluxo típico:

1. O estado inicial é **`"light"`**.
2. Ao usuário alternar o tema, a store atualiza `theme`.
3. O ThemeRegistry reage a essa mudança e:
   - Atualiza o MUI Theme.
   - Atualiza as CSS vars `--valtaris-*`.
   - Atualiza classes no `body` (ex.: `.theme-dark`).

---

## 4. Light/Dark Mode — fluxo padrão opt-in

Regras:

- **Light** é o tema padrão (sempre).
- **Dark** é opt-in (usuário escolhe).
- Nenhum componente deve supor “tema fixo”.
  Todos devem usar:
  - Tokens MUI (`theme.palette.*`)
  - Ou `var(--valtaris-*)`

Fluxo:

1. App inicia em **Light**.
2. Se o usuário tiver preferência salva (ex.: em localStorage), Theme Store aplica logo na inicialização.
3. Botão de alternância de tema dispara `toggleTheme()`.
4. ThemeRegistry sincroniza MUI Theme + CSS vars.

---

## 5. CSS Variables e Integração com DayPilot / externos

Para qualquer biblioteca externa (ex.: DayPilot, grids, charts, etc.):

- É **proibido** usar CSS inline com hex fixo.
- A estilização deve usar:

```css
:root {
  --valtaris-primary: #3e5bff;
  --valtaris-bg: #f5f7fa;
  --valtaris-surface: #ffffff;
  /* ... */
}

body.theme-dark {
  --valtaris-bg: #0b0d12;
  --valtaris-surface: #12141c;
  /* ... */
}
```
