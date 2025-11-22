-- Configurações de precificação padrão por tenant (margem, markup, impostos, comissões)
-- Tabela: precificacao_config

CREATE TABLE IF NOT EXISTS precificacao_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE,

    margem_desejada NUMERIC(5,2) DEFAULT 30.00 CHECK (margem_desejada >= 5 AND margem_desejada <= 100),
    markup_alvo NUMERIC(5,2) DEFAULT 1.5 CHECK (markup_alvo >= 1),
    imposto_percentual NUMERIC(5,2) DEFAULT 0.00 CHECK (imposto_percentual >= 0 AND imposto_percentual <= 100),
    comissao_percentual_default NUMERIC(5,2) DEFAULT 30.00 CHECK (comissao_percentual_default >= 0 AND comissao_percentual_default <= 100),

    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT precificacao_config_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS precificacao_config_tenant_id_key ON precificacao_config(tenant_id);
CREATE INDEX IF NOT EXISTS idx_precificacao_config_tenant ON precificacao_config(tenant_id);

COMMENT ON TABLE precificacao_config IS 'Configurações de precificação padrão por tenant (margem, markup, impostos, comissões)';
