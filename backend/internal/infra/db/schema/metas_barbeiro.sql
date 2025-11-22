-- Metas individuais por barbeiro (serviços gerais, extras, produtos)
-- Tabela: metas_barbeiro

CREATE TABLE IF NOT EXISTS metas_barbeiro (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    barbeiro_id UUID NOT NULL,
    mes_ano VARCHAR(7) NOT NULL, -- Formato YYYY-MM

    meta_servicos_gerais NUMERIC(15,2) DEFAULT 0,
    meta_servicos_extras NUMERIC(15,2) DEFAULT 0,
    meta_produtos NUMERIC(15,2) DEFAULT 0,

    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT metas_barbeiro_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT metas_barbeiro_barbeiro_id_fkey FOREIGN KEY (barbeiro_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT metas_barbeiro_tenant_id_barbeiro_id_mes_ano_key UNIQUE (tenant_id, barbeiro_id, mes_ano)
);

CREATE UNIQUE INDEX IF NOT EXISTS metas_barbeiro_tenant_id_barbeiro_id_mes_ano_key ON metas_barbeiro(tenant_id, barbeiro_id, mes_ano);
CREATE INDEX IF NOT EXISTS idx_metas_barbeiro_tenant ON metas_barbeiro(tenant_id);
CREATE INDEX IF NOT EXISTS idx_metas_barbeiro_mes_ano ON metas_barbeiro(tenant_id, mes_ano DESC);
CREATE INDEX IF NOT EXISTS idx_metas_barbeiro_barbeiro ON metas_barbeiro(barbeiro_id);

COMMENT ON TABLE metas_barbeiro IS 'Metas individuais por barbeiro (serviços gerais, extras, produtos)';
