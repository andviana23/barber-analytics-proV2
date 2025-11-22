# Fluxo de Permissões (RBAC) — VALTARIS (MVP 1.0)
```
[Início]
   ↓
[Usuário faz login]
   ↓
[Sistema identifica papel: Proprietario / Gerente / Barbeiro / Recepção / Contador]
   ↓
[Carregar permissões do papel]
   ↓
Acesso permitido?
   → Não → [Exibir erro 403]
   → Sim → [Liberar rota/tela]
   ↓
[Fim]
```
