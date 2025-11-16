-- Migrates MVP 1.0 financial data into Barber Analytics Pro v2
-- Ajuste os nomes (mvp_tables) caso o banco legacy use outras nomenclaturas.
-- O script deve ser executado dentro do mesmo tenant_id atual.

BEGIN;

-- Controle de staging para evitar duplicação
WITH
staging_receitas AS (
    SELECT
        id,
        tenant_id,
        descricao,
        valor,
        categoria,
        metodo_pagamento,
        data,
        status,
        observacoes,
        origem_manual,
        criado_em,
        atualizado_em
    FROM mvp_receitas -- substitua pelo nome original
    WHERE migrated_to_v2 IS NULL
),
staging_despesas AS (
    SELECT
        id,
        tenant_id,
        descricao,
        valor,
        categoria,
        fornecedor,
        metodo_pagamento,
        data,
        status,
        observacoes,
        origem_manual,
        criado_em,
        atualizado_em
    FROM mvp_despesas -- substitua pelo nome original
    WHERE migrated_to_v2 IS NULL
),
staging_assinaturas AS (
    SELECT
        id,
        tenant_id,
        plano,
        barbeiro_id,
        valor,
        periodicidade,
        status,
        data_inicio,
        data_fim,
        proxima_fatura_data,
        origem_manual,
        criado_em,
        atualizado_em
    FROM mvp_assinaturas -- substitua pelo nome original
    WHERE migrated_to_v2 IS NULL
)

-- Inserções de receitas
INSERT INTO receitas (
    id, tenant_id, usuario_id, descricao, valor, categoria_id,
    metodo_pagamento, data, status, observacoes, criado_em, atualizado_em, manual, origem_dado
)
SELECT
    id,
    tenant_id,
    NULL,
    descricao,
    valor,
    (SELECT id FROM categorias WHERE tenant_id = staging_receitas.tenant_id AND nome = staging_receitas.categoria LIMIT 1),
    metodo_pagamento,
    data,
    status,
    observacoes,
    criado_em,
    atualizado_em,
    origem_manual::boolean,
    'mvp'
FROM staging_receitas
ON CONFLICT (id) DO UPDATE SET
    descricao = EXCLUDED.descricao,
    valor = EXCLUDED.valor,
    categoria_id = EXCLUDED.categoria_id,
    status = EXCLUDED.status,
    observacoes = EXCLUDED.observacoes,
    atualizado_em = EXCLUDED.atualizado_em,
    manual = EXCLUDED.manual,
    origem_dado = EXCLUDED.origem_dado;

-- Inserções de despesas
INSERT INTO despesas (
    id, tenant_id, usuario_id, descricao, valor, categoria_id, fornecedor,
    metodo_pagamento, data, status, observacoes, criado_em, atualizado_em, manual, origem_dado
)
SELECT
    id,
    tenant_id,
    NULL,
    descricao,
    valor,
    (SELECT id FROM categorias WHERE tenant_id = staging_despesas.tenant_id AND nome = staging_despesas.categoria LIMIT 1),
    fornecedor,
    metodo_pagamento,
    data,
    status,
    observacoes,
    criado_em,
    atualizado_em,
    origem_manual::boolean,
    'mvp'
FROM staging_despesas
ON CONFLICT (id) DO UPDATE SET
    descricao = EXCLUDED.descricao,
    valor = EXCLUDED.valor,
    categoria_id = EXCLUDED.categoria_id,
    fornecedor = EXCLUDED.fornecedor,
    status = EXCLUDED.status,
    observacoes = EXCLUDED.observacoes,
    atualizado_em = EXCLUDED.atualizado_em,
    manual = EXCLUDED.manual,
    origem_dado = EXCLUDED.origem_dado;

-- Inserções de assinaturas manuais
INSERT INTO assinaturas (
    id, tenant_id, plan_id, barbeiro_id, asaas_subscription_id,
    status, data_inicio, data_fim, proxima_fatura_data, criado_em, atualizado_em, data_proximo_pagamento, origem_dado
)
SELECT
    id,
    tenant_id,
    (SELECT id FROM planos_assinatura WHERE tenant_id = staging_assinaturas.tenant_id AND nome = staging_assinaturas.plano LIMIT 1),
    barbeiro_id,
    NULL,
    status,
    data_inicio,
    data_fim,
    proxima_fatura_data,
    criado_em,
    atualizado_em,
    proxima_fatura_data,
    'mvp'
FROM staging_assinaturas
ON CONFLICT (id) DO UPDATE SET
    plan_id = EXCLUDED.plan_id,
    barbeiro_id = EXCLUDED.barbeiro_id,
    status = EXCLUDED.status,
    data_inicio = EXCLUDED.data_inicio,
    data_fim = EXCLUDED.data_fim,
    atualizado_em = EXCLUDED.atualizado_em,
    origem_dado = EXCLUDED.origem_dado;

-- Marca os registros legacy como migrados (ajuste para o nome correto da tabela)
UPDATE mvp_receitas SET migrated_to_v2 = TRUE WHERE id IN (SELECT id FROM staging_receitas);
UPDATE mvp_despesas SET migrated_to_v2 = TRUE WHERE id IN (SELECT id FROM staging_despesas);
UPDATE mvp_assinaturas SET migrated_to_v2 = TRUE WHERE id IN (SELECT id FROM staging_assinaturas);

COMMIT;

-- Validações pós-migração
-- 1. Verificar contagens por tenant
SELECT 'receitas' AS origem, tenant_id, COUNT(*) AS total FROM staging_receitas GROUP BY tenant_id;
SELECT 'despesas' AS origem, tenant_id, COUNT(*) AS total FROM staging_despesas GROUP BY tenant_id;
SELECT 'assinaturas' AS origem, tenant_id, COUNT(*) AS total FROM staging_assinaturas GROUP BY tenant_id;

-- 2. Validar somas e disparidades entre MVP e v2
WITH mvp_totals AS (
    SELECT tenant_id, SUM(valor) AS total_mvp FROM mvp_receitas GROUP BY tenant_id
), v2_totals AS (
    SELECT tenant_id, SUM(valor) AS total_v2 FROM receitas GROUP BY tenant_id
)
SELECT mvp_totals.tenant_id, total_mvp, total_v2, (total_mvp - total_v2) AS diff
FROM mvp_totals
LEFT JOIN v2_totals ON mvp_totals.tenant_id = v2_totals.tenant_id;

-- 3. Smoke test do CalculateCashflowUseCase (ajuste se necessário) - comparando saldo antes/depois
-- Antes da migração execute:
-- SELECT tenant_id, SUM(valor) - SUM(valor) AS saldo_mvp FROM mvp_receitas JOIN mvp_despesas ON ... (defina a regra aplicada no MVP).
-- Após a migração execute:
-- SELECT tenant_id, entradas, saidas, saldo FROM financial_snapshots WHERE origem_dado = 'mvp' ORDER BY periodo_inicio DESC LIMIT 10;

-- Checklist de execução/rollback
-- [ ] 1. Desabilitar jobs concorrentes (cron) e freeze de operações no MVP.
-- [ ] 2. Executar o script em um tenant tester (ex: tenant_dev) usando `psql scripts/sql/migrate_mvp_to_v2.sql`.
-- [ ] 3. Validar as contagens/somas acima após o COMMIT.
-- [ ] 4. Rever o `CalculateCashflowUseCase` para garantir que as somas utilizadas aparecem em `financial_snapshots`.
-- [ ] 5. Em caso de regressão, executar `ROLLBACK` (ou restaurar backup) e revisar logs.
