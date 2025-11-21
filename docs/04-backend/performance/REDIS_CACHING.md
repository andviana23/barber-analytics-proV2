> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 📦 Redis Caching — Barber Analytics Pro

## 📋 Visão Geral

Sistema de cache distribuído usando Redis para melhorar performance de endpoints críticos.

**Objetivo:** Cache hit rate > 70% para recursos frequentes

---

## 🚀 Setup

### 1. Iniciar Redis (instalação local ou serviço gerenciado)

- Opção A (local): instalar `redis-server` (ex.: `sudo apt install redis-server`), editar `redis.conf` para `requirepass barber123`, `maxmemory 256mb`, `maxmemory-policy allkeys-lru`, e subir com `sudo systemctl restart redis-server`.
- Opção B (serviço gerenciado): provisionar Redis 7+ com senha e política LRU equivalente, expor `REDIS_URL`/`REDIS_PASSWORD`.

### 2. Variáveis de Ambiente

Adicionar no `.env` ou variáveis de ambiente:

```env
REDIS_URL=localhost:6379
REDIS_PASSWORD=barber123
REDIS_DB=0
CACHE_ENABLED=true
```

### 3. Verificar Conexão

```bash
redis-cli -a barber123 ping
# Resposta: PONG
```

---

## 🗂️ Estrutura de Chaves

### Convenções

| Chave | Pattern | TTL | Exemplo |
|-------|---------|-----|---------|
| Dashboard KPIs | `dashboard:kpis:{tenant_id}` | 1 hora | `dashboard:kpis:123e4567-e89b-12d3-a456-426614174000` |
| Planos Assinatura | `subscription:plans:{tenant_id}` | 24 horas | `subscription:plans:123e4567-e89b-12d3-a456-426614174000` |
| Categorias | `categorias:{tenant_id}` | 7 dias | `categorias:123e4567-e89b-12d3-a456-426614174000` |
| Usuário | `user:{user_id}` | 15 minutos | `user:987f6543-e21b-12d3-a456-426614174000` |
| Tenant | `tenant:{tenant_id}` | 1 hora | `tenant:123e4567-e89b-12d3-a456-426614174000` |

### Namespaces

Todos os namespaces são automaticamente extraídos da primeira parte da chave (antes de `:`) para métricas Prometheus.

---

## 📊 Métricas Prometheus

### Métricas Coletadas

```prometheus
# Cache hits por namespace
cache_hits_total{namespace="dashboard"} 1250

# Cache misses por namespace
cache_misses_total{namespace="dashboard"} 150

# Erros de cache por namespace e operação
cache_errors_total{namespace="dashboard",operation="get"} 2

# Latência de operações de cache
cache_operation_duration_seconds{namespace="dashboard",operation="get",quantile="0.95"} 0.002
```

### Queries Úteis

```prometheus
# Hit rate por namespace
sum(rate(cache_hits_total[5m])) by (namespace) /
(sum(rate(cache_hits_total[5m])) by (namespace) + sum(rate(cache_misses_total[5m])) by (namespace))

# Top namespaces com mais hits
topk(5, sum(rate(cache_hits_total[5m])) by (namespace))

# Taxa de erro de cache
sum(rate(cache_errors_total[5m])) / sum(rate(cache_hits_total[5m]) + rate(cache_misses_total[5m]))
```

---

## 🔄 Invalidação de Cache

### Manual (via código)

```go
import "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/cache"

// Criar invalidador
invalidator := cache.NewInvalidator(cacheClient)

// Invalidar dashboard de um tenant
err := invalidator.InvalidateDashboard(ctx, tenantID)

// Invalidar planos de assinatura
err := invalidator.InvalidateSubscriptionPlans(ctx, tenantID)

// Invalidar categorias
err := invalidator.InvalidateCategorias(ctx, tenantID)

// Invalidar tudo de um tenant
err := invalidator.InvalidateAll(ctx, tenantID)
```

### Automática (após mutations)

**Recomendação:** Invalidar cache nos handlers de CREATE/UPDATE/DELETE:

```go
// Exemplo: handler de criar receita
func (h *ReceitaHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
    // ... criar receita ...

    // Invalidar cache do dashboard
    h.invalidator.InvalidateDashboard(r.Context(), tenantID)
}
```

**Tabela de Invalidação:**

| Operação | Cache Invalidado |
|----------|------------------|
| CREATE/UPDATE/DELETE Receita | `dashboard:kpis:{tenant_id}` |
| CREATE/UPDATE/DELETE Despesa | `dashboard:kpis:{tenant_id}` |
| CREATE/UPDATE/DELETE Assinatura | `dashboard:kpis:{tenant_id}` + `subscription:plans:{tenant_id}` |
| CREATE/UPDATE/DELETE Plano | `subscription:plans:{tenant_id}` |
| CREATE/UPDATE/DELETE Categoria | `categorias:{tenant_id}` |

---

## 🧪 Testes

### Testar Cache Hit/Miss

```bash
# 1. Fazer primeira requisição (cache miss)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/dashboard

# 2. Verificar no Redis
redis-cli -a barber123 KEYS dashboard:*
redis-cli -a barber123 TTL dashboard:kpis:123e4567-e89b-12d3-a456-426614174000

# 3. Fazer segunda requisição (cache hit)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/dashboard

# 4. Verificar métricas Prometheus
curl http://localhost:8080/metrics | grep cache_hits_total
```

### Invalidar Manualmente

```bash
# Via redis-cli
redis-cli -a barber123 DEL dashboard:kpis:123e4567-e89b-12d3-a456-426614174000
# Limpar tudo (cuidado em produção!)
redis-cli -a barber123 FLUSHDB
```

---

## 📈 Performance Esperada

### Antes do Cache (baseline)

- Dashboard KPIs: ~850ms (query pesada)
- Lista de Planos: ~120ms
- Lista de Categorias: ~80ms

### Depois do Cache (hit)

- Dashboard KPIs: ~5-10ms (170x mais rápido)
- Lista de Planos: ~2-5ms (40x mais rápido)
- Lista de Categorias: ~2-5ms (20x mais rápido)

### Meta

- **Cache Hit Rate:** > 70%
- **Latência p95 (cache hit):** < 10ms
- **Latência p95 (cache miss):** < 500ms (original + overhead de cache)

---

## 🛠️ Troubleshooting

### Redis não conecta

```bash
# Testar conexão
redis-cli -a barber123 ping

# Ver serviço
sudo systemctl status redis-server

# Logs (Ubuntu)
sudo journalctl -u redis-server -f
```

### Cache não invalida

1. Verificar se `CACHE_ENABLED=true`
2. Conferir logs do backend para erros de cache
3. Verificar se invalidador está sendo chamado após mutations
4. Checar se tenant_id está correto

### Hit rate baixo (< 70%)

1. Aumentar TTLs (se dados mudam pouco)
2. Verificar se endpoints cacheados estão sendo usados frequentemente
3. Analisar padrão de acesso (Grafana)
4. Considerar cache adicional (ex: receitas recentes)

### Memory usage alto

```bash
# Ver uso de memória
redis-cli -a barber123 INFO memory

# Ver chaves por namespace
redis-cli -a barber123 SCAN 0 MATCH dashboard:* COUNT 100
redis-cli -a barber123 SCAN 0 MATCH subscription:* COUNT 100

# Ajustar maxmemory se necessário
redis-cli -a barber123 CONFIG SET maxmemory 256mb
redis-cli -a barber123 CONFIG SET maxmemory-policy allkeys-lru
```

---

## 🔐 Segurança

### Produção

**Recomendações:**

1. **Password forte:** Usar variável de ambiente, não hardcode
2. **TLS/SSL:** Configurar Redis com SSL em produção
3. **Network isolation:** Redis em rede privada, não expor porta 6379 publicamente
4. **ACL:** Configurar ACLs do Redis 6+ para limitar comandos perigosos (FLUSHDB, CONFIG)

```bash
# Exemplo de ACL
# redis.conf
user default on >senha_forte ~* &* +@all -@dangerous
```

### Monitoramento

- Alertas: Memory usage > 80%
- Alertas: Connection errors > 10 em 5 min
- Alertas: Hit rate < 50% (possível problema de invalidação)

---

## 📚 Referências

- [Redis Documentation](https://redis.io/documentation)
- [go-redis Client](https://redis.uptrace.dev/)
- [Redis Best Practices](https://redis.io/docs/management/optimization/)
- [Eviction Policies](https://redis.io/docs/reference/eviction/)

---

**Última Atualização:** 15/11/2025
**Responsável:** Backend Team
**Status:** ✅ Implementado
