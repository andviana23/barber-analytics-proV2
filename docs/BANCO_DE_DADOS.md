# 🗄️ Design do Banco de Dados

**Versão:** 2.0  
**Data:** 14/11/2025  
**Status:** Finalizado

---

## 📋 Índice

1. [Diagrama ER](#diagrama-er)
2. [Tabelas Core](#tabelas-core)
3. [Índices & Performance](#índices--performance)
4. [Migrations](#migrations)
5. [Backup & Recovery](#backup--recovery)

---

## 🎯 Diagrama ER

```
┌─────────────┐
│   tenants   │
├─────────────┤
│ id (PK)     │
│ nome        │
│ cnpj        │
│ ativo       │
│ plano       │
└─────────────┘
      │ 1
      │
      │ n
      ↓
┌─────────────┐        ┌──────────────┐        ┌──────────────┐
│    users    │        │  categorias  │        │ receitas     │
├─────────────┤        ├──────────────┤        ├──────────────┤
│ id (PK)     │        │ id (PK)      │        │ id (PK)      │
│ tenant_id   │        │ tenant_id    │        │ tenant_id    │
│ email       │        │ nome         │        │ descricao    │
│ password    │        │ tipo         │        │ valor        │
│ nome        │        │ ativa        │        │ categoria_id │
│ role        │        └──────────────┘        │ data         │
│ ativo       │               ↑                └──────────────┘
└─────────────┘               │
      │                       │
      │ n                     │ n
      │                       │
      ├─────────────────────────
      │
      │ 1 ┌──────────────┐     ┌───────────────────┐
      │───│  despesas    │     │ planos_assinatura │
      │   ├──────────────┤     ├───────────────────┤
      │   │ id (PK)      │     │ id (PK)           │
      │   │ tenant_id    │     │ tenant_id         │
      │   │ descricao    │     │ nome              │
      │   │ valor        │     │ valor             │
      │   │ categoria_id │     │ periodicidade     │
      │   │ data         │     └───────────────────┘
      │   └──────────────┘              │
      │                                 │ n
      │                                 │
      │                                 ↓
      │                         ┌──────────────────┐
      │                         │  assinaturas     │
      │                         ├──────────────────┤
      │                         │ id (PK)          │
      │                         │ tenant_id        │
      │                         │ plan_id          │
      │                         │ barbeiro_id      │
      │                         │ asaas_sub_id     │
      │                         │ status           │
      │                         │ data_inicio      │
      │                         └──────────────────┘
      │                                 │
      │                                 │ n
      │                                 │
      │                                 ↓
      │                      ┌─────────────────────┐
      │                      │assinatura_invoices  │
      │                      ├─────────────────────┤
      │                      │ id (PK)             │
      │                      │ tenant_id           │
      │                      │ assinatura_id       │
      │                      │ asaas_invoice_id    │
      │                      │ valor               │
      │                      │ status              │
      │                      │ data_vencimento     │
      │                      └─────────────────────┘
      │
      └────→ ┌────────────────┐
             │  audit_logs    │
             ├────────────────┤
             │ id (PK)        │
             │ tenant_id      │
             │ user_id        │
             │ action         │
             │ resource       │
             │ timestamp      │
             └────────────────┘
```

---

## 🏗️ Tabelas Core

### tenants

Tabela de barbearias/clientes do SaaS.

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    cnpj VARCHAR(14) UNIQUE,
    ativo BOOLEAN DEFAULT true,
    plano VARCHAR(50) DEFAULT 'free', -- free, pro, enterprise
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE tenants IS 'Cada barbearia é um tenant';
COMMENT ON COLUMN tenants.cnpj IS 'Opcional, pode ser NULL inicialmente';
COMMENT ON COLUMN tenants.plano IS 'Define features disponíveis';
```

### users

Usuários do sistema, sempre pertencendo a um tenant.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee',
        -- owner, manager, accountant, employee, barbeiro
    ativo BOOLEAN DEFAULT true,
    ultimo_login TIMESTAMP,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
```

### categorias

Categorias de receitas/despesas.

```sql
CREATE TABLE categorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('RECEITA', 'DESPESA')),
    cor VARCHAR(7) DEFAULT '#000000',
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, nome)
);

CREATE INDEX idx_categorias_tenant_tipo ON categorias(tenant_id, tipo);
```

### receitas

Registros de receita.

```sql
CREATE TABLE receitas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(15, 2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'CONFIRMADO',
    observacoes TEXT,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data DESC);
CREATE INDEX idx_receitas_tenant_categoria ON receitas(tenant_id, categoria_id);
CREATE INDEX idx_receitas_tenant_status ON receitas(tenant_id, status);
```

### despesas

Registros de despesa.

```sql
CREATE TABLE despesas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor DECIMAL(15, 2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    fornecedor VARCHAR(255),
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'PENDENTE',
    observacoes TEXT,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_despesas_tenant_id ON despesas(tenant_id);
CREATE INDEX idx_despesas_tenant_data ON despesas(tenant_id, data DESC);
CREATE INDEX idx_despesas_tenant_status ON despesas(tenant_id, status);
```

### planos_assinatura

Planos oferecidos pelo sistema.

```sql
CREATE TABLE planos_assinatura (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    valor DECIMAL(10, 2) NOT NULL CHECK (valor > 0),
    periodicidade VARCHAR(50) NOT NULL,
    quantidade_servicos INT DEFAULT 0,
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, nome)
);
```

### assinaturas

Assinaturas ativas/históricas.

```sql
CREATE TABLE assinaturas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES planos_assinatura(id) ON DELETE RESTRICT,
    barbeiro_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    asaas_subscription_id VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ATIVA',
    data_inicio DATE NOT NULL,
    data_fim DATE,
    proxima_fatura_data DATE NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_assinaturas_tenant ON assinaturas(tenant_id);
CREATE INDEX idx_assinaturas_status ON assinaturas(tenant_id, status);
CREATE INDEX idx_assinaturas_barbeiro ON assinaturas(barbeiro_id);
```

### assinatura_invoices

Faturas de assinatura sincronizadas com Asaas.

```sql
CREATE TABLE assinatura_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    assinatura_id UUID NOT NULL REFERENCES assinaturas(id) ON DELETE CASCADE,
    asaas_invoice_id VARCHAR(255) UNIQUE NOT NULL,
    valor DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDENTE',
    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    processada BOOLEAN DEFAULT false,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_invoices_tenant ON assinatura_invoices(tenant_id);
CREATE INDEX idx_invoices_status ON assinatura_invoices(status);
CREATE INDEX idx_invoices_vencimento ON assinatura_invoices(data_vencimento);
```

### audit_logs

Auditoria de ações críticas.

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL, -- CREATE, UPDATE, DELETE, READ
    resource VARCHAR(100) NOT NULL, -- receita, despesa, assinatura
    resource_id VARCHAR(255),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_timestamp ON audit_logs(tenant_id, timestamp DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource, resource_id);
```

---

## ⚡ Índices & Performance

### Estratégia de Indexação

```sql
-- Índices obrigatórios
CREATE INDEX idx_receitas_tenant_id ON receitas(tenant_id);
CREATE INDEX idx_receitas_tenant_data ON receitas(tenant_id, data DESC);

-- Índices para queries comuns
CREATE INDEX idx_receitas_tenant_categoria ON receitas(tenant_id, categoria_id);
CREATE INDEX idx_receitas_status ON receitas(tenant_id, status);

-- Índices para join
CREATE INDEX idx_receitas_usuario ON receitas(usuario_id);
```

### Query Optimization

```go
// ❌ Lento: Sem índice
SELECT * FROM receitas WHERE data = '2024-11-14' AND valor > 50;

// ✅ Rápido: Com índice composto
SELECT * FROM receitas 
WHERE tenant_id = 'abc' AND data = '2024-11-14' 
  AND valor > 50;
```

---

## 📜 Migrations

Usar `golang-migrate`:

```
migrations/
├── 001_create_tenants.up.sql
├── 001_create_tenants.down.sql
├── 002_create_users.up.sql
├── 002_create_users.down.sql
├── 003_create_categorias.up.sql
├── 003_create_categorias.down.sql
├── 004_create_receitas.up.sql
├── 004_create_receitas.down.sql
└── ...
```

Executar:
```bash
migrate -path ./migrations -database "postgres://..." up
```

---

## 💾 Backup & Recovery

### Provedor: Neon

- ✅ Backup automático diário
- ✅ Point-in-time recovery (últimos 7 dias)
- ✅ Replicação automática
- ✅ Snapshots para branches

### Manual Backup

```bash
# Exportar dados
pg_dump -Fc $DATABASE_URL > backup-$(date +%Y%m%d).dump

# Restaurar
pg_restore -d $DATABASE_URL backup-20241114.dump
```

### RTO/RPO

- **RTO (Recovery Time Objective):** < 2 horas
- **RPO (Recovery Point Objective):** 24 horas

---

**Status:** ✅ Pronto para implementação
