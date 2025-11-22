# 🔎 Contexto — Estoque

## Estado atual
- Tabelas core de produtos/serviços existem, mas não há tabelas de movimentação/inventário nem lógica de estoque.
- Tasks detalhadas estão em `Tarefas/ESTOQUE/*.md`, mas dependem de financeiro (contas a pagar) e de consumo por serviço.

## Dependências
- Payables/receivables prontos para opcionalmente gerar financeiro a partir de entradas/saídas.
- Metas/precificação usarão custo médio/estoque mínimo; precisam dos saldos corretos.

## Referências
- `Tarefas/ESTOQUE/01-entrada.md`
- `Tarefas/ESTOQUE/02-saida.md`
- `Tarefas/ESTOQUE/03-consumo-automatico.md`
- `Tarefas/ESTOQUE/04-inventario.md`
- `Tarefas/ESTOQUE/05-curva-abc.md`
- `Tarefas/ESTOQUE/06-estoque-minimo.md`
