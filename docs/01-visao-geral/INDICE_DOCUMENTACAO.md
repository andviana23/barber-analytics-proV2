> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🗺️ Índice de Documentação - Barber Analytics Pro v2.0

**Propósito:** Navegação rápida entre documentos técnicos e estratégicos.

## 📊 Status & Visão Geral
- [RESUMO_EXECUTIVO.md](./RESUMO_EXECUTIVO.md) — status atual e próximos passos.
- [ROADMAP_IMPLEMENTACAO_V2.md](./ROADMAP_IMPLEMENTACAO_V2.md) — cronograma macro da entrega v2.
- [VISAO_GERAL_PRODUTO.md](./VISAO_GERAL_PRODUTO.md) — visão do SaaS multi-tenant.
- [PERSONAS_E_PERFIS.md](./PERSONAS_E_PERFIS.md) — perfis principais.
- [GLOSSARIO.md](./GLOSSARIO.md) — termos de domínio.

## 🏗️ Arquitetura
- [ARQUITETURA.md](../02-arquitetura/ARQUITETURA.md) — visão de camadas e princípios.
- [FLUXOS_CRITICOS_SISTEMA.md](../02-arquitetura/FLUXOS_CRITICOS_SISTEMA.md) — fluxos essenciais.
- [MODELO_DE_DADOS.md](../02-arquitetura/MODELO_DE_DADOS.md).
- [MODELO_MULTI_TENANT.md](../02-arquitetura/MODELO_MULTI_TENANT.md).
- [DOMAIN_MODELS.md](../02-arquitetura/DOMAIN_MODELS.md).
- ADRs: [ADR-0001-exemplo.md](../02-arquitetura/ADR/ADR-0001-exemplo.md).

## 💻 Backend
- [GUIA_DEV_BACKEND.md](../04-backend/GUIA_DEV_BACKEND.md) — padrões Go e organização.
- [SERVICOS_E_MODULOS.md](../04-backend/SERVICOS_E_MODULOS.md).
- [API_PUBLICA.md](../04-backend/API_PUBLICA.md) e [API_INTERNA.md](../04-backend/API_INTERNA.md).
- [EVENTOS_E_WEBHOOKS.md](../04-backend/EVENTOS_E_WEBHOOKS.md).
- [FEATURE_FLAGS_API.md](../04-backend/FEATURE_FLAGS_API.md).
- Performance: [REDIS_CACHING.md](../04-backend/performance/REDIS_CACHING.md), [QUERY_OPTIMIZATION.md](../04-backend/performance/QUERY_OPTIMIZATION.md), [PERFORMANCE_TASKS_COMPLETE.md](../04-backend/performance/PERFORMANCE_TASKS_COMPLETE.md).

## 🖥️ Frontend
- [GUIA_FRONTEND.md](../03-frontend/GUIA_FRONTEND.md) — padrões React/Next.
- [DESIGN_SYSTEM.md](../03-frontend/DESIGN_SYSTEM.md) — tokens e componentes.
- [MAPA_TELAS.md](../03-frontend/MAPA_TELAS.md).
- [COMPONENTES_CRITICOS.md](../03-frontend/COMPONENTES_CRITICOS.md).

## ⚙️ Ops & SRE
- [INFRA_OVERVIEW.md](../05-ops-sre/INFRA_OVERVIEW.md).
- [CI_CD_PIPELINE.md](../05-ops-sre/CI_CD_PIPELINE.md).
- [GUIA_DEVOPS.md](../05-ops-sre/GUIA_DEVOPS.md) e [GITHUB_SECRETS_SETUP.md](../05-ops-sre/GITHUB_SECRETS_SETUP.md).
- [FLUXO_CRONS.md](../05-ops-sre/FLUXO_CRONS.md).
- [MONITORING_E_ALERTAS.md](../05-ops-sre/MONITORING_E_ALERTAS.md).
- [AUDIT_LOGS.md](../05-ops-sre/AUDIT_LOGS.md).
- [BACKUP_E_RESTORE.md](../05-ops-sre/BACKUP_E_RESTORE.md).
- [PROXIMOS_PASSOS_E2E.md](../05-ops-sre/PROXIMOS_PASSOS_E2E.md).
- Observabilidade: [RUNBOOK_ALERTS.md](../05-ops-sre/observability/RUNBOOK_ALERTS.md), [prometheus/alert-rules.yml](../05-ops-sre/observability/prometheus/alert-rules.yml), [grafana/README.md](../05-ops-sre/observability/grafana/README.md).

## 🔐 Segurança & Compliance
- [ARQUITETURA_SEGURANCA.md](../06-seguranca/ARQUITETURA_SEGURANCA.md).
- [POLITICA_DE_DADOS.md](../06-seguranca/POLITICA_DE_DADOS.md).
- [CONTROLES_DE_ACESSO_INTERNO.md](../06-seguranca/CONTROLES_DE_ACESSO_INTERNO.md).
- [RBAC.md](../06-seguranca/RBAC.md).
- [COMPLIANCE_LGPD.md](../06-seguranca/COMPLIANCE_LGPD.md).
- [SECURITY_TESTING.md](../06-seguranca/SECURITY_TESTING.md).

## 🛠️ Produto & Funcionalidades
- [CATALOGO_FUNCIONALIDADES.md](../07-produto-e-funcionalidades/CATALOGO_FUNCIONALIDADES.md).
- [FLUXOS_DE_NEGOCIO.md](../07-produto-e-funcionalidades/FLUXOS_DE_NEGOCIO.md).
- [ONBOARDING_CLIENTE.md](../07-produto-e-funcionalidades/ONBOARDING_CLIENTE.md).
- [ASSINATURAS.md](../07-produto-e-funcionalidades/ASSINATURAS.md).
- [MANUAL_SUBSCRIPTION_FLOW.md](../07-produto-e-funcionalidades/MANUAL_SUBSCRIPTION_FLOW.md).
- [ONBOARDING_FLOW_REVIEW.md](../07-produto-e-funcionalidades/ONBOARDING_FLOW_REVIEW.md) e [ONBOARDING_WIZARD_IMPLEMENTATION.md](../07-produto-e-funcionalidades/ONBOARDING_WIZARD_IMPLEMENTATION.md).
- [PLANO_CONTINUACAO_ONBOARDING.md](../07-produto-e-funcionalidades/PLANO_CONTINUACAO_ONBOARDING.md).
- [ROADMAP_PRODUTO.md](../07-produto-e-funcionalidades/ROADMAP_PRODUTO.md) e [CHANGELOG.md](../07-produto-e-funcionalidades/CHANGELOG.md).
- [FINANCEIRO.md](../07-produto-e-funcionalidades/FINANCEIRO.md), [ESTOQUE.md](../07-produto-e-funcionalidades/ESTOQUE.md), [listadavez.md](../07-produto-e-funcionalidades/listadavez.md).
- [HELP_CENTER/README.md](../07-produto-e-funcionalidades/HELP_CENTER/README.md).

## 📈 Negócio & Métricas
- [PLANOS_E_PRECOS.md](../08-negocio-e-metricas/PLANOS_E_PRECOS.md).
- [METRICAS_DE_NEGOCIO.md](../08-negocio-e-metricas/METRICAS_DE_NEGOCIO.md).
- [POLITICAS_COMERCIAIS.md](../08-negocio-e-metricas/POLITICAS_COMERCIAIS.md).

## 🤖 AI & Agentes
- [PROMPTS_OFICIAIS.md](../09-ai-e-agentes/PROMPTS_OFICIAIS.md).
- [AGENTES_E_REGRAS.md](../09-ai-e-agentes/AGENTES_E_REGRAS.md).
- [TOM_2026_KB_INDEX.md](../09-ai-e-agentes/TOM_2026_KB_INDEX.md).

## 📚 Fluxos de Leitura Recomendados
- Onboarding rápido: RESUMO_EXECUTIVO → ARQUITETURA → GUIA_DEV_BACKEND ou GUIA_FRONTEND → DESIGN_SYSTEM.
- Trabalhar em onboarding: ONBOARDING_FLOW_REVIEW → PLANO_CONTINUACAO_ONBOARDING.
- Performance backend: REDIS_CACHING → QUERY_OPTIMIZATION → PERFORMANCE_TASKS_COMPLETE.
- Deploy/ops: INFRA_OVERVIEW → CI_CD_PIPELINE → GUIA_DEVOPS → GITHUB_SECRETS_SETUP.

## 🗂️ Estrutura de Diretórios (alto nível)
```
docs/
├── 01-visao-geral/                  ← índice e visão geral
├── 02-arquitetura/                  ← arquitetura, ADRs, integrações
├── 03-frontend/                     ← design system, guias, mapa de telas
├── 04-backend/                      ← guias, APIs, performance
├── 05-ops-sre/                      ← infra, CI/CD, observabilidade, runbooks
├── 06-seguranca/                    ← arquitetura de segurança, LGPD, RBAC
├── 07-produto-e-funcionalidades/    ← módulos de negócio, onboarding, roadmap
├── 08-negocio-e-metricas/           ← planos, métricas, políticas comerciais
└── 09-ai-e-agentes/                 ← prompts e regras de agentes
```

## 🚀 Atalhos Rápidos
- Nova feature backend: ARQUITETURA → GUIA_DEV_BACKEND → API_PUBLICA/API_INTERNA → MODELO_DE_DADOS.
- Componente UI: DESIGN_SYSTEM → GUIA_FRONTEND → MAPA_TELAS.
- Investigar incidente: MONITORING_E_ALERTAS → observability/RUNBOOK_ALERTS → RUNBOOKS/README.
- Melhorar segurança: ARQUITETURA_SEGURANCA → RBAC → SECURITY_TESTING.
