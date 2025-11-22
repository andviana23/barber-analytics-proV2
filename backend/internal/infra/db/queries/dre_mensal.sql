-- name: CreateDREMensal :one
INSERT INTO dre_mensal (
    tenant_id,
    mes_ano,
    receita_servicos,
    receita_produtos,
    receita_planos,
    receita_total,
    custo_comissoes,
    custo_insumos,
    custo_variavel_total,
    despesa_fixa,
    despesa_variavel,
    despesa_total,
    resultado_bruto,
    resultado_operacional,
    margem_bruta,
    margem_operacional,
    lucro_liquido,
    processado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;

-- name: GetDREMensalByID :one
SELECT * FROM dre_mensal
WHERE id = $1 AND tenant_id = $2;

-- name: GetDREMensalByMesAno :one
SELECT * FROM dre_mensal
WHERE tenant_id = $1 AND mes_ano = $2;

-- name: ListDREMensalByTenant :many
SELECT * FROM dre_mensal
WHERE tenant_id = $1
ORDER BY mes_ano DESC
LIMIT $2 OFFSET $3;

-- name: ListDREMensalByPeriod :many
SELECT * FROM dre_mensal
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3
ORDER BY mes_ano DESC;

-- name: UpdateDREMensal :one
UPDATE dre_mensal
SET
    receita_servicos = $3,
    receita_produtos = $4,
    receita_planos = $5,
    receita_total = $6,
    custo_comissoes = $7,
    custo_insumos = $8,
    custo_variavel_total = $9,
    despesa_fixa = $10,
    despesa_variavel = $11,
    despesa_total = $12,
    resultado_bruto = $13,
    resultado_operacional = $14,
    margem_bruta = $15,
    margem_operacional = $16,
    lucro_liquido = $17,
    processado_em = $18,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteDREMensal :exec
DELETE FROM dre_mensal
WHERE id = $1 AND tenant_id = $2;

-- name: SumReceitasByPeriod :one
SELECT
    COALESCE(SUM(receita_total), 0) as total_receitas
FROM dre_mensal
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3;

-- name: SumDespesasByPeriod :one
SELECT
    COALESCE(SUM(despesa_total), 0) as total_despesas
FROM dre_mensal
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3;

-- name: AvgMargemBrutaByPeriod :one
SELECT
    COALESCE(AVG(margem_bruta), 0) as media_margem_bruta
FROM dre_mensal
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3;

-- name: AvgMargemOperacionalByPeriod :one
SELECT
    COALESCE(AVG(margem_operacional), 0) as media_margem_operacional
FROM dre_mensal
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3;

-- name: CountDREMensalByTenant :one
SELECT COUNT(*) FROM dre_mensal
WHERE tenant_id = $1;
