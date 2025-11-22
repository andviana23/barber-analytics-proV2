-- name: CreateFluxoCaixaDiario :one
INSERT INTO fluxo_caixa_diario (
    tenant_id,
    data,
    saldo_inicial,
    saldo_final,
    entradas_confirmadas,
    entradas_previstas,
    saidas_pagas,
    saidas_previstas,
    processado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetFluxoCaixaDiarioByID :one
SELECT * FROM fluxo_caixa_diario
WHERE id = $1 AND tenant_id = $2;

-- name: GetFluxoCaixaDiarioByData :one
SELECT * FROM fluxo_caixa_diario
WHERE tenant_id = $1 AND data = $2;

-- name: ListFluxoCaixaDiarioByTenant :many
SELECT * FROM fluxo_caixa_diario
WHERE tenant_id = $1
ORDER BY data DESC
LIMIT $2 OFFSET $3;

-- name: ListFluxoCaixaDiarioByPeriod :many
SELECT * FROM fluxo_caixa_diario
WHERE tenant_id = $1
  AND data >= $2
  AND data <= $3
ORDER BY data DESC;

-- name: UpdateFluxoCaixaDiario :one
UPDATE fluxo_caixa_diario
SET
    saldo_inicial = $3,
    saldo_final = $4,
    entradas_confirmadas = $5,
    entradas_previstas = $6,
    saidas_pagas = $7,
    saidas_previstas = $8,
    processado_em = $9,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteFluxoCaixaDiario :exec
DELETE FROM fluxo_caixa_diario
WHERE id = $1 AND tenant_id = $2;

-- name: SumEntradasByPeriod :one
SELECT
    COALESCE(SUM(entradas_confirmadas), 0) as total_entradas
FROM fluxo_caixa_diario
WHERE tenant_id = $1
  AND data >= $2
  AND data <= $3;

-- name: SumSaidasByPeriod :one
SELECT
    COALESCE(SUM(saidas_pagas), 0) as total_saidas
FROM fluxo_caixa_diario
WHERE tenant_id = $1
  AND data >= $2
  AND data <= $3;

-- name: GetUltimoSaldo :one
SELECT saldo_final
FROM fluxo_caixa_diario
WHERE tenant_id = $1
  AND data < $2
ORDER BY data DESC
LIMIT 1;

-- name: CountFluxoCaixaDiarioByTenant :one
SELECT COUNT(*) FROM fluxo_caixa_diario
WHERE tenant_id = $1;
