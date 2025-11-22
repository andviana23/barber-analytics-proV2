-- Metas de faturamento mensal por tenant
-- Tabela: metas_mensais

CREATE TABLE IF NOT EXISTS metas_mensais (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    mes_ano VARCHAR(7) NOT NULL, -- Formato YYYY-MM
    meta_faturamento NUMERIC(15,2) NOT NULL CHECK (meta_faturamento > 0),
    origem VARCHAR(20) DEFAULT 'MANUAL' CHECK (origem IN ('MANUAL', 'AUTOMATICA')),
    status VARCHAR(20) DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'ACEITA', 'REJEITADA')),
    criado_por UUID,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT metas_mensais_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT metas_mensais_criado_por_fkey FOREIGN KEY (criado_por) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT metas_mensais_tenant_id_mes_ano_key UNIQUE (tenant_id, mes_ano)
);

CREATE UNIQUE INDEX IF NOT EXISTS metas_mensais_tenant_id_mes_ano_key ON metas_mensais(tenant_id, mes_ano);
CREATE INDEX IF NOT EXISTS idx_metas_mensais_tenant ON metas_mensais(tenant_id);
CREATE INDEX IF NOT EXISTS idx_metas_mensais_mes_ano ON metas_mensais(tenant_id, mes_ano DESC);

COMMENT ON TABLE metas_mensais IS 'Metas de faturamento mensal por tenant';
