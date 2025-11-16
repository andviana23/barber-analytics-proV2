# 🚩 Feature Flags — Quick Start Guide

## ✅ O que foi feito (15/11/2025)

### Backend (100% pronto)
- ✅ Migration `011_create_feature_flags` já aplicada
- ✅ Repository, Usecases e Handler conectados
- ✅ Endpoints REST funcionando:
  - `GET /api/v1/feature-flags` (tenant)
  - `GET/PATCH /api/v1/admin/feature-flags` (admin)
- ✅ Middleware `FeatureFlagMiddleware` implementado
- ✅ Testes unitários passando (6/6)
- ✅ Documentação completa (`docs/FEATURE_FLAGS_API.md`)
- ✅ Seeds criados (`backend/scripts/sql/seed_feature_flags.sql`)
- ✅ Migration script MVP→v2 (`backend/scripts/sql/migrate_mvp_to_v2.sql`)

### Frontend (80% pronto)
- ✅ Hook `useFeatureFlags` criado (`frontend/app/lib/hooks/useFeatureFlags.ts`)
- ✅ Hook `useFeature` para verificação simples
- ✅ Padrão dual-read documentado (`frontend/app/lib/hooks/useDualRead.example.ts`)
- ⏳ **Pendente:** Integrar nas páginas (ReceitasPage, DespesasPage, Dashboard)

---

## 🚀 Como usar agora

### 1. Popular feature flags no banco

```bash
cd backend
psql $DATABASE_URL -f scripts/sql/seed_feature_flags.sql
```

**O que faz:**
- Cria flags `use_v2_financial`, `use_v2_subscriptions`, `use_v2_inventory` para todos os tenants
- Todas inicialmente `enabled = false`

---

### 2. Verificar flags via API

```bash
# Listar flags do tenant e2e
curl -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001" \
  http://localhost:8080/api/v1/feature-flags | jq
```

**Esperado:**
```json
{
  "code": 200,
  "data": [
    {
      "feature": "use_v2_financial",
      "enabled": false
    },
    {
      "feature": "use_v2_subscriptions",
      "enabled": false
    }
  ]
}
```

---

### 3. Habilitar flag para tenant beta

```bash
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "e2e00000-0000-0000-0000-000000000001",
    "feature": "use_v2_financial",
    "enabled": true
  }' | jq
```

**Esperado:**
```json
{
  "code": 200,
  "message": "Feature flag updated successfully",
  "data": {
    "feature": "use_v2_financial",
    "enabled": true,
    "updated_at": "2024-11-15T..."
  }
}
```

---

### 4. Usar no frontend (exemplo)

```tsx
// pages/receitas/page.tsx
import { useFeature } from '@/lib/hooks/useFeatureFlags';

export default function ReceitasPage() {
  const { enabled, isLoading } = useFeature('use_v2_financial');

  if (isLoading) return <Loading />;

  if (!enabled) {
    return (
      <Alert severity="warning">
        Feature Financeiro v2 não habilitada para seu tenant.
      </Alert>
    );
  }

  return <ReceitasV2View />;
}
```

---

## 🧪 Testes validados

### Backend
```bash
cd backend
go test ./tests/unit/usecase/featureflag/... -v
```

**Resultado:**
```
=== RUN   TestCheckFeatureFlagUseCase_Success
--- PASS: TestCheckFeatureFlagUseCase_Success (0.00s)
=== RUN   TestCheckFeatureFlagUseCase_FlagDisabled
--- PASS: TestCheckFeatureFlagUseCase_FlagDisabled (0.00s)
=== RUN   TestCheckFeatureFlagUseCase_FlagNotFound
--- PASS: TestCheckFeatureFlagUseCase_FlagNotFound (0.00s)
=== RUN   TestCheckFeatureFlagUseCase_InvalidTenantID
--- PASS: TestCheckFeatureFlagUseCase_InvalidTenantID (0.00s)
=== RUN   TestCheckFeatureFlagUseCase_EmptyFeature
--- PASS: TestCheckFeatureFlagUseCase_EmptyFeature (0.00s)
=== RUN   TestCheckFeatureFlagUseCase_RepositoryError
--- PASS: TestCheckFeatureFlagUseCase_RepositoryError (0.00s)
PASS
```

✅ **6/6 testes passando**

---

## 📋 Próximos passos (implementação final)

### Frontend (20% restante)
1. **Integrar `useFeatureFlags` nas páginas principais**
   - ReceitasPage
   - DespesasPage
   - DashboardPage
   - AssinaturasPage

2. **Criar client Supabase** (se necessário para buscar dados MVP durante dual-read)
   ```typescript
   // frontend/app/lib/supabase/client.ts
   import { createClient } from '@supabase/supabase-js';

   export const supabase = createClient(
     process.env.NEXT_PUBLIC_SUPABASE_URL!,
     process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
   );
   ```

3. **Adaptar hooks existentes** para seguir padrão dual-read
   - Ver exemplo em `frontend/app/lib/hooks/useDualRead.example.ts`
   - Implementar strategy: buscar MVP vs v2 baseado em flag
   - Adicionar comparação de totais (logs + Sentry)

4. **Criar ValidationDashboard**
   ```typescript
   // frontend/app/(private)/admin/validation/page.tsx
   import { ValidationDashboard } from '@/lib/hooks/useDualRead.example';

   export default function ValidationPage() {
     return <ValidationDashboard />;
   }
   ```

5. **Testes e2e com feature flags**
   ```typescript
   // e2e/feature-flags.spec.ts
   test('deve mostrar MVP quando flag desabilitada', async ({ page }) => {
     await page.goto('/receitas');
     await expect(page.getByText('⚠️ MVP')).toBeVisible();
   });

   test('deve mostrar v2 quando flag habilitada', async ({ page }) => {
     // Habilitar flag via API
     await apiToggleFlag('use_v2_financial', true);
     await page.goto('/receitas');
     await expect(page.getByText('✅ v2.0')).toBeVisible();
   });
   ```

---

## 🔄 Workflow de Migração

### Fase 1: Beta (Semana 1)
1. Selecionar 3-5 tenants beta
2. Habilitar `use_v2_financial` via PATCH endpoint
3. Monitorar:
   - Sentry (erros 500)
   - Logs (divergências de cálculo)
   - Feedback dos usuários

### Fase 2: Expansão (Semanas 2-4)
- **Semana 2:** 50% dos tenants
- **Semana 3:** 75% dos tenants
- **Semana 4:** 100% dos tenants

### Rollback (se necessário)
```bash
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "<uuid>", "feature": "use_v2_financial", "enabled": false}'
```

---

## 📚 Documentação Completa

- **API Reference:** `docs/FEATURE_FLAGS_API.md`
- **Migration Guide:** `backend/scripts/MIGRATION_GUIDE.md`
- **Dual-Read Pattern:** `frontend/app/lib/hooks/useDualRead.example.ts`
- **Roadmap:** `Tarefas/FASE_5_MIGRACAO.md` (50% concluído)

---

## 🆘 Troubleshooting

### Erro: "Feature flag not found"
**Solução:** Executar seeds
```bash
psql $DATABASE_URL -f backend/scripts/sql/seed_feature_flags.sql
```

### Flag não responde a PATCH
**Causa:** Middleware de autenticação admin não implementado (dev mode aceitando qualquer request)
**Status:** OK para dev, adicionar RBAC em produção

### Frontend não carrega flags
**Causa:** Endpoint `/feature-flags` retornando 403 (MISSING_TENANT)
**Solução:** Verificar header `X-Tenant-ID` no interceptor do axios

---

**✅ Sistema pronto para uso em 80%**
**⏳ Implementação final no frontend (estimativa: 2-3h)**
