# VALTARIS Design System v1.0

> Identidade visual: **Cyber Luxury** – precisão tecnológica, atmosfera de vidro/metal e accents de neon contidos.
> Objetivo: tornar o frontend do VALTARIS coerente, rápido de evoluir e fácil de manter, sem “componentes soltos”.

---

## 1. Visão Geral

O **VALTARIS Design System v1.0** é o conjunto de **tokens, temas, componentes e padrões** que padronizam toda a interface do produto.

Ele garante que:

- Todas as telas do sistema tenham a **mesma identidade visual** (“Cyber Luxury”).
- Todo novo componente use as **mesmas cores, espaçamentos, tipografia e comportamento**.
- O time consiga **evoluir o produto com segurança**, sem quebrar visual nem acessibilidade.
- Light e Dark mode funcionem de forma consistente, com **Light como tema padrão**.

Se você está criando ou alterando qualquer tela do frontend, **você é obrigado a seguir este Design System.**

---

## 2. Stack Oficial do Frontend VALTARIS

O frontend do VALTARIS é construído com o seguinte stack:

### 2.1 Core

- **Framework:** Next.js 16 (App Router)
- **Linguagem:** TypeScript
- **UI Runtime:** React 19

### 2.2 Design System & UI

- **Base de Componentes:** **MUI 5**
- **Estilização de Tema:** Emotion + MUI Theme
- **Variáveis de Tema:** CSS Variables `--valtaris-*` para cores, surfaces e materiais
- **Identidade Visual:** “Cyber Luxury” (vidro + metal + neon contido)
- **Temas:**
  - **Light** → tema padrão (opt-in para Dark)
  - **Dark** → ativado via Theme Store (Zustand) + ThemeRegistry (MUI + Emotion)

> Regra de ouro: **nunca hardcodar cores**. Sempre usar tokens do tema MUI ou CSS vars `--valtaris-*`.

### 2.3 Estado, Formulários e Dados

- **Estado de UI global:** Zustand (tema, preferências, layout, etc.)
- **Formulários:** React Hook Form + Zod (schema-first)
- **Dados remotos:** TanStack Query (fetch, cache, status, refetch)
- **Agenda / Scheduler:** DayPilot (ou similar) estilizado **exclusivamente via CSS Vars do VALTARIS** (sem CSS inline).

### 2.4 Testes & Qualidade

- **Testes de unidade:** Jest + Testing Library
- **Testes E2E:** Playwright em fluxos críticos (assinaturas, financeiro, agenda)
- **Acessibilidade:** foco em WCAG 2.1 – contraste, focus ring, hit area mínima, navegação por teclado.

---

## 3. Como o Design System é organizado

A documentação do Design System do VALTARIS está dividida nos seguintes arquivos:

1. **Foundations** → `01-FOUNDATIONS.md`
   Paleta (light/dark), tokens → MUI → CSS vars, tipografia, espaçamentos, materiais (vidro/metal), sombras, breakpoints.

2. **Architecture** → `02-ARCHITECTURE.md`
   ThemeRegistry (MUI + Emotion), Theme Store (Zustand), fluxo de light/dark, CSS vars para DayPilot/externos, ordem de providers e passo a passo para novas páginas.

3. **Components** → `03-COMPONENTS.md`
   Botões, inputs/selects, cards, modais (vidro), badges/status, tabelas (DataTable) – com exemplos em MUI usando tokens do VALTARIS.

4. **Patterns** → `04-PATTERNS.md`
   Padrões obrigatórios para formulários (RHF + Zod), acessibilidade, feedback (erro/sucesso/loading), layout responsivo e erros a evitar.

5. **Componentes Críticos** → `COMPONENTES_CRITICOS.md`
   Conjunto mínimo que **não pode quebrar**: ThemeProvider/Layout, formulários chave, hooks de dados, Scheduler/DayPilot, base de componentes do DS.

6. **Guia de Uso no Projeto** → `GUIA_FRONTEND.md`
   Guia prático de como criar telas usando o Design System, estrutura de pastas, convenções, exemplos de fluxo de desenvolvimento diário.

7. **Mapa de Telas** → `MAPA_TELAS.md`
   Lista de telas-chave do produto e como cada uma consome o Design System (cards, modais, tabelas, scheduler, forms), sempre em Light/Dark.

---

## 4. Regras Essenciais do Design System

### 4.1 Paleta e Tokens

- **Proibido** criar novas cores fora da paleta oficial do VALTARIS.
- Todas as cores devem vir de:
  - Tokens MUI (`theme.palette.*`)
  - CSS Variables `--valtaris-*` definidas em Foundations.

### 4.2 Light & Dark Mode

- Tema padrão: **Light**.
- Dark é **opt-in**, controlado pela Theme Store (Zustand) + ThemeRegistry.
- Toda página nova deve:
  - Respeitar o tema atual (light/dark) automaticamente.
  - Não depender de hex fixo ou “gambiarras” que quebrem no dark mode.

### 4.3 Componentes e Patterns

- Use SEMPRE os componentes documentados em `03-COMPONENTS.md`.
- Qualquer novo componente deve:
  - Usar tokens de cores, tipografia e espaçamento já existentes.
  - Ter estados claros: hover, focus, disabled, loading.
  - Manter a estética “Cyber Luxury” (vidro/metal/neon contido).

### 4.4 Formulários

- Todo formulário deve seguir o padrão:
  - **React Hook Form + Zod**
  - Componentes de campo (inputs/selects) baseados no Design System
  - Tratamento padrão de erro, sucesso e loading, conforme `04-PATTERNS.md`.

---

## 5. Começando: como criar uma nova tela usando o VALTARIS DS

Quando for criar uma nova tela:

1. **Leia rapidamente**:
   - `01-FOUNDATIONS.md`
   - `02-ARCHITECTURE.md`
   - `03-COMPONENTS.md`
   - `04-PATTERNS.md`
2. No código:
   - Use o layout padrão do dashboard.
   - Use apenas componentes do Design System (botões, inputs, cards, modais, badges, tabelas).
   - Conecte formulários com RHF + Zod.
   - Evite qualquer CSS inline que fuja dos tokens.

Se você está escrevendo código que não respeita esses pontos, muito provavelmente está **quebrando o Design System** e deve parar, ajustar ou discutir antes de seguir.

---

## 6. Benefícios práticos

Seguir o VALTARIS Design System traz:

- **Velocidade**: menos tempo discutindo UI, mais tempo focando em regra de negócio.
- **Consistência**: o produto parece uma coisa só, não um Frankenstein de telas.
- **Escala**: fácil desligar recursos, dividir trabalho entre devs e revisar PRs.
- **Confiabilidade**: menos bugs visuais, menos “modo dark quebrado”, mais previsibilidade.

O Design System não é opcional – ele é a base visual e de UX de todo o frontend do VALTARIS.
