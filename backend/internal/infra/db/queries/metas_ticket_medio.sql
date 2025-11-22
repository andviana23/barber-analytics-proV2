-- name: CreateMetaTicketMedio :one
INSERT INTO metas_ticket_medio (
    tenant_id,
    mes_ano,
    meta_valor,
    tipo,
    barbeiro_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetMetaTicketMedioByID :one
SELECT * FROM metas_ticket_medio
WHERE id = $1 AND tenant_id = $2;

-- name: GetMetaTicketMedioGeralByMesAno :one
SELECT * FROM metas_ticket_medio
WHERE tenant_id = $1
  AND mes_ano = $2
  AND tipo = 'GERAL';

-- name: GetMetaTicketMedioBarbeiroByMesAno :one
SELECT * FROM metas_ticket_medio
WHERE tenant_id = $1
  AND mes_ano = $2
  AND barbeiro_id = $3
  AND tipo = 'BARBEIRO';

-- name: ListMetasTicketMedioByTenant :many
SELECT * FROM metas_ticket_medio
WHERE tenant_id = $1
ORDER BY mes_ano DESC, tipo, barbeiro_id
LIMIT $2 OFFSET $3;

-- name: ListMetasTicketMedioByMesAno :many
SELECT * FROM metas_ticket_medio
WHERE tenant_id = $1 AND mes_ano = $2
ORDER BY tipo, barbeiro_id;

-- name: ListMetasTicketMedioByBarbeiro :many
SELECT * FROM metas_ticket_medio
WHERE tenant_id = $1
  AND barbeiro_id = $2
  AND tipo = 'BARBEIRO'
ORDER BY mes_ano DESC
LIMIT $3 OFFSET $4;

-- name: UpdateMetaTicketMedio :one
UPDATE metas_ticket_medio
SET
    meta_valor = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteMetaTicketMedio :exec
DELETE FROM metas_ticket_medio
WHERE id = $1 AND tenant_id = $2;

-- name: CountMetasTicketMedioByTenant :one
SELECT COUNT(*) FROM metas_ticket_medio
WHERE tenant_id = $1;

-- name: CountMetasTicketMedioByBarbeiro :one
SELECT COUNT(*) FROM metas_ticket_medio
WHERE tenant_id = $1 AND barbeiro_id = $2;
