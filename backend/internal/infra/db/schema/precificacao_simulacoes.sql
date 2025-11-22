-- Histórico de simulações de precificação com parâmetros e resultados
-- Tabela: precificacao_simulacoes

CREATE TABLE IF NOT EXISTS precificacao_simulacoes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    item_id UUID,
    tipo_item VARCHAR(20) CHECK (tipo_item IN ('SERVICO', 'PRODUTO')),

    -- Custos detalhados (match com entity)
    custo_materiais NUMERIC(15,2) DEFAULT 0.00,
    custo_mao_de_obra NUMERIC(15,2) DEFAULT 0.00,
    custo_total NUMERIC(15,2) DEFAULT 0.00,

    -- Percentuais
    margem_desejada NUMERIC(5,2) NOT NULL,
    comissao_percentual NUMERIC(5,2) NOT NULL,
    imposto_percentual NUMERIC(5,2) NOT NULL,

    -- Preços e resultados
    preco_sugerido NUMERIC(15,2) NOT NULL,
    preco_atual NUMERIC(15,2) DEFAULT 0.00,
    diferenca_percentual NUMERIC(5,2) DEFAULT 0.00,

    -- Lucro e margem final
    lucro_estimado NUMERIC(15,2) DEFAULT 0.00,
    margem_final NUMERIC(5,2) DEFAULT 0.00,

    -- Campos legados (manter compatibilidade)
    custo_insumos NUMERIC(15,2) DEFAULT 0.00,
    markup_aplicado NUMERIC(5,2) DEFAULT 0.00,
    margem_resultante NUMERIC(5,2) DEFAULT 0.00,

    parametros_json JSONB,
    criado_por UUID,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT precificacao_simulacoes_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT precificacao_simulacoes_criado_por_fkey FOREIGN KEY (criado_por) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_precificacao_simulacoes_tenant ON precificacao_simulacoes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_precificacao_simulacoes_item ON precificacao_simulacoes(item_id) WHERE item_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_precificacao_simulacoes_criado_em ON precificacao_simulacoes(tenant_id, criado_em DESC);

COMMENT ON TABLE precificacao_simulacoes IS 'Histórico de simulações de precificação com parâmetros e resultados';
