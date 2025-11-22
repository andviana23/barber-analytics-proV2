# 🧪 05 — Checklist de QA (Quality Assurance)

**Última Atualização:** 21/11/2025
**Objetivo:** Validar qualidade, segurança e consistência antes do deploy

---

## 🔒 Segurança e Multi-Tenant

### Tenant Isolation

**Testes Obrigatórios:**

- [ ] **Financeiro:** Criar conta a pagar no tenant A, tentar acessar do tenant B → retorna `404 Not Found`
- [ ] **Metas:** Criar meta mensal no tenant A, listar do tenant B → lista vazia
- [ ] **Precificação:** Salvar config no tenant A, buscar do tenant B → retorna `404 Not Found`
- [ ] **Estoque:** Registrar movimentação no tenant A, listar do tenant B → lista vazia
- [ ] **DRE:** Gerar DRE no tenant A, buscar do tenant B → retorna `404 Not Found`

**Critério de Sucesso:**

- ✅ Nenhum dado vaza entre tenants
- ✅ Queries sempre filtram por `tenant_id`
- ✅ Testes automáticos garantem isolamento

---

### RBAC (Role-Based Access Control)

**Matriz de Permissões:**

| Endpoint                              | Owner | Manager | Accountant | Employee | Barber        |
| ------------------------------------- | ----- | ------- | ---------- | -------- | ------------- |
| **Financeiro (Payables/Receivables)** |
| GET /financial/payables               | ✅    | ✅      | ✅         | ❌       | ❌            |
| POST /financial/payables              | ✅    | ✅      | ❌         | ❌       | ❌            |
| PUT /financial/payables               | ✅    | ✅      | ❌         | ❌       | ❌            |
| DELETE /financial/payables            | ✅    | ❌      | ❌         | ❌       | ❌            |
| GET /financial/dre                    | ✅    | ✅      | ✅         | ❌       | ❌            |
| **Metas**                             |
| GET /metas/mensais                    | ✅    | ✅      | ✅         | ❌       | ✅ (próprias) |
| POST /metas/mensais                   | ✅    | ✅      | ❌         | ❌       | ❌            |
| **Precificação**                      |
| GET /pricing/config                   | ✅    | ✅      | ❌         | ❌       | ❌            |
| PUT /pricing/config                   | ✅    | ❌      | ❌         | ❌       | ❌            |
| POST /pricing/simulate                | ✅    | ✅      | ❌         | ❌       | ❌            |
| **Estoque**                           |
| GET /stock/movimentacoes              | ✅    | ✅      | ✅         | ✅       | ❌            |
| POST /stock/entrada                   | ✅    | ✅      | ❌         | ❌       | ❌            |
| POST /stock/saida                     | ✅    | ✅      | ❌         | ❌       | ❌            |

**Testes Obrigatórios:**

- [ ] Accountant tenta criar conta a pagar → `403 Forbidden`
- [ ] Employee tenta ver DRE → `403 Forbidden`
- [ ] Barber tenta editar meta mensal → `403 Forbidden`
- [ ] Manager tenta deletar conta a pagar → `403 Forbidden`

---

### Autenticação

**Testes:**

- [ ] Request sem token JWT → `401 Unauthorized`
- [ ] Request com token expirado → `401 Unauthorized`
- [ ] Request com token inválido → `401 Unauthorized`
- [ ] Request com token de outro tenant → `404 Not Found`

---

## ✅ Validação de Dados

### Constraints de Banco

**UNIQUE Constraints:**

- [ ] **DRE:** Criar DRE duplicado (mesmo `tenant_id`, `mes_ano`) → `409 Conflict`
- [ ] **Meta Mensal:** Criar meta duplicada (mesmo `tenant_id`, `mes_ano`) → `409 Conflict`
- [ ] **Meta Ticket:** Criar meta ticket duplicada (mesmo `tenant_id`, `mes_ano`) → `409 Conflict`
- [ ] **Compensação:** Criar compensação duplicada (mesmo `asaas_invoice_id`) → `409 Conflict`

**Foreign Keys:**

- [ ] Criar conta a pagar com `meio_pagamento_id` inexistente → `400 Bad Request`
- [ ] Criar meta barbeiro com `barbeiro_id` inexistente → `400 Bad Request`
- [ ] Criar movimentação de estoque com `produto_id` inexistente → `400 Bad Request`

**NOT NULL:**

- [ ] Criar conta a pagar sem `descricao` → `400 Bad Request`
- [ ] Criar meta mensal sem `valor_meta` → `400 Bad Request`
- [ ] Criar compensação sem `valor_previsto` → `400 Bad Request`

---

### Validações de Negócio

**Valores:**

- [ ] Criar conta a pagar com `valor` negativo → `400 Bad Request`
- [ ] Criar meta mensal com `valor_meta` = 0 → `400 Bad Request`
- [ ] Simular preço com margem de lucro 150% → `400 Bad Request`
- [ ] Registrar saída de estoque com quantidade negativa → `400 Bad Request`

**Status:**

- [ ] Marcar conta já PAGA como PAGA novamente → `400 Bad Request`
- [ ] Compensação com status inválido ("CONFIRMADOO") → `400 Bad Request`
- [ ] Conta com status inválido ("PAGANDO") → `400 Bad Request`

**Datas:**

- [ ] Criar DRE com `mes_ano` inválido ("2025-13") → `400 Bad Request`
- [ ] Criar compensação com `data_compensacao` no passado (quando D+ ainda não passou) → `400 Bad Request`
- [ ] Criar conta a pagar com `data_vencimento` no formato errado → `400 Bad Request`

---

## 📊 Funcionalidade (Fluxos de Negócio)

### Financeiro

**Contas a Pagar:**

- [ ] Criar conta a pagar → status inicial PENDENTE
- [ ] Marcar como PAGO → status muda para PAGO, `data_pagamento` preenchida
- [ ] Listar contas filtradas por status PENDENTE → retorna apenas pendentes
- [ ] Listar contas filtradas por período → retorna apenas do período

**Contas a Receber:**

- [ ] Criar conta a receber → status inicial PENDENTE
- [ ] Marcar como PAGO → status muda para PAGO, `data_recebimento` preenchida
- [ ] Vincular com assinatura → `assinatura_id` preenchido corretamente

**Compensação Bancária:**

- [ ] Criar compensação PREVISTO → `data_compensacao` calculada com D+ do meio de pagamento
- [ ] Marcar como COMPENSADO → status muda, valor confirmado
- [ ] Listar compensadas por período → filtro funciona

**DRE Mensal:**

- [ ] Gerar DRE do mês anterior → agrega receitas e despesas corretamente
- [ ] DRE calculado automaticamente (cron job) → aparece no banco
- [ ] Buscar DRE por `mes_ano` → retorna correto

**Fluxo de Caixa Diário:**

- [ ] Gerar fluxo diário → projeta entradas/saídas compensadas
- [ ] Buscar fluxo compensado por intervalo de datas → retorna lista

---

### Metas

**Meta Mensal:**

- [ ] Criar meta mensal → `mes_ano`, `valor_meta` salvos
- [ ] Calcular progresso → `valor_realizado` / `valor_meta` correto
- [ ] Alerta de desvio >= 20% → notificação enviada

**Meta Barbeiro:**

- [ ] Criar meta por barbeiro → vincula `barbeiro_id` corretamente
- [ ] Calcular progresso individual → valores corretos

**Meta Ticket Médio:**

- [ ] Criar meta de ticket médio → `valor_meta` salvo
- [ ] Calcular ticket médio realizado → média correta

---

### Precificação

**Configuração:**

- [ ] Salvar config → `margem_lucro`, `percentual_comissao` salvos
- [ ] Buscar config → retorna valores corretos

**Simulação:**

- [ ] Simular preço com custos fixos/variáveis → fórmula aplicada corretamente
- [ ] Salvar simulação → histórico gravado
- [ ] Buscar simulações → lista histórico

---

### Estoque

**Movimentações:**

- [ ] Registrar entrada → quantidade aumenta
- [ ] Registrar saída → quantidade diminui
- [ ] Consumo automático por serviço → quantidade consumida conforme configuração
- [ ] Alerta de estoque mínimo → notificação quando `quantidade_atual <= quantidade_minima`

**Validações:**

- [ ] Não permitir estoque negativo → `400 Bad Request`
- [ ] Apenas produtos ativos podem ser consumidos automaticamente

---

## ⏰ Automação (Cron Jobs)

### Execução Manual (Testes)

**GenerateDREMonthly:**

- [ ] Executar manualmente → DRE gerado no banco
- [ ] Log em `cron_run_logs` → start/end/duration registrados
- [ ] Métrica Prometheus → `cron_job_duration_seconds` atualizada

**GenerateFluxoDiario:**

- [ ] Executar manualmente → fluxo diário gerado
- [ ] Logs corretos
- [ ] Métricas corretas

**MarcarCompensacoes:**

- [ ] Executar manualmente → compensações marcadas conforme D+
- [ ] Logs corretos
- [ ] Métricas corretas

**NotifyPayables:**

- [ ] Executar manualmente → notificações enviadas (D-5, D-1, D0)
- [ ] Logs corretos

**CheckEstoqueMinimo:**

- [ ] Executar manualmente → alertas enviados
- [ ] Logs corretos

**CalculateComissoes:**

- [ ] Executar manualmente → comissões calculadas
- [ ] Logs corretos

---

### Configuração

**Feature Flags:**

- [ ] Desabilitar cron job via flag → job não executa
- [ ] Habilitar cron job via flag → job executa

**ENV Vars:**

- [ ] Mudar schedule via ENV → job executa no horário correto
- [ ] Desabilitar via ENV (`ENABLED=false`) → job não executa

---

## 🎨 Frontend (React Query + Services)

### Services

**Testes:**

- [ ] `dreService.getDRE(mes_ano)` → retorna DRE correto
- [ ] `payablesService.list(filters)` → retorna lista paginada
- [ ] `pricingService.simulate(params)` → retorna preço calculado
- [ ] `stockService.registrarEntrada(data)` → entrada registrada

**Erros:**

- [ ] Request sem auth → erro tratado, mensagem clara
- [ ] Request inválido → erro tratado, mensagem clara
- [ ] Timeout → retry automático (3x)

---

### Hooks React Query

**Estados:**

- [ ] `useDRE(mes_ano)` → `loading` true durante fetch
- [ ] Após fetch → `data` preenchido, `loading` false
- [ ] Erro → `error` preenchido, mensagem exibida

**Cache:**

- [ ] Criar conta a pagar → cache de `useContasPagar` invalidado
- [ ] Editar meta mensal → cache de `useMetasMensais` invalidado
- [ ] Cache keys incluem tenant ID → diferentes tenants têm cache separado

**Mutations:**

- [ ] `useCreateContaPagar` → sucesso atualiza cache
- [ ] `useSetMetaMensal` → sucesso invalida cache
- [ ] `useRegistrarEntrada` → sucesso invalida cache de estoque

---

## 📈 Performance

### Endpoints

**Latência (p95):**

- [ ] GET /financial/payables → < 500ms (dev)
- [ ] GET /financial/dre → < 500ms (dev)
- [ ] GET /metas/mensais → < 500ms (dev)
- [ ] POST /pricing/simulate → < 300ms (dev)
- [ ] GET /stock/movimentacoes → < 500ms (dev)

**Paginação:**

- [ ] Lista com 1000 registros → paginação funciona, não retorna tudo de uma vez
- [ ] `page_size=10` → retorna max 10 registros
- [ ] `page=2` → retorna próximos 10 registros

---

### Queries SQL

**EXPLAIN ANALYZE:**

- [ ] Queries de listagem usam índices corretos
- [ ] Agregações (SUM, AVG) usam índices
- [ ] Queries filtradas por `tenant_id` usam índice composto

**Otimizações:**

- [ ] Queries N+1 eliminadas (eager loading)
- [ ] Consultas desnecessárias eliminadas (caching)

---

## 📚 Documentação

**Swagger/OpenAPI:**

- [ ] Novos endpoints documentados em `/docs` ou Swagger
- [ ] Exemplos de request/response incluídos
- [ ] Códigos de status documentados (200, 400, 401, 403, 404, 409, 500)

**Postman/Insomnia:**

- [ ] Collection atualizada com novos endpoints
- [ ] Variáveis de ambiente configuradas (tenant_id, auth token)
- [ ] Exemplos de payloads válidos/inválidos

**README:**

- [ ] Instruções de configuração de cron jobs atualizadas
- [ ] Feature flags documentadas
- [ ] ENV vars obrigatórias listadas

---

## ✅ Checklist Final de QA

**Antes de aprovar para produção:**

- [ ] Todos os testes de segurança passando (tenant isolation, RBAC, auth)
- [ ] Todos os testes de validação passando (constraints, negócio, formatos)
- [ ] Todos os testes de funcionalidade passando (fluxos completos)
- [ ] Todos os testes de automação passando (cron jobs, logs, métricas)
- [ ] Todos os testes de frontend passando (services, hooks, estados)
- [ ] Performance aceitável (p95 < 500ms)
- [ ] Documentação completa e atualizada
- [ ] Code coverage >= 80%
- [ ] Linter sem warnings
- [ ] Deploy em staging funcionando
- [ ] Smoke tests em staging passando

---

**Resultado Final:** Se todos os itens estiverem ✅, **APROVAR** para produção.

**Próximo:** Após QA aprovado, seguir para `README.md` para visão geral dos bloqueios concluídos
