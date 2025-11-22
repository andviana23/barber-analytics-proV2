-- Migration: Ajustar schema das tabelas de precificação para match com entities
-- Data: 2025-11-22
-- Autor: System
-- Descrição: Renomear/adicionar colunas para compatibilidade com domain entities

-- ============================================================
-- TABELA: precificacao_simulacoes
-- ============================================================
-- Adicionar colunas faltantes conforme entity PrecificacaoSimulacao

ALTER TABLE precificacao_simulacoes
  ADD COLUMN IF NOT EXISTS custo_materiais NUMERIC(15,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS custo_mao_de_obra NUMERIC(15,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS custo_total NUMERIC(15,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS preco_atual NUMERIC(15,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS diferenca_percentual NUMERIC(5,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS lucro_estimado NUMERIC(15,2) DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS margem_final NUMERIC(5,2) DEFAULT 0.00;

-- Migrar dados existentes (se houver)
-- custo_insumos → custo_materiais
UPDATE precificacao_simulacoes
SET custo_materiais = custo_insumos,
    custo_total = custo_insumos
WHERE custo_materiais = 0.00;

-- Renomear coluna margem_resultante → (mantém para compatibilidade)
-- Nota: margem_final é novo campo calculado

-- ============================================================
-- COMENTÁRIOS
-- ============================================================
COMMENT ON COLUMN precificacao_simulacoes.custo_materiais IS 'Custo de materiais/insumos do serviço ou produto';
COMMENT ON COLUMN precificacao_simulacoes.custo_mao_de_obra IS 'Custo estimado de mão de obra';
COMMENT ON COLUMN precificacao_simulacoes.custo_total IS 'Custo total (materiais + mão de obra)';
COMMENT ON COLUMN precificacao_simulacoes.preco_atual IS 'Preço atualmente praticado (para comparação)';
COMMENT ON COLUMN precificacao_simulacoes.diferenca_percentual IS 'Diferença percentual entre preço sugerido e atual';
COMMENT ON COLUMN precificacao_simulacoes.lucro_estimado IS 'Lucro estimado após impostos e comissões';
COMMENT ON COLUMN precificacao_simulacoes.margem_final IS 'Margem de lucro final calculada';

-- ============================================================
-- NOTAS
-- ============================================================
-- 1. parametros_json já é JSONB no schema, entity usa string → converter na aplicação
-- 2. custo_insumos mantido por compatibilidade (pode ser removido em v2)
-- 3. markup_aplicado e margem_resultante mantidos por compatibilidade
