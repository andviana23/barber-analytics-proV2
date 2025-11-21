> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📚 Runbook - Alertas Barber Analytics Pro

## 🎯 Visão Geral

Este documento contém procedimentos operacionais (runbooks) para responder a cada tipo de alerta do sistema Barber Analytics Pro.

**Princípios:**
- ⏱️ **Tempo de Resposta:** Critical (< 15 min), Warning (< 1h)
- 👥 **Escalonamento:** DevOps → Backend Team → CTO
- 📊 **Diagnóstico:** Usar Grafana dashboards + logs + métricas Prometheus

---

## 🚨 CRITICAL ALERTS

### 1. ServiceCompletelyDown

**Alerta:** Backend completamente fora do ar (up == 0)

**Impacto:** 🔴 TOTAL - Sistema inacessível para todos os usuários

**Tempo de Resposta:** < 5 minutos

#### Checklist de Diagnóstico

```bash
# 1. Verificar se o processo está rodando
pm2 list
pm2 logs barber-api --lines 50

# 2. Verificar se a porta está acessível
curl -I http://localhost:8080/health

# 3. Verificar logs do sistema
journalctl -u barber-api -n 100 --no-pager

# 4. Verificar recursos do servidor
free -h
df -h
top -bn1 | head -20
```

#### Ações Corretivas

**Cenário A: Processo PM2 parado**
```bash
cd /path/to/backend
pm2 restart barber-api
pm2 save
```

**Cenário B: Porta em uso / conflito**
```bash
# Identificar processo usando porta 8080
lsof -i :8080
# Matar processo conflitante
kill -9 <PID>
# Reiniciar aplicação
pm2 restart barber-api
```

**Cenário C: Falta de memória/disco**
```bash
# Limpar logs antigos
journalctl --vacuum-time=7d
# Limpar cache do sistema
sync; echo 3 > /proc/sys/vm/drop_caches
# Reiniciar aplicação
pm2 restart barber-api
```

**Cenário D: Erro fatal no código (panic)**
```bash
# Fazer rollback para versão anterior
cd /path/to/backend
git log --oneline -5
git checkout <commit-anterior>
make build
pm2 restart barber-api
```

#### Escalonamento
- **5 min sem resolução:** Chamar Backend Lead
- **15 min sem resolução:** Chamar CTO + considerar comunicado aos clientes

---

### 2. ServiceDowntime (Uptime < 99.5% em 24h)

**Alerta:** Disponibilidade abaixo de 99.5% nas últimas 24 horas

**Impacto:** 🟠 ALTO - SLA violado, usuários afetados intermitentemente

**Tempo de Resposta:** < 15 minutos

#### Checklist de Diagnóstico

```bash
# 1. Verificar uptime do Prometheus
curl http://localhost:9090/api/v1/query?query='avg_over_time(up{job="barber-backend"}[24h])'

# 2. Verificar histórico de restarts
pm2 describe barber-api

# 3. Analisar logs de erros nas últimas 24h
journalctl -u barber-api --since "24 hours ago" | grep -i "error\|fatal\|panic"

# 4. Verificar dashboard Grafana Overview → Uptime
```

#### Ações Corretivas

**Investigar causa raiz:**
- Restarts frequentes → Memory leak ou panic recorrente
- Timeouts → Queries lentas ou serviços externos
- Erros 5xx → Bugs em produção

**Próximos passos:**
1. Analisar padrão de downtime (horário específico?)
2. Correlacionar com deploys recentes
3. Verificar alertas relacionados (latência, DB, crons)
4. Implementar correção e monitorar por 48h

---

### 3. HighErrorRate (Error rate 5xx > 1% em 5 min)

**Alerta:** Taxa de erros 5xx acima de 1% por mais de 5 minutos

**Impacto:** 🔴 ALTO - Múltiplos usuários afetados, possível bug crítico

**Tempo de Resposta:** < 10 minutos

#### Checklist de Diagnóstico

```bash
# 1. Identificar endpoints com mais erros
curl 'http://localhost:9090/api/v1/query?query=topk(5, sum(rate(http_requests_total{status=~"5.."}[5m])) by (endpoint))'

# 2. Ver logs de erros recentes
pm2 logs barber-api --lines 100 | grep "ERROR\|500\|502\|503"

# 3. Verificar se banco está acessível
psql -h <NEON_HOST> -U <USER> -d <DB> -c "SELECT 1;"

# 4. Verificar serviços externos (Asaas, etc)
curl -I https://sandbox.asaas.com/api/v3/customers
```

#### Ações Corretivas

**Cenário A: Banco de dados inacessível**
```bash
# Verificar credenciais e conectividade
ping <NEON_HOST>
# Verificar se pool de conexões não esgotou
# (ver dashboard Database → Connections)
```

**Cenário B: Deploy recente com bugs**
```bash
# Rollback imediato
cd /path/to/backend
git log --oneline -5
git checkout <commit-anterior-estavel>
make build
pm2 restart barber-api
```

**Cenário C: Serviço externo fora do ar**
- Ativar circuit breaker (se implementado)
- Retornar respostas cached ou degradadas
- Comunicar usuários sobre funcionalidade limitada

#### Escalonamento
- **10 min sem resolução:** Chamar Backend Lead
- **30 min sem resolução:** Considerar comunicado aos clientes

---

### 4. HighLatencyP95 (Latência p95 > 500ms em 5 min)

**Alerta:** 95% das requisições levam mais de 500ms

**Impacto:** 🟠 MÉDIO - Experiência degradada, usuários percebem lentidão

**Tempo de Resposta:** < 30 minutos

#### Checklist de Diagnóstico

```bash
# 1. Identificar endpoints mais lentos
curl 'http://localhost:9090/api/v1/query?query=topk(5, histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (endpoint, le)))'

# 2. Verificar dashboard Grafana Backend → Latency Heatmap

# 3. Verificar queries lentas no banco
# (ver dashboard Database → Slow Queries)

# 4. Verificar recursos do servidor
top -bn1 | head -20
iostat -x 1 5
```

#### Ações Corretivas

**Cenário A: Queries lentas (ver DB dashboard)**
```sql
-- Conectar ao banco e identificar queries ativas
SELECT pid, usename, state, query_start, query
FROM pg_stat_activity
WHERE state = 'active' AND query_start < now() - interval '5 seconds';

-- Matar queries problemáticas
SELECT pg_terminate_backend(<pid>);
```

**Cenário B: CPU/Memória alta**
```bash
# Identificar processos pesados
ps aux --sort=-%cpu | head -10
ps aux --sort=-%mem | head -10

# Considerar escalar recursos (vertical scaling)
# ou adicionar réplicas (horizontal scaling)
```

**Cenário C: Tráfego anormalmente alto**
- Verificar se há spike de requisições (dashboard Overview)
- Verificar rate limiting (deve estar ativo)
- Considerar ativar cache Redis (T-PERF-002)

---

### 5. CronNotExecuted (Cron não executou em 25h)

**Alerta:** Job agendado não executou com sucesso nas últimas 25 horas

**Impacto:** 🟠 VARIÁVEL - Depende do job (backups críticos vs. relatórios)

**Tempo de Resposta:** < 30 minutos

#### Checklist de Diagnóstico

```bash
# 1. Verificar qual job falhou
# (ver alerta ou dashboard Grafana Crons)

# 2. Verificar logs do scheduler
pm2 logs barber-api | grep -i "cron\|scheduler\|job"

# 3. Executar job manualmente para testar
# (depende da implementação do scheduler)

# 4. Verificar dependências do job
# - Banco de dados acessível?
# - Serviços externos disponíveis?
```

#### Ações Corretivas

**Cenário A: Scheduler travado**
```bash
# Reiniciar aplicação
pm2 restart barber-api

# Verificar se scheduler reiniciou
pm2 logs barber-api --lines 50 | grep -i "scheduler"
```

**Cenário B: Job falhando silenciosamente**
- Revisar código do job
- Adicionar logs detalhados
- Implementar retry mechanism
- Configurar timeout adequado

**Cenário C: Dependência externa indisponível**
- Verificar status da API externa (Asaas, email, etc)
- Implementar fallback ou queue para retry posterior

#### Jobs Críticos (Prioridade Alta)
- **Backups diários:** Executar manualmente imediatamente
- **Sincronização Asaas:** Verificar se há transações perdidas
- **Relatórios financeiros:** Regenerar e enviar manualmente

---

### 6. DatabaseConnectionsExhausted (Pool esgotado)

**Alerta:** Requisições aguardando conexões disponíveis (waiting > 5)

**Impacto:** 🔴 CRÍTICO - Sistema praticamente inacessível, timeouts generalizados

**Tempo de Resposta:** < 5 minutos

#### Checklist de Diagnóstico

```bash
# 1. Verificar pool stats
curl 'http://localhost:9090/api/v1/query?query=db_connections_open'
curl 'http://localhost:9090/api/v1/query?query=db_connections_in_use'
curl 'http://localhost:9090/api/v1/query?query=db_connections_waiting'

# 2. Identificar queries travadas
psql -h <NEON_HOST> -U <USER> -d <DB> -c "
SELECT pid, usename, state, query_start, query
FROM pg_stat_activity
WHERE state = 'active'
ORDER BY query_start
LIMIT 20;
"

# 3. Verificar locks
psql -h <NEON_HOST> -U <USER> -d <DB> -c "
SELECT l.pid, l.mode, l.granted, a.query
FROM pg_locks l
JOIN pg_stat_activity a ON l.pid = a.pid
WHERE NOT l.granted;
"
```

#### Ações Corretivas

**Ação Imediata: Matar queries travadas**
```sql
-- Matar queries ativas há mais de 30 segundos
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'active'
  AND query_start < now() - interval '30 seconds'
  AND query NOT LIKE '%pg_stat_activity%';
```

**Correção de Médio Prazo:**
1. Aumentar pool size no backend (arquivo de config)
2. Implementar timeouts mais agressivos em queries
3. Otimizar queries lentas (T-PERF-001)
4. Implementar connection pooling com PgBouncer (opcional)

---

## ⚠️ WARNING ALERTS

### 7. High4xxErrorRate (Error rate 4xx > 5% em 10 min)

**Alerta:** Taxa de erros 4xx acima de 5% por mais de 10 minutos

**Impacto:** 🟡 MÉDIO - Possível problema de validação ou autenticação

**Tempo de Resposta:** < 1 hora

#### Checklist de Diagnóstico

```bash
# 1. Identificar quais status codes 4xx
curl 'http://localhost:9090/api/v1/query?query=sum(rate(http_requests_total{status=~"4.."}[10m])) by (status)'

# 2. Identificar endpoints afetados
curl 'http://localhost:9090/api/v1/query?query=topk(5, sum(rate(http_requests_total{status=~"4.."}[10m])) by (endpoint))'

# 3. Analisar logs
pm2 logs barber-api --lines 200 | grep -E "401|403|404|422"
```

#### Ações Corretivas

**401 Unauthorized:** Problema de autenticação
- Verificar se JWT secret mudou
- Verificar expiração de tokens
- Verificar integração com sistema de auth

**403 Forbidden:** Problema de autorização (RBAC)
- Verificar se roles/permissions mudaram recentemente
- Verificar middleware de autorização

**404 Not Found:** Rota não encontrada
- Verificar se houve mudança nas rotas (breaking change)
- Comunicar frontend se necessário

**422 Unprocessable Entity:** Validação falhando
- Verificar se regras de validação mudaram
- Analisar payloads sendo enviados pelos clients

---

### 8. HighMemoryUsage / HighGoroutineCount

**Alerta:** Uso de memória ou goroutines acima do threshold

**Impacto:** 🟡 MÉDIO - Possível memory/goroutine leak

**Tempo de Resposta:** < 2 horas (monitorar evolução)

#### Checklist de Diagnóstico

```bash
# 1. Coletar heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# 2. Coletar goroutine profile
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof

# 3. Analisar com pprof
go tool pprof -http=:8081 heap.prof
go tool pprof -http=:8082 goroutine.prof
```

#### Ações Corretivas

**Se houver leak confirmado:**
1. Identificar código problemático via pprof
2. Criar issue no GitHub com evidências
3. Deploy de hotfix assim que disponível
4. Reiniciar aplicação periodicamente (workaround temporário)

**Se for crescimento natural:**
- Considerar aumentar memória do servidor
- Implementar memory limits no Go (GOMEMLIMIT)
- Otimizar estruturas de dados em memória

---

## 📞 Contatos de Escalonamento

| Severidade | Tempo | Contato | Canal |
|------------|-------|---------|-------|
| Critical | Imediato | DevOps On-Call | PagerDuty + Telefone |
| Critical | +15 min | Backend Lead | Slack DM + Telefone |
| Critical | +30 min | CTO | Telefone |
| Warning | +1h | Backend Team | Slack #backend |
| Warning | +4h | DevOps Lead | Slack DM |

---

## 🔗 Links Úteis

- **Grafana:** http://grafana.barberanalytics.com.br
- **Prometheus:** http://prometheus.barberanalytics.com.br
- **Alertmanager:** http://alertmanager.barberanalytics.com.br
- **Logs:** Acesso SSH ao servidor de produção
- **Documentação Técnica:** https://docs.barberanalytics.com.br

---

**Versão:** 1.0
**Última Atualização:** 15/11/2025
**Mantenedor:** DevOps Team
