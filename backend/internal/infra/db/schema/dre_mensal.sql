-- Demonstrativo de Resultado do Exercício mensal por tenant
-- Tabela: dre_mensal

CREATE TABLE IF NOT EXISTS dre_mensal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    mes_ano VARCHAR(7) NOT NULL, -- Formato YYYY-MM

    -- Receitas
    receita_servicos NUMERIC(15,2) DEFAULT 0,
    receita_produtos NUMERIC(15,2) DEFAULT 0,
    receita_planos NUMERIC(15,2) DEFAULT 0,
    receita_total NUMERIC(15,2) DEFAULT 0,

    -- Custos variáveis
    custo_comissoes NUMERIC(15,2) DEFAULT 0,
    custo_insumos NUMERIC(15,2) DEFAULT 0,
    custo_variavel_total NUMERIC(15,2) DEFAULT 0,

    -- Despesas
    despesa_fixa NUMERIC(15,2) DEFAULT 0,
    despesa_variavel NUMERIC(15,2) DEFAULT 0,
    despesa_total NUMERIC(15,2) DEFAULT 0,

    -- Resultados
    resultado_bruto NUMERIC(15,2) DEFAULT 0,
    resultado_operacional NUMERIC(15,2) DEFAULT 0,
    margem_bruta NUMERIC(5,2) DEFAULT 0,
    margem_operacional NUMERIC(5,2) DEFAULT 0,
    lucro_liquido NUMERIC(15,2) DEFAULT 0,

    processado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT dre_mensal_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT dre_mensal_tenant_id_mes_ano_key UNIQUE (tenant_id, mes_ano)
);

CREATE UNIQUE INDEX IF NOT EXISTS dre_mensal_tenant_id_mes_ano_key ON dre_mensal(tenant_id, mes_ano);
CREATE INDEX IF NOT EXISTS idx_dre_mensal_tenant ON dre_mensal(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dre_mensal_mes_ano ON dre_mensal(tenant_id, mes_ano DESC);

COMMENT ON TABLE dre_mensal IS 'Demonstrativo de Resultado do Exercício mensal por tenant';
