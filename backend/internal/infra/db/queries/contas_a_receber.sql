-- name: CreateContaReceber :one
INSERT INTO contas_a_receber (
    tenant_id,
    origem,
    assinatura_id,
    servico_id,
    descricao,
    valor,
    valor_pago,
    data_vencimento,
    data_recebimento,
    status,
    observacoes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetContaReceberByID :one
SELECT * FROM contas_a_receber
WHERE id = $1 AND tenant_id = $2;

-- name: ListContasReceberByTenant :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
ORDER BY data_vencimento DESC
LIMIT $2 OFFSET $3;

-- name: ListContasReceberByStatus :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND status = $2
ORDER BY data_vencimento ASC
LIMIT $3 OFFSET $4;

-- name: ListContasReceberByPeriod :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
  AND data_vencimento >= $2
  AND data_vencimento <= $3
ORDER BY data_vencimento ASC;

-- name: ListContasReceberVencidas :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
  AND status IN ('PENDENTE', 'ATRASADO')
  AND data_vencimento < $2
ORDER BY data_vencimento ASC;

-- name: ListContasReceberByAssinatura :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND assinatura_id = $2
ORDER BY data_vencimento DESC;

-- name: ListContasReceberByOrigem :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND origem = $2
ORDER BY data_vencimento DESC
LIMIT $3 OFFSET $4;

-- name: UpdateContaReceber :one
UPDATE contas_a_receber
SET
    descricao = $3,
    valor = $4,
    valor_pago = $5,
    data_vencimento = $6,
    data_recebimento = $7,
    status = $8,
    observacoes = $9,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaReceberComoRecebida :one
UPDATE contas_a_receber
SET
    status = 'RECEBIDO',
    data_recebimento = $3,
    valor_pago = $4,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaReceberComoAtrasada :exec
UPDATE contas_a_receber
SET
    status = 'ATRASADO',
    atualizado_em = NOW()
WHERE tenant_id = $1
  AND status = 'PENDENTE'
  AND data_vencimento < $2;

-- name: DeleteContaReceber :exec
DELETE FROM contas_a_receber
WHERE id = $1 AND tenant_id = $2;

-- name: SumContasReceberByPeriod :one
SELECT
    COALESCE(SUM(valor), 0) as total_a_receber
FROM contas_a_receber
WHERE tenant_id = $1
  AND data_vencimento >= $2
  AND data_vencimento <= $3
  AND status != 'CANCELADO';

-- name: SumContasRecebidasByPeriod :one
SELECT
    COALESCE(SUM(valor_pago), 0) as total_recebido
FROM contas_a_receber
WHERE tenant_id = $1
  AND data_recebimento >= $2
  AND data_recebimento <= $3
  AND status = 'RECEBIDO';

-- name: CountContasReceberByStatus :one
SELECT COUNT(*) FROM contas_a_receber
WHERE tenant_id = $1 AND status = $2;

-- name: CountContasReceberByTenant :one
SELECT COUNT(*) FROM contas_a_receber
WHERE tenant_id = $1;
