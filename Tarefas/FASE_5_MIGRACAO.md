# 🟦 FASE 5 — Migração Progressiva do MVP 1.0

**Objetivo:** Desativar gradualmente MVP 1.0, migrar para v2
**Duração:** 14-28 dias
**Dependências:** ✅ Fase 3 + Fase 4 completas
**Sprint:** Sprint 7-9

---

## 📊 Progresso Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 5: MIGRAÇÃO PROGRESSIVA MVP 1.0 → V2                  │
├─────────────────────────────────────────────────────────────┤
│  Progresso:  ██████████████░░░░░░  70% (2.8/4 concluídas)  │
│  Status:     🟡 EM PROGRESSO (ESTRATÉGIA SIMPLIFICADA)     │
│  Prioridade: 🔴 ALTA                                        │
│  Estimativa: 20 horas (14h concluídas, 6h restantes)        │
│  Sprint:     Sprint 7-9                                     │
│  Mudança:    ⚠️  SEM DUAL-READ - APENAS V2                 │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Checklist de Tarefas

### 🔴 T-INFRA-015 — Feature flags (Beta mode)
- **Responsável:** DevOps / Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 7
- **Status:** ✅ 100% CONCLUÍDO
- **Deliverable:** Sistema de feature flags por tenant

#### Critérios de Aceitação
- [x] Tabela `feature_flags` criada (migration 011)
  - [x] tenant_id, feature, enabled, created_at, updated_at
  - [x] Unique index (tenant_id, feature)
- [x] Repository PostgresFeatureFlagRepository conectado
- [x] Usecases: ListFeatureFlags, SetFeatureFlag, CheckFeatureFlag
- [x] Middleware FeatureFlagMiddleware implementado
- [x] Exemplo: `use_v2_financial = true/false` por tenant
- [x] Admin endpoint: `PATCH /admin/feature-flags` + GET
- [x] Public endpoint: `GET /api/v1/feature-flags`
- [x] Seed script: `backend/scripts/sql/seed_feature_flags.sql`
- [x] Documentação: `docs/FEATURE_FLAGS_API.md`
- [x] Migration guide: `backend/scripts/MIGRATION_GUIDE.md`
- [x] Testes unitários: CheckFeatureFlagUseCase (6/6 passing)
- [x] Migration 011 aplicada em banco Neon
- [x] Seeds executados (3 flags por tenant, todos disabled)
- [x] Middleware aplicado nas rotas financeiras e assinaturas
- [x] Backend compilando sem erros
- [x] Validação via @pgsql: flags criadas para tenant E2E
- [ ] Frontend consome feature flags (provider criado, integração pending)
- [ ] Validação em staging com flag habilitada/desabilitada

**Files Created/Modified:**
- ✅ `backend/migrations/011_create_feature_flags.{up,down}.sql`
- ✅ `backend/internal/infrastructure/repository/postgres_feature_flag_repository.go`
- ✅ `backend/internal/application/usecase/featureflag/*.go` (3 usecases)
- ✅ `backend/internal/infrastructure/http/handler/feature_flag_handler.go`
- ✅ `backend/internal/infrastructure/http/middleware/feature_flag_middleware.go`
- ✅ `backend/cmd/api/main.go` (feature flag integration)
- ✅ `backend/scripts/sql/seed_feature_flags.sql`
- ✅ `backend/scripts/sql/migrate_mvp_to_v2.sql`
- ✅ `backend/scripts/MIGRATION_GUIDE.md`
- ✅ `docs/FEATURE_FLAGS_API.md`
- ✅ `backend/tests/unit/usecase/featureflag/check_feature_flag_usecase_test.go`

**Backend API Ready:**
```bash
# Listar flags do tenant
curl -H "X-Tenant-ID: e2e00000-0000-0000-0000-000000000001" \
  http://localhost:8080/api/v1/feature-flags

# Habilitar flag (admin)
curl -X PATCH http://localhost:8080/api/v1/admin/feature-flags \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "...", "feature": "use_v2_financial", "enabled": true}'
```

### 🟡 T-FE-013 — Integração Feature Flags (SIMPLIFICADO - APENAS V2)
- **Responsável:** Frontend
- **Prioridade:** 🟡 Média
- **Estimativa:** 2h (simplificado)
- **Sprint:** Sprint 8
- **Status:** ✅ 80% CONCLUÍDO (hooks criados, implementação pending)
- **Deliverable:** Frontend verifica flags e protege rotas v2

#### ⚠️ MUDANÇA DE ESTRATÉGIA
**NÃO haverá dual-read (MVP + v2).**
**Decisão:** Usar apenas v2. Feature flags controlam **acesso**, não **fonte de dados**.

#### Critérios de Aceitação (SIMPLIFICADOS)
- [x] Hook `useFeatureFlags` criado e documentado
- [x] Hook `useFeature` para verificação simples
- [x] Provider `FeatureFlagsProvider` criado
- [ ] Frontend verifica feature flag `use_v2_financial`
- [ ] Se `false`: Exibir mensagem "Feature não disponível para seu tenant"
- [ ] Se `true`: Renderizar normalmente (sempre lendo de v2)
- [ ] Adicionar badge "v2.0" nas páginas protegidas
- [ ] ValidationDashboard simplificado (apenas verificar se API v2 responde)

**Fluxo Simplificado:**
```tsx
function ReceitasPage() {
  const { enabled, isLoading } = useFeature('use_v2_financial');

  if (isLoading) return <Loading />;

  if (!enabled) {
    return (
      <Alert severity="warning">
        O módulo Financeiro v2 não está disponível para seu tenant.
        Entre em contato com o suporte.
      </Alert>
    );
  }

  // Sempre lê de v2 (useReceitas já aponta para /api/v1/receitas)
  return <ReceitasV2View />;
}
```

**Files Created:**
- ✅ `frontend/app/lib/hooks/useFeatureFlags.ts`
- ✅ `frontend/app/lib/providers/FeatureFlagsProvider.tsx`
- ⏳ `frontend/app/components/FeatureGate.tsx` (componente protetor)

**Próximos Passos:**
1. ~~Criar client Supabase~~ ❌ NÃO NECESSÁRIO
2. ~~Adaptar hooks para dual-read~~ ❌ NÃO NECESSÁRIO
3. ✅ Adicionar `<FeatureFlagsProvider>` no layout privado
4. ✅ Criar componente `<FeatureGate>` para proteger páginas
5. ✅ Adicionar badges "v2.0" nas páginas
6. ✅ ValidationDashboard simplificado (apenas health check v2)
7. ✅ Testes e2e com toggle ON/OFF

---

### 🔴 T-QA-004 — Testes de regressão
- **Responsável:** QA
- **Prioridade:** 🔴 Alta
- **Estimativa:** 8h
- **Sprint:** Sprint 8
- **Status:** ⏳ Não iniciado
- **Deliverable:** Suite de testes de regressão

#### Critérios de Aceitação
- [ ] Teste: Totais de receita/despesa batem (MVP vs v2)
- [ ] Teste: Assinaturas ativas corretas
- [ ] Teste: Cálculos de comissão corretos
- [ ] Teste: Relatórios geram corretamente
- [ ] Teste: Fluxo de caixa idêntico
- [ ] Teste: E2E completo (login → dashboard → CRUD)
- [ ] Relatório de diferenças (se houver)

---

### 🟡 T-DOM-010 — Desativar MVP 1.0 (gradualmente)
- **Responsável:** DevOps / Product
- **Prioridade:** 🟡 Média
- **Estimativa:** 4h
- **Sprint:** Sprint 9
- **Status:** ⏳ Não iniciado
- **Deliverable:** Rollout gradual v2 para 100% dos tenants

#### Critérios de Aceitação
- [ ] **Semana 1:** 25% dos tenants usam v2
  - [ ] Monitorar: errors, latência, feedback
- [ ] **Semana 2:** 50% dos tenants usam v2
  - [ ] Validar métricas
- [ ] **Semana 3:** 75% dos tenants usam v2
- [ ] **Semana 4:** 100% dos tenants usam v2
- [ ] MVP 1.0 desativado (read-only por 30 dias)
- [ ] Comunicação aos usuários

---

## 📈 Métricas de Sucesso

### Fase 5 completa quando:
- [ ] ✅ Todos os 4 tasks concluídos (100%)
- [ ] ✅ MVP 1.0 e v2 rodando em paralelo
- [ ] ✅ Feature flags controlam o acesso ao Financeiro v2
- [ ] ✅ Beta phase completa e validada
- [ ] ✅ 100% dos tenants migrados para v2
- [ ] ✅ MVP 1.0 desativado (somente backup)

---

## 🎯 Deliverables da Fase 5

| # | Deliverable | Status |
|---|-------------|--------|
| 1 | Feature flags sistema implementado | ✅ 100% CONCLUÍDO |
| 2 | Integração frontend (simplificada) | 🟡 90% (provider + FeatureGate criados) |
| 3 | Testes de regressão passando | ⏳ Pendente |
| 4 | Rollout gradual concluído (100%) | ⏳ Pendente (playbook criado) |

---

## 📝 Resumo de Progresso (15/11/2025 - ATUALIZADO)

### ⚠️ MUDANÇA DE ESTRATÉGIA

**Decisão:** **NÃO** usar dual-read (MVP + v2).
**Novo fluxo:** Feature flags controlam apenas **acesso** às rotas v2, não fonte de dados.
**Impacto:** Simplificou implementação de 4h → 2h (frontend).

### ✅ Concluído Nesta Sessão

#### Backend (T-INFRA-015) - 100%
- ✅ Migration 011 aplicada via @pgsql no Neon
- ✅ Seeds executados: 3 flags (financial, subscriptions, inventory) para tenant E2E
- ✅ Middleware `FeatureFlagMiddleware` aplicado nas rotas:
  - Financeiro v2: `/api/v1/receitas`, `/despesas`, `/fluxo-caixa`, `/dashboard`
  - Assinaturas v2: `/api/v1/assinaturas`
- ✅ Backend compilando sem erros
- ✅ Teste manual: flags retornam corretamente via API

#### Frontend (T-FE-013) - 90%
- ✅ Hook `useFeatureFlags` e `useFeature` criados
- ✅ Provider `FeatureFlagsProvider` criado (context global)
- ✅ Componente `FeatureGate` criado (proteção de páginas)
- ✅ Badge `V2Badge` para indicar versão
- ✅ Hook `useMultipleFeatures` para verificar múltiplas flags
- ✅ Arquivo `useDualRead.example.ts` **OBSOLETO** (estratégia descartada)

#### DevOps (T-DOM-010) - 50%
- ✅ Playbook de rollout criado (`backend/scripts/ROLLOUT_PLAYBOOK.md`):
  - Cronograma 4 semanas (25% → 50% → 75% → 100%)
  - Scripts de habilitação em massa
  - Procedimentos de rollback (< 1min)
  - Queries de validação
  - Checklist de execução

### ⏳ Pendente

#### T-FE-013 (10% restante - ~1h)
- [ ] Integrar `FeatureFlagsProvider` no layout privado (`app/(private)/layout.tsx`)
- [ ] Envolver páginas principais com `<FeatureGate>`:
  - ReceitasPage, DespesasPage, DashboardPage
  - AssinaturasPage (com flag `use_v2_subscriptions`)
- [ ] Adicionar `<V2Badge />` nos headers das páginas
- [ ] Testar toggle manual: desabilitar flag → ver mensagem de indisponibilidade

#### T-QA-004
- [ ] Suite de testes de regressão
- [ ] Validação de totais (MVP vs v2)
- [ ] Testes de cálculos (comissões, fluxo de caixa, etc.)

#### T-DOM-010
- [ ] Rollout gradual (25% → 50% → 75% → 100%)
- [ ] Monitoramento de métricas (errors, latência, feedback)
- [ ] Desativação do MVP após 100% migrado

---

## 🚀 Próximos Passos

Após completar **100%** da Fase 5:

👉 **Iniciar FASE 6 — Hardening** (`Tarefas/FASE_6_HARDENING.md`)

**Resumo Fase 6:**
- Segurança (rate limiting avançado, auditoria, RBAC)
- Observabilidade (Prometheus, Grafana, Sentry)
- Performance (query optimization, caching Redis)
- Compliance (LGPD, backup, DR)

---

## 📝 Plano de Rollout Detalhado

### Semana 1 — 25% dos tenants
```
Tenants selecionados:
- Tenants com menor volume de dados
- Tenants beta testers (voluntários)
- Total: ~5-10 tenants

Monitoramento:
- Errors: < 0.1%
- Latência p95: < 500ms
- Feedback: Positivo

Rollback: Se error rate > 1% → voltar para MVP
```

### Semana 2 — 50% dos tenants
```
Adicionar:
- Tenants de médio porte
- Total acumulado: ~20-25 tenants

Validação:
- Comparar totais financeiros (MVP vs v2)
- Verificar crons executando corretamente
```

### Semana 3 — 75% dos tenants
```
Adicionar:
- Tenants maiores
- Total acumulado: ~35-40 tenants

Validação:
- Performance sob carga
- Backup/restore testado
```

### Semana 4 — 100% dos tenants
```
Migrar restantes:
- Total: 50+ tenants

Ações:
- MVP 1.0 → Read-only (30 dias)
- Comunicar usuários: "Migração completa"
- Monitorar por 7 dias
```

---

**Última Atualização:** 14/11/2025
**Status:** 🔴 Não Iniciado (0%)
**Próxima Revisão:** Após completar 50% das tarefas
