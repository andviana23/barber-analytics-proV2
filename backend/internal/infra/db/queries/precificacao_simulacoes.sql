-- name: CreatePrecificacaoSimulacao :one
INSERT INTO precificacao_simulacoes (
    tenant_id,
    item_id,
    tipo_item,
    custo_materiais,
    custo_mao_de_obra,
    custo_total,
    margem_desejada,
    comissao_percentual,
    imposto_percentual,
    preco_sugerido,
    preco_atual,
    diferenca_percentual,
    lucro_estimado,
    margem_final,
    parametros_json,
    criado_por
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;

-- name: GetPrecificacaoSimulacaoByID :one
SELECT * FROM precificacao_simulacoes
WHERE id = $1 AND tenant_id = $2;

-- name: ListSimulacoesByTenant :many
SELECT * FROM precificacao_simulacoes
WHERE tenant_id = $1
ORDER BY criado_em DESC
LIMIT $2 OFFSET $3;

-- name: ListSimulacoesByItem :many
SELECT * FROM precificacao_simulacoes
WHERE tenant_id = $1 AND item_id = $2
ORDER BY criado_em DESC
LIMIT $3 OFFSET $4;

-- name: ListSimulacoesByTipoItem :many
SELECT * FROM precificacao_simulacoes
WHERE tenant_id = $1 AND tipo_item = $2
ORDER BY criado_em DESC
LIMIT $3 OFFSET $4;

-- name: ListSimulacoesByUsuario :many
SELECT * FROM precificacao_simulacoes
WHERE tenant_id = $1 AND criado_por = $2
ORDER BY criado_em DESC
LIMIT $3 OFFSET $4;

-- name: DeletePrecificacaoSimulacao :exec
DELETE FROM precificacao_simulacoes
WHERE id = $1 AND tenant_id = $2;

-- name: CountSimulacoesByTenant :one
SELECT COUNT(*) FROM precificacao_simulacoes
WHERE tenant_id = $1;

-- name: CountSimulacoesByItem :one
SELECT COUNT(*) FROM precificacao_simulacoes
WHERE tenant_id = $1 AND item_id = $2;

-- name: GetUltimaSimulacaoByItem :one
SELECT * FROM precificacao_simulacoes
WHERE tenant_id = $1 AND item_id = $2
ORDER BY criado_em DESC
LIMIT 1;
