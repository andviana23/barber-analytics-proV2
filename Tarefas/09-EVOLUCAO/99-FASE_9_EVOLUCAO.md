# 🚀 FASE 9 — Evolução & Novas Funcionalidades

**Objetivo:** Expandir produto com base em feedback, escalar operação e preparar crescimento
**Duração:** Contínua (3-6 meses iniciais)
**Dependências:** ✅ Fase 8 completa (Sistema estável em produção)
**Sprint:** Sprints 17+ (Ciclos de 2 semanas)
**Período:** Fevereiro 2026 em diante

---

## 📊 Visão Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 9: EVOLUÇÃO CONTÍNUA DO PRODUTO                       │
├─────────────────────────────────────────────────────────────┤
│  Modelo:     🔄 Sprints de 2 semanas (iterativo)            │
│  Status:     ⏳ Aguardando Fase 8                           │
│  Prioridade: 🟡 MÉDIA (features não críticas)              │
│  Período:    Fev-Jul 2026 (inicial)                        │
│  Estratégia: Product-Led Growth + Data-Driven              │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 Objetivos Estratégicos (Q1-Q2 2026)

### 1. **Product Market Fit** (PMF)
- ✅ Validar PMF via NPS >8 e Churn <10%
- ✅ Identificar segmentos de clientes mais engajados
- ✅ Refinar proposta de valor com base em dados

### 2. **Crescimento Sustentável**
- 🎯 Meta: 100 clientes ativos até Julho 2026
- 🎯 MRR: R$ 10.000 (média R$ 100/cliente)
- 🎯 CAC < R$ 20, LTV > R$ 1.500
- 🎯 Payback period < 2 meses

### 3. **Excelência Operacional**
- 📊 Uptime >99.5% consistente
- 🚀 Deploy frequency: 2-3x/semana
- 🐛 Bug resolution time: Médio <48h
- ⚡ Performance: p95 <300ms (melhoria contínua)

### 4. **Preparação para Escala**
- 💰 Unit Economics validados
- 🏗️ Infraestrutura auto-scaling
- 📈 Processo de onboarding automatizado
- 🤖 Suporte parcialmente automatizado (chatbot)

---

## 📋 Roadmap por Trimestre

### Q1 2026 (Fev-Mar-Abr): Funcionalidades Essenciais

**Prioridade 1 (Must Have):**
1. **Relatórios PDF** (Sprint 17-18)
   - Export de receitas, despesas, fluxo de caixa
   - Personalização com logo da barbearia
   - Envio automático por email (mensal)

2. **Gráficos de Evolução** (Sprint 18)
   - Receita mensal (últimos 12 meses)
   - Despesa mensal (últimos 12 meses)
   - Comparativo mês a mês
   - Top 5 categorias

3. **Notificações In-App** (Sprint 19)
   - Assinatura próxima de vencer
   - Meta de receita atingida
   - Despesa acima da média
   - Novo usuário adicionado ao tenant

**Prioridade 2 (Should Have):**
4. **Dashboard Personalizado** (Sprint 20)
   - Widgets customizáveis
   - Filtros persistentes
   - Modo comparação (mês vs mês)

5. **Exportação de Dados** (Sprint 20)
   - CSV de receitas, despesas
   - Excel com múltiplas abas
   - API para integrações

**Prioridade 3 (Nice to Have):**
6. **Dark Mode Melhorado** (Sprint 21)
   - Preferência por usuário
   - Troca automática por horário

7. **Mobile App (PWA)** (Sprint 22-23)
   - Instalar como app nativo
   - Offline-first (sync quando online)
   - Push notifications

---

### Q2 2026 (Mai-Jun-Jul): Expansão & Diferenciação

**Prioridade 1 (Diferenciadores):**
1. **Agendamento Online** (Sprint 24-27) - GRANDE FEATURE
   - Calendário de horários
   - Booking página pública
   - Confirmação por WhatsApp
   - Lembretes automáticos

2. **Integração WhatsApp Business** (Sprint 28-29)
   - Envio de recibos automáticos
   - Lembretes de assinaturas
   - Notificações de pagamento

3. **Comissões Avançadas** (Sprint 30)
   - Configuração de % por barbeiro
   - Cálculo automático
   - Relatório de comissões
   - Exportação para folha de pagamento

**Prioridade 2 (Escalabilidade):**
4. **Multi-unidade** (Sprint 31-32)
   - Tenant com múltiplas lojas
   - Dashboard consolidado
   - Comparativo entre unidades

5. **Planos de Assinatura Flexíveis** (Sprint 33)
   - Freemium (básico grátis)
   - Pro (R$ 79/mês)
   - Enterprise (R$ 149/mês + onboarding dedicado)

**Prioridade 3 (Integrações):**
6. **API Pública** (Sprint 34)
   - Webhooks para eventos
   - REST API documentada
   - Rate limiting por tenant

7. **Integração Contábil** (Sprint 35)
   - Conta Azul
   - Omie
   - Exportação XML NF-e

---

## ✅ Features Detalhadas (Top 10)

### 🔴 1. Relatórios PDF (Alta Prioridade)

**Problema:** Clientes precisam enviar relatórios mensais para contador

**Solução:**
- Gerar PDF com receitas, despesas, fluxo de caixa
- Personalização com logo
- Download + Email automático

**Estimativa:** 12h (Sprint 17-18)

**Tasks:**
- T-PDF-001: Setup biblioteca pdfmake/wkhtmltopdf (2h)
- T-PDF-002: Template de relatório (header, footer, tabelas) (3h)
- T-PDF-003: Endpoint GET /reports/pdf?type=financial&period=2026-01 (3h)
- T-PDF-004: Frontend: Botão "Exportar PDF" (2h)
- T-PDF-005: Testes E2E (2h)

**Critérios de Aceitação:**
- [ ] PDF gerado em <3s
- [ ] Tamanho <2 MB
- [ ] Layout profissional (logo, cores)
- [ ] Totalizadores corretos
- [ ] Download + email funcionando

---

### 🔴 2. Gráficos de Evolução (Alta Prioridade)

**Problema:** Clientes querem ver tendências mês a mês

**Solução:**
- Chart.js ou Recharts
- Gráfico de linha: Receita vs Despesa (12 meses)
- Gráfico de barras: Top 5 categorias

**Estimativa:** 8h (Sprint 18)

**Tasks:**
- T-GRAPH-001: Endpoint GET /analytics/monthly-evolution (2h)
- T-GRAPH-002: Endpoint GET /analytics/top-categories (2h)
- T-GRAPH-003: Frontend: MonthlyEvolutionChart component (2h)
- T-GRAPH-004: Frontend: TopCategoriesChart component (2h)

**Critérios de Aceitação:**
- [ ] Gráficos responsivos (mobile + desktop)
- [ ] Tooltip com valores exatos
- [ ] Comparação visual clara
- [ ] Cores acessíveis (WCAG AA)

---

### 🟡 3. Notificações In-App (Média Prioridade)

**Problema:** Clientes perdem eventos importantes (assinatura vencendo, etc)

**Solução:**
- Bell icon no topbar
- Badge com contador
- Dropdown com últimas 10 notificações
- Cron diário gerando notificações

**Estimativa:** 10h (Sprint 19)

**Tasks:**
- T-NOTIF-001: Tabela notifications (migration) (1h)
- T-NOTIF-002: NotificationService (create, list, markAsRead) (2h)
- T-NOTIF-003: Cron job: Gerar notificações (3h)
- T-NOTIF-004: Endpoint GET /notifications (1h)
- T-NOTIF-005: Frontend: NotificationBell component (3h)

**Tipos de Notificações:**
- Assinatura vence em 7 dias
- Meta de receita atingida (configurável)
- Despesa acima da média mensal
- Novo usuário adicionado

---

### 🟡 4. Dashboard Personalizado (Média Prioridade)

**Problema:** Cada cliente tem KPIs diferentes (alguns querem comissões, outros estoque)

**Solução:**
- Drag & drop de widgets
- Configuração salva por usuário
- Widgets disponíveis: KPI Cards, Charts, Recent Activity

**Estimativa:** 14h (Sprint 20)

**Libraries:**
- react-grid-layout (drag & drop)
- Zustand (state de configuração)

---

### 🟢 5. Agendamento Online (Baixa Prioridade, Alto Impacto)

**Problema:** Clientes querem agendar horários online

**Solução:**
- Calendário semanal
- Booking page pública (`/b/barbearia-x`)
- Confirmação por WhatsApp

**Estimativa:** 24h (Sprint 24-27) - ÉPICO

**Complexidade:**
- Domain novo: Agendamento (Appointment, TimeSlot)
- Integração WhatsApp Business API
- Conflito de horários
- Notificações

**Faseamento:**
1. Sprint 24: Backend (domain + CRUD) - 8h
2. Sprint 25: Frontend Admin (calendário) - 6h
3. Sprint 26: Booking page pública - 6h
4. Sprint 27: WhatsApp integration - 4h

---

### 🟢 6. Integração WhatsApp Business (Baixa Prioridade)

**Problema:** Comunicação com clientes é manual

**Solução:**
- WhatsApp Business API (via Twilio/MessageBird)
- Envio automático de recibos
- Lembretes de assinatura

**Estimativa:** 12h (Sprint 28-29)

**Custos:**
- Twilio: ~R$ 0,10/msg
- MessageBird: ~R$ 0,08/msg

---

### 🟡 7. Comissões Avançadas (Média Prioridade)

**Problema:** Cálculo manual de comissões por barbeiro

**Solução:**
- Configuração de % por usuário
- Relatório mensal de comissões
- Exportação para folha

**Estimativa:** 10h (Sprint 30)

**Tasks:**
- T-COMM-001: Campo commission_percentage em users (1h)
- T-COMM-002: Cálculo automático (receita * % / 100) (2h)
- T-COMM-003: Endpoint GET /reports/commissions (2h)
- T-COMM-004: Frontend: Relatório de comissões (3h)
- T-COMM-005: Exportação CSV (2h)

---

### 🟢 8. Multi-unidade (Baixa Prioridade, Alto Esforço)

**Problema:** Redes com múltiplas barbearias querem consolidar dados

**Solução:**
- Tenant com múltiplas "locations"
- Dashboard consolidado
- Filtro por unidade

**Estimativa:** 18h (Sprint 31-32)

**Mudanças de Schema:**
- Nova tabela: locations (tenant_id, nome, endereco)
- Adicionar location_id em receitas, despesas, etc.
- Migração complexa para clientes existentes

---

### 🟡 9. Planos de Assinatura Flexíveis (Média Prioridade)

**Problema:** Barreira de entrada alta (R$ 79/mês)

**Solução:**
- Freemium: Grátis até 50 receitas/mês
- Pro: R$ 79/mês (ilimitado + suporte)
- Enterprise: R$ 149/mês (multi-unidade + onboarding)

**Estimativa:** 8h (Sprint 33)

**Tasks:**
- T-PLAN-001: Tabela subscription_plans (3 tiers) (1h)
- T-PLAN-002: Middleware: Check usage limits (3h)
- T-PLAN-003: Upgrade/downgrade flow (2h)
- T-PLAN-004: Billing com Stripe/Asaas (2h)

---

### 🟢 10. API Pública (Baixa Prioridade)

**Problema:** Clientes querem integrar com outros sistemas

**Solução:**
- API REST documentada (Swagger/OpenAPI)
- API Keys por tenant
- Rate limiting: 100 req/min (Free), 1000 req/min (Pro)
- Webhooks para eventos

**Estimativa:** 16h (Sprint 34)

**Eventos Webhook:**
- `receita.created`
- `despesa.created`
- `assinatura.created`
- `assinatura.cancelled`

---

## 🔄 Processo de Desenvolvimento

### Sprint Planning (a cada 2 semanas)

**Agenda (2h):**
1. Review de métricas do sprint anterior (30 min)
2. Apresentação de feedback de clientes (30 min)
3. Priorização de features (30 min)
4. Task breakdown e estimativas (30 min)

**Saídas:**
- Backlog priorizado
- Sprint goal definido
- Tasks assignadas

### Daily Standup (15 min)

**Perguntas:**
1. O que fiz ontem?
2. O que farei hoje?
3. Algum bloqueio?

### Sprint Review (1h)

**Agenda:**
1. Demo de features concluídas
2. Feedback do time
3. Validação de critérios de aceitação

### Retrospectiva (1h)

**Agenda:**
1. O que foi bem?
2. O que pode melhorar?
3. Ações para próximo sprint

---

## 📊 Métricas de Produto (KPIs)

### Engajamento
- **DAU (Daily Active Users):** Target >30
- **WAU (Weekly):** Target >60
- **MAU (Monthly):** Target >80
- **Stickiness (DAU/MAU):** Target >40%

### Retenção
- **Churn mensal:** Target <10%
- **Retenção D7:** Target >70% (7 dias após signup)
- **Retenção D30:** Target >50%
- **Retenção D90:** Target >40%

### Monetização
- **MRR (Monthly Recurring Revenue):** Track crescimento
- **ARPU (Average Revenue Per User):** Target R$ 100
- **CAC (Customer Acquisition Cost):** Target <R$ 20
- **LTV (Lifetime Value):** Target >R$ 1.500
- **LTV/CAC Ratio:** Target >3

### Feature Adoption
- **Dashboard:** >90% usage
- **Fluxo de Caixa:** >60% usage
- **Relatórios PDF:** Target >40% usage (após lançamento)
- **Agendamento:** Target >30% usage (após lançamento)

---

## 🏗️ Melhorias de Infraestrutura

### Escalabilidade

**Horizontal Scaling:**
- [ ] Backend: Docker Swarm ou Kubernetes
- [ ] Load Balancer: NGINX ou AWS ALB
- [ ] Database: Read replicas (Neon suporta)
- [ ] Redis: Cluster mode (alta disponibilidade)

**Performance:**
- [ ] CDN: Cloudflare para assets estáticos
- [ ] Edge Functions: Vercel Edge para frontend
- [ ] Database: Partitioning por tenant_id (se >1000 tenants)

### Observabilidade

**Logs:**
- [ ] Centralização: Loki ou CloudWatch
- [ ] Retention: 30 dias (compliance)

**Traces:**
- [ ] OpenTelemetry para distributed tracing
- [ ] Jaeger ou Tempo para visualização

**Custos:**
- [ ] Cost Explorer: Monitorar spend AWS/Neon
- [ ] Alertas: Budget >R$ 1.000/mês

### Segurança

**Penetration Testing:**
- [ ] Contratar pentest anual
- [ ] Correção de vulnerabilidades encontradas

**OWASP Top 10:**
- [ ] Revisão trimestral de security checklist
- [ ] Dependency updates automáticos (Dependabot)

---

## 💡 Experimentos & Inovação

### A/B Testing (quando >100 usuários)

**Experimentos Sugeridos:**
1. Onboarding flow: Guiado vs Livre
2. Pricing: R$ 79 vs R$ 99 vs Freemium
3. Dashboard layout: Cards vs Tabs
4. Call-to-Action: "Experimente Grátis" vs "Comece Agora"

**Ferramenta:** PostHog ou Split.io

### Beta Program

**Objetivo:** Testar features antes do lançamento geral

**Processo:**
1. Selecionar 5-10 clientes engajados
2. Apresentar roadmap e convidar para beta
3. Deploy em ambiente `/beta`
4. Coletar feedback semanal
5. Iterar antes do lançamento

---

## 🎓 Aprendizados & Best Practices

### Product Management

**Framework RICE para Priorização:**
- **R**each: Quantos clientes serão impactados?
- **I**mpact: Qual o impacto no negócio? (Low/Medium/High)
- **C**onfidence: Qual nossa confiança na estimativa? (%)
- **E**ffort: Quantas horas de dev?

**Score RICE = (Reach * Impact * Confidence) / Effort**

**Exemplo:**
- Relatórios PDF: (50 * High * 80%) / 12h = 3.33
- Agendamento: (80 * High * 60%) / 24h = 2.00
- Dark Mode: (100 * Low * 90%) / 3h = 3.00

**Prioridade:** PDF > Dark Mode > Agendamento

### Engenharia

**Technical Debt:**
- [ ] Reservar 20% do sprint para refatoração
- [ ] Code review obrigatório (aprovação de 1 dev)
- [ ] Cobertura de testes: Manter >70%

**Documentation:**
- [ ] ADRs (Architecture Decision Records) para decisões importantes
- [ ] Changelogs públicos para transparência
- [ ] API docs auto-geradas (Swagger)

---

## 🚀 Plano de Crescimento (6 meses)

### Fevereiro 2026
- Clientes: 40 → 50 (+25%)
- MRR: R$ 2.900 → R$ 4.000 (+38%)
- Features: PDF + Gráficos

### Março 2026
- Clientes: 50 → 65 (+30%)
- MRR: R$ 4.000 → R$ 5.500 (+38%)
- Features: Notificações + Dashboard Personalizado

### Abril 2026
- Clientes: 65 → 80 (+23%)
- MRR: R$ 5.500 → R$ 7.000 (+27%)
- Features: Exportação + Mobile PWA

### Maio 2026
- Clientes: 80 → 95 (+19%)
- MRR: R$ 7.000 → R$ 8.500 (+21%)
- Features: Agendamento (início)

### Junho 2026
- Clientes: 95 → 110 (+16%)
- MRR: R$ 8.500 → R$ 10.000 (+18%)
- Features: Agendamento (completo) + WhatsApp

### Julho 2026
- Clientes: 110 → 130 (+18%)
- MRR: R$ 10.000 → R$ 12.000 (+20%)
- Features: Comissões + Multi-unidade

**Meta Semestral:**
- ✅ 130 clientes ativos
- ✅ R$ 12.000 MRR
- ✅ Churn <10%
- ✅ 10+ features lançadas

---

## 🎯 Visão de Longo Prazo (12-24 meses)

### Ano 1 (2026)
- **Clientes:** 200+ ativos
- **MRR:** R$ 20.000
- **Equipe:** 3 devs + 1 designer + 1 suporte
- **Produto:** Completo (agendamento, integrações, multi-unidade)

### Ano 2 (2027)
- **Clientes:** 500+ ativos
- **MRR:** R$ 50.000
- **Equipe:** 5 devs + 2 designers + 3 suporte + 1 PM
- **Produto:** Marketplace de integrações, Mobile App nativo, IA preditiva

### Ano 3+ (2028)
- **Clientes:** 1000+ ativos
- **MRR:** R$ 100.000+
- **Expansão:** Outros nichos (salões, pet shops, academias)
- **Exit:** Possível aquisição ou IPO

---

## 📝 Documentação Contínua

### Criações Necessárias
- [ ] `docs/ROADMAP_2026_Q1.md` - Roadmap Q1 detalhado
- [ ] `docs/ROADMAP_2026_Q2.md` - Roadmap Q2 detalhado
- [ ] `docs/ADR/` - Architecture Decision Records
- [ ] `docs/EXPERIMENTS/` - A/B tests e resultados
- [ ] `docs/CHANGELOG_PUBLIC.md` - Changelog público

---

## ✅ Checklist de Transição (Fase 8 → Fase 9)

### Pré-requisitos
- [ ] Fase 8 100% completa
- [ ] 4 semanas de operação estável
- [ ] Uptime médio ≥99%
- [ ] NPS ≥7/10
- [ ] ≥30 clientes ativos
- [ ] Roadmap Fase 9 validado pelo time

### Setup
- [ ] Sprint 17 planejado (tasks criadas no GitHub)
- [ ] Features priorizadas via framework RICE
- [ ] Estimativas técnicas validadas
- [ ] Time alocado (devs + designer + PM)

### Comunicação
- [ ] Roadmap comunicado aos clientes (email + blog)
- [ ] Feedback loop estabelecido (survey mensal)
- [ ] Beta program ativo

---

**Última Atualização:** 17/11/2025 23:59
**Status:** ⏳ Aguardando Fase 8 (Sistema em produção estável)
**Início Previsto:** Fevereiro 2026
**Duração:** Contínua (evolução do produto)
**Modelo:** Sprints de 2 semanas + releases incrementais

---

**🎯 Objetivo Final:** Construir o sistema de gestão #1 para barbearias no Brasil, com 1000+ clientes ativos e produto que vende sozinho (Product-Led Growth).

---

**👥 Time Necessário (Fase 9):**
- 1-2 Backend Developers (Go)
- 1-2 Frontend Developers (Next.js/React)
- 1 Product Manager
- 1 UX/UI Designer
- 1-2 Customer Success (Suporte + Onboarding)
- 1 DevOps/SRE (part-time ou consultoria)
- 1 Marketing/Growth (part-time inicialmente)

**💰 Budget Estimado (mensal):**
- Salários: R$ 25.000 - R$ 40.000
- Cloud (AWS/Neon/Vercel): R$ 1.000 - R$ 3.000
- Ferramentas (GitHub, Slack, Notion): R$ 500
- Marketing/Ads: R$ 2.000 - R$ 5.000
- **Total:** R$ 28.500 - R$ 48.500/mês

**🎯 Break-even:** ~400 clientes (R$ 40.000 MRR)
