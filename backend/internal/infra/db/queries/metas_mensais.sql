-- name: CreateMetaMensal :one
INSERT INTO metas_mensais (
    tenant_id,
    mes_ano,
    meta_faturamento,
    origem,
    status,
    criado_por
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMetaMensalByID :one
SELECT * FROM metas_mensais
WHERE id = $1 AND tenant_id = $2;

-- name: GetMetaMensalByMesAno :one
SELECT * FROM metas_mensais
WHERE tenant_id = $1 AND mes_ano = $2;

-- name: ListMetasMensaisByTenant :many
SELECT * FROM metas_mensais
WHERE tenant_id = $1
ORDER BY mes_ano DESC
LIMIT $2 OFFSET $3;

-- name: ListMetasMensaisByStatus :many
SELECT * FROM metas_mensais
WHERE tenant_id = $1 AND status = $2
ORDER BY mes_ano DESC
LIMIT $3 OFFSET $4;

-- name: ListMetasMensaisByPeriod :many
SELECT * FROM metas_mensais
WHERE tenant_id = $1
  AND mes_ano >= $2
  AND mes_ano <= $3
ORDER BY mes_ano DESC;

-- name: UpdateMetaMensal :one
UPDATE metas_mensais
SET
    meta_faturamento = $3,
    origem = $4,
    status = $5,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: AprovarMetaMensal :one
UPDATE metas_mensais
SET
    status = 'ACEITA',
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: RejeitarMetaMensal :one
UPDATE metas_mensais
SET
    status = 'REJEITADA',
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteMetaMensal :exec
DELETE FROM metas_mensais
WHERE id = $1 AND tenant_id = $2;

-- name: CountMetasMensaisByTenant :one
SELECT COUNT(*) FROM metas_mensais
WHERE tenant_id = $1;
