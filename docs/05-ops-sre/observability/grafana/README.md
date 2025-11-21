> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Grafana Dashboards - Barber Analytics Pro

## 📊 Visão Geral

Este diretório contém 4 dashboards Grafana profissionais para monitoramento completo do Barber Analytics Pro, baseados nas métricas coletadas pelo Prometheus.

---

## 📁 Arquivos

### 1. `datasource.yaml`
Configuração da fonte de dados Prometheus para o Grafana.

**Configuração:**
- URL: `http://localhost:9090`
- Método HTTP: POST
- Intervalo de scrape: 15s

### 2. `dashboard-overview.json`
Dashboard principal com visão geral do sistema.

**Painéis:**
- ✅ **Uptime (24h)** - Disponibilidade do sistema
- ✅ **Total Requests (24h)** - Volume total de requisições
- ✅ **Error Rate (5m)** - Taxa de erro em tempo real
- ✅ **Active Tenants** - Total de tenants ativos
- ✅ **Requests per Second** - Throughput do sistema
- ✅ **Error Rate Over Time** - Erros 4xx/5xx ao longo do tempo
- ✅ **Top 10 Endpoints** - Endpoints mais acessados

**Alertas Configurados:**
- Error Rate > 1% → Estado crítico (vermelho)
- Uptime < 99.5% → Estado de atenção (amarelo)

---

### 3. `dashboard-backend.json`
Dashboard focado em performance do backend Go.

**Painéis:**
- ✅ **Request Latency (p50/p95/p99)** - Distribuição de latência
- ✅ **Throughput** - Req/s por status code (2xx, 4xx, 5xx)
- ✅ **In-Flight Requests** - Requisições concorrentes
- ✅ **Response Size Distribution** - Tamanho das respostas
- ✅ **Memory Usage** - Memória alocada, heap, stack
- ✅ **Goroutines** - Goroutines ativas
- ✅ **GC Pause Duration** - Tempo de pausa do Garbage Collector
- ✅ **Latency Heatmap** - Heatmap de latência por endpoint

**Alertas Configurados:**
- Latency p95 > 500ms → Alerta de alta latência

**Métricas Go Runtime:**
- `go_memstats_alloc_bytes` - Memória alocada
- `go_memstats_heap_inuse_bytes` - Heap em uso
- `go_goroutines` - Total de goroutines
- `go_gc_duration_seconds` - Duração do GC

---

### 4. `dashboard-crons.json`
Dashboard para monitoramento de jobs agendados (cron jobs).

**Painéis:**
- ✅ **Last Execution Time** - Última execução de cada job (tabela)
- ✅ **Cron Execution Status** - Status de sucesso/falha
- ✅ **Execution Duration** - Duração média (p50/p95) por job
- ✅ **Cron Executions Over Time** - Execuções ao longo do tempo
- ✅ **Failed Executions (24h)** - Tabela de jobs com falhas
- ✅ **Duration Heatmap** - Distribuição de duração
- ✅ **Jobs Not Executed (ALERT)** - Jobs que não executaram em 25h

**Alertas Críticos:**
- Job não executado em 25 horas → Alerta vermelho
- Detecção automática de jobs silenciosos

**Use Cases:**
- Identificar crons travados
- Monitorar tempo de execução de backups
- Detectar falhas em processamento batch

---

### 5. `dashboard-database.json`
Dashboard para monitoramento do PostgreSQL via métricas do backend.

**Painéis:**
- ✅ **Database Connections** - Open, In Use, Idle, Waiting
- ✅ **Connection Pool Stats** - Estatísticas em tempo real
- ✅ **Query Count by Operation** - SELECT, INSERT, UPDATE, DELETE
- ✅ **Query Count by Table** - Top 10 tabelas mais acessadas
- ✅ **Query Duration (p50/p95/p99)** - Distribuição de latência
- ✅ **Slow Queries (>1s)** - Tabela de queries lentas
- ✅ **Query Duration by Operation** - Latência por tipo de operação
- ✅ **Query Duration Heatmap** - Visualização de distribuição

**Alertas Configurados:**
- Connections > 20 → Alerta de pool esgotado
- Query p99 > 1s → Alerta de queries lentas

**Detecção de Problemas:**
- N+1 queries (alto volume em curto período)
- Queries lentas (>1s)
- Pool de conexões esgotado
- Operações bloqueantes

---

## 🚀 Como Usar

### 1. Instalar Grafana

**Via Docker:**
```bash
docker run -d \
  --name=grafana \
  -p 3000:3000 \
  -v grafana-storage:/var/lib/grafana \
  grafana/grafana:latest
```

**Via Helm (Kubernetes):**
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm install grafana grafana/grafana
```

### 2. Configurar Datasource

**Opção A - Via arquivo:**
```bash
# Copiar datasource.yaml para Grafana
cp datasource.yaml /etc/grafana/provisioning/datasources/
```

**Opção B - Via UI:**
1. Acessar Grafana: http://localhost:3000
2. Login padrão: admin/admin
3. Configuration → Data Sources → Add data source
4. Selecionar Prometheus
5. URL: `http://localhost:9090`
6. Save & Test

### 3. Importar Dashboards

**Opção A - Via arquivo (provisioning):**
```bash
# Copiar dashboards para Grafana
cp dashboard-*.json /etc/grafana/provisioning/dashboards/
```

**Opção B - Via UI:**
1. Dashboards → Import
2. Upload JSON file ou copiar/colar conteúdo
3. Selecionar datasource "Prometheus"
4. Import

### 4. Validar

Verificar se os 4 dashboards estão visíveis:
- ✅ Barber Analytics - Overview
- ✅ Barber Analytics - Backend Performance
- ✅ Barber Analytics - Cron Jobs
- ✅ Barber Analytics - Database

---

## 📈 Métricas Utilizadas

### HTTP Metrics (do PrometheusMiddleware)
- `http_requests_total{method, endpoint, status}`
- `http_request_duration_seconds_bucket{method, endpoint, status, le}`
- `http_requests_in_flight`
- `http_response_size_bytes_bucket{method, endpoint, le}`
- `http_errors_total{method, endpoint, status}`

### Database Metrics
- `db_connections_open`
- `db_connections_idle`
- `db_connections_in_use`
- `db_connections_waiting`
- `db_queries_total{operation, table}`
- `db_queries_duration_seconds_bucket{operation, table, le}`

### Cron Metrics
- `cron_executions_total{job_name, status}`
- `cron_execution_duration_seconds_bucket{job_name, le}`
- `cron_last_success_timestamp{job_name}`

### Business Metrics
- `barber_tenants_total`
- `barber_users_total`
- `barber_receitas_created_total{tenant_id}`
- `barber_despesas_created_total{tenant_id}`

### Go Runtime Metrics (automáticas)
- `go_memstats_alloc_bytes`
- `go_memstats_heap_inuse_bytes`
- `go_goroutines`
- `go_gc_duration_seconds`
- `up{job}`

---

## 🎯 Alertas Recomendados (T-OPS-004)

Os dashboards já incluem alertas básicos. Para configuração completa de notificações:

### 1. Configurar Notification Channels
```bash
# Slack
Configuration → Notification Channels → Add Channel
Type: Slack
Webhook URL: <seu-webhook-slack>
```

### 2. Alertas Configurados nos Dashboards

**Dashboard: Backend**
- ⚠️ Latency p95 > 500ms (5 min window)

**Dashboard: Database**
- ⚠️ Connections > 20
- ⚠️ Query p99 > 1s

**Dashboard: Crons**
- ⚠️ Job não executou em 25h (detecção automática)

### 3. Alertas Adicionais (via Prometheus rules)
Ver `T-OPS-004` para regras completas de alerting.

---

## 🔍 Troubleshooting

### Dashboard não mostra dados?
1. Verificar se Prometheus está rodando: `curl http://localhost:9090/metrics`
2. Verificar se backend expõe /metrics: `curl http://localhost:8080/metrics`
3. Verificar scrape config em `prometheus.yml`
4. Verificar logs do Prometheus: `docker logs prometheus`

### Queries retornam "No data"?
1. Verificar nome das métricas: `http_requests_total` vs `http_request_total`
2. Verificar labels: `{job="barber-backend"}`
3. Aguardar 1-2 minutos para primeira coleta
4. Verificar time range do dashboard (last 1h, last 24h, etc)

### Alertas não disparam?
1. Verificar se Alertmanager está configurado
2. Verificar notification channels
3. Testar manualmente: Dashboard → Edit → Alert tab → Test Rule

---

## 📚 Referências

- **Prometheus Querying:** https://prometheus.io/docs/prometheus/latest/querying/basics/
- **Grafana Dashboards:** https://grafana.com/docs/grafana/latest/dashboards/
- **PromQL Functions:** https://prometheus.io/docs/prometheus/latest/querying/functions/

---

## ✅ Checklist de Validação

- [ ] Grafana instalado e acessível (http://localhost:3000)
- [ ] Datasource Prometheus configurado
- [ ] Dashboard Overview importado e funcional
- [ ] Dashboard Backend importado e funcional
- [ ] Dashboard Crons importado e funcional
- [ ] Dashboard Database importado e funcional
- [ ] Todas as queries retornam dados
- [ ] Alertas testados manualmente
- [ ] Notification channels configurados (Slack/Email)
- [ ] Documentação revisada e atualizada

---

**Criado:** 15/11/2025
**Autor:** Andrey Viana
**Versão:** 1.0
**Projeto:** Barber Analytics Pro v2.0
