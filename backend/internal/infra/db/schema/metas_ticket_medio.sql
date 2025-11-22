-- Metas de ticket médio (geral ou por barbeiro)
-- Tabela: metas_ticket_medio

CREATE TABLE IF NOT EXISTS metas_ticket_medio (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    mes_ano VARCHAR(7) NOT NULL, -- Formato YYYY-MM
    meta_valor NUMERIC(10,2) NOT NULL CHECK (meta_valor > 0),
    tipo VARCHAR(20) DEFAULT 'GERAL' CHECK (tipo IN ('GERAL', 'BARBEIRO')),
    barbeiro_id UUID,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT metas_ticket_medio_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT metas_ticket_medio_barbeiro_id_fkey FOREIGN KEY (barbeiro_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_barbeiro_tipo CHECK (
        (tipo = 'GERAL' AND barbeiro_id IS NULL) OR
        (tipo = 'BARBEIRO' AND barbeiro_id IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_metas_ticket_tenant ON metas_ticket_medio(tenant_id);
CREATE INDEX IF NOT EXISTS idx_metas_ticket_mes_ano ON metas_ticket_medio(tenant_id, mes_ano DESC);
CREATE INDEX IF NOT EXISTS idx_metas_ticket_barbeiro ON metas_ticket_medio(barbeiro_id) WHERE barbeiro_id IS NOT NULL;

COMMENT ON TABLE metas_ticket_medio IS 'Metas de ticket médio (geral ou por barbeiro)';
