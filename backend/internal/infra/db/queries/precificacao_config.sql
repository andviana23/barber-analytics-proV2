-- name: CreatePrecificacaoConfig :one
INSERT INTO precificacao_config (
    tenant_id,
    margem_desejada,
    markup_alvo,
    imposto_percentual,
    comissao_percentual_default
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPrecificacaoConfigByID :one
SELECT * FROM precificacao_config
WHERE id = $1 AND tenant_id = $2;

-- name: GetPrecificacaoConfigByTenant :one
SELECT * FROM precificacao_config
WHERE tenant_id = $1;

-- name: UpdatePrecificacaoConfig :one
UPDATE precificacao_config
SET
    margem_desejada = $3,
    markup_alvo = $4,
    imposto_percentual = $5,
    comissao_percentual_default = $6,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeletePrecificacaoConfig :exec
DELETE FROM precificacao_config
WHERE id = $1 AND tenant_id = $2;
