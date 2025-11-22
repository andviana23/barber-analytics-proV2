# VALTARIS v 1.0 — Components

Componentes base com estética “Cyber Luxury”: borda fina metálica, vidro em overlays e acentos neon contidos. Use sempre tokens do tema (MUI ou `--valtaris-*`). Light é o default.

## Botões
Quando usar: ações primárias/secundárias, CTAs de fluxo. Variantes:
- **Primary**: ações de confirmação; fundo `primary`, hover `primaryDark`.
- **Secondary**: ações neutras; fundo `surface-subtle`, borda `border`.
- **Ghost**: ações de navegação/filtragem; fundo transparente, borda metálica leve.

```tsx
import { Button, alpha, useTheme } from "@mui/material";

export function VButton({ kind = "primary", ...props }: { kind?: "primary" | "secondary" | "ghost" } & React.ComponentProps<typeof Button>) {
  const theme = useTheme();
  const map = {
    primary: {
      bg: "var(--valtaris-primary)",
      hover: "var(--valtaris-primary-dark)",
      text: "#FFFFFF",
      border: "transparent",
    },
    secondary: {
      bg: "var(--valtaris-surface-subtle)",
      hover: "var(--valtaris-surface)",
      text: "var(--valtaris-text)",
      border: "var(--valtaris-border)",
    },
    ghost: {
      bg: "transparent",
      hover: alpha(theme.palette.primary.main, 0.08),
      text: "var(--valtaris-text)",
      border: alpha("#8A7CFF", 0.25),
    },
  }[kind];

  return (
    <Button
      {...props}
      variant="contained"
      sx={{
        px: 2.5,
        py: 1.25,
        minHeight: 40,
        borderRadius: "10px",
        textTransform: "none",
        fontWeight: 600,
        letterSpacing: 0.4,
        backgroundColor: map.bg,
        color: map.text,
        border: `1px solid ${map.border}`,
        boxShadow: "0 10px 30px rgba(62, 91, 255, 0.12)",
        "&:hover": {
          backgroundColor: map.hover,
          boxShadow: "0 14px 36px rgba(34, 211, 238, 0.18)",
        },
        "&:focus-visible": {
          outline: "2px solid var(--valtaris-accent-aqua)",
          outlineOffset: 2,
        },
        "&.Mui-disabled": {
          opacity: 0.48,
          boxShadow: "none",
        },
      }}
    />
  );
}
```

## Inputs e Selects
- Sempre envolver com wrappers RHF + Zod (ver `04-PATTERNS.md`).
- Erros em `var(--valtaris-danger)`, borda metálica preservada.
- Prefira `TextField`/`Select` MUI usando `sx`.

```tsx
import { TextField } from "@mui/material";

export function VTextField(props: React.ComponentProps<typeof TextField>) {
  return (
    <TextField
      {...props}
      fullWidth
      sx={{
        "& .MuiOutlinedInput-root": {
          borderRadius: "10px",
          backgroundColor: "var(--valtaris-surface-subtle)",
        },
        "& .MuiOutlinedInput-notchedOutline": {
          borderColor: "var(--valtaris-border)",
        },
        "&:hover .MuiOutlinedInput-notchedOutline": {
          borderColor: "var(--valtaris-primary)",
        },
        "& .Mui-error .MuiOutlinedInput-notchedOutline": {
          borderColor: "var(--valtaris-danger)",
        },
      }}
    />
  );
}
```

## Cards e Superfícies
- Cards padrão: fundo `surface`, borda 1px, raio 14px, sombra leve.
- Cards de destaque: aplicar `backdrop-filter: blur(10px)` para vidro.

```tsx
import { Card, CardContent, Typography, Box } from "@mui/material";

export function VCard({ title, subtitle, children }: { title: string; subtitle?: string; children: React.ReactNode }) {
  return (
    <Card
      elevation={0}
      sx={{
        background: "var(--valtaris-surface)",
        border: "1px solid var(--valtaris-border)",
        borderRadius: "14px",
        boxShadow: "0 12px 30px rgba(0,0,0,0.12)",
        backdropFilter: "blur(10px)",
      }}
    >
      <CardContent sx={{ p: 3, display: "grid", gap: 1.5 }}>
        <Box>
          <Typography variant="h6" sx={{ fontWeight: 600, letterSpacing: -0.2 }}>
            {title}
          </Typography>
          {subtitle && (
            <Typography variant="body2" color="text.secondary">
              {subtitle}
            </Typography>
          )}
        </Box>
        {children}
      </CardContent>
    </Card>
  );
}
```

## Modais (vidro obrigatório)
- Backdrop com blur e opacidade controlada.
- Paper com borda metálica e sombra forte.

```tsx
import { Dialog, DialogTitle, DialogContent, IconButton } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";

export function VModal({ open, onClose, title, children }) {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      slotProps={{
        backdrop: {
          sx: {
            backdropFilter: "blur(16px)",
            backgroundColor: "rgba(11,13,18,0.65)",
          },
        },
        paper: {
          sx: {
            borderRadius: "16px",
            border: "1px solid rgba(255,255,255,0.08)",
            background: "var(--valtaris-surface)",
            boxShadow: "0 24px 60px rgba(0,0,0,0.5)",
          },
        },
      }}
    >
      <DialogTitle sx={{ display: "flex", alignItems: "center", gap: 1 }}>
        {title}
        <IconButton
          aria-label="Fechar"
          onClick={onClose}
          sx={{
            ml: "auto",
            color: "var(--valtaris-text)",
            "&:hover": { color: "var(--valtaris-primary)" },
          }}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      </DialogTitle>
      <DialogContent sx={{ pb: 3 }}>{children}</DialogContent>
    </Dialog>
  );
}
```

## Badges e Status
- Usar cores de status oficiais: success `#38D69B`, warning `#F4B23E`, danger `#EF4444`.
- Texto em `text.primary` no light; em dark, preferir texto claro ou preto no badge gold.

```tsx
import { Chip } from "@mui/material";

export function StatusBadge({ status }: { status: "success" | "warning" | "danger" }) {
  const map = {
    success: { bg: "rgba(56,214,155,0.14)", color: "#38D69B" },
    warning: { bg: "rgba(244,178,62,0.16)", color: "#F4B23E" },
    danger: { bg: "rgba(239,68,68,0.16)", color: "#EF4444" },
  }[status];
  return <Chip label={status} size="small" sx={{ bgcolor: map.bg, color: map.color, borderRadius: "999px" }} />;
}
```

## DataTable
- Base: `Table`/`TableContainer` MUI; para grandes volumes, migrar para DataGrid mantendo tokens.
- Regras: linhas hover com `surface-subtle`, header com `text.secondary` e borda `border`.

```tsx
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from "@mui/material";

export function VTable({ rows }: { rows: { id: string; name: string; value: string }[] }) {
  return (
    <TableContainer component={Paper} sx={{ background: "var(--valtaris-surface)", border: "1px solid var(--valtaris-border)" }}>
      <Table size="small">
        <TableHead>
          <TableRow sx={{ "& th": { color: "var(--valtaris-text-muted)", borderBottom: "1px solid var(--valtaris-border)" } }}>
            <TableCell>Nome</TableCell>
            <TableCell align="right">Valor</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row) => (
            <TableRow
              key={row.id}
              hover
              sx={{
                "& td": { borderBottom: "1px solid var(--valtaris-border)" },
                "&:hover": { backgroundColor: "var(--valtaris-surface-subtle)" },
              }}
            >
              <TableCell>{row.name}</TableCell>
              <TableCell align="right">{row.value}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
```

## Interações gerais
- Focus ring sempre visível (`accent.aqua`).
- Hover curto (<=200ms) e suave; evitar delays.
- Loadings: `CircularProgress` com `accent.purple` ou `primary`.
- Empty states: ícones traço 2px, texto `text.secondary`, botão secundário com borda metálica.
