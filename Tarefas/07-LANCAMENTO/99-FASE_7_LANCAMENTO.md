# 🚀 FASE 7 — Lançamento & Go-Live

**Objetivo:** Deploy em produção e início de operação comercial do MVP 2.0
**Duração:** 3-5 dias
**Dependências:** ✅ Fase 6 completa (100%)
**Sprint:** Sprint 12
**Data Prevista:** Dezembro 20-26, 2025

---

## 📊 Progresso Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 7: LANÇAMENTO & GO-LIVE                               │
├─────────────────────────────────────────────────────────────┤
│  Progresso:  ░░░░░░░░░░░░░░░░░░░░  0% (0/6 concluídas)     │
│  Status:     ⏳ Aguardando Fase 6                           │
│  Prioridade: 🔴 CRÍTICA                                     │
│  Estimativa: 20 horas                                       │
│  Sprint:     Sprint 12                                      │
│  Go-Live:    26/12/2025 (estimado)                         │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Checklist de Tarefas

### 🔴 T-LAUNCH-001 — Completar Pendências da Fase 6

- **Responsável:** Backend + DevOps + Legal
- **Prioridade:** 🔴 Crítica
- **Estimativa:** 14h
- **Sprint:** Sprint 12
- **Status:** ⏳ Bloqueado (dependência: Fase 6)
- **Deliverable:** T-LGPD-001 + T-OPS-005 concluídos

#### Critérios de Aceitação

- [ ] **T-LGPD-001 Completo** (8h)
  - [ ] DELETE /api/v1/me implementado e testado
  - [ ] GET /api/v1/me/export implementado e testado
  - [ ] Banner de consentimento no frontend funcionando
  - [ ] Job de limpeza de dados (90 dias) agendado
  - [ ] Privacy Policy publicada em `/privacy`
  - [ ] Testes E2E para endpoints LGPD passando
- [ ] **T-OPS-005 Completo** (6h)
  - [ ] GitHub Actions workflow `backup-database.yml` funcionando
  - [ ] S3 bucket configurado com lifecycle policies
  - [ ] Primeiro teste de restore executado com sucesso
  - [ ] Alertas de backup configurados no Prometheus
  - [ ] Documentação de DR validada pela equipe
  - [ ] RTO/RPO targets validados (<2h / <1h)

**Bloqueadores Potenciais:**

- Acesso a AWS S3 (criar bucket + IAM)
- Aprovação legal da Privacy Policy
- Testes de restore podem revelar problemas de schema

---

### 🔴 T-LAUNCH-002 — Checklist Pré-Lançamento

- **Responsável:** Tech Lead + DevOps + QA
- **Prioridade:** 🔴 Crítica
- **Estimativa:** 3h
- **Sprint:** Sprint 12
- **Status:** ⏳ Não iniciado
- **Deliverable:** Checklist 100% validado

#### Critérios de Aceitação

**1. Código & Testes**

- [ ] Backend: `go test ./...` - 100% passando (48+ testes)
- [ ] Frontend: `pnpm test:unit` - 67 testes passando
- [ ] Frontend: `pnpm test:a11y` - 25 testes passando
- [ ] Frontend: `pnpm test:e2e` - ≥80% passando (21/26)
- [ ] Security tests: 35/35 passando
- [ ] Load test k6: p95 <500ms, error rate <0.1%
- [ ] Code coverage: Backend ≥30%, Frontend ≥60%

**2. Infraestrutura**

- [ ] Neon Database: 13 migrations aplicadas
- [ ] NGINX: SSL/TLS válido (Let's Encrypt)
- [ ] Systemd Service: backend rodando estável
- [ ] Redis: Funcionando (opcional, para cache)
- [ ] Prometheus: Coletando métricas
- [ ] Grafana: 4 dashboards configurados
- [ ] Alertmanager: 5 alertas configurados

**3. Funcionalidades Core**

- [ ] Login/Logout funcionando
- [ ] Dashboard com KPIs corretos
- [ ] CRUD Receitas: Create, Read, Update, Delete OK
- [ ] CRUD Despesas: Create, Read, Update, Delete OK
- [ ] Fluxo de Caixa: Cálculo correto de saldo
- [ ] Assinaturas: Listagem e filtros OK
- [ ] Multi-tenant: Isolamento validado (não vaza dados)
- [ ] RBAC: Permissões funcionando (Owner, Manager, Accountant, Employee)

**4. Segurança**

- [ ] JWT RS256: Tokens válidos expirando corretamente
- [ ] Rate limiting: NGINX + Backend ativo
- [ ] HTTPS: Forçado em todas as requisições
- [ ] SQL Injection: Testes passando (não vulnerável)
- [ ] XSS: Testes passando (não vulnerável)
- [ ] CSRF: Proteção ativa
- [ ] Audit logs: Registrando ações críticas
- [ ] LGPD: Endpoints /me/delete e /me/export funcionando

**5. Performance**

- [ ] GET /receitas: <100ms (com índices)
- [ ] GET /dashboard: <200ms (com cache)
- [ ] POST /auth/login: <50ms
- [ ] Database connections: <10 simultâneas (média)
- [ ] Redis cache: Hit rate >70%

**6. Observabilidade**

- [ ] Logs estruturados: JSON formatado, níveis corretos
- [ ] Métricas Prometheus: HTTP, DB, Cron, Business
- [ ] Grafana dashboards: Dados corretos e atualizados
- [ ] Alertas: Testados manualmente (simular falha)
- [ ] Runbook: Equipe treinada em procedimentos

**7. Dados & Backup**

- [ ] Seeds produção: Categorias padrão criadas
- [ ] Backup automático: Rodou nas últimas 24h
- [ ] Restore testado: Sucesso em <2h
- [ ] Neon PITR: Habilitado (7 dias)

**Arquivo Criado:**

- `docs/CHECKLIST_PRE_LAUNCH.md` - Checklist completo copiável

---

### 🔴 T-LAUNCH-003 — Deploy Produção

- **Responsável:** DevOps
- **Prioridade:** 🔴 Crítica
- **Estimativa:** 2h
- **Sprint:** Sprint 12
- **Status:** ⏳ Não iniciado
- **Deliverable:** Sistema rodando em produção

#### Critérios de Aceitação

**1. Preparação de Ambiente**

- [ ] VPS Produção: Ubuntu 22.04 atualizado
- [ ] Variáveis de ambiente configuradas (`.env.production`)
  - [ ] `DATABASE_URL` (Neon prod)
  - [ ] `JWT_PRIVATE_KEY` (RS256, gerado)
  - [ ] `JWT_PUBLIC_KEY`
  - [ ] `REDIS_URL` (se aplicável)
  - [ ] `SENTRY_DSN` (se aplicável)
  - [ ] `CORS_ALLOWED_ORIGINS` (frontend prod)
- [ ] Firewall configurado (portas 80, 443, 22)
- [ ] SSH key-only authentication

**2. Deploy Backend**

```bash
# Executar no servidor
cd /opt/barber-analytics-backend-v2
git pull origin main
make build
sudo systemctl restart barber-api
sudo systemctl status barber-api

# Verificar
curl https://api.barberpro.dev/health
```

- [ ] Build compilado sem erros
- [ ] Systemd service reiniciado
- [ ] Health check retornando 200 OK
- [ ] Logs sem erros críticos

**3. Deploy Frontend**

```bash
# Executar no servidor (ou Vercel)
cd /opt/barber-analytics-frontend-v2
git pull origin main
pnpm install
pnpm build
pnpm start --port 3000

# OU via Vercel:
vercel --prod
```

- [ ] Build Next.js sucesso
- [ ] Apontamento DNS correto (barberpro.dev)
- [ ] SSL/TLS ativo
- [ ] Lighthouse Score: Performance >80, Accessibility >95

**4. Configuração NGINX**

```nginx
# /etc/nginx/sites-available/barber-api
server {
    listen 443 ssl http2;
    server_name api.barberpro.dev;

    ssl_certificate /etc/letsencrypt/live/api.barberpro.dev/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.barberpro.dev/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

- [ ] NGINX configuration validada: `nginx -t`
- [ ] NGINX recarregado: `systemctl reload nginx`
- [ ] SSL válido: `curl https://api.barberpro.dev/health`

**5. Smoke Test Produção**

- [ ] Acesso `https://barberpro.dev` - Página login carrega
- [ ] Login com credenciais demo - Sucesso
- [ ] Dashboard exibe KPIs corretos
- [ ] Criar receita de teste - Sucesso
- [ ] Listar receitas - Aparece na tabela
- [ ] Deletar receita de teste - Removida
- [ ] Logout - Redirecionado para login

**6. Monitoramento 24h**

- [ ] Grafana dashboard aberto em monitor dedicado
- [ ] Slack/Telegram alertas ativos
- [ ] Equipe on-call escalada
- [ ] Runbook impresso/acessível

**Rollback Plan:**

```bash
# Se deploy falhar:
cd /opt/barber-analytics-backend-v2
git checkout <commit-anterior>
make build
sudo systemctl restart barber-api
```

---

### 🟡 T-LAUNCH-004 — Comunicação aos Stakeholders

- **Responsável:** Product Manager + Marketing
- **Prioridade:** 🟡 Média
- **Estimativa:** 1h
- **Sprint:** Sprint 12
- **Status:** ⏳ Não iniciado
- **Deliverable:** Comunicados enviados

#### Critérios de Aceitação

**1. Email para Clientes Beta (se houver)**

- [ ] Assunto: "Barber Analytics Pro v2.0 está no ar! 🚀"
- [ ] Corpo:
  - [ ] Agradecimento pela paciência
  - [ ] Principais novidades (dashboard, fluxo de caixa, assinaturas)
  - [ ] Link para acesso: https://barberpro.dev
  - [ ] Link para tutorial/onboarding
  - [ ] Canal de suporte (email/WhatsApp)
- [ ] Enviado via plataforma de email (SendGrid/Mailchimp)

**2. Post Redes Sociais**

- [ ] Instagram: Carrossel com screenshots do dashboard
- [ ] LinkedIn: Post profissional sobre lançamento
- [ ] WhatsApp Status: Anúncio breve
- [ ] Hashtags: #BarbeariaDigital #Gestão #SaaS

**3. Landing Page Atualizada**

- [ ] Seção "Novidades" com v2.0
- [ ] Depoimentos de beta testers (se houver)
- [ ] CTA: "Experimente Grátis por 14 dias"
- [ ] FAQ atualizado

**4. Documentação Pública**

- [ ] Help Center: Artigos básicos (Como criar receita, Como usar dashboard)
- [ ] Vídeo Tutorial: 3-5 minutos no YouTube
- [ ] Changelog público: `/changelog` com histórico de versões

---

### 🟢 T-LAUNCH-005 — Onboarding de Primeiros Clientes

- **Responsável:** Suporte + Product
- **Prioridade:** 🟢 Baixa (pós-lançamento)
- **Estimativa:** Contínuo
- **Sprint:** Pós-lançamento
- **Status:** ⏳ Não iniciado
- **Deliverable:** Clientes onboarded com sucesso

#### Critérios de Aceitação

**1. Primeiro Contato (em 24h)**

- [ ] Email de boas-vindas automático ao signup
- [ ] Agendamento de call de onboarding (15 min)
- [ ] Envio de tutorial em PDF/vídeo

**2. Call de Onboarding (15-30 min)**

- [ ] Apresentação do dashboard
- [ ] Configuração de categorias personalizadas
- [ ] Cadastro de primeira receita/despesa
- [ ] Explicação de assinaturas (se aplicável)
- [ ] Q&A

**3. Acompanhamento (Semana 1)**

- [ ] Email dia 3: "Como está sendo sua experiência?"
- [ ] Email dia 7: Dicas avançadas (filtros, relatórios)
- [ ] Oferta de call adicional se necessário

**4. Coleta de Feedback**

- [ ] NPS Survey após 14 dias
- [ ] Perguntas específicas:
  - [ ] Qual feature mais útil?
  - [ ] Qual feature faltando?
  - [ ] Dificuldades encontradas?
  - [ ] Nota geral (1-10)
- [ ] Compilar feedback em documento

---

### 🟢 T-LAUNCH-006 — Monitoramento Pós-Lançamento (7 dias)

- **Responsável:** DevOps + Tech Lead
- **Prioridade:** 🟢 Baixa (contínuo)
- **Estimativa:** 1h/dia x 7 dias = 7h
- **Sprint:** Pós-lançamento
- **Status:** ⏳ Não iniciado
- **Deliverable:** Relatório de estabilidade

#### Critérios de Aceitação

**1. Métricas Diárias (Check às 9h, 15h, 21h)**

- [ ] Uptime: >99.5%
- [ ] Latência p95: <500ms
- [ ] Error rate: <0.1%
- [ ] Requests/dia: Monitorar crescimento
- [ ] Novos signups/dia: Monitorar
- [ ] Active users: Monitorar

**2. Alertas Críticos**

- [ ] Responder em <30 min a alertas críticos
- [ ] Registrar incidentes em log
- [ ] Post-mortem se downtime >5 min

**3. Logs Review**

- [ ] Revisar erros 5xx diariamente
- [ ] Investigar erros recorrentes
- [ ] Criar issues para bugs encontrados

**4. Database Health**

- [ ] Conexões ativas: <20
- [ ] Queries lentas: Identificar e otimizar
- [ ] Tamanho do banco: Monitorar crescimento
- [ ] Backup: Validar execução diária

**5. Relatório Semanal**

```markdown
# Relatório Pós-Lançamento - Semana 1

## Métricas

- Uptime: 99.8%
- Latência média: 120ms (p95: 380ms)
- Total requests: 15.234
- Novos signups: 12
- Active users: 8

## Incidentes

- 1 alerta crítico (high latency) - resolvido em 15 min
- Causa: Query N+1 em listagem de assinaturas
- Fix: Implementado eager loading

## Feedback Clientes

- 5 respostas NPS: Média 8.2/10
- Feature mais pedida: Relatórios PDF
- 1 bug reportado: Filtro de data não funciona (corrigido)

## Ações Próxima Semana

- Implementar relatórios PDF
- Otimizar listagem de assinaturas
- Adicionar tutorial in-app
```

- [ ] Relatório publicado em `docs/POST_LAUNCH_REPORTS/`
- [ ] Apresentado ao time em reunião

---

## 📈 Métricas de Sucesso

### Fase 7 completa quando:

- [ ] ✅ Todas as 6 tasks concluídas (100%)
- [ ] ✅ Fase 6 100% completa (LGPD + Backup)
- [ ] ✅ Checklist pré-lançamento 100% validado
- [ ] ✅ Deploy produção realizado com sucesso
- [ ] ✅ Smoke tests passando em produção
- [ ] ✅ Comunicação enviada aos stakeholders
- [ ] ✅ Primeiros clientes onboarded
- [ ] ✅ Monitoramento 24/7 ativo
- [ ] ✅ Zero incidentes críticos nas primeiras 72h

---

## 🎯 Deliverables da Fase 7

| #   | Deliverable                       | Status      | Validação          |
| --- | --------------------------------- | ----------- | ------------------ |
| 1   | T-LGPD-001 + T-OPS-005 completos  | ⏳ Pendente | Testes passando    |
| 2   | Checklist pré-lançamento validado | ⏳ Pendente | 100% checado       |
| 3   | Deploy produção realizado         | ⏳ Pendente | Smoke tests OK     |
| 4   | Comunicação enviada               | ⏳ Pendente | Emails enviados    |
| 5   | Onboarding funcionando            | ⏳ Pendente | Primeiros clientes |
| 6   | Monitoramento ativo               | ⏳ Pendente | Relatório semanal  |

---

## 🚨 Plano de Contingência

### Se Deploy Falhar

**Sintomas:** Build error, systemd service não inicia, health check 500

**Ações:**

1. Rollback imediato (git checkout commit anterior)
2. Rebuild + restart service
3. Validar health check
4. Investigar logs
5. Fix em branch separada
6. Re-deploy após validação em staging

### Se Database Corrompido

**Sintomas:** Queries falhando, migrations quebradas, dados inconsistentes

**Ações:**

1. Ativar Neon PITR (restore to 1h ago)
2. Ou restore de S3 backup
3. Validar integridade com `validate_schema.sh`
4. Re-deploy backend se necessário
5. Post-mortem para identificar causa

### Se Feedback Negativo

**Sintomas:** NPS <5, reclamações recorrentes, churn alto

**Ações:**

1. Call emergencial com clientes insatisfeitos
2. Priorizar bugs críticos reportados
3. Sprint de correções (1 semana)
4. Re-onboarding com melhorias
5. Compensação se necessário (desconto, extensão trial)

---

## 📅 Timeline Detalhado

```
Dia 1 (20/12):
  - Completar T-LGPD-001 (8h)
  - Iniciar T-OPS-005 (3h)

Dia 2 (21/12):
  - Completar T-OPS-005 (3h)
  - Executar checklist pré-lançamento (3h)

Dia 3 (22/12):
  - Corrigir issues encontradas no checklist (4h)
  - Re-executar checklist (1h)

Dia 4 (23/12):
  - Deploy produção (2h)
  - Smoke tests (1h)
  - Comunicação stakeholders (1h)

Dia 5 (26/12):
  - 🎉 GO-LIVE OFICIAL
  - Monitoramento intensivo (8h)
  - Primeiros onboardings

Semana seguinte (27/12 - 02/01):
  - Monitoramento diário (1h/dia)
  - Suporte a clientes
  - Coleta de feedback
  - Correções críticas
```

---

## 🎉 Comunicado de Lançamento (Template)

```
🚀 ESTAMOS NO AR!

Barber Analytics Pro v2.0 acaba de ser lançado! 🎊

Depois de 12 semanas de desenvolvimento intenso, temos orgulho de apresentar
a nova plataforma de gestão para barbearias mais completa do mercado.

✨ Novidades:
- Dashboard interativo com KPIs em tempo real
- Fluxo de caixa automatizado
- Gestão de assinaturas (Clube do Trato)
- Sistema multi-usuário com permissões
- Relatórios avançados
- 100% seguro e LGPD compliant

🔗 Acesse: https://barberpro.dev
📧 Suporte: suporte@barberpro.dev
📱 WhatsApp: (11) 9xxxx-xxxx

Experimente GRÁTIS por 14 dias! 🎁

#BarbeariaDigital #Gestão #SaaS #Empreendedorismo
```

---

## 🏁 Critérios de Sucesso Final

### MVP 2.0 Lançado quando:

- ✅ Deploy em produção estável por 72h
- ✅ ≥5 clientes ativos usando sistema
- ✅ NPS médio ≥7/10
- ✅ Uptime ≥99%
- ✅ Zero bugs críticos em aberto
- ✅ Equipe de suporte treinada e atuando

### Celebração 🎉

- [ ] Retrospectiva do projeto (2h)
- [ ] Apresentação de métricas ao time
- [ ] Reconhecimento individual das conquistas
- [ ] Planejamento Fase 8 (Evolução)

---

**Última Atualização:** 20/11/2025 09:30
**Status:** ⏳ Aguardando conclusão da Fase 6 (LGPD + Backup)
**Data Go-Live Prevista:** 26/12/2025
**Responsável:** Tech Lead + DevOps Lead
**Stakeholders:** CEO, Product Manager, Marketing, Suporte

---

**🎯 Próximo:** FASE 8 — Monitoramento & Estabilização (Semanas 1-4 pós-lançamento)
