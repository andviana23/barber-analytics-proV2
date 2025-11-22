-- name: CreateCompensacaoBancaria :one
INSERT INTO compensacoes_bancarias (
    tenant_id,
    receita_id,
    data_transacao,
    data_compensacao,
    data_compensado,
    valor_bruto,
    taxa_percentual,
    taxa_fixa,
    valor_liquido,
    meio_pagamento_id,
    d_mais,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetCompensacaoBancariaByID :one
SELECT * FROM compensacoes_bancarias
WHERE id = $1 AND tenant_id = $2;

-- name: ListCompensacoesBancariasByTenant :many
SELECT * FROM compensacoes_bancarias
WHERE tenant_id = $1
ORDER BY data_compensacao DESC
LIMIT $2 OFFSET $3;

-- name: ListCompensacoesByStatus :many
SELECT * FROM compensacoes_bancarias
WHERE tenant_id = $1 AND status = $2
ORDER BY data_compensacao DESC
LIMIT $3 OFFSET $4;

-- name: ListCompensacoesByDataCompensacao :many
SELECT * FROM compensacoes_bancarias
WHERE tenant_id = $1
  AND data_compensacao >= $2
  AND data_compensacao <= $3
ORDER BY data_compensacao ASC;

-- name: ListCompensacoesByReceita :many
SELECT * FROM compensacoes_bancarias
WHERE tenant_id = $1 AND receita_id = $2
ORDER BY data_compensacao DESC;

-- name: UpdateCompensacaoBancaria :one
UPDATE compensacoes_bancarias
SET
    data_compensacao = $3,
    data_compensado = $4,
    valor_bruto = $5,
    taxa_percentual = $6,
    taxa_fixa = $7,
    valor_liquido = $8,
    status = $9,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarComoCompensado :one
UPDATE compensacoes_bancarias
SET
    status = 'COMPENSADO',
    data_compensado = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteCompensacaoBancaria :exec
DELETE FROM compensacoes_bancarias
WHERE id = $1 AND tenant_id = $2;

-- name: SumValorLiquidoByPeriod :one
SELECT
    COALESCE(SUM(valor_liquido), 0) as total_liquido
FROM compensacoes_bancarias
WHERE tenant_id = $1
  AND data_compensacao >= $2
  AND data_compensacao <= $3
  AND status IN ('COMPENSADO', 'CONFIRMADO');

-- name: CountCompensacoesByStatus :one
SELECT COUNT(*) FROM compensacoes_bancarias
WHERE tenant_id = $1 AND status = $2;

-- name: CountCompensacoesByTenant :one
SELECT COUNT(*) FROM compensacoes_bancarias
WHERE tenant_id = $1;
