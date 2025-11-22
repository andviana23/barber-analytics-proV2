# VALTARIS v1.0 — Foundations

> Esta é a base visual do VALTARIS: cores, tipografia, espaçamentos, materiais (vidro/metal), sombras e breakpoints.
> Se você está mexendo com UI, começa por aqui.

---

## 1. Identidade “Cyber Luxury”

A identidade visual do VALTARIS é definida como **Cyber Luxury**:

- **Precisão tecnológica** – tudo parece calculado, alinhado, sem excessos.
- **Materiais híbridos** – mistura de vidro (blur), metal (bordas finas) e neon (accents controlados).
- **Sofisticação contida** – nada é chamativo demais; os destaques são cirúrgicos.

Essa identidade deve guiar todas as decisões de cor, tipografia e layout.

---

## 2. Paleta Oficial (Tokens)

A paleta de cores do VALTARIS é **imutável**.
Você **não pode criar novos hex**, apenas usar os tokens existentes.

Todos os tokens são expostos via **tema MUI** e replicados como **CSS Variables `--valtaris-*`**, para uso em:

- Componentes MUI (`theme.palette.*`)
- Estilização de componentes externos (ex.: DayPilot)
- Camadas utilitárias (ex.: classes globais)

### 2.1 Light — “Cyber Luxury Light” (tema padrão)

Tema padrão do sistema.

- `primary` `#3E5BFF` · `primaryDark` `#2A42D9`
- `background.default` `#F5F7FA` · `paper` `#FFFFFF` · `subtle` `#F0F2F6`
- `border` `#E2E5EC`
- `text.primary` `#0E1015` · `text.secondary` `#707784` · `text.tertiary` `#A4AAB5`
- Accents:
  - `purple` `#8A7CFF`
  - `aqua` `#22D3EE`
  - `gold` `#D4AF37`
- Status:
  - `success` `#38D69B`
  - `danger` `#EF4444`
  - `warning` `#F4B23E`

> Light deve ser sempre o **tema inicial**. Dark é opt-in e nunca substitui o Light como base de design.

### 2.2 Dark — “Cyber Luxury Dark”

O tema Dark espelha a mesma paleta, mas reposiciona fundo, surfaces e contraste para ambientes escuros:

- Fundos mais profundos, com contraste forte entre background e surfaces.
- Texto sempre com contraste adequado (mínimo 4.5:1).
- Accents (aqua/purple/gold) usados com moderação, principalmente para estados ativos e foco.

> As cores base são **as mesmas da paleta oficial**; apenas os mapeamentos de `background`, `paper` e `text` mudam no MUI Theme + CSS vars.

---

## 3. Tokens → MUI → CSS Variables

Os tokens de cor passam por três camadas:

1. **Paleta base** (esta seção): hexes oficiais.
2. **MUI Theme**:
   - `theme.palette.primary.main = #3E5BFF`
   - `theme.palette.background.default = #F5F7FA`
   - `theme.palette.text.primary = #0E1015`
   - etc.
3. **CSS Variables** (exemplo):
   - `--valtaris-primary: #3E5BFF;`
   - `--valtaris-bg: #F5F7FA;`
   - `--valtaris-surface: #FFFFFF;`
   - `--valtaris-border: #E2E5EC;`
   - `--valtaris-text: #0E1015;`
   - `--valtaris-accent-aqua: #22D3EE;`

Regra prática:

- **No MUI**, use `theme.palette.*`.
- **Fora do MUI**, use `var(--valtaris-*)`.

---

## 4. Tipografia

### 4.1 Fontes

- **Primária:** `Space Grotesk` (ou equivalente geométrica)
  Usada em títulos, métricas, botões e elementos que exigem presença.

- **Monoespaçada (opcional):** `JetBrains Mono`
  Para códigos, IDs técnicos, números críticos ou labels em dashboards.

### 4.2 Hierarquia

Sugestão de níveis (podem ser mapeados para `Typography` do MUI):

- `h1` – Títulos de página / visão global.
- `h2` – Seções principais da página.
- `h3` – Blocos internos (cards, agrupamentos).
- `subtitle1` / `subtitle2` – Metadados e descrições curtas.
- `body1` – Texto padrão.
- `body2` – Texto auxiliar, rótulos longos, mensagens de sistema.
- `caption` – Labels discretos, legendas e tooltips.

Princípios:

- Sempre garantir legibilidade em Light e Dark.
- Evitar mais de 3 níveis diferentes na mesma tela.

---

## 5. Espaçamento e Grid

O VALTARIS usa uma base de espaçamento de **4px**:

- Unidade base: `4px`
- Escala sugerida: `4, 8, 12, 16, 20, 24, 32, 40, 48px`

Guidelines:

- Cards e seções principais: padding interno de 20–24px.
- Blocos de formulário: 16–24px entre grupos de campos.
- Distância entre cards em dashboards: 16–24px.

Grid:

- Em `md+`: preferir layout com **2 ou 3 colunas** de cards, evitando listas infinitas.
- Em `sm`: usar colunas únicas, com collapses/accordions para reduzir ruído visual.

---

## 6. Materiais: Vidro, Metal e Neon

### 6.1 Vidro

- Usado em modais, overlays e painéis especiais.
- Características:
  - `backdrop-filter: blur(...)`
  - Superficie semi-transparente
  - Borda sutil com cor metálica

### 6.2 Metal

- Representado por bordas finas, linhas divisórias e ícones.
- Combina com:
  - Bordas `1px`
  - Cores derivadas de `border` + accents sutis.

### 6.3 Neon (Accents)

- Principalmente `aqua` e `purple`.
- Uso recomendado:
  - Estados de foco
  - Highlights em gráficos e status
  - Ações de alta relevância (com moderação)

Evitar telas “cheias de neon”: o luxo vem da contenção.

---

## 7. Sombras e Elevação

- **Shadow Light (card):** `0 10px 30px rgba(62, 91, 255, 0.06)`
- **Shadow Dark (card):** `0 12px 32px rgba(0,0,0,0.35)`
- **Hover (accent aqua):** `0 6px 18px rgba(34, 211, 238, 0.12)`
- **Modais:** `0 24px 60px rgba(0,0,0,0.5)` + efeito de vidro

Regras:

- Não exagerar na quantidade de sombras.
- Usar sombras mais fortes apenas para elementos realmente importantes (ex.: modal, drawer principal).

---

## 8. Breakpoints

Breakpoints compatíveis com o grid do MUI:

- `xs: 0`
- `sm: 600`
- `md: 900`
- `lg: 1200`
- `xl: 1536`

Diretrizes:

- Ações primárias devem aparecer **acima da dobra** em `md`.
- Em `sm`, reduzir ruído visual com:
  - Accordions
  - Steppers
  - Resumos compactos

---

## 9. Princípios “Cyber Luxury”

- **Contraste inteligente:** fundos suaves, superfícies de vidro, cor viva apenas onde importa.
- **Materiais híbridos:** vidro (blur) + metal (borda fina) + neon (accents).
- **Legibilidade técnica:** tipografia geométrica, espaçamento generoso e hierarquia clara.
- **Consistência:** se um padrão visual foi usado em uma tela, deve ser repetido em telas da mesma família.

Esses princípios somados garantem que o VALTARIS tenha uma cara única – não é um painel genérico qualquer.
