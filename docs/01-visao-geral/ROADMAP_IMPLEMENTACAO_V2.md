> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📋 Roadmap de Implementação 2.0 - Barber Analytics Pro

**Versão:** 2.0  
**Data:** 14/11/2025  
**Timeline Estimada:** 8-12 semanas  
**Status:** Planejamento

---

## 📊 Visão Geral do Roadmap

```
FASE 0: Fundamentos (1-3 dias)
   ↓
FASE 1: Infra & DevOps (3-7 dias)
   ↓
FASE 2: Backend Go Core (1-2 semanas)
   ↓
FASE 3: Módulos Críticos (2-4 semanas)
   ↓
FASE 4: Frontend 2.0 (2-4 semanas) [Paralelo a FASE 3]
   ↓
FASE 5: Migração Progressiva (2-4 semanas)
   ↓
FASE 6: Hardening & Segurança (1-2 semanas)
   ↓
✅ MVP 2.0 em Produção
```

---

## 🟦 FASE 0 - Fundamentos e Organização (1-3 dias)

**Objetivo:** Preparar o terreno para v2 sem quebrar o MVP 1.0

### Checklist

- [ ] **Criar repositório GitHub para backend v2**
  - [ ] Nome: `barber-analytics-backend-v2`
  - [ ] Template: Go + Clean Architecture
  - [ ] Branches: `main`, `develop`, `staging`
  - [ ] Proteção: Require PR reviews em `main`

- [ ] **Documentação Base**
  - [ ] README.md com stack e objetivos
  - [ ] CONTRIBUTING.md com padrões de código
  - [ ] Arquivo .gitignore específico para Go
  - [ ] LICENSE (MIT ou Apache 2.0)

- [ ] **Escolher provedor de PostgreSQL**
  - [ ] [ ] Opção A: Neon (recomendado)
    - [ ] Criar conta Neon
    - [ ] Projeto free tier
    - [ ] DATABASE_URL gerada
  - [ ] [ ] Opção B: Supabase (DB-only mode)
    - [ ] Criar projeto
    - [ ] Desabilitar Auth supabase
    - [ ] DATABASE_URL gerada
  
- [ ] **Decidir modelo Multi-Tenant**
  - [ ] Aprovado: **Column-Based** (tenant_id por linha)
  - [ ] Documentar em MODELO_MULTI_TENANT.md
  - [ ] Criar script de seed com tabelas base

- [ ] **Inicial Tooling**
  - [ ] Go 1.22+ instalado
  - [ ] `go mod init github.com/seu-usuario/barber-analytics-backend-v2`
  - [ ] `tools.go` com dependências de build
  - [ ] Makefile com targets: `build`, `run`, `test`, `lint`

### Entregas

- ✅ Repositório configurado e versionado
- ✅ Ambiente local testado
- ✅ Decisões técnicas documentadas

---

## 🟦 FASE 1 - Infra & DevOps Base (3-7 dias)

**Objetivo:** Ambiente pronto para rodar backend Go profissionalmente

### Checklist

- [ ] **Serviço backend (systemd)**
  - [ ] Build binário Go (make build) com Go 1.24.0
  - [ ] Serviço systemd configurado (`/etc/systemd/system/barber-api.service`)
  - [ ] Logs via journalctl + logrotate configurado
  - [ ] Variáveis de ambiente documentadas (.env + systemd EnvironmentFile)

- [ ] **NGINX & Reverse Proxy**
  - [ ] Instalação no VPS
  - [ ] Arquivo de configuração
    - [ ] `api.seudominio.com` → `:8080`
    - [ ] `app.seudominio.com` → frontend (Vercel/outro)
    - [ ] Compressão gzip ativada
    - [ ] Rate limiting configurado
  - [ ] SSL/TLS
    - [ ] Certbot instalado
    - [ ] Certificados Let's Encrypt
    - [ ] Auto-renewal com systemd timer
    - [ ] HSTS header configurado

- [ ] **CI/CD com GitHub Actions**
  - [ ] Workflow `.github/workflows/ci.yml`
    - [ ] Trigger: push em `develop` e `main`
    - [ ] Steps:
      - [ ] Checkout código
      - [ ] Setup Go
      - [ ] Cache Go modules
      - [ ] `go mod tidy` check
      - [ ] `go fmt` check
      - [ ] `go vet` check
      - [ ] Tests com coverage
      - [ ] Build binário
  - [ ] Workflow `.github/workflows/deploy.yml`
    - [ ] Trigger: push em `main`
    - [ ] Deploy via SSH/rsync para VPS + restart systemd
    - [ ] Health check pós-deploy

- [ ] **Logs & Monitoring Base**
  - [ ] Backend escreve logs em JSON (estruturado)
  - [ ] Sistema de coleta:
    - [ ] journalctl para systemd
  - [ ] Scripts úteis:
    - [ ] `scripts/tail-logs.sh` para desenvolvimento
    - [ ] `scripts/health-check.sh` para monitoramento

- [ ] **Health Check**
  - [ ] Endpoint `/health` que retorna `{"status": "ok"}`
  - [ ] Endpoint `/health/db` que testa conexão DB
  - [ ] NGINX monitora e redireciona erros 5xx

### Entregas

- ✅ Backend Go rodando via systemd
- ✅ NGINX configurado e operacional
- ✅ CI/CD pipeline testado
- ✅ Deploy automatizado testado em staging

---

## 🟦 FASE 2 - Backend Go Core (1-2 semanas)

**Objetivo:** Levantar a espinha dorsal do backend (auth, tenants, financeiro base)

### Checklist

#### 2.1 - Configuração & Boot

- [ ] **Config Management**
  - [ ] Arquivo `internal/config/config.go`
  - [ ] Leitura de variáveis de ambiente:
    - [ ] `DATABASE_URL`
    - [ ] `JWT_SECRET` ou `JWT_PRIVATE_KEY`
    - [ ] `APP_ENV` (dev, staging, prod)
    - [ ] `PORT` (default 8080)
    - [ ] `LOG_LEVEL` (debug, info, warn, error)
  - [ ] Validação de config obrigatória
  - [ ] Arquivo `.env.example` documentado

- [ ] **Main.go & Server Bootstrap**
  - [ ] Criação de `cmd/api/main.go`
  - [ ] Inicialização:
    - [ ] Config loading
    - [ ] Database connection
    - [ ] Logger initialization
    - [ ] Router setup
    - [ ] Server start
  - [ ] Graceful shutdown com timeout

- [ ] **Database Connection Pool**
  - [ ] `sql.Open()` com driver pq
  - [ ] Connection pool tuning
    - [ ] `MaxOpenConns: 25`
    - [ ] `MaxIdleConns: 5`
    - [ ] `ConnMaxLifetime: 5m`
  - [ ] Health check na startup
  - [ ] Retry logic para falhas transitórias

#### 2.2 - Estrutura de Camadas

- [ ] **Domain Layer** (`internal/domain/`)
  - [ ] Entidades base:
    - [ ] `User` (id, email, password_hash, tenant_id, role)
    - [ ] `Tenant` (id, nome, cnpj, ativo)
  - [ ] Value Objects:
    - [ ] `Email` (validação imutável)
    - [ ] `Role` (enum: owner, manager, accountant, employee)
  - [ ] Repository interfaces:
    - [ ] `UserRepository`
    - [ ] `TenantRepository`

- [ ] **Application Layer** (`internal/application/`)
  - [ ] DTOs de entrada/saída
  - [ ] Mappers (Domain ↔ DTO)
  - [ ] Use Cases base:
    - [ ] `LoginUseCase`
    - [ ] `RefreshTokenUseCase`
    - [ ] `CreateUserUseCase`

- [ ] **Infrastructure Layer** (`internal/infrastructure/`)
  - [ ] Repository implementations (PostgreSQL)
  - [ ] HTTP handlers e rotas
  - [ ] Middlewares

#### 2.3 - Módulo de Autenticação & Multi-Tenant

- [ ] **Auth Domain**
  - [ ] Entidade `User`
  - [ ] Value Object `Role`
  - [ ] Service `PasswordHasher` (bcrypt)
  - [ ] Service `TokenGenerator` (JWT)

- [ ] **Auth Repositories**
  - [ ] `PostgresUserRepository.SaveUser()`
  - [ ] `PostgresUserRepository.FindByEmail()`
  - [ ] `PostgresUserRepository.FindByID()`

- [ ] **Auth Use Cases**
  - [ ] Login Use Case
    - [ ] Validar email/senha
    - [ ] Gerar access + refresh tokens
    - [ ] Registrar último login
  - [ ] Refresh Token Use Case
    - [ ] Validar refresh token
    - [ ] Gerar novo access token
  - [ ] Create User Use Case
    - [ ] Validar dados
    - [ ] Hash password
    - [ ] Salvar user com tenant

- [ ] **Auth HTTP Layer**
  - [ ] POST `/auth/login` - Login
  - [ ] POST `/auth/refresh` - Refresh token
  - [ ] POST `/auth/register` - Criar nova conta (opcional)
  - [ ] POST `/auth/logout` - Logout (opcional)

- [ ] **Middleware de Auth**
  - [ ] JWT parsing
  - [ ] Token validation
  - [ ] Error handling (401, 403)
  - [ ] Context population

- [ ] **Middleware de Tenant**
  - [ ] Extrair tenant_id do token
  - [ ] Validar existência do tenant
  - [ ] Injetar no contexto da request
  - [ ] Garantir query sempre filtra por tenant

- [ ] **Migrations**
  - [ ] Tabela `users`
  - [ ] Tabela `tenants`
  - [ ] Índices e constraints

#### 2.4 - Base de Domínio Financeiro

- [ ] **Domain Financeiro**
  - [ ] Entidades:
    - [ ] `Receita` (id, tenant_id, descricao, valor, data, categoria, criado_em)
    - [ ] `Despesa` (idem)
    - [ ] `Categoria` (id, nome, tipo: RECEITA/DESPESA)
    - [ ] `MetodoPagamento` (dinheiro, débito, crédito, pix, etc)
  - [ ] Value Objects:
    - [ ] `Dinheiro` (valor, moeda)
    - [ ] `Periodo` (from, to)
  - [ ] Services:
    - [ ] `CalculoComissao` (cálculo de repasse barbeiro)
    - [ ] `CalculoFluxoDeCaixa` (projeção)

- [ ] **Repositories Financeiro**
  - [ ] `ReceitaRepository`
  - [ ] `DespesaRepository`
  - [ ] `CategoriaRepository`

- [ ] **Use Cases Financeiro Básicos**
  - [ ] `CreateReceitaUseCase`
  - [ ] `ListReceitasUseCase` (com filtro por período, categoria)
  - [ ] `CreateDespesaUseCase`
  - [ ] `ListDespesasUseCase`
  - [ ] `GetTotalReceitaPeriodoUseCase`
  - [ ] `GetTotalDespesaPeriodoUseCase`

- [ ] **HTTP Handlers Financeiro**
  - [ ] POST `/financial/receitas` - Criar
  - [ ] GET `/financial/receitas?from=...&to=...` - Listar
  - [ ] PUT `/financial/receitas/{id}` - Atualizar
  - [ ] DELETE `/financial/receitas/{id}` - Deletar
  - [ ] (Idem para despesas)

- [ ] **Migrations Financeiro**
  - [ ] Tabelas `receitas`, `despesas`, `categorias`, `metodos_pagamento`
  - [ ] Índices em `tenant_id`, `data`

#### 2.5 - Padronização de DTO e Responses

- [ ] **Estrutura de Response**
  - [ ] Sucesso: `{ "data": {...}, "meta": {...} }`
  - [ ] Erro: `{ "error": {...}, "trace_id": "..." }`
  - [ ] Paginação: `{ "data": [...], "pagination": { "total", "page", "per_page" } }`

- [ ] **Mapeamento Domain ↔ DTO**
  - [ ] Mapper functions centralizadas
  - [ ] Exemplo: `ReceitaMapper.ToDTO()`, `ReceitaMapper.ToDomain()`

#### 2.6 - Testes & Qualidade

- [ ] **Unit Tests**
  - [ ] Domain layer (100% coverage)
  - [ ] Use cases (>80% coverage)
  - [ ] Repositories (>70% coverage)

- [ ] **Integration Tests**
  - [ ] Handlers com database real
  - [ ] Auth flow completo

- [ ] **Code Quality**
  - [ ] `go fmt` ok
  - [ ] `go vet` passando
  - [ ] golangci-lint com zero erros
  - [ ] Coverage >70%

### Entregas

- ✅ Backend Go estruturado em Clean Architecture
- ✅ Autenticação JWT funcional
- ✅ Multi-tenant implementado
- ✅ Módulo financeiro básico
- ✅ Testes com boa cobertura
- ✅ Deploy em staging validado

---

## 🟦 FASE 3 - Módulos Críticos (2-4 semanas)

**Objetivo:** Portar funcionalidades críticas (Financeiro + Assinaturas)

### Checklist

#### 3.1 - Módulo de Assinaturas

- [ ] **Domain Assinaturas**
  - [ ] Entidades:
    - [ ] `PlanoDeassinatura` (id, nome, valor, periodicidade)
    - [ ] `Assinatura` (id, tenant_id, plan_id, barbeiro_id, status, data_inicio, data_fim)
    - [ ] `AssinaturaInvoice` (id, assinatura_id, valor, status_asaas, data_vencimento)
  - [ ] Value Objects:
    - [ ] `Periodicidade` (enum: MENSAL, TRIMESTRAL, ANUAL)
    - [ ] `StatusAssinatura` (enum: ATIVA, CANCELADA, SUSPENSA)

- [ ] **Asaas Integration**
  - [ ] Biblioteca HTTP client
  - [ ] Serviço `AsaasClient`
    - [ ] `CreateSubscription()`
    - [ ] `CancelSubscription()`
    - [ ] `ListInvoices()`
    - [ ] `GetInvoiceDetails()`
  - [ ] Tratamento de erros específicos
  - [ ] Retry logic com backoff exponencial

- [ ] **Repositories Assinaturas**
  - [ ] `PlanoDeassinaturaRepository`
  - [ ] `AssinaturaRepository`
  - [ ] `AssinaturaInvoiceRepository`

- [ ] **Use Cases Assinaturas**
  - [ ] `CreateAssinaturaUseCase`
  - [ ] `ListAssinaturasUseCase`
  - [ ] `CancelAssinaturaUseCase`
  - [ ] `SincronizarAssinaturasComAsaasUseCase`
  - [ ] `SincronizarFaturasAsaasUseCase`

- [ ] **HTTP Handlers Assinaturas**
  - [ ] POST `/subscriptions/plans` - Criar plano
  - [ ] GET `/subscriptions/plans` - Listar planos
  - [ ] POST `/subscriptions` - Criar assinatura
  - [ ] GET `/subscriptions` - Listar assinaturas do tenant
  - [ ] PUT `/subscriptions/{id}` - Atualizar
  - [ ] DELETE `/subscriptions/{id}` - Cancelar
  - [ ] POST `/subscriptions/sync-asaas` - Sincronizar (admin)

- [ ] **Migrations Assinaturas**
  - [ ] `planos_assinatura`, `assinaturas`, `faturas_assinatura`

#### 3.2 - Módulo de Fluxo de Caixa

- [ ] **Domain Fluxo de Caixa**
  - [ ] Entidades:
    - [ ] `FluxoDeCaixa` (id, tenant_id, saldo_inicial, entradas, saidas, saldo_final, periodo)
  - [ ] Services:
    - [ ] `CalculoFluxoCaixa` (agregar receitas + despesas + assinaturas)

- [ ] **Use Cases Fluxo de Caixa**
  - [ ] `GetFluxoDeCaixaPeriodoUseCase`
  - [ ] `GetFluxoDeCaixaProjecaoUseCase` (próximos 30 dias)

- [ ] **HTTP Handlers**
  - [ ] GET `/financial/cashflow?from=...&to=...` - Fluxo histórico
  - [ ] GET `/financial/cashflow/projection` - Projeção

#### 3.3 - Cron Jobs em Go

- [ ] **Scheduler Setup**
  - [ ] Biblioteca `robfig/cron/v3`
  - [ ] Arquivo `internal/infrastructure/scheduler/scheduler.go`
  - [ ] Inicialização no main.go
  - [ ] Graceful shutdown

- [ ] **Cron Job: Sincronizar Asaas (Diário)**
  - [ ] Schedule: `0 2 * * *` (2h da manhã)
  - [ ] Função:
    - [ ] Buscar faturas não sincronizadas no Asaas
    - [ ] Mapear para Receitas
    - [ ] Persistir no banco
    - [ ] Log de execução
  - [ ] Retry em caso de falha
  - [ ] Notificação se falhar 3x

- [ ] **Cron Job: Gerar Snapshot Financeiro (Diário)**
  - [ ] Schedule: `0 3 * * *` (3h da manhã)
  - [ ] Por tenant:
    - [ ] Calcular fluxo do dia anterior
    - [ ] Armazenar em tabela de snapshots
    - [ ] Detectar anomalias (queda > 50%)

- [ ] **Cron Job: Alertas (Diário)**
  - [ ] Schedule: `0 8 * * *` (8h da manhã)
  - [ ] Regras:
    - [ ] Receita 0 no período
    - [ ] Despesas > receitas
    - [ ] Faturas vencidas não pagas
  - [ ] Enviar (futuro): Email ou Telegram

- [ ] **Cron Job: Limpeza de Dados (Semanal)**
  - [ ] Schedule: `0 4 * * 0` (segundas 4h)
  - [ ] Limpar logs antigos
  - [ ] Arquivar dados históricos (opcional)

- [ ] **Monitoring de Cron**
  - [ ] Log estruturado de cada execução
  - [ ] Tabela `cron_executions` para auditoria
  - [ ] Alert se cron não rodar em 25h

#### 3.4 - Integração com Asaas

- [ ] **Setup**
  - [ ] Conta Asaas criada
  - [ ] API Key obtida
  - [ ] Armazenar em variável de env: `ASAAS_API_KEY`
  - [ ] Documentar em `INTEGRACOES_ASAAS.md`

- [ ] **Fluxos de Integração**
  - [ ] Criar assinatura barbeiro → Asaas
  - [ ] Cancelar → Asaas
  - [ ] Sincronizar faturas pendentes
  - [ ] Webhook (futuro) para confirmações

- [ ] **Error Handling**
  - [ ] 401: API key inválida
  - [ ] 429: Rate limit
  - [ ] 5xx: Retry com backoff
  - [ ] Timeout > 30s: Registrar e alertar

#### 3.5 - Módulo Lista da Vez (Novo)

- [ ] **Domain Lista da Vez**
  - [ ] Entidade `BarbersTurnList`
  - [ ] Tabela `barbers_turn_list`
  - [ ] Tabela `barber_turn_history`
  - [ ] Lógica de pontuação e reordenação

- [ ] **Use Cases**
  - [ ] `GetNextBarberUseCase`
  - [ ] `RecordTurnUseCase`
  - [ ] `ResetTurnListUseCase` (Cron Mensal)

- [ ] **Integração Frontend**
  - [ ] Página de visualização da fila
  - [ ] Botão de registrar atendimento
  - [ ] Histórico de atendimentos

### Entregas

- ✅ Módulo de Assinaturas funcional
- ✅ Integração com Asaas operacional
- ✅ Fluxo de Caixa calculado
- ✅ Cron jobs em execução
- ✅ Sincronizações automáticas testadas
- ✅ Monitoramento de crons configurado

---

## 🟦 FASE 4 - Frontend 2.0 (2-4 semanas) [Paralelo a FASE 3]

**Objetivo:** Frontend em Next.js apontando para novo backend Go

### Checklist

#### 4.1 - Setup Next.js

- [ ] **Criar projeto**
  - [ ] `npx create-next-app@latest barber-analytics-frontend`
  - [ ] Opções:
    - [ ] TypeScript: Sim
    - [ ] ESLint: Sim
    - [ ] Tailwind: Sim
    - [ ] App Router: Sim
    - [ ] Src dir: Sim

- [ ] **Configuração Base**
  - [ ] `next.config.js` com:
    - [ ] Image optimization
    - [ ] Compression
    - [ ] Headers de segurança
  - [ ] `.env.local.example`
    - [ ] `NEXT_PUBLIC_API_URL`
    - [ ] `NEXT_PUBLIC_APP_ENV`

#### 4.2 - Estrutura de Projeto

- [ ] Criar diretórios:
  - [ ] `app/(auth)` - Páginas públicas
  - [ ] `app/(dashboard)` - Páginas protegidas
  - [ ] `app/api` - API routes (optional)
  - [ ] `components/` - Componentes reutilizáveis
  - [ ] `lib/` - Utils, types, constantes
  - [ ] `hooks/` - Custom React hooks
  - [ ] `types/` - TypeScript types globais

#### 4.3 - Autenticação

- [ ] **Login Page** (`app/(auth)/login`)
  - [ ] Form com email/senha
  - [ ] Validação com Zod
  - [ ] Chamada para `/auth/login`
  - [ ] Armazenar tokens (localStorage ou cookies httpOnly)
  - [ ] Redirect para dashboard se já autenticado

- [ ] **Refresh Token**
  - [ ] Interceptor Axios/Fetch
  - [ ] Detectar 401
  - [ ] Chamar `/auth/refresh`
  - [ ] Retry requisição original

- [ ] **Protected Routes**
  - [ ] Middleware Next.js
  - [ ] Verificar token na navegação
  - [ ] Redirect para login se inválido

#### 4.4 - Layout & Navigation

- [ ] **Root Layout**
  - [ ] Providers (TanStack Query, etc)
  - [ ] Fonts Google
  - [ ] Metadata SEO base

- [ ] **Dashboard Layout**
  - [ ] Sidebar com menu
  - [ ] Topbar com tenant selector
  - [ ] User dropdown (profile, logout)
  - [ ] Responsive design

- [ ] **Menu Items**
  - [ ] Dashboard
  - [ ] Financeiro (Receitas, Despesas, Fluxo de Caixa)
  - [ ] Assinaturas
  - [ ] Estoque (link, não implementado ainda)
  - [ ] Configurações

#### 4.5 - Páginas Críticas

- [ ] **Dashboard** (`app/(dashboard)/dashboard`)
  - [ ] Cards com KPIs:
    - [ ] Receita total mês
    - [ ] Despesa total mês
    - [ ] Saldo em caixa
    - [ ] Assinantes ativos
  - [ ] Gráficos (Chart.js ou Recharts):
    - [ ] Receita x Despesa últimos 12 meses
    - [ ] Fluxo de caixa diário

- [ ] **Receitas** (`app/(dashboard)/financial/receitas`)
  - [ ] Tabela com paginação
  - [ ] Filtros: período, categoria
  - [ ] Ações: criar, editar, deletar
  - [ ] Form modal para criação/edição

- [ ] **Despesas** (`app/(dashboard)/financial/despesas`)
  - [ ] Idem receitas

- [ ] **Fluxo de Caixa** (`app/(dashboard)/financial/cashflow`)
  - [ ] Tabela com período selecionável
  - [ ] Colcolunas: Saldo Inicial, Entradas, Saídas, Saldo Final
  - [ ] Visualização por dia/semana/mês

- [ ] **Assinaturas** (`app/(dashboard)/subscriptions`)
  - [ ] Lista de planos disponíveis
  - [ ] Lista de assinantes
  - [ ] Status por assinante
  - [ ] Ações: cancelar, renovar

#### 4.6 - Integração com Backend

- [ ] **API Client**
  - [ ] `lib/api/client.ts` com Axios ou Fetch
  - [ ] Base URL configurável
  - [ ] Interceptor para JWT
  - [ ] Error handling centralizado

- [ ] **React Query (TanStack Query)**
  - [ ] `lib/queries/` para queries reutilizáveis
  - [ ] `lib/mutations/` para mutações
  - [ ] Cache invalidation automática
  - [ ] Retry logic

- [ ] **Hooks Customizados**
  - [ ] `useAuth()` - estado de autenticação
  - [ ] `useTenant()` - tenant atual
  - [ ] `useReceitas()` - lista de receitas
  - [ ] `useDespesas()` - lista de despesas
  - [ ] etc.

#### 4.7 - UI Components

- [ ] Usar `shadcn/ui` para componentes base:
  - [ ] Button, Input, Form
  - [ ] Table, Dialog, Sheet
  - [ ] Card, Badge, Alert
  - [ ] Skeleton (loading states)

- [ ] Customizar com Tailwind:
  - [ ] Cores da marca
  - [ ] Tipografia
  - [ ] Espaçamentos

#### 4.8 - Testes & QA

- [ ] **Unit Tests** (Jest)
  - [ ] Hooks (>80% coverage)
  - [ ] Utils

- [ ] **Integration Tests** (Cypress ou Playwright)
  - [ ] Login flow
  - [ ] Criar receita flow
  - [ ] Filtros e paginação

### Entregas

- ✅ Frontend Next.js estruturado
- ✅ Páginas críticas implementadas
- ✅ Autenticação e proteção de rotas
- ✅ Integração com backend Go validada
- ✅ Responsividade testada
- ✅ Deploy em staging

---

## 🟦 FASE 5 - Migração Progressiva do MVP 1.0 (2-4 semanas)

**Objetivo:** Desligar gradualmente funcionalidades do MVP e migrar para v2

### Checklist

#### 5.1 - Estratégia de Migração

- [ ] **Modo Beta**
  - [ ] Feature flag para ativar v2 por tenant
  - [ ] Admin dashboard mostra quem usa v1 vs v2
  - [ ] Rollback rápido se necessário

- [ ] **Data Migration**
  - [ ] Script para copiar dados do MVP → v2
  - [ ] Validação de integridade
  - [ ] Backup antes da migração

#### 5.2 - Módulo por Módulo

**Financeiro (Receitas & Despesas)**

- [ ] Copiar dados históricos
- [ ] Ativar v2 para novos registros
- [ ] Frontend mostra dados de ambas as versões (transitoriamente)
- [ ] Depois de 1 mês estável: desativar leitura do MVP

**Assinaturas**

- [ ] Copiar planos e assinantes
- [ ] Sincronizar com Asaas
- [ ] Migrar para interface v2
- [ ] Testar crons de sincronização

**Estoque (futuro)**

- [ ] Mapear para novo schema
- [ ] Migrar incrementalmente

#### 5.3 - Validações

- [ ] Verificar:
  - [ ] Totais de receita/despesa coincidem
  - [ ] Assinaturas ativas corretas
  - [ ] Saldos batem

- [ ] Testes de regressão:
  - [ ] Relatórios geram correto
  - [ ] Cálculos de comissão corretos
  - [ ] Filtros funcionam

### Entregas

- ✅ MVP 1.0 e v2 rodando em paralelo
- ✅ Dados migrados com integridade
- ✅ Beta phase completa
- ✅ Pronto para desativação do MVP

---

## 🟦 FASE 6 - Hardening: Segurança, Observabilidade e Escala (1-2 semanas)

**Objetivo:** SaaS profissional, pronto para vender

### Checklist

#### 6.1 - Segurança Aplicacional

- [ ] **Auth Hardening**
  - [ ] Rate limit login (3 tentativas/15min)
  - [ ] 2FA (futuro)
  - [ ] Auditoria de login (IP, device)
  - [ ] Sessões concorrentes limitadas

- [ ] **Autorização**
  - [ ] RBAC implementado e testado
  - [ ] Policies por contexto
  - [ ] Exemplo: Barbeiro só vê seu próprio histórico

- [ ] **Auditoria**
  - [ ] Tabela `audit_logs`
  - [ ] Registrar: WHAT, WHO, WHEN, WHERE
  - [ ] Exemplo: `user:123 UPDATED receita:456 FROM 100.00 TO 150.00 AT 2024-11-14 10:30:00`
  - [ ] Retenção: 90 dias (configurável)

- [ ] **Data Protection**
  - [ ] Encriptação de dados sensíveis (CPF, etc)
  - [ ] PII masking em logs
  - [ ] Backup criptografado

#### 6.2 - Rate Limiting & DDoS

- [ ] **NGINX Rate Limiting**
  - [ ] Limite global: 100 req/s
  - [ ] Por IP: 30 req/s
  - [ ] Por usuário: 50 req/s para endpoints sensíveis

- [ ] **Aplicação**
  - [ ] Endpoint-level rate limit customizado
  - [ ] Queue para processamento heavy

#### 6.3 - Observabilidade

- [ ] **Prometheus**
  - [ ] Métricas Go built-in
  - [ ] Métricas customizadas:
    - [ ] Requisições por status
    - [ ] Latência por endpoint
    - [ ] Erros por tipo
    - [ ] Cron executions
  - [ ] Scrape interval: 15s

- [ ] **Grafana**
  - [ ] Dashboards:
    - [ ] Overview (uptime, requests, errors)
    - [ ] Backend (latência, throughput, memory)
    - [ ] Crons (última execução, duração)
    - [ ] Database (conexões, queries lentas)
  - [ ] Alertas para anomalias

- [ ] **Logs Centralizados**
  - [ ] Opção 1: Loki + Grafana
  - [ ] Opção 2: Axiom ou Datadog
  - [ ] Estruturado em JSON
  - [ ] Trace ID em requests

- [ ] **Sentry (APM)**
  - [ ] Backend Go integrado
  - [ ] Frontend integrado
  - [ ] Alertas para exceções
  - [ ] Performance monitoring

#### 6.4 - Testes & Validação

- [ ] **Load Testing**
  - [ ] Simular 100 concurrent users
  - [ ] Verificar latência < 500ms p95
  - [ ] Verificar error rate < 0.1%
  - [ ] Ferramentas: k6 ou Locust

- [ ] **Security Testing**
  - [ ] OWASP Top 10 checklist
  - [ ] SQL Injection: não vulnerável
  - [ ] XSS: não vulnerável
  - [ ] CSRF: proteção ativa
  - [ ] Auth bypass: não possível

- [ ] **Backup & Disaster Recovery**
  - [ ] Backup automático diário (Neon)
  - [ ] Restore testado semanalmente
  - [ ] RTO: < 2h
  - [ ] RPO: 24h

#### 6.5 - Performance

- [ ] **Database Optimization**
  - [ ] Índices em colunas de filtro
  - [ ] Queries com EXPLAIN ANALYZE
  - [ ] Nenhuma query > 1s sem motivo
  - [ ] Caching de leitura (Redis futuro)

- [ ] **Backend**
  - [ ] Paginação em listas (default 50)
  - [ ] Lazy loading
  - [ ] Bulk operations para imports
  - [ ] Compressão gzip ativa

- [ ] **Frontend**
  - [ ] Code splitting automático
  - [ ] Lazy load components
  - [ ] Image optimization
  - [ ] CSS-in-JS minimizado

#### 6.6 - Documentação & Runbooks

- [ ] **Runbooks para Incidentes**
  - [ ] Database down
  - [ ] API slow
  - [ ] Cron failing
  - [ ] Memory leak

- [ ] **Playbook de Escalação**
  - [ ] P1: Inacessível, perda de dados
  - [ ] P2: Lento, não funciona feature crítica
  - [ ] P3: Bug menor, performance degradada

#### 6.7 - Compliance & Regulamentações

- [ ] **LGPD** (Lei Geral de Proteção de Dados - Brasil)
  - [ ] Privacidade policy
  - [ ] Consentimento para coleta
  - [ ] Right to be forgotten implementado
  - [ ] Data portability (exportar dados)

- [ ] **Termos de Serviço**
  - [ ] SLA definido
  - [ ] Retenção de dados
  - [ ] Responsabilidades

### Entregas

- ✅ Plataforma com alta segurança
- ✅ Observabilidade completa
- ✅ Performance otimizada
- ✅ Compliance atendido
- ✅ Pronto para produção em escala

---

## 📈 Timeline Consolidada

| Fase | Duração | Datas Estimadas | Status |
|------|---------|-----------------|--------|
| 0 | 1-3 dias | Nov 14-17 | ⏳ Próxima |
| 1 | 3-7 dias | Nov 17-24 | ⏳ Próxima |
| 2 | 7-14 dias | Nov 24 - Dec 8 | ⏳ Próxima |
| 3 | 14-28 dias | Dec 8 - Jan 5 | ⏳ Próxima |
| 4 | 14-28 dias | Dec 1 - Jan 5 | ⏳ (paralelo a 3) |
| 5 | 14-28 dias | Jan 5 - Feb 2 | ⏳ Próxima |
| 6 | 7-14 dias | Feb 2 - Feb 16 | ⏳ Próxima |
| **TOTAL** | **~8-12 semanas** | Nov 14 - Feb 16 | 🎯 |

---

## 🎯 Critérios de Sucesso

- ✅ MVP 2.0 em produção
- ✅ Zero dados perdidos na migração
- ✅ Performance: p95 latency < 500ms
- ✅ Uptime: 99.5% SLA
- ✅ Segurança: OWASP compliant
- ✅ Compliance: LGPD atendida
- ✅ Observabilidade: Todos os erros rastreados
- ✅ Documentação: 100% das APIs documentadas

---

**Documento criado em:** 14/11/2025  
**Próxima revisão:** 28/11/2025  
**Status:** ✅ PRONTO PARA IMPLEMENTAÇÃO
