# 🔎 01 — Contexto e Objetivo

**Última Atualização:** 21/11/2025
**Status:** 🔴 CRÍTICO - Bloqueador de todos os módulos

---

## 📊 Estado Atual do Projeto

### ✅ Banco de Dados (100% Completo)

- **42 tabelas** criadas e migradas
- **Migrations 026-038** aplicadas com sucesso
- **11 tabelas novas** para módulos financeiros, metas e precificação

**Tabelas Novas (Migrations 026-038):**

1. `user_preferences` (LGPD)
2. `dre_mensal` (Financeiro)
3. `fluxo_caixa_diario` (Financeiro)
4. `compensacoes_bancarias` (Financeiro)
5. `metas_mensais` (Metas)
6. `metas_barbeiro` (Metas)
7. `metas_ticket_medio` (Metas)
8. `precificacao_config` (Precificação)
9. `precificacao_simulacoes` (Precificação)
10. `contas_a_pagar` (Financeiro)
11. `contas_a_receber` (Financeiro)

### 🟡 Backend Go (~40% Completo)

**O que existe:**

- ✅ Estrutura base (Clean Architecture + DDD)
- ✅ Autenticação e autorização
- ✅ ~40% das entidades antigas
- ✅ Alguns repositórios básicos

**O que falta:**

- ❌ **19 entidades novas** não existem no domínio
- ❌ **Repositórios** para as 11 tabelas novas
- ❌ **Use cases** para módulos novos
- ❌ **DTOs e Handlers HTTP** para novos endpoints
- ❌ **Cron jobs** configuráveis

### 🟡 Frontend Next.js (~30% Completo)

**O que existe:**

- ✅ Estrutura base (App Router)
- ✅ Autenticação
- ✅ Algumas páginas básicas

**O que falta:**

- ❌ **Services** para consumir novos endpoints
- ❌ **Hooks React Query** para módulos novos
- ❌ **Páginas** para Financeiro, Metas, Precificação

---

## 🔴 Lacunas Críticas (Bloqueadores)

### 1. **Domínio Ausente**

- ❌ Sem entidades Go para as 11 tabelas novas
- ❌ Sem Value Objects (`Money`, `Percentage`, `DMais`, `MesAno`)
- ❌ Sem invariantes de negócio
- ❌ Sem validação multi-tenant

### 2. **Repositórios Inexistentes**

- ❌ Sem implementações PostgreSQL para tabelas novas
- ❌ Sem métodos para novas colunas:
  - `meios_pagamento.d_mais` (compensação)
  - `categorias.tipo_custo` (DRE)
  - `receitas.subtipo` (DRE)

### 3. **Use Cases Não Implementados**

- ❌ Backend não expõe endpoints necessários
- ❌ DTOs e mappers não existem
- ❌ Handlers HTTP não criados

### 4. **Cron Jobs Incorretos**

- ❌ Acessam repositórios diretamente (violação A3)
- ❌ Sem configuração externa
- ❌ Sem logs estruturados

### 5. **Frontend Sem Camada de Consumo**

- ❌ Sem clients para novos endpoints
- ❌ Sem hooks React Query
- ❌ Sem tratamento de erros padronizado

---

## 🎯 Objetivo desta Etapa

**Entregar a base completa** de backend e frontend para **desbloquear** todos os módulos:

- 03-FINANCEIRO
- 04-ESTOQUE
- 05-METAS
- 06-PRECIFICACAO

---

## 📦 Entregas Esperadas

### Backend (Go):

| #   | Entrega              | Descrição                               |
| --- | -------------------- | --------------------------------------- |
| 1   | **Domínio Completo** | 19 entidades + Value Objects            |
| 2   | **Ports**            | Interfaces de repositório               |
| 3   | **Repositórios**     | PostgreSQL via sqlc + testes            |
| 4   | **Use Cases**        | Lógica de negócio para todos os módulos |
| 5   | **HTTP**             | DTOs + Handlers + Rotas                 |
| 6   | **Cron Jobs**        | Jobs agendados configuráveis            |

### Frontend (Next.js):

| #   | Entrega      | Descrição                         |
| --- | ------------ | --------------------------------- |
| 7   | **Services** | Clients HTTP para novos endpoints |
| 8   | **Hooks**    | React Query hooks com cache       |

---

## 📋 Módulos que Serão Desbloqueados

### Financeiro:

- ✅ DRE Mensal (Demonstrativo de Resultados)
- ✅ Fluxo de Caixa Compensado (D+)
- ✅ Contas a Pagar/Receber
- ✅ Compensações Bancárias
- ✅ Comissões Automáticas

### Metas:

- ✅ Metas Mensais
- ✅ Metas por Barbeiro
- ✅ Metas de Ticket Médio
- ✅ Cálculo de Progresso

### Precificação:

- ✅ Configuração de Preços
- ✅ Simulador de Precificação
- ✅ Histórico de Simulações

### Estoque:

- ✅ Entrada/Saída
- ✅ Consumo Automático
- ✅ Inventário

---

## 🔗 Referências Técnicas

### Documentação Principal:

- `../CONCLUIR/00-ANALISE_SISTEMA_ATUAL.md` - Análise completa
- `../CONCLUIR/01-backend-domain-entities.md` - Entidades (3-4 dias)
- `../CONCLUIR/02-backend-repository-interfaces.md` - Ports (2 dias)
- `../CONCLUIR/03-08-resumo-tarefas-restantes.md` - Resumo (17 dias)

### Arquitetura:

- `docs/02-arquitetura/ARQUITETURA.md`
- `docs/02-arquitetura/FLUXOS_CRITICOS_SISTEMA.md`
- `docs/02-arquitetura/MODELO_DE_DADOS.md`
- `docs/DIAGRAMA_DEPENDENCIAS_COMPLETO.md`

### Banco de Dados:

- `../DATABASE_MIGRATIONS_COMPLETED.md`

---

## ⏱️ Estimativa

**Total:** ~23 dias úteis (~5 semanas)

Detalhado em `03-sprint-plan.md`

---

**Próximo:** Leia `02-backlog.md` para ver todas as tarefas detalhadas
