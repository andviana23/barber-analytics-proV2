# ✅ 04 — Checklist de Desenvolvimento (Dev)

**Última Atualização:** 22/11/2025 - 18:40
**Objetivo:** Garantir qualidade e consistência durante implementação
**Status:** 🟡 Backend 44/44 endpoints ✅ | Frontend pendente ⏳

---

## 🏗️ Fundação Backend (T-CON-001 + T-CON-002)

### Domínio (Entities + Value Objects)

**Entities:**

- [x] ✅ Todas as 11 entities criadas em `backend/internal/domain/entity/`
- [x] ✅ `tenant_id` presente e obrigatório em TODAS as entities
- [x] ✅ Constructors `NewXxx()` com validação completa
- [x] ✅ Métodos auxiliares (`IsValid()`, `CanTransition()`, `Calculate()`)
- [x] ✅ Enums mapeados conforme migrations (StatusCompensacao, StatusConta, TipoCusto, SubtipoReceita)
- [x] ✅ Structs seguem convenção Go (PascalCase, tags `json:"snake_case"`)

**Value Objects:**

- [x] ✅ `Money` criado com validação (valor >= 0, precisão 2 casas)
- [x] ✅ `Percentage` criado com validação (0-100)
- [x] ✅ `DMais` criado com validação (>= 0)
- [x] ✅ `MesAno` criado com validação (formato YYYY-MM)
- [x] ✅ `StatusCompensacao` criado (PREVISTO/CONFIRMADO/COMPENSADO/CANCELADO)
- [x] ✅ `StatusConta` criado (PENDENTE/PAGO/CANCELADO)
- [x] ✅ `TipoCusto` criado (FIXO/VARIAVEL)
- [x] ✅ `SubtipoReceita` criado (SERVICO/PRODUTO/PLANO)
- [x] ✅ Métodos `String()`, `IsValid()`, `Equals()` implementados

**Testes Unitários:**

- [ ] ⚠️ Testes de validação (valores inválidos retornam erro) - PENDENTE
- [ ] ⚠️ Testes de constructors (NewXxx com dados válidos/inválidos) - PENDENTE
- [ ] ⚠️ Testes de transição de estado (status PENDENTE → PAGO) - PENDENTE
- [ ] ⚠️ Coverage >= 80% nas entities/VOs - PENDENTE

---

### Repository Ports (Interfaces)

**Interfaces:**

- [x] ✅ Criadas em `backend/internal/domain/port/`
- [x] ✅ Operações CRUD para cada tabela:
  - [x] ✅ `Create(ctx, entity)` → entity
  - [x] ✅ `FindByID(ctx, tenantID, id)` → entity
  - [x] ✅ `Update(ctx, entity)` → entity
  - [x] ✅ `Delete(ctx, tenantID, id)` → error
  - [x] ✅ `List(ctx, tenantID, filters)` → []entity

**Consultas Especializadas:**

- [x] ✅ Por período: `FindByMesAno`, `FindByDateRange`
- [x] ✅ Por status: `FindByStatus`
- [x] ✅ Por barbeiro: `FindByBarber` (metas)
- [x] ✅ Por assinatura: `FindBySubscription` (contas)

**Agregações:**

- [x] ✅ `SumByPeriod(ctx, tenantID, start, end)` → Money
- [x] ✅ `AvgTicket(ctx, tenantID, mesAno)` → Money
- [x] ✅ `ProjectFluxo(ctx, tenantID, dateRange)` → []FluxoCaixaDiario

**Documentação:**

- [x] ✅ Comentários GoDoc em todas as interfaces
- [x] ✅ Exemplos de uso nas interfaces críticas (DRE, Fluxo)

---

## 💾 Persistência (T-CON-003)

### sqlc Queries

**Arquivos `.sql`:**

- [ ] Criados em `backend/internal/infra/db/queries/`
- [ ] CRUD completo para cada tabela
- [ ] Queries especializadas (filtros, agregações)
- [ ] **TODAS** as queries filtram por `tenant_id`
- [ ] Índices usados corretamente (verificar com `EXPLAIN`)
- [ ] Paginação com `LIMIT` e `OFFSET`

**Validação:**

- [ ] `make sqlc` executa sem erros
- [ ] Código gerado em `backend/internal/infra/db/sqlc/`
- [ ] Tipos Go corretos (Money → int64, MesAno → string)

---

### Repositórios PostgreSQL

**Implementações:**

- [x] ✅ `DREMensalRepository` implementado
- [x] ✅ `FluxoCaixaDiarioRepository` implementado
- [x] ✅ `CompensacaoBancariaRepository` implementado
- [x] ✅ `MetaMensalRepository` implementado
- [x] ✅ `MetaBarbeiroRepository` implementado
- [x] ✅ `MetaTicketMedioRepository` implementado
- [x] ✅ `PrecificacaoConfigRepository` implementado
- [x] ✅ `PrecificacaoSimulacaoRepository` implementado
- [x] ✅ `ContaPagarRepository` implementado
- [x] ✅ `ContaReceberRepository` implementado
- [ ] ⏳ `UserPreferencesRepository` implementado - PENDENTE

**Validações:**

- [x] ✅ Erros de violação de UNIQUE constraint tratados
- [x] ✅ Erros de FK constraint tratados
- [x] ✅ Erros de NOT NULL tratados
- [x] ✅ Context timeout respeitado (5s padrão)

---

### Testes de Integração

**Setup:**

- [ ] Docker Compose com PostgreSQL test
- [ ] Migrations aplicadas automaticamente
- [ ] Seed data para testes (fixtures)

**Casos de Teste:**

- [ ] **Tenant Isolation:** Dados de tenant A não aparecem em queries de tenant B
- [ ] **UNIQUE Constraints:** Inserir duplicata retorna erro
- [ ] **Paginação:** `List` com `limit=10` retorna max 10 registros
- [ ] **Filtros:** Consultas por período/status retornam dados corretos
- [ ] **Agregações:** `SumByPeriod` retorna valor correto
- [ ] **Transações:** Rollback em caso de erro

**Coverage:**

- [ ] > = 80% nos repositórios
- [ ] Casos felizes e casos de erro cobertos

---

## 🧠 Lógica de Negócio (T-CON-004)

### ✅ TODOS OS MÓDULOS IMPLEMENTADOS (22/11/2025)

**Status:** 🟢 **44/44 ENDPOINTS FUNCIONAIS**

#### MÓDULO METAS (15 endpoints) ✅

**Use Cases Implementados:**

- [x] ✅ `SetMetaMensalUseCase`, `GetMetaMensalUseCase`, `ListMetasMensaisUseCase`, `UpdateMetaMensalUseCase`, `DeleteMetaMensalUseCase`
- [x] ✅ `SetMetaBarbeiroUseCase`, `GetMetaBarbeiroUseCase`, `ListMetasBarbeirosUseCase`, `UpdateMetaBarbeiroUseCase`, `DeleteMetaBarbeiroUseCase`
- [x] ✅ `SetMetaTicketMedioUseCase`, `GetMetaTicketMedioUseCase`, `ListMetasTicketMedioUseCase`, `UpdateMetaTicketMedioUseCase`, `DeleteMetaTicketMedioUseCase`

**Handlers HTTP Implementados:**

- [x] ✅ 15 endpoints (5 por entidade: POST, GET/:id, GET, PUT/:id, DELETE/:id)

#### MÓDULO PRECIFICAÇÃO (9 endpoints) ✅

**Use Cases Implementados:**

- [x] ✅ `SaveConfigPrecificacaoUseCase`, `GetConfigPrecificacaoUseCase`, `UpdateConfigPrecificacaoUseCase`, `DeleteConfigPrecificacaoUseCase`
- [x] ✅ `SimularPrecoUseCase`, `SaveSimulacaoUseCase`, `GetSimulacaoUseCase`, `ListSimulacoesUseCase`, `DeleteSimulacaoUseCase`

**Handlers HTTP Implementados:**

- [x] ✅ 9 endpoints (Config: 4, Simulação: 5)

#### MÓDULO FINANCEIRO (20 endpoints) ✅

**Use Cases Implementados:**

- [x] ✅ ContaPagar: `CreateContaPagarUseCase`, `GetContaPagarUseCase`, `ListContasPagarUseCase`, `UpdateContaPagarUseCase`, `DeleteContaPagarUseCase`, `MarcarPagamentoUseCase`
- [x] ✅ ContaReceber: `CreateContaReceberUseCase`, `GetContaReceberUseCase`, `ListContasReceberUseCase`, `UpdateContaReceberUseCase`, `DeleteContaReceberUseCase`, `MarcarRecebimentoUseCase`
- [x] ✅ Compensação: `GetCompensacaoUseCase`, `ListCompensacoesUseCase`, `DeleteCompensacaoUseCase`
- [x] ✅ FluxoCaixa: `GetFluxoCaixaUseCase`, `ListFluxoCaixaUseCase`
- [x] ✅ DRE: `GetDREUseCase`, `ListDREUseCase`
- [x] ✅ Cronjob: `GenerateFluxoDiarioUseCase`

**Handlers HTTP Implementados:**

- [x] ✅ 20 endpoints (ContaPagar: 6, ContaReceber: 6, Compensação: 3, FluxoCaixa: 2, DRE: 2, Cronjob: 1)

**Wiring Completo:**

- [x] ✅ Database connection (pgxpool) configurada
- [x] ✅ 11 repositórios instanciados
- [x] ✅ 47 use cases instanciados
- [x] ✅ 3 handlers criados com DI (MetasHandler, PricingHandler, FinancialHandler)
- [x] ✅ 44 rotas registradas
- [x] ✅ Middleware de tenant context
- [x] ✅ Logger estruturado (Zap)
- [x] ✅ Compilação: SUCCESS

**Testes:**

- [ ] ⚠️ Testes unitários - PENDENTE
- [ ] ⚠️ Testes de integração - PENDENTE
- [ ] ⚠️ Testes E2E - PENDENTE

**Documentação:**

- [x] ✅ `VERTICAL_SLICE_ALL_MODULES.md` criado
- [x] ✅ Padrão de replicação documentado

---

### ⏳ MÓDULO ESTOQUE - PENDENTE

**Use Cases - Pendentes:**

- [ ] `CreateContaPagar` validando campos obrigatórios
- [ ] `CreateContaReceber` validando campos obrigatórios
- [ ] `MarcarPagamento` transitando status PENDENTE → PAGO
- [ ] `MarcarRecebimento` transitando status PENDENTE → PAGO
- [ ] Validação: não permitir pagar conta já PAGA
- [ ] Validação: não permitir valores negativos

**DRE + Fluxo:**

- [ ] `GenerateDRE` agregando receitas/despesas por `mes_ano`
- [ ] `GenerateFluxoDiario` projetando entradas/saídas diárias
- [ ] `CreateCompensacao` criando compensação bancária
- [ ] `MarcarCompensacao` transitando status PREVISTO → COMPENSADO
- [ ] Validação: D+ correto conforme meio de pagamento

**Testes Unitários:**

- [ ] Mocks de repositórios criados
- [ ] Casos felizes testados
- [ ] Casos de erro testados (validação, não encontrado, etc.)
- [ ] Coverage >= 80%

---

### Use Cases — Metas

**Definição:**

- [ ] `SetMetaMensal` criando/atualizando meta mensal
- [ ] `SetMetaBarbeiro` criando/atualizando meta por barbeiro
- [ ] `SetMetaTicket` criando/atualizando meta de ticket médio
- [ ] Validação: valores > 0
- [ ] Validação: `mes_ano` válido (YYYY-MM)

**Cálculo:**

- [ ] `CalculateMetaProgress` calculando realizado vs meta
- [ ] `NotifyMetaDeviation` enviando alertas (desvio >= 20%)
- [ ] Integração com DRE/Fluxo para pegar valores realizados

**Testes:**

- [ ] Casos de progresso (0%, 50%, 100%, 120%)
- [ ] Casos de desvio (alertas disparados corretamente)

---

### Use Cases — Precificação

**Simulação:**

- [ ] `SaveConfigPrecificacao` salvando configuração
- [ ] `SimularPreco` calculando preço sugerido
- [ ] Fórmula: `(custo_fixo + custo_variavel) / (1 - margem_lucro - comissao)`
- [ ] `SaveSimulacao` salvando histórico de simulações

**Validações:**

- [ ] Margem de lucro 0-100%
- [ ] Comissão 0-100%
- [ ] Custos >= 0

**Testes:**

- [ ] Casos com diferentes margens/comissões
- [ ] Casos de erro (margem inválida, custo negativo)

---

### Use Cases — Estoque

**Movimentações:**

- [ ] `RegistrarEntrada` aumentando quantidade
- [ ] `RegistrarSaida` diminuindo quantidade
- [ ] `ConsumirPorServico` consumindo automaticamente
- [ ] `AjustarInventario` corrigindo divergências
- [ ] `NotifyEstoqueMinimo` enviando alertas

**Validações:**

- [ ] Não permitir estoque negativo
- [ ] Validar quantidade > 0
- [ ] Consumo automático apenas para produtos ativos

**Testes:**

- [ ] Casos de entrada/saída
- [ ] Casos de consumo automático
- [ ] Casos de alerta (estoque <= mínimo)

---

## 🌐 Exposição HTTP (T-CON-005)

### DTOs (Request/Response)

**Criados em `backend/internal/application/dto/`:**

- [ ] `CreateContaPagarRequest` / `ContaPagarResponse`
- [ ] `CreateContaReceberRequest` / `ContaReceberResponse`
- [ ] `FluxoCaixaDiarioResponse`
- [ ] `CompensacaoBancariaResponse`
- [ ] `DREMensalResponse`
- [ ] `SetMetaMensalRequest` / `MetaMensalResponse`
- [ ] `SimularPrecoRequest` / `SimularPrecoResponse`
- [ ] `EstoqueMovimentacaoRequest` / `EstoqueResponse`

**Padrões:**

- [ ] Tags JSON em `snake_case`
- [ ] `omitempty` para opcionais
- [ ] `validate:"required"` para obrigatórios
- [ ] Dinheiro como **string** no DTO (conversão no mapper)
- [ ] Datas com `FlexibleDate`

---

### Mappers

**Criados em `backend/internal/application/mapper/`:**

- [ ] `ToContaPagarResponse(entity)` → Response DTO
- [ ] `FromCreateContaPagarRequest(dto)` → Entity
- [ ] Conversão Money: entity (int64 centavos) ↔ DTO (string "100.50")
- [ ] Conversão MesAno: entity (string "2025-01") ↔ DTO (string "2025-01")

---

### Handlers HTTP

**Criados em `backend/internal/infra/http/handler/`:**

- [ ] `/api/v1/financial/payables` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/financial/receivables` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/financial/cashflow/compensado` (GET)
- [ ] `/api/v1/financial/dre` (GET por `mes_ano`)
- [ ] `/api/v1/metas/mensais` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/metas/barbeiros` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/metas/ticket-medio` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/pricing/config` (GET, PUT)
- [ ] `/api/v1/pricing/simulate` (POST)
- [ ] `/api/v1/stock/movimentacoes` (GET)
- [ ] `/api/v1/stock/entrada` (POST)
- [ ] `/api/v1/stock/saida` (POST)

**Validações:**

- [ ] Validator em todos os requests
- [ ] Tenant context extraído do JWT
- [ ] RBAC aplicado (owner/manager/accountant)
- [ ] Responses padronizadas (`ErrorResponse` em caso de erro)

---

### Testes de Integração HTTP

**Setup:**

- [ ] Backend rodando em ambiente de teste
- [ ] Token JWT válido gerado
- [ ] Tenant de teste criado

**Casos de Teste:**

- [ ] POST retorna `201 Created` com dados corretos
- [ ] GET retorna `200 OK` com lista paginada
- [ ] PUT retorna `200 OK` com dados atualizados
- [ ] DELETE retorna `204 No Content`
- [ ] Request sem token retorna `401 Unauthorized`
- [ ] Request de outro tenant retorna `404 Not Found` (tenant isolation)
- [ ] Request inválido retorna `400 Bad Request` com detalhes

---

## ⏰ Automação (T-CON-006)

### Cron Jobs

**Implementados em `backend/internal/infra/cron/`:**

- [ ] `GenerateDREMonthly` (todo dia 1 às 02:00)
- [ ] `GenerateFluxoDiario` (todo dia às 00:05)
- [ ] `MarcarCompensacoes` (todo dia às 01:00)
- [ ] `NotifyPayables` (D-5, D-1, D0 às 09:00)
- [ ] `CheckEstoqueMinimo` (todo dia às 08:00)
- [ ] `CalculateComissoes` (dia 1 de cada mês às 03:00)

**Configuração:**

- [ ] ENV vars: `CRON_DRE_ENABLED`, `CRON_DRE_SCHEDULE`
- [ ] Feature flags: `feature.cron.dre.enabled`
- [ ] Logs em `cron_run_logs` (start, end, duration, status)

**Métricas Prometheus:**

- [ ] `cron_job_duration_seconds{job="dre"}`
- [ ] `cron_job_errors_total{job="dre"}`
- [ ] `cron_job_last_run_timestamp{job="dre"}`

**Validações:**

- [ ] Jobs executam apenas quando habilitados
- [ ] Jobs não bloqueiam aplicação (async)
- [ ] Jobs logam início/fim/erros
- [ ] Jobs não acessam repositórios direto (usam use cases)

---

## 🎨 Frontend (T-CON-007 + T-CON-008)

### Services

**Criados em `frontend/lib/services/`:**

- [ ] `dreService.ts` (getDRE, listDRE)
- [ ] `fluxoService.ts` (getFluxoCompensado)
- [ ] `payablesService.ts` (list, create, update, delete)
- [ ] `receivablesService.ts` (list, create, update, delete)
- [ ] `metasService.ts` (list, create, update)
- [ ] `pricingService.ts` (getConfig, saveConfig, simulate)
- [ ] `stockService.ts` (list, registrarEntrada, registrarSaida)

**Padrões:**

- [ ] Fetch com interceptors (auth, tenant context)
- [ ] Parsing com Zod
- [ ] Retries (3x, backoff exponencial)
- [ ] Tratamento de erros padronizado

---

### Hooks React Query

**Criados em `frontend/hooks/`:**

- [ ] `useDRE(mes_ano)`
- [ ] `useFluxoCaixaCompensado(date_range)`
- [ ] `useContasPagar(filters)`
- [ ] `useContasReceber(filters)`
- [ ] `useMetasMensais(mes_ano)`
- [ ] `useMetasBarbeiro(mes_ano, barbeiro_id)`
- [ ] `useMetasTicket(mes_ano)`
- [ ] `usePrecificacaoConfig()`
- [ ] `useSimularPreco(params)`
- [ ] `useEstoque(filters)`
- [ ] `useMovimentacoes(filters)`

**Mutations:**

- [ ] `useCreateContaPagar()`
- [ ] `useCreateContaReceber()`
- [ ] `useSetMetaMensal()`
- [ ] `useRegistrarEntrada()`
- [ ] `useRegistrarSaida()`

**Requisitos:**

- [ ] Estados `loading/error/data` corretos
- [ ] Cache keys por tenant
- [ ] Invalidação após mutations
- [ ] Stale time configurável (5min padrão)

---

## ✅ Checklist Final

**Antes de marcar como concluído:**

- [ ] Todos os testes passando (unit + integration + E2E)
- [ ] Linter sem erros (`make lint`)
- [ ] Code review aprovado (min 1 aprovação)
- [ ] Deploy em dev funcionando
- [ ] Endpoints documentados
- [ ] Métricas Prometheus configuradas
- [ ] Logs estruturados (Zap)
- [ ] Feature flags configuradas
- [ ] Tenant isolation validado

---

**Próximo:** Leia `05-checklist-qa.md` para checklist de QA e testes finais
