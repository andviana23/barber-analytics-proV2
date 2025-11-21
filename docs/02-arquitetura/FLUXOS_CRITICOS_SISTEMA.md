> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Fluxos Críticos do Sistema

## Assinatura (Asaas)
1. Frontend cria/edita assinatura → backend valida tenant/planos → Asaas API v3 (sandbox/prod) cria/atualiza assinatura.
2. Faturas geradas em Asaas → cron de sync (FLUXO_CRONS.md) busca invoices → grava em PostgreSQL (Neon) → atualiza status e repasse.
3. Repasse/comissão: calculado no backend → registrado no financeiro → enviado a dashboards.

## Agendamento / Lista da Vez (Barber Turns)
1. Profissionais BARBEIRO entram na fila (lista da vez) → regras de pontos/ordem aplicadas.
2. Atendimento registrado → histórico e estatísticas atualizados (barber_turn_history) → UI reflete próximo barbeiro.
3. Alertas/cron podem recalcular/resetar lista conforme regras (vide docs/listadavez.md).

## Relatórios Financeiros
1. Receitas/Despesas lançadas via frontend → backend persiste em Neon com tenant_id.
2. Fluxo de caixa calcula saldo e snapshots (cron) → dashboards consumem via TanStack Query.
3. Export/visualização: endpoints em API pública (API_PUBLICA.md) com filtros e RBAC.

Referências: FLUXO_CRONS.md, FINANCEIRO.md, ASSINATURAS.md, listadavez.md, API_REFERENCE.md.
