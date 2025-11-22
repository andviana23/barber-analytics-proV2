# VALTARIS v 1.0 — Patterns

Padrões obrigatórios para formularios, acessibilidade e consistência “Cyber Luxury”. Light é padrão; Dark segue mesmo contrato de tokens.

## Formulários (RHF + Zod + MUI)
- Validação: schema-first com Zod; não validar em handlers soltos.
- UI: wrappers de input/select exibem erro em `var(--valtaris-danger)` e mantêm borda metálica.
- Submissão: usar `handleSubmit`; botão primário com estado de loading.

### Wrappers base
```tsx
import { Controller, useFormContext } from "react-hook-form";
import { TextField } from "@mui/material";

type FormInputProps = { name: string; label: string; helperText?: string } & Omit<
  React.ComponentProps<typeof TextField>,
  "name" | "label"
>;

export function FormInput({ name, label, helperText, ...rest }: FormInputProps) {
  const { control } = useFormContext();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <TextField
          {...rest}
          {...field}
          label={label}
          fullWidth
          error={!!fieldState.error}
          helperText={fieldState.error?.message || helperText}
          sx={{
            "& .MuiOutlinedInput-root": {
              borderRadius: "10px",
              backgroundColor: "var(--valtaris-surface-subtle)",
            },
            "& .MuiOutlinedInput-notchedOutline": {
              borderColor: fieldState.error ? "var(--valtaris-danger)" : "var(--valtaris-border)",
            },
            "&:hover .MuiOutlinedInput-notchedOutline": {
              borderColor: fieldState.error ? "var(--valtaris-danger)" : "var(--valtaris-primary)",
            },
          }}
        />
      )}
    />
  );
}
```

### Exemplo completo (RHF + Zod)
```tsx
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { FormProvider, useForm } from "react-hook-form";
import { VButton } from "../components/ui/v-button";
import { FormInput } from "../components/forms/form-input";

const schema = z.object({
  name: z.string().min(2, "Nome muito curto"),
  email: z.string().email("E-mail inválido"),
  value: z.number().positive("Valor deve ser positivo"),
});

export function SampleForm({ onSubmit }: { onSubmit: (data: any) => void }) {
  const methods = useForm({ resolver: zodResolver(schema), defaultValues: { name: "", email: "", value: 0 } });
  const { handleSubmit, formState } = methods;

  return (
    <FormProvider {...methods}>
      <form onSubmit={handleSubmit(onSubmit)} style={{ display: "grid", gap: 16 }}>
        <FormInput name="name" label="Nome" />
        <FormInput name="email" label="E-mail" />
        <FormInput name="value" label="Valor" type="number" />
        <VButton type="submit" disabled={formState.isSubmitting}>
          {formState.isSubmitting ? "Enviando..." : "Salvar"}
        </VButton>
      </form>
    </FormProvider>
  );
}
```

## Acessibilidade (WCAG)
- Contraste mínimo 4.5:1; `primary` com `#FFFFFF` cumpre; `accent.gold` deve vir com texto escuro.
- Focus ring sempre visível (`accent.aqua`); não remover outline.
- Hit area mínima 40px; padding 12–16px.
- Semântica: `aria-describedby` conecta helper/error; MUI `error` já seta `aria-invalid`.
- Teclado: Enter/Espaço em toggles/menus; ordem de tab lógica.

## Feedback e Estados
- Erro: mensagem curta, cor `var(--valtaris-danger)`, ícone opcional.
- Sucesso: `status.success`; manter fundo suave (`rgba(56,214,155,0.14)`).
- Loading: skeleton para tabelas/listas; spinner em botões apenas.
- Empty: texto `text.secondary`, CTA secundário com borda metálica.

## Layout Responsivo
- Form grids 2 colunas em `md+`, 1 coluna em `sm`.
- Section gaps 24–32px; agrupar campos relacionados em cards com título.
- Tabelas: oferecer “compact mode” via prop, não CSS manual.

## Erros a evitar
- CSS inline stringificado (prejudica DayPilot/hidratação).
- Cores fora dos tokens (`#fff` solto); sempre `--valtaris-*` ou `theme.palette`.
- Inputs sem schema Zod; perder mensagens consistentes.
