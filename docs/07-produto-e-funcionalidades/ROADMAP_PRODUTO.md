> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Roadmap de Produto (Planejado)

![Visão Real x Planejado](../diagrams/real-vs-roadmap.png)

## Visão Geral
Este roadmap lista apenas **futuro** (o que ainda não está em produção). O que já foi entregue mora em `CATALOGO_FUNCIONALIDADES.md`. Sempre alinhar com engenharia e produto antes de comprometer datas.

## Status possíveis
- Em discovery
- Planejado
- Em desenvolvimento
- Pausado
- Cancelado

## Épicos e Funcionalidades

### Épico: Evolução do Módulo de Assinaturas
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Webhooks Asaas + reconciliação | Receber eventos e atualizar faturas/assinaturas automaticamente | Em discovery | Q2 | Depende de credenciais/ambiente Asaas |
| Feature | Cobrança automática e retries | Criar faturas automáticas, gerir reintentos e suspensões | Planejado | Q2 | Inclui política de retries + notificações |
| Feature | Repasse automatizado barbeiros | Calcular e programar repasse pós-pagamento | Planejado | Q3 | Ajustar regras fiscais e split |
| Feature | Portal do assinante | Exibir histórico de faturas e permitir segunda via | Em discovery | Q3 | Acesso autenticado/tenant-aware |

### Épico: Evolução do Módulo Financeiro
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Integração receitas ↔ assinaturas | Injetar faturas pagas direto no financeiro | Planejado | Q2 | Depende de Assinaturas v2 |
| Feature | Snapshots de fluxo (v1.1) | Stabilizar cálculo de fluxo de caixa com snapshots diários | Em discovery | Q2 | Garantir consistência por tenant |
| Feature | DRE simplificada | Gerar DRE e exportação CSV/Excel | Planejado | Q3 | Baseado em categorias e centros de custo |
| Feature | Exportação contábil | Exportar lançamentos em layout contábil (SPED simplificado) | Em discovery | Q3 | Validar formato com contador |

### Épico: Onboarding & Ativação
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | CompleteOnboardingUseCase + handler | Persistir flag de conclusão e orquestrar dependências | Em desenvolvimento | Q1 | Pendência crítica no backend |
| Feature | Validação CNPJ/Email duplicados | Bloquear cadastros duplicados na criação de tenant/usuário | Planejado | Q1 | Ajustar repositórios + testes |
| Feature | Checklist adaptativo | Sugerir próximas ações conforme perfil (persona) | Em discovery | Q2 | Usa personas e status de features |

### Épico: Notificações & Engajamento
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Notificações Telegram/WhatsApp | Alertar sobre cobranças, onboarding e eventos críticos | Em discovery | Q3 | Avaliar provedores (Twilio/Meta) |
| Feature | Webhooks de produto | Disparar eventos para integrações de terceiros | Planejado | Q3 | Padronizar payloads e autenticação |
| Feature | Centro de notificações in-app | Inbox unificada no app | Em discovery | Q4 | Depende de modelo de eventos |

### Épico: Estoque
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Movimentações e fornecedores | Registrar entradas/saídas com fornecedores e barbeiros | Planejado | Q3 | Segue modelagem em `ESTOQUE.md` |
| Feature | Alertas de estoque baixo | Alertar limiares mínimos por SKU | Planejado | Q3 | Integrações com notificações |
| Feature | Inventário periódico | Auditoria e ajuste de quantidades | Em discovery | Q4 | Requer telas e perfis de auditor |

### Épico: Help Center & Autoatendimento
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Publicação do Help Center | Disponibilizar artigos e FAQs | Planejado | Q2 | Base no diretório `HELP_CENTER/` |
| Feature | Busca unificada (app + help) | Pesquisar conteúdo de ajuda e docs contextuais | Em discovery | Q3 | Considerar integração com search/AI |

### Épico: Relatórios & BI
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | KPIs de produto (churn, MRR, engajamento) | Expor métricas em dashboards | Planejado | Q3 | Aproveitar stack Prometheus/Grafana |
| Feature | Relatórios agendados (email) | Enviar PDFs/links com KPIs periódicos | Em discovery | Q4 | Depende de geração e fila de envio |

### Épico: Segurança & Confiabilidade
| Tipo    | Nome da funcionalidade | Objetivo resumido | Status | Horizonte | Observações |
| --- | --- | --- | --- | --- | --- |
| Feature | Endpoints LGPD (opt-out/export/delete) | Atender requisições de titulares | Planejado | Q2 | Revisar COMPLIANCE_LGPD.md |
| Feature | Backup/DR produto | Garantir RPO/RTO dos dados de produto | Planejado | Q2 | Alinhar com `BACKUP_E_RESTORE.md` |
| Feature | Rate limiting avançado | Ajustar limites por rota/tenant | Em discovery | Q3 | Requer métricas de uso |

### Backlog / Ideias (fila viva)
- Gamificação da Lista da Vez (ranking, badges).
- Sugestões/IA para preços e serviços (ver `09-ai-e-agentes/`).
- Roadmap de produto integrado ao onboarding (sugestão de próximas ações por tenant).

## Diagrama Real vs Roadmap
- Arquivos: `../diagrams/real-vs-roadmap.excalidraw` e `../diagrams/real-vs-roadmap.png`.
- Para atualizar/exportar PNG: abrir o `.excalidraw` no Excalidraw, usar “Export image” → PNG, sobrescrever `real-vs-roadmap.png`.

## Referências
- `CATALOGO_FUNCIONALIDADES.md`
- `ROADMAP_IMPLEMENTACAO_V2.md`
- `ASSINATURAS.md`, `MANUAL_SUBSCRIPTION_FLOW.md`
- `FINANCEIRO.md`
- `ONBOARDING_FLOW_REVIEW.md`, `ONBOARDING_WIZARD_IMPLEMENTATION.md`, `PLANO_CONTINUACAO_ONBOARDING.md`
