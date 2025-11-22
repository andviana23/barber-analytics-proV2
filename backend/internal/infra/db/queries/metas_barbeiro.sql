-- name: CreateMetaBarbeiro :one
INSERT INTO metas_barbeiro (
    tenant_id,
    barbeiro_id,
    mes_ano,
    meta_servicos_gerais,
    meta_servicos_extras,
    meta_produtos
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMetaBarbeiroByID :one
SELECT * FROM metas_barbeiro
WHERE id = $1 AND tenant_id = $2;

-- name: GetMetaBarbeiroByMesAno :one
SELECT * FROM metas_barbeiro
WHERE tenant_id = $1
  AND barbeiro_id = $2
  AND mes_ano = $3;

-- name: ListMetasBarbeiroByTenant :many
SELECT * FROM metas_barbeiro
WHERE tenant_id = $1
ORDER BY mes_ano DESC, barbeiro_id
LIMIT $2 OFFSET $3;

-- name: ListMetasBarbeiroByBarbeiro :many
SELECT * FROM metas_barbeiro
WHERE tenant_id = $1 AND barbeiro_id = $2
ORDER BY mes_ano DESC
LIMIT $3 OFFSET $4;

-- name: ListMetasBarbeiroByMesAno :many
SELECT * FROM metas_barbeiro
WHERE tenant_id = $1 AND mes_ano = $2
ORDER BY barbeiro_id;

-- name: UpdateMetaBarbeiro :one
UPDATE metas_barbeiro
SET
    meta_servicos_gerais = $3,
    meta_servicos_extras = $4,
    meta_produtos = $5,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteMetaBarbeiro :exec
DELETE FROM metas_barbeiro
WHERE id = $1 AND tenant_id = $2;

-- name: CountMetasBarbeiroByTenant :one
SELECT COUNT(*) FROM metas_barbeiro
WHERE tenant_id = $1;

-- name: CountMetasBarbeiroByBarbeiro :one
SELECT COUNT(*) FROM metas_barbeiro
WHERE tenant_id = $1 AND barbeiro_id = $2;
