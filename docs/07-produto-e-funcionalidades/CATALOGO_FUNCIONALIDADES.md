> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Catálogo de Funcionalidades (Produção)

![Visão Real x Planejado](../diagrams/real-vs-roadmap.png)

## Visão Geral
Fotografia das funcionalidades **já entregues em produção** (ou em beta interno). Use este catálogo para validar escopo com times de produto/engenharia antes de prometer novas entregas.

## Módulos e Funcionalidades

### Módulo: Assinaturas
Baseado em `ASSINATURAS.md` (modelagem) e `MANUAL_SUBSCRIPTION_FLOW.md` (fluxo manual atual).

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Cadastro de planos | Criação de planos recorrentes com nome, descrição, valor e periodicidade | `/assinaturas/planos` | Em produção | Fluxo manual, sem gateway |
| Criação de assinatura manual | Associa barbeiro a um plano com datas de início e fatura | `/assinaturas` | Em produção | `origem_dado = manual`, status inicial ATIVA |
| Geração de faturas | Cron `ValidateSubscriptions` gera automaticamente faturas pendentes; também há geração manual | Cron 02:00 + `/assinaturas/{id}/faturas` | Em produção | Calcula `proxima_fatura_data`; sem cobrança automática |
| Registro de pagamento | Registro manual de pagamentos e atualização de status da fatura | `/assinaturas/{id}/faturas` | Em produção | Reconciliação manual (sem webhooks Asaas) |
| Alertas de vencimento | Alertas internos para faturas vencidas ou próximas | App (toast) | Beta | Integração externa de cobrança ainda não aplicada |

### Módulo: Financeiro
Baseado em `FINANCEIRO.md`.

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Receitas CRUD | Cadastro/edição de receitas com categoria, método de pagamento e status | `/financeiro/receitas` | Em produção (v1) | Validações de valor, data e tenant |
| Despesas CRUD | Cadastro/edição de despesas com fornecedor e método de pagamento | `/financeiro/despesas` | Em produção (v1) | Suporta categorias de despesa |
| Categorias e métodos | Gestão de categorias (receita/despesa) e métodos de pagamento | `/financeiro/categorias` | Em produção (v1) | Value Objects de categoria/metodologia |
| Fluxo de caixa | Visão consolidada por período (entradas, saídas, saldo) | `/financeiro/fluxo-caixa` | Beta | Usa snapshots/read model descrito no domínio |
| Integração assinaturas | Receita proveniente de faturas de assinaturas | `/financeiro/receitas` | Beta | Hoje alimentação é manual a partir das faturas |

### Módulo: Onboarding
Baseado em `ONBOARDING_FLOW_REVIEW.md`, `ONBOARDING_WIZARD_IMPLEMENTATION.md` e `PLANO_CONTINUACAO_ONBOARDING.md`.

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Signup multi-tenant | Criação de tenant + usuário owner com login automático (JWT) | `POST /auth/signup` | Em produção | Campos validados, JWT RS256 |
| Wizard de onboarding | Fluxo multi-etapas (bem-vindo → checklist → concluir) | `/onboarding` | Em produção | Hook `useCompleteOnboarding`, UI responsiva |
| Marcar onboarding concluído | Registro do flag `onboarding_completed` para o tenant | `/tenants/onboarding/complete` | Beta (backend pendente) | Endpoint/DI e testes ainda a concluir |

### Módulo: Lista da Vez (Fila de Barbeiros)
Baseado em `listadavez.md`.

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Fila rotativa por pontos | Ordenação automática por `current_points`, `last_turn_at`, nome | `/lista-da-vez` | Beta | Barbeiro com menos pontos atende próximo |
| Registro de atendimento | Incremento de pontos, atualização de `last_turn_at` e reordenação | `/lista-da-vez/atender` | Beta | Persistência por tenant |
| Inicialização da fila | Criação inicial de fila com barbeiros ativos (0 pontos) | `/lista-da-vez` | Em produção | Ordenação inicial por data de cadastro |

### Módulo: Cadastros Operacionais
Origem: catálogos citados no catálogo anterior.

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Cadastro de clientes | CRUD de clientes | `/clientes` | Em produção |  |
| Cadastro de profissionais | CRUD de barbeiros/profissionais | `/profissionais` | Em produção | Usado pela Lista da Vez |
| Cadastro de serviços/produtos | CRUD de serviços vendidos e itens de apoio | `/servicos`, `/produtos` | Em produção | Base para financeiro e estoque futuro |
| Meios de pagamento e cupons | Catálogo de meios de pagamento e cupons/regras | `/configuracoes/pagamentos` | Em produção | Alimenta financeiro/assinaturas |

### Identidade e Segurança
Baseado em `RBAC.md`, `SECURITY_TESTING.md`, `AUDIT_LOGS.md`.

| Funcionalidade | Descrição curta | Telas / Rotas principais | Status | Observações |
| --- | --- | --- | --- | --- |
| Auth JWT RS256 | Login/signup com tokens assinados | `/auth/*` | Em produção | Chaves RSA; expiração configurável |
| RBAC por tenant | Perfis e permissões por role | Global | Em produção | Owners/Managers/Barbers |
| Audit Logs | Registro de ações críticas (CRUD) | Backend | Em produção | Persistência multi-tenant |
| Rate limiting | Proteção básica em rotas sensíveis | Edge/API | Beta | Ajuste fino planejado |

## Referências Rápidas
- `ASSINATURAS.md`, `MANUAL_SUBSCRIPTION_FLOW.md`
- `FINANCEIRO.md`
- `ONBOARDING_FLOW_REVIEW.md`, `ONBOARDING_WIZARD_IMPLEMENTATION.md`, `PLANO_CONTINUACAO_ONBOARDING.md`
- `listadavez.md`
- `RBAC.md`, `AUDIT_LOGS.md`, `SECURITY_TESTING.md`
