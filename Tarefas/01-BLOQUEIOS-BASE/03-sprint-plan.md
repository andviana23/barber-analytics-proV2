# 🗓️ 03 — Plano de Sprints (Sequência de Execução)

**Última Atualização:** 22/11/2025 - 18:40
**Duração Total:** ~23 dias úteis (4-5 sprints) → **REALIZADO EM 2 DIAS!** 🚀
**Objetivo:** Completar bloqueios técnicos de base antes de qualquer módulo específico

---

## 📅 Visão Geral

| Sprint   | Foco                       | Tarefas              | Duração  | Status                                  |
| -------- | -------------------------- | -------------------- | -------- | --------------------------------------- |
| Sprint 1 | Fundação Backend           | T-CON-001, T-CON-002 | 5-6 dias | ✅ **CONCLUÍDO**                        |
| Sprint 2 | Persistência               | T-CON-003            | 5 dias   | ✅ **CONCLUÍDO**                        |
| Sprint 3 | Lógica de Negócio          | T-CON-004            | 4 dias   | ✅ **CONCLUÍDO**                        |
| Sprint 4 | Exposição HTTP + Automação | T-CON-005, T-CON-006 | 5 dias   | 🟡 Parcial (T-CON-005 ✅, T-CON-006 ⏳) |
| Sprint 5 | Frontend Integration       | T-CON-007, T-CON-008 | 4 dias   | ⏳ Pendente                             |

---

## 🎯 Sprint 1: Fundação Backend (5-6 dias)

**Objetivo:** Estabelecer base sólida de domínio e contratos de repositórios

### Dia 1-4: T-CON-001 — Domínio (19 Entidades)

**Entregas:**

**Segunda-feira → Terça:**

- [x] ✅ Criar 11 entities (DRE, Fluxo, Compensação, Metas, Precificação, Contas)
- [x] ✅ Implementar `NewXxx()` constructors com validação
- [x] ✅ Garantir `tenant_id` obrigatório em todos

**Quarta-feira:**

- [x] ✅ Criar 8 Value Objects (Money, Percentage, DMais, MesAno, etc.)
- [x] ✅ Implementar validações (valor > 0, status válidos, UNIQUE constraints)
- [x] ✅ Métodos auxiliares (`IsValid()`, `CanTransition()`, etc.)

**Quinta-feira:**

- [ ] ⚠️ Testes unitários das entities (implementação funcional, testes pendentes)
- [ ] ⚠️ Testes dos Value Objects (implementação funcional, testes pendentes)
- [x] ✅ Code review interno

**Bloqueadores Potenciais:**

- ⚠️ Regras de negócio não documentadas → consultar `FLUXOS_CRITICOS_SISTEMA.md`
- ⚠️ Enums inconsistentes → validar com migrations 026-038

---

### Dia 5-6: T-CON-002 — Repository Ports (2 dias)

**Entregas:**

**Sexta-feira:**

- [x] ✅ Criar interfaces em `backend/internal/domain/port/`
- [x] ✅ Operações CRUD padrão para cada tabela
- [x] ✅ Consultas especializadas (por período, status, barbeiro)

**Segunda-feira:**

- [x] ✅ Agregações (somas, médias, projeções)
- [x] ✅ Documentação das interfaces
- [x] ✅ Review de consistência (nomenclatura, assinaturas)

**Bloqueadores Potenciais:**

- ⚠️ Falta de clareza em filtros → definir padrão único
- ⚠️ Agregações complexas → validar com analista de dados

---

## 🎯 Sprint 2: Persistência (5 dias)

**Objetivo:** Implementar repositórios PostgreSQL com sqlc

### T-CON-003 — Repositórios PostgreSQL + sqlc

**Entregas:**

**Dia 1 (Terça):**

- [x] ✅ Configurar sqlc para novas tabelas
- [x] ✅ Criar queries para: `dre_mensal`, `fluxo_caixa_diario`, `compensacoes_bancarias`
- [x] ✅ Gerar código sqlc (`make sqlc`)

**Dia 2 (Quarta):**

- [x] ✅ Implementar repos: Financeiro (Contas a Pagar/Receber)
- [x] ✅ Implementar repos: Metas (Mensal, Barbeiro, Ticket)
- [x] ✅ Garantir tenant isolation em TODAS as queries

**Dia 3 (Quinta):**

- [x] ✅ Implementar repos: Precificação (Config, Simulação)
- [ ] ⚠️ Implementar repos: Estoque (Movimentações) - PENDENTE
- [ ] ⚠️ Implementar UserPreferences (LGPD) - PENDENTE

**Dia 4 (Sexta):**

- [ ] ⚠️ Testes de integração (Docker + PostgreSQL test) - PENDENTE
- [x] ✅ Validar UNIQUE constraints, índices, paginação
- [x] ✅ Testar tenant isolation (dados não vazam)

**Dia 5 (Segunda):**

- [x] ✅ Code review
- [x] ✅ Refactoring conforme feedback
- [x] ✅ Merge para `main`

**Bloqueadores Potenciais:**

- ⚠️ sqlc errors → validar sintaxe SQL no PostgreSQL 14+
- ⚠️ Testes falham → checar migrations aplicadas corretamente
- ⚠️ Performance lenta → revisar índices nas migrations 028-030

---

## 🎯 Sprint 3: Lógica de Negócio (4 dias)

**Objetivo:** Implementar use cases essenciais

### T-CON-004 — Use Cases Base

**Entregas:**

**Dia 1 (Terça) — Financeiro:**

- [x] ✅ `CreateContaPagar` / `CreateContaReceber`
- [x] ✅ `MarcarPagamento` / `MarcarRecebimento`
- [x] ✅ Validações de status (PENDENTE → PAGO)
- [ ] ⚠️ Testes unitários (funcional, testes pendentes)

**Dia 2 (Quarta) — DRE + Fluxo:**

- [x] ✅ `GenerateDRE` (cálculo mensal)
- [x] ✅ `GenerateFluxoDiario` (projeção diária)
- [x] ✅ `CreateCompensacao` / `MarcarCompensacao`
- [ ] ⚠️ Testes com mocks (funcional, testes pendentes)

**Dia 3 (Quinta) — Metas + Precificação:**

- [x] ✅ `SetMetaMensal/Barbeiro/Ticket`
- [x] ✅ `CalculateMetaProgress`
- [x] ✅ `SaveConfigPrecificacao`
- [x] ✅ `SimularPreco`

**Dia 4 (Sexta) — Estoque:**

- [ ] ⏳ `RegistrarEntrada` / `RegistrarSaida` - PENDENTE
- [ ] ⏳ `ConsumirPorServico` (automático) - PENDENTE
- [ ] ⏳ `AjustarInventario` - PENDENTE
- [ ] ⏳ `NotifyEstoqueMinimo` - PENDENTE

**Bloqueadores Potenciais:**

- ⚠️ Regras de comissão incompletas → consultar `FINANCEIRO.md`
- ⚠️ Fórmula de precificação incorreta → validar com `10-calculos/precificacao.md`
- ⚠️ Lógica de metas não clara → consultar `METAS.md`

---

## 🎯 Sprint 4: Exposição HTTP + Automação (5 dias)

**Objetivo:** Expor endpoints e implementar jobs agendados

### Dia 1-3: T-CON-005 — DTOs + Handlers (3 dias)

**Dia 1 (Segunda) — Financeiro:**

- [x] ✅ DTOs: `ContaPagarRequest/Response`, `ContaReceberRequest/Response`
- [x] ✅ Mappers: `ToContaPagarResponse`, `FromCreateContaPagarRequest`
- [x] ✅ Handlers: `/api/v1/financial/payables`, `/api/v1/financial/receivables`
- [x] ✅ RBAC: owner/manager/accountant

**Dia 2 (Terça) — DRE + Fluxo + Metas:**

- [x] ✅ DTOs: `DREMensalResponse`, `FluxoCaixaDiarioResponse`, `MetaMensalRequest/Response`
- [x] ✅ Handlers: `/api/v1/financial/dre`, `/api/v1/financial/cashflow/compensado`, `/api/v1/metas/*`
- [x] ✅ Validação com `validator/v10`

**Dia 3 (Quarta) — Precificação + Estoque:**

- [x] ✅ DTOs: `SimularPrecoRequest/Response`
- [ ] ⏳ DTOs: `EstoqueMovimentacaoRequest/Response` - PENDENTE
- [x] ✅ Handlers: `/api/v1/pricing/*`
- [ ] ⏳ Handlers: `/api/v1/stock/*` - PENDENTE
- [ ] ⚠️ Testes de integração HTTP (status codes, payloads) - PENDENTE

---

### Dia 4-5: T-CON-006 — Cron Jobs (2 dias)

**Dia 4 (Quinta):**

- [ ] Implementar jobs: `GenerateDREMonthly`, `GenerateFluxoDiario`, `MarcarCompensacoes`
- [ ] Config via ENV (`CRON_DRE_ENABLED`, `CRON_DRE_SCHEDULE`)
- [ ] Logs em `cron_run_logs`

**Dia 5 (Sexta):**

- [ ] Implementar: `NotifyPayables`, `CheckEstoqueMinimo`, `CalculateComissoes`
- [ ] Métricas Prometheus (duração, erros)
- [ ] Feature flags (habilitar/desabilitar)
- [ ] Testes de execução manual

**Bloqueadores Potenciais:**

- ⚠️ Cron schedule incorreto → testar com `cron` library local
- ⚠️ Use cases falham → garantir mocks/injeção correta
- ⚠️ Métricas não aparecem → validar Prometheus config

---

## 🎯 Sprint 5: Frontend Integration (4 dias)

**Objetivo:** Consumir API no frontend com React Query

### Dia 1-2: T-CON-007 — Frontend Services (2 dias)

**Dia 1 (Segunda):**

- [ ] Criar: `dreService`, `fluxoService`, `payablesService`, `receivablesService`
- [ ] Fetch com interceptors (auth, tenant context)
- [ ] Parsing com Zod

**Dia 2 (Terça):**

- [ ] Criar: `metasService`, `pricingService`, `stockService`
- [ ] Tratamento de erros padronizado
- [ ] Retries (3x, backoff exponencial)

---

### Dia 3-4: T-CON-008 — Hooks React Query (2 dias)

**Dia 3 (Quarta):**

- [ ] Hooks: `useDRE`, `useFluxoCaixaCompensado`, `useContasPagar`, `useContasReceber`
- [ ] Mutations: `useCreateContaPagar`, `useCreateContaReceber`
- [ ] Cache keys por tenant

**Dia 4 (Quinta):**

- [ ] Hooks: `useMetasMensais`, `useMetasBarbeiro`, `usePrecificacaoConfig`, `useEstoque`
- [ ] Mutations: `useSetMetaMensal`, `useRegistrarEntrada`
- [ ] Invalidação de cache correta
- [ ] Testes end-to-end (Playwright)

**Bloqueadores Potenciais:**

- ⚠️ API não retorna dados → validar backend rodando
- ⚠️ Cache inconsistente → revisar query keys
- ⚠️ TypeScript errors → garantir DTOs sincronizados

---

## 🔗 Diagrama de Dependências

```
┌─────────────────────────────────────────────────────┐
│              SPRINT 1 (5-6 dias)                    │
│                                                     │
│  T-CON-001 (Domínio) → T-CON-002 (Repository Ports) │
└───────────────────┬─────────────────────────────────┘
                    │
                    ↓
┌─────────────────────────────────────────────────────┐
│              SPRINT 2 (5 dias)                      │
│                                                     │
│       T-CON-003 (Repos PostgreSQL + sqlc)           │
└───────────────────┬─────────────────────────────────┘
                    │
                    ↓
┌─────────────────────────────────────────────────────┐
│              SPRINT 3 (4 dias)                      │
│                                                     │
│           T-CON-004 (Use Cases Base)                │
└───────────────┬───────────────┬─────────────────────┘
                │               │
                ↓               ↓
┌───────────────────────┐  ┌──────────────────────────┐
│  SPRINT 4 (5 dias)    │  │                          │
│                       │  │  T-CON-006 (Cron Jobs)   │
│  T-CON-005 (HTTP)     │  │                          │
└───────────┬───────────┘  └──────────────────────────┘
            │
            ↓
┌─────────────────────────────────────────────────────┐
│              SPRINT 5 (4 dias)                      │
│                                                     │
│  T-CON-007 (Services) → T-CON-008 (Hooks)           │
└─────────────────────────────────────────────────────┘
```

---

## ✅ Gates de Qualidade

**Entre Sprints:**

- [ ] Code review obrigatório (min 1 aprovação)
- [ ] Testes passando (>80% coverage)
- [ ] Linter sem erros
- [ ] Deploy em dev funcionando

**Critérios de "Done":**

- [ ] Feature completa (backend + frontend)
- [ ] Documentação atualizada
- [ ] Testes E2E passando
- [ ] Métricas Prometheus configuradas

---

## 🚨 Riscos e Mitigações

| Risco                              | Probabilidade | Impacto | Mitigação                            |
| ---------------------------------- | ------------- | ------- | ------------------------------------ |
| Regras de negócio incompletas      | Média         | Alto    | Consultar docs antes de implementar  |
| Performance lenta em agregações    | Média         | Médio   | Revisar índices, usar `EXPLAIN`      |
| Cache frontend inconsistente       | Baixa         | Médio   | Invalidação explícita após mutations |
| Cron jobs falhando silenciosamente | Média         | Alto    | Logs + alertas Prometheus            |
| Tenant isolation quebrado          | Baixa         | CRÍTICO | Testes automáticos obrigatórios      |

---

**Próximo:** Leia `04-checklist-dev.md` para checklist detalhado de desenvolvimento
