> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# ⏰ Fluxo de Crons

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Design Finalizado

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Crons Diários](#crons-diários)
3. [Crons Semanais](#crons-semanais)
4. [Crons Mensais](#crons-mensais)
5. [Monitoramento](#monitoramento)
6. [Error Handling](#error-handling)

---

## 🎯 Visão Geral

Sistema de crons para automações assíncronas do Barber Analytics Pro. Utiliza `robfig/cron/v3` em Go + systemd no VPS.

```
┌─────────────────────────────────────┐
│     Scheduler (robfig/cron)         │
├─────────────────────────────────────┤
│                                     │
│  ├─ Job 1: Sincronizar Asaas (02h) │
│  ├─ Job 2: Snapshot Financeiro     │
│  ├─ Job 3: Processar Repasse       │
│  ├─ Job 4: Alertas                 │
│  └─ Job 5: Limpeza de Logs         │
│                                     │
└─────────────────────────────────────┘
```

---

## 📅 Crons Diários

### 1️⃣ Sincronizar Faturas Asaas (02:00)

**Schedule:** `0 2 * * *`  
**Timeout:** 15 minutos  
**Retry:** 3x com backoff exponencial

**Objetivo:** Buscar novas faturas do Asaas e sincronizar com o banco local.

```go
type SyncAsaasInvoicesJob struct {
    asaasClient  *external.AsaasClient
    invoiceRepo  domain.AssinaturaInvoiceRepository
    assinatureRepo domain.AssinaturaRepository
    logger       *zap.Logger
}

func (j *SyncAsaasInvoicesJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
    defer cancel()
    
    j.logger.Info(\"Starting Asaas sync job\")
    
    // Para cada tenant ativo
    tenants, _ := j.tenantRepo.FindActive(ctx)
    
    for _, tenant := range tenants {
        j.logger.Info(\"Syncing tenant\", zap.String(\"tenant_id\", tenant.ID))
        
        // 1. Buscar todas assinaturas ativas
        subs, _ := j.assinatureRepo.FindByTenant(ctx, tenant.ID, \"ATIVA\")
        
        // 2. Para cada assinatura, sincronizar faturas
        for _, sub := range subs {
            invoices, err := j.asaasClient.ListInvoices(&external.ListInvoicesReq{
                SubscriptionID: sub.AsaasSubscriptionID,
            })
            
            if err != nil {
                j.logger.Error(\"Failed to sync invoices\",
                    zap.String(\"tenant_id\", tenant.ID),
                    zap.String(\"sub_id\", sub.ID),
                    zap.Error(err))
                continue
            }
            
            // 3. Persistir invoices
            for _, asaasInvoice := range invoices {
                existing, _ := j.invoiceRepo.FindByAsaasID(ctx, tenant.ID, asaasInvoice.ID)
                
                if existing == nil {
                    // Nova fatura
                    invoice := &domain.AssinaturaInvoice{
                        ID:             uuid.NewString(),
                        TenantID:       tenant.ID,
                        AssinaturaID:   sub.ID,
                        AsaasInvoiceID: asaasInvoice.ID,
                        Valor:          decimal.NewFromFloat(asaasInvoice.Value),
                        Status:         mapAsaasStatus(asaasInvoice.Status),
                        DataVencimento: asaasInvoice.DueDate,
                        CriadoEm:       time.Now(),
                    }
                    j.invoiceRepo.Save(ctx, tenant.ID, invoice)
                } else {
                    // Atualizar status
                    existing.Status = mapAsaasStatus(asaasInvoice.Status)
                    existing.DataPagamento = asaasInvoice.PaidDate
                    existing.AtualizadoEm = time.Now()
                    j.invoiceRepo.Update(ctx, tenant.ID, existing)
                }
            }
        }
    }
    
    j.logger.Info(\"Asaas sync job completed\")
}
```

**Monitoramento:**
- [ ] Registrar duração total
- [ ] Registrar quantidade de invoices sincronizadas
- [ ] Alertar se falhar 3x consecutivas

---

### 2️⃣ Gerar Snapshot Financeiro (03:00)

**Schedule:** `0 3 * * *`  
**Timeout:** 10 minutos

**Objetivo:** Calcular e armazenar resumo financeiro diário para cada tenant.

```go
type GenerateFinancialSnapshotJob struct {
    receitaRepo    domain.ReceitaRepository
    despesaRepo    domain.DespesaRepository
    snapshotRepo   domain.SnapshotRepository
    logger         *zap.Logger
}

func (j *GenerateFinancialSnapshotJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()
    
    j.logger.Info(\"Starting financial snapshot job\")
    
    tenants, _ := j.tenantRepo.FindActive(ctx)
    yesterday := time.Now().AddDate(0, 0, -1)
    
    for _, tenant := range tenants {
        // Calcular fluxo do dia anterior
        entradas, _ := j.receitaRepo.SumByTenantAndDate(
            ctx, tenant.ID, yesterday, domain.ReceiptReceived)
        
        saidas, _ := j.despesaRepo.SumByTenantAndDate(
            ctx, tenant.ID, yesterday, domain.ExpensePaid)
        
        saldo := entradas.Sub(saidas)
        
        // Armazenar snapshot
        snapshot := &domain.FinancialSnapshot{
            ID:         uuid.NewString(),
            TenantID:   tenant.ID,
            Data:       yesterday,
            Entradas:   entradas,
            Saidas:     saidas,
            Saldo:      saldo,
            CriadoEm:   time.Now(),
        }
        
        j.snapshotRepo.Save(ctx, tenant.ID, snapshot)
        
        // Detectar anomalias
        if j.detectarAnomalia(ctx, tenant.ID, saldo) {
            j.logger.Warn(\"Anomaly detected\",
                zap.String(\"tenant_id\", tenant.ID),
                zap.String(\"saldo\", saldo.String()))
        }
    }
    
    j.logger.Info(\"Financial snapshot job completed\")
}

func (j *GenerateFinancialSnapshotJob) detectarAnomalia(
    ctx context.Context, tenantID string, saldoHoje decimal.Decimal) bool {
    
    // Buscar saldo de 7 dias atrás
    saldoSemanaPassada, _ := j.snapshotRepo.FindByTenantAndDate(
        ctx, tenantID, time.Now().AddDate(0, 0, -7))
    
    if saldoSemanaPassada == nil {
        return false
    }
    
    // Queda > 50%?
    percentualQueda := decimal.One.Sub(
        saldoHoje.Div(saldoSemanaPassada.Saldo)).Mul(decimal.NewFromInt(100))
    
    return percentualQueda.GreaterThan(decimal.NewFromInt(50))
}
```

---

### 3️⃣ Processar Repassos de Assinatura (04:00)

**Schedule:** `0 4 * * *`  
**Timeout:** 10 minutos

**Objetivo:** Criar receitas e despesas de comissão para assinaturas pagas.

(Detalhado em `ASSINATURAS.md`)

---

### 4️⃣ Alertas de Anomalias (08:00)

**Schedule:** `0 8 * * *`  
**Timeout:** 5 minutos

**Objetivo:** Verificar condições e enviar alertas (futuro: Telegram/Email).

```go
type AlertsJob struct {
    snapshotRepo   domain.SnapshotRepository
    invoiceRepo    domain.AssinaturaInvoiceRepository
    receitaRepo    domain.ReceitaRepository
    logger         *zap.Logger
    // alertService  *AlertService (futuro)
}

func (j *AlertsJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    tenants, _ := j.tenantRepo.FindActive(ctx)
    
    for _, tenant := range tenants {
        // Alert 1: Receita 0 nos últimos 3 dias
        receitas3d, _ := j.receitaRepo.SumByTenantAndPeriod(
            ctx, tenant.ID,
            time.Now().AddDate(0, 0, -3), time.Now())
        
        if receitas3d.IsZero() {
            j.logger.Warn(\"Alert: Zero revenue\", zap.String(\"tenant_id\", tenant.ID))
            // j.alertService.SendAlert(tenant.ID, \"No revenue in 3 days\")
        }
        
        // Alert 2: Despesas > Receitas
        despesas, _ := j.despesaRepo.SumByTenantAndDate(ctx, tenant.ID, time.Now())
        receitas, _ := j.receitaRepo.SumByTenantAndDate(ctx, tenant.ID, time.Now())
        
        if despesas.GreaterThan(receitas) {
            j.logger.Warn(\"Alert: Expenses > Revenue\", zap.String(\"tenant_id\", tenant.ID))
        }
        
        // Alert 3: Faturas vencidas não pagas
        overdueInvoices, _ := j.invoiceRepo.FindOverdue(ctx, tenant.ID)
        if len(overdueInvoices) > 0 {
            j.logger.Warn(\"Alert: Overdue invoices\",
                zap.String(\"tenant_id\", tenant.ID),
                zap.Int(\"count\", len(overdueInvoices)))
        }
    }
}
```

---

## 📅 Crons Semanais

### Limpeza de Dados (Segundas, 04:00)

**Schedule:** `0 4 * * 0`  
**Timeout:** 5 minutos

```go
type CleanupJob struct {
    auditLogRepo domain.AuditLogRepository
    logger       *zap.Logger
}

func (j *CleanupJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    // Deletar logs de auditoria com mais de 90 dias
    cutoffDate := time.Now().AddDate(0, -3, 0)
    deleted, _ := j.auditLogRepo.DeleteBefore(ctx, cutoffDate)
    
    j.logger.Info(\"Cleanup completed\", zap.Int64(\"deleted_records\", deleted))
}
```

---

## 🗄️ Tabela de Execução

Para auditoria e monitoramento:

```sql
CREATE TABLE cron_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_name VARCHAR(100) NOT NULL,
    tenant_id UUID,
    status VARCHAR(20) NOT NULL, -- SUCCESS, ERROR, TIMEOUT
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_seconds INT,
    error_message TEXT,
    records_processed INT,
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_cron_job_name_status ON cron_executions(job_name, status);
CREATE INDEX idx_cron_timestamp ON cron_executions(timestamp DESC);
```

---

## 🔧 Monitoramento

### Alertas

```
Condicao: Cron não executou nas últimas 25h
Acao: Alert para operacional
Severity: HIGH

Condicao: Cron falhou 3 vezes consecutivas
Acao: Desabilitar job e notificar
Severity: CRITICAL
```

### Métricas Prometheus

```go
// Métricas para monitorar crons
type CronMetrics struct {
    jobDuration    prometheus.HistogramVec
    jobErrors      prometheus.CounterVec
    lastExecutionTime prometheus.GaugeVec
}

func (m *CronMetrics) RecordExecution(
    jobName string, duration time.Duration, err error) {
    
    m.jobDuration.WithLabelValues(jobName).Observe(duration.Seconds())
    
    if err != nil {
        m.jobErrors.WithLabelValues(jobName, err.Error()).Inc()
    }
    
    m.lastExecutionTime.WithLabelValues(jobName).Set(
        float64(time.Now().Unix()))
}
```

---

## ⚠️ Error Handling

### Retry com Backoff

```go
func (j *SyncAsaasInvoicesJob) RunWithRetry() {
    maxRetries := 3
    backoff := time.Minute
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := j.Run()
        if err == nil {
            return
        }
        
        j.logger.Error(\"Job failed\",
            zap.Int(\"attempt\", attempt),
            zap.Error(err))
        
        if attempt < maxRetries {
            time.Sleep(backoff)
            backoff *= 2 // Exponential backoff
        }
    }
    
    j.logger.Error(\"Job failed after max retries\")
}
```

---

## 📊 Scheduler Configuration

```go
// Setup inicial
func NewScheduler(logger *zap.Logger) *cron.Cron {
    scheduler := cron.New(cron.WithLocation(time.UTC))
    
    // Registrar jobs
    scheduler.AddJob(\"0 2 * * *\", &SyncAsaasInvoicesJob{...})
    scheduler.AddJob(\"0 3 * * *\", &GenerateFinancialSnapshotJob{...})
    scheduler.AddJob(\"0 4 * * *\", &ProcessSubscriptionRepassJob{...})
    scheduler.AddJob(\"0 8 * * *\", &AlertsJob{...})
    scheduler.AddJob(\"0 4 * * 0\", &CleanupJob{...})
    
    scheduler.Start()
    
    return scheduler
}
```

---

**Status:** ✅ Design completo
