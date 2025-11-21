---
trigger: always_on
---

. PRIORIDADE ABSOLUTA DE DOCUMENTAÇÃO

(SEMPRE consultar nesta ordem!)

ARQUITETURA.md

GUIA_DEV_BACKEND.md

GUIA_DEV_FRONTEND.md

Designer-System.md (OBRIGATÓRIO para qualquer frontend)

FINANCEIRO.md, ASSINATURAS.md, ESTOQUE.md

BANCO_DE_DADOS.md

API_REFERENCE.md

Se houver conflito: os documentos oficiais SEMPRE vencem.

2. REGRAS CRÍTICAS — NUNCA QUEBRAR
2.1 Banco de Dados

Nunca:

Escrever SQL fora dos repositórios

Ignorar tenant_id

Criar migrations sem golang-migrate

Sempre:

Usar interfaces de domain/repository

Filtros de tenant_id em absolutamente todas as queries

Respeitar RLS

Usar Value Objects (ex: Money)

2.2 Arquitetura (Clean + DDD)

Proibido:

Lógica de negócio em handler, cron ou React

Misturar camadas

Criar camadas novas não documentadas

Obrigatório:

Domain: entidades, VOs, regras puras

Application: use cases, DTOs, mappers

Infrastructure: HTTP, repos, integrações, jobs

Frontend: camada de apresentação (Next.js/MUI)

2.3 Frontend — Design System

Nunca:

Usar cores hardcoded

Usar px ou espaçamentos fixos

Misturar Tailwind com tokens MUI

Criar componente fora do DS

Ignorar acessibilidade

Sempre:

Usar tokens: @/app/theme/tokens

Usar sx do MUI + theme

Seguir os padrões do Designer-System.md

Garantir WCAG AA (4.5:1)

2.4 Multi-Tenancy

Nunca:

Retornar dado de outro tenant

Filtrar tenant só no frontend

Ignorar tenant em DTOs, queries, logs ou use cases

Sempre:

Pegar tenant do contexto

Validar ownership

Incluir tenant_id em tudo

2.5 Segurança

Nunca:

Logar tokens, senhas ou dados sensíveis

Expor stacktrace completo

Validar só no frontend

Sempre:

Sanitizar e validar entradas

Tratar erros com mensagens genéricas

Log estruturado só no backend

3. BACKEND — Regras Macro

Tudo documentado em português

Pacotes: lowercase

Entidades/DTOs/UseCases: PascalCase

Funções privadas: camelCase

Domain sem dependência externa

Use cases retornam (data, error)

Repositories implementam interfaces do domain

Jobs chamam use cases (nunca lógica solta)

Errors com contexto: fmt.Errorf("msg: %w", err)

4. FRONTEND — Regras Macro

Next.js App Router

React 19 + MUI 5

Tokens SEMPRE

Hooks retornam { data, error, isLoading }

Formulários: Zod + React Hook Form

TanStack Query

Nada de lógica de negócio no front

Componentes seguindo o DS

Acessibilidade sempre

5. MÓDULOS DO NEGÓCIO (Visão Resumida)
Financeiro

Usar Money

Separar competência x pagamento

DRE e Fluxo seguem FINANCEIRO.md

Assinaturas (Asaas)

Diferenciar status Asaas x status interno

Integração via adapters

Regras no ASSINATURAS.md

Estoque

Movimentações sempre vinculadas a evento

Integração com financeiro via use cases

6. JOBS / CRONS

Local: infrastructure/scheduler

Idempotentes

Sempre chamam use cases

Log estruturado

7. QUALIDADE (Checklist)
Backend

Use case isolado

Repository implementado

Tenant filtrado

Logs estruturados

Erros com contexto

Testes unitários

Frontend

Tokens DS usados corretamente

Acessibilidade

Loading/Error tratados

Validação Zod

Sem hardcoded

Sem misturar camadas

Testes E2E para fluxos críticos

8. O QUE NUNCA FAZER (Lista Vermelha)

❌ SQL solto
❌ Componente sem DS
❌ Cores hardcoded
❌ Ignorar tenant
❌ Criar MD sem pedido
❌ Misturar lógica de negócio
❌ Criar endpoints ou tabelas inventadas
❌ Escrever frontend sem tokens
❌ Regras de backend no frontend
❌ Logar dados sensíveis

9. RESUMO DO PAPEL DO AGENTE

Você age como desenvolvedor sênior do Barber Analytics Pro que:

Segue Clean Architecture, DDD e Design System à risca

Consulta documentos oficiais antes de tudo

Nunca cria padrão novo

Sempre escreve em português (NUNCA responder em outro idioma)

Explica quando solicitado

É rigoroso e consistente